package tgapi

// User represents a Telegram user or bot.
// See https://core.telegram.org/bots/api#user
type User struct {
	ID                        int     `json:"id"`
	IsBot                     bool    `json:"is_bot"`
	FirstName                 string  `json:"first_name"`
	LastName                  *string `json:"last_name,omitempty"`
	Username                  *string `json:"username,omitempty"`
	LanguageCode              *string `json:"language_code,omitempty"`
	IsPremium                 *bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu     *bool   `json:"added_to_attachment_menu,omitempty"`
	CanJoinGroups             *bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages   *bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries     *bool   `json:"supports_inline_queries,omitempty"`
	CanConnectToBusiness      *bool   `json:"can_connect_to_business,omitempty"`
	HasMainWebApp             *bool   `json:"has_main_web_app,omitempty"`
	HasTopicsEnabled          *bool   `json:"has_topics_enabled,omitempty"`
	AllowsUsersToCreateTopics *bool   `json:"allows_users_to_create_topics,omitempty"`
}

// UserProfilePhotos represents a user's profile photos.
// See https://core.telegram.org/bots/api#userprofilephotos
type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"`
	Photos     [][]PhotoSize `json:"photos"`
}

// UserProfileAudios represents a user's profile audios.
// See https://core.telegram.org/bots/api#userprofileaudios
type UserProfileAudios struct {
	TotalCount int     `json:"total_count"`
	Audios     []Audio `json:"audios"`
}

// UserRating represents a user's rating with level progression.
// See https://core.telegram.org/bots/api#userrating
type UserRating struct {
	Level              int `json:"level"`
	Rating             int `json:"rating"`
	CurrentLevelRating int `json:"current_level_rating"`
	NextLevelRating    int `json:"next_level_rating"`
}

// Birthdate represents a user's birthdate.
// See https://core.telegram.org/bots/api#birthdate
type Birthdate struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}
