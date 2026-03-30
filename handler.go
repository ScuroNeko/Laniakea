package laniakea

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

// ErrInvalidPayloadType is returned when callback payload encoding type is unknown.
var ErrInvalidPayloadType = errors.New("invalid payload type")

func (bot *Bot[T]) handle(parentCtx context.Context, u *tgapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			bot.logger.Errorln(fmt.Sprintf("panic in handle: %v", r))
		}
	}()

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	msgCtx := &MsgContext{
		Update: *u, Api: bot.api,
		Logger:        bot.logger,
		errorTemplate: bot.errorTemplate,
		l10n:          bot.l10n,
		draftProvider: bot.draftProvider,
		sceneRuntime:  bot,
		payloadType:   bot.payloadType,
		ctx:           ctx,
	}
	bot.prepareUpdateCtx(u, msgCtx)

	for _, middleware := range bot.middlewares {
		if !middleware.Execute(msgCtx, bot.appData) {
			return
		}
	}

	sceneHandled, err := bot.tryHandleScene(msgCtx)
	if err != nil {
		bot.logger.Errorln(err)
		return
	}
	if sceneHandled {
		return
	}

	switch u.Type {
	case tgapi.UpdateTypeMessage, tgapi.UpdateTypeChannelPost:
		bot.handleMessage(u, msgCtx)
	case tgapi.UpdateTypeCallbackQuery:
		bot.handleCallback(u, msgCtx)
	default:
		bot.handleUpdate(u, msgCtx)
	}
}

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

func (bot *Bot[T]) handleUpdate(u *tgapi.Update, ctx *MsgContext) {
	for _, plugin := range bot.plugins {
		handler, ok := plugin.handlers[u.Type]
		if !ok {
			continue
		}

		pluginCtx := cloneMsgContext(ctx)
		if plugin.logger != nil {
			pluginCtx.Logger = plugin.logger
		}
		if !plugin.executeMiddlewares(pluginCtx, bot.appData) {
			continue
		}
		if err := handler(pluginCtx, bot.appData); err != nil {
			pluginCtx.error(err)
		}
	}
}

func cloneMsgContext(src *MsgContext) *MsgContext {
	cloned := *src
	if src.Args != nil {
		cloned.Args = append([]string(nil), src.Args...)
	}
	return &cloned
}

