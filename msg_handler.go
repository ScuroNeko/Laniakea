package laniakea

import (
	"strings"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

func (bot *Bot[T]) handleMessage(update *tgapi.Update, ctx *MsgContext) {
	var msg *tgapi.Message
	if update.Message != nil {
		msg = update.Message
	} else if update.ChannelPost != nil {
		msg = update.ChannelPost
	} else {
		return
	}

	var text string
	if len(msg.Text) > 0 {
		text = msg.Text
	} else if len(msg.Caption) > 0 {
		text = msg.Caption
	} else {
		return
	}

	prefix, cmd, args := bot.parseCommand(text)
	if cmd == "" {
		return
	}
	ctx.Prefix = prefix

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

			if plugin.logger != nil {
				ctx.Logger = plugin.logger
			}
			if !plugin.executeMiddlewares(ctx, bot.appData) {
				return
			}
			plugin.executeCmd(cmd, ctx, bot.appData)
			return
		}
	}
}

func (bot *Bot[T]) handleCallback(update *tgapi.Update, ctx *MsgContext) {
	data, err := bot.decodePayload(update.CallbackQuery.Data)
	if err != nil {
		bot.logger.Errorln(err)
		return
	}

	ctx.Args = data.Args

	for _, plugin := range bot.plugins {
		_, ok := plugin.payloads[data.Command]
		if !ok {
			continue
		}

		ctx.Logger = plugin.logger
		if ctx.Logger == nil {
			ctx.Logger = bot.logger
		}
		if !plugin.executeMiddlewares(ctx, bot.appData) {
			return
		}
		plugin.executePayload(data.Command, ctx, bot.appData)
		return
	}
}

func (bot *Bot[T]) checkPrefixes(text string) (string, bool) {
	for _, prefix := range bot.prefixes {
		if prefix == "" {
			if bot.logger != nil {
				bot.logger.Warnln("empty prefix is not allowed")
			}
			continue
		}
		if strings.HasPrefix(text, prefix) {
			return prefix, true
		}
	}
	return "", false
}

func (bot *Bot[T]) parseCommand(text string) (prefix, cmd, args string) {
	if prefix, hasPrefix := bot.checkPrefixes(text); hasPrefix {
		text = strings.TrimSpace(text[len(prefix):])
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
		return prefix, cmd, args
	}
	return "", "", ""
}
