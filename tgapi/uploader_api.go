package tgapi

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/utils"
	"git.scuroneko.dev/scuroneko/sneklog/v2"
)

const (
	// UploaderPhotoType is the multipart field name for photo uploads.
	UploaderPhotoType UploaderFileType = "photo"
	// UploaderVideoType is the multipart field name for video uploads.
	UploaderVideoType UploaderFileType = "video"
	// UploaderAudioType is the multipart field name for audio uploads.
	UploaderAudioType UploaderFileType = "audio"
	// UploaderDocumentType is the multipart field name for document uploads.
	UploaderDocumentType UploaderFileType = "document"
	// UploaderVoiceType is the multipart field name for voice uploads.
	UploaderVoiceType UploaderFileType = "voice"
	// UploaderVideoNoteType is the multipart field name for video-note uploads.
	UploaderVideoNoteType UploaderFileType = "video_note"
	// UploaderThumbnailType is the multipart field name for thumbnail uploads.
	UploaderThumbnailType UploaderFileType = "thumbnail"
	// UploaderStickerType is the multipart field name for sticker uploads.
	UploaderStickerType UploaderFileType = "sticker"
	// UploaderCertificateType is the multipart field name for webhook certificate uploads.
	UploaderCertificateType UploaderFileType = "certificate"
)

// UploaderFileType represents the Telegram form field name for a file upload.
type UploaderFileType string

// UploaderFile holds the data and metadata for a single file to be uploaded.
type UploaderFile struct {
	filename string
	data     []byte
	field    UploaderFileType
}

// NewUploaderFile creates a new UploaderFile, auto-detecting the field type from the file extension.
// If detection is incorrect, use SetType to override.
func NewUploaderFile(name string, data []byte) UploaderFile {
	t := uploaderTypeByExt(name)
	return UploaderFile{filename: name, data: data, field: t}
}

// SetType overrides the auto-detected upload field type.
// For example, use it when a voice file is detected as audio.
func (f UploaderFile) SetType(t UploaderFileType) UploaderFile {
	f.field = t
	return f
}

// Uploader is a Telegram Bot API client specialized for multipart file uploads.
//
// Use Uploader methods when you need to upload binary files directly
// (InputFile/multipart). For JSON-only calls (file_id, URL, plain params), use API.
type Uploader struct {
	api    *API
	logger *sneklog.Logger
}

// NewUploader creates a multipart uploader bound to an API client.
func NewUploader(api *API) *Uploader {
	if api == nil {
		return nil
	}
	logger := utils.CreateLogger(
		"UPLOADER", utils.GetLoggerLevel(),
		api.logFormat, api.logFormatter,
	)
	return &Uploader{api, logger}
}

// Close flushes and closes uploader logger resources.
// See https://core.telegram.org/bots/api
func (u *Uploader) Close() error { return u.logger.Close() }

// GetLogger returns uploader logger instance.
// See https://core.telegram.org/bots/api
func (u *Uploader) GetLogger() *sneklog.Logger { return u.logger }

// UploaderRequest is a low-level multipart upload request wrapper.
//
// Prefer method-specific helpers such as SendPhoto or SetWebhook. UploaderRequest
// is intended for advanced use cases where callers manage the method name, files,
// and request/response types themselves. In that sense it is an unsafe escape
// hatch compared with the typed uploader API.
type UploaderRequest[R, P any] struct {
	method string
	files  []UploaderFile
	params P
	chatID int64
}

// NewUploaderRequest creates a low-level multipart upload request with no associated chat ID.
func NewUploaderRequest[R, P any](method string, params P, files ...UploaderFile) UploaderRequest[R, P] {
	return UploaderRequest[R, P]{method: method, files: files, params: params, chatID: 0}
}

// NewUploaderRequestWithChatID creates a low-level multipart upload request with an associated chat ID.
// The chat ID is used for per-chat rate limiting.
func NewUploaderRequestWithChatID[R, P any](method string, params P, chatID int64, files ...UploaderFile) UploaderRequest[R, P] {
	return UploaderRequest[R, P]{method: method, files: files, params: params, chatID: chatID}
}

