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

	UpdateTypeManagedBot UpdateType = "managed_bot"
)

// Update represents an incoming update from Telegram.
// See https://core.telegram.org/bots/api#update
type Update struct {
	Type UpdateType `json:"-"`

	UpdateID          int      `json:"update_id"`
	Message           *Message `json:"message,omitempty"`
	EditedMessage     *Message `json:"edited_message,omitempty"`
	ChannelPost       *Message `json:"channel_post,omitempty"`
	EditedChannelPost *Message `json:"edited_channel_post,omitempty"`

	BusinessConnection      *BusinessConnection          `json:"business_connection,omitempty"`
	BusinessMessage         *Message                     `json:"business_message,omitempty"`
	EditedBusinessMessage   *Message                     `json:"edited_business_message,omitempty"`
	DeletedBusinessMessages *BusinessMessagesDeleted     `json:"deleted_business_messages,omitempty"`
	MessageReaction         *MessageReactionUpdated      `json:"message_reaction,omitempty"`
	MessageReactionCount    *MessageReactionCountUpdated `json:"message_reaction_count,omitempty"`

	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	CallbackQuery      *CallbackQuery      `json:"callback_query,omitempty"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`
	PurchasedPaidMedia *PaidMediaPurchased `json:"purchased_paid_media,omitempty"`

	Poll             *Poll              `json:"poll,omitempty"`
	PollAnswer       *PollAnswer        `json:"poll_answer,omitempty"`
	MyChatMember     *ChatMemberUpdated `json:"my_chat_member,omitempty"`
	ChatMember       *ChatMemberUpdated `json:"chat_member,omitempty"`
	ChatJoinRequest  *ChatJoinRequest   `json:"chat_join_request,omitempty"`
	ChatBoost        *ChatBoostUpdated  `json:"chat_boost,omitempty"`
	RemovedChatBoost *ChatBoostRemoved  `json:"removed_chat_boost,omitempty"`

	ManagedBot *ManagedBotUpdated `json:"managed_bot,omitempty"`
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

type ProximityAlertTriggered struct {
	Traveler User `json:"traveler"`
	Watcher  User `json:"watcher"`
	Distance int  `json:"distance"`
}

// InlineQuery represents an incoming inline query.
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
// See https://core.telegram.org/bots/api#choseninlineresult
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`
	From            User      `json:"from"`
	Location        *Location `json:"location,omitempty"`
	InlineMessageID string    `json:"inline_message_id"`
	Query           string    `json:"query"`
}

// File represents a file ready to be downloaded.
// See https://core.telegram.org/bots/api#file
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

// ChatMemberUpdated represents changes in the status of a chat member.
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
// See https://core.telegram.org/bots/api#chatjoinrequest
type ChatJoinRequest struct {
	Chat       Chat            `json:"chat"`
	From       User            `json:"from"`
	UserChatID int64           `json:"user_chat_id"`
	Date       int64           `json:"date"`
	Bio        *string         `json:"bio,omitempty"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

// Location represents a point on the map.
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
type LocationAddress struct {
	CountryCode string  `json:"country_code"`
	State       *string `json:"state,omitempty"`
	City        *string `json:"city,omitempty"`
	Street      *string `json:"street,omitempty"`
}

// Venue represents a venue.
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
// See https://core.telegram.org/bots/api#webappinfo
type WebAppInfo struct {
	URL string `json:"url"`
}

type WebAppData struct {
	Data       string `json:"data"`
	ButtonText string `json:"button_text"`
}

// StarAmount represents an amount of Telegram Stars.
type StarAmount struct {
	Amount         int `json:"amount"`
	NanostarAmount int `json:"nanostar_amount"`
}

// AcceptedGiftTypes represents the types of gifts accepted by a user or chat.
type AcceptedGiftTypes struct {
	UnlimitedGifts      bool `json:"unlimited_gifts"`
	LimitedGifts        bool `json:"limited_gifts"`
	UniqueGifts         bool `json:"unique_gifts"`
	PremiumSubscription bool `json:"premium_subscription"`
	GiftsFromChannels   bool `json:"gifts_from_channels"`
}

// GiftBackground represents the background of a gift.
type GiftBackground struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	TextColor   int `json:"text_color"`
}

