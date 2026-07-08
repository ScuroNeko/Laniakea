package tgapi

import "encoding/json"

// UpdateType represents the type of incoming update.
type UpdateType string

const (
	// UpdateTypeUnknown marks an update whose payload does not match a known Telegram update kind.
	UpdateTypeUnknown UpdateType = "unknown"

	// UpdateTypeMessage is a regular message update.
	UpdateTypeMessage UpdateType = "message"
	// UpdateTypeEditedMessage is an edited message update.
	UpdateTypeEditedMessage UpdateType = "edited_message"
	// UpdateTypeChannelPost is a channel post update.
	UpdateTypeChannelPost UpdateType = "channel_post"
	// UpdateTypeEditedChannelPost is an edited channel post update.
	UpdateTypeEditedChannelPost UpdateType = "edited_channel_post"
	// UpdateTypeMessageReaction is a message reaction update.
	UpdateTypeMessageReaction UpdateType = "message_reaction"
	// UpdateTypeMessageReactionCount is a message reaction count update.
	UpdateTypeMessageReactionCount UpdateType = "message_reaction_count"

	// UpdateTypeBusinessConnection is a business connection update.
	UpdateTypeBusinessConnection UpdateType = "business_connection"
	// UpdateTypeBusinessMessage is a business message update.
	UpdateTypeBusinessMessage UpdateType = "business_message"
	// UpdateTypeEditedBusinessMessage is an edited business message update.
	UpdateTypeEditedBusinessMessage UpdateType = "edited_business_message"
	// UpdateTypeDeletedBusinessMessages is a deleted business messages update.
	UpdateTypeDeletedBusinessMessages UpdateType = "deleted_business_messages"

	// UpdateTypeInlineQuery is an inline query update.
	UpdateTypeInlineQuery UpdateType = "inline_query"
	// UpdateTypeChosenInlineResult is a chosen inline result update.
	UpdateTypeChosenInlineResult UpdateType = "chosen_inline_result"
	// UpdateTypeCallbackQuery is a callback query update.
	UpdateTypeCallbackQuery UpdateType = "callback_query"
	// UpdateTypeShippingQuery is a shipping query update.
	UpdateTypeShippingQuery UpdateType = "shipping_query"
	// UpdateTypePreCheckoutQuery is a pre-checkout query update.
	UpdateTypePreCheckoutQuery UpdateType = "pre_checkout_query"
	// UpdateTypePurchasedPaidMedia is a purchased paid media update.
	UpdateTypePurchasedPaidMedia UpdateType = "purchased_paid_media"
	// UpdateTypePoll is a poll update.
	UpdateTypePoll UpdateType = "poll"
	// UpdateTypePollAnswer is a poll answer update.
	UpdateTypePollAnswer UpdateType = "poll_answer"
	// UpdateTypeMyChatMember is a my chat member update.
	UpdateTypeMyChatMember UpdateType = "my_chat_member"
	// UpdateTypeChatMember is a chat member update.
	UpdateTypeChatMember UpdateType = "chat_member"
	// UpdateTypeChatJoinRequest is a chat join request update.
	UpdateTypeChatJoinRequest UpdateType = "chat_join_request"
	// UpdateTypeChatBoost is a chat boost update.
	UpdateTypeChatBoost UpdateType = "chat_boost"
	// UpdateTypeRemovedChatBoost is a removed chat boost update.
	UpdateTypeRemovedChatBoost UpdateType = "removed_chat_boost"

	// UpdateTypeManagedBot is a managed bot update.
	UpdateTypeManagedBot UpdateType = "managed_bot"

	// UpdateTypeGuestMessage is a guest message update.
	UpdateTypeGuestMessage UpdateType = "guest_message"
)

