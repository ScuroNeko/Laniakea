package laniakea

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

// ErrInvalidPayloadType is returned when callback payload encoding type is unknown.
var ErrInvalidPayloadType = errors.New("invalid payload type")

func (bot *Bot[T]) handle(u *tgapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			bot.logger.Errorln(fmt.Sprintf("panic in handle: %v", r))
		}
	}()

	ctx := &MsgContext{
		Update: *u, Api: bot.api,
		Logger:        bot.logger,
		errorTemplate: bot.errorTemplate,
		l10n:          bot.l10n,
		draftProvider: bot.draftProvider,
		payloadType:   bot.payloadType,
	}
	bot.prepareUpdateCtx(u, ctx)

	for _, middleware := range bot.middlewares {
		if !middleware.Execute(ctx, bot.dbContext) {
			return
		}
	}

	switch u.Type {
	case tgapi.UpdateTypeMessage, tgapi.UpdateTypeChannelPost:
		bot.handleMessage(u, ctx)
	case tgapi.UpdateTypeCallbackQuery:
		bot.handleCallback(u, ctx)
	default:
		bot.handleUpdate(u, ctx)
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

	text = strings.TrimSpace(text)
	prefix, hasPrefix := bot.checkPrefixes(text)
	if !hasPrefix {
		return
	}

	ctx.Prefix = prefix
	ctx.Update = *update

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

			if plugin.logger != nil {
				ctx.Logger = plugin.logger
			}
			if !plugin.executeMiddlewares(ctx, bot.dbContext) {
				return
			}
			plugin.executeCmd(cmd, ctx, bot.dbContext)
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
		if !plugin.executeMiddlewares(ctx, bot.dbContext) {
			return
		}
		plugin.executePayload(data.Command, ctx, bot.dbContext)
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
		if !plugin.executeMiddlewares(pluginCtx, bot.dbContext) {
			continue
		}
		handler(pluginCtx, bot.dbContext)
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
	ctx.Update = *u
}

func (bot *Bot[T]) checkPrefixes(text string) (string, bool) {
	for _, prefix := range bot.prefixes {
		if prefix == "" {
			continue
		}
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
func decodePayload(payloadType BotPayloadType, s string) (CallbackData, error) {
	switch payloadType {
	case BotPayloadBase64:
		return decodeBase64Payload(s)
	case BotPayloadJson:
		return decodeJsonPayload(s)
	}
	return CallbackData{}, ErrInvalidPayloadType
}

//	func (bot *Bot[T]) encodePayload(d CallbackData) (string, error) {
//		return encodePayload(bot.payloadType, d)
//	}
func (bot *Bot[T]) decodePayload(s string) (CallbackData, error) {
	return decodePayload(bot.payloadType, s)
}
