package tgapi

import "context"

// SetPassportDataErrors holds parameters for the setPassportDataErrors method.
// Since: Bot API 4.0
// See https://core.telegram.org/bots/api#setpassportdataerrors
type SetPassportDataErrors struct {
	UserID int64                  `json:"user_id"`
	Errors []PassportElementError `json:"errors"`
}

// SetPassportDataErrors informs a user about Telegram Passport data errors.
// Since: Bot API 4.0
// Returns true on success.
// See https://core.telegram.org/bots/api#setpassportdataerrors
func (api *API) SetPassportDataErrors(params SetPassportDataErrors) (bool, error) {
	req := NewRequest[bool]("setPassportDataErrors", params)
	return req.Do(api)
}

// SetPassportDataErrorsWithContext is the context-aware variant of SetPassportDataErrors.
// Since: Bot API 4.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#setpassportdataerrors
func (api *API) SetPassportDataErrorsWithContext(ctx context.Context, params SetPassportDataErrors) (bool, error) {
	req := NewRequest[bool]("setPassportDataErrors", params)
	return req.DoWithContext(ctx, api)
}
