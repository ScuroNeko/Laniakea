package tgapi

// ForumTopic represents a forum topic.
// Since: Bot API 6.3
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
// Since: Bot API 6.3
// See https://core.telegram.org/bots/api#forumtopiciconcolor
type ForumTopicIconColor int

const (
	// ForumTopicIconColorBlue is the blue color for forum topic icons (value 7322096).
	ForumTopicIconColorBlue ForumTopicIconColor = 7322096
)

// ForumTopicCreated represents a service message about a new forum topic created.
// Since: Bot API 6.3
type ForumTopicCreated struct {
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
	IsNameImplicit    bool   `json:"is_name_implicit,omitempty"`
}

// ForumTopicEdited represents a service message about an edited forum topic.
// Since: Bot API 6.4
type ForumTopicEdited struct {
	Name              string `json:"name,omitempty"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicClosed represents a service message about a forum topic closed.
// Since: Bot API 6.3
type ForumTopicClosed struct{}

// ForumTopicReopened represents a service message about a forum topic reopened.
// Since: Bot API 6.3
type ForumTopicReopened struct{}

// GeneralForumTopicHidden represents a service message about the General forum topic hidden.
// Since: Bot API 6.4
type GeneralForumTopicHidden struct{}

// GeneralForumTopicUnhidden represents a service message about the General forum topic unhidden.
// Since: Bot API 6.4
type GeneralForumTopicUnhidden struct{}
