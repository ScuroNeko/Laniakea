package tgapi

// BusinessIntro contains information about the business intro.
// See https://core.telegram.org/bots/api#businessintro
type BusinessIntro struct {
	Title   string   `json:"title,omitempty"`
	Message string   `json:"message,omitempty"`
	Sticker *Sticker `json:"sticker,omitempty"`
}

// BusinessLocation contains information about the business location.
// See https://core.telegram.org/bots/api#businesslocation
type BusinessLocation struct {
	Address  string    `json:"address"`
	Location *Location `json:"location,omitempty"`
}

// BusinessOpeningHoursInterval represents an interval of opening hours.
// See https://core.telegram.org/bots/api#businessopeninghoursinterval
type BusinessOpeningHoursInterval struct {
	OpeningMinute int `json:"opening_minute"`
	ClosingMinute int `json:"closing_minute"`
}

// BusinessOpeningHours represents the opening hours of a business.
// See https://core.telegram.org/bots/api#businessopeninghours
type BusinessOpeningHours struct {
	TimeZoneName string                         `json:"time_zone_name"`
	OpeningHours []BusinessOpeningHoursInterval `json:"opening_hours"`
}

// BusinessBotRights represents the rights of a business bot.
// All fields are optional booleans that, when present, are always true.
// See https://core.telegram.org/bots/api#businessbotrights
type BusinessBotRights struct {
	CanReply                   *bool `json:"can_reply,omitempty"`
	CanReadMessages            *bool `json:"can_read_messages,omitempty"`
	CanDeleteSentMessages      *bool `json:"can_delete_sent_messages,omitempty"`
	CanDeleteAllMessages       *bool `json:"can_delete_all_messages,omitempty"`
	CanEditName                *bool `json:"can_edit_name,omitempty"`
	CanEditBio                 *bool `json:"can_edit_bio,omitempty"`
	CanEditProfilePhoto        *bool `json:"can_edit_profile_photo,omitempty"`
	CanEditUsername            *bool `json:"can_edit_username,omitempty"`
	CanChangeGiftSettings      *bool `json:"can_change_gift_settings,omitempty"`
	CanViewGiftsAndStars       *bool `json:"can_view_gifts_and_stars,omitempty"`
	CanConvertGiftsToStars     *bool `json:"can_convert_gifts_to_stars,omitempty"`
	CanTransferAndUpgradeGifts *bool `json:"can_transfer_and_upgrade_gifts,omitempty"`
	CanTransferStars           *bool `json:"can_transfer_stars,omitempty"`
	CanManageStories           *bool `json:"can_manage_stories,omitempty"`
}

// BusinessConnection contains information about a business connection.
// See https://core.telegram.org/bots/api#businessconnection
type BusinessConnection struct {
	ID         string             `json:"id"`
	User       User               `json:"user"`
	UserChatID int                `json:"user_chat_id"`
	Date       int                `json:"date"`
	Rights     *BusinessBotRights `json:"rights,omitempty"`
	IsEnabled  bool               `json:"is_enabled"`
}

// InputStoryContentType indicates the type of input story content.
type InputStoryContentType string

const (
	InputStoryContentPhotoType InputStoryContentType = "photo"
	InputStoryContentVideoType InputStoryContentType = "video"
)

// InputStoryContent represents the content of a story to be posted.
// See https://core.telegram.org/bots/api#inputstorycontent
type InputStoryContent struct {
	Type InputStoryContentType `json:"type"`

	// Photo fields
	Photo *string `json:"photo,omitempty"`

	// Video fields
	Video               *string  `json:"video,omitempty"`
	Duration            *float64 `json:"duration,omitempty"`
	CoverFrameTimestamp *float64 `json:"cover_frame_timestamp,omitempty"`
	IsAnimation         *bool    `json:"is_animation,omitempty"`
}

// StoryAreaPosition describes the position of a clickable area on a story.
// See https://core.telegram.org/bots/api#storyareaposition
type StoryAreaPosition struct {
	XPercentage            float64 `json:"x_percentage"`
	YPercentage            float64 `json:"y_percentage"`
	WidthPercentage        float64 `json:"width_percentage"`
	HeightPercentage       float64 `json:"height_percentage"`
	RotationAngle          float64 `json:"rotation_angle"`
	CornerRadiusPercentage float64 `json:"corner_radius_percentage"`
}

// StoryAreaTypeType indicates the type of story area.
type StoryAreaTypeType string

const (
	StoryAreaTypeLocationType   StoryAreaTypeType = "location"
	StoryAreaTypeReactionType   StoryAreaTypeType = "suggested_reaction"
	StoryAreaTypeLinkType       StoryAreaTypeType = "link"
	StoryAreaTypeWeatherType    StoryAreaTypeType = "weather"
	StoryAreaTypeUniqueGiftType StoryAreaTypeType = "unique_gift"
)

// StoryAreaType describes the type of a clickable area on a story.
// Fields should be set according to the Type.
// See https://core.telegram.org/bots/api#storyareatype
type StoryAreaType struct {
	Type StoryAreaTypeType `json:"type"`

	// Location
	Latitude  *float64         `json:"latitude,omitempty"`
	Longitude *float64         `json:"longitude,omitempty"`
	Address   *LocationAddress `json:"address,omitempty"`

	// Suggested reaction
	ReactionType *ReactionType `json:"reaction_type,omitempty"`
	IsDark       *bool         `json:"is_dark,omitempty"`
	IsFlipped    *bool         `json:"is_flipped,omitempty"`

	// Link
	URL *string `json:"url,omitempty"`

	// Weather
	Temperature     *float64 `json:"temperature,omitempty"`
	Emoji           *string  `json:"emoji,omitempty"`
	BackgroundColor *int     `json:"background_color,omitempty"`

	// Unique gift
	Name *string `json:"name,omitempty"`
}

// StoryArea represents a clickable area on a story.
// See https://core.telegram.org/bots/api#storyarea
type StoryArea struct {
	Position StoryAreaPosition `json:"position"`
	Type     StoryAreaType     `json:"type"`
}
