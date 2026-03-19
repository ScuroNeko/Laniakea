package tgapi

import "context"

// SendInvoiceP holds parameters for the sendInvoice method.
// See https://core.telegram.org/bots/api#sendinvoice
type SendInvoiceP struct {
	ChatID                int64 `json:"chat_id"`
	MessageThreadID       int   `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int   `json:"direct_messages_topic_id,omitempty"`

	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Payload       string         `json:"payload"`
	ProviderToken string         `json:"provider_token,omitempty"`
	Currency      string         `json:"currency"`
	Prices        []LabeledPrice `json:"prices"`

	MaxTipAmount        int    `json:"max_tip_amount,omitempty"`
	SuggestedTipAmounts []int  `json:"suggested_tip_amounts,omitempty"`
	StartParameter      string `json:"start_parameter,omitempty"`
	ProviderData        string `json:"provider_data,omitempty"`
	PhotoURL            string `json:"photo_url,omitempty"`
	PhotoSize           int    `json:"photo_size,omitempty"`
	PhotoWidth          int    `json:"photo_width,omitempty"`
	PhotoHeight         int    `json:"photo_height,omitempty"`
	NeedName            bool   `json:"need_name,omitempty"`
	NeedPhoneNumber     bool   `json:"need_phone_number,omitempty"`
	NeedEmail           bool   `json:"need_email,omitempty"`
	NeedShippingAddress bool   `json:"need_shipping_address,omitempty"`
	SendPhoneToProvider bool   `json:"send_phone_number_to_provider,omitempty"`
	SendEmailToProvider bool   `json:"send_email_to_provider,omitempty"`
	IsFlexible          bool   `json:"is_flexible,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
	AllowPaidBroadcast  bool   `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID     string `json:"message_effect_id,omitempty"`

	SuggestedPostParameters *SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             *InlineKeyboardMarkup    `json:"reply_markup,omitempty"`
}

// SendInvoice sends an invoice.
// See https://core.telegram.org/bots/api#sendinvoice
func (api *API) SendInvoice(params SendInvoiceP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendInvoice", params, params.ChatID)
	return req.Do(api)
}

// SendInvoiceWithContext is the context-aware variant of SendInvoice.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#sendinvoice
func (api *API) SendInvoiceWithContext(ctx context.Context, params SendInvoiceP) (Message, error) {
	req := NewRequestWithChatID[Message]("sendInvoice", params, params.ChatID)
	return req.DoWithContext(ctx, api)
}

// CreateInvoiceLinkP holds parameters for the createInvoiceLink method.
// See https://core.telegram.org/bots/api#createinvoicelink
type CreateInvoiceLinkP struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`

	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Payload       string         `json:"payload"`
	ProviderToken string         `json:"provider_token,omitempty"`
	Currency      string         `json:"currency"`
	Prices        []LabeledPrice `json:"prices"`

	SubscriptionPeriod  int    `json:"subscription_period,omitempty"`
	MaxTipAmount        int    `json:"max_tip_amount,omitempty"`
	SuggestedTipAmounts []int  `json:"suggested_tip_amounts,omitempty"`
	ProviderData        string `json:"provider_data,omitempty"`
	PhotoURL            string `json:"photo_url,omitempty"`
	PhotoSize           int    `json:"photo_size,omitempty"`
	PhotoWidth          int    `json:"photo_width,omitempty"`
	PhotoHeight         int    `json:"photo_height,omitempty"`
	NeedName            bool   `json:"need_name,omitempty"`
	NeedPhoneNumber     bool   `json:"need_phone_number,omitempty"`
	NeedEmail           bool   `json:"need_email,omitempty"`
	NeedShippingAddress bool   `json:"need_shipping_address,omitempty"`
	SendPhoneToProvider bool   `json:"send_phone_number_to_provider,omitempty"`
	SendEmailToProvider bool   `json:"send_email_to_provider,omitempty"`
	IsFlexible          bool   `json:"is_flexible,omitempty"`
}

// CreateInvoiceLink creates an invoice link.
// See https://core.telegram.org/bots/api#createinvoicelink
func (api *API) CreateInvoiceLink(params CreateInvoiceLinkP) (string, error) {
	req := NewRequest[string]("createInvoiceLink", params)
	return req.Do(api)
}

// CreateInvoiceLinkWithContext is the context-aware variant of CreateInvoiceLink.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#createinvoicelink
func (api *API) CreateInvoiceLinkWithContext(ctx context.Context, params CreateInvoiceLinkP) (string, error) {
	req := NewRequest[string]("createInvoiceLink", params)
	return req.DoWithContext(ctx, api)
}

// AnswerShippingQueryP holds parameters for the answerShippingQuery method.
// See https://core.telegram.org/bots/api#answershippingquery
type AnswerShippingQueryP struct {
	ShippingQueryID string           `json:"shipping_query_id"`
	OK              bool             `json:"ok"`
	ShippingOptions []ShippingOption `json:"shipping_options,omitempty"`
	ErrorMessage    string           `json:"error_message,omitempty"`
}

// AnswerShippingQuery answers a shipping query.
// Returns true on success.
// See https://core.telegram.org/bots/api#answershippingquery
func (api *API) AnswerShippingQuery(params AnswerShippingQueryP) (bool, error) {
	req := NewRequest[bool]("answerShippingQuery", params)
	return req.Do(api)
}

// AnswerShippingQueryWithContext is the context-aware variant of AnswerShippingQuery.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#answershippingquery
func (api *API) AnswerShippingQueryWithContext(ctx context.Context, params AnswerShippingQueryP) (bool, error) {
	req := NewRequest[bool]("answerShippingQuery", params)
	return req.DoWithContext(ctx, api)
}

// AnswerPreCheckoutQueryP holds parameters for the answerPreCheckoutQuery method.
// See https://core.telegram.org/bots/api#answerprecheckoutquery
type AnswerPreCheckoutQueryP struct {
	PreCheckoutQueryID string `json:"pre_checkout_query_id"`
	OK                 bool   `json:"ok"`
	ErrorMessage       string `json:"error_message,omitempty"`
}

// AnswerPreCheckoutQuery answers a pre-checkout query.
// Returns true on success.
// See https://core.telegram.org/bots/api#answerprecheckoutquery
func (api *API) AnswerPreCheckoutQuery(params AnswerPreCheckoutQueryP) (bool, error) {
	req := NewRequest[bool]("answerPreCheckoutQuery", params)
	return req.Do(api)
}

// AnswerPreCheckoutQueryWithContext is the context-aware variant of AnswerPreCheckoutQuery.
// It executes the same request but uses ctx for cancellation and deadlines.
// See https://core.telegram.org/bots/api#answerprecheckoutquery
func (api *API) AnswerPreCheckoutQueryWithContext(ctx context.Context, params AnswerPreCheckoutQueryP) (bool, error) {
	req := NewRequest[bool]("answerPreCheckoutQuery", params)
	return req.DoWithContext(ctx, api)
}
