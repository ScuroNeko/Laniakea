package tgapi

import "context"

// VerifyUser holds parameters for the verifyUser method.
// Since: Bot API 8.0
// See https://core.telegram.org/bots/api#verifyuser
type VerifyUser struct {
	UserID            int64  `json:"user_id"`
	CustomDescription string `json:"custom_description,omitempty"`
}

// VerifyUser verifies a user.
// Since: Bot API 8.0
// Returns true on success.
// See https://core.telegram.org/bots/api#verifyuser
func (api *API) VerifyUser(params VerifyUser) (bool, error) {
	req := NewRequest[bool]("verifyUser", params)
	return req.Do(api)
}

// VerifyUserWithContext is the context-aware variant of VerifyUser.
// Since: Bot API 8.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#verifyuser
func (api *API) VerifyUserWithContext(ctx context.Context, params VerifyUser) (bool, error) {
	req := NewRequest[bool]("verifyUser", params)
	return req.DoWithContext(ctx, api)
}

// VerifyChat holds parameters for the verifyChat method.
// Since: Bot API 8.0
// See https://core.telegram.org/bots/api#verifychat
type VerifyChat struct {
	ChatID            int64  `json:"chat_id"`
	CustomDescription string `json:"custom_description,omitempty"`
}

// VerifyChat verifies a chat.
// Since: Bot API 8.0
// Returns true on success.
// See https://core.telegram.org/bots/api#verifychat
func (api *API) VerifyChat(params VerifyChat) (bool, error) {
	req := NewRequest[bool]("verifyChat", params)
	return req.Do(api)
}

// VerifyChatWithContext is the context-aware variant of VerifyChat.
// Since: Bot API 8.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#verifychat
func (api *API) VerifyChatWithContext(ctx context.Context, params VerifyChat) (bool, error) {
	req := NewRequest[bool]("verifyChat", params)
	return req.DoWithContext(ctx, api)
}

// RemoveUserVerification holds parameters for the removeUserVerification method.
// Since: Bot API 8.0
// See https://core.telegram.org/bots/api#removeuserverification
type RemoveUserVerification struct {
	UserID int64 `json:"user_id"`
}

// RemoveUserVerification removes a user's verification.
// Since: Bot API 8.0
// Returns true on success.
// See https://core.telegram.org/bots/api#removeuserverification
func (api *API) RemoveUserVerification(params RemoveUserVerification) (bool, error) {
	req := NewRequest[bool]("removeUserVerification", params)
	return req.Do(api)
}

// RemoveUserVerificationWithContext is the context-aware variant of RemoveUserVerification.
// Since: Bot API 8.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#removeuserverification
func (api *API) RemoveUserVerificationWithContext(ctx context.Context, params RemoveUserVerification) (bool, error) {
	req := NewRequest[bool]("removeUserVerification", params)
	return req.DoWithContext(ctx, api)
}

// RemoveChatVerification holds parameters for the removeChatVerification method.
// Since: Bot API 8.0
// See https://core.telegram.org/bots/api#removechatverification
type RemoveChatVerification struct {
	ChatID int64 `json:"chat_id"`
}

// RemoveChatVerification removes a chat's verification.
// Since: Bot API 8.0
// Returns true on success.
// See https://core.telegram.org/bots/api#removechatverification
func (api *API) RemoveChatVerification(params RemoveChatVerification) (bool, error) {
	req := NewRequest[bool]("removeChatVerification", params)
	return req.Do(api)
}

// RemoveChatVerificationWithContext is the context-aware variant of RemoveChatVerification.
// Since: Bot API 8.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#removechatverification
func (api *API) RemoveChatVerificationWithContext(ctx context.Context, params RemoveChatVerification) (bool, error) {
	req := NewRequest[bool]("removeChatVerification", params)
	return req.DoWithContext(ctx, api)
}

// ReadBusinessMessage holds parameters for the readBusinessMessage method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#readbusinessmessage
type ReadBusinessMessage struct {
	BusinessConnectionID string `json:"business_connection_id"`
	ChatID               int64  `json:"chat_id"`
	MessageID            int    `json:"message_id"`
}

// ReadBusinessMessage marks a business message as read.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#readbusinessmessage
func (api *API) ReadBusinessMessage(params ReadBusinessMessage) (bool, error) {
	req := NewRequest[bool]("readBusinessMessage", params)
	return req.Do(api)
}

