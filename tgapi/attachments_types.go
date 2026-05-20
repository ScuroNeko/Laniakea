package tgapi

// Animation represents an animation file (GIF or H.264/MPEG-4 AVC without sound).
// Since: Bot API 4.0
type Animation struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Duration     int    `json:"duration"`

	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
	FileName  string     `json:"file_name"`
	MimeType  string     `json:"mime_type"`
	FileSize  int        `json:"file_size"`
}

// Audio represents an audio file to be treated as music by the Telegram clients.
// Since: Bot API 1.2
// See https://core.telegram.org/bots/api#audio
type Audio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`

	Performer string     `json:"performer,omitempty"`
	Title     string     `json:"title,omitempty"`
	FileName  string     `json:"file_name,omitempty"` // Since: Bot API 5.0
	MimeType  string     `json:"mime_type,omitempty"`
	FileSize  int64      `json:"file_size,omitempty"`
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
}

// Document represents a general file (as opposed to photos, voice messages and audio files).
// Since: Bot API 1.0
type Document struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name"`
	MimeType     string     `json:"mime_type"`
	FileSize     int        `json:"file_size,omitempty"`
}

// Story represents a story.
// Since: Bot API 6.8
type Story struct {
	Chat Chat `json:"chat"`
	ID   int  `json:"id"`
}

// Video represents a video file.
// Since: Bot API 1.0
type Video struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Duration     int    `json:"duration"`

	Thumbnail      *PhotoSize     `json:"thumbnail,omitempty"`
	Cover          []PhotoSize    `json:"cover,omitempty"`     // Since: Bot API 8.3
	StartTimestamp int64          `json:"start_timestamp"`     // Since: Bot API 8.3
	Qualities      []VideoQuality `json:"qualities,omitempty"` // Since: Bot API 9.4
	FileName       string         `json:"file_name,omitempty"`
	MimeType       string         `json:"mime_type,omitempty"`
	FileSize       int64          `json:"file_size,omitempty"`
}

// VideoQuality describes an alternative quality for a video.
// Since: Bot API 9.4
// See https://core.telegram.org/bots/api#videoquality
type VideoQuality struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Codec        string `json:"codec"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// VideoNote represents a video message.
// Since: Bot API 3.0
type VideoNote struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Length       int        `json:"length"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

// Voice represents a voice note.
// Since: Bot API 1.2
type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

// PaidMediaInfo describes paid media.
// Since: Bot API 7.6
type PaidMediaInfo struct {
	StarCount int         `json:"star_count"`
	PaidMedia []PaidMedia `json:"paid_media"`
}

// PaidMediaType represents the type of paid media.
// Since: Bot API 7.6
type PaidMediaType string

const (
	PaidMediaPreviewType   PaidMediaType = "preview"
	PaidMediaPhotoType     PaidMediaType = "photo"
	PaidMediaVideoType     PaidMediaType = "video"
	PaidMediaLivePhotoType PaidMediaType = "live_photo" // Since: Bot API 10.0
)

// PaidMedia describes paid media content.
// Since: Bot API 7.6
type PaidMedia struct {
	Type PaidMediaType `json:"type,omitempty"`

	Width    int `json:"width,omitempty"`
	Height   int `json:"height,omitempty"`
	Duration int `json:"duration,omitempty"`

	Photo []PhotoSize `json:"photo,omitempty"`

	Video     *Video     `json:"video,omitempty"`
	LivePhoto *LivePhoto `json:"live_photo,omitempty"` // Since: Bot API 10.0
}

// Contact represents a phone contact.
// Since: Bot API 1.0
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
	Vcard       string `json:"vcard,omitempty"`
}

// Dice represents an animated emoji with a random value.
// Since: Bot API 4.7
type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

