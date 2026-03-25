package tgapi

// InputMediaType represents the type of input media.
type InputMediaType string

const (
	// InputMediaTypeAnimation is a GIF or H.264/MPEG-4 AVC video without sound.
	InputMediaTypeAnimation InputMediaType = "animation"
	// InputMediaTypeDocument is a general file.
	InputMediaTypeDocument InputMediaType = "document"
	// InputMediaTypePhoto is a photo.
	InputMediaTypePhoto InputMediaType = "photo"
	// InputMediaTypeVideo is a video.
	InputMediaTypeVideo InputMediaType = "video"
	// InputMediaTypeAudio is an audio file.
	InputMediaTypeAudio InputMediaType = "audio"
)

// InputMedia represents the content of a media message to be sent.
// It is a union type described in https://core.telegram.org/bots/api#inputmedia.
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

// InputPaidMediaType represents the type of paid media.
type InputPaidMediaType string

const (
	// InputPaidMediaTypeVideo represents a paid video.
	InputPaidMediaTypeVideo InputPaidMediaType = "video"
	// InputPaidMediaTypePhoto represents a paid photo.
	InputPaidMediaTypePhoto InputPaidMediaType = "photo"
)

// InputPaidMedia describes the paid media to be sent.
// See https://core.telegram.org/bots/api#inputpaidmedia
type InputPaidMedia struct {
	Type  InputPaidMediaType `json:"type"`
	Media string             `json:"media"`

	Cover             *string `json:"cover,omitempty"`
	StartTimestamp    *int64  `json:"start_timestamp,omitempty"`
	Width             *int    `json:"width,omitempty"`
	Height            *int    `json:"height,omitempty"`
	Duration          *int    `json:"duration,omitempty"`
	SupportsStreaming *bool   `json:"supports_streaming,omitempty"`
}

// PhotoSize represents one size of a photo or a file/sticker thumbnail.
// See https://core.telegram.org/bots/api#photosize
type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int64  `json:"file_size,omitempty"`
}
