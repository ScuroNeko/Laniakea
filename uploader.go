package laniakea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"git.nix13.pw/scuroneko/slog"
)

type Uploader struct {
	api    *Api
	logger *slog.Logger
}

func NewUploader(api *Api) *Uploader {
	logger := slog.CreateLogger().Level(GetLoggerLevel()).Prefix("UPLOADER")
	logger.AddWriter(logger.CreateJsonStdoutWriter())
	return &Uploader{api, logger}
}
func (u *Uploader) Close() {
	u.logger.Close()
}

type UploaderFileType string

const (
	UploaderPhotoType     UploaderFileType = "photo"
	UploaderVideoType     UploaderFileType = "video"
	UploaderAudioType     UploaderFileType = "audio"
	UploaderDocumentType  UploaderFileType = "document"
	UploaderVoiceType     UploaderFileType = "voice"
	UploaderVideoNoteType UploaderFileType = "video_note"
)

type UploaderFile struct {
	filename string
	data     []byte
	t        UploaderFileType
}

func NewUploaderFile(name string, data []byte) UploaderFile {
	t := uploaderTypeByExt(name)
	return UploaderFile{filename: name, data: data, t: t}
}

// SetType used when auto-detect failed. I.e. you sending a voice message, but it detects as audio
func (f UploaderFile) SetType(t UploaderFileType) UploaderFile {
	f.t = t
	return f
}

type UploaderRequest[R, P any] struct {
	method string
	file   UploaderFile
	params P
}

func NewUploaderRequest[R, P any](method string, file UploaderFile, params P) UploaderRequest[R, P] {
	return UploaderRequest[R, P]{method, file, params}
}

func (u UploaderRequest[R, P]) Do(up *Uploader) (*R, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", up.api.token, u.method)

	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)

	fw, err := w.CreateFormFile(string(u.file.t), u.file.filename)
	if err != nil {
		w.Close()
		return nil, err
	}
	_, err = fw.Write(u.file.data)
	if err != nil {
		w.Close()
		return nil, err
	}
	err = Encode(w, u.params)
	if err != nil {
		w.Close()
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	up.logger.Debugln("UPLOADER REQ", u.method)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	up.logger.Debugln("UPLOADER RES", u.method, string(body))
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[%d] %s", res.StatusCode, string(body))
	}

	response := new(ApiResponse[*R])
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, fmt.Errorf("[%d] %s", response.ErrorCode, response.Description)
	}
	return response.Result, nil
}

func (u *Uploader) UploadPhoto(file UploaderFile, params SendPhotoBaseP) (*Message, error) {
	req := NewUploaderRequest[Message]("sendPhoto", file, params)
	return req.Do(u)
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
