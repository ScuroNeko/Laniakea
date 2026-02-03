package laniakea

import (
	"encoding/json"
	"strings"
)

func (b *Bot) handleMessage(update *Update, ctx *MsgContext) {
	if update.Message == nil {
		return
	}

	var text string
	if len(update.Message.Text) > 0 {
		text = update.Message.Text
	} else {
		text = update.Message.Caption
	}

	text = strings.TrimSpace(text)
	prefix, hasPrefix := b.checkPrefixes(text)
	if !hasPrefix {
		return
	}
	ctx.Prefix = prefix
	ctx.FromID = update.Message.From.ID
	ctx.From = update.Message.From
	ctx.Msg = update.Message

	text = strings.TrimSpace(text[len(prefix):])

	for _, plugin := range b.plugins {
		// Check every command
		for cmd := range plugin.Commands {
			if !strings.HasPrefix(text, cmd) {
				continue
			}

			ctx.Text = strings.TrimSpace(text[len(cmd):])
			ctx.Args = strings.Split(ctx.Text, " ")

			go plugin.Execute(cmd, ctx, b.dbContext)
			return
		}
	}
}

func (b *Bot) handleCallback(update *Update, ctx *MsgContext) {
	data := new(CallbackData)
	err := json.Unmarshal([]byte(update.CallbackQuery.Data), data)
	if err != nil {
		b.logger.Errorln(err)
		return
	}

	ctx.FromID = update.CallbackQuery.From.ID
	ctx.From = update.CallbackQuery.From
	ctx.Msg = update.CallbackQuery.Message
	ctx.CallbackMsgId = update.CallbackQuery.Message.MessageID
	ctx.CallbackQueryId = update.CallbackQuery.ID
	ctx.Args = data.Args

	for _, plugin := range b.plugins {
		_, ok := plugin.Payloads[data.Command]
		if !ok {
			continue
		}
		go plugin.ExecutePayload(data.Command, ctx, b.dbContext)
		return
	}
}

func (b *Bot) checkPrefixes(text string) (string, bool) {
	for _, prefix := range b.prefixes {
		if strings.HasPrefix(text, prefix) {
			return prefix, true
		}
	}
	return "", false
}
