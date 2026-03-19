package tgapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"git.nix13.pw/scuroneko/laniakea/utils"
)

// UpdateParams holds parameters for the getUpdates method.
// See https://core.telegram.org/bots/api#getupdates
type UpdateParams struct {
	Offset         *int         `json:"offset,omitempty"`
	Limit          *int         `json:"limit,omitempty"`
	Timeout        *int         `json:"timeout,omitempty"`
	AllowedUpdates []UpdateType `json:"allowed_updates,omitempty"`
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

func (api *API) GetUpdatesWithContext(ctx context.Context, params UpdateParams) ([]Update, error) {
	req := NewRequest[[]Update]("getUpdates", params)
	return req.DoWithContext(ctx, api)
}

// SetWebhookP holds parameters for the setWebhook method.
// See https://core.telegram.org/bots/api#setwebhook
type SetWebhookP struct {
	URL                string       `json:"url"`
	Certificate        string       `json:"certificate,omitempty"`
	IPAddress          string       `json:"ip_address,omitempty"`
	MaxConnections     int          `json:"max_connections,omitempty"`
	AllowedUpdates     []UpdateType `json:"allowed_updates,omitempty"`
	DropPendingUpdates bool         `json:"drop_pending_updates,omitempty"`
	SecretToken        string       `json:"secret_token,omitempty"`
}

// SetWebhook sets a webhook URL for incoming updates.
// Returns true on success.
// See https://core.telegram.org/bots/api#setwebhook
func (api *API) SetWebhook(params SetWebhookP) (bool, error) {
	req := NewRequest[bool]("setWebhook", params)
	return req.Do(api)
}

// DeleteWebhookP holds parameters for the deleteWebhook method.
// See https://core.telegram.org/bots/api#deletewebhook
type DeleteWebhookP struct {
	DropPendingUpdates bool `json:"drop_pending_updates,omitempty"`
}

// DeleteWebhook removes the current webhook integration.
// Returns true on success.
// See https://core.telegram.org/bots/api#deletewebhook
func (api *API) DeleteWebhook(params DeleteWebhookP) (bool, error) {
	req := NewRequest[bool]("deleteWebhook", params)
	return req.Do(api)
}

// GetWebhookInfo returns the current webhook status.
// See https://core.telegram.org/bots/api#getwebhookinfo
func (api *API) GetWebhookInfo() (WebhookInfo, error) {
	req := NewRequest[WebhookInfo]("getWebhookInfo", NoParams)
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
	methodPrefix := ""
	if api.useTestServer {
		methodPrefix = "/test"
	}
	u := fmt.Sprintf("%s/file/bot%s%s/%s", api.apiUrl, api.token, methodPrefix, link)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		body, readErr := io.ReadAll(io.LimitReader(res.Body, 4<<10))
		if readErr != nil {
			return nil, fmt.Errorf("unexpected status %d", res.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(body))
	}
	return io.ReadAll(res.Body)
}
