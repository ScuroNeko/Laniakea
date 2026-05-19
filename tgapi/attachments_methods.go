package tgapi

import "context"

// SendPhoto holds parameters for the sendPhoto method.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendphoto
type SendPhoto struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Photo           string          `json:"photo"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	ShowCaptionAboveMedia bool   `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool   `json:"has_spoiler,omitempty"`
	DisableNotifications  bool   `json:"disable_notification,omitempty"`
	ProtectContent        bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

// SendPhoto sends a photo.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendphoto
func (api *API) SendPhoto(params SendPhoto) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPhoto", params, params.ChatID)
	return req.Do(api)
}

// SendPhotoWithContext is the context-aware variant of SendPhoto.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendphoto
func (api *API) SendPhotoWithContext(ctx context.Context, params SendPhoto) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPhoto", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendAudio holds parameters for the sendAudio method.
// Since: Bot API 1.2
// See https://core.telegram.org/bots/api#sendaudio
type SendAudio struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Audio           string          `json:"audio"`
	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Duration        int             `json:"duration,omitempty"`
	Performer       string          `json:"performer,omitempty"`
	Title           string          `json:"title,omitempty"`
	Thumbnail       string          `json:"thumbnail,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

// SendAudio sends an audio file.
// Since: Bot API 1.2
// See https://core.telegram.org/bots/api#sendaudio
func (api *API) SendAudio(params SendAudio) (Message, error) {
	req := NewRequestWithChatID[Message]("sendAudio", params, params.ChatID)
	return req.Do(api)
}

// SendAudioWithContext is the context-aware variant of SendAudio.
// Since: Bot API 1.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendaudio
func (api *API) SendAudioWithContext(ctx context.Context, params SendAudio) (Message, error) {
	req := NewRequestWithChatID[Message]("sendAudio", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendDocument holds parameters for the sendDocument method.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#senddocument
type SendDocument struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Document                    string          `json:"document"`
	Thumbnail                   string          `json:"thumbnail,omitempty"`
	Caption                     string          `json:"caption,omitempty"`
	ParseMode                   ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities             []MessageEntity `json:"caption_entities,omitempty"`
	DisableContentTypeDetection bool            `json:"disable_content_type_detection,omitempty"`

	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

// SendDocument sends a document.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#senddocument
func (api *API) SendDocument(params SendDocument) (Message, error) {
	req := NewRequestWithChatID[Message]("sendDocument", params, params.ChatID)
	return req.Do(api)
}

// SendDocumentWithContext is the context-aware variant of SendDocument.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#senddocument
func (api *API) SendDocumentWithContext(ctx context.Context, params SendDocument) (Message, error) {
	req := NewRequestWithChatID[Message]("sendDocument", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendVideo holds parameters for the sendVideo method.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendvideo
type SendVideo struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Video     string `json:"video"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Cover     string `json:"cover,omitempty"`

	StartTimestamp  int             `json:"start_timestamp,omitempty"`
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

// SendVideo sends a video.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sendvideo
func (api *API) SendVideo(params SendVideo) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVideo", params, params.ChatID)
	return req.Do(api)
}

// SendVideoWithContext is the context-aware variant of SendVideo.
// Since: Bot API 1.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvideo
func (api *API) SendVideoWithContext(ctx context.Context, params SendVideo) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVideo", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendAnimation holds parameters for the sendAnimation method.
// Since: Bot API 4.0
// See https://core.telegram.org/bots/api#sendanimation
type SendAnimation struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Animation string `json:"animation"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`

	Caption               string          `json:"caption,omitempty"`
	ParseMode             ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool            `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool            `json:"has_spoiler,omitempty"`
	DisableNotification   bool            `json:"disable_notification,omitempty"`
	ProtectContent        bool            `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool            `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string          `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

// SendAnimation sends an animation file (GIF or H.264/MPEG-4 AVC video without sound).
// Since: Bot API 4.0
// See https://core.telegram.org/bots/api#sendanimation
func (api *API) SendAnimation(params SendAnimation) (Message, error) {
	req := NewRequestWithChatID[Message]("sendAnimation", params, params.ChatID)
	return req.Do(api)
}

