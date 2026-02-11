package tgapi

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

type OwnedGiftType string
type BaseOwnedGift struct {
	Type OwnedGiftType `json:"type"`
}
