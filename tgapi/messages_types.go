package tgapi

import (
	"encoding/json"

	"git.scuroneko.dev/scuroneko/extypes"
)

// MessageID represents a message identifier wrapper returned by some API methods.
// Since: Bot API 7.0
type MessageID struct {
	MessageID int `json:"message_id"`
}

// DirectMessageTopic represents a forum topic in a direct message.
// Since: Bot API 9.2
type DirectMessageTopic struct {
	TopicID int64 `json:"topic_id"`
	User    *User `json:"user,omitempty"`
}

// MessageOriginType represents the type of a message origin.
type MessageOriginType string

const (
	MessageOriginUserType       = "user"
	MessageOriginHiddenUserType = "hidden_user"
	MessageOriginChatType       = "chat"
	MessageOriginChannel        = "channel"
)

// MessageOrigin describes the origin of a message.
// Since: Bot API 7.0
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

// ExternalReplyInfo contains information about a message that is being replied to.
// Since: Bot API 7.0
type ExternalReplyInfo struct {
	Origin             MessageOrigin       `json:"origin"`
	Chat               *Chat               `json:"chat,omitempty"`
	MessageID          int                 `json:"message_id,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	Animation          *Animation          `json:"animation,omitempty"`
	Audio              *Audio              `json:"audio,omitempty"`
	Document           *Document           `json:"document,omitempty"`
	PaidMedia          *PaidMediaInfo      `json:"paid_media,omitempty"` // Since: Bot API 7.6
	Photo              []PhotoSize         `json:"photo,omitempty"`
	LivePhoto          *LivePhoto          `json:"live_photo,omitempty"` // Since: Bot API 10.0
	Sticker            *Sticker            `json:"sticker,omitempty"`
	Story              *Story              `json:"story,omitempty"`
	Video              *Video              `json:"video,omitempty"`
	VideoNote          *VideoNote          `json:"video_note,omitempty"`
	Voice              *Voice              `json:"voice,omitempty"`
	HasMediaSpoiler    bool                `json:"has_media_spoiler,omitempty"`
	Checklist          *Checklist          `json:"checklist,omitempty"` // Since: Bot API 9.1
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

// TextQuote contains information about the quoted part of a message.
// Since: Bot API 7.0
type TextQuote struct {
	Text     string          `json:"text"`
	Entities []MessageEntity `json:"entities"`
	Position int             `json:"position"`
	IsManual bool            `json:"is_manual,omitempty"`
}

// MessageAutoDeleteTimerChanged represents a service message about a change in auto-delete timer settings.
// Since: Bot API 5.1
type MessageAutoDeleteTimerChanged struct {
	MessageAutoDeleteTime int `json:"message_auto_delete_time"`
}

// DirectMessagePriceChanged represents a service message about a change in the price of direct messages.
// Since: Bot API 9.1
type DirectMessagePriceChanged struct {
	AreDirectMessagesEnabled bool `json:"are_direct_messages_enabled"`
	DirectMessageStarCount   int  `json:"direct_message_star_count,omitempty"`
}

// PaidMessagePriceChanged represents a service message about a change in the price of paid messages.
// Since: Bot API 9.x
type PaidMessagePriceChanged struct {
	PaidMessageStarCount int `json:"paid_message_star_count"`
}

// Message represents a Telegram message.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#message
type Message struct {
	MessageID          int                 `json:"message_id"`
	MessageThreadID    int                 `json:"message_thread_id,omitempty"`    // Since: Bot API 6.3
	DirectMessageTopic *DirectMessageTopic `json:"direct_message_topic,omitempty"` // Since: Bot API 9.2
	From               *User               `json:"from,omitempty"`

	SenderChat           *Chat          `json:"sender_chat,omitempty"`         // Since: Bot API 5.0
	SenderBoostCount     int            `json:"sender_boost_count,omitempty"`  // Since: Bot API 7.1
	SenderBusinessBot    *User          `json:"sender_business_bot,omitempty"` // Since: Bot API 7.2
	SenderTag            string         `json:"sender_tag,omitempty"`          // Since: Bot API 9.5
	Date                 int            `json:"date"`
	GuestQueryID         string         `json:"guest_query_id,omitempty"`         // Since: Bot API 10.0
	BusinessConnectionID string         `json:"business_connection_id,omitempty"` // Since: Bot API 7.2
	Chat                 *Chat          `json:"chat,omitempty"`
	ForwardOrigin        *MessageOrigin `json:"forward_origin,omitempty"` // Since: Bot API 7.0

	IsTopicMessage     bool               `json:"is_topic_message,omitempty"`     // Since: Bot API 6.3
	IsAutomaticForward bool               `json:"is_automatic_forward,omitempty"` // Since: Bot API 5.5
	ReplyToMessage     *Message           `json:"reply_to_message,omitempty"`
	ExternalReply      *ExternalReplyInfo `json:"external_reply,omitempty"` // Since: Bot API 7.0
	Quote              *TextQuote         `json:"quote,omitempty"`          // Since: Bot API 7.0

	ReplyToStory           *Story `json:"reply_to_story,omitempty"`             // Since: Bot API 7.1
	ReplyToChecklistTaskID int    `json:"reply_to_checklist_task_id,omitempty"` // Since: Bot API 9.1
	ReplyToPollOptionID    string `json:"reply_to_poll_option_id,omitempty"`    // Since: Bot API 9.6
	ViaBot                 *User  `json:"via_bot,omitempty"`
	GuestBotCallerUser     *User  `json:"guest_bot_caller_user,omitempty"` // Since: Bot API 10.0
	GuestBotCallerChat     *Chat  `json:"guest_bot_caller_chat,omitempty"` // Since: Bot API 10.0
	EditDate               int    `json:"edit_date,omitempty"`             // Since: Bot API 2.1
	HasProtectedContent    bool   `json:"has_protected_content,omitempty"` // Since: Bot API 5.5
	IsFromOffline          bool   `json:"is_from_offline,omitempty"`       // Since: Bot API 7.2
	IsPaidPost             bool   `json:"is_paid_post,omitempty"`          // Since: Bot API 9.1
	MediaGroupID           string `json:"media_group_id,omitempty"`        // Since: Bot API 3.5
	AuthorSignature        string `json:"author_signature,omitempty"`
	PaidStarCount          int    `json:"paid_star_count,omitempty"` // Since: Bot API 8.3

	Text               string              `json:"text"`
	Entities           []MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	SuggestedPostInfo  *SuggestedPostInfo  `json:"suggested_post_info,omitempty"` // Since: Bot API 9.1
	EffectID           string              `json:"effect_id,omitempty"`           // Since: Bot API 7.4

	Animation             *Animation               `json:"animation,omitempty"` // Since: Bot API 4.0
	Audio                 *Audio                   `json:"audio,omitempty"`
	Document              *Document                `json:"document,omitempty"`
	PaidMedia             *PaidMediaInfo           `json:"paid_media,omitempty"` // Since: Bot API 7.6
	Photo                 extypes.Slice[PhotoSize] `json:"photo,omitempty"`
	LivePhoto             *LivePhoto               `json:"live_photo,omitempty"` // Since: Bot API 10.0
	Sticker               *Sticker                 `json:"sticker,omitempty"`
	Story                 *Story                   `json:"story,omitempty"`
	Video                 *Video                   `json:"video,omitempty"`
	VideoNote             *VideoNote               `json:"video_note,omitempty"`               // Since: Bot API 3.0
	Voice                 *Voice                   `json:"voice,omitempty"`                    // Since: Bot API 1.2
	Caption               string                   `json:"caption,omitempty"`                  // Since: Bot API 3.4
	CaptionEntities       []MessageEntity          `json:"caption_entities,omitempty"`         // Since: Bot API 3.4
	ShowCaptionAboveMedia bool                     `json:"show_caption_above_media,omitempty"` // Since: Bot API 7.4
	HasMediaSpoiler       bool                     `json:"has_media_spoiler,omitempty"`        // Since: Bot API 6.4
	Checklist             *Checklist               `json:"checklist,omitempty"`                // Since: Bot API 9.1
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
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"` // Since: Bot API 5.1
	MigrateToChatID               int64                          `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID             int64                          `json:"migrate_from_chat_id,omitempty"`
	PinnedMessage                 *MaybeInaccessibleMessage      `json:"pinned_message,omitempty"`

	Invoice           *Invoice           `json:"invoice,omitempty"`            // Since: Bot API 3.0
	SuccessfulPayment *SuccessfulPayment `json:"successful_payment,omitempty"` // Since: Bot API 3.0
	RefundedPayment   *RefundedPayment   `json:"refunded_payment,omitempty"`   // Since: Bot API 7.7
	UsersShared       *UsersShared       `json:"users_shared,omitempty"`       // Since: Bot API 6.5
	ChatShared        *ChatShared        `json:"chat_shared,omitempty"`        // Since: Bot API 6.5
	Gift              *GiftInfo          `json:"gift,omitempty"`               // Since: Bot API 9.0
	UniqueGift        *UniqueGiftInfo    `json:"unique_gift,omitempty"`        // Since: Bot API 9.0
	GiftUpgradeSent   *GiftInfo          `json:"gift_upgrade_sent,omitempty"`  // Since: Bot API 9.3

	ConnectedWebsite        string                   `json:"connected_website,omitempty"`
	WriteAccessAllowed      *WriteAccessAllowed      `json:"write_access_allowed,omitempty"` // Since: Bot API 6.4
	PassportData            *PassportData            `json:"passport_data,omitempty"`
	ProximityAlertTriggered *ProximityAlertTriggered `json:"proximity_alert_triggered,omitempty"` // Since: Bot API 5.0
	BoostAdded              *ChatBoostAdded          `json:"boost_added,omitempty"`               // Since: Bot API 7.1
	ChatBackgroundSet       *ChatBackground          `json:"chat_background_set,omitempty"`       // Since: Bot API 7.5

	ChecklistTaskDone         *ChecklistTaskDone         `json:"checklist_task_done,omitempty"`          // Since: Bot API 9.1
	ChecklistTasksAdded       *ChecklistTasksAdded       `json:"checklist_tasks_added,omitempty"`        // Since: Bot API 9.1
	DirectMessagePriceChanged *DirectMessagePriceChanged `json:"direct_message_price_changed,omitempty"` // Since: Bot API 9.1
	PaidMessagePriceChanged   *PaidMessagePriceChanged   `json:"paid_message_price_changed,omitempty"`   // Since: Bot API 9.x
	ForumTopicCreated         *ForumTopicCreated         `json:"forum_topic_created,omitempty"`          // Since: Bot API 6.3
	ForumTopicEdited          *ForumTopicEdited          `json:"forum_topic_edited,omitempty"`           // Since: Bot API 6.4
	ForumTopicClosed          *ForumTopicClosed          `json:"forum_topic_closed,omitempty"`           // Since: Bot API 6.3
	ForumTopicReopened        *ForumTopicReopened        `json:"forum_topic_reopened,omitempty"`         // Since: Bot API 6.3
	GeneralForumTopicHidden   *GeneralForumTopicHidden   `json:"general_forum_topic_hidden,omitempty"`   // Since: Bot API 6.4
	GeneralForumTopicUnhidden *GeneralForumTopicUnhidden `json:"general_forum_topic_unhidden,omitempty"` // Since: Bot API 6.4

	GiveawayCreated   *GiveawayCreated   `json:"giveaway_created,omitempty"`   // Since: Bot API 7.0
	Giveaway          *Giveaway          `json:"giveaway,omitempty"`           // Since: Bot API 7.0
	GiveawayWinners   *GiveawayWinners   `json:"giveaway_winners,omitempty"`   // Since: Bot API 7.0
	GiveawayCompleted *GiveawayCompleted `json:"giveaway_completed,omitempty"` // Since: Bot API 7.0

	ManagedBotCreated *ManagedBotCreated `json:"managed_bot_created,omitempty"` // Since: Bot API 9.6
	PollOptionAdded   *PollOptionAdded   `json:"poll_option_added,omitempty"`   // Since: Bot API 9.6
	PollOptionDeleted *PollOptionDeleted `json:"poll_option_deleted,omitempty"` // Since: Bot API 9.6

	SuggestedPostApproved       *SuggestedPostApproved       `json:"suggested_post_approved,omitempty"`        // Since: Bot API 9.1
	SuggestedPostApprovalFailed *SuggestedPostApprovalFailed `json:"suggested_post_approval_failed,omitempty"` // Since: Bot API 9.1
	SuggestedPostDeclined       *SuggestedPostDeclined       `json:"suggested_post_declined,omitempty"`        // Since: Bot API 9.1
	SuggestedPostPaid           *SuggestedPostPaid           `json:"suggested_post_paid,omitempty"`            // Since: Bot API 9.1
	SuggestedPostRefunded       *SuggestedPostRefunded       `json:"suggested_post_refunded,omitempty"`        // Since: Bot API 9.1

	VideoChatScheduled           *VideoChatScheduled           `json:"video_chat_scheduled,omitempty"`            // Since: Bot API 6.0
	VideoChatStarted             *VideoChatStarted             `json:"video_chat_started,omitempty"`              // Since: Bot API 5.1
	VideoChatEnded               *VideoChatEnded               `json:"video_chat_ended,omitempty"`                // Since: Bot API 5.1
	VideoChatParticipantsInvited *VideoChatParticipantsInvited `json:"video_chat_participants_invited,omitempty"` // Since: Bot API 5.1

	WebAppData  *WebAppData           `json:"web_app_data,omitempty"` // Since: Bot API 6.0
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"` // Since: Bot API 4.3
}

