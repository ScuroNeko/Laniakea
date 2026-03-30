package laniakea

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

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

func cloneMsgContext(src *MsgContext) *MsgContext {
	cloned := *src
	if src.Args != nil {
		cloned.Args = append([]string(nil), src.Args...)
	}
	return &cloned
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
