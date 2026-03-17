package tgapi

// SetPassportDataErrorsP holds parameters for the setPassportDataErrors method.
// See https://core.telegram.org/bots/api#setpassportdataerrors
type SetPassportDataErrorsP struct {
	UserID int64                  `json:"user_id"`
	Errors []PassportElementError `json:"errors"`
}

// SetPassportDataErrors informs a user about Telegram Passport data errors.
// Returns true on success.
// See https://core.telegram.org/bots/api#setpassportdataerrors
func (api *API) SetPassportDataErrors(params SetPassportDataErrorsP) (bool, error) {
	req := NewRequest[bool]("setPassportDataErrors", params)
	return req.Do(api)
}
