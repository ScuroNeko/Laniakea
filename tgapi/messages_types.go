package tgapi

import "git.nix13.pw/scuroneko/extypes"

type MessageReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type Message struct {
	MessageID            int    `json:"message_id"`
	MessageThreadID      int    `json:"message_thread_id,omitempty"`
	BusinessConnectionId string `json:"business_connection_id,omitempty"`
	From                 *User  `json:"from,omitempty"`

	SenderChat        *Chat  `json:"sender_chat,omitempty"`
	SenderBoostCount  int    `json:"sender_boost_count,omitempty"`
	SenderBusinessBot *User  `json:"sender_business_bot,omitempty"`
	SenderTag         string `json:"sender_tag,omitempty"`
	Chat              *Chat  `json:"chat,omitempty"`

	IsTopicMessage     bool     `json:"is_topic_message,omitempty"`
	IsAutomaticForward bool     `json:"is_automatic_forward,omitempty"`
	IsFromOffline      bool     `json:"is_from_offline,omitempty"`
	IsPaidPost         bool     `json:"is_paid_post,omitempty"`
	MediaGroupId       string   `json:"media_group_id,omitempty"`
	AuthorSignature    string   `json:"author_signature,omitempty"`
	PaidStarCount      int      `json:"paid_star_count,omitempty"`
	ReplyToMessage     *Message `json:"reply_to_message,omitempty"`

	Text string `json:"text"`

	Photo           extypes.Slice[*PhotoSize] `json:"photo,omitempty"`
	Caption         string                    `json:"caption,omitempty"`
	CaptionEntities []MessageEntity           `json:"caption_entities,omitempty"`

	Date     int `json:"date"`
	EditDate int `json:"edit_date"`

	ReplyMarkup *MessageReplyMarkup `json:"reply_markup,omitempty"`

	Entities           []MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	SuggestedPostInfo  *SuggestedPostInfo  `json:"suggested_post_info,omitempty"`

	EffectID string `json:"effect_id,omitempty"`
}

type InaccessibleMessage struct {
	Chat      Chat `json:"chat"`
	MessageID int  `json:"message_id"`
	Date      int  `json:"date"`
}

type MaybeInaccessibleMessage interface{ Message | InaccessibleMessage }

type MessageEntityType string

const (
	MessageEntityMention              MessageEntityType = "mention"
	MessageEntityHashtag              MessageEntityType = "hashtag"
	MessageEntityCashtag              MessageEntityType = "cashtag"
	MessageEntityBotCommand           MessageEntityType = "bot_command"
	MessageEntityUrl                  MessageEntityType = "url"
	MessageEntityEmail                MessageEntityType = "email"
	MessageEntityPhoneNumber          MessageEntityType = "phone_number"
	MessageEntityBold                 MessageEntityType = "bold"
	MessageEntityItalic               MessageEntityType = "italic"
	MessageEntityUnderline            MessageEntityType = "underline"
	MessageEntityStrike               MessageEntityType = "strikethrough"
	MessageEntitySpoiler              MessageEntityType = "spoiler"
	MessageEntityBlockquote           MessageEntityType = "blockquote"
	MessageEntityExpandableBlockquote MessageEntityType = "expandable_blockquote"
	MessageEntityCode                 MessageEntityType = "code"
	MessageEntityPre                  MessageEntityType = "pre"
	MessageEntityTextLink             MessageEntityType = "text_link"
	MessageEntityTextMention          MessageEntityType = "text_mention"
	MessageEntityCustomEmoji          MessageEntityType = "custom_emoji"
	MessageEntityDateTime             MessageEntityType = "date_time"
)

type MessageEntity struct {
	Type MessageEntityType `json:"type"`

	Offset        int    `json:"offset"`
	Length        int    `json:"length"`
	URL           string `json:"url,omitempty"`
	User          *User  `json:"user,omitempty"`
	Language      string `json:"language,omitempty"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`

	UnixTime       int    `json:"unix_time,omitempty"`
	DateTimeFormat string `json:"date_time_format,omitempty"`
}

type ReplyParameters struct {
	MessageID int `json:"message_id"`
	ChatID    int `json:"chat_id,omitempty"`

	AllowSendingWithoutReply bool             `json:"allow_sending_without_reply,omitempty"`
	Quote                    string           `json:"quote,omitempty"`
	QuoteParsingMode         string           `json:"quote_parsing_mode,omitempty"`
	QuoteEntities            []*MessageEntity `json:"quote_entities,omitempty"`
	QuotePosition            int              `json:"quote_position,omitempty"`
	ChecklistTaskID          int              `json:"checklist_task_id,omitempty"`
}

