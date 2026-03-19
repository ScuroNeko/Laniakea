package tgapi

import "context"

// VerifyUserP holds parameters for the verifyUser method.
// See https://core.telegram.org/bots/api#verifyuser
type VerifyUserP struct {
	UserID            int64  `json:"user_id"`
	CustomDescription string `json:"custom_description,omitempty"`
}

// VerifyUser verifies a user.
// Returns true on success.
// See https://core.telegram.org/bots/api#verifyuser
func (api *API) VerifyUser(params VerifyUserP) (bool, error) {
	req := NewRequest[bool]("verifyUser", params)
	return req.Do(api)
}

// VerifyUserWithContext is the context-aware variant of VerifyUser.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#verifyuser
func (api *API) VerifyUserWithContext(ctx context.Context, params VerifyUserP) (bool, error) {
	req := NewRequest[bool]("verifyUser", params)
	return req.DoWithContext(ctx, api)
}

// VerifyChatP holds parameters for the verifyChat method.
// See https://core.telegram.org/bots/api#verifychat
type VerifyChatP struct {
	ChatID            int64  `json:"chat_id"`
	CustomDescription string `json:"custom_description,omitempty"`
}

// VerifyChat verifies a chat.
// Returns true on success.
// See https://core.telegram.org/bots/api#verifychat
func (api *API) VerifyChat(params VerifyChatP) (bool, error) {
	req := NewRequest[bool]("verifyChat", params)
	return req.Do(api)
}

// VerifyChatWithContext is the context-aware variant of VerifyChat.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#verifychat
func (api *API) VerifyChatWithContext(ctx context.Context, params VerifyChatP) (bool, error) {
	req := NewRequest[bool]("verifyChat", params)
	return req.DoWithContext(ctx, api)
}

// RemoveUserVerificationP holds parameters for the removeUserVerification method.
// See https://core.telegram.org/bots/api#removeuserverification
type RemoveUserVerificationP struct {
	UserID int64 `json:"user_id"`
}

// RemoveUserVerification removes a user's verification.
// Returns true on success.
// See https://core.telegram.org/bots/api#removeuserverification
func (api *API) RemoveUserVerification(params RemoveUserVerificationP) (bool, error) {
	req := NewRequest[bool]("removeUserVerification", params)
	return req.Do(api)
}

// RemoveUserVerificationWithContext is the context-aware variant of RemoveUserVerification.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#removeuserverification
func (api *API) RemoveUserVerificationWithContext(ctx context.Context, params RemoveUserVerificationP) (bool, error) {
	req := NewRequest[bool]("removeUserVerification", params)
	return req.DoWithContext(ctx, api)
}

// RemoveChatVerificationP holds parameters for the removeChatVerification method.
// See https://core.telegram.org/bots/api#removechatverification
type RemoveChatVerificationP struct {
	ChatID int64 `json:"chat_id"`
}

// RemoveChatVerification removes a chat's verification.
// Returns true on success.
// See https://core.telegram.org/bots/api#removechatverification
func (api *API) RemoveChatVerification(params RemoveChatVerificationP) (bool, error) {
	req := NewRequest[bool]("removeChatVerification", params)
	return req.Do(api)
}

// RemoveChatVerificationWithContext is the context-aware variant of RemoveChatVerification.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#removechatverification
func (api *API) RemoveChatVerificationWithContext(ctx context.Context, params RemoveChatVerificationP) (bool, error) {
	req := NewRequest[bool]("removeChatVerification", params)
	return req.DoWithContext(ctx, api)
}

// ReadBusinessMessageP holds parameters for the readBusinessMessage method.
// See https://core.telegram.org/bots/api#readbusinessmessage
type ReadBusinessMessageP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	ChatID               int64  `json:"chat_id"`
	MessageID            int    `json:"message_id"`
}

// ReadBusinessMessage marks a business message as read.
// Returns true on success.
// See https://core.telegram.org/bots/api#readbusinessmessage
func (api *API) ReadBusinessMessage(params ReadBusinessMessageP) (bool, error) {
	req := NewRequest[bool]("readBusinessMessage", params)
	return req.Do(api)
}

// ReadBusinessMessageWithContext is the context-aware variant of ReadBusinessMessage.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#readbusinessmessage
func (api *API) ReadBusinessMessageWithContext(ctx context.Context, params ReadBusinessMessageP) (bool, error) {
	req := NewRequest[bool]("readBusinessMessage", params)
	return req.DoWithContext(ctx, api)
}

