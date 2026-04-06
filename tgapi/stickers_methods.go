package tgapi

import "context"

// SendSticker holds parameters for the sendSticker method.
// See https://core.telegram.org/bots/api#sendsticker
type SendSticker struct {
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

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *ReplyMarkup             `json:"reply_markup,omitempty"`
}

// SendSticker sends a static .WEBP, animated .TGS, or video .WEBM sticker.
// See https://core.telegram.org/bots/api#sendsticker
func (api *API) SendSticker(params SendSticker) (Message, error) {
	req := NewRequestWithChatID[Message]("sendSticker", params, params.ChatID)
	return req.Do(api)
}

// SendStickerWithContext is the context-aware variant of SendSticker.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendsticker
func (api *API) SendStickerWithContext(ctx context.Context, params SendSticker) (Message, error) {
	req := NewRequestWithChatID[Message]("sendSticker", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// GetStickerSet holds parameters for the getStickerSet method.
// See https://core.telegram.org/bots/api#getstickerset
type GetStickerSet struct {
	Name string `json:"name"`
}

// GetStickerSet returns a sticker set by its name.
// See https://core.telegram.org/bots/api#getstickerset
func (api *API) GetStickerSet(params GetStickerSet) (StickerSet, error) {
	req := NewRequest[StickerSet]("getStickerSet", params)
	return req.Do(api)
}

// GetStickerSetWithContext is the context-aware variant of GetStickerSet.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getstickerset
func (api *API) GetStickerSetWithContext(ctx context.Context, params GetStickerSet) (StickerSet, error) {
	req := NewRequest[StickerSet]("getStickerSet", params)
	return req.DoWithContext(ctx, api)
}

// GetCustomEmojiStickers holds parameters for the getCustomEmojiStickers method.
// See https://core.telegram.org/bots/api#getcustomemojistickers
type GetCustomEmojiStickers struct {
	CustomEmojiIDs []string `json:"custom_emoji_ids"`
}

// GetCustomEmojiStickers returns information about custom emoji stickers by their IDs.
// See https://core.telegram.org/bots/api#getcustomemojistickers
func (api *API) GetCustomEmojiStickers(params GetCustomEmojiStickers) ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getCustomEmojiStickers", params)
	return req.Do(api)
}

// GetCustomEmojiStickersWithContext is the context-aware variant of GetCustomEmojiStickers.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getcustomemojistickers
func (api *API) GetCustomEmojiStickersWithContext(ctx context.Context, params GetCustomEmojiStickers) ([]Sticker, error) {
	req := NewRequest[[]Sticker]("getCustomEmojiStickers", params)
	return req.DoWithContext(ctx, api)
}

// UploadStickerFile holds parameters for the uploadStickerFile method.
// See https://core.telegram.org/bots/api#uploadstickerfile
type UploadStickerFile struct {
	UserID        int64              `json:"user_id"`
	StickerFormat InputStickerFormat `json:"sticker_format"`
}

// UploadStickerFile uploads a sticker file for later use in sticker set methods.
// sticker is the file to upload.
// See https://core.telegram.org/bots/api#uploadstickerfile
func (api *API) UploadStickerFile(params UploadStickerFile, sticker UploaderFile) (File, error) {
	uploader := NewUploader(api)
	defer func() {
		_ = uploader.Close()
	}()
	req := NewUploaderRequest[File]("uploadStickerFile", params, sticker.SetType(UploaderStickerType))
	return req.Do(uploader)
}

// UploadStickerFileWithContext is the context-aware variant of UploadStickerFile.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#uploadstickerfile
func (api *API) UploadStickerFileWithContext(ctx context.Context, params UploadStickerFile, sticker UploaderFile) (File, error) {
	uploader := NewUploader(api)
	defer func() {
		_ = uploader.Close()
	}()
	req := NewUploaderRequest[File]("uploadStickerFile", params, sticker.SetType(UploaderStickerType))
	return req.DoWithContext(ctx, uploader)
}

// CreateNewStickerSet holds parameters for the createNewStickerSet method.
// See https://core.telegram.org/bots/api#createnewstickerset
type CreateNewStickerSet struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Title  string `json:"title"`

	Stickers        []InputSticker `json:"stickers"`
	StickerType     StickerType    `json:"sticker_type,omitempty"`
	NeedsRepainting bool           `json:"needs_repainting,omitempty"`
}

