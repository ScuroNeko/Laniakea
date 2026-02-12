package tgapi

type SetMyCommandsP struct {
	Commands []BotCommand     `json:"commands"`
	Scope    *BotCommandScope `json:"scope,omitempty"`
	Language string           `json:"language_code,omitempty"`
}

func (api *Api) SetMyCommands(params SetMyCommandsP) (bool, error) {
	req := NewRequest[bool]("setMyCommands", params)
	return req.Do(api)
}

type DeleteMyCommandsP struct {
	Scope    *BotCommandScope `json:"scope,omitempty"`
	Language string           `json:"language_code,omitempty"`
}

func (api *Api) DeleteMyCommands(params DeleteMyCommandsP) (bool, error) {
	req := NewRequest[bool]("deleteMyCommands", params)
	return req.Do(api)
}

type GetMyCommands struct {
	Scope    *BotCommandScope `json:"scope,omitempty"`
	Language string           `json:"language_code,omitempty"`
}

func (api *Api) GetMyCommands(params GetMyCommands) ([]BotCommand, error) {
	req := NewRequest[[]BotCommand]("getMyCommands", params)
	return req.Do(api)
}

type SetMyName struct {
	Name     string `json:"name"`
	Language string `json:"language_code,omitempty"`
}

func (api *Api) SetMyName(params SetMyName) (bool, error) {
	req := NewRequest[bool]("setMyName", params)
	return req.Do(api)
}

type GetMyName struct {
	Language string `json:"language_code,omitempty"`
}

func (api *Api) GetMyName(params GetMyName) (BotName, error) {
	req := NewRequest[BotName]("getMyName", params)
	return req.Do(api)
}

type SetMyDescription struct {
	Description string `json:"description"`
	Language    string `json:"language_code,omitempty"`
}

func (api *Api) SetMyDescription(params SetMyDescription) (bool, error) {
	req := NewRequest[bool]("setMyDescription", params)
	return req.Do(api)
}

type GetMyDescription struct {
	Language string `json:"language_code,omitempty"`
}

func (api *Api) GetMyDescription(params GetMyDescription) (BotDescription, error) {
	req := NewRequest[BotDescription]("getMyDescription", params)
	return req.Do(api)
}

type SetMyShortDescription struct {
	ShortDescription string `json:"short_description,omitempty"`
	Language         string `json:"language_code,omitempty"`
}

func (api *Api) SetMyShortDescription(params SetMyShortDescription) (bool, error) {
	req := NewRequest[bool]("setMyShortDescription", params)
	return req.Do(api)
}

type GetMyShortDescription struct {
	Language string `json:"language_code,omitempty"`
}

func (api *Api) GetMyShortDescription(params GetMyShortDescription) (BotShortDescription, error) {
	req := NewRequest[BotShortDescription]("getMyShortDescription", params)
	return req.Do(api)
}

type SetMyProfilePhotoP struct {
	Photo InputProfilePhoto `json:"photo"`
}

func (api *Api) SetMyProfilePhoto(params SetMyProfilePhotoP) (bool, error) {
	req := NewRequest[bool]("setMyProfilePhoto", params)
	return req.Do(api)
}
func (api *Api) RemoveMyProfilePhoto() (bool, error) {
	req := NewRequest[bool]("removeMyProfilePhoto", NoParams)
	return req.Do(api)
}

type SetChatMenuButtonP struct {
	ChatID     int            `json:"chat_id"`
	MenuButton MenuButtonType `json:"menu_button"`
}

func (api *Api) SetChatMenuButton(params SetChatMenuButtonP) (bool, error) {
	req := NewRequest[bool]("setChatMenuButton", params)
	return req.Do(api)
}

type GetChatMenuButtonP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) GetChatMenuButton(params GetChatMenuButtonP) (BaseMenuButton, error) {
	req := NewRequest[BaseMenuButton]("getChatMenuButton", params)
	return req.Do(api)
}

type SetMyDefaultAdministratorRightsP struct {
	Rights      *ChatAdministratorRights `json:"rights"`
	ForChannels bool                     `json:"for_channels"`
}

func (api *Api) SetMyDefaultAdministratorRights(params SetMyDefaultAdministratorRightsP) (bool, error) {
	req := NewRequest[bool]("setMyDefaultAdministratorRights", params)
	return req.Do(api)
}

type GetMyDefaultAdministratorRightsP struct {
	ForChannels bool `json:"for_channels"`
}

func (api *Api) GetMyDefaultAdministratorRights(params GetMyDefaultAdministratorRightsP) (ChatAdministratorRights, error) {
	req := NewRequest[ChatAdministratorRights]("getMyDefaultAdministratorRights", params)
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

func (api *Api) SendGift(params SendGiftP) (bool, error) {
	req := NewRequest[bool]("sendGift", params)
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

func (api *Api) GiftPremiumSubscription(params GiftPremiumSubscriptionP) (bool, error) {
	req := NewRequest[bool]("giftPremiumSubscription", params)
	return req.Do(api)
}
