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
	UploaderThumbnailType UploaderFileType = "thumbnail"
)

type UploaderFile struct {
	filename string
	data     []byte
	field    UploaderFileType
}

func NewUploaderFile(name string, data []byte) UploaderFile {
	t := uploaderTypeByExt(name)
	return UploaderFile{filename: name, data: data, field: t}
}

// SetType used when auto-detect failed. I.e. you sending a voice message, but it detects as audio, or if you send audio with thumbnail
func (f UploaderFile) SetType(t UploaderFileType) UploaderFile {
	f.field = t
	return f
}

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
			w.Close()
			return zero, err
		}
		_, err = fw.Write(file.data)
		if err != nil {
			w.Close()
			return zero, err
		}
	}

	err := utils.Encode(w, u.params)
	if err != nil {
		w.Close()
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
	req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))

	up.logger.Debugln("UPLOADER REQ", u.method)
	res, err := up.api.client.Do(req)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return zero, err
	}
	up.logger.Debugln("UPLOADER RES", u.method, string(body))
	if res.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("[%d] %s", res.StatusCode, string(body))
	}

	var response ApiResponse[R]
	err = json.Unmarshal(body, &response)
	if err != nil {
		return zero, err
	}
	if !response.Ok {
		return zero, fmt.Errorf("[%d] %s", response.ErrorCode, response.Description)
	}
	return response.Result, nil
}
func (u UploaderRequest[R, P]) Do(up *Uploader) (R, error) {
	ctx := context.Background()
	return u.DoWithContext(ctx, up)
}

type UploadPhotoP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int    `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	ShowCaptionAboveMedia bool   `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool   `json:"has_spoiler,omitempty"`
	DisableNotification   bool   `json:"disable_notification,omitempty"`
	ProtectContent        bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (u *Uploader) UploadPhoto(params UploadPhotoP, file UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendPhoto", params, file)
	return req.Do(u)
}

type UploadAudioP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int    `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	Duration  int    `json:"duration,omitempty"`
	Performer string `json:"performer,omitempty"`
	Title     string `json:"title,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (u *Uploader) UploadAudio(params UploadAudioP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendAudio", params, files...)
	return req.Do(u)
}

type UploadDocumentP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int    `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	DisableContentTypeDetection bool   `json:"disable_content_type_detection,omitempty"`
	DisableNotification         bool   `json:"disable_notification,omitempty"`
	ProtectContent              bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast          bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID             string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (u *Uploader) UploadDocument(params UploadDocumentP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendDocument", params, files...)
	return req.Do(u)
}

type UploadVideoP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int    `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Duration int `json:"duration,omitempty"`
	Width    int `json:"width,omitempty"`
	Height   int `json:"height,omitempty"`

	StartTimestamp  int64           `json:"start_timestamp,omitempty"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	ShowCaptionAboveMedia bool   `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool   `json:"has_spoiler,omitempty"`
	SupportsStreaming     bool   `json:"supports_streaming,omitempty"`
	DisableNotification   bool   `json:"disable_notification,omitempty"`
	ProtectContent        bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (u *Uploader) UploadVideo(params UploadVideoP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendVideo", params, files...)
	return req.Do(u)
}

type UploadAnimationP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int    `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Duration int `json:"duration,omitempty"`
	Width    int `json:"width,omitempty"`
	Height   int `json:"height,omitempty"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	ShowCaptionAboveMedia bool   `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool   `json:"has_spoiler,omitempty"`
	DisableNotification   bool   `json:"disable_notification,omitempty"`
	ProtectContent        bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (u *Uploader) UploadAnimation(params UploadAnimationP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendAnimation", params, files...)
	return req.Do(u)
}

type UploadVoiceP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int    `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Duration        int             `json:"duration,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (u *Uploader) UploadVoice(params UploadVoiceP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendVoice", params, files...)
	return req.Do(u)
}

type UploadVideoNoteP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int    `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Duration int `json:"duration,omitempty"`
	Length   int `json:"length,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

func (u *Uploader) UploadVideoNote(params UploadVideoNoteP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendVideoNote", params, files...)
	return req.Do(u)
}

// setChatPhoto https://core.telegram.org/bots/api#setchatphoto
type UploadChatPhotoP struct {
	ChatID int `json:"chat_id"`
}

func (u *Uploader) UploadChatPhoto(params UploadChatPhotoP, photo UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendChatPhoto", params, photo)
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
