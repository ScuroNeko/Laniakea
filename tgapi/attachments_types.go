package tgapi

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
// See https://core.telegram.org/bots/api#audio
type Audio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`

	Performer string     `json:"performer,omitempty"`
	Title     string     `json:"title,omitempty"`
	FileName  string     `json:"file_name,omitempty"`
	MimeType  string     `json:"mime_type,omitempty"`
	FileSize  int64      `json:"file_size,omitempty"`
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
}

type Document struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name"`
	MimeType     string     `json:"mime_type"`
	FileSize     int        `json:"file_size,omitempty"`
}

// Story represents a story.
type Story struct {
	Chat Chat `json:"chat"`
	ID   int  `json:"id"`
}

type Video struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Duration     int    `json:"duration"`

	Thumbnail      *PhotoSize     `json:"thumbnail,omitempty"`
	Cover          []PhotoSize    `json:"cover,omitempty"`
	StartTimestamp int64          `json:"start_timestamp"`
	Qualities      []VideoQuality `json:"qualities,omitempty"`
	FileName       string         `json:"file_name,omitempty"`
	MimeType       string         `json:"mime_type,omitempty"`
	FileSize       int64          `json:"file_size,omitempty"`
}

// VideoQuality describes an alternative quality for a video.
// See https://core.telegram.org/bots/api#videoquality
type VideoQuality struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Codec        string `json:"codec"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type VideoNote struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Length       int        `json:"length"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

type PaidMediaInfo struct {
	StarCount int         `json:"star_count"`
	PaidMedia []PaidMedia `json:"paid_media"`
}
type PaidMediaType string

const (
	PaidMediaPreviewType PaidMediaType = "preview"
	PaidMediaPhotoType   PaidMediaType = "photo"
	PaidMediaVideoType   PaidMediaType = "video"
)

type PaidMedia struct {
	Type PaidMediaType `json:"type,omitempty"`

	Width    int `json:"width,omitempty"`
	Height   int `json:"height,omitempty"`
	Duration int `json:"duration,omitempty"`

	Photo []PhotoSize `json:"photo,omitempty"`

	Video *Video `json:"video,omitempty"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
	Vcard       string `json:"vcard,omitempty"`
}

type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

// PollOption contains information about one answer option in a poll.
// See https://core.telegram.org/bots/api#polloption
type PollOption struct {
	PersistentID string          `json:"persistent_id"`
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities"`
	VoterCount   int             `json:"voter_count"`

	AddedByUser  *User `json:"added_by_user,omitempty"`
	AddedByChat  *Chat `json:"added_by_chat,omitempty"`
	AdditionDate int   `json:"addition_date,omitempty"`
}

// InputPollOption contains information about one answer option in a poll to be sent.
// See https://core.telegram.org/bots/api#inputpolloption
type InputPollOption struct {
	Text          string          `json:"text"`
	TextParseMode ParseMode       `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity `json:"text_entities,omitempty"`
}

type PollOptionAdded struct {
	PollMessage        *InaccessibleMessage `json:"poll_message,omitempty"`
	OptionPersistentID string               `json:"option_persistent_id"`
	OptionText         string               `json:"option_text"`
	OptionTextEntities []MessageEntity      `json:"option_text_entities,omitempty"`
}

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
// See https://core.telegram.org/bots/api#pollanswer
type PollAnswer struct {
	PollID              string   `json:"poll_id"`
	VoterChat           Chat     `json:"voter_chat"`
	User                User     `json:"user"`
	OptionIDs           []int    `json:"option_ids"`
	OptionPersistentIDs []string `json:"option_persistent_ids"`
}

// Poll contains information about a poll.
// See https://core.telegram.org/bots/api#poll
type Poll struct {
	ID               string          `json:"id"`
	Question         string          `json:"question"`
	QuestionEntities []MessageEntity `json:"question_entities"`
	Options          []PollOption    `json:"options"`
	TotalVoterCount  int             `json:"total_voter_count"`
	IsClosed         bool            `json:"is_closed"`
	IsAnonymous      bool            `json:"is_anonymous"`
	Type             PollType        `json:"type"`

	AllowsMultipleAnswers bool            `json:"allows_multiple_answers"`
	AllowsRevoting        bool            `json:"allows_revoting"`
	CorrectOptionIDs      []int           `json:"correct_option_ids,omitempty"`
	Explanation           string          `json:"explanation,omitempty"`
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"`
	OpenPeriod            int             `json:"open_period,omitempty"`
	CloseDate             int             `json:"close_date,omitempty"`
	Description           string          `json:"description,omitempty"`
	DescriptionEntities   []MessageEntity `json:"description_entities,omitempty"`
}

type ChecklistTask struct {
	ID              int             `json:"id"`
	Text            string          `json:"text"`
	TextEntities    []MessageEntity `json:"text_entities,omitempty"`
	CompletedByUser *User           `json:"completed_by_user,omitempty"`
	CompletedByChat *Chat           `json:"completed_by_chat,omitempty"`
	CompletionDate  int             `json:"completion_date,omitempty"`
}

type Checklist struct {
	Title                    string          `json:"title"`
	TitleEntities            []MessageEntity `json:"title_entities,omitempty"`
	Tasks                    []ChecklistTask `json:"tasks"`
	OthersCanAddTasks        bool            `json:"others_can_add_tasks,omitempty"`
	OthersCanMarkTasksAsDone bool            `json:"others_can_mark_tasks_as_done,omitempty"`
}

// InputChecklistTask describes a task in a checklist.
type InputChecklistTask struct {
	ID           int             `json:"id"`
	Text         string          `json:"text"`
	ParseMode    ParseMode       `json:"parse_mode,omitempty"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
}

// InputChecklist represents a checklist to be sent.
type InputChecklist struct {
	Title                   string               `json:"title"`
	ParseMode               ParseMode            `json:"parse_mode,omitempty"`
	TitleEntities           []MessageEntity      `json:"title_entities,omitempty"`
	Tasks                   []InputChecklistTask `json:"tasks"`
	OtherCanAddTasks        bool                 `json:"other_can_add_tasks,omitempty"`
	OtherCanMarkTasksAsDone bool                 `json:"other_can_mark_tasks_as_done,omitempty"`
}

type ChecklistTaskDone struct {
	ChecklistMessage       *Message `json:"checklist_message,omitempty"`
	MarkedAsDoneTaskIDs    []int    `json:"marked_as_done_task_ids,omitempty"`
	MarkedAsNotDoneTaskIDs []int    `json:"marked_as_not_done_task_ids,omitempty"`
}

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
