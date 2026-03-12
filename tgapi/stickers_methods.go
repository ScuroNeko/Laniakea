package tgapi

// SendStickerP holds parameters for the sendSticker method.
// See https://core.telegram.org/bots/api#sendsticker
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

// SendSticker sends a static .WEBP, animated .TGS, or video .WEBM sticker.
// See https://core.telegram.org/bots/api#sendsticker
func (api *API) SendSticker(params SendStickerP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendSticker", params, params.ChatID)
	return req.Do(api)
}

// GetStickerSetP holds parameters for the getStickerSet method.
// See https://core.telegram.org/bots/api#getstickerset
type GetStickerSetP struct {
	Name string `json:"name"`
}

// GetStickerSet returns a sticker set by its name.
// See https://core.telegram.org/bots/api#getstickerset
func (api *API) GetStickerSet(params GetStickerSetP) (StickerSet, error) {
	req := NewRequest[StickerSet]("getStickerSet", params)
	return req.Do(api)
}

// GetCustomEmojiStickersP holds parameters for the getCustomEmojiStickers method.
// See https://core.telegram.org/bots/api#getcustomemojistickers
type GetCustomEmojiStickersP struct {
	CustomEmojiIDs []string `json:"custom_emoji_ids"`
}

// GetCustomEmojiStickers returns information about custom emoji stickers by their IDs.
// See https://core.telegram.org/bots/api#getcustomemojistickers
func (api *API) GetCustomEmojiStickers(params GetCustomEmojiStickersP) ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getCustomEmojiStickers", params)
	return req.Do(api)
}

// CreateNewStickerSetP holds parameters for the createNewStickerSet method.
// See https://core.telegram.org/bots/api#createnewstickerset
type CreateNewStickerSetP struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Title  string `json:"title"`

	Stickers        []InputSticker `json:"stickers"`
	StickerType     StickerType    `json:"sticker_type,omitempty"`
	NeedsRepainting bool           `json:"needs_repainting,omitempty"`
}

// CreateNewStickerSet creates a new sticker set owned by a user.
// Returns True on success.
// See https://core.telegram.org/bots/api#createnewstickerset
func (api *API) CreateNewStickerSet(params CreateNewStickerSetP) (bool, error) {
	req := NewRequest[bool]("createNewStickerSet", params)
	return req.Do(api)
}

// AddStickerToSetP holds parameters for the addStickerToSet method.
// See https://core.telegram.org/bots/api#addstickertoset
type AddStickerToSetP struct {
	UserID  int          `json:"user_id"`
	Name    string       `json:"name"`
	Sticker InputSticker `json:"sticker"`
}

// AddStickerToSet adds a new sticker to a set created by the bot.
// Returns True on success.
// See https://core.telegram.org/bots/api#addstickertoset
func (api *API) AddStickerToSet(params AddStickerToSetP) (bool, error) {
	req := NewRequest[bool]("addStickerToSet", params)
	return req.Do(api)
}

// SetStickerPositionInSetP holds parameters for the setStickerPositionInSet method.
// See https://core.telegram.org/bots/api#setstickerpositioninset
type SetStickerPositionInSetP struct {
	Sticker  string `json:"sticker"`
	Position int    `json:"position"`
}

// SetStickerPositionInSet moves a sticker in a set to a specific position.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickerpositioninset
func (api *API) SetStickerPositionInSet(params SetStickerPositionInSetP) (bool, error) {
	req := NewRequest[bool]("setStickerPositionInSet", params)
	return req.Do(api)
}

// DeleteStickerFromSetP holds parameters for the deleteStickerFromSet method.
// See https://core.telegram.org/bots/api#deletestickerfromset
type DeleteStickerFromSetP struct {
	Sticker string `json:"sticker"`
}

// DeleteStickerFromSet deletes a sticker from a set created by the bot.
// Returns True on success.
// See https://core.telegram.org/bots/api#deletestickerfromset
func (api *API) DeleteStickerFromSet(params DeleteStickerFromSetP) (bool, error) {
	req := NewRequest[bool]("deleteStickerFromSet", params)
	return req.Do(api)
}

// ReplaceStickerInSetP holds parameters for the replaceStickerInSet method.
// See https://core.telegram.org/bots/api#replacestickerinset
type ReplaceStickerInSetP struct {
	UserID     int          `json:"user_id"`
	Name       string       `json:"name"`
	OldSticker string       `json:"old_sticker"`
	Sticker    InputSticker `json:"sticker"`
}

