package tgapi

// InlineQueryResult is a JSON-serializable inline query result object.
// Since: Bot API 1.7
// See https://core.telegram.org/bots/api#inlinequeryresult
type InlineQueryResult map[string]any

// InlineQueryResultsButton represents a button shown above inline query results.
// Since: Bot API 6.3
// See https://core.telegram.org/bots/api#inlinequeryresultsbutton
type InlineQueryResultsButton struct {
	Text           string      `json:"text"`
	WebApp         *WebAppInfo `json:"web_app,omitempty"`
	StartParameter string      `json:"start_parameter,omitempty"`
}

// SentWebAppMessage describes an inline message sent by a Web App on behalf of a user.
// Since: Bot API 8.0
// See https://core.telegram.org/bots/api#sentwebappmessage
type SentWebAppMessage struct {
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// PreparedInlineMessage describes a prepared inline message.
// Since: Bot API 8.0
// See https://core.telegram.org/bots/api#preparedinlinemessage
type PreparedInlineMessage struct {
	ID             string `json:"id"`
	ExpirationDate int    `json:"expiration_date"`
}

// PreparedKeyboardButton describes a prepared keyboard button.
// Since: Bot API 8.0
// See https://core.telegram.org/bots/api#preparedkeyboardbutton
type PreparedKeyboardButton struct {
	ID string `json:"id"`
}