// InaccessibleMessage describes a message that was deleted or is otherwise inaccessible.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#inaccessiblemessage
type InaccessibleMessage struct {
	Chat      Chat `json:"chat"`
	MessageID int  `json:"message_id"`
	Date      int  `json:"date"`
}

// MaybeInaccessibleMessage is a union type that can be either Message or InaccessibleMessage.
// Since: Bot API 7.0
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
	// MessageEntityURL identifies a URL entity.
	MessageEntityURL MessageEntityType = "url"
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
	MessageEntitySpoiler MessageEntityType = "spoiler" // Since: Bot API 5.6
	// MessageEntityBlockquote identifies a blockquote entity.
	MessageEntityBlockquote MessageEntityType = "blockquote"
	// MessageEntityExpandableBlockquote identifies an expandable blockquote entity.
	MessageEntityExpandableBlockquote MessageEntityType = "expandable_blockquote" // Since: Bot API 7.5
	// MessageEntityCode identifies inline code.
	MessageEntityCode MessageEntityType = "code"
	// MessageEntityPre identifies a preformatted block.
	MessageEntityPre MessageEntityType = "pre"
	// MessageEntityTextLink identifies linked text.
	MessageEntityTextLink MessageEntityType = "text_link"
	// MessageEntityTextMention identifies a text mention.
	MessageEntityTextMention MessageEntityType = "text_mention"
	// MessageEntityCustomEmoji identifies a custom emoji entity.
	MessageEntityCustomEmoji MessageEntityType = "custom_emoji" // Since: Bot API 6.2
	// MessageEntityDateTime identifies a date-time entity.
	MessageEntityDateTime MessageEntityType = "date_time" // Since: Bot API 9.5
)