// GetBusinessConnectionP holds parameters for the getBusinessConnection method.
// See https://core.telegram.org/bots/api#getbusinessconnection
type GetBusinessConnectionP struct {
	BusinessConnectionID string `json:"business_connection_id"`
}

// GetBusinessConnection returns information about a business connection.
// See https://core.telegram.org/bots/api#getbusinessconnection
func (api *API) GetBusinessConnection(params GetBusinessConnectionP) (BusinessConnection, error) {
	req := NewRequest[BusinessConnection]("getBusinessConnection", params)
	return req.Do(api)
}

// GetBusinessConnectionWithContext is the context-aware variant of GetBusinessConnection.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getbusinessconnection
func (api *API) GetBusinessConnectionWithContext(ctx context.Context, params GetBusinessConnectionP) (BusinessConnection, error) {
	req := NewRequest[BusinessConnection]("getBusinessConnection", params)
	return req.DoWithContext(ctx, api)
}

// DeleteBusinessMessagesP holds parameters for the deleteBusinessMessages method.
// See https://core.telegram.org/bots/api#deletebusinessmessages
type DeleteBusinessMessagesP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	MessageIDs           []int  `json:"message_ids"`
}

// DeleteBusinessMessages deletes business messages.
// Returns true on success.
// See https://core.telegram.org/bots/api#deletebusinessmessages
func (api *API) DeleteBusinessMessages(params DeleteBusinessMessagesP) (bool, error) {
	req := NewRequest[bool]("deleteBusinessMessages", params)
	return req.Do(api)
}

// DeleteBusinessMessagesWithContext is the context-aware variant of DeleteBusinessMessages.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletebusinessmessages
func (api *API) DeleteBusinessMessagesWithContext(ctx context.Context, params DeleteBusinessMessagesP) (bool, error) {
	req := NewRequest[bool]("deleteBusinessMessages", params)
	return req.DoWithContext(ctx, api)
}

// SetBusinessAccountNameP holds parameters for the setBusinessAccountName method.
// See https://core.telegram.org/bots/api#setbusinessaccountname
type SetBusinessAccountNameP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name,omitempty"`
}

// SetBusinessAccountName sets the first and last name of a business account.
// Returns true on success.
// See https://core.telegram.org/bots/api#setbusinessaccountname
func (api *API) SetBusinessAccountName(params SetBusinessAccountNameP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountName", params)
	return req.Do(api)
}

// SetBusinessAccountNameWithContext is the context-aware variant of SetBusinessAccountName.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setbusinessaccountname
func (api *API) SetBusinessAccountNameWithContext(ctx context.Context, params SetBusinessAccountNameP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountName", params)
	return req.DoWithContext(ctx, api)
}

// SetBusinessAccountUsernameP holds parameters for the setBusinessAccountUsername method.
// See https://core.telegram.org/bots/api#setbusinessaccountusername
type SetBusinessAccountUsernameP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Username             string `json:"username,omitempty"`
}

// SetBusinessAccountUsername sets the username of a business account.
// Returns true on success.
// See https://core.telegram.org/bots/api#setbusinessaccountusername
func (api *API) SetBusinessAccountUsername(params SetBusinessAccountUsernameP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountUsername", params)
	return req.Do(api)
}

// SetBusinessAccountUsernameWithContext is the context-aware variant of SetBusinessAccountUsername.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setbusinessaccountusername
func (api *API) SetBusinessAccountUsernameWithContext(ctx context.Context, params SetBusinessAccountUsernameP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountUsername", params)
	return req.DoWithContext(ctx, api)
}

// SetBusinessAccountBioP holds parameters for the setBusinessAccountBio method.
// See https://core.telegram.org/bots/api#setbusinessaccountbio
type SetBusinessAccountBioP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Bio                  string `json:"bio,omitempty"`
}

// SetBusinessAccountBio sets the bio of a business account.
// Returns true on success.
// See https://core.telegram.org/bots/api#setbusinessaccountbio
func (api *API) SetBusinessAccountBio(params SetBusinessAccountBioP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountBio", params)
	return req.Do(api)
}

