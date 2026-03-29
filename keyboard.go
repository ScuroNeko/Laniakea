package laniakea

import (
	"fmt"

	"git.scuroneko.dev/scuroneko/extypes"
	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

const (
	// ButtonStyleDanger marks a destructive inline keyboard action.
	ButtonStyleDanger tgapi.KeyboardButtonStyle = "danger"
	// ButtonStyleSuccess marks a confirmatory inline keyboard action.
	ButtonStyleSuccess tgapi.KeyboardButtonStyle = "success"
	// ButtonStylePrimary marks a primary inline keyboard action.
	ButtonStylePrimary tgapi.KeyboardButtonStyle = "primary"
)

// InlineKbButtonBuilder is a fluent builder for creating a single inline keyboard button.
//
// Use NewInlineKbButton() to start, then chain methods to configure:
//   - SetIconCustomEmojiId() — adds a custom emoji icon
//   - SetStyle() — sets visual style (danger/success/primary)
//   - SetUrl() — makes button open a URL
//   - SetCallbackDataJson() — attaches structured command + args for bot handling
//
// Call build() to produce the final tgapi.InlineKeyboardButton.
// Builder methods are immutable — each returns a copy.
type InlineKbButtonBuilder struct {
	text              string
	iconCustomEmojiID string
	style             tgapi.KeyboardButtonStyle
	url               string
	callbackData      string
}

// NewInlineKbButton creates a new button builder with the given display text.
// The button will have no URL, no style, and no callback data by default.
func NewInlineKbButton(text string) InlineKbButtonBuilder {
	return InlineKbButtonBuilder{text: text}
}

// SetIconCustomEmojiId sets a custom emoji ID to display as the button's icon.
// This is a Telegram Bot API feature for custom emoji icons.
func (b InlineKbButtonBuilder) SetIconCustomEmojiId(id string) InlineKbButtonBuilder {
	b.iconCustomEmojiID = id
	return b
}

// SetStyle sets the visual style of the button.
// Valid values: ButtonStyleDanger, ButtonStyleSuccess, ButtonStylePrimary.
// If not set, the button uses the default style.
func (b InlineKbButtonBuilder) SetStyle(style tgapi.KeyboardButtonStyle) InlineKbButtonBuilder {
	b.style = style
	return b
}

// SetUrl sets a URL that will be opened when the button is pressed.
// If both URL and CallbackData are set, Telegram will prioritize URL.
func (b InlineKbButtonBuilder) SetUrl(url string) InlineKbButtonBuilder {
	b.url = url
	return b
}

// SetCallbackDataJson sets a structured callback payload that will be sent to the bot
// when the button is pressed. The command and arguments are serialized as JSON.
//
// Args are converted to strings using fmt.Sprint. Non-string types (e.g., int, bool)
// are safely serialized, but complex structs may not serialize usefully.
//
// Example: SetCallbackDataJson("delete_user", 123, "confirm") → {"cmd":"delete_user","args":["123","confirm"]}.
func (b InlineKbButtonBuilder) SetCallbackDataJson(cmd string, args ...any) InlineKbButtonBuilder {
	b.callbackData = NewCallbackData(cmd, args...).ToJson()
	return b
}

// SetCallbackDataBase64 sets a structured callback payload encoded as Base64.
// This can be useful when the JSON payload exceeds Telegram's callback data length limit.
// Args are converted to strings using fmt.Sprint.
func (b InlineKbButtonBuilder) SetCallbackDataBase64(cmd string, args ...any) InlineKbButtonBuilder {
	b.callbackData = NewCallbackData(cmd, args...).ToBase64()
	return b
}

// Internal helper that converts the builder state into a Telegram button.
func (b InlineKbButtonBuilder) build() tgapi.InlineKeyboardButton {
	return tgapi.InlineKeyboardButton{
		Text:              b.text,
		URL:               b.url,
		Style:             b.style,
		IconCustomEmojiID: b.iconCustomEmojiID,
		CallbackData:      b.callbackData,
	}
}

// InlineKeyboard is a stateful builder for constructing Telegram inline keyboard layouts.
//
// Buttons are added row-by-row. When a row reaches maxRow, it is automatically flushed.
// Call AddLine() to manually end a row, or Get() to finalize and retrieve the markup.
//
// The keyboard is not thread-safe. Build it in a single goroutine.
type InlineKeyboard struct {
	CurrentLine extypes.Slice[tgapi.InlineKeyboardButton] // Current row being built
	Lines       [][]tgapi.InlineKeyboardButton            // Completed rows
	maxRow      int                                       // Max buttons per row (e.g., 3 or 4)

	payloadType BotPayloadType // Serialization format for callback data (JSON or Base64)
}

// NewInlineKeyboardJson creates a new keyboard builder with the specified maximum
// number of buttons per row.
//
// Example: NewInlineKeyboardJson(3) creates a keyboard with at most 3 buttons per line.
func NewInlineKeyboardJson(maxRow int) *InlineKeyboard {
	return NewInlineKeyboard(BotPayloadJson, maxRow)
}

// NewInlineKeyboardBase64 creates a new keyboard builder with the specified maximum
// number of buttons per row, using Base64 encoding for button payloads.
//
// Example: NewInlineKeyboardBase64(3) creates a keyboard with at most 3 buttons per line.
func NewInlineKeyboardBase64(maxRow int) *InlineKeyboard {
	return NewInlineKeyboard(BotPayloadBase64, maxRow)
}

// NewInlineKeyboard creates a new keyboard builder with the specified payload encoding
// type and maximum number of buttons per row.
//
// Use NewInlineKeyboardJson or NewInlineKeyboardBase64 for the common cases.
func NewInlineKeyboard(payloadType BotPayloadType, maxRow int) *InlineKeyboard {
	return &InlineKeyboard{
		CurrentLine: make(extypes.Slice[tgapi.InlineKeyboardButton], 0),
		Lines:       make([][]tgapi.InlineKeyboardButton, 0),
		maxRow:      maxRow,
		payloadType: payloadType,
	}
}

// SetPayloadType sets the keyboard-local serialization format for callback data added via
// AddCallbackButton and AddCallbackButtonStyle methods.
// It overrides the bot's default payload type for this keyboard only.
func (in *InlineKeyboard) SetPayloadType(t BotPayloadType) *InlineKeyboard {
	in.payloadType = t
	return in
}

// GetPayloadType returns the keyboard-local callback payload encoding type.
func (in *InlineKeyboard) GetPayloadType() BotPayloadType { return in.payloadType }

func (in *InlineKeyboard) SetMaxRow(maxRow int) *InlineKeyboard {
	in.maxRow = maxRow
	return in
}

// Internal helper that appends a button and auto-flushes a full row.
func (in *InlineKeyboard) append(button tgapi.InlineKeyboardButton) *InlineKeyboard {
	if in.CurrentLine.Len() == in.maxRow {
		in.AddLine()
	}
	in.CurrentLine = in.CurrentLine.Push(button)
	return in
}

// AddUrlButton adds a button that opens a URL when pressed.
// No callback data is attached.
func (in *InlineKeyboard) AddUrlButton(text, url string) *InlineKeyboard {
	return in.append(tgapi.InlineKeyboardButton{Text: text, URL: url})
}

// AddUrlButtonStyle adds a button with a visual style that opens a URL.
// Style must be one of: ButtonStyleDanger, ButtonStyleSuccess, ButtonStylePrimary.
func (in *InlineKeyboard) AddUrlButtonStyle(text string, style tgapi.KeyboardButtonStyle, url string) *InlineKeyboard {
	return in.append(tgapi.InlineKeyboardButton{Text: text, Style: style, URL: url})
}

// AddCallbackButton adds a button that sends a structured callback payload to the bot.
// The command and args are serialized according to the current payloadType.
func (in *InlineKeyboard) AddCallbackButton(text string, cmd string, args ...any) *InlineKeyboard {
	return in.append(tgapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: NewCallbackData(cmd, args...).Encode(in.payloadType),
	})
}

