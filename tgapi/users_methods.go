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
