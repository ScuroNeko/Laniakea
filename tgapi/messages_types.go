package tgapi

import "git.nix13.pw/scuroneko/extypes"

// MessageID represents a message identifier wrapper returned by some API methods.
type MessageID struct {
	MessageID int `json:"message_id"`
}

// MessageReplyMarkup represents an inline keyboard markup for a message.
// It is used in the Message type.
type MessageReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// DirectMessageTopic represents a forum topic in a direct message.
type DirectMessageTopic struct {
	TopicID int64 `json:"topic_id"`
	User    *User `json:"user,omitempty"`
}

// Message represents a Telegram message.
// See https://core.telegram.org/bots/api#message
type Message struct {
	MessageID            int                 `json:"message_id"`
	MessageThreadID      int                 `json:"message_thread_id,omitempty"`
	DirectMessageTopic   *DirectMessageTopic `json:"direct_message_topic,omitempty"`
	BusinessConnectionId string              `json:"business_connection_id,omitempty"`
	From                 *User               `json:"from,omitempty"`

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

// InaccessibleMessage describes a message that was deleted or is otherwise inaccessible.
// See https://core.telegram.org/bots/api#inaccessiblemessage
type InaccessibleMessage struct {
	Chat      Chat `json:"chat"`
	MessageID int  `json:"message_id"`
	Date      int  `json:"date"`
}

// MaybeInaccessibleMessage is a union type that can be either Message or InaccessibleMessage.
// See https://core.telegram.org/bots/api#maybeinaccessiblemessage
type MaybeInaccessibleMessage interface{ Message | InaccessibleMessage }

// MessageEntityType represents the type of a message entity.
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

// MessageEntity represents one special entity in a text message.
// See https://core.telegram.org/bots/api#messageentity
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

// ReplyParameters describes the parameters to use when replying to a message.
// See https://core.telegram.org/bots/api#replyparameters
type ReplyParameters struct {
	MessageID int   `json:"message_id"`
	ChatID    int64 `json:"chat_id,omitempty"`

	AllowSendingWithoutReply bool             `json:"allow_sending_without_reply,omitempty"`
	Quote                    string           `json:"quote,omitempty"`
	QuoteParsingMode         string           `json:"quote_parsing_mode,omitempty"`
	QuoteEntities            []*MessageEntity `json:"quote_entities,omitempty"`
	QuotePosition            int              `json:"quote_position,omitempty"`
	ChecklistTaskID          int              `json:"checklist_task_id,omitempty"`
}

// LinkPreviewOptions describes the options used for link preview generation.
// See https://core.telegram.org/bots/api#linkpreviewoptions
type LinkPreviewOptions struct {
	IsDisabled       bool   `json:"is_disabled,omitempty"`
	URL              string `json:"url,omitempty"`
	PreferSmallMedia bool   `json:"prefer_small_media,omitempty"`
	PreferLargeMedia bool   `json:"prefer_large_media,omitempty"`
	ShowAboveText    bool   `json:"show_above_text,omitempty"`
}

// ReplyMarkup represents a custom keyboard or inline keyboard.
// See https://core.telegram.org/bots/api#replymarkup
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

// InlineKeyboardMarkup represents an inline keyboard that appears right next to the message it belongs to.
// See https://core.telegram.org/bots/api#inlinekeyboardmarkup
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
}

// KeyboardButtonStyle represents the style of a keyboard button.
type KeyboardButtonStyle string

// InlineKeyboardButton represents one button of an inline keyboard.
// See https://core.telegram.org/bots/api#inlinekeyboardbutton
type InlineKeyboardButton struct {
	Text              string              `json:"text"`
	URL               string              `json:"url,omitempty"`
	CallbackData      string              `json:"callback_data,omitempty"`
	Style             KeyboardButtonStyle `json:"style,omitempty"`
	IconCustomEmojiID string              `json:"icon_custom_emoji_id,omitempty"`
}

// ReplyKeyboardMarkup represents a custom keyboard with reply options.
// See https://core.telegram.org/bots/api#replykeyboardmarkup
type ReplyKeyboardMarkup struct {
	Keyboard [][]int `json:"keyboard"`
}

