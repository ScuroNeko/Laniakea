package tgapi

type Update struct {
	UpdateID          int      `json:"update_id"`
	Message           *Message `json:"message,omitempty"`
	EditedMessage     *Message `json:"edited_message,omitempty"`
	ChannelPost       *Message `json:"channel_post,omitempty"`
	EditedChannelPost *Message `json:"edited_channel_post,omitempty"`

	BusinessConnection     *BusinessConnection          `json:"business_connection,omitempty"`
	BusinessMessage        *Message                     `json:"business_message,omitempty"`
	EditedBusinessMessage  *Message                     `json:"edited_business_message,omitempty"`
	DeletedBusinessMessage *Message                     `json:"deleted_business_messages,omitempty"`
	MessageReaction        *MessageReactionUpdated      `json:"message_reaction,omitempty"`
	MessageReactionCount   *MessageReactionCountUpdated `json:"message_reaction_count,omitempty"`

	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	CallbackQuery      *CallbackQuery      `json:"callback_query,omitempty"`
	ShippingQuery      ShippingQuery       `json:"shipping_query,omitempty"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`
	PurchasedPaidMedia *PaidMediaPurchased `json:"purchased_paid_media,omitempty"`

	Poll             any `json:"poll,omitempty"`
	PollAnswer       any `json:"poll_answer,omitempty"`
	MyChatMember     any `json:"my_chat_member,omitempty"`
	ChatMember       any `json:"chat_member,omitempty"`
	ChatJoinRequest  any `json:"chat_join_request,omitempty"`
	ChatBoost        any `json:"chat_boost,omitempty"`
	RemovedChatBoost any `json:"removed_chat_boost,omitempty"`
}

type InlineQuery struct {
	ID       string    `json:"id"`
	From     User      `json:"from"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
	ChatType *ChatType `json:"chat_type,omitempty"`
	Location *Location `json:"location,omitempty"`
}
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`
	From            User      `json:"from"`
	Location        *Location `json:"location,omitempty"`
	InlineMessageID string    `json:"inline_message_id"`
	Query           string    `json:"query"`
}

type ShippingQuery struct {
	ID              string          `json:"id"`
	From            User            `json:"from"`
	InvoicePayload  string          `json:"invoice_payload"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

type OrderInfo struct {
	Name            string          `json:"name"`
	PhoneNumber     string          `json:"phone_number"`
	Email           string          `json:"email"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             User       `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

type PaidMediaPurchased struct {
	From             User   `json:"from"`
	PaidMediaPayload string `json:"paid_media_payload"`
}

type File struct {
	FileId       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

type Audio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`

	Performer string     `json:"performer,omitempty"`
	Title     string     `json:"title,omitempty"`
	FileName  string     `json:"file_name,omitempty"`
	MimeType  string     `json:"mime_type,omitempty"`
	FileSize  int        `json:"file_size,omitempty"`
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
}

