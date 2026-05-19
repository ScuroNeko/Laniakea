package tgapi

import "context"

// UploadPhoto holds parameters for uploading a photo using the Uploader.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendphoto
type UploadPhoto struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
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

// SendPhoto uploads a photo via multipart and sends it as a message.
// Since: Bot API 1.0
// file is the photo file to upload.
// See https://core.telegram.org/bots/api#sendphoto
func (u *Uploader) SendPhoto(params UploadPhoto, file UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendPhoto", params, params.ChatID, file)
	return req.Do(u)
}

// SendPhotoWithContext is the context-aware variant of SendPhoto.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendphoto
func (u *Uploader) SendPhotoWithContext(ctx context.Context, params UploadPhoto, file UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendPhoto", params, params.ChatID, file)
	return req.DoWithContext(ctx, u)
}

// UploadAudio holds parameters for uploading an audio file using the Uploader.
// Since: Bot API 1.2
// See https://core.telegram.org/bots/api#sendaudio
type UploadAudio struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
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

// SendAudio uploads an audio file via multipart and sends it as a message.
// Since: Bot API 1.2
// files are the audio file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendaudio
func (u *Uploader) SendAudio(params UploadAudio, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendAudio", params, params.ChatID, files...)
	return req.Do(u)
}

// SendAudioWithContext is the context-aware variant of SendAudio.
// Since: Bot API 1.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendaudio
func (u *Uploader) SendAudioWithContext(ctx context.Context, params UploadAudio, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendAudio", params, params.ChatID, files...)
	return req.DoWithContext(ctx, u)
}

// UploadDocument holds parameters for uploading a document using the Uploader.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#senddocument
type UploadDocument struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
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

// SendDocument uploads a document via multipart and sends it as a message.
// Since: Bot API 1.0
// files are the document file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#senddocument
func (u *Uploader) SendDocument(params UploadDocument, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendDocument", params, params.ChatID, files...)
	return req.Do(u)
}

// SendDocumentWithContext is the context-aware variant of SendDocument.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#senddocument
func (u *Uploader) SendDocumentWithContext(ctx context.Context, params UploadDocument, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendDocument", params, params.ChatID, files...)
	return req.DoWithContext(ctx, u)
}

// UploadVideo holds parameters for uploading a video using the Uploader.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendvideo
type UploadVideo struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
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

// SendVideo uploads a video via multipart and sends it as a message.
// Since: Bot API 1.0
// files are the video file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendvideo
func (u *Uploader) SendVideo(params UploadVideo, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendVideo", params, params.ChatID, files...)
	return req.Do(u)
}

// SendVideoWithContext is the context-aware variant of SendVideo.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvideo
func (u *Uploader) SendVideoWithContext(ctx context.Context, params UploadVideo, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendVideo", params, params.ChatID, files...)
	return req.DoWithContext(ctx, u)
}

// UploadAnimation holds parameters for uploading an animation using the Uploader.
// Since: Bot API 4.0
// See https://core.telegram.org/bots/api#sendanimation
type UploadAnimation struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
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

// SendAnimation uploads an animation via multipart and sends it as a message.
// Since: Bot API 4.0
// files are the animation file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendanimation
func (u *Uploader) SendAnimation(params UploadAnimation, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendAnimation", params, params.ChatID, files...)
	return req.Do(u)
}

// SendAnimationWithContext is the context-aware variant of SendAnimation.
// Since: Bot API 4.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendanimation
func (u *Uploader) SendAnimationWithContext(ctx context.Context, params UploadAnimation, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendAnimation", params, params.ChatID, files...)
	return req.DoWithContext(ctx, u)
}

// UploadVoice holds parameters for uploading a voice note using the Uploader.
// Since: Bot API 1.2
// See https://core.telegram.org/bots/api#sendvoice
type UploadVoice struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
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

// SendVoice uploads a voice note via multipart and sends it as a message.
// Since: Bot API 1.2
// files are the voice file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendvoice
func (u *Uploader) SendVoice(params UploadVoice, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendVoice", params, params.ChatID, files...)
	return req.Do(u)
}

// SendVoiceWithContext is the context-aware variant of SendVoice.
// Since: Bot API 1.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvoice
func (u *Uploader) SendVoiceWithContext(ctx context.Context, params UploadVoice, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendVoice", params, params.ChatID, files...)
	return req.DoWithContext(ctx, u)
}