// Gift represents a gift that can be sent.
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
type Gifts struct {
	Gifts []Gift `json:"gifts"`
}

type UniqueGiftModel struct {
	Name           string  `json:"name"`
	Sticker        Sticker `json:"sticker"`
	RarityPerMille int     `json:"rarity_per_mille"`
	Rarity         string  `json:"rarity,omitempty"`
}
type UniqueGiftSymbol struct {
	Name           string  `json:"name"`
	Sticker        Sticker `json:"sticker"`
	RarityPerMille int     `json:"rarity_per_mille"`
}
type UniqueGiftBackdropColors struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	SymbolColor int `json:"symbol_color"`
	TextColor   int `json:"text_color"`
}
type UniqueGiftBackdrop struct {
	Name           string                   `json:"name"`
	Colors         UniqueGiftBackdropColors `json:"colors"`
	RarityPerMille int                      `json:"rarity_per_mille"`
}

// UniqueGiftColors represents color information for a unique gift.
type UniqueGiftColors struct {
	ModelCustomEmojiID    string `json:"model_custom_emoji_id"`
	SymbolCustomEmojiID   string `json:"symbol_custom_emoji_id"`
	LightThemeMainColor   int    `json:"light_theme_main_color"`
	LightThemeOtherColors []int  `json:"light_theme_other_colors"`
	DarkThemeMainColor    int    `json:"dark_theme_main_color"`
	DarkThemeOtherColors  []int  `json:"dark_theme_other_colors"`
}

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
type OwnedGiftType string

const (
	// OwnedGiftRegularType identifies a regular owned gift.
	OwnedGiftRegularType OwnedGiftType = "regular"
	// OwnedGiftUniqueType identifies a unique owned gift.
	OwnedGiftUniqueType OwnedGiftType = "unique"
)

// OwnedGift represents a gift owned by a user or chat.
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
type OwnedGifts struct {
	TotalCount int         `json:"total_count"`
	Gifts      []OwnedGift `json:"gifts"`
	NextOffset string      `json:"next_offset"`
}

type GiveawayCreated struct {
	PrizeStarCount int `json:"prize_star_count,omitempty"`
}

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

type GiveawayCompleted struct {
	WinnerCount         int      `json:"winner_count"`
	UnclaimedPrizeCount int      `json:"unclaimed_prize_count,omitempty"`
	GiveawayMessage     *Message `json:"giveaway_message,omitempty"`
	IsStarGiveaway      bool     `json:"is_star_giveaway,omitempty"`
}

type WriteAccessAllowed struct {
	FromRequest        bool   `json:"from_request,omitempty"`
	WebAppName         string `json:"web_app_name,omitempty"`
	FromAttachmentMenu bool   `json:"from_attachment_menu,omitempty"`
}

type BackgroundFillType string

const (
	BackgroundFillSolidType            BackgroundFillType = "solid"
	BackgroundFillGradientType         BackgroundFillType = "gradient"
	BackgroundFillFreeformGradientType BackgroundFillType = "freeform_gradient"
)

type BackgroundFill struct {
	Type BackgroundFillType `json:"type"`

	Color int `json:"color,omitempty"`

	TopColor      int `json:"top_color,omitempty"`
	BottomColor   int `json:"bottom_color,omitempty"`
	RotationAngle int `json:"rotation_angle,omitempty"`

	Colors []int `json:"colors,omitempty"`
}

type BackgroundTypeType string

const (
	BackgroundTypeFillType      BackgroundTypeType = "fill"
	BackgroundTypeWallpaperType BackgroundTypeType = "wallpaper"
	BackgroundTypePatternType   BackgroundTypeType = "pattern"
	BackgroundTypeChatThemeType BackgroundTypeType = "chat_theme"
)

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