// ReadBusinessMessageWithContext is the context-aware variant of ReadBusinessMessage.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#readbusinessmessage
func (api *API) ReadBusinessMessageWithContext(ctx context.Context, params ReadBusinessMessage) (bool, error) {
	req := NewRequest[bool]("readBusinessMessage", params)
	return req.DoWithContext(ctx, api)
}

// GetBusinessConnection holds parameters for the getBusinessConnection method.
// Since: Bot API 7.2
// See https://core.telegram.org/bots/api#getbusinessconnection
type GetBusinessConnection struct {
	BusinessConnectionID string `json:"business_connection_id"`
}

// GetBusinessConnection returns information about a business connection.
// Since: Bot API 7.2
// See https://core.telegram.org/bots/api#getbusinessconnection
func (api *API) GetBusinessConnection(params GetBusinessConnection) (BusinessConnection, error) {
	req := NewRequest[BusinessConnection]("getBusinessConnection", params)
	return req.Do(api)
}

// GetBusinessConnectionWithContext is the context-aware variant of GetBusinessConnection.
// Since: Bot API 7.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getbusinessconnection
func (api *API) GetBusinessConnectionWithContext(ctx context.Context, params GetBusinessConnection) (BusinessConnection, error) {
	req := NewRequest[BusinessConnection]("getBusinessConnection", params)
	return req.DoWithContext(ctx, api)
}

// DeleteBusinessMessages holds parameters for the deleteBusinessMessages method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#deletebusinessmessages
type DeleteBusinessMessages struct {
	BusinessConnectionID string `json:"business_connection_id"`
	MessageIDs           []int  `json:"message_ids"`
}

// DeleteBusinessMessages deletes business messages.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#deletebusinessmessages
func (api *API) DeleteBusinessMessages(params DeleteBusinessMessages) (bool, error) {
	req := NewRequest[bool]("deleteBusinessMessages", params)
	return req.Do(api)
}

// DeleteBusinessMessagesWithContext is the context-aware variant of DeleteBusinessMessages.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletebusinessmessages
func (api *API) DeleteBusinessMessagesWithContext(ctx context.Context, params DeleteBusinessMessages) (bool, error) {
	req := NewRequest[bool]("deleteBusinessMessages", params)
	return req.DoWithContext(ctx, api)
}

// SetBusinessAccountName holds parameters for the setBusinessAccountName method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#setbusinessaccountname
type SetBusinessAccountName struct {
	BusinessConnectionID string `json:"business_connection_id"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name,omitempty"`
}

// SetBusinessAccountName sets the first and last name of a business account.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#setbusinessaccountname
func (api *API) SetBusinessAccountName(params SetBusinessAccountName) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountName", params)
	return req.Do(api)
}

// SetBusinessAccountNameWithContext is the context-aware variant of SetBusinessAccountName.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setbusinessaccountname
func (api *API) SetBusinessAccountNameWithContext(ctx context.Context, params SetBusinessAccountName) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountName", params)
	return req.DoWithContext(ctx, api)
}

// SetBusinessAccountUsername holds parameters for the setBusinessAccountUsername method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#setbusinessaccountusername
type SetBusinessAccountUsername struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Username             string `json:"username,omitempty"`
}

// SetBusinessAccountUsername sets the username of a business account.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#setbusinessaccountusername
func (api *API) SetBusinessAccountUsername(params SetBusinessAccountUsername) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountUsername", params)
	return req.Do(api)
}

// SetBusinessAccountUsernameWithContext is the context-aware variant of SetBusinessAccountUsername.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setbusinessaccountusername
func (api *API) SetBusinessAccountUsernameWithContext(ctx context.Context, params SetBusinessAccountUsername) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountUsername", params)
	return req.DoWithContext(ctx, api)
}

// SetBusinessAccountBio holds parameters for the setBusinessAccountBio method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#setbusinessaccountbio
type SetBusinessAccountBio struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Bio                  string `json:"bio,omitempty"`
}

// SetBusinessAccountBio sets the bio of a business account.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#setbusinessaccountbio
func (api *API) SetBusinessAccountBio(params SetBusinessAccountBio) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountBio", params)
	return req.Do(api)
}

// SetBusinessAccountBioWithContext is the context-aware variant of SetBusinessAccountBio.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setbusinessaccountbio
func (api *API) SetBusinessAccountBioWithContext(ctx context.Context, params SetBusinessAccountBio) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountBio", params)
	return req.DoWithContext(ctx, api)
}

