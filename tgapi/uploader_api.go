package tgapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

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
func (u *Uploader) Close() error { return u.logger.Close() }

type UploaderRequest[R, P any] struct {
	method string
	files  []UploaderFile
	params P
}

func NewUploaderRequest[R, P any](method string, params P, files ...UploaderFile) UploaderRequest[R, P] {
	return UploaderRequest[R, P]{method, files, params}
}
func (u UploaderRequest[R, P]) DoWithContext(ctx context.Context, up *Uploader) (R, error) {
	var zero R
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", up.api.token, u.method)

	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)

	for _, file := range u.files {
		fw, err := w.CreateFormFile(string(file.field), file.filename)
		if err != nil {
			_ = w.Close()
			return zero, err
		}
		_, err = fw.Write(file.data)
		if err != nil {
			_ = w.Close()
			return zero, err
		}
	}

	err := utils.Encode(w, u.params)
	if err != nil {
		_ = w.Close()
		return zero, err
	}
	err = w.Close()
	if err != nil {
		return zero, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, buf)
	if err != nil {
		return zero, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))

	up.logger.Debugln("UPLOADER REQ", u.method)
	res, err := up.api.client.Do(req)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()

	reader := io.LimitReader(res.Body, 10<<20)
	body, err := io.ReadAll(reader)
	if err != nil {
		return zero, err
	}
	up.logger.Debugln("UPLOADER RES", u.method, string(body))
	if res.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("unexpected status code: %d, %s", res.StatusCode, string(body))
	}

	var resp ApiResponse[R]
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return zero, err
	}
	if !resp.Ok {
		return zero, fmt.Errorf("[%d] %s", resp.ErrorCode, resp.Description)
	}
	return resp.Result, nil
}
func (u UploaderRequest[R, P]) Do(up *Uploader) (R, error) {
	return u.DoWithContext(context.Background(), up)
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
