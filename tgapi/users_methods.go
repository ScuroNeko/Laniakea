package tgapi

// GetUserProfilePhotosP holds parameters for the GetUserProfilePhotos method.
// See https://core.telegram.org/bots/api#getuserprofilephotos
type GetUserProfilePhotosP struct {
	UserID int `json:"user_id"`
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// GetUserProfilePhotos returns a list of profile pictures for a user.
// See https://core.telegram.org/bots/api#getuserprofilephotos
func (api *API) GetUserProfilePhotos(params GetUserProfilePhotosP) (UserProfilePhotos, error) {
	req := NewRequest[UserProfilePhotos]("getUserProfilePhotos", params)
	return req.Do(api)
}

// GetUserProfileAudiosP holds parameters for the GetUserProfileAudios method.
// See https://core.telegram.org/bots/api#getuserprofileaudios
type GetUserProfileAudiosP struct {
	UserID int `json:"user_id"`
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// GetUserProfileAudios returns a list of profile audios for a user.
// See https://core.telegram.org/bots/api#getuserprofileaudios
func (api *API) GetUserProfileAudios(params GetUserProfileAudiosP) (UserProfileAudios, error) {
	req := NewRequest[UserProfileAudios]("getUserProfileAudios", params)
	return req.Do(api)
}

// SetUserEmojiStatusP holds parameters for the SetUserEmojiStatus method.
// See https://core.telegram.org/bots/api#setuseremojistatus
type SetUserEmojiStatusP struct {
	UserID         int    `json:"user_id"`
	EmojiID        string `json:"emoji_status_custom_emoji_id,omitempty"`
	ExpirationDate int    `json:"emoji_status_expiration_date,omitempty"`
}

// SetUserEmojiStatus sets a custom emoji status for a user.
// Returns true on success.
// See https://core.telegram.org/bots/api#setuseremojistatus
func (api *API) SetUserEmojiStatus(params SetUserEmojiStatusP) (bool, error) {
	req := NewRequest[bool]("setUserEmojiStatus", params)
	return req.Do(api)
}

// GetUserGiftsP holds parameters for the GetUserGifts method.
// See https://core.telegram.org/bots/api#getusergifts
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

// GetUserGifts returns gifts owned by a user.
// See https://core.telegram.org/bots/api#getusergifts
func (api *API) GetUserGifts(params GetUserGiftsP) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getUserGifts", params)
	return req.Do(api)
}
