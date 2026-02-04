package laniakea

type Update struct {
	UpdateID               int                          `json:"update_id"`
	Message                *Message                     `json:"message"`
	EditedMessage          *Message                     `json:"edited_message,omitempty"`
	ChannelPost            *Message                     `json:"channel_post,omitempty"`
	EditedChannelPost      *Message                     `json:"edited_channel_post,omitempty"`
	BusinessConnection     *BusinessConnection          `json:"business_connection,omitempty"`
	BusinessMessage        *Message                     `json:"business_message,omitempty"`
	EditedBusinessMessage  *Message                     `json:"edited_business_message,omitempty"`
	DeletedBusinessMessage *Message                     `json:"deleted_business_messages,omitempty"`
	MessageReaction        *MessageReactionUpdated      `json:"message_reaction,omitempty"`
	MessageReactionCount   *MessageReactionCountUpdated `json:"message_reaction_count,omitempty"`
	CallbackQuery          *CallbackQuery               `json:"callback_query,omitempty"`
	InlineQuery            int
	ChosenInlineResult     int
}

type User struct {
	ID                      int    `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name,omitempty"`
	Username                string `json:"username,omitempty"`
	LanguageCode            string `json:"language_code,omitempty"`
	IsPremium               bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
	CanConnectToBusiness    bool   `json:"can_connect_to_business,omitempty"`
	HasMainWebApp           bool   `json:"has_main_web_app,omitempty"`
}

type Chat struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	IsForum   bool   `json:"is_forum,omitempty"`
}

type MessageReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type Message struct {
	MessageID       int    `json:"message_id"`
	MessageThreadID int    `json:"message_thread_id,omitempty"`
	From            *User  `json:"from,omitempty"`
	Chat            *Chat  `json:"chat,omitempty"`
	Text            string `json:"text"`

	Photo          Slice[*PhotoSize] `json:"photo,omitempty"`
	Caption        string            `json:"caption,omitempty"`
	ReplyToMessage *Message          `json:"reply_to_message"`

	ReplyMarkup *MessageReplyMarkup `json:"reply_markup,omitempty"`
}

type InaccessableMessage struct {
	Chat      *Chat `json:"chat"`
	MessageID int   `json:"message_id"`
	Date      int   `json:"date"`
}

type MaybeInaccessibleMessage struct {
}

type MessageEntity struct {
	Type          string `json:"type"`
	Offset        int    `json:"offset"`
	Length        int    `json:"length"`
	URL           string `json:"url,omitempty"`
	User          *User  `json:"user,omitempty"`
	Language      string `json:"language,omitempty"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

type ReplyParameters struct {
	MessageID                int              `json:"message_id"`
	ChatID                   int              `json:"chat_id,omitempty"`
	AllowSendingWithoutReply bool             `json:"allow_sending_without_reply,omitempty"`
	Quote                    string           `json:"quote,omitempty"`
	QuoteParsingMode         string           `json:"quote_parsing_mode,omitempty"`
	QuoteEntities            []*MessageEntity `json:"quote_entities,omitempty"`
	QuotePosition            int              `json:"quote_postigin,omitempty"`
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}

type LinkPreviewOptions struct {
	IsDisabled       bool   `json:"is_disabled,omitempty"`
	URL              string `json:"url,omitempty"`
	PreferSmallMedia bool   `json:"prefer_small_media,omitempty"`
	PreferLargeMedia bool   `json:"prefer_large_media,omitempty"`
	ShowAboveText    bool   `json:"show_above_text,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]int `json:"keyboard"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    *User    `json:"from"`
	Message *Message `json:"message"`

	Data string `json:"data"`
}

type BusinessConnection struct {
	ID         string `json:"id"`
	User       *User  `json:"user"`
	UserChatID int    `json:"user_chat_id"`
	Date       int    `json:"date"`
	CanReply   bool   `json:"can_reply"`
	IsEnabled  bool   `json:"id_enabled"`
}

type MessageReactionUpdated struct {
	Chat        *Chat           `json:"chat"`
	MessageID   int             `json:"message_id"`
	User        *User           `json:"user,omitempty"`
	ActorChat   *Chat           `json:"actor_chat"`
	Date        int             `json:"date"`
	OldReaction []*ReactionType `json:"old_reaction"`
	NewReaction []*ReactionType `json:"new_reaction"`
}

type MessageReactionCountUpdated struct {
	Chat      *Chat            `json:"chat"`
	MessageID int              `json:"message_id"`
	Date      int              `json:"date"`
	Reactions []*ReactionCount `json:"reactions"`
}

type ReactionType struct {
	Type string `json:"type"`
}
type ReactionTypeEmoji struct {
	ReactionType
	Emoji string `json:"emoji"`
}
type ReactionTypeCustomEmoji struct {
	ReactionType
	CustomEmojiID string `json:"custom_emoji_id"`
}
type ReactionTypePaid struct {
	ReactionType
}
type ReactionCount struct {
	Type       *ReactionType `json:"type"`
	TotalCount int           `json:"total_count"`
}

type File struct {
	FileId       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

type ChatActions string

const (
	ChatActionTyping          ChatActions = "typing"
	ChatActionUploadPhoto     ChatActions = "upload_photo"
	ChatActionUploadVideo     ChatActions = "upload_video"
	ChatActionUploadVoice     ChatActions = "upload_voice"
	ChatActionUploadDocument  ChatActions = "upload_document"
	ChatActionChooseSticker   ChatActions = "choose_sticker"
	ChatActionFindLocation    ChatActions = "find_location"
	ChatActionUploadVideoNone ChatActions = "upload_video_none"
)