// MessageEntity represents one special entity in a text message.
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#messageentity
type MessageEntity struct {
	Type MessageEntityType `json:"type"`

	Offset        int    `json:"offset"`
	Length        int    `json:"length"`
	URL           string `json:"url,omitempty"`
	User          *User  `json:"user,omitempty"`
	Language      string `json:"language,omitempty"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"` // Since: Bot API 6.2

	UnixTime       int64  `json:"unix_time,omitempty"`
	DateTimeFormat string `json:"date_time_format,omitempty"` // Since: Bot API 9.5
}

// ReplyParameters describes the parameters to use when replying to a message.
// Since: Bot API 7.0
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
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#linkpreviewoptions
type LinkPreviewOptions struct {
	IsDisabled       bool   `json:"is_disabled,omitempty"`
	URL              string `json:"url,omitempty"`
	PreferSmallMedia bool   `json:"prefer_small_media,omitempty"`
	PreferLargeMedia bool   `json:"prefer_large_media,omitempty"`
	ShowAboveText    bool   `json:"show_above_text,omitempty"`
}

// ReplyMarkup represents a custom keyboard or inline keyboard.
// Since: Bot API 1.0
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
// Since: Bot API 2.0
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
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#keyboardbutton
type KeyboardButton struct {
	Text              string                           `json:"text"`
	IconCustomEmojiID string                           `json:"icon_custom_emoji_id,omitempty"` // Since: Bot API 9.4
	Style             KeyboardButtonStyle              `json:"style,omitempty"`                // Since: Bot API 9.4
	RequestUsers      *KeyboardButtonRequestUsers      `json:"request_users,omitempty"`        // Since: Bot API 7.0
	RequestChat       *KeyboardButtonRequestChat       `json:"request_chat,omitempty"`         // Since: Bot API 6.5
	RequestManagedBot *KeyboardButtonRequestManagedBot `json:"request_managed_bot,omitempty"`  // Since: Bot API 9.6
	RequestContact    bool                             `json:"request_contact,omitempty"`
	RequestLocation   bool                             `json:"request_location,omitempty"`
	RequestPoll       *KeyboardButtonPollType          `json:"request_poll,omitempty"` // Since: Bot API 4.6
	WebApp            *WebAppInfo                      `json:"web_app,omitempty"`      // Since: Bot API 6.0
}

