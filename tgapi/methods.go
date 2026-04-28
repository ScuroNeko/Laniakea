package tgapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"git.scuroneko.dev/scuroneko/laniakea/utils"
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
	req := NewRequest[User]("getMe", NoParams)
	return req.Do(api)
}

// GetMeWithContext is the context-aware variant of GetMe.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getme
func (api *API) GetMeWithContext(ctx context.Context) (User, error) {
	req := NewRequest[User]("getMe", NoParams)
	return req.DoWithContext(ctx, api)
}

// GetManagedBotToken holds parameters for the getManagedBotToken method.
// See https://core.telegram.org/bots/api#getmanagedbottoken
type GetManagedBotToken struct {
	UserID int64 `json:"user_id"`
}

// GetManagedBotToken returns the current token of a managed bot.
// See https://core.telegram.org/bots/api#getmanagedbottoken
func (api *API) GetManagedBotToken(params GetManagedBotToken) (string, error) {
	req := NewRequest[string]("getManagedBotToken", params)
	return req.Do(api)
}

// GetManagedBotTokenWithContext is the context-aware variant of GetManagedBotToken.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getmanagedbottoken
func (api *API) GetManagedBotTokenWithContext(ctx context.Context, params GetManagedBotToken) (string, error) {
	req := NewRequest[string]("getManagedBotToken", params)
	return req.DoWithContext(ctx, api)
}

// ReplaceManagedBotToken holds parameters for the replaceManagedBotToken method.
// See https://core.telegram.org/bots/api#replacemanagedbottoken
type ReplaceManagedBotToken struct {
	UserID int64 `json:"user_id"`
}

// ReplaceManagedBotToken replaces and returns the token of a managed bot.
// See https://core.telegram.org/bots/api#replacemanagedbottoken
func (api *API) ReplaceManagedBotToken(params ReplaceManagedBotToken) (string, error) {
	req := NewRequest[string]("replaceManagedBotToken", params)
	return req.Do(api)
}

// ReplaceManagedBotTokenWithContext is the context-aware variant of ReplaceManagedBotToken.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#replacemanagedbottoken
func (api *API) ReplaceManagedBotTokenWithContext(ctx context.Context, params ReplaceManagedBotToken) (string, error) {
	req := NewRequest[string]("replaceManagedBotToken", params)
	return req.DoWithContext(ctx, api)
}

// LogOut logs the bot out from the cloud Bot API server.
// Returns true on success.
// See https://core.telegram.org/bots/api#logout
func (api *API) LogOut() (bool, error) {
	req := NewRequest[bool, EmptyParams]("logOut", NoParams)
	return req.Do(api)
}

// LogOutWithContext is the context-aware variant of LogOut.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#logout
func (api *API) LogOutWithContext(ctx context.Context) (bool, error) {
	req := NewRequest[bool, EmptyParams]("logOut", NoParams)
	return req.DoWithContext(ctx, api)
}

// CloseRemote closes the bot instance on the local server.
// Returns true on success.
// See https://core.telegram.org/bots/api#close
func (api *API) CloseRemote() (bool, error) {
	req := NewRequest[bool, EmptyParams]("close", NoParams)
	return req.Do(api)
}

// CloseRemoteWithContext is the context-aware variant of CloseRemote.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#close
func (api *API) CloseRemoteWithContext(ctx context.Context) (bool, error) {
	req := NewRequest[bool, EmptyParams]("close", NoParams)
	return req.DoWithContext(ctx, api)
}

// GetUpdates receives incoming updates using long polling.
// See https://core.telegram.org/bots/api#getupdates
func (api *API) GetUpdates(params UpdateParams) ([]Update, error) {
	req := NewRequest[[]Update]("getUpdates", params)
	return req.Do(api)
}

// GetUpdatesWithContext is the context-aware variant of GetUpdates.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getupdates
func (api *API) GetUpdatesWithContext(ctx context.Context, params UpdateParams) ([]Update, error) {
	req := NewRequest[[]Update]("getUpdates", params)
	return req.DoWithContext(ctx, api)
}

// SetWebhook holds parameters for the setWebhook method.
// To upload a self-signed certificate, use Uploader.SetWebhook.
// See https://core.telegram.org/bots/api#setwebhook
type SetWebhook struct {
	URL                string       `json:"url"`
	IPAddress          string       `json:"ip_address,omitempty"`
	MaxConnections     int8         `json:"max_connections,omitempty"`
	AllowedUpdates     []UpdateType `json:"allowed_updates,omitempty"`
	DropPendingUpdates bool         `json:"drop_pending_updates,omitempty"`
	SecretToken        string       `json:"secret_token,omitempty"`
}

// SetWebhook sets a webhook URL for incoming updates.
// For certificate upload, use Uploader.SetWebhook.
// Returns true on success.
// See https://core.telegram.org/bots/api#setwebhook
func (api *API) SetWebhook(params SetWebhook) (bool, error) {
	req := NewRequest[bool]("setWebhook", params)
	return req.Do(api)
}