type LinkPreviewOptions struct {
	IsDisabled       bool   `json:"is_disabled,omitempty"`
	URL              string `json:"url,omitempty"`
	PreferSmallMedia bool   `json:"prefer_small_media,omitempty"`
	PreferLargeMedia bool   `json:"prefer_large_media,omitempty"`
	ShowAboveText    bool   `json:"show_above_text,omitempty"`
}

type ReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`

	Keyboard              [][]int `json:"keyboard,omitempty"`
	IsPersistent          bool    `json:"is_persistent,omitempty"`
	ResizeKeyboard        bool    `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool    `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string  `json:"input_field_placeholder,omitempty"`
	Selective             bool    `json:"selective,omitempty"`

	RemoveKeyboard bool `json:"remove_keyboard,omitempty"`

	ForceReply bool `json:"force_reply,omitempty"`
}
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
}

type KeyboardButtonStyle string
type InlineKeyboardButton struct {
	Text              string              `json:"text"`
	URL               string              `json:"url,omitempty"`
	CallbackData      string              `json:"callback_data,omitempty"`
	Style             KeyboardButtonStyle `json:"style,omitempty"`
	IconCustomEmojiID string              `json:"icon_custom_emoji_id,omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]int `json:"keyboard"`
}

type CallbackQuery struct {
	ID      string  `json:"id"`
	From    User    `json:"from"`
	Message Message `json:"message"`

	Data string `json:"data"`
}

type InputPollOption struct {
	Text          string           `json:"text"`
	TextParseMode ParseMode        `json:"text_parse_mode,omitempty"`
	TextEntities  []*MessageEntity `json:"text_entities,omitempty"`
}
type PollType string

const (
	PollTypeRegular PollType = "regular"
	PollTypeQuiz    PollType = "quiz"
)

type InputChecklistTask struct {
	ID           int              `json:"id"`
	Text         string           `json:"text"`
	ParseMode    ParseMode        `json:"parse_mode,omitempty"`
	TextEntities []*MessageEntity `json:"text_entities,omitempty"`
}
type InputChecklist struct {
	Title                   string               `json:"title"`
	ParseMode               ParseMode            `json:"parse_mode,omitempty"`
	TitleEntities           []*MessageEntity     `json:"title_entities,omitempty"`
	Tasks                   []InputChecklistTask `json:"tasks"`
	OtherCanAddTasks        bool                 `json:"other_can_add_tasks,omitempty"`
	OtherCanMarkTasksAsDone bool                 `json:"other_can_mark_tasks_as_done,omitempty"`
}

type ChatActionType string

const (
	ChatActionTyping          ChatActionType = "typing"
	ChatActionUploadPhoto     ChatActionType = "upload_photo"
	ChatActionUploadVideo     ChatActionType = "upload_video"
	ChatActionUploadVoice     ChatActionType = "upload_voice"
	ChatActionUploadDocument  ChatActionType = "upload_document"
	ChatActionChooseSticker   ChatActionType = "choose_sticker"
	ChatActionFindLocation    ChatActionType = "find_location"
	ChatActionUploadVideoNone ChatActionType = "upload_video_none"
)

type MessageReactionUpdated struct {
	Chat        *Chat          `json:"chat"`
	MessageID   int            `json:"message_id"`
	User        *User          `json:"user,omitempty"`
	ActorChat   *Chat          `json:"actor_chat"`
	Date        int            `json:"date"`
	OldReaction []ReactionType `json:"old_reaction"`
	NewReaction []ReactionType `json:"new_reaction"`
}

type MessageReactionCountUpdated struct {
	Chat      *Chat            `json:"chat"`
	MessageID int              `json:"message_id"`
	Date      int              `json:"date"`
	Reactions []*ReactionCount `json:"reactions"`
}
type ReactionType struct {
	Type string `json:"type"`
	// ReactionTypeEmoji
	Emoji *string `json:"emoji,omitempty"`
	// ReactionTypeCustomEmoji
	CustomEmojiID *string `json:"custom_emoji_id,omitempty"`
}
type ReactionCount struct {
	Type       ReactionType `json:"type"`
	TotalCount int          `json:"total_count"`
}

type SuggestedPostPrice struct {
	Currency string `json:"currency"`
	Amount   int    `json:"amount"`
}
type SuggestedPostInfo struct {
	State    string             `json:"state"` //State of the suggested post. Currently, it can be one of “pending”, “approved”, “declined”.
	Price    SuggestedPostPrice `json:"price"`
	SendDate int                `json:"send_date"`
}
type SuggestedPostParameters struct {
	Price    SuggestedPostPrice `json:"price"`
	SendDate int                `json:"send_date"`
}
