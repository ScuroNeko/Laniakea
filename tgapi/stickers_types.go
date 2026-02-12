package tgapi

type MaskPositionPoint string

const (
	MaskPositionForehead MaskPositionPoint = "forehead"
	MaskPositionEyes     MaskPositionPoint = "eyes"
	MaskPositionMouth    MaskPositionPoint = "mouth"
	MaskPositionChin     MaskPositionPoint = "chin"
)

type MaskPosition struct {
	Point  MaskPositionPoint `json:"point"`
	XShift float32           `json:"x_shift"`
	YShift float32           `json:"y_shift"`
	Scale  float32           `json:"scale"`
}

type StickerType string

const (
	StickerTypeRegular     StickerType = "regular"
	StickerTypeMask        StickerType = "mask"
	StickerTypeCustomEmoji StickerType = "custom_emoji"
)

type Sticker struct {
	FileId       string      `json:"file_id"`
	FileUniqueId string      `json:"file_unique_id"`
	Type         StickerType `json:"type"`
	Width        int         `json:"width"`
	Height       int         `json:"height"`
	IsAnimated   bool        `json:"is_animated"`
	IsVideo      bool        `json:"is_video"`

	Thumbnail      *PhotoSize    `json:"thumbnail,omitempty"`
	Emoji          *string       `json:"emoji,omitempty"`
	SetName        *string       `json:"set_name,omitempty"`
	MaskPosition   *MaskPosition `json:"mask_position,omitempty"`
	CustomEmojiID  *string       `json:"custom_emoji_id,omitempty"`
	NeedRepainting *bool         `json:"need_repainting,omitempty"`
	FileSize       *int          `json:"file_size,omitempty"`
}
type StickerSet struct {
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	StickerType StickerType `json:"sticker_type"`
	Stickers    []Sticker   `json:"stickers"`
	Thumbnail   *PhotoSize  `json:"thumbnail,omitempty"`
}
type InputStickerFormat string

const (
	InputStickerFormatStatic   InputStickerFormat = "static"
	InputStickerFormatAnimated InputStickerFormat = "animated"
	InputStickerFormatVideo    InputStickerFormat = "video"
)

type InputSticker struct {
	Sticker      string             `json:"sticker"`
	Format       InputStickerFormat `json:"format"`
	EmojiList    []string           `json:"emoji_list"`
	MaskPosition *MaskPosition      `json:"mask_position,omitempty"`
	Keywords     []string           `json:"keywords,omitempty"`
}
