package laniakea

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

func (bot *Bot[T]) handle(u *tgapi.Update) {
	ctx := &MsgContext{Update: *u, Api: bot.api, botLogger: bot.logger, errorTemplate: bot.errorTemplate, l10n: bot.l10n}
	for _, middleware := range bot.middlewares {
		middleware.Execute(ctx, bot.dbContext)
	}

	if u.CallbackQuery != nil {
		bot.handleCallback(u, ctx)
	} else {
		bot.handleMessage(u, ctx)
	}
}

func (bot *Bot[T]) handleMessage(update *tgapi.Update, ctx *MsgContext) {
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
	prefix, hasPrefix := bot.checkPrefixes(text)
	if !hasPrefix {
		return
	}
	ctx.Prefix = prefix
	ctx.FromID = update.Message.From.ID
	ctx.From = update.Message.From
	ctx.Msg = update.Message

	text = strings.TrimSpace(text[len(prefix):])

	for _, plugin := range bot.plugins {
		for cmd := range plugin.commands {
			if !strings.HasPrefix(text, cmd) {
				continue
			}
			requestParts := strings.Split(text, " ")
			cmdParts := strings.Split(cmd, " ")
			isValid := true
			for i, part := range cmdParts {
				if part != requestParts[i] {
					isValid = false
					break
				}
			}

			if !isValid {
				continue
			}

			ctx.Text = strings.TrimSpace(text[len(cmd):])
			if ctx.Text == "" {
				ctx.Args = []string{}
			} else {
				ctx.Args = strings.Split(ctx.Text, " ")
			}

			if !plugin.executeMiddlewares(ctx, bot.dbContext) {
				return
			}
			go plugin.executeCmd(cmd, ctx, bot.dbContext)
			return
		}
	}
}

func (bot *Bot[T]) handleCallback(update *tgapi.Update, ctx *MsgContext) {
	data := new(CallbackData)
	err := json.Unmarshal([]byte(update.CallbackQuery.Data), data)
	if err != nil {
		bot.logger.Errorln(err)
		return
	}

	ctx.FromID = update.CallbackQuery.From.ID
	ctx.From = &update.CallbackQuery.From
	ctx.Msg = &update.CallbackQuery.Message
	ctx.CallbackMsgId = update.CallbackQuery.Message.MessageID
	ctx.CallbackQueryId = update.CallbackQuery.ID
	ctx.Args = data.Args

	for _, plugin := range bot.plugins {
		_, ok := plugin.payloads[data.Command]
		if !ok {
			continue
		}

		if !plugin.executeMiddlewares(ctx, bot.dbContext) {
			return
		}
		go plugin.executePayload(data.Command, ctx, bot.dbContext)
		return
	}
}

func (bot *Bot[T]) checkPrefixes(text string) (string, bool) {
	for _, prefix := range bot.prefixes {
		if strings.HasPrefix(text, prefix) {
			return prefix, true
		}
	}
	return "", false
}

func encodeJsonPayload(d CallbackData) (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func decodeJsonPayload(s string) (CallbackData, error) {
	var data CallbackData
	err := json.Unmarshal([]byte(s), &data)
	return data, err
}
func encodeBase64Payload(d CallbackData) (string, error) {
	data, err := encodeJsonPayload(d)
	if err != nil {
		return "", err
	}
	dst := make([]byte, base64.StdEncoding.EncodedLen(len([]byte(data))))
	base64.StdEncoding.Encode(dst, []byte(data))
	return string(dst), nil
}
func decodeBase64Payload(s string) (CallbackData, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return CallbackData{}, err
	}
	return decodeJsonPayload(string(b))
}
