package laniakea

import (
	"encoding/json"
	"fmt"
)

type EmptyParams struct{}

var NoParams = EmptyParams{}

type UpdateParams struct {
	Offset         int      `json:"offset"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

func (b *Bot) Updates() ([]*Update, error) {
	params := UpdateParams{
		Offset:         b.updateOffset,
		Timeout:        30,
		AllowedUpdates: b.updateTypes,
	}

	req := NewRequest[[]*Update]("getUpdates", params)
	res, err := req.Do(b.api)
	if err != nil {
		return []*Update{}, err
	}
	updates := *res

	for _, u := range updates {
		b.updateOffset = u.UpdateID + 1
		err = b.updateQueue.Enqueue(u)
		if err != nil {
			return updates, err
		}

		if b.debug && b.requestLogger != nil {
			j, err := json.Marshal(u)
			if err != nil {
				b.logger.Error(err)
			}
			b.requestLogger.Debugln(fmt.Sprintf("UPDATE %s", j))
		}
	}
	return updates, err
}

func (api *Api) GetMe() (*User, error) {
	req := NewRequest[User, EmptyParams]("getMe", NoParams)
	return req.Do(api)
}
func (api *Api) LogOut() (bool, error) {
	req := NewRequest[bool, EmptyParams]("logOut", NoParams)
	res, err := req.Do(api)
	if err != nil {
		return false, err
	}
	return *res, nil
}
func (api *Api) Close() (bool, error) {
	req := NewRequest[bool, EmptyParams]("close", NoParams)
	res, err := req.Do(api)
	if err != nil {
		return false, err
	}
	return *res, nil
}

type SendMessageP struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               int    `json:"chat_id"`
	MessageThreadID      int    `json:"message_thread_id,omitempty"`
	DirectMessageTopicID int    `json:"direct_message_topic_id,omitempty"`

	Text                 string              `json:"text"`
	ParseMode            ParseMode           `json:"parse_mode,omitempty"`
	Entities             []*MessageEntity    `json:"entities,omitempty"`
	LinkPreviewOptions   *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	DisableNotifications bool                `json:"disable_notifications,omitempty"`
	ProtectContent       bool                `json:"protect_content,omitempty"`
	AllowPaidBroadcast   bool                `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID      string              `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *InlineKeyboardMarkup    `json:"reply_markup,omitempty"`
}

func (api *Api) SendMessage(params *SendMessageP) (*Message, error) {
	req := NewRequest[Message, SendMessageP]("sendMessage", *params)
	return req.Do(api)
}

type ForwardMessageP struct {
	ChatID               int `json:"chat_id"`
	MessageThreadID      int `json:"message_thread_id,omitempty"`
	DirectMessageTopicID int `json:"direct_message_topic_id,omitempty"`

	MessageID           int  `json:"message_id,omitempty"`
	FromChatID          int  `json:"from_chat_id,omitempty"`
	VideoStartTimestamp int  `json:"video_start_timestamp,omitempty"`
	DisableNotification bool `json:"disable_notification,omitempty"`
	ProtectContent      bool `json:"protect_content,omitempty"`

	MessageEffectID         string                   `json:"message_effect_id,omitempty"`
	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
}

func (api *Api) ForwardMessage(params ForwardMessageP) (*Message, error) {
	req := NewRequest[Message]("forwardMessage", params)
	return req.Do(api)
}

type ForwardMessagesP struct {
	ChatID               int `json:"chat_id"`
	MessageThreadID      int `json:"message_thread_id,omitempty"`
	DirectMessageTopicID int `json:"direct_message_topic_id,omitempty"`

	FromChatID          int   `json:"from_chat_id,omitempty"`
	MessageIDs          []int `json:"message_ids,omitempty"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
	ProtectContent      bool  `json:"protect_content,omitempty"`
}

func (api *Api) ForwardMessages(params ForwardMessagesP) ([]int, error) {
	req := NewRequest[[]int]("forwardMessages", params)
	res, err := req.Do(api)
	if err != nil {
		return []int{}, err
	}
	return *res, nil
}

