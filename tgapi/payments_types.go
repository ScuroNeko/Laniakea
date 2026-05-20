package tgapi

// LabeledPrice represents a price portion.
// Since: Bot API 3.0
// See https://core.telegram.org/bots/api#labeledprice
type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

// Invoice contains basic information about an invoice.
// Since: Bot API 3.0
type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
}

// ShippingQuery represents an incoming shipping query.
// Since: Bot API 3.0
// See https://core.telegram.org/bots/api#shippingquery
type ShippingQuery struct {
	ID              string          `json:"id"`
	From            User            `json:"from"`
	InvoicePayload  string          `json:"invoice_payload"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

// ShippingAddress represents a shipping address.
// Since: Bot API 3.0
// See https://core.telegram.org/bots/api#shippingaddress
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

// OrderInfo represents information about an order.
// Since: Bot API 3.0
// See https://core.telegram.org/bots/api#orderinfo
type OrderInfo struct {
	Name            string          `json:"name"`
	PhoneNumber     string          `json:"phone_number"`
	Email           string          `json:"email"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

// PreCheckoutQuery represents an incoming pre-checkout query.
// Since: Bot API 3.0
// See https://core.telegram.org/bots/api#precheckoutquery
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             User       `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

// PaidMediaPurchased represents a purchased paid media.
// Since: Bot API 7.10
// See https://core.telegram.org/bots/api#paidmediapurchased
type PaidMediaPurchased struct {
	From             User   `json:"from"`
	PaidMediaPayload string `json:"paid_media_payload"`
}

// ShippingOption represents one shipping option.
// Since: Bot API 3.0
// See https://core.telegram.org/bots/api#shippingoption
type ShippingOption struct {
	ID     string         `json:"id"`
	Title  string         `json:"title"`
	Prices []LabeledPrice `json:"prices"`
}

// SuccessfulPayment contains basic information about a successful payment.
// Since: Bot API 3.0
type SuccessfulPayment struct {
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
	InvoicePayload string `json:"invoice_payload"`

	SubscriptionExpirationDate int        `json:"subscription_expiration_date,omitempty"` // Since: Bot API 8.0
	IsRecurring                bool       `json:"is_recurring,omitempty"`                 // Since: Bot API 8.0
	IsFirstRecurring           bool       `json:"is_first_recurring,omitempty"`           // Since: Bot API 8.0
	ShippingOptionID           string     `json:"shipping_option_id,omitempty"`
	OrderInfo                  *OrderInfo `json:"order_info,omitempty"`

	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string `json:"proviced_payment_charge_id"`
}

// RefundedPayment contains basic information about a refunded payment.
// Since: Bot API 7.7
type RefundedPayment struct {
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
	InvoicePayload string `json:"invoice_payload"`

	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string `json:"proviced_payment_charge_id,omitempty"`
}
