package tgapi

type VerifyUserP struct {
	UserID            int    `json:"user_id"`
	CustomDescription string `json:"custom_description,omitempty"`
}

func (api *Api) VerifyUser(params VerifyUserP) (bool, error) {
	req := NewRequest[bool]("verifyUser", params)
	return req.Do(api)
}

type VerifyChatP struct {
	ChatID            int    `json:"chat_id"`
	CustomDescription string `json:"custom_description,omitempty"`
}

func (api *Api) VerifyChat(params VerifyChatP) (bool, error) {
	req := NewRequest[bool]("verifyChat", params)
	return req.Do(api)
}

type RemoveUserVerificationP struct {
	UserID int `json:"user_id"`
}

func (api *Api) RemoveUserVerification(params RemoveUserVerificationP) (bool, error) {
	req := NewRequest[bool]("removeUserVerification", params)
	return req.Do(api)
}

type RemoveChatVerificationP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) RemoveChatVerification(params RemoveChatVerificationP) (bool, error) {
	req := NewRequest[bool]("removeChatVerification", params)
	return req.Do(api)
}

type ReadBusinessMessageP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	ChatID               int    `json:"chat_id"`
	MessageID            int    `json:"message_id"`
}

func (api *Api) ReadBusinessMessage(params ReadBusinessMessageP) (bool, error) {
	req := NewRequest[bool]("readBusinessMessage", params)
	return req.Do(api)
}

type DeleteBusinessMessageP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	MessageIDs           []int  `json:"message_ids"`
}

func (api *Api) DeleteBusinessMessage(params DeleteBusinessMessageP) (bool, error) {
	req := NewRequest[bool]("deleteBusinessMessage", params)
	return req.Do(api)
}

type SetBusinessAccountNameP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name,omitempty"`
}

func (api *Api) SetBusinessAccountName(params SetBusinessAccountNameP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountName", params)
	return req.Do(api)
}

type SetBusinessAccountUsernameP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Username             string `json:"username,omitempty"`
}

func (api *Api) SetBusinessAccountUsername(params SetBusinessAccountUsernameP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountUsername", params)
	return req.Do(api)
}

type SetBusinessAccountBioP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Bio                  string `json:"bio,omitempty"`
}

func (api *Api) SetBusinessAccountBio(params SetBusinessAccountBioP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountBio", params)
	return req.Do(api)
}

type SetBusinessAccountProfilePhoto struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	Photo                InputProfilePhoto `json:"photo,omitempty"`
	IsPublic             bool              `json:"is_public,omitempty"`
}

func (api *Api) SetBusinessAccountProfilePhoto(params SetBusinessAccountProfilePhoto) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountProfilePhoto", params)
	return req.Do(api)
}

type RemoveBusinessAccountProfilePhotoP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	IsPublic             bool   `json:"is_public,omitempty"`
}

func (api *Api) RemoveBusinessAccountProfilePhoto(params RemoveBusinessAccountProfilePhotoP) (bool, error) {
	req := NewRequest[bool]("removeBusinessAccountProfilePhoto", params)
	return req.Do(api)
}

type SetBusinessAccountGiftSettingsP struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	ShowGiftButton       bool              `json:"show_gift_button"`
	AcceptedGiftTypes    AcceptedGiftTypes `json:"accepted_gift_types"`
}

func (api *Api) SetBusinessAccountGiftSettings(params SetBusinessAccountGiftSettingsP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountGiftSettings", params)
	return req.Do(api)
}

type GetBusinessAccountStarBalanceP struct {
	BusinessConnectionID string `json:"business_connection_id"`
}

func (api *Api) GetBusinessAccountStarBalance(params GetBusinessAccountStarBalanceP) (StarAmount, error) {
	req := NewRequest[StarAmount]("getBusinessAccountGiftSettings", params)
	return req.Do(api)
}

type TransferBusinessAccountStartP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	StarCount            int    `json:"star_count"`
}

func (api *Api) TransferBusinessAccountStart(params TransferBusinessAccountStartP) (bool, error) {
	req := NewRequest[bool]("transferBusinessAccountStart", params)
	return req.Do(api)
}

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

func (api *Api) GetBusinessAccountGifts(params GetBusinessAccountGiftsP) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getBusinessAccountGifts", params)
	return req.Do(api)
}

type ConvertGiftToStarsP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
}

func (api *Api) ConvertGiftToStars(params ConvertGiftToStarsP) (bool, error) {
	req := NewRequest[bool]("convertGiftToStars", params)
	return req.Do(api)
}

type UpgradeGiftP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
	KeepOriginalDetails  bool   `json:"keep_original_details,omitempty"`
	StarCount            int    `json:"star_count,omitempty"`
}

func (api *Api) UpgradeGift(params UpgradeGiftP) (bool, error) {
	req := NewRequest[bool]("upgradeGift", params)
	return req.Do(api)
}

type TransferGiftP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
	NewOwnerChatID       int    `json:"new_owner_chat_id"`
	StarCount            int    `json:"star_count,omitempty"`
}

func (api *Api) TransferGift(params TransferGiftP) (bool, error) {
	req := NewRequest[bool]("transferGift", params)
	return req.Do(api)
}

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

func (api *Api) PostStoryPhoto(params PostStoryP) (Story, error) {
	req := NewRequest[Story]("postStory", params)
	return req.Do(api)
}
func (api *Api) PostStoryVideo(params PostStoryP) (Story, error) {
	req := NewRequest[Story]("postStory", params)
	return req.Do(api)
}

type RepostStoryP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	FromChatID           int    `json:"from_chat_id"`
	FromStoryID          int    `json:"from_story_id"`
	ActivePeriod         int    `json:"active_period"`
	PostToChatPage       bool   `json:"post_to_chat_page,omitempty"`
	ProtectContent       bool   `json:"protect_content,omitempty"`
}

func (api *Api) RepostStory(params RepostStoryP) (Story, error) {
	req := NewRequest[Story]("repostStory", params)
	return req.Do(api)
}

type EditStoryP struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	StoryID              int               `json:"story_id"`
	Content              InputStoryContent `json:"content"`

	Caption         string          `json:"caption,omitempty"`
	ParseMode       ParseMode       `json:"parse_mode,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	Areas           []StoryArea     `json:"areas,omitempty"`
}

func (api *Api) EditStory(params EditStoryP) (Story, error) {
	req := NewRequest[Story]("editStory", params)
	return req.Do(api)
}

type DeleteStoryP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	StoryID              int    `json:"story_id"`
}

func (api *Api) DeleteStory(params DeleteStoryP) (bool, error) {
	req := NewRequest[bool]("deleteStory", params)
	return req.Do(api)
}