// SetBusinessAccountBioWithContext is the context-aware variant of SetBusinessAccountBio.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setbusinessaccountbio
func (api *API) SetBusinessAccountBioWithContext(ctx context.Context, params SetBusinessAccountBioP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountBio", params)
	return req.DoWithContext(ctx, api)
}

// SetBusinessAccountProfilePhoto holds parameters for the setBusinessAccountProfilePhoto method.
// See https://core.telegram.org/bots/api#setbusinessaccountprofilephoto
type SetBusinessAccountProfilePhoto struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	Photo                InputProfilePhoto `json:"photo,omitempty"`
	IsPublic             bool              `json:"is_public,omitempty"`
}

// SetBusinessAccountProfilePhoto sets the profile photo of a business account.
// Returns true on success.
// See https://core.telegram.org/bots/api#setbusinessaccountprofilephoto
func (api *API) SetBusinessAccountProfilePhoto(params SetBusinessAccountProfilePhoto) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountProfilePhoto", params)
	return req.Do(api)
}

// SetBusinessAccountProfilePhotoWithContext is the context-aware variant of SetBusinessAccountProfilePhoto.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setbusinessaccountprofilephoto
func (api *API) SetBusinessAccountProfilePhotoWithContext(ctx context.Context, params SetBusinessAccountProfilePhoto) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountProfilePhoto", params)
	return req.DoWithContext(ctx, api)
}

// RemoveBusinessAccountProfilePhotoP holds parameters for the removeBusinessAccountProfilePhoto method.
// See https://core.telegram.org/bots/api#removebusinessaccountprofilephoto
type RemoveBusinessAccountProfilePhotoP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	IsPublic             bool   `json:"is_public,omitempty"`
}

// RemoveBusinessAccountProfilePhoto removes the profile photo of a business account.
// Returns true on success.
// See https://core.telegram.org/bots/api#removebusinessaccountprofilephoto
func (api *API) RemoveBusinessAccountProfilePhoto(params RemoveBusinessAccountProfilePhotoP) (bool, error) {
	req := NewRequest[bool]("removeBusinessAccountProfilePhoto", params)
	return req.Do(api)
}

// RemoveBusinessAccountProfilePhotoWithContext is the context-aware variant of RemoveBusinessAccountProfilePhoto.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#removebusinessaccountprofilephoto
func (api *API) RemoveBusinessAccountProfilePhotoWithContext(ctx context.Context, params RemoveBusinessAccountProfilePhotoP) (bool, error) {
	req := NewRequest[bool]("removeBusinessAccountProfilePhoto", params)
	return req.DoWithContext(ctx, api)
}

// SetBusinessAccountGiftSettingsP holds parameters for the setBusinessAccountGiftSettings method.
// See https://core.telegram.org/bots/api#setbusinessaccountgiftsettings
type SetBusinessAccountGiftSettingsP struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	ShowGiftButton       bool              `json:"show_gift_button"`
	AcceptedGiftTypes    AcceptedGiftTypes `json:"accepted_gift_types"`
}

// SetBusinessAccountGiftSettings sets gift settings for a business account.
// Returns true on success.
// See https://core.telegram.org/bots/api#setbusinessaccountgiftsettings
func (api *API) SetBusinessAccountGiftSettings(params SetBusinessAccountGiftSettingsP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountGiftSettings", params)
	return req.Do(api)
}

// SetBusinessAccountGiftSettingsWithContext is the context-aware variant of SetBusinessAccountGiftSettings.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setbusinessaccountgiftsettings
func (api *API) SetBusinessAccountGiftSettingsWithContext(ctx context.Context, params SetBusinessAccountGiftSettingsP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountGiftSettings", params)
	return req.DoWithContext(ctx, api)
}

// GetBusinessAccountStarBalanceP holds parameters for the getBusinessAccountStarBalance method.
// See https://core.telegram.org/bots/api#getbusinessaccountstarbalance
type GetBusinessAccountStarBalanceP struct {
	BusinessConnectionID string `json:"business_connection_id"`
}

// GetBusinessAccountStarBalance returns the star balance of a business account.
// See https://core.telegram.org/bots/api#getbusinessaccountstarbalance
func (api *API) GetBusinessAccountStarBalance(params GetBusinessAccountStarBalanceP) (StarAmount, error) {
	req := NewRequest[StarAmount]("getBusinessAccountStarBalance", params)
	return req.Do(api)
}

