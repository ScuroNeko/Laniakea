package laniakea

import "fmt"

var NO_PARAMS = make(map[string]interface{})

func (b *Bot) Updates() ([]*Update, error) {
	params := make(map[string]interface{})
	params["offset"] = b.updateOffset
	params["timeout"] = 30
	params["allowed_updates"] = b.updateTypes

	data, err := b.request("getUpdates", params)
	if err != nil {
		return nil, err
	}
	res := make([]*Update, 0)
	for _, u := range data["data"].([]interface{}) {
		updateObj := new(Update)
		err = MapToStruct(u.(map[string]interface{}), updateObj)
		if err != nil {
			return res, err
		}
		b.updateOffset = updateObj.UpdateID + 1
		err = b.updateQueue.Enqueue(updateObj)
		if err != nil {
			return res, err
		}
		res = append(res, updateObj)

		if b.debug && b.requestLogger != nil {
			j, err := MapToJson(u.(map[string]interface{}))
			if err != nil {
				b.logger.Error(err)
			}
			b.requestLogger.Debug(fmt.Sprintf("UPDATE %s", j))
		}
	}
	return res, err
}

func (b *Bot) GetMe() (*User, error) {
	data, err := b.request("getMe", NO_PARAMS)
	if err != nil {
		return nil, err
	}
	user := new(User)
	err = MapToStruct(data, user)
	return user, err
}

type SendMessageP struct {
	BusinessConnectionID string                `json:"business_connection_id,omitempty"`
	ChatID               int                   `json:"chat_id"`
	MessageThreadID      int                   `json:"message_thread_id,omitempty"`
	Text                 string                `json:"text"`
	ParseMode            ParseMode             `json:"parse_mode,omitempty"`
	Entities             []*MessageEntity      `json:"entities,omitempty"`
	LinkPreviewOptions   *LinkPreviewOptions   `json:"link_preview_options,omitempty"`
	DisableNotifications bool                  `json:"disable_notifications,omitempty"`
	ProtectContent       bool                  `json:"protect_content,omitempty"`
	AllowPaidBroadcast   bool                  `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID      string                `json:"message_effect_id,omitempty"`
	ReplyParameters      *ReplyParameters      `json:"reply_parameters,omitempty"`
	InlineKeyboardMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	// ReplyKeyboardMarkup  *ReplyKeyboardMarkup  `json:"reply_markup,omitempty"`
}

func (b *Bot) SendMessage(params *SendMessageP) (*Message, error) {
	dataP, err := StructToMap(params)
	if err != nil {
		return nil, err
	}
	data, err := b.request("sendMessage", dataP)
	if err != nil {
		return nil, err
	}
	message := new(Message)
	err = MapToStruct(data, message)
	return message, err
}

type SendPhotoP struct {
	BusinessConnectionID  string           `json:"business_connection_id,omitempty"`
	ChatID                int              `json:"chat_id"`
	MessageThreadID       int              `json:"message_thread_id,omitempty"`
	Photo                 string           `json:"photo"`
	Caption               string           `json:"caption,omitempty"`
	ParseMode             ParseMode        `json:"parse_mode,omitempty"`
	CaptionEntities       []*MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool             `json:"show_caption_above_media"`
	HasSpoiler            bool             `json:"has_spoiler"`
	DisableNotifications  bool             `json:"disable_notifications,omitempty"`
	ProtectContent        bool             `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool             `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string           `json:"message_effect_id,omitempty"`
}

func (b *Bot) SendPhoto(params *SendPhotoP) (*Message, error) {
	dataP, err := StructToMap(params)
	if err != nil {
		return nil, err
	}
	data, err := b.request("sendPhoto", dataP)
	if err != nil {
		return nil, err
	}
	message := new(Message)
	err = MapToStruct(data, message)
	return message, err
}
