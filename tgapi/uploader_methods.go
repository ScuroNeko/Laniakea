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

// UploadPhoto uploads a photo and sends it as a message.
// file is the photo file to upload.
// See https://core.telegram.org/bots/api#sendphoto
func (u *Uploader) UploadPhoto(params UploadPhotoP, file UploaderFile) (Message, error) {
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

// UploadAudio uploads an audio file and sends it as a message.
// files are the audio file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendaudio
func (u *Uploader) UploadAudio(params UploadAudioP, files ...UploaderFile) (Message, error) {
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

// UploadDocument uploads a document and sends it as a message.
// files are the document file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#senddocument
func (u *Uploader) UploadDocument(params UploadDocumentP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendDocument", params, files...)
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

// UploadVideo uploads a video and sends it as a message.
// files are the video file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendvideo
func (u *Uploader) UploadVideo(params UploadVideoP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendVideo", params, files...)
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

// UploadAnimation uploads an animation (GIF or H.264/MPEG-4 AVC video without sound) and sends it as a message.
// files are the animation file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendanimation
func (u *Uploader) UploadAnimation(params UploadAnimationP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendAnimation", params, files...)
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

// UploadVoice uploads a voice note and sends it as a message.
// files are the voice file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendvoice
func (u *Uploader) UploadVoice(params UploadVoiceP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendVoice", params, files...)
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

// UploadVideoNote uploads a video note (rounded video) and sends it as a message.
// files are the video note file(s) to upload (typically one file).
// See https://core.telegram.org/bots/api#sendvideonote
func (u *Uploader) UploadVideoNote(params UploadVideoNoteP, files ...UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendVideoNote", params, files...)
	return req.Do(u)
}

// UploadChatPhotoP holds parameters for uploading a chat photo using the Uploader.
// See https://core.telegram.org/bots/api#setchatphoto
type UploadChatPhotoP struct {
	ChatID int64 `json:"chat_id"`
}

// UploadChatPhoto uploads a new chat photo.
// photo is the photo file to upload.
// See https://core.telegram.org/bots/api#setchatphoto
func (u *Uploader) UploadChatPhoto(params UploadChatPhotoP, photo UploaderFile) (Message, error) {
	req := NewUploaderRequest[Message]("sendChatPhoto", params, photo)
	return req.Do(u)
}
