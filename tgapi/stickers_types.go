package tgapi

// MaskPositionPoint represents the part of the face where a mask should be placed.
type MaskPositionPoint string

const (
	// MaskPositionForehead places the mask on the forehead.
	MaskPositionForehead MaskPositionPoint = "forehead"
	// MaskPositionEyes places the mask on the eyes.
	MaskPositionEyes MaskPositionPoint = "eyes"
	// MaskPositionMouth places the mask on the mouth.
	MaskPositionMouth MaskPositionPoint = "mouth"
	// MaskPositionChin places the mask on the chin.
	MaskPositionChin MaskPositionPoint = "chin"
)

// MaskPosition describes the position on faces where a mask should be placed by default.
// Since: Bot API 3.2
// See https://core.telegram.org/bots/api#maskposition
type MaskPosition struct {
	Point  MaskPositionPoint `json:"point"`
	XShift float32           `json:"x_shift"`
	YShift float32           `json:"y_shift"`
	Scale  float32           `json:"scale"`
}

// StickerType represents the type of a sticker.
type StickerType string

const (
	// StickerTypeRegular is a regular sticker.
	StickerTypeRegular StickerType = "regular"
	// StickerTypeMask is a mask sticker that can be placed on faces.
	StickerTypeMask StickerType = "mask"
	// StickerTypeCustomEmoji is a custom emoji sticker.
	StickerTypeCustomEmoji StickerType = "custom_emoji"
)

// Sticker represents a sticker.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#sticker
type Sticker struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`

	Type           StickerType   `json:"type"`                // Since: Bot API 6.2
	IsAnimated     bool          `json:"is_animated"`         // Since: Bot API 4.4
	IsVideo        bool          `json:"is_video"`            // Since: Bot API 5.7
	Thumbnail      *PhotoSize    `json:"thumbnail,omitempty"` // Since: Bot API 6.6
	Emoji          *string       `json:"emoji,omitempty"`
	SetName        *string       `json:"set_name,omitempty"`        // Since: Bot API 3.2
	MaskPosition   *MaskPosition `json:"mask_position,omitempty"`   // Since: Bot API 3.2
	CustomEmojiID  *string       `json:"custom_emoji_id,omitempty"` // Since: Bot API 6.2
	NeedRepainting *bool         `json:"need_repainting,omitempty"` // Since: Bot API 6.6
	FileSize       *int64        `json:"file_size,omitempty"`
}

// StickerSet represents a sticker set.
// Since: Bot API 3.2
// See https://core.telegram.org/bots/api#stickerset
type StickerSet struct {
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	StickerType StickerType `json:"sticker_type"`
	Stickers    []Sticker   `json:"stickers"`
	Thumbnail   *PhotoSize  `json:"thumbnail,omitempty"`
}

// InputStickerFormat represents the format of an input sticker.
type InputStickerFormat string

const (
	// InputStickerFormatStatic is a static sticker (WEBP).
	InputStickerFormatStatic InputStickerFormat = "static"
	// InputStickerFormatAnimated is an animated sticker (TGS).
	InputStickerFormatAnimated InputStickerFormat = "animated"
	// InputStickerFormatVideo is a video sticker (WEBM).
	InputStickerFormatVideo InputStickerFormat = "video"
)

// InputSticker describes a sticker to be added to a sticker set.
// Since: Bot API 6.6
// See https://core.telegram.org/bots/api#inputsticker
type InputSticker struct {
	Sticker      string             `json:"sticker"`
	Format       InputStickerFormat `json:"format"`
	EmojiList    []string           `json:"emoji_list"`
	MaskPosition *MaskPosition      `json:"mask_position,omitempty"`
	Keywords     []string           `json:"keywords,omitempty"`
}
