package tgapi

import (
	"fmt"
	"io"
	"net/http"
)

// ParseMode represents the text formatting mode for message parsing.
type ParseMode string

const (
	// ParseMDV2 enables MarkdownV2 style parsing.
	ParseMDV2 ParseMode = "MarkdownV2"
	// ParseHTML enables HTML style parsing.
	ParseHTML ParseMode = "HTML"
	// ParseMD enables legacy Markdown style parsing.
	ParseMD ParseMode = "Markdown"
	// ParseNone disables any parsing.
	ParseNone ParseMode = "None"
)

// EmptyParams is a placeholder for methods that take no parameters.
type EmptyParams struct{}

// NoParams is a convenient instance of EmptyParams.
var NoParams = EmptyParams{}

// UpdateParams holds parameters for the getUpdates method.
// See https://core.telegram.org/bots/api#getupdates
type UpdateParams struct {
	Offset         *int         `json:"offset,omitempty"`
	Limit          *int         `json:"limit,omitempty"`
	Timeout        *int         `json:"timeout,omitempty"`
	AllowedUpdates []UpdateType `json:"allowed_updates"`
}

// GetMe returns basic information about the bot.
// See https://core.telegram.org/bots/api#getme
func (api *API) GetMe() (User, error) {
	req := NewRequest[User, EmptyParams]("getMe", NoParams)
	return req.Do(api)
}

// LogOut logs the bot out from the cloud Bot API server.
// Returns true on success.
// See https://core.telegram.org/bots/api#logout
func (api *API) LogOut() (bool, error) {
	req := NewRequest[bool, EmptyParams]("logOut", NoParams)
	return req.Do(api)
}

// Close closes the bot instance on the local server.
// Returns true on success.
// See https://core.telegram.org/bots/api#close
func (api *API) Close() (bool, error) {
	req := NewRequest[bool, EmptyParams]("close", NoParams)
	return req.Do(api)
}

// GetUpdates receives incoming updates using long polling.
// See https://core.telegram.org/bots/api#getupdates
func (api *API) GetUpdates(params UpdateParams) ([]Update, error) {
	req := NewRequest[[]Update]("getUpdates", params)
	return req.Do(api)
}

// GetFileP holds parameters for the getFile method.
// See https://core.telegram.org/bots/api#getfile
type GetFileP struct {
	FileId string `json:"file_id"`
}

// GetFile returns basic information about a file and prepares it for downloading.
// See https://core.telegram.org/bots/api#getfile
func (api *API) GetFile(params GetFileP) (File, error) {
	req := NewRequest[File]("getFile", params)
	return req.Do(api)
}

// GetFileByLink downloads a file from Telegram's file server using the provided file link.
// The link is usually obtained from File.FilePath.
// See https://core.telegram.org/bots/api#file
func (api *API) GetFileByLink(link string) ([]byte, error) {
	u := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", api.token, link)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	return io.ReadAll(res.Body)
}