// CreateNewStickerSet creates a new sticker set owned by a user.
// Returns True on success.
// See https://core.telegram.org/bots/api#createnewstickerset
func (api *API) CreateNewStickerSet(params CreateNewStickerSet) (bool, error) {
	req := NewRequest[bool]("createNewStickerSet", params)
	return req.Do(api)
}

// CreateNewStickerSetWithContext is the context-aware variant of CreateNewStickerSet.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#createnewstickerset
func (api *API) CreateNewStickerSetWithContext(ctx context.Context, params CreateNewStickerSet) (bool, error) {
	req := NewRequest[bool]("createNewStickerSet", params)
	return req.DoWithContext(ctx, api)
}

// AddStickerToSet holds parameters for the addStickerToSet method.
// See https://core.telegram.org/bots/api#addstickertoset
type AddStickerToSet struct {
	UserID  int64        `json:"user_id"`
	Name    string       `json:"name"`
	Sticker InputSticker `json:"sticker"`
}

// AddStickerToSet adds a new sticker to a set created by the bot.
// Returns True on success.
// See https://core.telegram.org/bots/api#addstickertoset
func (api *API) AddStickerToSet(params AddStickerToSet) (bool, error) {
	req := NewRequest[bool]("addStickerToSet", params)
	return req.Do(api)
}

// AddStickerToSetWithContext is the context-aware variant of AddStickerToSet.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#addstickertoset
func (api *API) AddStickerToSetWithContext(ctx context.Context, params AddStickerToSet) (bool, error) {
	req := NewRequest[bool]("addStickerToSet", params)
	return req.DoWithContext(ctx, api)
}

// SetStickerPositionInSet holds parameters for the setStickerPositionInSet method.
// See https://core.telegram.org/bots/api#setstickerpositioninset
type SetStickerPositionInSet struct {
	Sticker  string `json:"sticker"`
	Position int    `json:"position"`
}

// SetStickerPositionInSet moves a sticker in a set to a specific position.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickerpositioninset
func (api *API) SetStickerPositionInSet(params SetStickerPositionInSet) (bool, error) {
	req := NewRequest[bool]("setStickerPositionInSet", params)
	return req.Do(api)
}

// SetStickerPositionInSetWithContext is the context-aware variant of SetStickerPositionInSet.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setstickerpositioninset
func (api *API) SetStickerPositionInSetWithContext(ctx context.Context, params SetStickerPositionInSet) (bool, error) {
	req := NewRequest[bool]("setStickerPositionInSet", params)
	return req.DoWithContext(ctx, api)
}

// DeleteStickerFromSet holds parameters for the deleteStickerFromSet method.
// See https://core.telegram.org/bots/api#deletestickerfromset
type DeleteStickerFromSet struct {
	Sticker string `json:"sticker"`
}

// DeleteStickerFromSet deletes a sticker from a set created by the bot.
// Returns True on success.
// See https://core.telegram.org/bots/api#deletestickerfromset
func (api *API) DeleteStickerFromSet(params DeleteStickerFromSet) (bool, error) {
	req := NewRequest[bool]("deleteStickerFromSet", params)
	return req.Do(api)
}

// DeleteStickerFromSetWithContext is the context-aware variant of DeleteStickerFromSet.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletestickerfromset
func (api *API) DeleteStickerFromSetWithContext(ctx context.Context, params DeleteStickerFromSet) (bool, error) {
	req := NewRequest[bool]("deleteStickerFromSet", params)
	return req.DoWithContext(ctx, api)
}