// Update represents an incoming update from Telegram.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#update
type Update struct {
	Type UpdateType `json:"-"`

	UpdateID          int      `json:"update_id"`
	Message           *Message `json:"message,omitempty"`
	EditedMessage     *Message `json:"edited_message,omitempty"`
	ChannelPost       *Message `json:"channel_post,omitempty"`        // Since: Bot API 2.3
	EditedChannelPost *Message `json:"edited_channel_post,omitempty"` // Since: Bot API 2.3

	BusinessConnection      *BusinessConnection          `json:"business_connection,omitempty"`       // Since: Bot API 7.2
	BusinessMessage         *Message                     `json:"business_message,omitempty"`          // Since: Bot API 7.2
	EditedBusinessMessage   *Message                     `json:"edited_business_message,omitempty"`   // Since: Bot API 7.2
	DeletedBusinessMessages *BusinessMessagesDeleted     `json:"deleted_business_messages,omitempty"` // Since: Bot API 7.2
	GuestMessage            *Message                     `json:"guest_message,omitempty"`             // Since: Bot API 10.0
	MessageReaction         *MessageReactionUpdated      `json:"message_reaction,omitempty"`          // Since: Bot API 7.0
	MessageReactionCount    *MessageReactionCountUpdated `json:"message_reaction_count,omitempty"`    // Since: Bot API 7.0

	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`         // Since: Bot API 1.7
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"` // Since: Bot API 1.8
	CallbackQuery      *CallbackQuery      `json:"callback_query,omitempty"`       // Since: Bot API 2.0
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`       // Since: Bot API 3.0
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`   // Since: Bot API 3.0
	PurchasedPaidMedia *PaidMediaPurchased `json:"purchased_paid_media,omitempty"` // Since: Bot API 7.10

	Poll             *Poll              `json:"poll,omitempty"`               // Since: Bot API 4.2
	PollAnswer       *PollAnswer        `json:"poll_answer,omitempty"`        // Since: Bot API 4.6
	MyChatMember     *ChatMemberUpdated `json:"my_chat_member,omitempty"`     // Since: Bot API 5.1
	ChatMember       *ChatMemberUpdated `json:"chat_member,omitempty"`        // Since: Bot API 5.1
	ChatJoinRequest  *ChatJoinRequest   `json:"chat_join_request,omitempty"`  // Since: Bot API 5.4
	ChatBoost        *ChatBoostUpdated  `json:"chat_boost,omitempty"`         // Since: Bot API 7.0
	RemovedChatBoost *ChatBoostRemoved  `json:"removed_chat_boost,omitempty"` // Since: Bot API 7.0

	ManagedBot *ManagedBotUpdated `json:"managed_bot,omitempty"` // Since: Bot API 9.6
}

// UnmarshalJSON decodes an update and derives its Type from the populated payload field.
func (u *Update) UnmarshalJSON(data []byte) error {
	type Alias Update

	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*u = Update(aux)

	switch {
	case u.Message != nil:
		u.Type = UpdateTypeMessage
	case u.EditedMessage != nil:
		u.Type = UpdateTypeEditedMessage
	case u.ChannelPost != nil:
		u.Type = UpdateTypeChannelPost
	case u.EditedChannelPost != nil:
		u.Type = UpdateTypeEditedChannelPost

	case u.BusinessConnection != nil:
		u.Type = UpdateTypeBusinessConnection
	case u.BusinessMessage != nil:
		u.Type = UpdateTypeBusinessMessage
	case u.EditedBusinessMessage != nil:
		u.Type = UpdateTypeEditedBusinessMessage
	case u.DeletedBusinessMessages != nil:
		u.Type = UpdateTypeDeletedBusinessMessages
	case u.GuestMessage != nil:
		u.Type = UpdateTypeGuestMessage
	case u.MessageReaction != nil:
		u.Type = UpdateTypeMessageReaction
	case u.MessageReactionCount != nil:
		u.Type = UpdateTypeMessageReactionCount

	case u.InlineQuery != nil:
		u.Type = UpdateTypeInlineQuery
	case u.ChosenInlineResult != nil:
		u.Type = UpdateTypeChosenInlineResult
	case u.CallbackQuery != nil:
		u.Type = UpdateTypeCallbackQuery
	case u.ShippingQuery != nil:
		u.Type = UpdateTypeShippingQuery
	case u.PreCheckoutQuery != nil:
		u.Type = UpdateTypePreCheckoutQuery
	case u.PurchasedPaidMedia != nil:
		u.Type = UpdateTypePurchasedPaidMedia

	case u.Poll != nil:
		u.Type = UpdateTypePoll
	case u.PollAnswer != nil:
		u.Type = UpdateTypePollAnswer
	case u.MyChatMember != nil:
		u.Type = UpdateTypeMyChatMember
	case u.ChatMember != nil:
		u.Type = UpdateTypeChatMember
	case u.ChatJoinRequest != nil:
		u.Type = UpdateTypeChatJoinRequest
	case u.ChatBoost != nil:
		u.Type = UpdateTypeChatBoost
	case u.RemovedChatBoost != nil:
		u.Type = UpdateTypeRemovedChatBoost
	case u.ManagedBot != nil:
		u.Type = UpdateTypeManagedBot
	default:
		u.Type = UpdateTypeUnknown
	}

	return nil
}