func (bot *Bot[T]) prepareUpdateCtx(u *tgapi.Update, ctx *MsgContext) {
	var from *tgapi.User
	switch u.Type {
	case tgapi.UpdateTypeMessage:
		if u.Message != nil {
			ctx.Msg = u.Message
		}
	case tgapi.UpdateTypeEditedMessage:
		if u.EditedMessage != nil {
			ctx.Msg = u.EditedMessage
		}
	case tgapi.UpdateTypeChannelPost:
		if u.ChannelPost != nil {
			ctx.Msg = u.ChannelPost
		}
	case tgapi.UpdateTypeEditedChannelPost:
		if u.EditedChannelPost != nil {
			ctx.Msg = u.EditedChannelPost
		}
	case tgapi.UpdateTypeBusinessMessage:
		if u.BusinessMessage != nil {
			ctx.Msg = u.BusinessMessage
		}
	case tgapi.UpdateTypeEditedBusinessMessage:
		if u.EditedBusinessMessage != nil {
			ctx.Msg = u.EditedBusinessMessage
		}
	case tgapi.UpdateTypeInlineQuery:
		if u.InlineQuery != nil {
			from = &u.InlineQuery.From
		}
	case tgapi.UpdateTypeChosenInlineResult:
		if u.ChosenInlineResult != nil {
			from = &u.ChosenInlineResult.From
		}
	case tgapi.UpdateTypeCallbackQuery:
		if u.CallbackQuery != nil {
			if u.CallbackQuery.Message != nil {
				ctx.Msg = u.CallbackQuery.Message
				ctx.CallbackMsgId = u.CallbackQuery.Message.MessageID
			}
			if u.CallbackQuery.InlineMessageID != nil {
				ctx.InlineMsgId = *u.CallbackQuery.InlineMessageID
			}
			ctx.CallbackQueryId = u.CallbackQuery.ID
			from = &u.CallbackQuery.From
		}
	case tgapi.UpdateTypeShippingQuery:
		if u.ShippingQuery != nil {
			from = &u.ShippingQuery.From
		}
	case tgapi.UpdateTypePreCheckoutQuery:
		if u.PreCheckoutQuery != nil {
			from = &u.PreCheckoutQuery.From
		}
	case tgapi.UpdateTypePurchasedPaidMedia:
		if u.PurchasedPaidMedia != nil {
			from = &u.PurchasedPaidMedia.From
		}
	case tgapi.UpdateTypeMyChatMember:
		if u.MyChatMember != nil {
			from = &u.MyChatMember.From
		}
	case tgapi.UpdateTypeChatMember:
		if u.ChatMember != nil {
			from = &u.ChatMember.From
		}
	case tgapi.UpdateTypeChatJoinRequest:
		if u.ChatJoinRequest != nil {
			from = &u.ChatJoinRequest.From
		}
	case tgapi.UpdateTypeBusinessConnection:
		if u.BusinessConnection != nil {
			from = &u.BusinessConnection.User
		}
	case tgapi.UpdateTypePollAnswer:
		if u.PollAnswer != nil {
			from = &u.PollAnswer.User
		}
	case tgapi.UpdateTypeMessageReaction:
		if u.MessageReaction != nil {
			from = u.MessageReaction.User
		}
	case tgapi.UpdateTypeChatBoost:
		if u.ChatBoost != nil {
			from = &u.ChatBoost.Boost.Source.User
		}
	case tgapi.UpdateTypeRemovedChatBoost:
		if u.RemovedChatBoost != nil {
			from = &u.RemovedChatBoost.Source.User
		}
	}
	if ctx.Msg != nil && from == nil {
		from = ctx.Msg.From
	}
	if from != nil {
		ctx.From = from
		ctx.FromID = from.ID
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
	dst := make([]byte, base64.RawURLEncoding.EncodedLen(len([]byte(data))))
	base64.RawURLEncoding.Encode(dst, []byte(data))
	return string(dst), nil
}

//	func encodePayload(payloadType BotPayloadType, d CallbackData) (string, error) {
//		switch payloadType {
//		case BotPayloadBase64:
//			return encodeBase64Payload(d)
//		case BotPayloadJson:
//			return encodeJsonPayload(d)
//		}
//		return "", ErrInvalidPayloadType
//	}
func decodeBase64Payload(s string) (CallbackData, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return CallbackData{}, err
	}
	return decodeJsonPayload(string(b))
}
func decodePayload(payloadType BotPayloadType, s string, strict bool) (CallbackData, BotPayloadType, error) {
	switch payloadType {
	case BotPayloadBase64:
		data, err := decodeBase64Payload(s)
		if err == nil {
			return data, BotPayloadBase64, nil
		}
		if strict {
			return CallbackData{}, "", fmt.Errorf("%w: expected %s", ErrPayloadTypeMismatch, BotPayloadBase64)
		}
		data, err = decodeJsonPayload(s)
		if err != nil {
			return CallbackData{}, "", err
		}
		return data, BotPayloadJson, nil
	case BotPayloadJson:
		data, err := decodeJsonPayload(s)
		if err == nil {
			return data, BotPayloadJson, nil
		}
		if strict {
			return CallbackData{}, "", fmt.Errorf("%w: expected %s", ErrPayloadTypeMismatch, BotPayloadJson)
		}
		data, err = decodeBase64Payload(s)
		if err != nil {
			return CallbackData{}, "", err
		}
		return data, BotPayloadBase64, nil
	}
	return CallbackData{}, "", ErrInvalidPayloadType
}

//	func (bot *Bot[T]) encodePayload(d CallbackData) (string, error) {
//		return encodePayload(bot.payloadType, d)
//	}
func (bot *Bot[T]) decodePayload(s string) (CallbackData, error) {
	data, decodedType, err := decodePayload(bot.payloadType, s, bot.strictPayloadType)
	if err != nil {
		return CallbackData{}, err
	}
	if decodedType == BotPayloadBase64 && bot.debug && bot.logger != nil {
		bot.logger.Debugf("decoded callback payload base64->json: raw=%q json=%s", s, data.ToJson())
	}
	return data, nil
}
