package tgapi

// SetMyCommandsP holds parameters for the setMyCommands method.
// See https://core.telegram.org/bots/api#setmycommands
type SetMyCommandsP struct {
	Commands []BotCommand     `json:"commands"`
	Scope    *BotCommandScope `json:"scope,omitempty"`
	Language string           `json:"language_code,omitempty"`
}

// SetMyCommands changes the list of the bot's commands.
// Returns true on success.
// See https://core.telegram.org/bots/api#setmycommands
func (api *API) SetMyCommands(params SetMyCommandsP) (bool, error) {
	req := NewRequest[bool]("setMyCommands", params)
	return req.Do(api)
}

// DeleteMyCommandsP holds parameters for the deleteMyCommands method.
// See https://core.telegram.org/bots/api#deletemycommands
type DeleteMyCommandsP struct {
	Scope    *BotCommandScope `json:"scope,omitempty"`
	Language string           `json:"language_code,omitempty"`
}

// DeleteMyCommands deletes the list of the bot's commands for the given scope and user language.
// Returns true on success.
// See https://core.telegram.org/bots/api#deletemycommands
func (api *API) DeleteMyCommands(params DeleteMyCommandsP) (bool, error) {
	req := NewRequest[bool]("deleteMyCommands", params)
	return req.Do(api)
}

// GetMyCommands holds parameters for the getMyCommands method.
// See https://core.telegram.org/bots/api#getmycommands
type GetMyCommands struct {
	Scope    *BotCommandScope `json:"scope,omitempty"`
	Language string           `json:"language_code,omitempty"`
}

// GetMyCommands returns the current list of the bot's commands for the given scope and user language.
// See https://core.telegram.org/bots/api#getmycommands
func (api *API) GetMyCommands(params GetMyCommands) ([]BotCommand, error) {
	req := NewRequest[[]BotCommand]("getMyCommands", params)
	return req.Do(api)
}

// SetMyName holds parameters for the setMyName method.
// See https://core.telegram.org/bots/api#setmyname
type SetMyName struct {
	Name     string `json:"name"`
	Language string `json:"language_code,omitempty"`
}

// SetMyName changes the bot's name.
// Returns true on success.
// See https://core.telegram.org/bots/api#setmyname
func (api *API) SetMyName(params SetMyName) (bool, error) {
	req := NewRequest[bool]("setMyName", params)
	return req.Do(api)
}

// GetMyName holds parameters for the getMyName method.
// See https://core.telegram.org/bots/api#getmyname
type GetMyName struct {
	Language string `json:"language_code,omitempty"`
}

// GetMyName returns the bot's name for the given language.
// See https://core.telegram.org/bots/api#getmyname
func (api *API) GetMyName(params GetMyName) (BotName, error) {
	req := NewRequest[BotName]("getMyName", params)
	return req.Do(api)
}

// SetMyDescription holds parameters for the setMyDescription method.
// See https://core.telegram.org/bots/api#setmydescription
type SetMyDescription struct {
	Description string `json:"description"`
	Language    string `json:"language_code,omitempty"`
}

// SetMyDescription changes the bot's description.
// Returns true on success.
// See https://core.telegram.org/bots/api#setmydescription
func (api *API) SetMyDescription(params SetMyDescription) (bool, error) {
	req := NewRequest[bool]("setMyDescription", params)
	return req.Do(api)
}

// GetMyDescription holds parameters for the getMyDescription method.
// See https://core.telegram.org/bots/api#getmydescription
type GetMyDescription struct {
	Language string `json:"language_code,omitempty"`
}

// GetMyDescription returns the bot's description for the given language.
// See https://core.telegram.org/bots/api#getmydescription
func (api *API) GetMyDescription(params GetMyDescription) (BotDescription, error) {
	req := NewRequest[BotDescription]("getMyDescription", params)
	return req.Do(api)
}

// SetMyShortDescription holds parameters for the setMyShortDescription method.
// See https://core.telegram.org/bots/api#setmyshortdescription
type SetMyShortDescription struct {
	ShortDescription string `json:"short_description,omitempty"`
	Language         string `json:"language_code,omitempty"`
}

// SetMyShortDescription changes the bot's short description.
// Returns true on success.
// See https://core.telegram.org/bots/api#setmyshortdescription
func (api *API) SetMyShortDescription(params SetMyShortDescription) (bool, error) {
	req := NewRequest[bool]("setMyShortDescription", params)
	return req.Do(api)
}

// GetMyShortDescription holds parameters for the getMyShortDescription method.
// See https://core.telegram.org/bots/api#getmyshortdescription
type GetMyShortDescription struct {
	Language string `json:"language_code,omitempty"`
}

// GetMyShortDescription returns the bot's short description for the given language.
// See https://core.telegram.org/bots/api#getmyshortdescription
func (api *API) GetMyShortDescription(params GetMyShortDescription) (BotShortDescription, error) {
	req := NewRequest[BotShortDescription]("getMyShortDescription", params)
	return req.Do(api)
}

// SetMyProfilePhotoP holds parameters for the setMyProfilePhoto method.
// See https://core.telegram.org/bots/api#setmyprofilephoto
type SetMyProfilePhotoP struct {
	Photo InputProfilePhoto `json:"photo"`
}

