package laniakea

import "encoding/json"

type InlineKeyboard struct {
	CurrentLine []InlineKeyboardButton
	Lines       [][]InlineKeyboardButton
}

func NewInlineKeyboard() *InlineKeyboard {
	return &InlineKeyboard{
		CurrentLine: make([]InlineKeyboardButton, 0),
		Lines:       make([][]InlineKeyboardButton, 0),
	}
}

func (in *InlineKeyboard) append(button InlineKeyboardButton) *InlineKeyboard {
	in.CurrentLine = append(in.CurrentLine, button)
	return in
}
func (in *InlineKeyboard) AddUrlButton(text, url string) *InlineKeyboard {
	return in.append(InlineKeyboardButton{Text: text, URL: url})
}
func (in *InlineKeyboard) AddCallbackButton(text string, data CallbackData) *InlineKeyboard {
	return in.append(InlineKeyboardButton{Text: text, CallbackData: data.ToJson()})
}

func (in *InlineKeyboard) AddLine() *InlineKeyboard {
	in.Lines = append(in.Lines, in.CurrentLine)
	in.CurrentLine = make([]InlineKeyboardButton, 0)
	return in
}
func (in *InlineKeyboard) Get() *InlineKeyboardMarkup {
	if len(in.CurrentLine) > 0 {
		in.Lines = append(in.Lines, in.CurrentLine)
	}
	return &InlineKeyboardMarkup{
		InlineKeyboard: in.Lines,
	}
}

type CallbackData struct {
	Command string `json:"cmd"`
	Args    []any  `json:"args"`
}

func (d *CallbackData) ToJson() string {
	data, err := json.Marshal(d)
	if err != nil {
		return `{"cmd":""}`
	}
	return string(data)
}
