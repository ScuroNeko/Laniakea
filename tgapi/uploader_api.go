package tgapi

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"git.nix13.pw/scuroneko/laniakea/utils"
	"git.nix13.pw/scuroneko/slog"
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
	logger *slog.Logger
}

// NewUploader creates a multipart uploader bound to an API client.
func NewUploader(api *API) *Uploader {
	logger := slog.CreateLogger().Level(utils.GetLoggerLevel()).Prefix("UPLOADER")
	logger.AddWriter(logger.CreateJsonStdoutWriter())
	return &Uploader{api, logger}
}

// Close flushes and closes uploader logger resources.
// See https://core.telegram.org/bots/api
func (u *Uploader) Close() error { return u.logger.Close() }

// GetLogger returns uploader logger instance.
// See https://core.telegram.org/bots/api
func (u *Uploader) GetLogger() *slog.Logger { return u.logger }

// UploaderRequest is a multipart file upload request to the Telegram API.
// Use NewUploaderRequest or NewUploaderRequestWithChatID to construct one.
type UploaderRequest[R, P any] struct {
	method string
	files  []UploaderFile
	params P
	chatId int64
}

// NewUploaderRequest creates a new multipart upload request with no associated chat ID.
func NewUploaderRequest[R, P any](method string, params P, files ...UploaderFile) UploaderRequest[R, P] {
	return UploaderRequest[R, P]{method: method, files: files, params: params, chatId: 0}
}

// NewUploaderRequestWithChatID creates a new multipart upload request with an associated chat ID.
// The chat ID is used for per-chat rate limiting.
func NewUploaderRequestWithChatID[R, P any](method string, params P, chatId int64, files ...UploaderFile) UploaderRequest[R, P] {
	return UploaderRequest[R, P]{method: method, files: files, params: params, chatId: chatId}
}
func (r UploaderRequest[R, P]) doRequest(ctx context.Context, up *Uploader) (R, error) {
	var zero R

	methodPrefix := ""
	if up.api.useTestServer {
		methodPrefix = "/test"
	}
	url := fmt.Sprintf("%s/bot%s%s/%s", up.api.apiUrl, up.api.token, methodPrefix, r.method)

	for {
		if up.api.Limiter != nil {
			if err := up.api.Limiter.Check(ctx, up.api.dropOverflowLimit, r.chatId); err != nil {
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
		req.Header.Set("Accept-Encoding", "gzip")
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
				up.logger.Warnf("Rate limited, retry after %d seconds (chat: %d)", after, r.chatId)
				if up.api.Limiter != nil {
					if r.chatId > 0 {
						up.api.Limiter.SetChatLock(r.chatId, after)
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

// prepareMultipart builds a multipart form body from the given files and params.
// Params are encoded via utils.Encode. The writer boundary is finalized before returning.
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

// uploaderTypeByExt infers the Telegram upload field name from a file extension.
// Falls back to UploaderDocumentType for unrecognized extensions.
func uploaderTypeByExt(filename string) UploaderFileType {
	ext := filepath.Ext(filename)
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