// WebhookInfo describes the current webhook status.
// Since: Bot API 2.2
// See https://core.telegram.org/bots/api#webhookinfo
type WebhookInfo struct {
	URL                          string   `json:"url"`
	HasCustomCertificate         bool     `json:"has_custom_certificate"`
	PendingUpdateCount           int      `json:"pending_update_count"`
	IPAddress                    string   `json:"ip_address,omitempty"`
	LastErrorDate                int      `json:"last_error_date,omitempty"`
	LastErrorMessage             string   `json:"last_error_message,omitempty"`
	LastSynchronizationErrorDate int      `json:"last_synchronization_error_date,omitempty"`
	MaxConnections               int      `json:"max_connections,omitempty"`
	AllowedUpdates               []string `json:"allowed_updates,omitempty"`
}

// ProximityAlertTriggered represents the content of a service message sent when a user triggers a proximity alert.
// Since: Bot API 5.0
type ProximityAlertTriggered struct {
	Traveler User `json:"traveler"`
	Watcher  User `json:"watcher"`
	Distance int  `json:"distance"`
}

// InlineQuery represents an incoming inline query.
// Since: Bot API 1.7
// See https://core.telegram.org/bots/api#inlinequery
type InlineQuery struct {
	ID       string    `json:"id"`
	From     User      `json:"from"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
	ChatType *ChatType `json:"chat_type,omitempty"`
	Location *Location `json:"location,omitempty"`
}

// ChosenInlineResult represents a result of an inline query that was chosen by the user.
// Since: Bot API 1.8
// See https://core.telegram.org/bots/api#choseninlineresult
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`
	From            User      `json:"from"`
	Location        *Location `json:"location,omitempty"`
	InlineMessageID string    `json:"inline_message_id"`
	Query           string    `json:"query"`
}

// File represents a file ready to be downloaded.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#file
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

// ChatMemberUpdated represents changes in the status of a chat member.
// Since: Bot API 5.1
// See https://core.telegram.org/bots/api#chatmemberupdated
type ChatMemberUpdated struct {
	Chat                    Chat            `json:"chat"`
	From                    User            `json:"from"`
	Date                    int64           `json:"date"`
	OldChatMember           ChatMember      `json:"old_chat_member"`
	NewChatMember           ChatMember      `json:"new_chat_member"`
	InviteLink              *ChatInviteLink `json:"invite_link,omitempty"`
	ViaJoinRequest          *bool           `json:"via_join_request,omitempty"`
	ViaChatFolderInviteLink *bool           `json:"via_chat_folder_invite_link,omitempty"`
}

// ChatJoinRequest represents a join request sent to a chat.
// Since: Bot API 5.4
// See https://core.telegram.org/bots/api#chatjoinrequest
type ChatJoinRequest struct {
	Chat       Chat            `json:"chat"`
	From       User            `json:"from"`
	UserChatID int64           `json:"user_chat_id"`
	Date       int64           `json:"date"`
	Bio        *string         `json:"bio,omitempty"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`

	// QueryID identifies the join request query; present only for bots
	// assigned to process join requests. When set, the bot must call
	// SendChatJoinRequestWebApp or AnswerChatJoinRequestQuery within 10 seconds.
	QueryID *string `json:"query_id,omitempty"` // Since: Bot API 10.1
}

// Location represents a point on the map.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#location
type Location struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy"`
	LivePeriod           int     `json:"live_period"`
	Heading              int     `json:"heading"`
	ProximityAlertRadius int     `json:"proximity_alert_radius"`
}

// LocationAddress represents a human-readable address of a location.
// Since: Bot API 8.0
type LocationAddress struct {
	CountryCode string  `json:"country_code"`
	State       *string `json:"state,omitempty"`
	City        *string `json:"city,omitempty"`
	Street      *string `json:"street,omitempty"`
}

