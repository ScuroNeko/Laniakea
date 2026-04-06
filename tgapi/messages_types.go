package tgapi

import (
	"encoding/json"

	"git.scuroneko.dev/scuroneko/extypes"
)

// MessageID represents a message identifier wrapper returned by some API methods.
type MessageID struct {
	MessageID int `json:"message_id"`
}

// DirectMessageTopic represents a forum topic in a direct message.
type DirectMessageTopic struct {
	TopicID int64 `json:"topic_id"`
	User    *User `json:"user,omitempty"`
}

type MessageOriginType string

const (
	MessageOriginUserType       = "user"
	MessageOriginHiddenUserType = "hidden_user"
	MessageOriginChatType       = "chat"
	MessageOriginChannel        = "channel"
)

type MessageOrigin struct {
	Type MessageOriginType `json:"type"`
	Date int64             `json:"date"`

	SenderUser *User `json:"sender_user,omitempty"`

	SenderUserName string `json:"sender_user_name,omitempty"`

	SenderChat *Chat `json:"sender_chat,omitempty"`

	Chat      *Chat `json:"chat,omitempty"`
	MessageID int   `json:"message_id"`

	AuthorSignature string `json:"author_signature,omitempty"`
}

type ExternalReplyInfo struct {
	Origin             MessageOrigin       `json:"origin"`
	Chat               *Chat               `json:"chat,omitempty"`
	MessageID          int                 `json:"message_id,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	Animation          *Animation          `json:"animation,omitempty"`
	Audio              *Audio              `json:"audio,omitempty"`
	Document           *Document           `json:"document,omitempty"`
	PaidMedia          *PaidMediaInfo      `json:"paid_media,omitempty"`
	Photo              []PhotoSize         `json:"photo,omitempty"`
	Sticker            *Sticker            `json:"sticker,omitempty"`
	Story              *Story              `json:"story,omitempty"`
	Video              *Video              `json:"video,omitempty"`
	VideoNote          *VideoNote          `json:"video_note,omitempty"`
	Voice              *Voice              `json:"voice,omitempty"`
	HasMediaSpoiler    bool                `json:"has_media_spoiler,omitempty"`
	Checklist          *Checklist          `json:"checklist,omitempty"`
	Contact            *Contact            `json:"contact,omitempty"`
	Dice               *Dice               `json:"dice,omitempty"`
	Game               *Game               `json:"game,omitempty"`
	Giveaway           *Giveaway           `json:"giveaway,omitempty"`
	GiveawayWinners    *GiveawayWinners    `json:"giveaway_winners,omitempty"`
	Invoice            *Invoice            `json:"invoice,omitempty"`
	Location           *Location           `json:"location,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	Venue              *Venue              `json:"venue,omitempty"`
}

type TextQuote struct {
	Text     string          `json:"text"`
	Entities []MessageEntity `json:"entities"`
	Position int             `json:"position"`
	IsManual bool            `json:"is_manual,omitempty"`
}

type MessageAutoDeleteTimerChanged struct {
	MessageAutoDeleteTime int `json:"message_auto_delete_time"`
}

type DirectMessagePriceChanged struct {
	AreDirectMessagesEnabled bool `json:"are_direct_messages_enabled"`
	DirectMessageStarCount   int  `json:"direct_message_star_count,omitempty"`
}

type PaidMessagePriceChanged struct {
	PaidMessageStarCount int `json:"paid_message_star_count"`
}

