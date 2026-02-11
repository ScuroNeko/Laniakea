package tgapi

type InputMedia struct {
	Type  InputMediaType `json:"type"`
	Media string         `json:"media"`

	Caption               *string         `json:"caption,omitempty"`
	ParseMode             *ParseMode      `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool           `json:"show_caption_above_media,omitempty"`
	HasSpoiler            *bool           `json:"has_spoiler,omitempty"`

	Cover             *string `json:"cover"`
	StartTimestamp    *int    `json:"start_timestamp"`
	Width             *int    `json:"width,omitempty"`
	Height            *int    `json:"height,omitempty"`
	Duration          *int    `json:"duration,omitempty"`
	SupportsStreaming *bool   `json:"supports_streaming,omitempty"`

	Performer *string `json:"performer,omitempty"`
	Title     *string `json:"title,omitempty"`
}

type InputMediaType string

const (
	InputMediaTypeAnimation InputMediaType = "animation"
	InputMediaTypeDocument  InputMediaType = "document"
	InputMediaTypePhoto     InputMediaType = "photo"
	InputMediaTypeVideo     InputMediaType = "video"
	InputMediaTypeAudio     InputMediaType = "audio"
)

type InputPaidMediaType string

const (
	InputPaidMediaTypeVideo InputPaidMediaType = "video"
	InputPaidMediaTypePhoto InputPaidMediaType = "photo"
)

type InputPaidMedia struct {
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