// SetWebhookWithContext is the context-aware variant of SetWebhook.
// It executes the same request but uses ctx for cancellation and deadlines.
// For certificate upload, use Uploader.SetWebhook.
// See https://core.telegram.org/bots/api#setwebhook
func (api *API) SetWebhookWithContext(ctx context.Context, params SetWebhook) (bool, error) {
	req := NewRequest[bool]("setWebhook", params)
	return req.DoWithContext(ctx, api)
}

// DeleteWebhook holds parameters for the deleteWebhook method.
// See https://core.telegram.org/bots/api#deletewebhook
type DeleteWebhook struct {
	DropPendingUpdates bool `json:"drop_pending_updates,omitempty"`
}

// DeleteWebhook removes the current webhook integration.
// Returns true on success.
// See https://core.telegram.org/bots/api#deletewebhook
func (api *API) DeleteWebhook(params DeleteWebhook) (bool, error) {
	req := NewRequest[bool]("deleteWebhook", params)
	return req.Do(api)
}

// DeleteWebhookWithContext is the context-aware variant of DeleteWebhook.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#deletewebhook
func (api *API) DeleteWebhookWithContext(ctx context.Context, params DeleteWebhook) (bool, error) {
	req := NewRequest[bool]("deleteWebhook", params)
	return req.DoWithContext(ctx, api)
}

// GetWebhookInfo returns the current webhook status.
// See https://core.telegram.org/bots/api#getwebhookinfo
func (api *API) GetWebhookInfo() (WebhookInfo, error) {
	req := NewRequest[WebhookInfo]("getWebhookInfo", NoParams)
	return req.Do(api)
}

// GetWebhookInfoWithContext is the context-aware variant of GetWebhookInfo.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getwebhookinfo
func (api *API) GetWebhookInfoWithContext(ctx context.Context) (WebhookInfo, error) {
	req := NewRequest[WebhookInfo]("getWebhookInfo", NoParams)
	return req.DoWithContext(ctx, api)
}

// GetFile holds parameters for the getFile method.
// See https://core.telegram.org/bots/api#getfile
type GetFile struct {
	FileID string `json:"file_id"`
}

// GetFile returns basic information about a file and prepares it for downloading.
// See https://core.telegram.org/bots/api#getfile
func (api *API) GetFile(params GetFile) (File, error) {
	req := NewRequest[File]("getFile", params)
	return req.Do(api)
}

// GetFileWithContext is the context-aware variant of GetFile.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getfile
func (api *API) GetFileWithContext(ctx context.Context, params GetFile) (File, error) {
	req := NewRequest[File]("getFile", params)
	return req.DoWithContext(ctx, api)
}

// GetFileByLink downloads a file from Telegram's file server using the provided file link.
// The link is usually obtained from File.FilePath.
// For large files, prefer OpenFileByLink or OpenFileByLinkWithContext to stream the response body.
// See https://core.telegram.org/bots/api#file
func (api *API) GetFileByLink(link string) ([]byte, error) {
	return api.getFileByLink(context.Background(), link)
}

// GetFileByLinkWithContext is the context-aware variant of GetFileByLink.
// It executes the same request but uses ctx for cancellation and deadlines.
// For large files, prefer OpenFileByLinkWithContext to stream the response body.
// See https://core.telegram.org/bots/api#file
func (api *API) GetFileByLinkWithContext(ctx context.Context, link string) ([]byte, error) {
	return api.getFileByLink(ctx, link)
}

// OpenFileByLink opens a streaming response body for a file hosted on Telegram's file server.
// The caller must close the returned ReadCloser.
// See https://core.telegram.org/bots/api#file
func (api *API) OpenFileByLink(link string) (io.ReadCloser, error) {
	return api.openFileByLink(context.Background(), link)
}

// OpenFileByLinkWithContext is the context-aware variant of OpenFileByLink.
// The caller must close the returned ReadCloser.
// See https://core.telegram.org/bots/api#file
func (api *API) OpenFileByLinkWithContext(ctx context.Context, link string) (io.ReadCloser, error) {
	return api.openFileByLink(ctx, link)
}

func (api *API) getFileByLink(ctx context.Context, link string) ([]byte, error) {
	body, err := api.openFileByLink(ctx, link)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = body.Close()
	}()
	return io.ReadAll(body)
}

func (api *API) openFileByLink(ctx context.Context, link string) (io.ReadCloser, error) {
	methodPrefix := ""
	if api.useTestServer {
		methodPrefix = "/test"
	}
	u := fmt.Sprintf("%s/file/bot%s%s/%s", api.apiURL, api.token, methodPrefix, link)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("Laniakea/%s", utils.VersionString))

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		defer func() {
			_ = res.Body.Close()
		}()
		body, readErr := io.ReadAll(io.LimitReader(res.Body, 4<<10))
		if readErr != nil {
			return nil, fmt.Errorf("unexpected status %d", res.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(body))
	}
	return res.Body, nil
}