// SetBusinessAccountProfilePhoto holds parameters for the setBusinessAccountProfilePhoto method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#setbusinessaccountprofilephoto
type SetBusinessAccountProfilePhoto struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	Photo                InputProfilePhoto `json:"photo,omitempty"`
	IsPublic             bool              `json:"is_public,omitempty"`
}

// SetBusinessAccountProfilePhoto sets the profile photo of a business account.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#setbusinessaccountprofilephoto
func (api *API) SetBusinessAccountProfilePhoto(params SetBusinessAccountProfilePhoto) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountProfilePhoto", params)
	return req.Do(api)
}

// SetBusinessAccountProfilePhotoWithContext is the context-aware variant of SetBusinessAccountProfilePhoto.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setbusinessaccountprofilephoto
func (api *API) SetBusinessAccountProfilePhotoWithContext(ctx context.Context, params SetBusinessAccountProfilePhoto) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountProfilePhoto", params)
	return req.DoWithContext(ctx, api)
}

// RemoveBusinessAccountProfilePhoto holds parameters for the removeBusinessAccountProfilePhoto method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#removebusinessaccountprofilephoto
type RemoveBusinessAccountProfilePhoto struct {
	BusinessConnectionID string `json:"business_connection_id"`
	IsPublic             bool   `json:"is_public,omitempty"`
}

// RemoveBusinessAccountProfilePhoto removes the profile photo of a business account.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#removebusinessaccountprofilephoto
func (api *API) RemoveBusinessAccountProfilePhoto(params RemoveBusinessAccountProfilePhoto) (bool, error) {
	req := NewRequest[bool]("removeBusinessAccountProfilePhoto", params)
	return req.Do(api)
}

// RemoveBusinessAccountProfilePhotoWithContext is the context-aware variant of RemoveBusinessAccountProfilePhoto.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#removebusinessaccountprofilephoto
func (api *API) RemoveBusinessAccountProfilePhotoWithContext(ctx context.Context, params RemoveBusinessAccountProfilePhoto) (bool, error) {
	req := NewRequest[bool]("removeBusinessAccountProfilePhoto", params)
	return req.DoWithContext(ctx, api)
}

// SetBusinessAccountGiftSettings holds parameters for the setBusinessAccountGiftSettings method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#setbusinessaccountgiftsettings
type SetBusinessAccountGiftSettings struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	ShowGiftButton       bool              `json:"show_gift_button"`
	AcceptedGiftTypes    AcceptedGiftTypes `json:"accepted_gift_types"`
}

// SetBusinessAccountGiftSettings sets gift settings for a business account.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#setbusinessaccountgiftsettings
func (api *API) SetBusinessAccountGiftSettings(params SetBusinessAccountGiftSettings) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountGiftSettings", params)
	return req.Do(api)
}

// SetBusinessAccountGiftSettingsWithContext is the context-aware variant of SetBusinessAccountGiftSettings.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setbusinessaccountgiftsettings
func (api *API) SetBusinessAccountGiftSettingsWithContext(ctx context.Context, params SetBusinessAccountGiftSettings) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountGiftSettings", params)
	return req.DoWithContext(ctx, api)
}

// GetBusinessAccountStarBalance holds parameters for the getBusinessAccountStarBalance method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#getbusinessaccountstarbalance
type GetBusinessAccountStarBalance struct {
	BusinessConnectionID string `json:"business_connection_id"`
}

// GetBusinessAccountStarBalance returns the star balance of a business account.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#getbusinessaccountstarbalance
func (api *API) GetBusinessAccountStarBalance(params GetBusinessAccountStarBalance) (StarAmount, error) {
	req := NewRequest[StarAmount]("getBusinessAccountStarBalance", params)
	return req.Do(api)
}

// GetBusinessAccountStarBalanceWithContext is the context-aware variant of GetBusinessAccountStarBalance.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getbusinessaccountstarbalance
func (api *API) GetBusinessAccountStarBalanceWithContext(ctx context.Context, params GetBusinessAccountStarBalance) (StarAmount, error) {
	req := NewRequest[StarAmount]("getBusinessAccountStarBalance", params)
	return req.DoWithContext(ctx, api)
}

// TransferBusinessAccountStars holds parameters for the transferBusinessAccountStars method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#transferbusinessaccountstars
type TransferBusinessAccountStars struct {
	BusinessConnectionID string `json:"business_connection_id"`
	StarCount            int    `json:"star_count"`
}