// CallbackQuery represents an incoming callback query from a callback button in an inline keyboard.
// See https://core.telegram.org/bots/api#callbackquery
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            User     `json:"from"`
	Message         *Message `json:"message,omitempty"`
	InlineMessageID *string  `json:"inline_message_id,omitempty"`
	ChatInstance    string   `json:"chat_instance,omitempty"`
	Data            string   `json:"data,omitempty"`
	GameShortName   string   `json:"game_short_name,omitempty"`
}

// InputPollOption contains information about one answer option in a poll to be sent.
// See https://core.telegram.org/bots/api#inputpolloption
type InputPollOption struct {
	Text          string           `json:"text"`
	TextParseMode ParseMode        `json:"text_parse_mode,omitempty"`
	TextEntities  []*MessageEntity `json:"text_entities,omitempty"`
}

// PollType represents the type of a poll.
type PollType string

const (
	PollTypeRegular PollType = "regular"
	PollTypeQuiz    PollType = "quiz"
)

// InputChecklistTask describes a task in a checklist.
type InputChecklistTask struct {
	ID           int              `json:"id"`
	Text         string           `json:"text"`
	ParseMode    ParseMode        `json:"parse_mode,omitempty"`
	TextEntities []*MessageEntity `json:"text_entities,omitempty"`
}

// InputChecklist represents a checklist to be sent.
type InputChecklist struct {
	Title                   string               `json:"title"`
	ParseMode               ParseMode            `json:"parse_mode,omitempty"`
	TitleEntities           []*MessageEntity     `json:"title_entities,omitempty"`
	Tasks                   []InputChecklistTask `json:"tasks"`
	OtherCanAddTasks        bool                 `json:"other_can_add_tasks,omitempty"`
	OtherCanMarkTasksAsDone bool                 `json:"other_can_mark_tasks_as_done,omitempty"`
}

// ChatActionType represents the type of chat action.
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

// MessageReactionUpdated represents a change of a reaction on a message.
// See https://core.telegram.org/bots/api#messagereactionupdated
type MessageReactionUpdated struct {
	Chat        *Chat          `json:"chat"`
	MessageID   int            `json:"message_id"`
	User        *User          `json:"user,omitempty"`
	ActorChat   *Chat          `json:"actor_chat"`
	Date        int            `json:"date"`
	OldReaction []ReactionType `json:"old_reaction"`
	NewReaction []ReactionType `json:"new_reaction"`
}

// MessageReactionCountUpdated represents a change in the count of reactions on a message.
// See https://core.telegram.org/bots/api#messagereactioncountupdated
type MessageReactionCountUpdated struct {
	Chat      *Chat            `json:"chat"`
	MessageID int              `json:"message_id"`
	Date      int              `json:"date"`
	Reactions []*ReactionCount `json:"reactions"`
}

// ReactionType describes the type of a reaction.
// See https://core.telegram.org/bots/api#reactiontype
type ReactionType struct {
	Type string `json:"type"`
	// ReactionTypeEmoji
	Emoji *string `json:"emoji,omitempty"`
	// ReactionTypeCustomEmoji
	CustomEmojiID *string `json:"custom_emoji_id,omitempty"`
}

// ReactionCount represents a reaction added to a message along with the number of times it was added.
// See https://core.telegram.org/bots/api#reactioncount
type ReactionCount struct {
	Type       ReactionType `json:"type"`
	TotalCount int          `json:"total_count"`
}

// SuggestedPostPrice represents the price of a suggested post.
type SuggestedPostPrice struct {
	Currency string `json:"currency"`
	Amount   int    `json:"amount"`
}

// SuggestedPostInfo contains information about a suggested post.
// See https://core.telegram.org/bots/api#suggestedpostinfo
type SuggestedPostInfo struct {
	State    string             `json:"state"` // "pending", "approved", or "declined"
	Price    SuggestedPostPrice `json:"price"`
	SendDate int                `json:"send_date"`
}

// SuggestedPostParameters holds parameters for suggesting a post.
type SuggestedPostParameters struct {
	Price    SuggestedPostPrice `json:"price"`
	SendDate int                `json:"send_date"`
}
