package tgapi

import "context"

// SetMyCommands holds parameters for the setMyCommands method.
// Since: Bot API 4.7
// See https://core.telegram.org/bots/api#setmycommands
type SetMyCommands struct {
	Commands []BotCommand     `json:"commands"`
	Scope    *BotCommandScope `json:"scope,omitempty"`
	Language string           `json:"language_code,omitempty"`
}

// SetMyCommands changes the list of the bot's commands.
// Since: Bot API 4.7
// Returns true on success.
// See https://core.telegram.org/bots/api#setmycommands
func (api *API) SetMyCommands(params SetMyCommands) (bool, error) {
	req := NewRequest[bool]("setMyCommands", params)
	return req.Do(api)
}

// SetMyCommandsWithContext is the context-aware variant of SetMyCommands.
// Since: Bot API 4.7
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setmycommands
func (api *API) SetMyCommandsWithContext(ctx context.Context, params SetMyCommands) (bool, error) {
	req := NewRequest[bool]("setMyCommands", params)
	return req.DoWithContext(ctx, api)
}

// DeleteMyCommands holds parameters for the deleteMyCommands method.
// Since: Bot API 5.3
// See https://core.telegram.org/bots/api#deletemycommands
type DeleteMyCommands struct {
	Scope    *BotCommandScope `json:"scope,omitempty"`
	Language string           `json:"language_code,omitempty"`
}

// DeleteMyCommands deletes the list of the bot's commands for the given scope and user language.
// Since: Bot API 5.3
// Returns true on success.
// See https://core.telegram.org/bots/api#deletemycommands
func (api *API) DeleteMyCommands(params DeleteMyCommands) (bool, error) {
	req := NewRequest[bool]("deleteMyCommands", params)
	return req.Do(api)
}

// DeleteMyCommandsWithContext is the context-aware variant of DeleteMyCommands.
// Since: Bot API 5.3
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletemycommands
func (api *API) DeleteMyCommandsWithContext(ctx context.Context, params DeleteMyCommands) (bool, error) {
	req := NewRequest[bool]("deleteMyCommands", params)
	return req.DoWithContext(ctx, api)
}

// GetMyCommands holds parameters for the getMyCommands method.
// Since: Bot API 4.7
// See https://core.telegram.org/bots/api#getmycommands
type GetMyCommands struct {
	Scope    *BotCommandScope `json:"scope,omitempty"`
	Language string           `json:"language_code,omitempty"`
}

// GetMyCommands returns the current list of the bot's commands for the given scope and user language.
// Since: Bot API 4.7
// See https://core.telegram.org/bots/api#getmycommands
func (api *API) GetMyCommands(params GetMyCommands) ([]BotCommand, error) {
	req := NewRequest[[]BotCommand]("getMyCommands", params)
	return req.Do(api)
}

// GetMyCommandsWithContext is the context-aware variant of GetMyCommands.
// Since: Bot API 4.7
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getmycommands
func (api *API) GetMyCommandsWithContext(ctx context.Context, params GetMyCommands) ([]BotCommand, error) {
	req := NewRequest[[]BotCommand]("getMyCommands", params)
	return req.DoWithContext(ctx, api)
}

// SetMyName holds parameters for the setMyName method.
// Since: Bot API 6.7
// See https://core.telegram.org/bots/api#setmyname
type SetMyName struct {
	Name     string `json:"name"`
	Language string `json:"language_code,omitempty"`
}

// SetMyName changes the bot's name.
// Since: Bot API 6.7
// Returns true on success.
// See https://core.telegram.org/bots/api#setmyname
func (api *API) SetMyName(params SetMyName) (bool, error) {
	req := NewRequest[bool]("setMyName", params)
	return req.Do(api)
}

// SetMyNameWithContext is the context-aware variant of SetMyName.
// Since: Bot API 6.7
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setmyname
func (api *API) SetMyNameWithContext(ctx context.Context, params SetMyName) (bool, error) {
	req := NewRequest[bool]("setMyName", params)
	return req.DoWithContext(ctx, api)
}