// SetMyProfilePhoto changes the bot's profile photo.
// Returns true on success.
// See https://core.telegram.org/bots/api#setmyprofilephoto
func (api *API) SetMyProfilePhoto(params SetMyProfilePhotoP) (bool, error) {
	req := NewRequest[bool]("setMyProfilePhoto", params)
	return req.Do(api)
}

// RemoveMyProfilePhoto removes the bot's profile photo.
// Returns true on success.
// See https://core.telegram.org/bots/api#removemyprofilephoto
func (api *API) RemoveMyProfilePhoto() (bool, error) {
	req := NewRequest[bool]("removeMyProfilePhoto", NoParams)
	return req.Do(api)
}

// SetChatMenuButtonP holds parameters for the setChatMenuButton method.
// See https://core.telegram.org/bots/api#setchatmenubutton
type SetChatMenuButtonP struct {
	ChatID     int            `json:"chat_id"`
	MenuButton MenuButtonType `json:"menu_button"`
}

// SetChatMenuButton changes the menu button for a given chat or the default menu button.
// Returns true on success.
// See https://core.telegram.org/bots/api#setchatmenubutton
func (api *API) SetChatMenuButton(params SetChatMenuButtonP) (bool, error) {
	req := NewRequest[bool]("setChatMenuButton", params)
	return req.Do(api)
}

// GetChatMenuButtonP holds parameters for the getChatMenuButton method.
// See https://core.telegram.org/bots/api#getchatmenubutton
type GetChatMenuButtonP struct {
	ChatID int `json:"chat_id"`
}

// GetChatMenuButton returns the current menu button for the given chat.
// See https://core.telegram.org/bots/api#getchatmenubutton
func (api *API) GetChatMenuButton(params GetChatMenuButtonP) (BaseMenuButton, error) {
	req := NewRequest[BaseMenuButton]("getChatMenuButton", params)
	return req.Do(api)
}

// SetMyDefaultAdministratorRightsP holds parameters for the setMyDefaultAdministratorRights method.
// See https://core.telegram.org/bots/api#setmydefaultadministratorrights
type SetMyDefaultAdministratorRightsP struct {
	Rights      *ChatAdministratorRights `json:"rights"`
	ForChannels bool                     `json:"for_channels"`
}

// SetMyDefaultAdministratorRights changes the default administrator rights for the bot.
// Returns true on success.
// See https://core.telegram.org/bots/api#setmydefaultadministratorrights
func (api *API) SetMyDefaultAdministratorRights(params SetMyDefaultAdministratorRightsP) (bool, error) {
	req := NewRequest[bool]("setMyDefaultAdministratorRights", params)
	return req.Do(api)
}

// GetMyDefaultAdministratorRightsP holds parameters for the getMyDefaultAdministratorRights method.
// See https://core.telegram.org/bots/api#getmydefaultadministratorrights
type GetMyDefaultAdministratorRightsP struct {
	ForChannels bool `json:"for_channels"`
}

// GetMyDefaultAdministratorRights returns the current default administrator rights for the bot.
// See https://core.telegram.org/bots/api#getmydefaultadministratorrights
func (api *API) GetMyDefaultAdministratorRights(params GetMyDefaultAdministratorRightsP) (ChatAdministratorRights, error) {
	req := NewRequest[ChatAdministratorRights]("getMyDefaultAdministratorRights", params)
	return req.Do(api)
}

// GetAvailableGifts returns the list of gifts that can be sent by the bot.
// See https://core.telegram.org/bots/api#getavailablegifts
func (api *API) GetAvailableGifts() (Gifts, error) {
	req := NewRequest[Gifts]("getAvailableGifts", NoParams)
	return req.Do(api)
}

// SendGiftP holds parameters for the sendGift method.
// See https://core.telegram.org/bots/api#sendgift
type SendGiftP struct {
	UserID        int             `json:"user_id,omitempty"`
	ChatID        int             `json:"chat_id,omitempty"`
	GiftID        string          `json:"gift_id"`
	PayForUpgrade bool            `json:"pay_for_upgrade"`
	Text          string          `json:"text"`
	TextParseMode ParseMode       `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity `json:"text_entities,omitempty"`
}

// SendGift sends a gift to the given user or chat.
// Returns true on success.
// See https://core.telegram.org/bots/api#sendgift
func (api *API) SendGift(params SendGiftP) (bool, error) {
	req := NewRequest[bool]("sendGift", params)
	return req.Do(api)
}

// GiftPremiumSubscriptionP holds parameters for the giftPremiumSubscription method.
// See https://core.telegram.org/bots/api#giftpremiumsubscription
type GiftPremiumSubscriptionP struct {
	UserID        int             `json:"user_id"`
	MonthCount    int             `json:"month_count"`
	StarCount     int             `json:"star_count"`
	Text          string          `json:"text,omitempty"`
	TextParseMode ParseMode       `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity `json:"text_entities,omitempty"`
}

// GiftPremiumSubscription gifts a Telegram Premium subscription to the user.
// Returns true on success.
// See https://core.telegram.org/bots/api#giftpremiumsubscription
func (api *API) GiftPremiumSubscription(params GiftPremiumSubscriptionP) (bool, error) {
	req := NewRequest[bool]("giftPremiumSubscription", params)
	return req.Do(api)
}