// KeyboardButtonRequestUsers defines criteria used to request suitable users.
// Since: Bot API 7.0
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
// Since: Bot API 6.5
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
// Since: Bot API 9.6
// See https://core.telegram.org/bots/api#keyboardbuttonrequestmanagedbot
type KeyboardButtonRequestManagedBot struct {
	RequestID         int32  `json:"request_id"`
	SuggestedName     string `json:"suggested_name,omitempty"`
	SuggestedUsername string `json:"suggested_username,omitempty"`
}

// KeyboardButtonPollType represents the type of a poll that may be created from a keyboard button.
// Since: Bot API 4.6
// See https://core.telegram.org/bots/api#keyboardbuttonpolltype
type KeyboardButtonPollType struct {
	Type PollType `json:"type,omitempty"`
}

// InlineKeyboardButton represents one button of an inline keyboard.
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#inlinekeyboardbutton
type InlineKeyboardButton struct {
	Text              string              `json:"text"`
	URL               string              `json:"url,omitempty"`
	CallbackData      string              `json:"callback_data,omitempty"`
	Style             KeyboardButtonStyle `json:"style,omitempty"`                // Since: Bot API 9.4
	IconCustomEmojiID string              `json:"icon_custom_emoji_id,omitempty"` // Since: Bot API 9.4
}

