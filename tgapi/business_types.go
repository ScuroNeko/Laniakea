package tgapi

type BusinessIntro struct {
	Title   string   `json:"title,omitempty"`
	Message string   `json:"message,omitempty"`
	Sticker *Sticker `json:"sticker,omitempty"`
}
type BusinessLocation struct {
	Address  string    `json:"address"`
	Location *Location `json:"location,omitempty"`
}
type BusinessOpeningHoursInterval struct {
	OpeningMinute int `json:"opening_minute"`
	ClosingMinute int `json:"closing_minute"`
}
type BusinessOpeningHours struct {
	TimeZoneName string      `json:"time_zone_name"`
	OpeningHours []Birthdate `json:"opening_hours"`
}

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
type BusinessConnection struct {
	ID         string             `json:"id"`
	User       User               `json:"user"`
	UserChatID int                `json:"user_chat_id"`
	Date       int                `json:"date"`
	Rights     *BusinessBotRights `json:"rights,omitempty"`
	IsEnabled  bool               `json:"id_enabled"`
}

const (
	InputStoryContentPhotoType InputStoryContentType = "photo"
	InputStoryContentVideoType InputStoryContentType = "video"
)

type InputStoryContentType string
type InputStoryContent struct {
	Type InputStoryContentType `json:"type"`
	// Photo
	Photo *string `json:"photo,omitempty"`

	// Video
	Video               *string  `json:"video,omitempty"`
	Duration            *float64 `json:"duration,omitempty"`
	CoverFrameTimestamp *float64 `json:"cover_frame_timestamp,omitempty"`
	IsAnimation         *bool    `json:"is_animation,omitempty"`
}

type StoryAreaPosition struct {
	XPercentage            float64 `json:"x_percentage"`
	YPercentage            float64 `json:"y_percentage"`
	WidthPercentage        float64 `json:"width_percentage"`
	HeightPercentage       float64 `json:"height_percentage"`
	RotationAngle          float64 `json:"rotation_angle"`
	CornerRadiusPercentage float64 `json:"corner_radius_percentage"`
}

const (
	StoryAreaTypeLocationType   StoryAreaTypeType = "location"
	StoryAreaTypeReactionType   StoryAreaTypeType = "suggested_reaction"
	StoryAreaTypeLinkType       StoryAreaTypeType = "link"
	StoryAreaTypeWeatherType    StoryAreaTypeType = "weather"
	StoryAreaTypeUniqueGiftType StoryAreaTypeType = "unique_gift"
)

type StoryAreaTypeType string
type StoryAreaType struct {
	Type StoryAreaTypeType `json:"type"`

	Latitude  *float64         `json:"latitude,omitempty"`
	Longitude *float64         `json:"longitude,omitempty"`
	Address   *LocationAddress `json:"address,omitempty"`

	ReactionType *ReactionType `json:"reaction_type,omitempty"`
	IsDark       *bool         `json:"is_dark,omitempty"`
	IsFlipped    *bool         `json:"is_flipped,omitempty"`

	URL *string `json:"url,omitempty"`

	Temperature     *float64 `json:"temperature,omitempty"`
	Emoji           *string  `json:"emoji"`
	BackgroundColor *int     `json:"background_color"`

	Name *string `json:"name,omitempty"`
}
type StoryArea struct {
	Position StoryAreaPosition `json:"position"`
	Type     StoryAreaType     `json:"type"`
}