// ReplaceStickerInSet holds parameters for the replaceStickerInSet method.
// See https://core.telegram.org/bots/api#replacestickerinset
type ReplaceStickerInSet struct {
	UserID     int64        `json:"user_id"`
	Name       string       `json:"name"`
	OldSticker string       `json:"old_sticker"`
	Sticker    InputSticker `json:"sticker"`
}

// ReplaceStickerInSet replaces an existing sticker in a set with a new one.
// Returns True on success.
// See https://core.telegram.org/bots/api#replacestickerinset
func (api *API) ReplaceStickerInSet(params ReplaceStickerInSet) (bool, error) {
	req := NewRequest[bool]("replaceStickerInSet", params)
	return req.Do(api)
}

// ReplaceStickerInSetWithContext is the context-aware variant of ReplaceStickerInSet.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#replacestickerinset
func (api *API) ReplaceStickerInSetWithContext(ctx context.Context, params ReplaceStickerInSet) (bool, error) {
	req := NewRequest[bool]("replaceStickerInSet", params)
	return req.DoWithContext(ctx, api)
}

// SetStickerEmojiList holds parameters for the setStickerEmojiList method.
// See https://core.telegram.org/bots/api#setstickeremojilist
type SetStickerEmojiList struct {
	Sticker   string   `json:"sticker"`
	EmojiList []string `json:"emoji_list"`
}

// SetStickerEmojiList changes the list of emoji associated with a sticker.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickeremojilist
func (api *API) SetStickerEmojiList(params SetStickerEmojiList) (bool, error) {
	req := NewRequest[bool]("setStickerEmojiList", params)
	return req.Do(api)
}

// SetStickerEmojiListWithContext is the context-aware variant of SetStickerEmojiList.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setstickeremojilist
func (api *API) SetStickerEmojiListWithContext(ctx context.Context, params SetStickerEmojiList) (bool, error) {
	req := NewRequest[bool]("setStickerEmojiList", params)
	return req.DoWithContext(ctx, api)
}

// SetStickerKeywords holds parameters for the setStickerKeywords method.
// See https://core.telegram.org/bots/api#setstickerkeywords
type SetStickerKeywords struct {
	Sticker  string   `json:"sticker"`
	Keywords []string `json:"keywords"`
}

// SetStickerKeywords changes the keywords of a sticker.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickerkeywords
func (api *API) SetStickerKeywords(params SetStickerKeywords) (bool, error) {
	req := NewRequest[bool]("setStickerKeywords", params)
	return req.Do(api)
}

// SetStickerKeywordsWithContext is the context-aware variant of SetStickerKeywords.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setstickerkeywords
func (api *API) SetStickerKeywordsWithContext(ctx context.Context, params SetStickerKeywords) (bool, error) {
	req := NewRequest[bool]("setStickerKeywords", params)
	return req.DoWithContext(ctx, api)
}

