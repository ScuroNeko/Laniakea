package laniakea

import (
	"strings"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

func (bot *Bot[T]) handleMessage(update *tgapi.Update, ctx *MessageContext) bool {
	text, ok := messageText(update)
	if !ok {
		return false
	}

	prefix, cmd, args := bot.parseCommand(text)
	if cmd == "" {
		return bot.handleFallback(update, ctx)
	}
	ctx.Prefix = prefix

	if strings.Contains(cmd, "@") {
		botUsername := bot.username
		if botUsername != "" && strings.HasSuffix(cmd, "@"+botUsername) {
			cmd = cmd[:len(cmd)-len("@"+botUsername)] // remove @botname
		}
	}
	// Ищем команду по точному совпадению
	for _, plugin := range bot.plugins {
		if _, exists := plugin.commands[cmd]; exists {

			ctx.Text = args
			ctx.Args = strings.Fields(args)

			if plugin.logger != nil {
				ctx.Logger = plugin.logger
			}
			if !plugin.executeMiddlewares(ctx, bot.appData) {
				return false
			}

			startTime := time.Now()
			bot.safeEmitEvent(ctx.Context(), HandlerStartedEvent{
				UpdateID:    update.UpdateID,
				UpdateType:  update.Type,
				Plugin:      plugin.name,
				HandlerKind: HandlerCommandKind,
				HandlerName: cmd,
				FromID:      ctx.FromID,
				ChatID:      ctx.ChatID,
			})

			err := plugin.executeCmd(cmd, ctx, bot.appData)
			handlerEndEvent := HandlerFinishedEvent{
				UpdateID:    update.UpdateID,
				UpdateType:  update.Type,
				Plugin:      plugin.name,
				HandlerKind: HandlerCommandKind,
				HandlerName: cmd,
				FromID:      ctx.FromID,
				ChatID:      ctx.ChatID,
				Duration:    time.Since(startTime),
			}

			var errorEvent *ErrorEvent = nil
			if err != nil {
				ctx.error(err)
				handlerEndEvent.Err = err
				handlerEndEvent.UserFacing = IsUserError(err)
				errorEvent = &ErrorEvent{
					UpdateID:    update.UpdateID,
					UpdateType:  update.Type,
					Plugin:      plugin.name,
					HandlerKind: HandlerCommandKind,
					HandlerName: cmd,
					FromID:      ctx.FromID,
					ChatID:      ctx.ChatID,
					Err:         err,
					UserFacing:  handlerEndEvent.UserFacing,
				}
			}
			bot.safeEmitEvent(ctx.Context(), handlerEndEvent)
			if errorEvent != nil {
				bot.safeEmitEvent(ctx.Context(), *errorEvent)
			}
			return true
		}
	}

	return bot.handleFallback(update, ctx)
}

func (bot *Bot[T]) handleFallback(update *tgapi.Update, ctx *MessageContext) bool {
	text, ok := messageText(update)
	if !ok {
		return false
	}

	prefix, _, _ := bot.parseCommand(text)
	handled := false
	for _, plugin := range bot.plugins {
		if plugin.messageFallback == nil {
			continue
		}

		pluginCtx := cloneMsgContext(ctx)
		pluginCtx.Prefix = prefix
		pluginCtx.Text = text
		pluginCtx.Args = strings.Fields(text)
		if plugin.logger != nil {
			pluginCtx.Logger = plugin.logger
		}
		if !plugin.executeMiddlewares(pluginCtx, bot.appData) {
			continue
		}

		startTime := time.Now()
		bot.safeEmitEvent(pluginCtx.Context(), HandlerStartedEvent{
			UpdateID:    update.UpdateID,
			UpdateType:  update.Type,
			Plugin:      plugin.name,
			HandlerKind: HandlerMessageKind,
			HandlerName: "message_fallback",
			FromID:      pluginCtx.FromID,
			ChatID:      pluginCtx.ChatID,
		})
		err := plugin.messageFallback(pluginCtx, bot.appData)
		endEvent := HandlerFinishedEvent{
			UpdateID:    update.UpdateID,
			UpdateType:  update.Type,
			Plugin:      plugin.name,
			HandlerKind: HandlerMessageKind,
			HandlerName: "message_fallback",
			FromID:      pluginCtx.FromID,
			ChatID:      pluginCtx.ChatID,
			Duration:    time.Since(startTime),
		}
		if err != nil {
			endEvent.Err = err
			endEvent.UserFacing = IsUserError(err)
		}
		bot.safeEmitEvent(pluginCtx.Context(), endEvent)
		if err != nil {
			pluginCtx.error(err)
			bot.safeEmitEvent(pluginCtx.Context(), ErrorEvent{
				UpdateID:    update.UpdateID,
				UpdateType:  update.Type,
				Plugin:      plugin.name,
				HandlerKind: HandlerMessageKind,
				HandlerName: "message_fallback",
				FromID:      pluginCtx.FromID,
				ChatID:      pluginCtx.ChatID,
				Err:         err,
				UserFacing:  IsUserError(err),
			})
		}
		handled = true
	}
	return handled
}

func messageText(update *tgapi.Update) (string, bool) {
	var msg *tgapi.Message
	if update.Message != nil {
		msg = update.Message
	} else if update.ChannelPost != nil {
		msg = update.ChannelPost
	} else {
		return "", false
	}

	var text string
	if len(msg.Text) > 0 {
		text = msg.Text
	} else if len(msg.Caption) > 0 {
		text = msg.Caption
	} else {
		return "", false
	}
	return text, true
}

func (bot *Bot[T]) handleCallback(update *tgapi.Update, ctx *MessageContext) bool {
	data, err := bot.decodePayload(update.CallbackQuery.Data)
	if err != nil {
		bot.logger.Errorln(err)
		bot.safeEmitEvent(ctx.Context(), ErrorEvent{
			UpdateID:    update.UpdateID,
			UpdateType:  update.Type,
			Plugin:      "bot",
			HandlerKind: HandlerPayloadKind,
			HandlerName: "decodePayload",
			FromID:      ctx.FromID,
			ChatID:      ctx.ChatID,
			Err:         err,
			UserFacing:  false,
		})
		return false
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
			return false
		}

		startTime := time.Now()
		bot.safeEmitEvent(ctx.Context(), HandlerStartedEvent{
			UpdateID:    update.UpdateID,
			UpdateType:  update.Type,
			Plugin:      plugin.name,
			HandlerKind: HandlerPayloadKind,
			HandlerName: data.Command,
			FromID:      ctx.FromID,
			ChatID:      ctx.ChatID,
		})
		err := plugin.executePayload(data.Command, ctx, bot.appData)

		endEvent := HandlerFinishedEvent{
			UpdateID:    update.UpdateID,
			UpdateType:  update.Type,
			Plugin:      plugin.name,
			HandlerKind: HandlerPayloadKind,
			HandlerName: data.Command,
			FromID:      ctx.FromID,
			ChatID:      ctx.ChatID,
			Duration:    time.Since(startTime),
		}
		var errorEvent *ErrorEvent = nil
		if err != nil {
			ctx.error(err)
			errorEvent = &ErrorEvent{
				UpdateID:    update.UpdateID,
				UpdateType:  update.Type,
				Plugin:      plugin.name,
				HandlerKind: HandlerPayloadKind,
				HandlerName: data.Command,
				FromID:      ctx.FromID,
				ChatID:      ctx.ChatID,
				Err:         err,
				UserFacing:  IsUserError(err),
			}
			endEvent.Err = err
			endEvent.UserFacing = errorEvent.UserFacing
		}
		bot.safeEmitEvent(ctx.Context(), endEvent)
		if errorEvent != nil {
			bot.safeEmitEvent(ctx.Context(), *errorEvent)
		}
		return true
	}
	return false
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
