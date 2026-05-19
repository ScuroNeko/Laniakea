package tgapi

// BotCommand represents a bot command.
// Since: Bot API 4.7
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
// Since: Bot API 5.3
// See https://core.telegram.org/bots/api#botcommandscope
type BotCommandScope struct {
	Type   BotCommandScopeType `json:"type"`
	ChatID *int64              `json:"chat_id,omitempty"`
	UserID *int64              `json:"user_id,omitempty"`
}

// BotName represents the bot's name.
// Since: Bot API 6.7
type BotName struct {
	Name string `json:"name"`
}

// BotDescription represents the bot's description.
// Since: Bot API 6.6
type BotDescription struct {
	Description string `json:"description"`
}

// BotShortDescription represents the bot's short description.
// Since: Bot API 6.6
type BotShortDescription struct {
	ShortDescription string `json:"short_description"`
}

// InputProfilePhotoType indicates the type of a profile photo input.
type InputProfilePhotoType string

const (
	// InputProfilePhotoStaticType identifies a static profile photo input.
	InputProfilePhotoStaticType InputProfilePhotoType = "static"
	// InputProfilePhotoAnimatedType identifies an animated profile photo input.
	InputProfilePhotoAnimatedType InputProfilePhotoType = "animated"
)

// InputProfilePhoto describes a profile photo to set.
// Since: Bot API 9.0
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
	// MenuButtonCommandsType identifies a commands menu button.
	MenuButtonCommandsType MenuButtonType = "commands"
	// MenuButtonWebAppType identifies a web app menu button.
	MenuButtonWebAppType MenuButtonType = "web_app"
	// MenuButtonDefaultType identifies Telegram's default menu button.
	MenuButtonDefaultType MenuButtonType = "default"
)

// MenuButton represents a menu button.
// Since: Bot API 6.0
// See https://core.telegram.org/bots/api#menubutton
type MenuButton struct {
	Type MenuButtonType `json:"type"`

	// WebApp fields (for web_app button)
	Text   *string     `json:"text"`
	WebApp *WebAppInfo `json:"web_app"`
}

// BotAccessSettings describes access settings of a managed bot.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#botaccesssettings
type BotAccessSettings struct {
	AllowAllPrivateChats bool `json:"allow_all_private_chats"`
}
