package laniakea

import (
	"encoding/json"
	"fmt"
)

type InlineKeyboard struct {
	CurrentLine []InlineKeyboardButton
	Lines       [][]InlineKeyboardButton
	maxRow      int
}

func NewInlineKeyboard(maxRow int) *InlineKeyboard {
	return &InlineKeyboard{
		CurrentLine: make([]InlineKeyboardButton, 0),
		Lines:       make([][]InlineKeyboardButton, 0),
		maxRow:      maxRow,
	}
}

func (in *InlineKeyboard) append(button InlineKeyboardButton) *InlineKeyboard {
	if len(in.CurrentLine) == in.maxRow {
		in.AddLine()
	}
	in.CurrentLine = append(in.CurrentLine, button)
	return in
}
func (in *InlineKeyboard) AddUrlButton(text, url string) *InlineKeyboard {
	return in.append(InlineKeyboardButton{Text: text, URL: url})
}
func (in *InlineKeyboard) AddCallbackButton(text string, cmd string, args ...any) *InlineKeyboard {
	return in.append(InlineKeyboardButton{Text: text, CallbackData: NewCallbackData(cmd, args...).ToJson()})
}

func (in *InlineKeyboard) AddLine() *InlineKeyboard {
	if len(in.CurrentLine) == 0 {
		return in
	}
	in.Lines = append(in.Lines, in.CurrentLine)
	in.CurrentLine = make([]InlineKeyboardButton, 0)
	return in
}
func (in *InlineKeyboard) Get() InlineKeyboardMarkup {
	if len(in.CurrentLine) > 0 {
		in.Lines = append(in.Lines, in.CurrentLine)
	}
	return InlineKeyboardMarkup{
		InlineKeyboard: in.Lines,
	}
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
