package tgapi

// GetStarTransactionsP holds parameters for the getStarTransactions method.
// See https://core.telegram.org/bots/api#getstartransactions
type GetStarTransactionsP struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// GetMyStarBalance returns the bot's Telegram Star balance.
// See https://core.telegram.org/bots/api#getmystarbalance
func (api *API) GetMyStarBalance() (StarAmount, error) {
	req := NewRequest[StarAmount]("getMyStarBalance", NoParams)
	return req.Do(api)
}

// GetStarTransactions returns Telegram Star transactions for the bot.
// See https://core.telegram.org/bots/api#getstartransactions
func (api *API) GetStarTransactions(params GetStarTransactionsP) (StarTransactions, error) {
	req := NewRequest[StarTransactions]("getStarTransactions", params)
	return req.Do(api)
}

// RefundStarPaymentP holds parameters for the refundStarPayment method.
// See https://core.telegram.org/bots/api#refundstarpayment
type RefundStarPaymentP struct {
	UserID                  int64  `json:"user_id"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
}

// RefundStarPayment refunds a successful Telegram Stars payment.
// Returns true on success.
// See https://core.telegram.org/bots/api#refundstarpayment
func (api *API) RefundStarPayment(params RefundStarPaymentP) (bool, error) {
	req := NewRequest[bool]("refundStarPayment", params)
	return req.Do(api)
}

// EditUserStarSubscriptionP holds parameters for the editUserStarSubscription method.
// See https://core.telegram.org/bots/api#edituserstarsubscription
type EditUserStarSubscriptionP struct {
	UserID                  int64  `json:"user_id"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	IsCanceled              bool   `json:"is_canceled"`
}

// EditUserStarSubscription cancels or re-enables a user star subscription extension.
// Returns true on success.
// See https://core.telegram.org/bots/api#edituserstarsubscription
func (api *API) EditUserStarSubscription(params EditUserStarSubscriptionP) (bool, error) {
	req := NewRequest[bool]("editUserStarSubscription", params)
	return req.Do(api)
}