// AddCallbackButtonStyle adds a styled callback button.
// Style affects visual appearance; callback data is sent to bot on press.
func (in *InlineKeyboard) AddCallbackButtonStyle(text string, style tgapi.KeyboardButtonStyle, cmd string, args ...any) *InlineKeyboard {
	return in.append(tgapi.InlineKeyboardButton{
		Text:         text,
		Style:        style,
		CallbackData: NewCallbackData(cmd, args...).Encode(in.payloadType),
	})
}

// AddButton adds a button pre-configured via InlineKbButtonBuilder.
// This is the most flexible way to create buttons with custom emoji, style, URL, and callback.
func (in *InlineKeyboard) AddButton(b InlineKbButtonBuilder) *InlineKeyboard {
	return in.append(b.build())
}

// AddLine manually ends the current row and starts a new one.
// If the current row is empty, nothing happens.
func (in *InlineKeyboard) AddLine() *InlineKeyboard {
	if in.CurrentLine.Len() == 0 {
		return in
	}
	in.Lines = append(in.Lines, in.CurrentLine)
	in.CurrentLine = make(extypes.Slice[tgapi.InlineKeyboardButton], 0)
	return in
}

// Get finalizes the keyboard and returns a tgapi.ReplyMarkup.
// Automatically flushes the current line if not empty.
//
// Returns a pointer to a ReplyMarkup suitable for use with tgapi.SendMessage.
func (in *InlineKeyboard) Get() *tgapi.ReplyMarkup {
	if in.CurrentLine.Len() > 0 {
		in.AddLine()
	}
	return &tgapi.ReplyMarkup{InlineKeyboard: in.Lines}
}

