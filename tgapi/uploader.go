package tgapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"git.nix13.pw/scuroneko/laniakea/utils"
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
func (u *Uploader) Close() error {
	return u.logger.Close()
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
	err = utils.Encode(w, u.params)
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

type UploadPhotoP struct {
	BusinessConnectionID  string                `json:"business_connection_id,omitempty"`
	ChatID                int                   `json:"chat_id"`
	MessageThreadID       int                   `json:"message_thread_id,omitempty"`
	ParseMode             ParseMode             `json:"parse_mode,omitempty"`
	Caption               string                `json:"caption,omitempty"`
	CaptionEntities       []*MessageEntity      `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                  `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool                  `json:"has_spoiler,omitempty"`
	DisableNotifications  bool                  `json:"disable_notifications,omitempty"`
	ProtectContent        bool                  `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool                  `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string                `json:"message_effect_id,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (u *Uploader) UploadPhoto(file UploaderFile, params UploadPhotoP) (*Message, error) {
	req := NewUploaderRequest[Message]("sendPhoto", file, params)
	return req.Do(u)
}

// setChatPhoto https://core.telegram.org/bots/api#setchatphoto
type UploadChatPhotoP struct {
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