type Location struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy"`
	LivePeriod           int     `json:"live_period"`
	Heading              int     `json:"heading"`
	ProximityAlertRadius int     `json:"proximity_alert_radius"`
}
type LocationAddress struct {
	CountryCode string  `json:"country_code"`
	State       *string `json:"state,omitempty"`
	City        *string `json:"city,omitempty"`
	Street      *string `json:"street,omitempty"`
}
type Venue struct {
	Location        Location `json:"location"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	FoursquareID    string   `json:"foursquare_id,omitempty"`
	FoursquareType  string   `json:"foursquare_type,omitempty"`
	GooglePlaceID   string   `json:"google_place_id,omitempty"`
	GooglePlaceType string   `json:"google_place_type,omitempty"`
}

type WebAppInfo struct {
	URL string `json:"url"`
}

type StarAmount struct {
	Amount         int `json:"amount"`
	NanostarAmount int `json:"nanostar_amount"`
}

type Story struct {
	Chat Chat `json:"chat"`
	ID   int  `json:"id"`
}

// Gifts

type AcceptedGiftTypes struct {
	UnlimitedGifts      bool `json:"unlimited_gifts"`
	LimitedGifts        bool `json:"limited_gifts"`
	UniqueGifts         bool `json:"unique_gifts"`
	PremiumSubscription bool `json:"premium_subscription"`
	GiftsFromChannels   bool `json:"gifts_from_channels"`
}

type UniqueGiftColors struct {
	ModelCustomEmojiID    string `json:"model_custom_emoji_id"`
	SymbolCustomEmojiID   string `json:"symbol_custom_emoji_id"`
	LightThemeMainColor   int    `json:"light_theme_main_color"`
	LightThemeOtherColors []int  `json:"light_theme_other_colors"`
	DarkThemeMainColor    int    `json:"dark_theme_main_color"`
	DarkThemeOtherColors  []int  `json:"dark_theme_other_colors"`
}

type GiftBackground struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	TextColor   int `json:"text_color"`
}
type Gift struct {
	ID                     string         `json:"id"`
	Sticker                Sticker        `json:"sticker"`
	StarCount              int            `json:"star_count"`
	UpdateStarCount        *int           `json:"update_star_count,omitempty"`
	IsPremium              *bool          `json:"is_premium,omitempty"`
	HasColors              *bool          `json:"has_colors,omitempty"`
	TotalCount             *int           `json:"total_count,omitempty"`
	RemainingCount         *int           `json:"remaining_count,omitempty"`
	PersonalTotalCount     *int           `json:"personal_total_count,omitempty"`
	PersonalRemainingCount *int           `json:"personal_remaining_count,omitempty"`
	Background             GiftBackground `json:"background,omitempty"`
	UniqueGiftVariantColor *int           `json:"unique_gift_variant_color,omitempty"`
	PublisherChat          *Chat          `json:"publisher_chat,omitempty"`
}
type Gifts struct {
	Gifts []Gift `json:"gifts"`
}

const (
	OwnedGiftRegularType OwnedGiftType = "regular"
	OwnedGiftUniqueType  OwnedGiftType = "unique"
)

type OwnedGiftType string
type BaseOwnedGift struct {
	Type        OwnedGiftType `json:"type"`
	OwnerGiftID *string       `json:"owner_gift_id,omitempty"`
	SendDate    *int          `json:"send_date,omitempty"`
	IsSaved     *bool         `json:"is_saved,omitempty"`
}
type OwnedGiftRegular struct {
	BaseOwnedGift
	Gift                    Gift            `json:"gift"`
	SenderUser              User            `json:"sender_user,omitempty"`
	Text                    string          `json:"text,omitempty"`
	Entities                []MessageEntity `json:"entities,omitempty"`
	IsPrivate               *bool           `json:"is_private,omitempty"`
	CanBeUpgraded           *bool           `json:"can_be_upgraded,omitempty"`
	WasRefunded             *bool           `json:"was_refunded,omitempty"`
	ConvertStarCount        *int            `json:"convert_star_count,omitempty"`
	PrepaidUpgradeStarCount *int            `json:"prepaid_upgrade_star_count,omitempty"`
	IsUpgradeSeparate       *bool           `json:"is_upgrade_separate,omitempty"`
	UniqueGiftNumber        *int            `json:"unique_gift_number,omitempty"`
}
type OwnedGiftUnique struct {
	BaseOwnedGift
	CanBeTransferred  *bool `json:"can_be_transferred,omitempty"`
	TransferStarCount *int  `json:"transfer_star_count,omitempty"`
	NextTransferDate  *int  `json:"next_transfer_date,omitempty"`
}
type OwnedGifts struct {
	TotalCount int `json:"total_count"`
	Gifts      []struct {
		OwnedGiftRegular
		OwnedGiftUnique
	} `json:"gifts"`
	NextOffset string `json:"next_offset"`
}
