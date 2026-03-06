package tgapi

type SendStickerP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int64  `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Sticker             string `json:"sticker"`
	Emoji               string `json:"emoji,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`
}

func (api *API) SendSticker(params SendStickerP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendSticker", params, params.ChatID)
	return req.Do(api)
}

type GetStickerSetP struct {
	Name string `json:"name"`
}

func (api *API) GetStickerSet(params GetStickerSetP) (StickerSet, error) {
	req := NewRequest[StickerSet]("getStickerSet", params)
	return req.Do(api)
}

type GetCustomEmojiStickersP struct {
	CustomEmojiIDs []string `json:"custom_emoji_ids"`
}

func (api *API) GetCustomEmojiStickers(params GetCustomEmojiStickersP) ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getCustomEmojiStickers", params)
	return req.Do(api)
}

type CreateNewStickerSetP struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Title  string `json:"title"`

	Stickers        []InputSticker `json:"stickers"`
	StickerType     StickerType    `json:"sticker_type,omitempty"`
	NeedsRepainting bool           `json:"needs_repainting,omitempty"`
}

func (api *API) CreateNewStickerSet(params CreateNewStickerSetP) (bool, error) {
	req := NewRequest[bool]("createNewStickerSet", params)
	return req.Do(api)
}

type AddStickerToSetP struct {
	UserID  int          `json:"user_id"`
	Name    string       `json:"name"`
	Sticker InputSticker `json:"sticker"`
}

func (api *API) AddStickerToSet(params AddStickerToSetP) (bool, error) {
	req := NewRequest[bool]("addStickerToSet", params)
	return req.Do(api)
}

type SetStickerPositionInSetP struct {
	Sticker  string `json:"sticker"`
	Position int    `json:"position"`
}

func (api *API) SetStickerPosition(params SetStickerPositionInSetP) (bool, error) {
	req := NewRequest[bool]("setStickerPosition", params)
	return req.Do(api)
}

type DeleteStickerFromSetP struct {
	Sticker string `json:"sticker"`
}

func (api *API) DeleteStickerFromSet(params DeleteStickerFromSetP) (bool, error) {
	req := NewRequest[bool]("deleteStickerFromSet", params)
	return req.Do(api)
}

type ReplaceStickerInSetP struct {
	UserID     int          `json:"user_id"`
	Name       string       `json:"name"`
	OldSticker string       `json:"old_sticker"`
	Sticker    InputSticker `json:"sticker"`
}

func (api *API) ReplaceStickerInSet(params ReplaceStickerInSetP) (bool, error) {
	req := NewRequest[bool]("replaceStickerInSet", params)
	return req.Do(api)
}

type SetStickerEmojiListP struct {
	Sticker   string   `json:"sticker"`
	EmojiList []string `json:"emoji_list"`
}

func (api *API) SetStickerEmojiList(params SetStickerEmojiListP) (bool, error) {
	req := NewRequest[bool]("setStickerEmojiList", params)
	return req.Do(api)
}

type SetStickerKeywordsP struct {
	Sticker  string   `json:"sticker"`
	Keywords []string `json:"keywords"`
}

func (api *API) SetStickerKeywords(params SetStickerKeywordsP) (bool, error) {
	req := NewRequest[bool]("setStickerKeywords", params)
	return req.Do(api)
}

type SetStickerMaskPositionP struct {
	Sticker      string        `json:"sticker"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
}

func (api *API) SetStickerMaskPosition(params SetStickerMaskPositionP) (bool, error) {
	req := NewRequest[bool]("setStickerMaskPosition", params)
	return req.Do(api)
}

type SetStickerSetTitleP struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func (api *API) SetStickerSetTitle(params SetStickerSetTitleP) (bool, error) {
	req := NewRequest[bool]("setStickerSetTitle", params)
	return req.Do(api)
}

type SetStickerSetThumbnailP struct {
	Name      string             `json:"name"`
	UserID    int                `json:"user_id"`
	Thumbnail string             `json:"thumbnail"`
	Format    InputStickerFormat `json:"format"`
}

func (api *API) SetStickerSetThumbnail(params SetStickerSetThumbnailP) (bool, error) {
	req := NewRequest[bool]("setStickerSetThumbnail", params)
	return req.Do(api)
}

type SetCustomEmojiStickerSetThumbnailP struct {
	Name          string `json:"name"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

func (api *API) SetCustomEmojiStickerSetThumbnail(params SetStickerSetThumbnailP) (bool, error) {
	req := NewRequest[bool]("setCustomEmojiStickerSetThumbnail", params)
	return req.Do(api)
}

type DeleteStickerSetP struct {
	Name string `json:"name"`
}

func (api *API) DeleteStickerSet(params DeleteStickerSetP) (bool, error) {
	req := NewRequest[bool]("deleteStickerSet", params)
	return req.Do(api)
}