// Venue represents a venue.
// Since: Bot API 2.0
// See https://core.telegram.org/bots/api#venue
type Venue struct {
	Location        Location `json:"location"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	FoursquareID    string   `json:"foursquare_id,omitempty"`
	FoursquareType  string   `json:"foursquare_type,omitempty"`
	GooglePlaceID   string   `json:"google_place_id,omitempty"`
	GooglePlaceType string   `json:"google_place_type,omitempty"`
}

// WebAppInfo contains information about a Web App.
// Since: Bot API 6.0
// See https://core.telegram.org/bots/api#webappinfo
type WebAppInfo struct {
	URL string `json:"url"`
}

// WebAppData represents data sent from a Web App to the bot.
// Since: Bot API 6.0
type WebAppData struct {
	Data       string `json:"data"`
	ButtonText string `json:"button_text"`
}

// StarAmount represents an amount of Telegram Stars.
// Since: Bot API 7.5
type StarAmount struct {
	Amount         int `json:"amount"`
	NanostarAmount int `json:"nanostar_amount"`
}

// AcceptedGiftTypes represents the types of gifts accepted by a user or chat.
// Since: Bot API 9.0
type AcceptedGiftTypes struct {
	UnlimitedGifts      bool `json:"unlimited_gifts"`
	LimitedGifts        bool `json:"limited_gifts"`
	UniqueGifts         bool `json:"unique_gifts"`
	PremiumSubscription bool `json:"premium_subscription"`
	GiftsFromChannels   bool `json:"gifts_from_channels"`
}

// GiftBackground represents the background of a gift.
// Since: Bot API 9.0
type GiftBackground struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	TextColor   int `json:"text_color"`
}

// Gift represents a gift that can be sent.
// Since: Bot API 9.0
type Gift struct {
	ID                     string          `json:"id"`
	Sticker                Sticker         `json:"sticker"`
	StarCount              int             `json:"star_count"`
	UpdateStarCount        *int            `json:"update_star_count,omitempty"`
	IsPremium              *bool           `json:"is_premium,omitempty"`
	HasColors              *bool           `json:"has_colors,omitempty"`
	TotalCount             *int            `json:"total_count,omitempty"`
	RemainingCount         *int            `json:"remaining_count,omitempty"`
	PersonalTotalCount     *int            `json:"personal_total_count,omitempty"`
	PersonalRemainingCount *int            `json:"personal_remaining_count,omitempty"`
	Background             *GiftBackground `json:"background,omitempty"`
	UniqueGiftVariantColor *int            `json:"unique_gift_variant_color,omitempty"`
	PublisherChat          *Chat           `json:"publisher_chat,omitempty"`
}

// Gifts represents a list of gifts.
// Since: Bot API 9.0
type Gifts struct {
	Gifts []Gift `json:"gifts"`
}

// UniqueGiftModel describes the model component of a unique gift.
// Since: Bot API 9.0
type UniqueGiftModel struct {
	Name           string  `json:"name"`
	Sticker        Sticker `json:"sticker"`
	RarityPerMille int     `json:"rarity_per_mille"`
	Rarity         string  `json:"rarity,omitempty"`
}

// UniqueGiftSymbol describes the symbol component of a unique gift.
// Since: Bot API 9.0
type UniqueGiftSymbol struct {
	Name           string  `json:"name"`
	Sticker        Sticker `json:"sticker"`
	RarityPerMille int     `json:"rarity_per_mille"`
}

// UniqueGiftBackdropColors describes the colors of a unique gift backdrop.
// Since: Bot API 9.0
type UniqueGiftBackdropColors struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	SymbolColor int `json:"symbol_color"`
	TextColor   int `json:"text_color"`
}

// UniqueGiftBackdrop describes the backdrop of a unique gift.
// Since: Bot API 9.0
type UniqueGiftBackdrop struct {
	Name           string                   `json:"name"`
	Colors         UniqueGiftBackdropColors `json:"colors"`
	RarityPerMille int                      `json:"rarity_per_mille"`
}

// UniqueGiftColors represents color information for a unique gift.
// Since: Bot API 9.3
type UniqueGiftColors struct {
	ModelCustomEmojiID    string `json:"model_custom_emoji_id"`
	SymbolCustomEmojiID   string `json:"symbol_custom_emoji_id"`
	LightThemeMainColor   int    `json:"light_theme_main_color"`
	LightThemeOtherColors []int  `json:"light_theme_other_colors"`
	DarkThemeMainColor    int    `json:"dark_theme_main_color"`
	DarkThemeOtherColors  []int  `json:"dark_theme_other_colors"`
}

// UniqueGift represents a unique gift.
// Since: Bot API 9.0
type UniqueGift struct {
	GiftID   string             `json:"gift_id"`
	BaseName string             `json:"base_name"`
	Name     string             `json:"name"`
	Number   int                `json:"number"`
	Model    UniqueGiftModel    `json:"model"`
	Symbol   UniqueGiftSymbol   `json:"symbol"`
	Backdrop UniqueGiftBackdrop `json:"backdrop"`

	IsPremium        bool              `json:"is_premium,omitempty"`
	IsBurned         bool              `json:"is_burned,omitempty"`
	IsFromBlockchain bool              `json:"is_from_blockchain,omitempty"`
	Colors           *UniqueGiftColors `json:"colors,omitempty"`
	PublisherChat    *Chat             `json:"publisher_chat,omitempty"`
}

// GiftInfo contains information about a received gift.
// Since: Bot API 9.0
type GiftInfo struct {
	Gift Gift `json:"gift"`

	OwnedGiftID             string          `json:"owned_gift_id,omitempty"`
	ConvertStarCount        int             `json:"convert_star_count,omitempty"`
	PrepaidUpgradeStarCount int             `json:"prepaid_upgrade_star_count,omitempty"`
	IsUpgradeSeparate       bool            `json:"is_upgrade_separate,omitempty"`
	CanBeUpgraded           bool            `json:"can_be_upgraded,omitempty"`
	Text                    string          `json:"text,omitempty"`
	Entities                []MessageEntity `json:"entities,omitempty"`
	IsPrivate               bool            `json:"is_private,omitempty"`
	UniqueGiftNumber        int             `json:"unique_gift_number,omitempty"`
}

// UniqueGiftInfo contains information about a received unique gift.
// Since: Bot API 9.0
type UniqueGiftInfo struct {
	Gift               UniqueGift `json:"gift"`
	Origin             string     `json:"origin"`
	LastResaleCurrency string     `json:"last_resale_currency,omitempty"`
	LastResaleAmount   int        `json:"last_resale_amount,omitempty"`
	OwnedGiftID        string     `json:"owned_gift_id,omitempty"`
	TransferStarCount  int        `json:"transfer_star_count,omitempty"`
	NextTransferDate   int        `json:"next_transfer_date,omitempty"`
}

// OwnedGiftType represents the type of an owned gift.
// Since: Bot API 9.0
type OwnedGiftType string

const (
	// OwnedGiftRegularType identifies a regular owned gift.
	OwnedGiftRegularType OwnedGiftType = "regular"
	// OwnedGiftUniqueType identifies a unique owned gift.
	OwnedGiftUniqueType OwnedGiftType = "unique"
)

// OwnedGift represents a gift owned by a user or chat.
// Since: Bot API 9.0
type OwnedGift struct {
	Type        OwnedGiftType `json:"type"`
	OwnedGiftID string        `json:"ownen_gift_id,omitempty"`
	SendDate    int           `json:"send_date,omitempty"`
	IsSaved     bool          `json:"is_saved,omitempty"`

	// Fields specific to "regular" type
	Gift                    Gift            `json:"gift"`
	SenderUser              *User           `json:"sender_user,omitempty"`
	Text                    string          `json:"text,omitempty"`
	Entities                []MessageEntity `json:"entities,omitempty"`
	IsPrivate               bool            `json:"is_private,omitempty"`
	CanBeUpgraded           bool            `json:"can_be_upgraded,omitempty"`
	WasRefunded             bool            `json:"was_refunded,omitempty"`
	ConvertStarCount        int             `json:"convert_star_count,omitempty"`
	PrepaidUpgradeStarCount int             `json:"prepaid_upgrade_star_count,omitempty"`
	IsUpgradeSeparate       bool            `json:"is_upgrade_separate,omitempty"`
	UniqueGiftNumber        int             `json:"unique_gift_number,omitempty"`

	// Fields specific to "unique" type
	CanBeTransferred  bool `json:"can_be_transferred,omitempty"`
	TransferStarCount int  `json:"transfer_star_count,omitempty"`
	NextTransferDate  int  `json:"next_transfer_date,omitempty"`
}

// OwnedGifts represents a list of owned gifts with pagination.
// Since: Bot API 9.0
type OwnedGifts struct {
	TotalCount int         `json:"total_count"`
	Gifts      []OwnedGift `json:"gifts"`
	NextOffset string      `json:"next_offset"`
}

// GiveawayCreated represents a service message about a giveaway being created.
// Since: Bot API 7.0
type GiveawayCreated struct {
	PrizeStarCount int `json:"prize_star_count,omitempty"`
}

// Giveaway represents a message about a scheduled giveaway.
// Since: Bot API 7.0
type Giveaway struct {
	Chats                []Chat `json:"chats"`
	WinnersSelectionDate int    `json:"winners_selection_date"`
	WinnerCount          int    `json:"winner_count"`

	OnlyNewMembers                bool     `json:"only_new_members,omitempty"`
	HasPublicWinners              bool     `json:"has_public_winners,omitempty"`
	PrizeDescription              string   `json:"prize_description,omitempty"`
	CountryCodes                  []string `json:"country_codes,omitempty"`
	PrizeStarCount                int      `json:"prize_star_count,omitempty"`
	PremiumSubscriptionMonthCount int      `json:"premium_subscription_month_count,omitempty"`
}

// GiveawayWinners represents a message about the completion of a giveaway with public winners.
// Since: Bot API 7.0
type GiveawayWinners struct {
	Chat                 Chat   `json:"chat"`
	GiveawayMessageID    int    `json:"giveaway_message_id"`
	WinnersSelectionDate int    `json:"winners_selection_date"`
	WinnerCount          int    `json:"winner_count"`
	Winners              []User `json:"winners"`

	AdditionalChatCount           int    `json:"additional_chat_count,omitempty"`
	PrizeStarCount                int    `json:"prize_star_count,omitempty"`
	PremiumSubscriptionMonthCount int    `json:"premium_subscription_month_count,omitempty"`
	UnclaimedPrizeCount           int    `json:"unclaimed_prize_count,omitempty"`
	OnlyNewMembers                bool   `json:"only_new_members,omitempty"`
	WasRefunded                   bool   `json:"was_refunded,omitempty"`
	PrizeDescription              string `json:"prize_description,omitempty"`
}

// GiveawayCompleted represents a service message about the completion of a giveaway without public winners.
// Since: Bot API 7.0
type GiveawayCompleted struct {
	WinnerCount         int      `json:"winner_count"`
	UnclaimedPrizeCount int      `json:"unclaimed_prize_count,omitempty"`
	GiveawayMessage     *Message `json:"giveaway_message,omitempty"`
	IsStarGiveaway      bool     `json:"is_star_giveaway,omitempty"`
}

// WriteAccessAllowed represents a service message about a user allowing a bot to write messages.
// Since: Bot API 6.4
type WriteAccessAllowed struct {
	FromRequest        bool   `json:"from_request,omitempty"`
	WebAppName         string `json:"web_app_name,omitempty"`
	FromAttachmentMenu bool   `json:"from_attachment_menu,omitempty"`
}

// BackgroundFillType represents the type of a background fill.
// Since: Bot API 7.5
type BackgroundFillType string

const (
	BackgroundFillSolidType            BackgroundFillType = "solid"
	BackgroundFillGradientType         BackgroundFillType = "gradient"
	BackgroundFillFreeformGradientType BackgroundFillType = "freeform_gradient"
)

// BackgroundFill describes the way a background is filled.
// Since: Bot API 7.5
type BackgroundFill struct {
	Type BackgroundFillType `json:"type"`

	Color int `json:"color,omitempty"`

	TopColor      int `json:"top_color,omitempty"`
	BottomColor   int `json:"bottom_color,omitempty"`
	RotationAngle int `json:"rotation_angle,omitempty"`

	Colors []int `json:"colors,omitempty"`
}

// BackgroundTypeType represents the type of a chat background.
// Since: Bot API 7.5
type BackgroundTypeType string

const (
	BackgroundTypeFillType      BackgroundTypeType = "fill"
	BackgroundTypeWallpaperType BackgroundTypeType = "wallpaper"
	BackgroundTypePatternType   BackgroundTypeType = "pattern"
	BackgroundTypeChatThemeType BackgroundTypeType = "chat_theme"
)

// BackgroundType describes the type of a background.
// Since: Bot API 7.5
type BackgroundType struct {
	Type BackgroundTypeType `json:"type"`

	Fill             *BackgroundFill `json:"fill,omitempty"`
	DarkThemeDimming int             `json:"dark_theme_dimming,omitempty"`

	Document  *Document `json:"document,omitempty"`
	IsBlurred bool      `json:"is_blurred,omitempty"`
	IsMoving  bool      `json:"is_moving,omitempty"`

	Intensity  int  `json:"intensity,omitempty"`
	IsInverted bool `json:"is_inverted,omitempty"`

	ThemeName string `json:"theme_name,omitempty"`
}
