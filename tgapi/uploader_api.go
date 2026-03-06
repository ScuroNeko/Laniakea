package tgapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"git.nix13.pw/scuroneko/laniakea/utils"
	"git.nix13.pw/scuroneko/slog"
)

const (
	UploaderPhotoType     UploaderFileType = "photo"
	UploaderVideoType     UploaderFileType = "video"
	UploaderAudioType     UploaderFileType = "audio"
	UploaderDocumentType  UploaderFileType = "document"
	UploaderVoiceType     UploaderFileType = "voice"
	UploaderVideoNoteType UploaderFileType = "video_note"
	UploaderThumbnailType UploaderFileType = "thumbnail"
)

type UploaderFileType string
type UploaderFile struct {
	filename string
	data     []byte
	field    UploaderFileType
}

func NewUploaderFile(name string, data []byte) UploaderFile {
	t := uploaderTypeByExt(name)
	return UploaderFile{filename: name, data: data, field: t}
}

// SetType used when auto-detect failed.
// i.e. you sending a voice message, but it detects as audio, or if you send audio with thumbnail
func (f UploaderFile) SetType(t UploaderFileType) UploaderFile {
	f.field = t
	return f
}

type Uploader struct {
	api    *API
	logger *slog.Logger
}

func NewUploader(api *API) *Uploader {
	logger := slog.CreateLogger().Level(utils.GetLoggerLevel()).Prefix("UPLOADER")
	logger.AddWriter(logger.CreateJsonStdoutWriter())
	return &Uploader{api, logger}
}
func (u *Uploader) Close() error            { return u.logger.Close() }
func (u *Uploader) GetLogger() *slog.Logger { return u.logger }

type UploaderRequest[R, P any] struct {
	method string
	files  []UploaderFile
	params P
	chatId int64
}

func NewUploaderRequest[R, P any](method string, params P, files ...UploaderFile) UploaderRequest[R, P] {
	return UploaderRequest[R, P]{method: method, files: files, params: params, chatId: 0}
}
func NewUploaderRequestWithChatID[R, P any](method string, params P, chatId int64, files ...UploaderFile) UploaderRequest[R, P] {
	return UploaderRequest[R, P]{method: method, files: files, params: params, chatId: chatId}
}
func (r UploaderRequest[R, P]) doRequest(ctx context.Context, up *Uploader) (R, error) {
	var zero R

	buf, contentType, err := prepareMultipart(r.files, r.params)
	if err != nil {
		return zero, err
	}

	methodPrefix := ""
	if up.api.useTestServer {
		methodPrefix = "/test"
	}
	url := fmt.Sprintf("%s/bot%s%s/%s", up.api.apiUrl, up.api.token, methodPrefix, r.method)
	req, err := http.NewRequestWithContext(ctx, "POST", url, buf)
	if err != nil {
		return zero, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))
	req.Header.Set("Accept-Encoding", "gzip")
	req.ContentLength = int64(buf.Len())

	for {
		if up.api.Limiter != nil {
			if up.api.dropOverflowLimit {
				if !up.api.Limiter.GlobalAllow() {
					return zero, errors.New("rate limited")
				}
			} else {
				if err := up.api.Limiter.GlobalWait(ctx); err != nil {
					return zero, err
				}
			}
		}

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
				if r.chatId > 0 {
					up.api.Limiter.SetChatLock(r.chatId, after)
				} else {
					up.api.Limiter.SetGlobalLock(after)
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
func (r UploaderRequest[R, P]) Do(up *Uploader) (R, error) {
	return r.DoWithContext(context.Background(), up)
}

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