// ReplyKeyboardMarkup represents a custom keyboard with reply options.
// Since: Bot API 1.0
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
// Since: Bot API 2.0
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
	ChatActionUploadVideoNone = ChatActionUploadVideoNote
)

// MessageReactionUpdated represents a change of a reaction on a message.
// Since: Bot API 7.0
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
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#messagereactioncountupdated
type MessageReactionCountUpdated struct {
	Chat      *Chat            `json:"chat"`
	MessageID int              `json:"message_id"`
	Date      int              `json:"date"`
	Reactions []*ReactionCount `json:"reactions"`
}

// ReactionType describes the type of a reaction.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#reactiontype
type ReactionType struct {
	Type string `json:"type"`
	// ReactionTypeEmoji
	Emoji *string `json:"emoji,omitempty"`
	// ReactionTypeCustomEmoji
	CustomEmojiID *string `json:"custom_emoji_id,omitempty"`
}

// ReactionCount represents a reaction added to a message along with the number of times it was added.
// Since: Bot API 7.0
// See https://core.telegram.org/bots/api#reactioncount
type ReactionCount struct {
	Type       ReactionType `json:"type"`
	TotalCount int          `json:"total_count"`
}

// SuggestedPostPrice represents the price of a suggested post.
// Since: Bot API 9.1
type SuggestedPostPrice struct {
	Currency string `json:"currency"`
	Amount   int    `json:"amount"`
}