// GetBusinessAccountStarBalanceWithContext is the context-aware variant of GetBusinessAccountStarBalance.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getbusinessaccountstarbalance
func (api *API) GetBusinessAccountStarBalanceWithContext(ctx context.Context, params GetBusinessAccountStarBalanceP) (StarAmount, error) {
	req := NewRequest[StarAmount]("getBusinessAccountStarBalance", params)
	return req.DoWithContext(ctx, api)
}

// TransferBusinessAccountStarsP holds parameters for the transferBusinessAccountStars method.
// See https://core.telegram.org/bots/api#transferbusinessaccountstars
type TransferBusinessAccountStarsP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	StarCount            int    `json:"star_count"`
}

// TransferBusinessAccountStars transfers stars from a business account.
// Returns true on success.
// See https://core.telegram.org/bots/api#transferbusinessaccountstars
func (api *API) TransferBusinessAccountStars(params TransferBusinessAccountStarsP) (bool, error) {
	req := NewRequest[bool]("transferBusinessAccountStars", params)
	return req.Do(api)
}

// TransferBusinessAccountStarsWithContext is the context-aware variant of TransferBusinessAccountStars.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#transferbusinessaccountstars
func (api *API) TransferBusinessAccountStarsWithContext(ctx context.Context, params TransferBusinessAccountStarsP) (bool, error) {
	req := NewRequest[bool]("transferBusinessAccountStars", params)
	return req.DoWithContext(ctx, api)
}

// GetBusinessAccountGiftsP holds parameters for the getBusinessAccountGifts method.
// See https://core.telegram.org/bots/api#getbusinessaccountgifts
type GetBusinessAccountGiftsP struct {
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
// See https://core.telegram.org/bots/api#getbusinessaccountgifts
func (api *API) GetBusinessAccountGifts(params GetBusinessAccountGiftsP) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getBusinessAccountGifts", params)
	return req.Do(api)
}

// GetBusinessAccountGiftsWithContext is the context-aware variant of GetBusinessAccountGifts.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getbusinessaccountgifts
func (api *API) GetBusinessAccountGiftsWithContext(ctx context.Context, params GetBusinessAccountGiftsP) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getBusinessAccountGifts", params)
	return req.DoWithContext(ctx, api)
}

// ConvertGiftToStarsP holds parameters for the convertGiftToStars method.
// See https://core.telegram.org/bots/api#convertgifttostars
type ConvertGiftToStarsP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
}

// ConvertGiftToStars converts a gift to Telegram Stars.
// Returns true on success.
// See https://core.telegram.org/bots/api#convertgifttostars
func (api *API) ConvertGiftToStars(params ConvertGiftToStarsP) (bool, error) {
	req := NewRequest[bool]("convertGiftToStars", params)
	return req.Do(api)
}

// ConvertGiftToStarsWithContext is the context-aware variant of ConvertGiftToStars.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#convertgifttostars
func (api *API) ConvertGiftToStarsWithContext(ctx context.Context, params ConvertGiftToStarsP) (bool, error) {
	req := NewRequest[bool]("convertGiftToStars", params)
	return req.DoWithContext(ctx, api)
}

// UpgradeGiftP holds parameters for the upgradeGift method.
// See https://core.telegram.org/bots/api#upgradegift
type UpgradeGiftP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
	KeepOriginalDetails  bool   `json:"keep_original_details,omitempty"`
	StarCount            int    `json:"star_count,omitempty"`
}

// UpgradeGift upgrades a gift.
// Returns true on success.
// See https://core.telegram.org/bots/api#upgradegift
func (api *API) UpgradeGift(params UpgradeGiftP) (bool, error) {
	req := NewRequest[bool]("upgradeGift", params)
	return req.Do(api)
}

// UpgradeGiftWithContext is the context-aware variant of UpgradeGift.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#upgradegift
func (api *API) UpgradeGiftWithContext(ctx context.Context, params UpgradeGiftP) (bool, error) {
	req := NewRequest[bool]("upgradeGift", params)
	return req.DoWithContext(ctx, api)
}

// TransferGiftP holds parameters for the transferGift method.
// See https://core.telegram.org/bots/api#transfergift
type TransferGiftP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
	NewOwnerChatID       int64  `json:"new_owner_chat_id"`
	StarCount            int    `json:"star_count,omitempty"`
}

