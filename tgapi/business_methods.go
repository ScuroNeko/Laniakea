package tgapi

type VerifyUserP struct {
	UserID            int    `json:"user_id"`
	CustomDescription string `json:"custom_description,omitempty"`
}

func (api *Api) VerifyUser(p VerifyUserP) (bool, error) {
	req := NewRequest[bool]("verifyUser", p)
	return req.Do(api)
}

type VerifyChatP struct {
	ChatID            int    `json:"chat_id"`
	CustomDescription string `json:"custom_description,omitempty"`
}

func (api *Api) VerifyChat(p VerifyChatP) (bool, error) {
	req := NewRequest[bool]("verifyChat", p)
	return req.Do(api)
}

type RemoveUserVerificationP struct {
	UserID int `json:"user_id"`
}

func (api *Api) RemoveUserVerification(p RemoveUserVerificationP) (bool, error) {
	req := NewRequest[bool]("removeUserVerification", p)
	return req.Do(api)
}

type RemoveChatVerificationP struct {
	ChatID int `json:"chat_id"`
}

func (api *Api) RemoveChatVerification(p RemoveChatVerificationP) (bool, error) {
	req := NewRequest[bool]("removeChatVerification", p)
	return req.Do(api)
}

type ReadBusinessMessageP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	ChatID               int    `json:"chat_id"`
	MessageID            int    `json:"message_id"`
}

func (api *Api) ReadBusinessMessage(p ReadBusinessMessageP) (bool, error) {
	req := NewRequest[bool]("readBusinessMessage", p)
	return req.Do(api)
}

type DeleteBusinessMessageP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	MessageIDs           []int  `json:"message_ids"`
}

func (api *Api) DeleteBusinessMessage(p DeleteBusinessMessageP) (bool, error) {
	req := NewRequest[bool]("deleteBusinessMessage", p)
	return req.Do(api)
}

type SetBusinessAccountNameP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name,omitempty"`
}

func (api *Api) SetBusinessAccountName(p SetBusinessAccountNameP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountName", p)
	return req.Do(api)
}

type SetBusinessAccountUsernameP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Username             string `json:"username,omitempty"`
}

func (api *Api) SetBusinessAccountUsername(p SetBusinessAccountUsernameP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountUsername", p)
	return req.Do(api)
}

type SetBusinessAccountBioP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Bio                  string `json:"bio,omitempty"`
}

func (api *Api) SetBusinessAccountBio(p SetBusinessAccountBioP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountBio", p)
	return req.Do(api)
}

type SetBusinessAccountProfilePhoto struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	Photo                InputProfilePhoto `json:"photo,omitempty"`
	IsPublic             bool              `json:"is_public,omitempty"`
}

func (api *Api) SetBusinessAccountProfilePhoto(p SetBusinessAccountProfilePhoto) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountProfilePhoto", p)
	return req.Do(api)
}

type RemoveBusinessAccountProfilePhotoP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	IsPublic             bool   `json:"is_public,omitempty"`
}

func (api *Api) RemoveBusinessAccountProfilePhoto(p RemoveBusinessAccountProfilePhotoP) (bool, error) {
	req := NewRequest[bool]("removeBusinessAccountProfilePhoto", p)
	return req.Do(api)
}

type SetBusinessAccountGiftSettingsP struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	ShowGiftButton       bool              `json:"show_gift_button"`
	AcceptedGiftTypes    AcceptedGiftTypes `json:"accepted_gift_types"`
}

func (api *Api) SetBusinessAccountGiftSettings(p SetBusinessAccountGiftSettingsP) (bool, error) {
	req := NewRequest[bool]("setBusinessAccountGiftSettings", p)
	return req.Do(api)
}

type GetBusinessAccountStarBalanceP struct {
	BusinessConnectionID string `json:"business_connection_id"`
}

func (api *Api) GetBusinessAccountStarBalance(p GetBusinessAccountStarBalanceP) (StarAmount, error) {
	req := NewRequest[StarAmount]("getBusinessAccountGiftSettings", p)
	return req.Do(api)
}

type TransferBusinessAccountStartP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	StarCount            int    `json:"star_count"`
}

func (api *Api) TransferBusinessAccountStart(p TransferBusinessAccountStartP) (bool, error) {
	req := NewRequest[bool]("transferBusinessAccountStart", p)
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

func (api *Api) GetBusinessAccountGifts(p GetBusinessAccountGiftsP) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getBusinessAccountGifts", p)
	return req.Do(api)
}

type ConvertGiftToStarsP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
}

func (api *Api) ConvertGiftToStars(p ConvertGiftToStarsP) (bool, error) {
	req := NewRequest[bool]("convertGiftToStars", p)
	return req.Do(api)
}

type UpgradeGiftP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
	KeepOriginalDetails  bool   `json:"keep_original_details,omitempty"`
	StarCount            int    `json:"star_count,omitempty"`
}

func (api *Api) UpgradeGift(p UpgradeGiftP) (bool, error) {
	req := NewRequest[bool]("upgradeGift", p)
	return req.Do(api)
}

type TransferGiftP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
	NewOwnerChatID       int    `json:"new_owner_chat_id"`
	StarCount            int    `json:"star_count,omitempty"`
}

func (api *Api) TransferGift(p TransferGiftP) (bool, error) {
	req := NewRequest[bool]("transferGift", p)
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

func (api *Api) PostStoryPhoto(p PostStoryP) (Story, error) {
	req := NewRequest[Story]("postStory", p)
	return req.Do(api)
}
func (api *Api) PostStoryVideo(p PostStoryP) (Story, error) {
	req := NewRequest[Story]("postStory", p)
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

func (api *Api) RepostStory(p RepostStoryP) (Story, error) {
	req := NewRequest[Story]("repostStory", p)
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

func (api *Api) EditStory(p EditStoryP) (Story, error) {
	req := NewRequest[Story]("editStory", p)
	return req.Do(api)
}

type DeleteStoryP struct {
	BusinessConnectionID string `json:"business_connection_id"`
	StoryID              int    `json:"story_id"`
}

func (api *Api) DeleteStory(p DeleteStoryP) (bool, error) {
	req := NewRequest[bool]("deleteStory", p)
	return req.Do(api)
}
