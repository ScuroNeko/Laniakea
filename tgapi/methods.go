package tgapi

type ParseMode string

const (
	ParseMDV2 ParseMode = "MarkdownV2"
	ParseHTML ParseMode = "HTML"
	ParseMD   ParseMode = "Markdown"
)

type EmptyParams struct{}

var NoParams = EmptyParams{}

type UpdateParams struct {
	Offset         *int         `json:"offset,omitempty"`
	Limit          *int         `json:"limit,omitempty"`
	Timeout        *int         `json:"timeout,omitempty"`
	AllowedUpdates []UpdateType `json:"allowed_updates"`
}

func (api *Api) GetMe() (User, error) {
	req := NewRequest[User, EmptyParams]("getMe", NoParams)
	return req.Do(api)
}
func (api *Api) LogOut() (bool, error) {
	req := NewRequest[bool, EmptyParams]("logOut", NoParams)
	return req.Do(api)
}
func (api *Api) Close() (bool, error) {
	req := NewRequest[bool, EmptyParams]("close", NoParams)
	return req.Do(api)
}
func (api *Api) GetUpdates(params UpdateParams) ([]Update, error) {
	req := NewRequest[[]Update]("getUpdates", params)
	return req.Do(api)
}

type GetFileP struct {
	FileId string `json:"file_id"`
}

func (api *Api) GetFile(params GetFileP) (File, error) {
	req := NewRequest[File]("getFile", params)
	return req.Do(api)
}
