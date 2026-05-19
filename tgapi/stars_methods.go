package tgapi

import "context"

// GetStarTransactions holds parameters for the getStarTransactions method.
// Since: Bot API 7.5
// See https://core.telegram.org/bots/api#getstartransactions
type GetStarTransactions struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// GetMyStarBalance returns the bot's Telegram Star balance.
// Since: Bot API 7.5
// See https://core.telegram.org/bots/api#getmystarbalance
func (api *API) GetMyStarBalance() (StarAmount, error) {
	req := NewRequest[StarAmount]("getMyStarBalance", NoParams)
	return req.Do(api)
}

// GetMyStarBalanceWithContext is the context-aware variant of GetMyStarBalance.
// Since: Bot API 7.5
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getmystarbalance
func (api *API) GetMyStarBalanceWithContext(ctx context.Context) (StarAmount, error) {
	req := NewRequest[StarAmount]("getMyStarBalance", NoParams)
	return req.DoWithContext(ctx, api)
}

// GetStarTransactions returns Telegram Star transactions for the bot.
// Since: Bot API 7.5
// See https://core.telegram.org/bots/api#getstartransactions
func (api *API) GetStarTransactions(params GetStarTransactions) (StarTransactions, error) {
	req := NewRequest[StarTransactions]("getStarTransactions", params)
	return req.Do(api)
}

// GetStarTransactionsWithContext is the context-aware variant of GetStarTransactions.
// Since: Bot API 7.5
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#getstartransactions
func (api *API) GetStarTransactionsWithContext(ctx context.Context, params GetStarTransactions) (StarTransactions, error) {
	req := NewRequest[StarTransactions]("getStarTransactions", params)
	return req.DoWithContext(ctx, api)
}

// RefundStarPayment holds parameters for the refundStarPayment method.
// Since: Bot API 7.4
// See https://core.telegram.org/bots/api#refundstarpayment
type RefundStarPayment struct {
	UserID                  int64  `json:"user_id"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
}

// RefundStarPayment refunds a successful Telegram Stars payment.
// Since: Bot API 7.4
// Returns true on success.
// See https://core.telegram.org/bots/api#refundstarpayment
func (api *API) RefundStarPayment(params RefundStarPayment) (bool, error) {
	req := NewRequest[bool]("refundStarPayment", params)
	return req.Do(api)
}

// RefundStarPaymentWithContext is the context-aware variant of RefundStarPayment.
// Since: Bot API 7.4
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#refundstarpayment
func (api *API) RefundStarPaymentWithContext(ctx context.Context, params RefundStarPayment) (bool, error) {
	req := NewRequest[bool]("refundStarPayment", params)
	return req.DoWithContext(ctx, api)
}

// EditUserStarSubscription holds parameters for the editUserStarSubscription method.
// Since: Bot API 8.0
// See https://core.telegram.org/bots/api#edituserstarsubscription
type EditUserStarSubscription struct {
	UserID                  int64  `json:"user_id"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	IsCanceled              bool   `json:"is_canceled"`
}

// EditUserStarSubscription cancels or re-enables a user star subscription extension.
// Since: Bot API 8.0
// Returns true on success.
// See https://core.telegram.org/bots/api#edituserstarsubscription
func (api *API) EditUserStarSubscription(params EditUserStarSubscription) (bool, error) {
	req := NewRequest[bool]("editUserStarSubscription", params)
	return req.Do(api)
}

// EditUserStarSubscriptionWithContext is the context-aware variant of EditUserStarSubscription.
// Since: Bot API 8.0
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#edituserstarsubscription
func (api *API) EditUserStarSubscriptionWithContext(ctx context.Context, params EditUserStarSubscription) (bool, error) {
	req := NewRequest[bool]("editUserStarSubscription", params)
	return req.DoWithContext(ctx, api)
}
