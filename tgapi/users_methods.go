package tgapi

import "context"

// GetUserProfilePhotos holds parameters for the GetUserProfilePhotos method.
// See https://core.telegram.org/bots/api#getuserprofilephotos
type GetUserProfilePhotos struct {
	UserID int64 `json:"user_id"`
	Offset int   `json:"offset,omitempty"`
	Limit  int   `json:"limit,omitempty"`
}

// GetUserProfilePhotos returns a list of profile pictures for a user.
// See https://core.telegram.org/bots/api#getuserprofilephotos
func (api *API) GetUserProfilePhotos(params GetUserProfilePhotos) (UserProfilePhotos, error) {
	req := NewRequest[UserProfilePhotos]("getUserProfilePhotos", params)
	return req.Do(api)
}

// GetUserProfilePhotosWithContext is the context-aware variant of GetUserProfilePhotos.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getuserprofilephotos
func (api *API) GetUserProfilePhotosWithContext(ctx context.Context, params GetUserProfilePhotos) (UserProfilePhotos, error) {
	req := NewRequest[UserProfilePhotos]("getUserProfilePhotos", params)
	return req.DoWithContext(ctx, api)
}

// GetUserProfileAudios holds parameters for the GetUserProfileAudios method.
// See https://core.telegram.org/bots/api#getuserprofileaudios
type GetUserProfileAudios struct {
	UserID int64 `json:"user_id"`
	Offset int   `json:"offset,omitempty"`
	Limit  int   `json:"limit,omitempty"`
}

// GetUserProfileAudios returns a list of profile audios for a user.
// See https://core.telegram.org/bots/api#getuserprofileaudios
func (api *API) GetUserProfileAudios(params GetUserProfileAudios) (UserProfileAudios, error) {
	req := NewRequest[UserProfileAudios]("getUserProfileAudios", params)
	return req.Do(api)
}

// GetUserProfileAudiosWithContext is the context-aware variant of GetUserProfileAudios.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getuserprofileaudios
func (api *API) GetUserProfileAudiosWithContext(ctx context.Context, params GetUserProfileAudios) (UserProfileAudios, error) {
	req := NewRequest[UserProfileAudios]("getUserProfileAudios", params)
	return req.DoWithContext(ctx, api)
}

// SetUserEmojiStatus holds parameters for the SetUserEmojiStatus method.
// See https://core.telegram.org/bots/api#setuseremojistatus
type SetUserEmojiStatus struct {
	UserID         int64  `json:"user_id"`
	EmojiID        string `json:"emoji_status_custom_emoji_id,omitempty"`
	ExpirationDate int    `json:"emoji_status_expiration_date,omitempty"`
}

// SetUserEmojiStatus sets a custom emoji status for a user.
// Returns true on success.
// See https://core.telegram.org/bots/api#setuseremojistatus
func (api *API) SetUserEmojiStatus(params SetUserEmojiStatus) (bool, error) {
	req := NewRequest[bool]("setUserEmojiStatus", params)
	return req.Do(api)
}

// SetUserEmojiStatusWithContext is the context-aware variant of SetUserEmojiStatus.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setuseremojistatus
func (api *API) SetUserEmojiStatusWithContext(ctx context.Context, params SetUserEmojiStatus) (bool, error) {
	req := NewRequest[bool]("setUserEmojiStatus", params)
	return req.DoWithContext(ctx, api)
}

// GetUserGifts holds parameters for the GetUserGifts method.
// See https://core.telegram.org/bots/api#getusergifts
type GetUserGifts struct {
	UserID                      int64  `json:"user_id"`
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
func (api *API) GetUserGifts(params GetUserGifts) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getUserGifts", params)
	return req.Do(api)
}

// GetUserGiftsWithContext is the context-aware variant of GetUserGifts.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getusergifts
func (api *API) GetUserGiftsWithContext(ctx context.Context, params GetUserGifts) (OwnedGifts, error) {
	req := NewRequest[OwnedGifts]("getUserGifts", params)
	return req.DoWithContext(ctx, api)
}