type CopyMessageP struct {
	ChatID               int `json:"chat_id"`
	MessageThreadID      int `json:"message_thread_id,omitempty"`
	DirectMessageTopicID int `json:"direct_message_topic_id,omitempty"`

	FromChatID          int       `json:"from_chat_id"`
	MessageID           int       `json:"message_id"`
	VideoStartTimestamp int       `json:"video_start_timestamp,omitempty"`
	Caption             string    `json:"caption,omitempty"`
	ParseMode           ParseMode `json:"parse_mode,omitempty"`

	CaptionEntities       []*MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool             `json:"show_caption_above_media,omitempty"`
	DisableNotification   bool             `json:"disable_notification,omitempty"`
	ProtectContent        bool             `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool             `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string           `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *InlineKeyboardMarkup    `json:"reply_markup,omitempty"`
}

func (api *Api) CopyMessage(params CopyMessageP) (int, error) {
	req := NewRequest[int]("copyMessage", params)
	res, err := req.Do(api)
	if err != nil {
		return 0, err
	}
	return *res, nil
}

type CopyMessagesP struct {
	ChatID               int `json:"chat_id"`
	MessageThreadID      int `json:"message_thread_id,omitempty"`
	DirectMessageTopicID int `json:"direct_message_topic_id,omitempty"`

	FromChatID          int   `json:"from_chat_id,omitempty"`
	MessageIDs          []int `json:"message_ids,omitempty"`
	DisableNotification bool  `json:"disable_notification,omitempty"`
	ProtectContent      bool  `json:"protect_content,omitempty"`
	RemoveCaption       bool  `json:"remove_caption,omitempty"`
}

func (api *Api) CopyMessages(params CopyMessagesP) ([]int, error) {
	req := NewRequest[[]int]("copyMessages", params)
	res, err := req.Do(api)
	if err != nil {
		return []int{}, err
	}
	return *res, nil
}

type SendPhotoBaseP struct {
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
type SendPhotoP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int    `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Photo           string           `json:"photo"`
	Caption         string           `json:"caption,omitempty"`
	ParseMode       ParseMode        `json:"parse_mode,omitempty"`
	CaptionEntities []*MessageEntity `json:"caption_entities,omitempty"`

	ShowCaptionAboveMedia bool   `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool   `json:"has_spoiler,omitempty"`
	DisableNotifications  bool   `json:"disable_notifications,omitempty"`
	ProtectContent        bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *InlineKeyboardMarkup    `json:"reply_markup,omitempty"`
}

func (api *Api) SendPhoto(params *SendPhotoP) (*Message, error) {
	req := NewRequest[Message]("sendPhoto", params)
	return req.Do(api)
}

type EditMessageTextP struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int                   `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	Text                 string                `json:"text"`
	ParseMode            ParseMode             `json:"parse_mode,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (api *Api) EditMessageText(params *EditMessageTextP) (*Message, error) {
	req := NewRequest[Message]("editMessageText", params)
	return req.Do(api)
}

type EditMessageCaptionP struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int                   `json:"chat_id,omitempty"`
	MessageID            int                   `json:"message_id,omitempty"`
	InlineMessageID      string                `json:"inline_message_id,omitempty"`
	Caption              string                `json:"caption"`
	ParseMode            ParseMode             `json:"parse_mode,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (api *Api) EditMessageCaption(params *EditMessageCaptionP) (*Message, error) {
	req := NewRequest[Message]("editMessageCaption", params)
	return req.Do(api)
}

type DeleteMessageP struct {
	ChatID    int `json:"chat_id"`
	MessageID int `json:"message_id"`
}

func (api *Api) DeleteMessage(params *DeleteMessageP) (bool, error) {
	req := NewRequest[bool]("deleteMessage", params)
	ok, err := req.Do(api)
	if err != nil {
		return false, err
	}
	return *ok, err
}

type AnswerCallbackQueryP struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int    `json:"cache_time,omitempty"`
}

func (api *Api) AnswerCallbackQuery(params *AnswerCallbackQueryP) (bool, error) {
	req := NewRequest[bool]("answerCallbackQuery", params)
	ok, err := req.Do(api)
	if err != nil {
		return false, err
	}
	return *ok, err
}

type GetFileP struct {
	FileId string `json:"file_id"`
}

func (api *Api) GetFile(params *GetFileP) (*File, error) {
	req := NewRequest[File]("getFile", params)
	return req.Do(api)
}

type SendChatActionP struct {
	BusinessConnectionID string      `json:"business_connection_id,omitempty"`
	ChatID               int         `json:"chat_id"`
	MessageThreadID      int         `json:"message_thread_id,omitempty"`
	Action               ChatActions `json:"action"`
}

func (api *Api) SendChatAction(params SendChatActionP) (bool, error) {
	req := NewRequest[bool]("sendChatAction", params)
	res, err := req.Do(api)
	if err != nil {
		return false, err
	}
	return *res, err
}

type SetMessageReactionP struct {
	ChatId    int  `json:"chat_id"`
	MessageId int  `json:"message_id"`
	IsBig     bool `json:"is_big,omitempty"`
}
type SetMessageReactionEmojiP struct {
	SetMessageReactionP
	Reaction []ReactionTypeEmoji `json:"reaction"`
}

func (api *Api) SetMessageReactionEmoji(params SetMessageReactionEmojiP) (bool, error) {
	req := NewRequest[bool]("setMessageReaction", params)
	res, err := req.Do(api)
	if err != nil {
		return false, err
	}
	return *res, err
}

type SetMessageReactionCustomEmojiP struct {
	SetMessageReactionP
	Reaction []ReactionTypeCustomEmoji `json:"reaction"`
}

func (api *Api) SetMessageReactionCustom(params SetMessageReactionCustomEmojiP) (bool, error) {
	req := NewRequest[bool]("setMessageReaction", params)
	res, err := req.Do(api)
	if err != nil {
		return false, err
	}
	return *res, err
}

type SetMessageReactionPaidP struct {
	SetMessageReactionP
}

func (api *Api) SetMessageReactionPaid(params SetMessageReactionPaidP) (bool, error) {
	req := NewRequest[bool]("setMessageReaction", params)
	res, err := req.Do(api)
	if err != nil {
		return false, err
	}
	return *res, err
}