// SuggestedPostInfo contains information about a suggested post.
// Since: Bot API 9.1
// See https://core.telegram.org/bots/api#suggestedpostinfo
type SuggestedPostInfo struct {
	State    string             `json:"state"` // "pending", "approved", or "declined"
	Price    SuggestedPostPrice `json:"price"`
	SendDate int                `json:"send_date"`
}

// SuggestedPostParameters holds parameters for suggesting a post.
// Since: Bot API 9.2
type SuggestedPostParameters struct {
	Price    SuggestedPostPrice `json:"price"`
	SendDate int                `json:"send_date"`
}

// ManagedBotCreated describes a service message about a newly created managed bot.
// Since: Bot API 9.6
// See https://core.telegram.org/bots/api#managedbotcreated
type ManagedBotCreated struct {
	Bot User `json:"bot"`
}

// ManagedBotUpdated describes an update about a managed bot and its manager.
// Since: Bot API 9.6
// See https://core.telegram.org/bots/api#managedbotupdated
type ManagedBotUpdated struct {
	User User `json:"user"`
	Bot  User `json:"bot"`
}

// SharedUser represents a user shared via a KeyboardButtonRequestUsers button.
// Since: Bot API 7.2
type SharedUser struct {
	UserID    int64       `json:"user_id"`
	FirstName string      `json:"first_name,omitempty"`
	LastName  string      `json:"last_name,omitempty"`
	Username  string      `json:"username,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
}

// UsersShared represents a service message about users shared via a KeyboardButtonRequestUsers button.
// Since: Bot API 6.5
type UsersShared struct {
	RequestID int          `json:"request_id"`
	Users     []SharedUser `json:"users"`
}

// ChatShared represents a service message about a chat shared via a KeyboardButtonRequestChat button.
// Since: Bot API 6.5
type ChatShared struct {
	RequestID int         `json:"request_id"`
	ChatID    int64       `json:"chat_id"`
	Title     string      `json:"title,omitempty"`
	Username  string      `json:"username,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
}

// SuggestedPostApproved is a service message about an approved suggested post.
// Since: Bot API 9.1
type SuggestedPostApproved struct {
	SuggestedPostMessage *Message           `json:"suggested_post_message,omitempty"`
	Price                SuggestedPostPrice `json:"price"`
	SendDate             int                `json:"send_date"`
}

// SuggestedPostApprovalFailed is a service message about a failed suggested post approval.
// Since: Bot API 9.1
type SuggestedPostApprovalFailed struct {
	SuggestedPostMessage *Message           `json:"suggested_post_message,omitempty"`
	Price                SuggestedPostPrice `json:"price"`
}

// SuggestedPostDeclined is a service message about a declined suggested post.
// Since: Bot API 9.1
type SuggestedPostDeclined struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`
	Comment              string   `json:"comment,omitempty"`
}

// SuggestedPostPaid is a service message about a paid suggested post.
// Since: Bot API 9.1
type SuggestedPostPaid struct {
	SuggestedPostMessage *Message    `json:"suggested_post_message,omitempty"`
	Currency             string      `json:"currency"`
	Amount               int         `json:"amount"`
	StarAmount           *StarAmount `json:"star_amount,omitempty"`
}

// SuggestedPostRefunded is a service message about a refunded suggested post.
// Since: Bot API 9.1
type SuggestedPostRefunded struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`
	Reason               string   `json:"reason,omitempty"`
}

// VideoChatScheduled represents a service message about a video chat scheduled in the chat.
// Since: Bot API 6.0
type VideoChatScheduled struct {
	StartDate int64 `json:"start_date"`
}

// VideoChatStarted represents a service message about a video chat started in the chat.
// Since: Bot API 5.1
type VideoChatStarted struct{}

// VideoChatEnded represents a service message about a video chat ended in the chat.
// Since: Bot API 5.1
type VideoChatEnded struct {
	Duration int64 `json:"duration"`
}

// VideoChatParticipantsInvited represents a service message about new members invited to a video chat.
// Since: Bot API 5.1
type VideoChatParticipantsInvited struct {
	Users []User `json:"users"`
}

// SentGuestMessage describes an inline message sent by a guest bot.
// Since: Bot API 10.0
type SentGuestMessage struct {
	InlineMessageID string `json:"inline_message_id"`
}
