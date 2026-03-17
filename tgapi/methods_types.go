package tgapi

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

// WebhookInfo describes the current webhook status.
// See https://core.telegram.org/bots/api#webhookinfo
type WebhookInfo struct {
	URL                          string   `json:"url"`
	HasCustomCertificate         bool     `json:"has_custom_certificate"`
	PendingUpdateCount           int      `json:"pending_update_count"`
	IPAddress                    string   `json:"ip_address,omitempty"`
	LastErrorDate                int      `json:"last_error_date,omitempty"`
	LastErrorMessage             string   `json:"last_error_message,omitempty"`
	LastSynchronizationErrorDate int      `json:"last_synchronization_error_date,omitempty"`
	MaxConnections               int      `json:"max_connections,omitempty"`
	AllowedUpdates               []string `json:"allowed_updates,omitempty"`
}