func (r UploaderRequest[R, P]) doRequest(ctx context.Context, up *Uploader) (R, error) {
	var zero R

	methodPrefix := ""
	if up.api.useTestServer {
		methodPrefix = "/test"
	}
	url := fmt.Sprintf("%s/bot%s%s/%s", up.api.apiURL, up.api.token, methodPrefix, r.method)

	for {
		if up.api.Limiter != nil {
			if err := up.api.Limiter.Check(ctx, up.api.dropOverflowLimit, r.chatID); err != nil {
				return zero, err
			}
		}

		buf, contentType, err := prepareMultipart(r.files, r.params)
		if err != nil {
			return zero, err
		}
		req, err := http.NewRequestWithContext(ctx, "POST", url, buf)
		if err != nil {
			return zero, err
		}
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))
		req.ContentLength = int64(buf.Len())

		up.logger.Debugln("UPLOADER REQ", r.method)
		resp, err := up.api.client.Do(req)
		if err != nil {
			return zero, err
		}

		body, err := readBody(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return zero, err
		}
		up.logger.Debugln("UPLOADER RES", r.method, string(body))

		response, err := parseBody[R](body)
		if err != nil {
			return zero, err
		}

		if !response.Ok {
			if response.ErrorCode == 429 && response.Parameters != nil && response.Parameters.RetryAfter != nil {
				after := *response.Parameters.RetryAfter
				up.logger.Warnf("Rate limited, retry after %d seconds (chat: %d)", after, r.chatID)
				if up.api.Limiter != nil {
					if r.chatID > 0 {
						up.api.Limiter.SetChatLock(r.chatID, after)
					} else {
						up.api.Limiter.SetGlobalLock(after)
					}
				}

				select {
				case <-ctx.Done():
					return zero, ctx.Err()
				case <-time.After(time.Duration(after) * time.Second):
					continue // Повторяем запрос
				}
			}
			return zero, fmt.Errorf("[%d] %s", response.ErrorCode, response.Description)
		}
		return response.Result, nil
	}
}

// DoWithContext executes the upload request asynchronously via the worker pool.
// Returns the result or error. Respects context cancellation.
func (r UploaderRequest[R, P]) DoWithContext(ctx context.Context, up *Uploader) (R, error) {
	var zero R

	result, err := up.api.pool.submit(ctx, func(ctx context.Context) (any, error) {
		return r.doRequest(ctx, up)
	})
	if err != nil {
		return zero, err
	}

	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case res := <-result:
		if res.err != nil {
			return zero, res.err
		}
		if val, ok := res.value.(R); ok {
			return val, nil
		}
		return zero, ErrPoolUnexpected
	}
}

// Do executes the upload request synchronously with a background context.
// Use only for simple, non-critical uploads.
func (r UploaderRequest[R, P]) Do(up *Uploader) (R, error) {
	return r.DoWithContext(context.Background(), up)
}

// Internal helper that builds a finalized multipart body from files and params.
func prepareMultipart[P any](files []UploaderFile, params P) (*bytes.Buffer, string, error) {
	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)

	for _, file := range files {
		fw, err := w.CreateFormFile(string(file.field), file.filename)
		if err != nil {
			_ = w.Close() // Закрываем, чтобы не было утечки
			return nil, "", err
		}

		_, err = fw.Write(file.data)
		if err != nil {
			_ = w.Close()
			return nil, "", err
		}
	}

	err := utils.Encode(w, params) // Предполагается, что это записывает в w
	if err != nil {
		_ = w.Close()
		return nil, "", err
	}

	err = w.Close() // ✅ ОБЯЗАТЕЛЬНО вызвать в конце — иначе запрос битый!
	if err != nil {
		return nil, "", err
	}

	return buf, w.FormDataContentType(), nil
}

// Internal helper that infers an upload field name from a file extension.
func uploaderTypeByExt(filename string) UploaderFileType {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".bmp":
		return UploaderPhotoType
	case ".mp4":
		return UploaderVideoType
	case ".mp3", ".m4a":
		return UploaderAudioType
	case ".ogg":
		return UploaderVoiceType
	default:
		return UploaderDocumentType
	}
}
