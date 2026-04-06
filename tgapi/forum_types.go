package tgapi

// ForumTopic represents a forum topic.
// See https://core.telegram.org/bots/api#forumtopic
type ForumTopic struct {
	MessageThreadID   int    `json:"message_thread_id"`
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
	IsNameImplicit    bool   `json:"is_name_implicit,omitempty"`
}

// ForumTopicIconColor represents the color of a forum topic icon.
// The value is an integer representing the color in RGB format.
// See https://core.telegram.org/bots/api#forumtopiciconcolor
type ForumTopicIconColor int

const (
	// ForumTopicIconColorBlue is the blue color for forum topic icons (value 7322096).
	ForumTopicIconColorBlue ForumTopicIconColor = 7322096
)

type ForumTopicCreated struct {
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
	IsNameImplicit    bool   `json:"is_name_implicit,omitempty"`
}
type ForumTopicEdited struct {
	Name              string `json:"name,omitempty"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}
type ForumTopicClosed struct{}
type ForumTopicReopened struct{}
type GeneralForumTopicHidden struct{}
type GeneralForumTopicUnhidden struct {
}
