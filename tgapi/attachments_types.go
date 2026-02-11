package tgapi

type InputMedia interface {
	InputMediaPhoto | InputMediaVideo | InputMediaAudio | InputMediaDocument
}

type InputMediaType string

const (
	InputMediaTypeAnimation InputMediaType = "animation"
	InputMediaTypeDocument  InputMediaType = "document"
	InputMediaTypePhoto     InputMediaType = "photo"
	InputMediaTypeVideo     InputMediaType = "video"
	InputMediaTypeAudio     InputMediaType = "audio"
)

type InputMediaPhoto struct {
	Type  InputMediaType `json:"type"`
	Media string         `json:"media"`

	Caption               string          `json:"caption,omitempty"`
	ParseMode             ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool            `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool            `json:"has_spoiler,omitempty"`
}
type InputMediaVideo struct {
	Type  InputMediaType `json:"type"`
	Media string         `json:"media"`

	Cover                 string          `json:"cover"`
	StartTimestamp        int             `json:"start_timestamp"`
	Caption               string          `json:"caption,omitempty"`
	ParseMode             ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool            `json:"show_caption_above_media,omitempty"`

	Width             int  `json:"width,omitempty"`
	Height            int  `json:"height,omitempty"`
	Duration          int  `json:"duration,omitempty"`
	SupportsStreaming bool `json:"supports_streaming,omitempty"`
	HasSpoiler        bool `json:"has_spoiler,omitempty"`
}
type InputMediaAnimation struct {
	Type  InputMediaType `json:"type"`
	Media string         `json:"media"`

	Caption               string          `json:"caption,omitempty"`
	ParseMode             ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool            `json:"show_caption_above_media,omitempty"`
	Width                 int             `json:"width,omitempty"`
	Height                int             `json:"height,omitempty"`
	Duration              int             `json:"duration,omitempty"`
	HasSpoiler            bool            `json:"has_spoiler,omitempty"`
}
type InputMediaAudio struct {
	Type  InputMediaType `json:"type"`
	Media string         `json:"media"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Duration        int             `json:"duration,omitempty"`
	Performer       string          `json:"performer,omitempty"`
	Title           string          `json:"title,omitempty"`
}
type InputMediaDocument struct {
	Type  InputMediaType `json:"type"`
	Media string         `json:"media"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
}

type InputPaidMediaType string

const (
	InputPaidMediaTypeVideo InputPaidMediaType = "video"
	InputPaidMediaTypePhoto InputPaidMediaType = "photo"
)

type InputPaidMedia interface {
	InputPaidMediaVideo | InputPaidMediaPhoto
}
type InputPaidMediaPhoto struct {
	Type  InputPaidMediaType `json:"type"`
	Media string             `json:"media"`
}
type InputPaidMediaVideo struct {
	Type  InputPaidMediaType `json:"type"`
	Media string             `json:"media"`

	Cover             string `json:"cover"`
	StartTimestamp    int64  `json:"start_timestamp"`
	Width             int    `json:"width"`
	Height            int    `json:"height"`
	Duration          int    `json:"duration"`
	SupportsStreaming bool   `json:"supports_streaming"`
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}