// ReplaceStickerInSet replaces an existing sticker in a set with a new one.
// Returns True on success.
// See https://core.telegram.org/bots/api#replacestickerinset
func (api *API) ReplaceStickerInSet(params ReplaceStickerInSetP) (bool, error) {
	req := NewRequest[bool]("replaceStickerInSet", params)
	return req.Do(api)
}

// SetStickerEmojiListP holds parameters for the setStickerEmojiList method.
// See https://core.telegram.org/bots/api#setstickeremojilist
type SetStickerEmojiListP struct {
	Sticker   string   `json:"sticker"`
	EmojiList []string `json:"emoji_list"`
}

// SetStickerEmojiList changes the list of emoji associated with a sticker.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickeremojilist
func (api *API) SetStickerEmojiList(params SetStickerEmojiListP) (bool, error) {
	req := NewRequest[bool]("setStickerEmojiList", params)
	return req.Do(api)
}

// SetStickerKeywordsP holds parameters for the setStickerKeywords method.
// See https://core.telegram.org/bots/api#setstickerkeywords
type SetStickerKeywordsP struct {
	Sticker  string   `json:"sticker"`
	Keywords []string `json:"keywords"`
}

// SetStickerKeywords changes the keywords of a sticker.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickerkeywords
func (api *API) SetStickerKeywords(params SetStickerKeywordsP) (bool, error) {
	req := NewRequest[bool]("setStickerKeywords", params)
	return req.Do(api)
}

// SetStickerMaskPositionP holds parameters for the setStickerMaskPosition method.
// See https://core.telegram.org/bots/api#setstickermaskposition
type SetStickerMaskPositionP struct {
	Sticker      string        `json:"sticker"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
}

// SetStickerMaskPosition changes the mask position of a mask sticker.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickermaskposition
func (api *API) SetStickerMaskPosition(params SetStickerMaskPositionP) (bool, error) {
	req := NewRequest[bool]("setStickerMaskPosition", params)
	return req.Do(api)
}

// SetStickerSetTitleP holds parameters for the setStickerSetTitle method.
// See https://core.telegram.org/bots/api#setstickersettitle
type SetStickerSetTitleP struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

// SetStickerSetTitle sets the title of a sticker set created by the bot.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickersettitle
func (api *API) SetStickerSetTitle(params SetStickerSetTitleP) (bool, error) {
	req := NewRequest[bool]("setStickerSetTitle", params)
	return req.Do(api)
}

// SetStickerSetThumbnailP holds parameters for the setStickerSetThumbnail method.
// See https://core.telegram.org/bots/api#setstickersetthumbnail
type SetStickerSetThumbnailP struct {
	Name      string             `json:"name"`
	UserID    int                `json:"user_id"`
	Thumbnail string             `json:"thumbnail"`
	Format    InputStickerFormat `json:"format"`
}

// SetStickerSetThumbnail sets the thumbnail of a sticker set.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickersetthumbnail
func (api *API) SetStickerSetThumbnail(params SetStickerSetThumbnailP) (bool, error) {
	req := NewRequest[bool]("setStickerSetThumbnail", params)
	return req.Do(api)
}

// SetCustomEmojiStickerSetThumbnailP holds parameters for the setCustomEmojiStickerSetThumbnail method.
// See https://core.telegram.org/bots/api#setcustomemojistickersetthumbnail
type SetCustomEmojiStickerSetThumbnailP struct {
	Name          string `json:"name"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

// SetCustomEmojiStickerSetThumbnail sets the thumbnail of a custom emoji sticker set.
// Returns True on success.
// See https://core.telegram.org/bots/api#setcustomemojistickersetthumbnail
//
// Note: This method uses SetStickerSetThumbnailP as its parameter type, which might be inconsistent.
func (api *API) SetCustomEmojiStickerSetThumbnail(params SetStickerSetThumbnailP) (bool, error) {
	req := NewRequest[bool]("setCustomEmojiStickerSetThumbnail", params)
	return req.Do(api)
}

// DeleteStickerSetP holds parameters for the deleteStickerSet method.
// See https://core.telegram.org/bots/api#deletestickerset
type DeleteStickerSetP struct {
	Name string `json:"name"`
}

// DeleteStickerSet deletes a sticker set created by the bot.
// Returns True on success.
// See https://core.telegram.org/bots/api#deletestickerset
func (api *API) DeleteStickerSet(params DeleteStickerSetP) (bool, error) {
	req := NewRequest[bool]("deleteStickerSet", params)
	return req.Do(api)
}
