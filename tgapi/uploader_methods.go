package tgapi

// UploadPhotoP holds parameters for uploading a photo using the Uploader.
// See https://core.telegram.org/bots/api#sendphoto
type UploadPhotoP struct {
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
// file is the photo file to upload.
// See https://core.telegram.org/bots/api#sendphoto
func (u *Uploader) SendPhoto(params UploadPhotoP, file UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendPhoto", params, params.ChatID, file)
	return req.Do(u)
}

// UploadAudioP holds parameters for uploading an audio file using the Uploader.
// See https://core.telegram.org/bots/api#sendaudio
type UploadAudioP struct {
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
// files are the audio file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendaudio
func (u *Uploader) SendAudio(params UploadAudioP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendAudio", params, params.ChatID, files...)
	return req.Do(u)
}

// UploadDocumentP holds parameters for uploading a document using the Uploader.
// See https://core.telegram.org/bots/api#senddocument
type UploadDocumentP struct {
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
// files are the document file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#senddocument
func (u *Uploader) SendDocument(params UploadDocumentP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendDocument", params, params.ChatID, files...)
	return req.Do(u)
}

// UploadVideoP holds parameters for uploading a video using the Uploader.
// See https://core.telegram.org/bots/api#sendvideo
type UploadVideoP struct {
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
// files are the video file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendvideo
func (u *Uploader) SendVideo(params UploadVideoP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendVideo", params, params.ChatID, files...)
	return req.Do(u)
}

// UploadAnimationP holds parameters for uploading an animation using the Uploader.
// See https://core.telegram.org/bots/api#sendanimation
type UploadAnimationP struct {
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
// files are the animation file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendanimation
func (u *Uploader) SendAnimation(params UploadAnimationP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendAnimation", params, params.ChatID, files...)
	return req.Do(u)
}

// UploadVoiceP holds parameters for uploading a voice note using the Uploader.
// See https://core.telegram.org/bots/api#sendvoice
type UploadVoiceP struct {
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
// files are the voice file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendvoice
func (u *Uploader) SendVoice(params UploadVoiceP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendVoice", params, params.ChatID, files...)
	return req.Do(u)
}

// UploadVideoNoteP holds parameters for uploading a video note (rounded video) using the Uploader.
// See https://core.telegram.org/bots/api#sendvideonote
type UploadVideoNoteP struct {
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
// files are the video note file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendvideonote
func (u *Uploader) SendVideoNote(params UploadVideoNoteP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequestWithChatID[Message]("sendVideoNote", params, params.ChatID, files...)
	return req.Do(u)
}

// UploadChatPhotoP holds parameters for uploading a chat photo using the Uploader.
// See https://core.telegram.org/bots/api#setchatphoto
type UploadChatPhotoP struct {
	ChatID int64 `json:"chat_id"`
}

// SetChatPhoto uploads a new chat photo.
// photo is the photo file to upload.
// See https://core.telegram.org/bots/api#setchatphoto
func (u *Uploader) SetChatPhoto(params UploadChatPhotoP, photo UploaderFile) (bool, error) {
	req := NewUploaderRequestWithChatID[bool]("setChatPhoto", params, params.ChatID, photo)
	return req.Do(u)
}