// GetMyName holds parameters for the getMyName method.
// Since: Bot API 6.7
// See https://core.telegram.org/bots/api#getmyname
type GetMyName struct {
	Language string `json:"language_code,omitempty"`
}

// GetMyName returns the bot's name for the given language.
// Since: Bot API 6.7
// See https://core.telegram.org/bots/api#getmyname
func (api *API) GetMyName(params GetMyName) (BotName, error) {
	req := NewRequest[BotName]("getMyName", params)
	return req.Do(api)
}

// GetMyNameWithContext is the context-aware variant of GetMyName.
// Since: Bot API 6.7
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getmyname
func (api *API) GetMyNameWithContext(ctx context.Context, params GetMyName) (BotName, error) {
	req := NewRequest[BotName]("getMyName", params)
	return req.DoWithContext(ctx, api)
}

// SetMyDescription holds parameters for the setMyDescription method.
// Since: Bot API 6.6
// See https://core.telegram.org/bots/api#setmydescription
type SetMyDescription struct {
	Description string `json:"description"`
	Language    string `json:"language_code,omitempty"`
}

// SetMyDescription changes the bot's description.
// Since: Bot API 6.6
// Returns true on success.
// See https://core.telegram.org/bots/api#setmydescription
func (api *API) SetMyDescription(params SetMyDescription) (bool, error) {
	req := NewRequest[bool]("setMyDescription", params)
	return req.Do(api)
}

// SetMyDescriptionWithContext is the context-aware variant of SetMyDescription.
// Since: Bot API 6.6
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setmydescription
func (api *API) SetMyDescriptionWithContext(ctx context.Context, params SetMyDescription) (bool, error) {
	req := NewRequest[bool]("setMyDescription", params)
	return req.DoWithContext(ctx, api)
}

// GetMyDescription holds parameters for the getMyDescription method.
// Since: Bot API 6.6
// See https://core.telegram.org/bots/api#getmydescription
type GetMyDescription struct {
	Language string `json:"language_code,omitempty"`
}

// GetMyDescription returns the bot's description for the given language.
// Since: Bot API 6.6
// See https://core.telegram.org/bots/api#getmydescription
func (api *API) GetMyDescription(params GetMyDescription) (BotDescription, error) {
	req := NewRequest[BotDescription]("getMyDescription", params)
	return req.Do(api)
}

// GetMyDescriptionWithContext is the context-aware variant of GetMyDescription.
// Since: Bot API 6.6
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getmydescription
func (api *API) GetMyDescriptionWithContext(ctx context.Context, params GetMyDescription) (BotDescription, error) {
	req := NewRequest[BotDescription]("getMyDescription", params)
	return req.DoWithContext(ctx, api)
}

// SetMyShortDescription holds parameters for the setMyShortDescription method.
// Since: Bot API 6.6
// See https://core.telegram.org/bots/api#setmyshortdescription
type SetMyShortDescription struct {
	ShortDescription string `json:"short_description,omitempty"`
	Language         string `json:"language_code,omitempty"`
}

// SetMyShortDescription changes the bot's short description.
// Since: Bot API 6.6
// Returns true on success.
// See https://core.telegram.org/bots/api#setmyshortdescription
func (api *API) SetMyShortDescription(params SetMyShortDescription) (bool, error) {
	req := NewRequest[bool]("setMyShortDescription", params)
	return req.Do(api)
}

// SetMyShortDescriptionWithContext is the context-aware variant of SetMyShortDescription.
// Since: Bot API 6.6
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setmyshortdescription
func (api *API) SetMyShortDescriptionWithContext(ctx context.Context, params SetMyShortDescription) (bool, error) {
	req := NewRequest[bool]("setMyShortDescription", params)
	return req.DoWithContext(ctx, api)
}

// GetMyShortDescription holds parameters for the getMyShortDescription method.
// Since: Bot API 6.6
// See https://core.telegram.org/bots/api#getmyshortdescription
type GetMyShortDescription struct {
	Language string `json:"language_code,omitempty"`
}