// PollOption contains information about one answer option in a poll.
// Since: Bot API 4.2
// See https://core.telegram.org/bots/api#polloption
type PollOption struct {
	PersistentID string          `json:"persistent_id"` // Since: Bot API 9.6
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities"`
	Media        *PollMedia      `json:"media,omitempty"` // Since: Bot API 10.0
	VoterCount   int             `json:"voter_count"`

	AddedByUser  *User `json:"added_by_user,omitempty"` // Since: Bot API 9.6
	AddedByChat  *Chat `json:"added_by_chat,omitempty"` // Since: Bot API 9.6
	AdditionDate int   `json:"addition_date,omitempty"` // Since: Bot API 9.6
}

// InputPollOptionMedia describes the media to attach to a poll option.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#inputpolloptionmedia
type InputPollOptionMedia struct {
	Type  string `json:"type"`
	Media string `json:"media"`
}

// InputPollOption contains information about one answer option in a poll to be sent.
// Since: Bot API 7.3
// See https://core.telegram.org/bots/api#inputpolloption
type InputPollOption struct {
	Text          string                `json:"text"`
	TextParseMode ParseMode             `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity       `json:"text_entities,omitempty"`
	Media         *InputPollOptionMedia `json:"media,omitempty"` // Since: Bot API 10.0
}

// InputPollMedia describes the media to attach to a poll or its explanation.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#inputpollmedia
type InputPollMedia struct {
	Type  string `json:"type"`
	Media string `json:"media"`
}

// PollOptionAdded describes a service message about a poll option being added.
// Since: Bot API 9.6
type PollOptionAdded struct {
	PollMessage        *InaccessibleMessage `json:"poll_message,omitempty"`
	OptionPersistentID string               `json:"option_persistent_id"`
	OptionText         string               `json:"option_text"`
	OptionTextEntities []MessageEntity      `json:"option_text_entities,omitempty"`
}

// PollOptionDeleted describes a service message about a poll option being deleted.
// Since: Bot API 9.6
type PollOptionDeleted struct {
	PollMessage        *InaccessibleMessage `json:"poll_message,omitempty"`
	OptionPersistentID string               `json:"option_persistent_id"`
	OptionText         string               `json:"option_text"`
	OptionTextEntities []MessageEntity      `json:"option_text_entities,omitempty"`
}

// PollType represents the type of a poll.
type PollType string

const (
	// PollTypeRegular identifies a regular poll.
	PollTypeRegular PollType = "regular"
	// PollTypeQuiz identifies a quiz poll.
	PollTypeQuiz PollType = "quiz"
)

// PollAnswer represents an answer of a user in a poll.
// Since: Bot API 4.6
// See https://core.telegram.org/bots/api#pollanswer
type PollAnswer struct {
	PollID              string   `json:"poll_id"`
	VoterChat           Chat     `json:"voter_chat"` // Since: Bot API 6.8
	User                User     `json:"user"`
	OptionIDs           []int    `json:"option_ids"`
	OptionPersistentIDs []string `json:"option_persistent_ids"` // Since: Bot API 9.6
}

// Poll contains information about a poll.
// Since: Bot API 4.2
// See https://core.telegram.org/bots/api#poll
type Poll struct {
	ID               string          `json:"id"`
	Question         string          `json:"question"`
	QuestionEntities []MessageEntity `json:"question_entities"` // Since: Bot API 7.3
	Options          []PollOption    `json:"options"`
	TotalVoterCount  int             `json:"total_voter_count"`
	IsClosed         bool            `json:"is_closed,omitempty"`
	IsAnonymous      bool            `json:"is_anonymous,omitempty"`
	Type             PollType        `json:"type"`

	AllowsMultipleAnswers bool            `json:"allows_multiple_answers,omitempty"` // Since: Bot API 4.6
	AllowsRevoting        bool            `json:"allows_revoting,omitempty"`         // Since: Bot API 9.6
	MembersOnly           bool            `json:"members_only,omitempty"`            // Since: Bot API 10.0
	CountryCodes          []string        `json:"country_codes,omitempty"`           // Since: Bot API 10.0
	CorrectOptionIDs      []int           `json:"correct_option_ids,omitempty"`      // Since: Bot API 9.6
	Explanation           string          `json:"explanation,omitempty"`             // Since: Bot API 4.8
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"`    // Since: Bot API 4.8
	ExplanationMedia      *PollMedia      `json:"explanation_media,omitempty"`       // Since: Bot API 10.0
	OpenPeriod            int             `json:"open_period,omitempty"`             // Since: Bot API 4.8
	CloseDate             int             `json:"close_date,omitempty"`              // Since: Bot API 4.8
	Description           string          `json:"description,omitempty"`             // Since: Bot API 9.6
	DescriptionEntities   []MessageEntity `json:"description_entities,omitempty"`    // Since: Bot API 9.6
	Media                 *PollMedia      `json:"media,omitempty"`                   // Since: Bot API 10.0
}