// Message represents a Telegram message.
// See https://core.telegram.org/bots/api#message
type Message struct {
	MessageID          int                 `json:"message_id"`
	MessageThreadID    int                 `json:"message_thread_id,omitempty"`
	DirectMessageTopic *DirectMessageTopic `json:"direct_message_topic,omitempty"`
	From               *User               `json:"from,omitempty"`

	SenderChat           *Chat          `json:"sender_chat,omitempty"`
	SenderBoostCount     int            `json:"sender_boost_count,omitempty"`
	SenderBusinessBot    *User          `json:"sender_business_bot,omitempty"`
	SenderTag            string         `json:"sender_tag,omitempty"`
	Date                 int            `json:"date"`
	BusinessConnectionId string         `json:"business_connection_id,omitempty"`
	Chat                 *Chat          `json:"chat,omitempty"`
	ForwardOrigin        *MessageOrigin `json:"forward_origin,omitempty"`

	IsTopicMessage     bool               `json:"is_topic_message,omitempty"`
	IsAutomaticForward bool               `json:"is_automatic_forward,omitempty"`
	ReplyToMessage     *Message           `json:"reply_to_message,omitempty"`
	ExternalReply      *ExternalReplyInfo `json:"external_reply,omitempty"`
	Quote              *TextQuote         `json:"quote,omitempty"`

	ReplyToStory           *Story `json:"reply_to_story,omitempty"`
	ReplyToChecklistTaskID int    `json:"reply_to_checklist_task_id,omitempty"`
	ReplyToPollOptionID    string `json:"reply_to_poll_option_id,omitempty"`
	ViaBot                 *User  `json:"via_bot,omitempty"`
	EditDate               int    `json:"edit_date,omitempty"`
	HasProtectedContent    bool   `json:"has_protected_content,omitempty"`
	IsFromOffline          bool   `json:"is_from_offline,omitempty"`
	IsPaidPost             bool   `json:"is_paid_post,omitempty"`
	MediaGroupId           string `json:"media_group_id,omitempty"`
	AuthorSignature        string `json:"author_signature,omitempty"`
	PaidStarCount          int    `json:"paid_star_count,omitempty"`

	Text               string              `json:"text"`
	Entities           []MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	SuggestedPostInfo  *SuggestedPostInfo  `json:"suggested_post_info,omitempty"`
	EffectID           string              `json:"effect_id,omitempty"`

	Animation             *Animation               `json:"animation,omitempty"`
	Audio                 *Audio                   `json:"audio,omitempty"`
	Document              *Document                `json:"document,omitempty"`
	PaidMedia             *PaidMediaInfo           `json:"paid_media,omitempty"`
	Photo                 extypes.Slice[PhotoSize] `json:"photo,omitempty"`
	Sticker               *Sticker                 `json:"sticker,omitempty"`
	Story                 *Story                   `json:"story,omitempty"`
	Video                 *Video                   `json:"video,omitempty"`
	VideoNote             *VideoNote               `json:"video_note,omitempty"`
	Voice                 *Voice                   `json:"voice,omitempty"`
	Caption               string                   `json:"caption,omitempty"`
	CaptionEntities       []MessageEntity          `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                     `json:"show_caption_above_media,omitempty"`
	HasMediaSpoiler       bool                     `json:"has_media_spoiler,omitempty"`
	Checklist             *Checklist               `json:"checklist,omitempty"`
	Contact               *Contact                 `json:"contact,omitempty"`
	Dice                  *Dice                    `json:"dice,omitempty"`
	Game                  *Game                    `json:"game,omitempty"`
	Poll                  *Poll                    `json:"poll,omitempty"`
	Venue                 *Venue                   `json:"venue,omitempty"`
	Location              *Location                `json:"location,omitempty"`

	NewChatMembers                []User                         `json:"new_chat_members,omitempty"`
	LeftChatMember                *User                          `json:"left_chat_member,omitempty"`
	ChatOwnerLeft                 *ChatOwnerLeft                 `json:"chat_owner_left,omitempty"`
	ChatOwnerChanged              *ChatOwnerChanged              `json:"chat_owner_changed,omitempty"`
	NewChatTitle                  string                         `json:"new_chat_title,omitempty"`
	NewChatPhoto                  []PhotoSize                    `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto               bool                           `json:"delete_chat_photo,omitempty"`
	GroupChatCreated              bool                           `json:"group_chat_created,omitempty"`
	SupergroupChatCreated         bool                           `json:"supergroup_chat_created,omitempty"`
	ChannelChatCreated            bool                           `json:"channel_chat_created,omitempty"`
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	MigrateToChatID               int64                          `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID             int64                          `json:"migrate_from_chat_id,omitempty"`
	PinnedMessage                 *MaybeInaccessibleMessage      `json:"pinned_message,omitempty"`

	Invoice           *Invoice           `json:"invoice,omitempty"`
	SuccessfulPayment *SuccessfulPayment `json:"successful_payment,omitempty"`
	RefundedPayment   *RefundedPayment   `json:"refunded_payment,omitempty"`
	UsersShared       *UsersShared       `json:"users_shared,omitempty"`
	ChatShared        *ChatShared        `json:"chat_shared,omitempty"`
	Gift              *GiftInfo          `json:"gift,omitempty"`
	UniqueGift        *UniqueGiftInfo    `json:"unique_gift,omitempty"`
	GiftUpgradeSent   *GiftInfo          `json:"gift_upgrade_sent,omitempty"`

	ConnectedWebsite        string                   `json:"connected_website,omitempty"`
	WriteAccessAllowed      *WriteAccessAllowed      `json:"write_access_allowed,omitempty"`
	PassportData            *PassportData            `json:"passport_data,omitempty"`
	ProximityAlertTriggered *ProximityAlertTriggered `json:"proximity_alert_triggered,omitempty"`
	BoostAdded              *ChatBoostAdded          `json:"boost_added,omitempty"`
	ChatBackgroundSet       *ChatBackground          `json:"chat_background_set,omitempty"`

	ChecklistTaskDone         *ChecklistTaskDone         `json:"checklist_task_done,omitempty"`
	ChecklistTasksAdded       *ChecklistTasksAdded       `json:"checklist_tasks_added,omitempty"`
	DirectMessagePriceChanged *DirectMessagePriceChanged `json:"direct_message_price_changed,omitempty"`
	ForumTopicCreated         *ForumTopicCreated         `json:"forum_topic_created,omitempty"`
	ForumTopicEdited          *ForumTopicEdited          `json:"forum_topic_edited,omitempty"`
	ForumTopicClosed          *ForumTopicClosed          `json:"forum_topic_closed,omitempty"`
	ForumTopicReopened        *ForumTopicReopened        `json:"forum_topic_reopened,omitempty"`
	GeneralForumTopicHidden   *GeneralForumTopicHidden   `json:"general_forum_topic_hidden,omitempty"`
	GeneralForumTopicUnhidden *GeneralForumTopicUnhidden `json:"general_forum_topic_unhidden,omitempty"`

	GiveawayCreated   *GiveawayCreated   `json:"giveaway_created,omitempty"`
	Giveaway          *Giveaway          `json:"giveaway,omitempty"`
	GiveawayWinners   *GiveawayWinners   `json:"giveaway_winners,omitempty"`
	GiveawayCompleted *GiveawayCompleted `json:"giveaway_completed,omitempty"`

	ManagedBotCreated       *ManagedBotCreated       `json:"managed_bot_created,omitempty"`
	PaidMessagePriceChanged *PaidMessagePriceChanged `json:"paid_message_price_changed,omitempty"`
	PollOptionAdded         *PollOptionAdded         `json:"poll_option_added,omitempty"`
	PollOptionDeleted       *PollOptionDeleted       `json:"poll_option_deleted,omitempty"`

	SuggestedPostApproved       *SuggestedPostApproved       `json:"suggested_post_approved,omitempty"`
	SuggestedPostApprovalFailed *SuggestedPostApprovalFailed `json:"suggested_post_approval_failed,omitempty"`
	SuggestedPostDeclined       *SuggestedPostDeclined       `json:"suggested_post_declined,omitempty"`
	SuggestedPostPaid           *SuggestedPostPaid           `json:"suggested_post_paid,omitempty"`
	SuggestedPostRefunded       *SuggestedPostRefunded       `json:"suggested_post_refunded,omitempty"`

	VideoChatScheduled           *VideoChatScheduled           `json:"video_chat_scheduled,omitempty"`
	VideoChatStarted             *VideoChatStarted             `json:"video_chat_started,omitempty"`
	VideoChatEnded               *VideoChatEnded               `json:"video_chat_ended,omitempty"`
	VideoChatParticipantsInvited *VideoChatParticipantsInvited `json:"video_chat_participants_invited,omitempty"`

	WebAppData  *WebAppData           `json:"web_app_data,omitempty"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
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
type MaybeInaccessibleMessage struct {
	msg *Message
	ina *InaccessibleMessage
}

// UnmarshalJSON decodes either an accessible Message or an InaccessibleMessage.
func (m *MaybeInaccessibleMessage) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Date int `json:"date"`
	}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	var err error
	if tmp.Date > 0 {
		err = json.Unmarshal(data, &m.msg)
	} else {
		err = json.Unmarshal(data, &m.ina)
	}
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON encodes the populated accessible or inaccessible message payload.
func (m *MaybeInaccessibleMessage) MarshalJSON() ([]byte, error) {
	if m.msg != nil {
		return json.Marshal(m.msg)
	} else if m.ina != nil {
		return json.Marshal(m.ina)
	}
	return json.Marshal(nil)
}