// GetMyShortDescription returns the bot's short description for the given language.
// Since: Bot API 6.6
// See https://core.telegram.org/bots/api#getmyshortdescription
func (api *API) GetMyShortDescription(params GetMyShortDescription) (BotShortDescription, error) {
	req := NewRequest[BotShortDescription]("getMyShortDescription", params)
	return req.Do(api)
}

// GetMyShortDescriptionWithContext is the context-aware variant of GetMyShortDescription.
// Since: Bot API 6.6
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getmyshortdescription
func (api *API) GetMyShortDescriptionWithContext(ctx context.Context, params GetMyShortDescription) (BotShortDescription, error) {
	req := NewRequest[BotShortDescription]("getMyShortDescription", params)
	return req.DoWithContext(ctx, api)
}

// SetMyProfilePhoto holds parameters for the setMyProfilePhoto method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#setmyprofilephoto
type SetMyProfilePhoto struct {
	Photo InputProfilePhoto `json:"photo"`
}

// SetMyProfilePhoto changes the bot's profile photo.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#setmyprofilephoto
func (api *API) SetMyProfilePhoto(params SetMyProfilePhoto) (bool, error) {
	req := NewRequest[bool]("setMyProfilePhoto", params)
	return req.Do(api)
}

// SetMyProfilePhotoWithContext is the context-aware variant of SetMyProfilePhoto.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setmyprofilephoto
func (api *API) SetMyProfilePhotoWithContext(ctx context.Context, params SetMyProfilePhoto) (bool, error) {
	req := NewRequest[bool]("setMyProfilePhoto", params)
	return req.DoWithContext(ctx, api)
}

// RemoveMyProfilePhoto removes the bot's profile photo.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#removemyprofilephoto
func (api *API) RemoveMyProfilePhoto() (bool, error) {
	req := NewRequest[bool]("removeMyProfilePhoto", NoParams)
	return req.Do(api)
}

// RemoveMyProfilePhotoWithContext is the context-aware variant of RemoveMyProfilePhoto.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#removemyprofilephoto
func (api *API) RemoveMyProfilePhotoWithContext(ctx context.Context) (bool, error) {
	req := NewRequest[bool]("removeMyProfilePhoto", NoParams)
	return req.DoWithContext(ctx, api)
}

// SetChatMenuButton holds parameters for the setChatMenuButton method.
// Since: Bot API 6.0
// See https://core.telegram.org/bots/api#setchatmenubutton
type SetChatMenuButton struct {
	ChatID     int64       `json:"chat_id,omitempty"`
	MenuButton *MenuButton `json:"menu_button,omitempty"`
}

// SetChatMenuButton changes the menu button for a given chat or the default menu button.
// Since: Bot API 6.0
// Returns true on success.
// See https://core.telegram.org/bots/api#setchatmenubutton
func (api *API) SetChatMenuButton(params SetChatMenuButton) (bool, error) {
	req := NewRequest[bool]("setChatMenuButton", params)
	return req.Do(api)
}

// SetChatMenuButtonWithContext is the context-aware variant of SetChatMenuButton.
// Since: Bot API 6.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setchatmenubutton
func (api *API) SetChatMenuButtonWithContext(ctx context.Context, params SetChatMenuButton) (bool, error) {
	req := NewRequest[bool]("setChatMenuButton", params)
	return req.DoWithContext(ctx, api)
}

// GetChatMenuButton holds parameters for the getChatMenuButton method.
// Since: Bot API 6.0
// See https://core.telegram.org/bots/api#getchatmenubutton
type GetChatMenuButton struct {
	ChatID int64 `json:"chat_id,omitempty"`
}

// GetChatMenuButton returns the current menu button for the given chat.
// Since: Bot API 6.0
// See https://core.telegram.org/bots/api#getchatmenubutton
func (api *API) GetChatMenuButton(params GetChatMenuButton) (MenuButton, error) {
	req := NewRequest[MenuButton]("getChatMenuButton", params)
	return req.Do(api)
}