// SetStickerMaskPosition holds parameters for the setStickerMaskPosition method.
// See https://core.telegram.org/bots/api#setstickermaskposition
type SetStickerMaskPosition struct {
	Sticker      string        `json:"sticker"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
}

// SetStickerMaskPosition changes the mask position of a mask sticker.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickermaskposition
func (api *API) SetStickerMaskPosition(params SetStickerMaskPosition) (bool, error) {
	req := NewRequest[bool]("setStickerMaskPosition", params)
	return req.Do(api)
}

// SetStickerMaskPositionWithContext is the context-aware variant of SetStickerMaskPosition.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setstickermaskposition
func (api *API) SetStickerMaskPositionWithContext(ctx context.Context, params SetStickerMaskPosition) (bool, error) {
	req := NewRequest[bool]("setStickerMaskPosition", params)
	return req.DoWithContext(ctx, api)
}

// SetStickerSetTitle holds parameters for the setStickerSetTitle method.
// See https://core.telegram.org/bots/api#setstickersettitle
type SetStickerSetTitle struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

// SetStickerSetTitle sets the title of a sticker set created by the bot.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickersettitle
func (api *API) SetStickerSetTitle(params SetStickerSetTitle) (bool, error) {
	req := NewRequest[bool]("setStickerSetTitle", params)
	return req.Do(api)
}

// SetStickerSetTitleWithContext is the context-aware variant of SetStickerSetTitle.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setstickersettitle
func (api *API) SetStickerSetTitleWithContext(ctx context.Context, params SetStickerSetTitle) (bool, error) {
	req := NewRequest[bool]("setStickerSetTitle", params)
	return req.DoWithContext(ctx, api)
}

// SetStickerSetThumbnail holds parameters for the setStickerSetThumbnail method.
// See https://core.telegram.org/bots/api#setstickersetthumbnail
type SetStickerSetThumbnail struct {
	Name      string             `json:"name"`
	UserID    int64              `json:"user_id"`
	Thumbnail string             `json:"thumbnail"`
	Format    InputStickerFormat `json:"format"`
}

// SetStickerSetThumbnail sets the thumbnail of a sticker set.
// Returns True on success.
// See https://core.telegram.org/bots/api#setstickersetthumbnail
func (api *API) SetStickerSetThumbnail(params SetStickerSetThumbnail) (bool, error) {
	req := NewRequest[bool]("setStickerSetThumbnail", params)
	return req.Do(api)
}

// SetStickerSetThumbnailWithContext is the context-aware variant of SetStickerSetThumbnail.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setstickersetthumbnail
func (api *API) SetStickerSetThumbnailWithContext(ctx context.Context, params SetStickerSetThumbnail) (bool, error) {
	req := NewRequest[bool]("setStickerSetThumbnail", params)
	return req.DoWithContext(ctx, api)
}

// SetCustomEmojiStickerSetThumbnail holds parameters for the setCustomEmojiStickerSetThumbnail method.
// See https://core.telegram.org/bots/api#setcustomemojistickersetthumbnail
type SetCustomEmojiStickerSetThumbnail struct {
	Name          string `json:"name"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

// SetCustomEmojiStickerSetThumbnail sets the thumbnail of a custom emoji sticker set.
// Returns True on success.
// See https://core.telegram.org/bots/api#setcustomemojistickersetthumbnail
func (api *API) SetCustomEmojiStickerSetThumbnail(params SetCustomEmojiStickerSetThumbnail) (bool, error) {
	req := NewRequest[bool]("setCustomEmojiStickerSetThumbnail", params)
	return req.Do(api)
}

// SetCustomEmojiStickerSetThumbnailWithContext is the context-aware variant of SetCustomEmojiStickerSetThumbnail.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setcustomemojistickersetthumbnail
func (api *API) SetCustomEmojiStickerSetThumbnailWithContext(ctx context.Context, params SetCustomEmojiStickerSetThumbnail) (bool, error) {
	req := NewRequest[bool]("setCustomEmojiStickerSetThumbnail", params)
	return req.DoWithContext(ctx, api)
}

// DeleteStickerSet holds parameters for the deleteStickerSet method.
// See https://core.telegram.org/bots/api#deletestickerset
type DeleteStickerSet struct {
	Name string `json:"name"`
}

// DeleteStickerSet deletes a sticker set created by the bot.
// Returns True on success.
// See https://core.telegram.org/bots/api#deletestickerset
func (api *API) DeleteStickerSet(params DeleteStickerSet) (bool, error) {
	req := NewRequest[bool]("deleteStickerSet", params)
	return req.Do(api)
}

// DeleteStickerSetWithContext is the context-aware variant of DeleteStickerSet.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletestickerset
func (api *API) DeleteStickerSetWithContext(ctx context.Context, params DeleteStickerSet) (bool, error) {
	req := NewRequest[bool]("deleteStickerSet", params)
	return req.DoWithContext(ctx, api)
}