// TransferGift transfers a gift to another chat.
// Returns true on success.
// See https://core.telegram.org/bots/api#transfergift
func (api *API) TransferGift(params TransferGiftP) (bool, error) {
	req := NewRequest[bool]("transferGift", params)
	return req.Do(api)
}

// TransferGiftWithContext is the context-aware variant of TransferGift.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#transfergift
func (api *API) TransferGiftWithContext(ctx context.Context, params TransferGiftP) (bool, error) {
	req := NewRequest[bool]("transferGift", params)
	return req.DoWithContext(ctx, api)
}

// PostStoryP holds parameters for the postStory method.
// See https://core.telegram.org/bots/api#poststory
type PostStoryP struct {
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

// PostStoryPhoto posts a story with a photo.
// See https://core.telegram.org/bots/api#poststory
func (api *API) PostStoryPhoto(params PostStoryP) (Story, error) {
	req := NewRequest[Story]("postStory", params)
	return req.Do(api)
}

// PostStoryPhotoWithContext is the context-aware variant of PostStoryPhoto.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#poststory
func (api *API) PostStoryPhotoWithContext(ctx context.Context, params PostStoryP) (Story, error) {
	req := NewRequest[Story]("postStory", params)
	return req.DoWithContext(ctx, api)
}

// PostStoryVideo posts a story with a video.
// See https://core.telegram.org/bots/api#poststory
func (api *API) PostStoryVideo(params PostStoryP) (Story, error) {
	req := NewRequest[Story]("postStory", params)
	return req.Do(api)
}

// PostStoryVideoWithContext is the context-aware variant of PostStoryVideo.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#poststory
func (api *API) PostStoryVideoWithContext(ctx context.Context, params PostStoryP) (Story, error) {
	req := NewRequest[Story]("postStory", params)
	return req.DoWithContext(ctx, api)
}

// RepostStoryP holds parameters for the repostStory method.
// See https://core.telegram.org/bots/api#repoststory
type RepostStoryP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	FromChatID           int64  `json:"from_chat_id"`
	FromStoryID          int    `json:"from_story_id"`
	ActivePeriod         int    `json:"active_period"`
	PostToChatPage       bool   `json:"post_to_chat_page,omitempty"`
	ProtectContent       bool   `json:"protect_content,omitempty"`
}

// RepostStory reposts a story from another chat.
// Returns the reposted story.
// See https://core.telegram.org/bots/api#repoststory
func (api *API) RepostStory(params RepostStoryP) (Story, error) {
	req := NewRequest[Story]("repostStory", params)
	return req.Do(api)
}

// RepostStoryWithContext is the context-aware variant of RepostStory.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#repoststory
func (api *API) RepostStoryWithContext(ctx context.Context, params RepostStoryP) (Story, error) {
	req := NewRequest[Story]("repostStory", params)
	return req.DoWithContext(ctx, api)
}

// EditStoryP holds parameters for the editStory method.
// See https://core.telegram.org/bots/api#editstory
type EditStoryP struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	StoryID              int               `json:"story_id"`
	Content              InputStoryContent `json:"content"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Areas           []StoryArea     `json:"areas,omitempty"`
}

// EditStory edits an existing story.
// Returns the updated story.
// See https://core.telegram.org/bots/api#editstory
func (api *API) EditStory(params EditStoryP) (Story, error) {
	req := NewRequest[Story]("editStory", params)
	return req.Do(api)
}

// EditStoryWithContext is the context-aware variant of EditStory.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#editstory
func (api *API) EditStoryWithContext(ctx context.Context, params EditStoryP) (Story, error) {
	req := NewRequest[Story]("editStory", params)
	return req.DoWithContext(ctx, api)
}

// DeleteStoryP holds parameters for the deleteStory method.
// See https://core.telegram.org/bots/api#deletestory
type DeleteStoryP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	StoryID              int    `json:"story_id"`
}

// DeleteStory deletes a story.
// Returns true on success.
// See https://core.telegram.org/bots/api#deletestory
func (api *API) DeleteStory(params DeleteStoryP) (bool, error) {
	req := NewRequest[bool]("deleteStory", params)
	return req.Do(api)
}

// DeleteStoryWithContext is the context-aware variant of DeleteStory.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletestory
func (api *API) DeleteStoryWithContext(ctx context.Context, params DeleteStoryP) (bool, error) {
	req := NewRequest[bool]("deleteStory", params)
	return req.DoWithContext(ctx, api)
}
