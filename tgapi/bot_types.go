package tgapi

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}
type BotCommandScopeType string

const (
	BotCommandScopeDefaultType               BotCommandScopeType = "default"
	BotCommandScopePrivateType               BotCommandScopeType = "all_private_chats"
	BotCommandScopeGroupType                 BotCommandScopeType = "all_group_chats"
	BotCommandScopeAllChatAdministratorsType BotCommandScopeType = "all_chat_administrators"
	BotCommandScopeChatType                  BotCommandScopeType = "chat"
	BotCommandScopeChatAdministratorsType    BotCommandScopeType = "chat_administrators"
	BotCommandScopeChatMemberType            BotCommandScopeType = "chat_member"
)

type BotCommandScope struct {
	Type   BotCommandScopeType `json:"type"`
	ChatID *int                `json:"chat_id,omitempty"`
	UserID *int                `json:"user_id,omitempty"`
}

type BotName struct {
	Name string `json:"name"`
}
type BotDescription struct {
	Description string `json:"description"`
}
type BotShortDescription struct {
	ShortDescription string `json:"short_description"`
}

const (
	InputProfilePhotoStaticType   InputProfilePhotoType = "static"
	InputProfilePhotoAnimatedType InputProfilePhotoType = "animated"
)

type InputProfilePhotoType string
type InputProfilePhoto struct {
	Type InputProfilePhotoType `json:"type"`

	// Static
	Photo *string `json:"photo,omitempty"`

	// Animated
	Animation          *string  `json:"animation,omitempty"`
	MainFrameTimestamp *float64 `json:"main_frame_timestamp,omitempty"`
}

const (
	MenuButtonCommandsType MenuButtonType = "commands"
	MenuButtonWebAppType   MenuButtonType = "web_app"
	MenuButtonDefaultType  MenuButtonType = "default"
)

type MenuButtonType string
type BaseMenuButton struct {
	Type MenuButtonType `json:"type"`
	// WebApp
	Text   string     `json:"text"`
	WebApp WebAppInfo `json:"web_app"`
}
