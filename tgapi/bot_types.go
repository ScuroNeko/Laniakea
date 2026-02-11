package tgapi

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}
type BotCommandScopeType string

const (
	BotCommandScopeDefaultType               BotCommandScopeType = "default"
	BotCommandScopePrivateType               BotCommandScopeType = "all_private_chats"
	BotCommandScopeGroupType                 BotCommandScopeType = "all_groups_chats"
	BotCommandScopeAllChatAdministratorsType BotCommandScopeType = "all_chat_administrators"
	BotCommandScopeChatType                  BotCommandScopeType = "chat"
	BotCommandScopeChatAdministratorsType    BotCommandScopeType = "chat_administrators"
	BotCommandScopeChatMemberType            BotCommandScopeType = "chat_member"
)

type BaseBotCommandScope struct {
	Type BotCommandScopeType `json:"type"`
}
type BotCommandScopeDefault struct {
	BaseBotCommandScope
}
type BotCommandScopePrivateChats struct {
	BaseBotCommandScope
}
type BotCommandScopeAllGroupChats struct {
	BaseBotCommandScope
}
type BotCommandScopeAllChatAdministrators struct {
	BaseBotCommandScope
}
type BotCommandScopeChat struct {
	BaseBotCommandScope
	ChatID int `json:"chat_id"`
}
type BotCommandScopeChatAdministrators struct {
	BaseBotCommandScope
	ChatID int `json:"chat_id"`
}
type BotCommandScopeChatMember struct {
	BaseBotCommandScope
	ChatID int `json:"chat_id"`
	UserID int `json:"user_id"`
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
type InputProfilePhoto interface {
	InputProfilePhotoStatic | InputProfilePhotoAnimated
}
type BaseInputProfilePhoto struct {
	Type InputProfilePhotoType `json:"type"`
}
type InputProfilePhotoStatic struct {
	BaseInputProfilePhoto
	Photo string `json:"photo"`
}
type InputProfilePhotoAnimated struct {
	BaseInputProfilePhoto
	Animation          string  `json:"animation"`
	MainFrameTimestamp float64 `json:"main_frame_timestamp"`
}

const (
	MenuButtonCommandsType MenuButtonType = "commands"
	MenuButtonWebAppType   MenuButtonType = "web_app"
	MenuButtonDefaultType  MenuButtonType = "default"
)

type MenuButtonType string
type MenuButton interface {
	MenuButtonCommands | MenuButtonDefault | MenuButtonWebApp
}
type BaseMenuButton struct {
	Type MenuButtonType `json:"type"`
}
type MenuButtonCommands struct{ BaseMenuButton }
type MenuButtonDefault struct{ BaseMenuButton }
type MenuButtonWebApp struct {
	BaseMenuButton
	Text   string     `json:"text"`
	WebApp WebAppInfo `json:"web_app"`
}
