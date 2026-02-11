package tgapi

type SetMyCommandsP struct {
	Commands []BotCommand         `json:"commands"`
	Scope    *BaseBotCommandScope `json:"scope,omitempty"`
	Language string               `json:"language_code,omitempty"`
}

func (api *Api) SetMyCommands(p SetMyCommandsP) (bool, error) {
	req := NewRequest[bool]("setMyCommands", p)
	return req.Do(api)
}

type DeleteMyCommandsP struct {
	Scope    *BaseBotCommandScope `json:"scope,omitempty"`
	Language string               `json:"language_code,omitempty"`
}

func (api *Api) DeleteMyCommands(p DeleteMyCommandsP) (bool, error) {
	req := NewRequest[bool]("deleteMyCommands", p)
	return req.Do(api)
}

type GetMyCommands struct {
	Scope    *BaseBotCommandScope `json:"scope,omitempty"`
	Language string               `json:"language_code,omitempty"`
}

func (api *Api) GetMyCommands(p GetMyCommands) ([]BotCommand, error) {
	req := NewRequest[[]BotCommand]("getMyCommands", p)
	return req.Do(api)
}

type SetMyName struct {
	Name     string `json:"name"`
	Language string `json:"language_code,omitempty"`
}

func (api *Api) SetMyName(p SetMyName) (bool, error) {
	req := NewRequest[bool]("setMyName", p)
	return req.Do(api)
}

type GetMyName struct {
	Language string `json:"language_code,omitempty"`
}

func (api *Api) GetMyName(p GetMyName) (BotName, error) {
	req := NewRequest[BotName]("getMyName", p)
	return req.Do(api)
}

type SetMyDescription struct {
	Description string `json:"description"`
	Language    string `json:"language_code,omitempty"`
}

func (api *Api) SetMyDescription(p SetMyDescription) (bool, error) {
	req := NewRequest[bool]("setMyDescription", p)
	return req.Do(api)
}

type GetMyDescription struct {
	Language string `json:"language_code,omitempty"`
}

func (api *Api) GetMyDescription(p GetMyDescription) (BotDescription, error) {
	req := NewRequest[BotDescription]("getMyDescription", p)
	return req.Do(api)
}

type SetMyShortDescription struct {
	ShortDescription string `json:"short_description,omitempty"`
	Language         string `json:"language_code,omitempty"`
}

func (api *Api) SetMyShortDescription(p SetMyShortDescription) (bool, error) {
	req := NewRequest[bool]("setMyShortDescription", p)
	return req.Do(api)
}

type GetMyShortDescription struct {
	Language string `json:"language_code,omitempty"`
}

func (api *Api) GetMyShortDescription(p GetMyShortDescription) (BotShortDescription, error) {
	req := NewRequest[BotShortDescription]("getMyShortDescription", p)
	return req.Do(api)
}

type SetMyProfilePhotoP[T InputProfilePhoto] struct {
	Photo T `json:"photo"`
}

func (api *Api) SetMyProfilePhotoStatic(p SetMyProfilePhotoP[InputProfilePhotoStatic]) (bool, error) {
	req := NewRequest[bool]("setMyProfilePhoto", p)
	return req.Do(api)
}
func (api *Api) SetMyProfilePhotoAnimated(p SetMyProfilePhotoP[InputProfilePhotoAnimated]) (bool, error) {
	req := NewRequest[bool]("setMyProfilePhoto", p)
	return req.Do(api)
}
func (api *Api) RemoveMyProfilePhoto() (bool, error) {
	req := NewRequest[bool]("removeMyProfilePhoto", NoParams)
	return req.Do(api)
}

type SetChatMenuButtonP[T MenuButton] struct {
	ChatID     int `json:"chat_id"`
	MenuButton T   `json:"menu_button"`
}

func (api *Api) SetChatMenuButtonCommands(p SetChatMenuButtonP[MenuButtonCommands]) (bool, error) {
	req := NewRequest[bool]("setChatMenuButton", p)
	return req.Do(api)
}
func (api *Api) SetChatMenuButtonWebApp(p SetChatMenuButtonP[MenuButtonWebApp]) (bool, error) {
	req := NewRequest[bool]("setChatMenuButton", p)
	return req.Do(api)
}
func (api *Api) SetChatMenuButtonDefault(p SetChatMenuButtonP[MenuButtonDefault]) (bool, error) {
	req := NewRequest[bool]("setChatMenuButton", p)
	return req.Do(api)
}

type GetChatMenuButtonP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) GetChatMenuButton(p GetChatMenuButtonP) (BaseMenuButton, error) {
	req := NewRequest[BaseMenuButton]("getChatMenuButton", p)
	return req.Do(api)
}

type SetMyDefaultAdministratorRightsP struct {
	Rights      *ChatAdministratorRights `json:"rights"`
	ForChannels bool                     `json:"for_channels"`
}

func (api *Api) SetMyDefaultAdministratorRights(p SetMyDefaultAdministratorRightsP) (bool, error) {
	req := NewRequest[bool]("setMyDefaultAdministratorRights", p)
	return req.Do(api)
}

type GetMyDefaultAdministratorRightsP struct {
	ForChannels bool `json:"for_channels"`
}

func (api *Api) GetMyDefaultAdministratorRights(p GetMyDefaultAdministratorRightsP) (ChatAdministratorRights, error) {
	req := NewRequest[ChatAdministratorRights]("getMyDefaultAdministratorRights", p)
	return req.Do(api)
}

func (api *Api) GetAvailableGifts() (Gifts, error) {
	req := NewRequest[Gifts]("getAvailableGifts", NoParams)
	return req.Do(api)
}

type SendGiftP struct {
	UserID        int             `json:"user_id,omitempty"`
	ChatID        int             `json:"chat_id,omitempty"`
	GiftID        string          `json:"gift_id"`
	PayForUpgrade bool            `json:"pay_for_upgrade"`
	Text          string          `json:"text"`
	TextParseMode ParseMode       `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity `json:"text_entities,omitempty"`
}

func (api *Api) SendGift(p SendGiftP) (bool, error) {
	req := NewRequest[bool]("sendGift", p)
	return req.Do(api)
}

type GiftPremiumSubscriptionP struct {
	UserID        int             `json:"user_id"`
	MonthCount    int             `json:"month_count"`
	StarCount     int             `json:"star_count"`
	Text          string          `json:"text,omitempty"`
	TextParseMode ParseMode       `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity `json:"text_entities,omitempty"`
}

func (api *Api) GiftPremiumSubscription(p GiftPremiumSubscriptionP) (bool, error) {
	req := NewRequest[bool]("giftPremiumSubscription", p)
	return req.Do(api)
}