// GetChatMenuButtonWithContext is the context-aware variant of GetChatMenuButton.
// Since: Bot API 6.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getchatmenubutton
func (api *API) GetChatMenuButtonWithContext(ctx context.Context, params GetChatMenuButton) (MenuButton, error) {
	req := NewRequest[MenuButton]("getChatMenuButton", params)
	return req.DoWithContext(ctx, api)
}

// SetMyDefaultAdministratorRights holds parameters for the setMyDefaultAdministratorRights method.
// Since: Bot API 6.0
// See https://core.telegram.org/bots/api#setmydefaultadministratorrights
type SetMyDefaultAdministratorRights struct {
	Rights      *ChatAdministratorRights `json:"rights"`
	ForChannels bool                     `json:"for_channels"`
}

// SetMyDefaultAdministratorRights changes the default administrator rights for the bot.
// Since: Bot API 6.0
// Returns true on success.
// See https://core.telegram.org/bots/api#setmydefaultadministratorrights
func (api *API) SetMyDefaultAdministratorRights(params SetMyDefaultAdministratorRights) (bool, error) {
	req := NewRequest[bool]("setMyDefaultAdministratorRights", params)
	return req.Do(api)
}

// SetMyDefaultAdministratorRightsWithContext is the context-aware variant of SetMyDefaultAdministratorRights.
// Since: Bot API 6.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setmydefaultadministratorrights
func (api *API) SetMyDefaultAdministratorRightsWithContext(ctx context.Context, params SetMyDefaultAdministratorRights) (bool, error) {
	req := NewRequest[bool]("setMyDefaultAdministratorRights", params)
	return req.DoWithContext(ctx, api)
}

// GetMyDefaultAdministratorRights holds parameters for the getMyDefaultAdministratorRights method.
// Since: Bot API 6.0
// See https://core.telegram.org/bots/api#getmydefaultadministratorrights
type GetMyDefaultAdministratorRights struct {
	ForChannels bool `json:"for_channels"`
}

// GetMyDefaultAdministratorRights returns the current default administrator rights for the bot.
// Since: Bot API 6.0
// See https://core.telegram.org/bots/api#getmydefaultadministratorrights
func (api *API) GetMyDefaultAdministratorRights(params GetMyDefaultAdministratorRights) (ChatAdministratorRights, error) {
	req := NewRequest[ChatAdministratorRights]("getMyDefaultAdministratorRights", params)
	return req.Do(api)
}

// GetMyDefaultAdministratorRightsWithContext is the context-aware variant of GetMyDefaultAdministratorRights.
// Since: Bot API 6.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getmydefaultadministratorrights
func (api *API) GetMyDefaultAdministratorRightsWithContext(ctx context.Context, params GetMyDefaultAdministratorRights) (ChatAdministratorRights, error) {
	req := NewRequest[ChatAdministratorRights]("getMyDefaultAdministratorRights", params)
	return req.DoWithContext(ctx, api)
}

// GetAvailableGifts returns the list of gifts that can be sent by the bot.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#getavailablegifts
func (api *API) GetAvailableGifts() (Gifts, error) {
	req := NewRequest[Gifts]("getAvailableGifts", NoParams)
	return req.Do(api)
}

// GetAvailableGiftsWithContext is the context-aware variant of GetAvailableGifts.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getavailablegifts
func (api *API) GetAvailableGiftsWithContext(ctx context.Context) (Gifts, error) {
	req := NewRequest[Gifts]("getAvailableGifts", NoParams)
	return req.DoWithContext(ctx, api)
}

