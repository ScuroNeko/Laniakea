package tgapi

// BotCommand represents a bot command.
// See https://core.telegram.org/bots/api#botcommand
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// BotCommandScopeType indicates the type of a command scope.
type BotCommandScopeType string

const (
	// BotCommandScopeDefaultType is the default command scope.
	BotCommandScopeDefaultType BotCommandScopeType = "default"
	// BotCommandScopePrivateType covers all private chats.
	BotCommandScopePrivateType BotCommandScopeType = "all_private_chats"
	// BotCommandScopeGroupType covers all group and supergroup chats.
	BotCommandScopeGroupType BotCommandScopeType = "all_group_chats"
	// BotCommandScopeAllChatAdministratorsType covers all chat administrators.
	BotCommandScopeAllChatAdministratorsType BotCommandScopeType = "all_chat_administrators"
	// BotCommandScopeChatType covers a specific chat.
	BotCommandScopeChatType BotCommandScopeType = "chat"
	// BotCommandScopeChatAdministratorsType covers administrators of a specific chat.
	BotCommandScopeChatAdministratorsType BotCommandScopeType = "chat_administrators"
	// BotCommandScopeChatMemberType covers a specific member of a specific chat.
	BotCommandScopeChatMemberType BotCommandScopeType = "chat_member"
)

// BotCommandScope represents the scope to which bot commands are applied.
// See https://core.telegram.org/bots/api#botcommandscope
type BotCommandScope struct {
	Type   BotCommandScopeType `json:"type"`
	ChatID *int64              `json:"chat_id,omitempty"`
	UserID *int64              `json:"user_id,omitempty"`
}

// BotName represents the bot's name.
type BotName struct {
	Name string `json:"name"`
}

// BotDescription represents the bot's description.
type BotDescription struct {
	Description string `json:"description"`
}

// BotShortDescription represents the bot's short description.
type BotShortDescription struct {
	ShortDescription string `json:"short_description"`
}

// InputProfilePhotoType indicates the type of a profile photo input.
type InputProfilePhotoType string

const (
	InputProfilePhotoStaticType   InputProfilePhotoType = "static"
	InputProfilePhotoAnimatedType InputProfilePhotoType = "animated"
)

// InputProfilePhoto describes a profile photo to set.
// See https://core.telegram.org/bots/api#inputprofilephoto
type InputProfilePhoto struct {
	Type InputProfilePhotoType `json:"type"`

	// Static fields (for static photos)
	Photo *string `json:"photo,omitempty"`

	// Animated fields (for animated profile videos)
	Animation          *string  `json:"animation,omitempty"`
	MainFrameTimestamp *float64 `json:"main_frame_timestamp,omitempty"`
}

// MenuButtonType indicates the type of a menu button.
type MenuButtonType string

const (
	MenuButtonCommandsType MenuButtonType = "commands"
	MenuButtonWebAppType   MenuButtonType = "web_app"
	MenuButtonDefaultType  MenuButtonType = "default"
)

// BaseMenuButton represents a menu button.
// See https://core.telegram.org/bots/api#menubutton
type BaseMenuButton struct {
	Type MenuButtonType `json:"type"`

	// WebApp fields (for web_app button)
	Text   string     `json:"text"`
	WebApp WebAppInfo `json:"web_app"`
}