// UploadVideoNote holds parameters for uploading a video note (rounded video) using the Uploader.
// Since: Bot API 3.0
// See https://core.telegram.org/bots/api#sendvideonote
type UploadVideoNote struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
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

// SendVideoNote uploads a video note via multipart and sends it as a message.
// Since: Bot API 3.0
// files are the video note file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendvideonote
func (u *Uploader) SendVideoNote(params UploadVideoNote, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendVideoNote", params, params.ChatID, files...)
	return req.Do(u)
}

// SendVideoNoteWithContext is the context-aware variant of SendVideoNote.
// Since: Bot API 3.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvideonote
func (u *Uploader) SendVideoNoteWithContext(ctx context.Context, params UploadVideoNote, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendVideoNote", params, params.ChatID, files...)
	return req.DoWithContext(ctx, u)
}

// UploadChatPhoto holds parameters for uploading a chat photo using the Uploader.
// Since: Bot API 3.1
// See https://core.telegram.org/bots/api#setchatphoto
type UploadChatPhoto struct {
	ChatID int64 `json:"chat_id"`
}

// SetChatPhoto uploads a new chat photo.
// Since: Bot API 3.1
// photo is the photo file to upload.
// See https://core.telegram.org/bots/api#setchatphoto
func (u *Uploader) SetChatPhoto(params UploadChatPhoto, photo UploaderFile) (bool, error) {
	req := NewUploaderRequestWithChatID[bool]("setChatPhoto", params, params.ChatID, photo)
	return req.Do(u)
}

// SetChatPhotoWithContext is the context-aware variant of SetChatPhoto.
// Since: Bot API 3.1
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setchatphoto
func (u *Uploader) SetChatPhotoWithContext(ctx context.Context, params UploadChatPhoto, photo UploaderFile) (bool, error) {
	req := NewUploaderRequestWithChatID[bool]("setChatPhoto", params, params.ChatID, photo)
	return req.DoWithContext(ctx, u)
}

// UploadSetWebhook holds multipart parameters for the setWebhook method.
// Since: Bot API 1.0
// Use this type when uploading a self-signed certificate file.
// See https://core.telegram.org/bots/api#setwebhook
type UploadSetWebhook struct {
	URL                string       `json:"url"`
	IPAddress          string       `json:"ip_address,omitempty"`
	MaxConnections     int8         `json:"max_connections,omitempty"`
	AllowedUpdates     []UpdateType `json:"allowed_updates,omitempty"`
	DropPendingUpdates bool         `json:"drop_pending_updates,omitempty"`
	SecretToken        string       `json:"secret_token,omitempty"`
}

// SetWebhook uploads a certificate and sets a webhook URL.
// Since: Bot API 1.0
// certificate maps to the multipart field "certificate".
// See https://core.telegram.org/bots/api#setwebhook
func (u *Uploader) SetWebhook(params UploadSetWebhook, certificate UploaderFile) (bool, error) {
	req := NewUploaderRequest[bool]("setWebhook", params, certificate.SetType(UploaderCertificateType))
	return req.Do(u)
}

// SetWebhookWithContext is the context-aware variant of SetWebhook.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setwebhook
func (u *Uploader) SetWebhookWithContext(ctx context.Context, params UploadSetWebhook, certificate UploaderFile) (bool, error) {
	req := NewUploaderRequest[bool]("setWebhook", params, certificate.SetType(UploaderCertificateType))
	return req.DoWithContext(ctx, u)
}

// UploadLivePhoto holds parameters for uploading a live photo using the Uploader.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#sendlivephoto
type UploadLivePhoto struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
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

// SendLivePhoto uploads a live photo via multipart and sends it as a message.
// Since: Bot API 10.0
// file is the live photo file to upload.
// See https://core.telegram.org/bots/api#sendlivephoto
func (u *Uploader) SendLivePhoto(params UploadLivePhoto, file UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendLivePhoto", params, params.ChatID, file.SetType(UploaderLivePhotoType))
	return req.Do(u)
}

// SendLivePhotoWithContext is the context-aware variant of SendLivePhoto.
// Since: Bot API 10.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendlivephoto
func (u *Uploader) SendLivePhotoWithContext(ctx context.Context, params UploadLivePhoto, file UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendLivePhoto", params, params.ChatID, file.SetType(UploaderLivePhotoType))
	return req.DoWithContext(ctx, u)
}
