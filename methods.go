package laniakea

import (
	"encoding/json"
	"fmt"
)

var NoParams = make(map[string]any)

func (b *Bot) Updates() ([]*Update, error) {
	params := make(map[string]any)
	params["offset"] = b.updateOffset
	params["timeout"] = 30
	params["allowed_updates"] = b.updateTypes

	data, err := b.request("getUpdates", params)
	if err != nil {
		return nil, err
	}
	res := make([]*Update, 0)
	err = AnyToStruct(data["data"], &res)
	if err != nil {
		return res, err
	}

	for _, u := range res {
		b.updateOffset = u.UpdateID + 1
		err = b.updateQueue.Enqueue(u)
		if err != nil {
			return res, err
		}
		res = append(res, u)

		if b.debug && b.requestLogger != nil {
			j, err := json.Marshal(u)
			if err != nil {
				b.logger.Error(err)
			}
			b.requestLogger.Debugln(fmt.Sprintf("UPDATE %s", j))
		}
	}
	return res, err
}

func (b *Bot) GetMe() (*User, error) {
	//req := NewRequest[User, EmptyParams]("getMe", EmptyParams{})
	//user, err := req.Do(b)
	data, err := b.request("getMe", NoParams)
	if err != nil {
		return nil, err
	}
	user := new(User)
	err = MapToStruct(data, user)
	return user, err
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

	data, err := b.request("sendMessage", params)
	if err != nil {
		return nil, err
	}
	message := new(Message)
	err = MapToStruct(data, message)
	return message, err
}

type SendPhotoP struct {
	BusinessConnectionID  string               `json:"business_connection_id,omitempty"`
	ChatID                int                  `json:"chat_id"`
	MessageThreadID       int                  `json:"message_thread_id,omitempty"`
	ParseMode             ParseMode            `json:"parse_mode,omitempty"`
	Photo                 string               `json:"photo"`
	Caption               string               `json:"caption,omitempty"`
	CaptionEntities       []*MessageEntity     `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                 `json:"show_caption_above_media"`
	HasSpoiler            bool                 `json:"has_spoiler"`
	DisableNotifications  bool                 `json:"disable_notifications,omitempty"`
	ProtectContent        bool                 `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool                 `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string               `json:"message_effect_id,omitempty"`
	ReplyMarkup           InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (b *Bot) SendPhoto(params *SendPhotoP) (*Message, error) {
	data, err := b.request("sendPhoto", params)
	if err != nil {
		return nil, err
	}
	message := new(Message)
	err = MapToStruct(data, message)
	return message, err
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
	data, err := b.request("editMessageText", params)
	if err != nil {
		return nil, err
	}
	message := new(Message)
	err = MapToStruct(data, message)
	return message, err
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
	data, err := b.request("editMessageCaption", params)
	if err != nil {
		return nil, err
	}
	message := new(Message)
	err = MapToStruct(data, message)
	return message, err
}

type DeleteMessageP struct {
	ChatID    int `json:"chat_id"`
	MessageID int `json:"message_id"`
}

func (b *Bot) DeleteMessage(params *DeleteMessageP) (*Message, error) {
	data, err := b.request("deleteMessage", params)
	if err != nil {
		return nil, err
	}
	message := new(Message)
	err = MapToStruct(data, message)
	return message, err
}