// SendGift holds parameters for the sendGift method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#sendgift
type SendGift struct {
	UserID        int64           `json:"user_id,omitempty"`
	ChatID        int64           `json:"chat_id,omitempty"`
	GiftID        string          `json:"gift_id"`
	PayForUpgrade bool            `json:"pay_for_upgrade"`
	Text          string          `json:"text"`
	TextParseMode ParseMode       `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity `json:"text_entities,omitempty"`
}

// SendGift sends a gift to the given user or chat.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#sendgift
func (api *API) SendGift(params SendGift) (bool, error) {
	req := NewRequest[bool]("sendGift", params)
	return req.Do(api)
}

// SendGiftWithContext is the context-aware variant of SendGift.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendgift
func (api *API) SendGiftWithContext(ctx context.Context, params SendGift) (bool, error) {
	req := NewRequest[bool]("sendGift", params)
	return req.DoWithContext(ctx, api)
}

// GiftPremiumSubscription holds parameters for the giftPremiumSubscription method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#giftpremiumsubscription
type GiftPremiumSubscription struct {
	UserID        int64           `json:"user_id"`
	MonthCount    int             `json:"month_count"`
	StarCount     int             `json:"star_count"`
	Text          string          `json:"text,omitempty"`
	TextParseMode ParseMode       `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity `json:"text_entities,omitempty"`
}

// GiftPremiumSubscription gifts a Telegram Premium subscription to the user.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#giftpremiumsubscription
func (api *API) GiftPremiumSubscription(params GiftPremiumSubscription) (bool, error) {
	req := NewRequest[bool]("giftPremiumSubscription", params)
	return req.Do(api)
}

// GiftPremiumSubscriptionWithContext is the context-aware variant of GiftPremiumSubscription.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#giftpremiumsubscription
func (api *API) GiftPremiumSubscriptionWithContext(ctx context.Context, params GiftPremiumSubscription) (bool, error) {
	req := NewRequest[bool]("giftPremiumSubscription", params)
	return req.DoWithContext(ctx, api)
}

// GetManagedBotAccessSettings holds parameters for the getManagedBotAccessSettings method.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#getmanagedbotaccesssettings
type GetManagedBotAccessSettings struct {
	BotUserID int64 `json:"bot_user_id"`
}

// GetManagedBotAccessSettings returns the access settings of a managed bot.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#getmanagedbotaccesssettings
func (api *API) GetManagedBotAccessSettings(params GetManagedBotAccessSettings) (BotAccessSettings, error) {
	req := NewRequest[BotAccessSettings]("getManagedBotAccessSettings", params)
	return req.Do(api)
}

// GetManagedBotAccessSettingsWithContext is the context-aware variant of GetManagedBotAccessSettings.
// Since: Bot API 10.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getmanagedbotaccesssettings
func (api *API) GetManagedBotAccessSettingsWithContext(ctx context.Context, params GetManagedBotAccessSettings) (BotAccessSettings, error) {
	req := NewRequest[BotAccessSettings]("getManagedBotAccessSettings", params)
	return req.DoWithContext(ctx, api)
}

// SetManagedBotAccessSettings holds parameters for the setManagedBotAccessSettings method.
// Since: Bot API 10.0
// See https://core.telegram.org/bots/api#setmanagedbotaccesssettings
type SetManagedBotAccessSettings struct {
	BotUserID      int64             `json:"bot_user_id"`
	AccessSettings BotAccessSettings `json:"access_settings"`
}

// SetManagedBotAccessSettings changes the access settings of a managed bot.
// Since: Bot API 10.0
// Returns True on success.
// See https://core.telegram.org/bots/api#setmanagedbotaccesssettings
func (api *API) SetManagedBotAccessSettings(params SetManagedBotAccessSettings) (bool, error) {
	req := NewRequest[bool]("setManagedBotAccessSettings", params)
	return req.Do(api)
}

// SetManagedBotAccessSettingsWithContext is the context-aware variant of SetManagedBotAccessSettings.
// Since: Bot API 10.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setmanagedbotaccesssettings
func (api *API) SetManagedBotAccessSettingsWithContext(ctx context.Context, params SetManagedBotAccessSettings) (bool, error) {
	req := NewRequest[bool]("setManagedBotAccessSettings", params)
	return req.DoWithContext(ctx, api)
}
