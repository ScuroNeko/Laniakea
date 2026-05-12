package laniakea

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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
	startTime := time.Now()

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	msgCtx := &MsgContext{
		Update: *u, API: bot.api,
		Logger:        bot.logger,
		errorTemplate: bot.errorTemplate,
		l10n:          bot.l10n,
		draftProvider: bot.draftProvider,
		sceneRuntime:  bot,
		observer:      bot.observer,
		payloadType:   bot.payloadType,
		ctx:           ctx,
	}
	bot.prepareUpdateCtx(u, msgCtx)
	bot.safeEmitEvent(ctx, UpdateReceivedEvent{
		UpdateID:   u.UpdateID,
		UpdateType: u.Type,
		FromID:     msgCtx.FromID,
		ChatID:     msgCtx.ChatID,
	})

	for _, middleware := range bot.middlewares {
		if !middleware.Execute(msgCtx, bot.appData) {
			bot.safeEmitEvent(ctx, UpdateHandledEvent{
				UpdateID:   u.UpdateID,
				UpdateType: u.Type,
				FromID:     msgCtx.FromID,
				ChatID:     msgCtx.ChatID,
				Duration:   time.Since(startTime),
				Handled:    false,
			})
			return
		}
	}

	sceneHandled, err := bot.tryHandleScene(msgCtx)
	if err != nil {
		bot.logger.Errorln(err)
		bot.safeEmitEvent(ctx, UpdateHandledEvent{
			UpdateID:   u.UpdateID,
			UpdateType: u.Type,
			FromID:     msgCtx.FromID,
			ChatID:     msgCtx.ChatID,
			Duration:   time.Since(startTime),
			Handled:    false,
		})
		bot.safeEmitEvent(ctx, ErrorEvent{
			UpdateID:    u.UpdateID,
			UpdateType:  u.Type,
			Plugin:      "bot",
			HandlerKind: HandlerSceneKind,
			HandlerName: "tryHandleScene",
			FromID:      msgCtx.FromID,
			ChatID:      msgCtx.ChatID,
			Err:         err,
			UserFacing:  false,
		})
		return
	}
	if sceneHandled {
		bot.safeEmitEvent(ctx, UpdateHandledEvent{
			UpdateID:   u.UpdateID,
			UpdateType: u.Type,
			FromID:     msgCtx.FromID,
			ChatID:     msgCtx.ChatID,
			Duration:   time.Since(startTime),
			Handled:    true,
		})
		return
	}

	handled := false
	switch u.Type {
	case tgapi.UpdateTypeMessage, tgapi.UpdateTypeChannelPost:
		handled = bot.handleMessage(u, msgCtx)
	case tgapi.UpdateTypeCallbackQuery:
		handled = bot.handleCallback(u, msgCtx)
	default:
		handled = bot.handleUpdate(u, msgCtx)
	}
	bot.safeEmitEvent(ctx, UpdateHandledEvent{
		UpdateID:   u.UpdateID,
		UpdateType: u.Type,
		FromID:     msgCtx.FromID,
		ChatID:     msgCtx.ChatID,
		Duration:   time.Since(startTime),
		Handled:    handled,
	})
}

func cloneMsgContext(src *MsgContext) *MsgContext {
	cloned := *src
	if src.Args != nil {
		cloned.Args = append([]string(nil), src.Args...)
	}
	return &cloned
}

func encodeJSONPayload(d CallbackData) (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func decodeJSONPayload(s string) (CallbackData, error) {
	var data CallbackData
	err := json.Unmarshal([]byte(s), &data)
	return data, err
}

func encodeBase64Payload(d CallbackData) (string, error) {
	data, err := encodeJSONPayload(d)
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
	return decodeJSONPayload(string(b))
}

func encodeCompactPayload(d CallbackData) (string, error) {
	args := strings.Join(d.Args, ",")
	return d.Command + "|" + args, nil
}
func decodeCompactPayload(s string) (CallbackData, error) {
	values := strings.SplitN(s, "|", 2)
	if len(values) != 2 {
		return CallbackData{}, errors.New("invalid payload")
	}
	cmd, argsRaw := values[0], values[1]
	var args []string
	if argsRaw != "" {
		args = strings.Split(argsRaw, ",")
	}
	return CallbackData{Command: cmd, Args: args}, nil
}
func encodeCompactBase64Payload(d CallbackData) (string, error) {
	payload, _ := encodeCompactPayload(d)
	return base64.RawURLEncoding.EncodeToString([]byte(payload)), nil
}
func decodeCompactBase64Payload(s string) (CallbackData, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return CallbackData{}, err
	}
	return decodeCompactPayload(string(b))
}

func decodePayloadAs(payloadType BotPayloadType, s string) (CallbackData, error) {
	switch payloadType {
	case BotPayloadBase64:
		return decodeBase64Payload(s)
	case BotPayloadJSON:
		return decodeJSONPayload(s)
	case BotPayloadCompact:
		return decodeCompactPayload(s)
	case BotPayloadCompactBase64:
		return decodeCompactBase64Payload(s)
	}
	return CallbackData{}, ErrInvalidPayloadType
}

func decodePayload(payloadType BotPayloadType, s string, strict bool) (CallbackData, BotPayloadType, error) {
	knownTypes := []BotPayloadType{
		BotPayloadBase64,
		BotPayloadJSON,
		BotPayloadCompact,
		BotPayloadCompactBase64,
	}
	if _, err := decodePayloadAs(payloadType, ""); errors.Is(err, ErrInvalidPayloadType) {
		return CallbackData{}, "", ErrInvalidPayloadType
	}

	data, err := decodePayloadAs(payloadType, s)
	if err == nil {
		return data, payloadType, nil
	}
	if strict {
		return CallbackData{}, "", fmt.Errorf("%w: expected %s", ErrPayloadTypeMismatch, payloadType)
	}

	for _, candidate := range knownTypes {
		if candidate == payloadType {
			continue
		}
		data, err = decodePayloadAs(candidate, s)
		if err == nil {
			return data, candidate, nil
		}
	}
	return CallbackData{}, "", err
}

func (bot *Bot[T]) decodePayload(s string) (CallbackData, error) {
	data, decodedType, err := decodePayload(bot.payloadType, s, bot.strictPayloadType)
	if err != nil {
		return CallbackData{}, err
	}
	if decodedType == BotPayloadBase64 && bot.debug && bot.logger != nil {
		bot.logger.Debugf("decoded callback payload base64->json: raw=%q json=%s", s, data.ToJSON())
	}
	return data, nil
}