// PollMedia represents media attached to a poll.
// Since: Bot API 10.0
type PollMedia struct {
	Animation *Animation  `json:"animation,omitempty"`
	Audio     *Audio      `json:"audio,omitempty"`
	Document  *Document   `json:"document,omitempty"`
	LivePhoto *LivePhoto  `json:"live_photo,omitempty"`
	Location  *Location   `json:"location,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
	Sticker   *Sticker    `json:"sticker,omitempty"`
	Venue     *Venue      `json:"venue,omitempty"`
	Video     *Video      `json:"video,omitempty"`
}

// ChecklistTask represents a single task in a checklist.
// Since: Bot API 9.1
type ChecklistTask struct {
	ID              int             `json:"id"`
	Text            string          `json:"text"`
	TextEntities    []MessageEntity `json:"text_entities,omitempty"`
	CompletedByUser *User           `json:"completed_by_user,omitempty"`
	CompletedByChat *Chat           `json:"completed_by_chat,omitempty"`
	CompletionDate  int             `json:"completion_date,omitempty"`
}

// Checklist represents a checklist.
// Since: Bot API 9.1
type Checklist struct {
	Title                    string          `json:"title"`
	TitleEntities            []MessageEntity `json:"title_entities,omitempty"`
	Tasks                    []ChecklistTask `json:"tasks"`
	OthersCanAddTasks        bool            `json:"others_can_add_tasks,omitempty"`
	OthersCanMarkTasksAsDone bool            `json:"others_can_mark_tasks_as_done,omitempty"`
}

// InputChecklistTask describes a task in a checklist.
// Since: Bot API 9.1
type InputChecklistTask struct {
	ID           int             `json:"id"`
	Text         string          `json:"text"`
	ParseMode    ParseMode       `json:"parse_mode,omitempty"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
}

// InputChecklist represents a checklist to be sent.
// Since: Bot API 9.1
type InputChecklist struct {
	Title                   string               `json:"title"`
	ParseMode               ParseMode            `json:"parse_mode,omitempty"`
	TitleEntities           []MessageEntity      `json:"title_entities,omitempty"`
	Tasks                   []InputChecklistTask `json:"tasks"`
	OtherCanAddTasks        bool                 `json:"other_can_add_tasks,omitempty"`
	OtherCanMarkTasksAsDone bool                 `json:"other_can_mark_tasks_as_done,omitempty"`
}

// ChecklistTaskDone describes a service message about checklist tasks being marked as done.
// Since: Bot API 9.1
type ChecklistTaskDone struct {
	ChecklistMessage       *Message `json:"checklist_message,omitempty"`
	MarkedAsDoneTaskIDs    []int    `json:"marked_as_done_task_ids,omitempty"`
	MarkedAsNotDoneTaskIDs []int    `json:"marked_as_not_done_task_ids,omitempty"`
}