// TransferBusinessAccountStars transfers stars from a business account.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#transferbusinessaccountstars
func (api *API) TransferBusinessAccountStars(params TransferBusinessAccountStars) (bool, error) {
	req := NewRequest[bool]("transferBusinessAccountStars", params)
	return req.Do(api)
}

// TransferBusinessAccountStarsWithContext is the context-aware variant of TransferBusinessAccountStars.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#transferbusinessaccountstars
func (api *API) TransferBusinessAccountStarsWithContext(ctx context.Context, params TransferBusinessAccountStars) (bool, error) {
	req := NewRequest[bool]("transferBusinessAccountStars", params)
	return req.DoWithContext(ctx, api)
}

// GetBusinessAccountGifts holds parameters for the getBusinessAccountGifts method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#getbusinessaccountgifts
type GetBusinessAccountGifts struct {
	BusinessConnectionID        string `json:"business_connection_id"`
	ExcludeUnsaved              bool   `json:"exclude_unsaved,omitempty"`
	ExcludeSaved                bool   `json:"exclude_saved,omitempty"`
	ExcludeUnlimited            bool   `json:"exclude_unlimited,omitempty"`
	ExcludeLimitedUpgradable    bool   `json:"exclude_limited_upgradable,omitempty"`
	ExcludeLimitedNonUpgradable bool   `json:"exclude_limited_non_upgradable,omitempty"`
	ExcludeUnique               bool   `json:"exclude_unique,omitempty"`
	ExcludeFromBlockchain       bool   `json:"exclude_from_blockchain,omitempty"`
	SortByPrice                 bool   `json:"sort_by_price,omitempty"`
	Offset                      string `json:"offset,omitempty"`
	Limit                       int    `json:"limit,omitempty"`
}

// GetBusinessAccountGifts returns gifts owned by a business account.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#getbusinessaccountgifts
func (api *API) GetBusinessAccountGifts(params GetBusinessAccountGifts) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getBusinessAccountGifts", params)
	return req.Do(api)
}

// GetBusinessAccountGiftsWithContext is the context-aware variant of GetBusinessAccountGifts.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getbusinessaccountgifts
func (api *API) GetBusinessAccountGiftsWithContext(ctx context.Context, params GetBusinessAccountGifts) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getBusinessAccountGifts", params)
	return req.DoWithContext(ctx, api)
}

// ConvertGiftToStars holds parameters for the convertGiftToStars method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#convertgifttostars
type ConvertGiftToStars struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
}

// ConvertGiftToStars converts a gift to Telegram Stars.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#convertgifttostars
func (api *API) ConvertGiftToStars(params ConvertGiftToStars) (bool, error) {
	req := NewRequest[bool]("convertGiftToStars", params)
	return req.Do(api)
}

// ConvertGiftToStarsWithContext is the context-aware variant of ConvertGiftToStars.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#convertgifttostars
func (api *API) ConvertGiftToStarsWithContext(ctx context.Context, params ConvertGiftToStars) (bool, error) {
	req := NewRequest[bool]("convertGiftToStars", params)
	return req.DoWithContext(ctx, api)
}

// UpgradeGift holds parameters for the upgradeGift method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#upgradegift
type UpgradeGift struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
	KeepOriginalDetails  bool   `json:"keep_original_details,omitempty"`
	StarCount            int    `json:"star_count,omitempty"`
}

// UpgradeGift upgrades a gift.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#upgradegift
func (api *API) UpgradeGift(params UpgradeGift) (bool, error) {
	req := NewRequest[bool]("upgradeGift", params)
	return req.Do(api)
}

// UpgradeGiftWithContext is the context-aware variant of UpgradeGift.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#upgradegift
func (api *API) UpgradeGiftWithContext(ctx context.Context, params UpgradeGift) (bool, error) {
	req := NewRequest[bool]("upgradeGift", params)
	return req.DoWithContext(ctx, api)
}

// TransferGift holds parameters for the transferGift method.
// Since: Bot API 9.0
// See https://core.telegram.org/bots/api#transfergift
type TransferGift struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
	NewOwnerChatID       int64  `json:"new_owner_chat_id"`
	StarCount            int    `json:"star_count,omitempty"`
}

// TransferGift transfers a gift to another chat.
// Since: Bot API 9.0
// Returns true on success.
// See https://core.telegram.org/bots/api#transfergift
func (api *API) TransferGift(params TransferGift) (bool, error) {
	req := NewRequest[bool]("transferGift", params)
	return req.Do(api)
}