// CallbackData represents the structured payload sent when an inline button
// with callback data is pressed.
//
// This structure is serialized to JSON and sent to the bot as a string.
// The bot should parse this back to determine the command and arguments.
//
// Example:
//
//	{"cmd":"delete_user","args":["123","confirm"]}
type CallbackData struct {
	Command string   `json:"cmd"`  // The command name to route to
	Args    []string `json:"args"` // Arguments passed as strings
}

// NewCallbackData creates a new CallbackData instance with the given command and args.
//
// All args are converted to strings using fmt.Sprint. This is safe for primitives
// (int, string, bool, float64) but may not serialize complex structs meaningfully.
//
// Use this to build callback payloads for bot command routing.
func NewCallbackData(command string, args ...any) CallbackData {
	stringArgs := make([]string, len(args))
	for i, arg := range args {
		stringArgs[i] = fmt.Sprint(arg)
	}
	return CallbackData{
		Command: command,
		Args:    stringArgs,
	}
}

// ToJson serializes the CallbackData to a JSON string.
//
// If serialization fails (e.g., due to unmarshalable fields), returns a fallback
// JSON object: {"cmd":""} to prevent breaking Telegram's API.
//
// This fallback ensures the bot receives a valid JSON payload even if internal
// errors occur — avoiding "invalid callback_data" errors from Telegram.
func (d CallbackData) ToJson() string {
	data, err := encodeJsonPayload(d)
	if err != nil {
		// Fallback: return minimal valid JSON to avoid Telegram API rejection
		return `{"cmd":""}`
	}
	return data
}

// ToBase64 serializes the CallbackData to a JSON string and then encodes it as Base64.
// Returns an empty string if serialization or encoding fails.
func (d CallbackData) ToBase64() string {
	s, err := encodeBase64Payload(d)
	if err != nil {
		return ``
	}
	return s
}

// Encode serializes the CallbackData according to the specified payload type.
// Supported types: BotPayloadJson and BotPayloadBase64.
// For unknown types, returns an empty string.
func (d CallbackData) Encode(t BotPayloadType) string {
	switch t {
	case BotPayloadBase64:
		return d.ToBase64()
	case BotPayloadJson:
		return d.ToJson()
	}
	return ""
}
