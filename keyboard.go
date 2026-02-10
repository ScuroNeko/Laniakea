package laniakea

import (
	"encoding/json"
	"fmt"

	"git.nix13.pw/scuroneko/extypes"
)

const (
	ButtonStyleDanger  KeyboardButtonStyle = "danger"
	ButtonStyleSuccess KeyboardButtonStyle = "success"
	ButtonStylePrimary KeyboardButtonStyle = "primary"
)

type InlineKbButtonBuilder struct {
	text              string
	iconCustomEmojiID string
	style             KeyboardButtonStyle
	url               string
	callbackData      string
}

func NewInlineKbButton(text string) InlineKbButtonBuilder {
	return InlineKbButtonBuilder{text: text}
}
func (b InlineKbButtonBuilder) SetIconCustomEmojiId(id string) InlineKbButtonBuilder {
	b.iconCustomEmojiID = id
	return b
}
func (b InlineKbButtonBuilder) SetStyle(style KeyboardButtonStyle) InlineKbButtonBuilder {
	b.style = style
	return b
}
func (b InlineKbButtonBuilder) SetUrl(url string) InlineKbButtonBuilder {
	b.url = url
	return b
}
func (b InlineKbButtonBuilder) SetCallbackData(cmd string, args ...any) InlineKbButtonBuilder {
	b.callbackData = NewCallbackData(cmd, args...).ToJson()
	return b
}

func (b InlineKbButtonBuilder) build() InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:              b.text,
		URL:               b.url,
		Style:             b.style,
		IconCustomEmojiID: b.iconCustomEmojiID,
		CallbackData:      b.callbackData,
	}
}

type InlineKeyboard struct {
	CurrentLine extypes.Slice[InlineKeyboardButton]
	Lines       [][]InlineKeyboardButton
	maxRow      int
}

func NewInlineKeyboard(maxRow int) *InlineKeyboard {
	return &InlineKeyboard{
		CurrentLine: make(extypes.Slice[InlineKeyboardButton], 0),
		Lines:       make([][]InlineKeyboardButton, 0),
		maxRow:      maxRow,
	}
}

func (in *InlineKeyboard) append(button InlineKeyboardButton) *InlineKeyboard {
	if in.CurrentLine.Len() == in.maxRow {
		in.AddLine()
	}
	in.CurrentLine = in.CurrentLine.Push(button)
	return in
}

func (in *InlineKeyboard) AddUrlButton(text, url string) *InlineKeyboard {
	return in.append(InlineKeyboardButton{Text: text, URL: url})
}
func (in *InlineKeyboard) AddUrlButtonStyle(text string, style KeyboardButtonStyle, url string) *InlineKeyboard {
	return in.append(InlineKeyboardButton{Text: text, Style: style, URL: url})
}
func (in *InlineKeyboard) AddCallbackButton(text string, cmd string, args ...any) *InlineKeyboard {
	return in.append(InlineKeyboardButton{
		Text: text, CallbackData: NewCallbackData(cmd, args...).ToJson(),
	})
}
func (in *InlineKeyboard) AddCallbackButtonStyle(text string, style KeyboardButtonStyle, cmd string, args ...any) *InlineKeyboard {
	return in.append(InlineKeyboardButton{
		Text: text, Style: style,
		CallbackData: NewCallbackData(cmd, args...).ToJson(),
	})
}
func (in *InlineKeyboard) AddButton(b InlineKbButtonBuilder) *InlineKeyboard {
	return in.append(b.build())
}

func (in *InlineKeyboard) AddLine() *InlineKeyboard {
	if in.CurrentLine.Len() == 0 {
		return in
	}
	in.Lines = append(in.Lines, in.CurrentLine)
	in.CurrentLine = make(extypes.Slice[InlineKeyboardButton], 0)
	return in
}
func (in *InlineKeyboard) Get() *InlineKeyboardMarkup {
	if in.CurrentLine.Len() > 0 {
		in.Lines = append(in.Lines, in.CurrentLine)
	}
	return &InlineKeyboardMarkup{InlineKeyboard: in.Lines}
}

type CallbackData struct {
	Command string   `json:"cmd"`
	Args    []string `json:"args"`
}

func NewCallbackData(command string, args ...any) *CallbackData {
	stringArgs := make([]string, len(args))
	for i, arg := range args {
		stringArgs[i] = fmt.Sprint(arg)
	}
	return &CallbackData{
		Command: command,
		Args:    stringArgs,
	}
}
func (d *CallbackData) ToJson() string {
	data, err := json.Marshal(d)
	if err != nil {
		return `{"cmd":""}`
	}
	return string(data)
}
