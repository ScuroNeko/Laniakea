package tgapi

type GetUserProfilePhotosP struct {
	UserID int `json:"user_id"`
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

func (api *Api) GetUserProfilePhotos(p GetUserProfilePhotosP) (UserProfilePhotos, error) {
	req := NewRequest[UserProfilePhotos]("getUserProfilePhotos", p)
	return req.Do(api)
}

type GetUserProfileAudiosP struct {
	UserID int `json:"user_id"`
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

func (api *Api) GetUserProfileAudios(p GetUserProfileAudiosP) (UserProfileAudios, error) {
	req := NewRequest[UserProfileAudios]("getUserProfileAudios", p)
	return req.Do(api)
}

type SetUserEmojiStatusP struct {
	UserID         int    `json:"user_id"`
	EmojiID        string `json:"emoji_status_custom_emoji_id,omitempty"`
	ExpirationDate int    `json:"emoji_status_expiration_date,omitempty"`
}

func (api *Api) SetUserEmojiStatus(p SetUserEmojiStatusP) (bool, error) {
	req := NewRequest[bool]("setUserEmojiStatus", p)
	return req.Do(api)
}

type GetUserGiftsP struct {
	UserID                      int    `json:"user_id"`
	ExcludeUnlimited            bool   `json:"exclude_unlimited,omitempty"`
	ExcludeLimitedUpgradable    bool   `json:"exclude_limited_upgradable,omitempty"`
	ExcludeLimitedNonUpgradable bool   `json:"exclude_limited_non_upgradable,omitempty"`
	ExcludeUnique               bool   `json:"exclude_unique,omitempty"`
	ExcludeFromBlockchain       bool   `json:"exclude_from_blockchain,omitempty"`
	SortByPrice                 bool   `json:"sort_by_price,omitempty"`
	Offset                      string `json:"offset,omitempty"`
	Limit                       int    `json:"limit,omitempty"`
}

func (api *Api) GetUserGifts(p GetUserGiftsP) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getUserGifts", p)
	return req.Do(api)
}