// Message returns the accessible message payload when present.
func (m *MaybeInaccessibleMessage) Message() *Message {
	return m.msg
}

// InaccessibleMessage returns the inaccessible message payload when present.
func (m *MaybeInaccessibleMessage) InaccessibleMessage() *InaccessibleMessage {
	return m.ina
}

// IsAccessible reports whether the payload is an accessible message.
func (m *MaybeInaccessibleMessage) IsAccessible() bool {
	return m.msg != nil
}

// IsInaccessible reports whether the payload is an inaccessible message.
func (m *MaybeInaccessibleMessage) IsInaccessible() bool {
	return m.ina != nil
}

// MessageID returns the message identifier from either payload form.
func (m *MaybeInaccessibleMessage) MessageID() int {
	if m.IsAccessible() {
		return m.msg.MessageID
	} else if m.IsInaccessible() {
		return m.ina.MessageID
	}
	return 0
}

// Chat returns the chat from either payload form.
func (m *MaybeInaccessibleMessage) Chat() *Chat {
	if m.IsAccessible() {
		return m.msg.Chat
	} else if m.IsInaccessible() {
		return &m.ina.Chat
	}
	return nil
}

// MessageEntityType represents the type of a message entity.
type MessageEntityType string

const (
	// MessageEntityMention identifies an @mention entity.
	MessageEntityMention MessageEntityType = "mention"
	// MessageEntityHashtag identifies a hashtag entity.
	MessageEntityHashtag MessageEntityType = "hashtag"
	// MessageEntityCashtag identifies a cashtag entity.
	MessageEntityCashtag MessageEntityType = "cashtag"
	// MessageEntityBotCommand identifies a bot command entity.
	MessageEntityBotCommand MessageEntityType = "bot_command"
	// MessageEntityUrl identifies a URL entity.
	MessageEntityUrl MessageEntityType = "url"
	// MessageEntityEmail identifies an email entity.
	MessageEntityEmail MessageEntityType = "email"
	// MessageEntityPhoneNumber identifies a phone number entity.
	MessageEntityPhoneNumber MessageEntityType = "phone_number"
	// MessageEntityBold identifies bold text.
	MessageEntityBold MessageEntityType = "bold"
	// MessageEntityItalic identifies italic text.
	MessageEntityItalic MessageEntityType = "italic"
	// MessageEntityUnderline identifies underlined text.
	MessageEntityUnderline MessageEntityType = "underline"
	// MessageEntityStrike identifies strikethrough text.
	MessageEntityStrike MessageEntityType = "strikethrough"
	// MessageEntitySpoiler identifies spoiler text.
	MessageEntitySpoiler MessageEntityType = "spoiler"
	// MessageEntityBlockquote identifies a blockquote entity.
	MessageEntityBlockquote MessageEntityType = "blockquote"
	// MessageEntityExpandableBlockquote identifies an expandable blockquote entity.
	MessageEntityExpandableBlockquote MessageEntityType = "expandable_blockquote"
	// MessageEntityCode identifies inline code.
	MessageEntityCode MessageEntityType = "code"
	// MessageEntityPre identifies a preformatted block.
	MessageEntityPre MessageEntityType = "pre"
	// MessageEntityTextLink identifies linked text.
	MessageEntityTextLink MessageEntityType = "text_link"
	// MessageEntityTextMention identifies a text mention.
	MessageEntityTextMention MessageEntityType = "text_mention"
	// MessageEntityCustomEmoji identifies a custom emoji entity.
	MessageEntityCustomEmoji MessageEntityType = "custom_emoji"
	// MessageEntityDateTime identifies a date-time entity.
	MessageEntityDateTime MessageEntityType = "date_time"
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

	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply,omitempty"`
	Quote                    string          `json:"quote,omitempty"`
	QuoteParsingMode         string          `json:"quote_parsing_mode,omitempty"`
	QuoteEntities            []MessageEntity `json:"quote_entities,omitempty"`
	QuotePosition            int             `json:"quote_position,omitempty"`
	ChecklistTaskID          int             `json:"checklist_task_id,omitempty"`
	PollOptionID             string          `json:"poll_option_id,omitempty"`
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

	Keyboard              [][]KeyboardButton `json:"keyboard,omitempty"`
	IsPersistent          bool               `json:"is_persistent,omitempty"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string             `json:"input_field_placeholder,omitempty"`
	Selective             bool               `json:"selective,omitempty"`

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

const (
	// KeyboardButtonStyleDanger marks a destructive keyboard button.
	KeyboardButtonStyleDanger KeyboardButtonStyle = "danger"
	// KeyboardButtonStyleSuccess marks a confirmatory keyboard button.
	KeyboardButtonStyleSuccess KeyboardButtonStyle = "success"
	// KeyboardButtonStylePrimary marks a primary keyboard button.
	KeyboardButtonStylePrimary KeyboardButtonStyle = "primary"
)

// KeyboardButton represents one button of the reply keyboard.
// See https://core.telegram.org/bots/api#keyboardbutton
type KeyboardButton struct {
	Text              string                           `json:"text"`
	IconCustomEmojiID string                           `json:"icon_custom_emoji_id,omitempty"`
	Style             KeyboardButtonStyle              `json:"style,omitempty"`
	RequestUsers      *KeyboardButtonRequestUsers      `json:"request_users,omitempty"`
	RequestChat       *KeyboardButtonRequestChat       `json:"request_chat,omitempty"`
	RequestManagedBot *KeyboardButtonRequestManagedBot `json:"request_managed_bot,omitempty"`
	RequestContact    bool                             `json:"request_contact,omitempty"`
	RequestLocation   bool                             `json:"request_location,omitempty"`
	RequestPoll       *KeyboardButtonPollType          `json:"request_poll,omitempty"`
	WebApp            *WebAppInfo                      `json:"web_app,omitempty"`
}

// KeyboardButtonRequestUsers defines criteria used to request suitable users.
// See https://core.telegram.org/bots/api#keyboardbuttonrequestusers
type KeyboardButtonRequestUsers struct {
	RequestID       int   `json:"request_id"`
	UserIsBot       *bool `json:"user_is_bot,omitempty"`
	UserIsPremium   *bool `json:"user_is_premium,omitempty"`
	MaxQuantity     int   `json:"max_quantity,omitempty"`
	RequestName     bool  `json:"request_name,omitempty"`
	RequestUsername bool  `json:"request_username,omitempty"`
	RequestPhoto    bool  `json:"request_photo,omitempty"`
}

// KeyboardButtonRequestChat defines criteria used to request a suitable chat.
// See https://core.telegram.org/bots/api#keyboardbuttonrequestchat
type KeyboardButtonRequestChat struct {
	RequestID               int                      `json:"request_id"`
	ChatIsChannel           bool                     `json:"chat_is_channel"`
	ChatIsForum             *bool                    `json:"chat_is_forum,omitempty"`
	ChatHasUsername         *bool                    `json:"chat_has_username,omitempty"`
	ChatIsCreated           *bool                    `json:"chat_is_created,omitempty"`
	UserAdministratorRights *ChatAdministratorRights `json:"user_administrator_rights,omitempty"`
	BotAdministratorRights  *ChatAdministratorRights `json:"bot_administrator_rights,omitempty"`
	BotIsMember             bool                     `json:"bot_is_member,omitempty"`
	RequestTitle            bool                     `json:"request_title,omitempty"`
	RequestUsername         bool                     `json:"request_username,omitempty"`
	RequestPhoto            bool                     `json:"request_photo,omitempty"`
}

// KeyboardButtonRequestManagedBot defines criteria used to request a managed bot.
// See https://core.telegram.org/bots/api#keyboardbuttonrequestmanagedbot
type KeyboardButtonRequestManagedBot struct {
	RequestID         int32  `json:"request_id"`
	SuggestedName     string `json:"suggested_name,omitempty"`
	SuggestedUsername string `json:"suggested_username,omitempty"`
}

// KeyboardButtonPollType represents the type of a poll that may be created from a keyboard button.
// See https://core.telegram.org/bots/api#keyboardbuttonpolltype
type KeyboardButtonPollType struct {
	Type PollType `json:"type,omitempty"`
}

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
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	IsPersistent          bool               `json:"is_persistent,omitempty"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string             `json:"input_field_placeholder,omitempty"`
	Selective             bool               `json:"selective,omitempty"`
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

// ChatActionType represents the type of chat action.
type ChatActionType string

const (
	// ChatActionTyping tells Telegram the bot is typing.
	ChatActionTyping ChatActionType = "typing"
	// ChatActionUploadPhoto tells Telegram the bot is uploading a photo.
	ChatActionUploadPhoto ChatActionType = "upload_photo"
	// ChatActionUploadVideo tells Telegram the bot is uploading a video.
	ChatActionUploadVideo ChatActionType = "upload_video"
	// ChatActionUploadVoice tells Telegram the bot is uploading a voice message.
	ChatActionUploadVoice ChatActionType = "upload_voice"
	// ChatActionUploadDocument tells Telegram the bot is uploading a document.
	ChatActionUploadDocument ChatActionType = "upload_document"
	// ChatActionChooseSticker tells Telegram the bot is choosing a sticker.
	ChatActionChooseSticker ChatActionType = "choose_sticker"
	// ChatActionFindLocation tells Telegram the bot is finding a location.
	ChatActionFindLocation ChatActionType = "find_location"
	// ChatActionUploadVideoNote tells Telegram the bot is uploading a video note.
	ChatActionUploadVideoNote ChatActionType = "upload_video_note"
	// ChatActionUploadVideoNone is a deprecated alias for ChatActionUploadVideoNote.
	ChatActionUploadVideoNone ChatActionType = ChatActionUploadVideoNote
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

// ManagedBotCreated describes a service message about a newly created managed bot.
// See https://core.telegram.org/bots/api#managedbotcreated
type ManagedBotCreated struct {
	Bot User `json:"bot"`
}

// ManagedBotUpdated describes an update about a managed bot and its manager.
// See https://core.telegram.org/bots/api#managedbotupdated
type ManagedBotUpdated struct {
	User User `json:"user"`
	Bot  User `json:"bot"`
}

type SharedUser struct {
	UserID    int64       `json:"user_id"`
	FirstName string      `json:"first_name,omitempty"`
	LastName  string      `json:"last_name,omitempty"`
	Username  string      `json:"username,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
}
type UsersShared struct {
	RequestID int          `json:"request_id"`
	Users     []SharedUser `json:"users"`
}
type ChatShared struct {
	RequestID int         `json:"request_id"`
	ChatID    int64       `json:"chat_id"`
	Title     string      `json:"title,omitempty"`
	Username  string      `json:"username,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
}

type SuggestedPostApproved struct {
	SuggestedPostMessage *Message           `json:"suggested_post_message,omitempty"`
	Price                SuggestedPostPrice `json:"price"`
	SendDate             int                `json:"send_date"`
}
type SuggestedPostApprovalFailed struct {
	SuggestedPostMessage *Message           `json:"suggested_post_message,omitempty"`
	Price                SuggestedPostPrice `json:"price"`
}
type SuggestedPostDeclined struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`
	Comment              string   `json:"comment,omitempty"`
}
type SuggestedPostPaid struct {
	SuggestedPostMessage *Message    `json:"suggested_post_message,omitempty"`
	Currency             string      `json:"currency"`
	Amount               int         `json:"amount"`
	StarAmount           *StarAmount `json:"star_amount,omitempty"`
}
type SuggestedPostRefunded struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`
	Reason               string   `json:"reason,omitempty"`
}

type VideoChatScheduled struct {
	StartDate int64 `json:"start_date"`
}
type VideoChatStarted struct{}
type VideoChatEnded struct {
	Duration int64 `json:"duration"`
}
type VideoChatParticipantsInvited struct {
	Users []User `json:"users"`
}