// TransferGiftWithContext is the context-aware variant of TransferGift.
// Since: Bot API 9.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#transfergift
func (api *API) TransferGiftWithContext(ctx context.Context, params TransferGift) (bool, error) {
	req := NewRequest[bool]("transferGift", params)
	return req.DoWithContext(ctx, api)
}

// PostStory holds parameters for the postStory method.
// Since: Bot API 7.2
// See https://core.telegram.org/bots/api#poststory
type PostStory struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	Content              InputStoryContent `json:"content"`
	ActivePeriod         int               `json:"active_period"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Areas           []StoryArea     `json:"areas"`

	PostToChatPage bool `json:"post_to_chat_page,omitempty"`
	ProtectContent bool `json:"protect_content,omitempty"`
}

// PostStory posts a story with a photo.
// Since: Bot API 7.2
// See https://core.telegram.org/bots/api#poststory
func (api *API) PostStory(params PostStory) (Story, error) {
	req := NewRequest[Story]("postStory", params)
	return req.Do(api)
}

// PostStoryWithContext is the context-aware variant of PostStory.
// Since: Bot API 7.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#poststory
func (api *API) PostStoryWithContext(ctx context.Context, params PostStory) (Story, error) {
	req := NewRequest[Story]("postStory", params)
	return req.DoWithContext(ctx, api)
}

// RepostStory holds parameters for the repostStory method.
// Since: Bot API 7.2
// See https://core.telegram.org/bots/api#repoststory
type RepostStory struct {
	BusinessConnectionID string `json:"business_connection_id"`
	FromChatID           int64  `json:"from_chat_id"`
	FromStoryID          int    `json:"from_story_id"`
	ActivePeriod         int    `json:"active_period"`
	PostToChatPage       bool   `json:"post_to_chat_page,omitempty"`
	ProtectContent       bool   `json:"protect_content,omitempty"`
}

// RepostStory reposts a story from another chat.
// Since: Bot API 7.2
// Returns the reposted story.
// See https://core.telegram.org/bots/api#repoststory
func (api *API) RepostStory(params RepostStory) (Story, error) {
	req := NewRequest[Story]("repostStory", params)
	return req.Do(api)
}

// RepostStoryWithContext is the context-aware variant of RepostStory.
// Since: Bot API 7.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#repoststory
func (api *API) RepostStoryWithContext(ctx context.Context, params RepostStory) (Story, error) {
	req := NewRequest[Story]("repostStory", params)
	return req.DoWithContext(ctx, api)
}

// EditStory holds parameters for the editStory method.
// Since: Bot API 7.2
// See https://core.telegram.org/bots/api#editstory
type EditStory struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	StoryID              int               `json:"story_id"`
	Content              InputStoryContent `json:"content"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Areas           []StoryArea     `json:"areas,omitempty"`
}

// EditStory edits an existing story.
// Since: Bot API 7.2
// Returns the updated story.
// See https://core.telegram.org/bots/api#editstory
func (api *API) EditStory(params EditStory) (Story, error) {
	req := NewRequest[Story]("editStory", params)
	return req.Do(api)
}

// EditStoryWithContext is the context-aware variant of EditStory.
// Since: Bot API 7.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editstory
func (api *API) EditStoryWithContext(ctx context.Context, params EditStory) (Story, error) {
	req := NewRequest[Story]("editStory", params)
	return req.DoWithContext(ctx, api)
}

// DeleteStory holds parameters for the deleteStory method.
// Since: Bot API 7.2
// See https://core.telegram.org/bots/api#deletestory
type DeleteStory struct {
	BusinessConnectionID string `json:"business_connection_id"`
	StoryID              int    `json:"story_id"`
}

// DeleteStory deletes a story.
// Since: Bot API 7.2
// Returns true on success.
// See https://core.telegram.org/bots/api#deletestory
func (api *API) DeleteStory(params DeleteStory) (bool, error) {
	req := NewRequest[bool]("deleteStory", params)
	return req.Do(api)
}

// DeleteStoryWithContext is the context-aware variant of DeleteStory.
// Since: Bot API 7.2
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletestory
func (api *API) DeleteStoryWithContext(ctx context.Context, params DeleteStory) (bool, error) {
	req := NewRequest[bool]("deleteStory", params)
	return req.DoWithContext(ctx, api)
}
