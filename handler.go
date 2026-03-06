package laniakea

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

func (bot *Bot[T]) handle(u *tgapi.Update) {
	ctx := &MsgContext{
		Update: *u, Api: bot.api,
		botLogger:     bot.logger,
		errorTemplate: bot.errorTemplate,
		l10n:          bot.l10n,
		draftProvider: bot.draftProvider,
	}
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

	// Убираем префикс
	text = strings.TrimSpace(text[len(prefix):])

	// Извлекаем команду как первое слово
	spaceIndex := strings.Index(text, " ")
	var cmd string
	var args string

	if spaceIndex == -1 {
		cmd = text
		args = ""
	} else {
		cmd = text[:spaceIndex]
		args = strings.TrimSpace(text[spaceIndex:])
	}

	if strings.Contains(cmd, "@") {
		botUsername := bot.username
		if botUsername != "" && strings.HasSuffix(cmd, "@"+botUsername) {
			cmd = cmd[:len(cmd)-len("@"+botUsername)] // убираем @botname
		}
	}

	// Ищем команду по точному совпадению
	for _, plugin := range bot.plugins {
		if _, exists := plugin.commands[cmd]; exists {
			ctx.Text = args
			ctx.Args = strings.Fields(args) // Убирает лишние пробелы
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