// ChecklistTasksAdded describes a service message about new checklist tasks being added.
// Since: Bot API 9.1
type ChecklistTasksAdded struct {
	ChecklistMessage *Message        `json:"checklist_message,omitempty"`
	Tasks            []ChecklistTask `json:"tasks"`
}

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

	InputMediaTypeSticker   InputMediaType = "sticker"
	InputMediaTypeLocation  InputMediaType = "location"
	InputMediaTypeVenue     InputMediaType = "venue"
	InputMediaTypeLivePhoto InputMediaType = "live_photo" // Since: Bot API 10.0
)

// InputMedia represents the content of a media message to be sent.
// Since: Bot API 4.0
// See https://core.telegram.org/bots/api#inputmedia
type InputMedia struct {
	Type  InputMediaType `json:"type"`
	Media string         `json:"media"`

	Caption               *string         `json:"caption,omitempty"`
	ParseMode             *ParseMode      `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool           `json:"show_caption_above_media,omitempty"` // Since: Bot API 7.4
	HasSpoiler            *bool           `json:"has_spoiler,omitempty"`              // Since: Bot API 6.4

	Cover             *string `json:"cover"`           // Since: Bot API 8.3
	StartTimestamp    *int    `json:"start_timestamp"` // Since: Bot API 8.3
	Width             *int    `json:"width,omitempty"`
	Height            *int    `json:"height,omitempty"`
	Duration          *int    `json:"duration,omitempty"`
	SupportsStreaming *bool   `json:"supports_streaming,omitempty"`

	Performer *string `json:"performer,omitempty"`
	Title     *string `json:"title,omitempty"`

	Emoji *string `json:"emoji,omitempty"`

	Latitude        *float64 `json:"latitude,omitempty"`
	Longitude       *float64 `json:"longitude,omitempty"`
	Address         *string  `json:"address,omitempty"`
	FoursquareID    *string  `json:"foursquare_id,omitempty"`
	FoursquareType  *string  `json:"foursquare_type,omitempty"`
	GooglePlaceID   *string  `json:"google_place_id,omitempty"`
	GooglePlaceType *string  `json:"google_place_type,omitempty"`

	HorizontalAccuracy *float64 `json:"horizontal_accuracy,omitempty"`
}

// InputPaidMediaType represents the type of paid media.
type InputPaidMediaType string

const (
	// InputPaidMediaTypeVideo represents a paid video.
	InputPaidMediaTypeVideo InputPaidMediaType = "video"
	// InputPaidMediaTypePhoto represents a paid photo.
	InputPaidMediaTypePhoto InputPaidMediaType = "photo"
	// InputPaidMediaTypeLivePhoto represents a paid live photo.
	InputPaidMediaTypeLivePhoto InputPaidMediaType = "live_photo" // Since: Bot API 10.0
)

// InputPaidMedia describes the paid media to be sent.
// Since: Bot API 7.6
// See https://core.telegram.org/bots/api#inputpaidmedia
type InputPaidMedia struct {
	Type  InputPaidMediaType `json:"type"`
	Media string             `json:"media"`

	Cover             *string `json:"cover,omitempty"`           // Since: Bot API 8.3
	StartTimestamp    *int64  `json:"start_timestamp,omitempty"` // Since: Bot API 8.3
	Width             *int    `json:"width,omitempty"`
	Height            *int    `json:"height,omitempty"`
	Duration          *int    `json:"duration,omitempty"`
	SupportsStreaming *bool   `json:"supports_streaming,omitempty"`
}

// PhotoSize represents one size of a photo or a file/sticker thumbnail.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#photosize
type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// LivePhoto represents a live photo (a photo with a short video attached).
// Since: Bot API 10.0
type LivePhoto struct {
	Photo        []PhotoSize `json:"photo,omitempty"`
	FileID       string      `json:"file_id"`
	FileUniqueID string      `json:"file_unique_id"`
	Width        int         `json:"width"`
	Height       int         `json:"height"`
	Duration     int         `json:"duration"`
	MIMEType     string      `json:"mime_type,omitempty"`
	FileSize     int64       `json:"file_size,omitempty"`
}
