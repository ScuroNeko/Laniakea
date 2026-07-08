package tgapi

// User represents a Telegram user or bot.
// Since: Bot API 1.0
// See https://core.telegram.org/bots/api#user
type User struct {
	ID        int64   `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name,omitempty"`
	Username  *string `json:"username,omitempty"`

	IsBot                     bool    `json:"is_bot"`                                  // Since: Bot API 3.3
	LanguageCode              *string `json:"language_code,omitempty"`                 // Since: Bot API 3.0
	IsPremium                 *bool   `json:"is_premium,omitempty"`                    // Since: Bot API 6.1
	AddedToAttachmentMenu     *bool   `json:"added_to_attachment_menu,omitempty"`      // Since: Bot API 6.1
	CanJoinGroups             *bool   `json:"can_join_groups,omitempty"`               // Since: Bot API 4.6
	CanReadAllGroupMessages   *bool   `json:"can_read_all_group_messages,omitempty"`   // Since: Bot API 4.6
	SupportsInlineQueries     *bool   `json:"supports_inline_queries,omitempty"`       // Since: Bot API 4.6
	CanConnectToBusiness      *bool   `json:"can_connect_to_business,omitempty"`       // Since: Bot API 7.2
	HasMainWebApp             *bool   `json:"has_main_web_app,omitempty"`              // Since: Bot API 7.8
	HasTopicsEnabled          *bool   `json:"has_topics_enabled,omitempty"`            // Since: Bot API 9.3
	AllowsUsersToCreateTopics *bool   `json:"allows_users_to_create_topics,omitempty"` // Since: Bot API 9.4
	CanManageBots             *bool   `json:"can_manage_bots,omitempty"`               // Since: Bot API 9.6
	SupportsGuestQueries      *bool   `json:"supports_guest_queries,omitempty"`        // Since: Bot API 10.0

	// SupportsJoinRequestQueries reports that the bot supports join request
	// queries and can be assigned to process them. Returned only in getMe.
	SupportsJoinRequestQueries *bool `json:"supports_join_request_queries,omitempty"` // Since: Bot API 10.1
}

// UserProfilePhotos represents a user's profile photos.
// Since: Bot API 1.4
// See https://core.telegram.org/bots/api#userprofilephotos
type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"`
	Photos     [][]PhotoSize `json:"photos"`
}

// UserProfileAudios represents a user's profile audios.
// Since: Bot API 9.3
// See https://core.telegram.org/bots/api#userprofileaudios
type UserProfileAudios struct {
	TotalCount int     `json:"total_count"`
	Audios     []Audio `json:"audios"`
}

// UserRating represents a user's rating with level progression.
// Since: Bot API 9.3
// See https://core.telegram.org/bots/api#userrating
type UserRating struct {
	Level              int `json:"level"`
	Rating             int `json:"rating"`
	CurrentLevelRating int `json:"current_level_rating"`
	NextLevelRating    int `json:"next_level_rating"`
}

// Birthdate represents a user's birthdate.
// Since: Bot API 7.2
// See https://core.telegram.org/bots/api#birthdate
type Birthdate struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}
