package tgapi

type GetUserProfilePhotosP struct {
	UserID int `json:"user_id"`
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

func (api *API) GetUserProfilePhotos(params GetUserProfilePhotosP) (UserProfilePhotos, error) {
	req := NewRequest[UserProfilePhotos]("getUserProfilePhotos", params)
	return req.Do(api)
}

type GetUserProfileAudiosP struct {
	UserID int `json:"user_id"`
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

func (api *API) GetUserProfileAudios(params GetUserProfileAudiosP) (UserProfileAudios, error) {
	req := NewRequest[UserProfileAudios]("getUserProfileAudios", params)
	return req.Do(api)
}

type SetUserEmojiStatusP struct {
	UserID         int    `json:"user_id"`
	EmojiID        string `json:"emoji_status_custom_emoji_id,omitempty"`
	ExpirationDate int    `json:"emoji_status_expiration_date,omitempty"`
}

func (api *API) SetUserEmojiStatus(params SetUserEmojiStatusP) (bool, error) {
	req := NewRequest[bool]("setUserEmojiStatus", params)
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

func (api *API) GetUserGifts(params GetUserGiftsP) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getUserGifts", params)
	return req.Do(api)
}
