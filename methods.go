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
	res, err := req.Do(b)
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

func (b *Bot) GetMe() (*User, error) {
	req := NewRequest[User, EmptyParams]("getMe", NoParams)
	return req.Do(b)
}

type SendMessageP struct {
	BusinessConnectionID string               `json:"business_connection_id,omitempty"`
	ChatID               int                  `json:"chat_id"`
	MessageThreadID      int                  `json:"message_thread_id,omitempty"`
	ParseMode            ParseMode            `json:"parse_mode,omitempty"`
	Text                 string               `json:"text"`
	Entities             []*MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions   *LinkPreviewOptions  `json:"link_preview_options,omitempty"`
	DisableNotifications bool                 `json:"disable_notifications,omitempty"`
	ProtectContent       bool                 `json:"protect_content,omitempty"`
	AllowPaidBroadcast   bool                 `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID      string               `json:"message_effect_id,omitempty"`
	ReplyParameters      *ReplyParameters     `json:"reply_parameters,omitempty"`
	ReplyMarkup          InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (b *Bot) SendMessage(params *SendMessageP) (*Message, error) {
	req := NewRequest[Message, SendMessageP]("sendMessage", *params)
	return req.Do(b)
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
	BusinessConnectionID  string               `json:"business_connection_id,omitempty"`
	ChatID                int                  `json:"chat_id"`
	MessageThreadID       int                  `json:"message_thread_id,omitempty"`
	ParseMode             ParseMode            `json:"parse_mode,omitempty"`
	Caption               string               `json:"caption,omitempty"`
	CaptionEntities       []*MessageEntity     `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                 `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool                 `json:"has_spoiler,omitempty"`
	DisableNotifications  bool                 `json:"disable_notifications,omitempty"`
	ProtectContent        bool                 `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool                 `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string               `json:"message_effect_id,omitempty"`
	ReplyMarkup           InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	Photo                 string               `json:"photo"`
}

func (b *Bot) SendPhoto(params *SendPhotoP) (*Message, error) {
	req := NewRequest[Message]("sendPhoto", params)
	return req.Do(b)
}

type EditMessageTextP struct {
	BusinessConnectionID string               `json:"business_connection_id,omitempty"`
	ChatID               int                  `json:"chat_id,omitempty"`
	MessageID            int                  `json:"message_id,omitempty"`
	InlineMessageID      string               `json:"inline_message_id,omitempty"`
	Text                 string               `json:"text"`
	ParseMode            ParseMode            `json:"parse_mode,omitempty"`
	ReplyMarkup          InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (b *Bot) EditMessageText(params *EditMessageTextP) (*Message, error) {
	req := NewRequest[Message]("editMessageText", params)
	return req.Do(b)
}

type EditMessageCaptionP struct {
	BusinessConnectionID string               `json:"business_connection_id,omitempty"`
	ChatID               int                  `json:"chat_id,omitempty"`
	MessageID            int                  `json:"message_id,omitempty"`
	InlineMessageID      string               `json:"inline_message_id,omitempty"`
	Caption              string               `json:"caption"`
	ParseMode            ParseMode            `json:"parse_mode,omitempty"`
	ReplyMarkup          InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (b *Bot) EditMessageCaption(params *EditMessageCaptionP) (*Message, error) {
	req := NewRequest[Message]("editMessageCaption", params)
	return req.Do(b)
}

type DeleteMessageP struct {
	ChatID    int `json:"chat_id"`
	MessageID int `json:"message_id"`
}

func (b *Bot) DeleteMessage(params *DeleteMessageP) (bool, error) {
	req := NewRequest[bool]("deleteMessage", params)
	ok, err := req.Do(b)
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

func (b *Bot) AnswerCallbackQuery(params *AnswerCallbackQueryP) (bool, error) {
	req := NewRequest[bool]("answerCallbackQuery", params)
	ok, err := req.Do(b)
	if err != nil {
		return false, err
	}
	return *ok, err
}

type GetFileP struct {
	FileId string `json:"file_id"`
}

func (b *Bot) GetFile(params *GetFileP) (*File, error) {
	req := NewRequest[File]("getFile", params)
	return req.Do(b)
}

type SendChatActionP struct {
	BusinessConnectionID string      `json:"business_connection_id,omitempty"`
	ChatID               int         `json:"chat_id"`
	MessageThreadID      int         `json:"message_thread_id,omitempty"`
	Action               ChatActions `json:"action"`
}

func (b *Bot) SendChatAction(params SendChatActionP) (bool, error) {
	req := NewRequest[bool]("sendChatAction", params)
	res, err := req.Do(b)
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

func (b *Bot) SetMessageReactionEmoji(params SetMessageReactionEmojiP) (bool, error) {
	req := NewRequest[bool]("setMessageReaction", params)
	res, err := req.Do(b)
	if err != nil {
		return false, err
	}
	return *res, err
}

type SetMessageReactionCustomEmojiP struct {
	SetMessageReactionP
	Reaction []ReactionTypeCustomEmoji `json:"reaction"`
}

func (b *Bot) SetMessageReactionCustom(params SetMessageReactionCustomEmojiP) (bool, error) {
	req := NewRequest[bool]("setMessageReaction", params)
	res, err := req.Do(b)
	if err != nil {
		return false, err
	}
	return *res, err
}

type SetMessageReactionPaidP struct {
	SetMessageReactionP
}

func (b *Bot) SetMessageReactionPaid(params SetMessageReactionPaidP) (bool, error) {
	req := NewRequest[bool]("setMessageReaction", params)
	res, err := req.Do(b)
	if err != nil {
		return false, err
	}
	return *res, err
}
