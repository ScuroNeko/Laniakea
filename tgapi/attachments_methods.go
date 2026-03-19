package tgapi

import "context"

// SendPhotoP holds parameters for the sendPhoto method.
// See https://core.telegram.org/bots/api#sendphoto
type SendPhotoP struct {
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
// See https://core.telegram.org/bots/api#sendphoto
func (api *API) SendPhoto(params SendPhotoP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPhoto", params, params.ChatID)
	return req.Do(api)
}

// SendPhotoWithContext is the context-aware variant of SendPhoto.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendphoto
func (api *API) SendPhotoWithContext(ctx context.Context, params SendPhotoP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPhoto", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendAudioP holds parameters for the sendAudio method.
// See https://core.telegram.org/bots/api#sendaudio
type SendAudioP struct {
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
// See https://core.telegram.org/bots/api#sendaudio
func (api *API) SendAudio(params SendAudioP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendAudio", params, params.ChatID)
	return req.Do(api)
}

// SendAudioWithContext is the context-aware variant of SendAudio.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendaudio
func (api *API) SendAudioWithContext(ctx context.Context, params SendAudioP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendAudio", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendDocumentP holds parameters for the sendDocument method.
// See https://core.telegram.org/bots/api#senddocument
type SendDocumentP struct {
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
// See https://core.telegram.org/bots/api#senddocument
func (api *API) SendDocument(params SendDocumentP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendDocument", params, params.ChatID)
	return req.Do(api)
}

// SendDocumentWithContext is the context-aware variant of SendDocument.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#senddocument
func (api *API) SendDocumentWithContext(ctx context.Context, params SendDocumentP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendDocument", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendVideoP holds parameters for the sendVideo method.
// See https://core.telegram.org/bots/api#sendvideo
type SendVideoP struct {
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
// See https://core.telegram.org/bots/api#sendvideo
func (api *API) SendVideo(params SendVideoP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVideo", params, params.ChatID)
	return req.Do(api)
}

// SendVideoWithContext is the context-aware variant of SendVideo.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvideo
func (api *API) SendVideoWithContext(ctx context.Context, params SendVideoP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVideo", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendAnimationP holds parameters for the sendAnimation method.
// See https://core.telegram.org/bots/api#sendanimation
type SendAnimationP struct {
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
// See https://core.telegram.org/bots/api#sendanimation
func (api *API) SendAnimation(params SendAnimationP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendAnimation", params, params.ChatID)
	return req.Do(api)
}

// SendAnimationWithContext is the context-aware variant of SendAnimation.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendanimation
func (api *API) SendAnimationWithContext(ctx context.Context, params SendAnimationP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendAnimation", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendVoiceP holds parameters for the sendVoice method.
// See https://core.telegram.org/bots/api#sendvoice
type SendVoiceP struct {
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
// See https://core.telegram.org/bots/api#sendvoice
func (api *API) SendVoice(params *SendVoiceP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVoice", params, params.ChatID)
	return req.Do(api)
}

// SendVoiceWithContext is the context-aware variant of SendVoice.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvoice
func (api *API) SendVoiceWithContext(ctx context.Context, params *SendVoiceP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVoice", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendVideoNoteP holds parameters for the sendVideoNote method.
// See https://core.telegram.org/bots/api#sendvideonote
type SendVideoNoteP struct {
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
// See https://core.telegram.org/bots/api#sendvideonote
func (api *API) SendVideoNote(params SendVideoNoteP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVideoNote", params, params.ChatID)
	return req.Do(api)
}

// SendVideoNoteWithContext is the context-aware variant of SendVideoNote.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendvideonote
func (api *API) SendVideoNoteWithContext(ctx context.Context, params SendVideoNoteP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendVideoNote", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendPaidMediaP holds parameters for the sendPaidMedia method.
// See https://core.telegram.org/bots/api#sendpaidmedia
type SendPaidMediaP struct {
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
// See https://core.telegram.org/bots/api#sendpaidmedia
func (api *API) SendPaidMedia(params SendPaidMediaP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPaidMedia", params, params.ChatID)
	return req.Do(api)
}

// SendPaidMediaWithContext is the context-aware variant of SendPaidMedia.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendpaidmedia
func (api *API) SendPaidMediaWithContext(ctx context.Context, params SendPaidMediaP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendPaidMedia", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// SendMediaGroupP holds parameters for the sendMediaGroup method.
// See https://core.telegram.org/bots/api#sendmediagroup
type SendMediaGroupP struct {
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
// See https://core.telegram.org/bots/api#sendmediagroup
func (api *API) SendMediaGroup(params SendMediaGroupP) ([]Message, error) {
	req := NewRequestWithChatID[[]Message]("sendMediaGroup", params, params.ChatID)
	return req.Do(api)
}

// SendMediaGroupWithContext is the context-aware variant of SendMediaGroup.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendmediagroup
func (api *API) SendMediaGroupWithContext(ctx context.Context, params SendMediaGroupP) ([]Message, error) {
	req := NewRequestWithChatID[[]Message]("sendMediaGroup", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}
