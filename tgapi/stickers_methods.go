package tgapi

type SendStickerP struct {
	BusinessConnectionID  string `json:"business_connection_id,omitempty"`
	ChatID                int    `json:"chat_id"`
	MessageThreadID       int    `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int    `json:"direct_messages_topic_id,omitempty"`

	Sticker             string `json:"sticker"`
	Emoji               string `json:"emoji,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`
}

func (api *Api) SendSticker(p SendStickerP) (Message, error) {
	req := NewRequest[Message]("sendSticker", p)
	return req.Do(api)
}

type GetStickerSetP struct {
	Name string `json:"name"`
}

func (api *Api) GetStickerSet(p GetStickerSetP) (StickerSet, error) {
	req := NewRequest[StickerSet]("getStickerSet", p)
	return req.Do(api)
}

type GetCustomEmojiStickersP struct {
	CustomEmojiIDs []string `json:"custom_emoji_ids"`
}

func (api *Api) GetCustomEmojiStickers(p GetCustomEmojiStickersP) ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getCustomEmojiStickers", p)
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

func (api *Api) CreateNewStickerSet(p CreateNewStickerSetP) (bool, error) {
	req := NewRequest[bool]("createNewStickerSet", p)
	return req.Do(api)
}

type AddStickerToSetP struct {
	UserID  int          `json:"user_id"`
	Name    string       `json:"name"`
	Sticker InputSticker `json:"sticker"`
}

func (api *Api) AddStickerToSet(p AddStickerToSetP) (bool, error) {
	req := NewRequest[bool]("addStickerToSet", p)
	return req.Do(api)
}

type SetStickerPositionInSetP struct {
	Sticker  string `json:"sticker"`
	Position int    `json:"position"`
}

func (api *Api) SetStickerPosition(p SetStickerPositionInSetP) (bool, error) {
	req := NewRequest[bool]("setStickerPosition", p)
	return req.Do(api)
}

type DeleteStickerFromSetP struct {
	Sticker string `json:"sticker"`
}

func (api *Api) DeleteStickerFromSet(p DeleteStickerFromSetP) (bool, error) {
	req := NewRequest[bool]("deleteStickerFromSet", p)
	return req.Do(api)
}

type ReplaceStickerInSetP struct {
	UserID     int          `json:"user_id"`
	Name       string       `json:"name"`
	OldSticker string       `json:"old_sticker"`
	Sticker    InputSticker `json:"sticker"`
}

func (api *Api) ReplaceStickerInSet(p ReplaceStickerInSetP) (bool, error) {
	req := NewRequest[bool]("replaceStickerInSet", p)
	return req.Do(api)
}

type SetStickerEmojiListP struct {
	Sticker   string   `json:"sticker"`
	EmojiList []string `json:"emoji_list"`
}

func (api *Api) SetStickerEmojiList(p SetStickerEmojiListP) (bool, error) {
	req := NewRequest[bool]("setStickerEmojiList", p)
	return req.Do(api)
}

type SetStickerKeywordsP struct {
	Sticker  string   `json:"sticker"`
	Keywords []string `json:"keywords"`
}

func (api *Api) SetStickerKeywords(p SetStickerKeywordsP) (bool, error) {
	req := NewRequest[bool]("setStickerKeywords", p)
	return req.Do(api)
}

type SetStickerMaskPositionP struct {
	Sticker      string        `json:"sticker"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
}

func (api *Api) SetStickerMaskPosition(p SetStickerMaskPositionP) (bool, error) {
	req := NewRequest[bool]("setStickerMaskPosition", p)
	return req.Do(api)
}

type SetStickerSetTitleP struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func (api *Api) SetStickerSetTitle(p SetStickerSetTitleP) (bool, error) {
	req := NewRequest[bool]("setStickerSetTitle", p)
	return req.Do(api)
}

type SetStickerSetThumbnailP struct {
	Name      string             `json:"name"`
	UserID    int                `json:"user_id"`
	Thumbnail string             `json:"thumbnail"`
	Format    InputStickerFormat `json:"format"`
}

func (api *Api) SetStickerSetThumbnail(p SetStickerSetThumbnailP) (bool, error) {
	req := NewRequest[bool]("setStickerSetThumbnail", p)
	return req.Do(api)
}

type SetCustomEmojiStickerSetThumbnailP struct {
	Name          string `json:"name"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

func (api *Api) SetCustomEmojiStickerSetThumbnail(p SetStickerSetThumbnailP) (bool, error) {
	req := NewRequest[bool]("setCustomEmojiStickerSetThumbnail", p)
	return req.Do(api)
}

type DeleteStickerSetP struct {
	Name string `json:"name"`
}

func (api *Api) DeleteStickerSet(p DeleteStickerSetP) (bool, error) {
	req := NewRequest[bool]("deleteStickerSet", p)
	return req.Do(api)
}