// SendAnimationWithContext is the context-aware variant of SendAnimation.
// Since: Bot API 4.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendanimation
func (api *API) SendAnimationWithContext(ctx context.Context, params SendAnimation) (Message, error) {
	req := NewRequestWithChatID[Message]("sendAnimation", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendVoice holds parameters for the sendVoice method.
// Since: Bot API 1.2
// See https://core.telegram.org/bots/api#sendvoice
type SendVoice struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Voice               string          `json:"voice"`
	Caption             string          `json:"caption,omitempty"`
	ParseMode           ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities     []MessageEntity `json:"caption_entities,omitempty"`
	Duration            int             `json:"duration,omitempty"`
	DisableNotification bool            `json:"disable_notification,omitempty"`
	ProtectContent      bool            `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool            `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string          `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

// SendVoice sends a voice note.
// Since: Bot API 1.2
// See https://core.telegram.org/bots/api#sendvoice
func (api *API) SendVoice(params SendVoice) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVoice", params, params.ChatID)
	return req.Do(api)
}

// SendVoiceWithContext is the context-aware variant of SendVoice.
// Since: Bot API 1.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvoice
func (api *API) SendVoiceWithContext(ctx context.Context, params SendVoice) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVoice", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendVideoNote holds parameters for the sendVideoNote method.
// Since: Bot API 3.0
// See https://core.telegram.org/bots/api#sendvideonote
type SendVideoNote struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	VideoNote           string `json:"video_note"`
	Thumbnail           string `json:"thumbnail,omitempty"`
	Duration            int    `json:"duration,omitempty"`
	Length              int    `json:"length,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

// SendVideoNote sends a video note (rounded video message).
// Since: Bot API 3.0
// See https://core.telegram.org/bots/api#sendvideonote
func (api *API) SendVideoNote(params SendVideoNote) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVideoNote", params, params.ChatID)
	return req.Do(api)
}

// SendVideoNoteWithContext is the context-aware variant of SendVideoNote.
// Since: Bot API 3.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvideonote
func (api *API) SendVideoNoteWithContext(ctx context.Context, params SendVideoNote) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVideoNote", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendPaidMedia holds parameters for the sendPaidMedia method.
// Since: Bot API 7.6
// See https://core.telegram.org/bots/api#sendpaidmedia
type SendPaidMedia struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`
	StarCount             int    `json:"star_count,omitempty"`

	Media                 []InputPaidMedia `json:"media"`
	Payload               string           `json:"payload,omitempty"`
	Caption               string           `json:"caption,omitempty"`
	ParseMode             ParseMode        `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity  `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool             `json:"show_caption_above_media,omitempty"`
	DisableNotification   bool             `json:"disable_notification,omitempty"`
	ProtectContent        bool             `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool             `json:"allow_paid_broadcast,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

// SendPaidMedia sends paid media.
// Since: Bot API 7.6
// See https://core.telegram.org/bots/api#sendpaidmedia
func (api *API) SendPaidMedia(params SendPaidMedia) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPaidMedia", params, params.ChatID)
	return req.Do(api)
}

// SendPaidMediaWithContext is the context-aware variant of SendPaidMedia.
// Since: Bot API 7.6
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendpaidmedia
func (api *API) SendPaidMediaWithContext(ctx context.Context, params SendPaidMedia) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPaidMedia", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendMediaGroup holds parameters for the sendMediaGroup method.
// Since: Bot API 3.5
// See https://core.telegram.org/bots/api#sendmediagroup
type SendMediaGroup struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Media               []InputMedia     `json:"media"`
	DisableNotification bool             `json:"disable_notification,omitempty"`
	ProtectContent      bool             `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool             `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string           `json:"message_effect_id,omitempty"`
	ReplyParameters     *ReplyParameters `json:"reply_parameters,omitempty"`
}

// SendMediaGroup sends a group of photos, videos, documents or audios as an album.
// Since: Bot API 3.5
// See https://core.telegram.org/bots/api#sendmediagroup
func (api *API) SendMediaGroup(params SendMediaGroup) ([]Message, error) {
	req := NewRequestWithChatID[[]Message]("sendMediaGroup", params, params.ChatID)
	return req.Do(api)
}

// SendMediaGroupWithContext is the context-aware variant of SendMediaGroup.
// Since: Bot API 3.5
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendmediagroup
func (api *API) SendMediaGroupWithContext(ctx context.Context, params SendMediaGroup) ([]Message, error) {
	req := NewRequestWithChatID[[]Message]("sendMediaGroup", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendLivePhoto holds parameters for the sendLivePhoto method.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#sendlivephoto
type SendLivePhoto struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	LivePhoto       string          `json:"live_photo"`
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

// SendLivePhoto sends a live photo.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#sendlivephoto
func (api *API) SendLivePhoto(params SendLivePhoto) (Message, error) {
	req := NewRequestWithChatID[Message]("sendLivePhoto", params, params.ChatID)
	return req.Do(api)
}

// SendLivePhotoWithContext is the context-aware variant of SendLivePhoto.
// Since: Bot API 10.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendlivephoto
func (api *API) SendLivePhotoWithContext(ctx context.Context, params SendLivePhoto) (Message, error) {
	req := NewRequestWithChatID[Message]("sendLivePhoto", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}
