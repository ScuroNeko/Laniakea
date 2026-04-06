package laniakea

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/slog"
)

// MsgContext holds the normalized per-update context passed to command, payload,
// scene, middleware, and generic update handlers.
//
// MsgContext is populated from the current Telegram update before handler routing.
// Not every field is guaranteed for every update kind. In particular:
//   - Update is always present.
//   - Msg is populated only for update kinds that carry a Telegram message object.
//   - From and FromID are populated only when the update exposes a user identity.
//   - Chat and ChatID are populated only when the update exposes a chat identity.
//   - Text, Args, and Prefix are populated only by command or scene command routing.
//   - CallbackQueryId, CallbackMsgId, and InlineMsgId are populated only for
//     callback query handling when the corresponding callback targets exist.
//
// Helper methods on MsgContext may require a message-backed context. For example,
// reply helpers need Msg, while inline callback edit helpers can work through
// InlineMsgId when there is no chat message.
type MsgContext struct {
	Api    *tgapi.API
	Update tgapi.Update

	// Msg is the normalized Telegram message for message-backed update kinds.
	// It is nil for updates that do not include a message object.
	Msg *tgapi.Message
	// From is the normalized Telegram user for update kinds that expose one.
	// It stays nil for sender-chat-only updates and update kinds without a user.
	From *tgapi.User
	// Chat is the normalized Telegram chat for update kinds that expose one.
	// It is nil for updates that do not include a chat identity.
	Chat *tgapi.Chat

	// Logger is the logger assigned by the matched plugin for the current handler call.
	// It may fall back to the bot logger when the plugin has no dedicated logger.
	Logger *slog.Logger

	// InlineMsgId is the inline message identifier for callback queries that target
	// an inline message instead of a chat message.
	InlineMsgId string
	// CallbackMsgId is the message ID targeted by the current callback query when
	// the callback comes from a chat message.
	CallbackMsgId int
	// CallbackQueryId is the Telegram callback query ID for payload handlers and
	// callback-backed scene handlers.
	CallbackQueryId string
	// FromID is the normalized sender ID when the current update exposes a user.
	// It is zero when the update has no user identity.
	FromID int64
	// ChatID is the normalized chat ID when the current update exposes a chat.
	// It is zero when the update has no chat identity.
	ChatID int64
	// Prefix is the matched command prefix for command routing and scene-local
	// command routing. It is empty outside those flows.
	Prefix string
	// Text is the parsed command tail for command routing, the parsed scene-command
	// tail for scene-local command routing, or the trimmed message text seen by a
	// scene step/message handler. It is empty when the current routing path does
	// not derive text input.
	Text string
	// Args contains parsed command or payload arguments for the current routing
	// path. It is nil or empty when no argument vector is derived.
	Args []string

	errorTemplate string
	l10n          *L10n
	draftProvider *DraftProvider
	payloadType   BotPayloadType
	sceneRuntime  sceneRuntime
	observer      Observer

	ctx context.Context
}

// AnswerMessage represents a message sent or edited via MsgContext.
// It holds metadata to allow further editing or deletion.
type AnswerMessage struct {
	MessageID int
	Text      string
	IsMedia   bool
	ctx       *MsgContext // internal back-reference
}

// Internal helper for text edits with optional keyboard and parse mode.
func (ctx *MsgContext) edit(messageId int, text string, keyboard *InlineKeyboard, parseMode tgapi.ParseMode) *AnswerMessage {
	if err := validateMessageText(text); err != nil {
		ctx.Logger.Errorln(err)
		return nil
	}
	params := tgapi.EditMessageText{
		Text:      text,
		ParseMode: parseMode,
	}
	switch {
	case messageId > 0 && ctx.Msg != nil:
		params.MessageID = messageId
		params.ChatID = ctx.Msg.Chat.ID
	case ctx.InlineMsgId != "":
		params.InlineMessageID = ctx.InlineMsgId
	default:
		ctx.Logger.Errorln(ErrEditTargetMissing)
		return nil
	}
	if keyboard != nil {
		params.ReplyMarkup = keyboard.Get()
	}
	msg, _, err := ctx.Api.EditMessageTextWithContext(ctx.Context(), params)
	if err != nil {
		ctx.Logger.Errorln(err)
		return nil
	}
	resultMessageID := messageId
	if msg.MessageID > 0 {
		resultMessageID = msg.MessageID
	}
	return &AnswerMessage{
		MessageID: resultMessageID, ctx: ctx, Text: text, IsMedia: false,
	}
}

// Edit replaces the text of the message without changing the keyboard or parse mode.
// Uses ParseNone (plain text).
func (m *AnswerMessage) Edit(text string) *AnswerMessage {
	return m.ctx.edit(m.MessageID, text, nil, tgapi.ParseNone)
}

// EditMarkdown replaces the text of the message using MarkdownV2 formatting.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
// Unescaped input may cause Telegram API errors or broken formatting.
func (m *AnswerMessage) EditMarkdown(text string) *AnswerMessage {
	return m.ctx.edit(m.MessageID, text, nil, tgapi.ParseMDV2)
}

// Internal helper for editing callback-linked messages.
func (ctx *MsgContext) editCallback(text string, keyboard *InlineKeyboard, parseMode tgapi.ParseMode) *AnswerMessage {
	if ctx.CallbackMsgId == 0 && ctx.InlineMsgId == "" {
		ctx.Logger.Errorln(ErrCallbackMessageMissing)
		return nil
	}
	return ctx.edit(ctx.CallbackMsgId, text, keyboard, parseMode)
}

// EditCallback edits the callback message using plain text (ParseNone).
func (ctx *MsgContext) EditCallback(text string, keyboard *InlineKeyboard) *AnswerMessage {
	return ctx.editCallback(text, keyboard, tgapi.ParseNone)
}

// EditCallbackMarkdown edits the callback message using MarkdownV2.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
func (ctx *MsgContext) EditCallbackMarkdown(text string, keyboard *InlineKeyboard) *AnswerMessage {
	return ctx.editCallback(text, keyboard, tgapi.ParseMDV2)
}

// EditCallbackf formats a string using fmt.Sprintf and edits the callback message with plain text.
func (ctx *MsgContext) EditCallbackf(format string, keyboard *InlineKeyboard, args ...any) *AnswerMessage {
	return ctx.editCallback(fmt.Sprintf(format, args...), keyboard, tgapi.ParseNone)
}

// EditCallbackfMarkdown formats a string using fmt.Sprintf and edits the callback message with MarkdownV2.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
func (ctx *MsgContext) EditCallbackfMarkdown(format string, keyboard *InlineKeyboard, args ...any) *AnswerMessage {
	return ctx.editCallback(fmt.Sprintf(format, args...), keyboard, tgapi.ParseMDV2)
}

// Internal helper for media-caption edits.
func (ctx *MsgContext) editPhotoText(messageId int, text string, kb *InlineKeyboard, parseMode tgapi.ParseMode) *AnswerMessage {
	if err := validateCaptionText(text); err != nil {
		ctx.Logger.Errorln(err)
		return nil
	}
	params := tgapi.EditMessageCaption{
		Caption:   text,
		ParseMode: parseMode,
	}
	switch {
	case messageId > 0 && ctx.Msg != nil:
		params.ChatID = ctx.Msg.Chat.ID
		params.MessageID = messageId
	case ctx.InlineMsgId != "":
		params.InlineMessageID = ctx.InlineMsgId
	default:
		ctx.Logger.Errorln(ErrEditTargetMissing)
		return nil
	}
	if kb != nil {
		params.ReplyMarkup = kb.Get()
	}

	msg, _, err := ctx.Api.EditMessageCaptionWithContext(ctx.Context(), params)
	if err != nil {
		ctx.Logger.Errorln(err)
		return nil
	}
	resultMessageID := messageId
	if msg.MessageID > 0 {
		resultMessageID = msg.MessageID
	}
	return &AnswerMessage{
		MessageID: resultMessageID, ctx: ctx, Text: text, IsMedia: true,
	}
}

// EditCaption edits the caption of a media message using plain text.
func (m *AnswerMessage) EditCaption(text string) *AnswerMessage {
	return m.ctx.editPhotoText(m.MessageID, text, nil, tgapi.ParseNone)
}

// EditCaptionMarkdown edits the caption of a media message using MarkdownV2.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
func (m *AnswerMessage) EditCaptionMarkdown(text string) *AnswerMessage {
	return m.ctx.editPhotoText(m.MessageID, text, nil, tgapi.ParseMDV2)
}

// EditCaptionKeyboard edits the caption of a media message with a new inline keyboard (plain text).
func (m *AnswerMessage) EditCaptionKeyboard(text string, kb *InlineKeyboard) *AnswerMessage {
	return m.ctx.editPhotoText(m.MessageID, text, kb, tgapi.ParseNone)
}

// EditCaptionKeyboardMarkdown edits the caption of a media message with a new inline keyboard using MarkdownV2.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
func (m *AnswerMessage) EditCaptionKeyboardMarkdown(text string, kb *InlineKeyboard) *AnswerMessage {
	return m.ctx.editPhotoText(m.MessageID, text, kb, tgapi.ParseMDV2)
}

// Internal helper for message replies with optional keyboard and parse mode.
func (ctx *MsgContext) answer(text string, keyboard *InlineKeyboard, parseMode tgapi.ParseMode) *AnswerMessage {
	if ctx.Msg == nil {
		ctx.Logger.Errorln(ErrMessageContextNil)
		return nil
	}
	if err := validateMessageText(text); err != nil {
		ctx.Logger.Errorln(err)
		return nil
	}
	params := tgapi.SendMessage{
		ChatID:    ctx.Msg.Chat.ID,
		Text:      text,
		ParseMode: parseMode,
	}
	if keyboard != nil {
		params.ReplyMarkup = keyboard.Get()
	}
	if ctx.Msg.MessageThreadID > 0 {
		params.MessageThreadID = ctx.Msg.MessageThreadID
	}
	if ctx.Msg.DirectMessageTopic != nil {
		params.DirectMessagesTopicID = ctx.Msg.DirectMessageTopic.TopicID
	}

	msg, err := ctx.Api.SendMessageWithContext(ctx.Context(), params)
	if err != nil {
		ctx.Logger.Errorln(err)
		return nil
	}
	return &AnswerMessage{
		MessageID: msg.MessageID, ctx: ctx, IsMedia: false, Text: text,
	}
}

// Answer sends a plain text message (ParseNone).
func (ctx *MsgContext) Answer(text string) *AnswerMessage {
	return ctx.answer(text, nil, tgapi.ParseNone)
}

// AnswerLong sends one or more plain-text messages if text exceeds Telegram's limit.
//
// The text is split into Telegram-safe chunks. Returned messages preserve send
// order. If a chunk fails to send, already-sent messages are returned.
func (ctx *MsgContext) AnswerLong(text string) []*AnswerMessage {
	return ctx.answerLong(text, nil, tgapi.ParseNone)
}

// AnswerMarkdown sends a message using MarkdownV2 formatting.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
func (ctx *MsgContext) AnswerMarkdown(text string) *AnswerMessage {
	return ctx.answer(text, nil, tgapi.ParseMDV2)
}

// Answerf formats a string using fmt.Sprintf and sends it as a plain text message.
func (ctx *MsgContext) Answerf(template string, args ...any) *AnswerMessage {
	return ctx.answer(fmt.Sprintf(template, args...), nil, tgapi.ParseNone)
}

// AnswerLongf formats a string using fmt.Sprintf and sends it as one or more plain-text messages.
func (ctx *MsgContext) AnswerLongf(template string, args ...any) []*AnswerMessage {
	return ctx.answerLong(fmt.Sprintf(template, args...), nil, tgapi.ParseNone)
}

// AnswerfMarkdown formats a string using fmt.Sprintf and sends it using MarkdownV2.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
func (ctx *MsgContext) AnswerfMarkdown(template string, args ...any) *AnswerMessage {
	return ctx.answer(fmt.Sprintf(template, args...), nil, tgapi.ParseMDV2)
}

// Keyboard sends a message with an inline keyboard (plain text).
func (ctx *MsgContext) Keyboard(text string, kb *InlineKeyboard) *AnswerMessage {
	return ctx.answer(text, kb, tgapi.ParseNone)
}

// KeyboardLong sends long plain text split across multiple messages.
//
// The inline keyboard is attached only to the final chunk.
func (ctx *MsgContext) KeyboardLong(text string, kb *InlineKeyboard) []*AnswerMessage {
	return ctx.answerLong(text, kb, tgapi.ParseNone)
}

// KeyboardMarkdown sends a message with an inline keyboard using MarkdownV2.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
func (ctx *MsgContext) KeyboardMarkdown(text string, keyboard *InlineKeyboard) *AnswerMessage {
	return ctx.answer(text, keyboard, tgapi.ParseMDV2)
}

func (ctx *MsgContext) answerLong(text string, keyboard *InlineKeyboard, parseMode tgapi.ParseMode) []*AnswerMessage {
	if parseMode != tgapi.ParseNone {
		ctx.Logger.Errorln(ErrMessageSplitImpossible)
		return nil
	}
	if ctx.Msg == nil {
		ctx.Logger.Errorln(ErrMessageContextNil)
		return nil
	}
	if err := validateMessageText(text); err == nil {
		msg := ctx.answer(text, keyboard, parseMode)
		if msg == nil {
			return nil
		}
		return []*AnswerMessage{msg}
	} else if !errors.Is(err, ErrMessageTooLong) {
		ctx.Logger.Errorln(err)
		return nil
	}

	parts := SplitMessageText(text)
	messages := make([]*AnswerMessage, 0, len(parts))
	for i, part := range parts {
		partKeyboard := (*InlineKeyboard)(nil)
		if i == len(parts)-1 {
			partKeyboard = keyboard
		}
		msg := ctx.answer(part, partKeyboard, parseMode)
		if msg == nil {
			break
		}
		messages = append(messages, msg)
	}
	if len(messages) == 0 {
		return nil
	}
	return messages
}

// Internal helper for photo replies with optional caption and keyboard.
func (ctx *MsgContext) answerPhoto(photoId, text string, kb *InlineKeyboard, parseMode tgapi.ParseMode) *AnswerMessage {
	if ctx.Msg == nil {
		ctx.Logger.Errorln(ErrMessageContextNil)
		return nil
	}
	if err := validateCaptionText(text); err != nil {
		ctx.Logger.Errorln(err)
		return nil
	}
	params := tgapi.SendPhoto{
		ChatID:    ctx.Msg.Chat.ID,
		Caption:   text,
		ParseMode: parseMode,
		Photo:     photoId,
	}
	if kb != nil {
		params.ReplyMarkup = kb.Get()
	}
	if ctx.Msg.MessageThreadID > 0 {
		params.MessageThreadID = ctx.Msg.MessageThreadID
	}
	if ctx.Msg.DirectMessageTopic != nil {
		params.DirectMessagesTopicID = int(ctx.Msg.DirectMessageTopic.TopicID)
	}

	msg, err := ctx.Api.SendPhotoWithContext(ctx.Context(), params)
	if err != nil {
		ctx.Logger.Errorln(err)
		return nil
	}
	return &AnswerMessage{
		MessageID: msg.MessageID, ctx: ctx, Text: text, IsMedia: true,
	}
}

// AnswerPhoto sends a photo with plain text caption.
func (ctx *MsgContext) AnswerPhoto(photoId, text string) *AnswerMessage {
	return ctx.answerPhoto(photoId, text, nil, tgapi.ParseNone)
}

// AnswerPhotoMarkdown sends a photo with MarkdownV2 caption.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
func (ctx *MsgContext) AnswerPhotoMarkdown(photoId, text string) *AnswerMessage {
	return ctx.answerPhoto(photoId, text, nil, tgapi.ParseMDV2)
}

// AnswerPhotoKeyboard sends a photo with caption and inline keyboard (plain text).
func (ctx *MsgContext) AnswerPhotoKeyboard(photoId, text string, kb *InlineKeyboard) *AnswerMessage {
	return ctx.answerPhoto(photoId, text, kb, tgapi.ParseNone)
}

// AnswerPhotoKeyboardMarkdown sends a photo with caption and inline keyboard using MarkdownV2.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
func (ctx *MsgContext) AnswerPhotoKeyboardMarkdown(photoId, text string, kb *InlineKeyboard) *AnswerMessage {
	return ctx.answerPhoto(photoId, text, kb, tgapi.ParseMDV2)
}

// AnswerPhotof formats a string and sends it as a photo caption (plain text).
func (ctx *MsgContext) AnswerPhotof(photoId, template string, args ...any) *AnswerMessage {
	return ctx.answerPhoto(photoId, fmt.Sprintf(template, args...), nil, tgapi.ParseNone)
}

// AnswerPhotofMarkdown formats a string and sends it as a photo caption using MarkdownV2.
//
// ⚠️ WARNING: User input must be escaped with laniakea.EscapeMarkdownV2() before passing here.
func (ctx *MsgContext) AnswerPhotofMarkdown(photoId, template string, args ...any) *AnswerMessage {
	return ctx.answerPhoto(photoId, fmt.Sprintf(template, args...), nil, tgapi.ParseMDV2)
}

// Internal helper that deletes a message by ID.
func (ctx *MsgContext) delete(messageId int) {
	if messageId == 0 {
		ctx.Logger.Errorln(ErrMessageIDZero)
		return
	}
	if ctx.Msg == nil {
		ctx.Logger.Errorln(ErrMessageContextNil)
		return
	}
	_, err := ctx.Api.DeleteMessageWithContext(ctx.Context(), tgapi.DeleteMessage{
		ChatID:    ctx.Msg.Chat.ID,
		MessageID: messageId,
	})
	if err != nil {
		ctx.Logger.Errorln(err)
	}
}

// Delete removes the message associated with this AnswerMessage.
func (m *AnswerMessage) Delete() { m.ctx.delete(m.MessageID) }

// CallbackDelete deletes the message that triggered the callback query.
func (ctx *MsgContext) CallbackDelete() {
	if ctx.CallbackMsgId == 0 {
		ctx.Logger.Errorln(ErrCallbackMessageMissing)
		return
	}
	ctx.delete(ctx.CallbackMsgId)
}

// Internal helper that answers a callback query with optional text, alert, or URL.
func (ctx *MsgContext) answerCallbackQuery(url, text string, showAlert bool) {
	if len(ctx.CallbackQueryId) == 0 {
		return
	}
	_, err := ctx.Api.AnswerCallbackQueryWithContext(ctx.Context(), tgapi.AnswerCallbackQuery{
		CallbackQueryID: ctx.CallbackQueryId,
		Text:            text, ShowAlert: showAlert, URL: url,
	})
	if err != nil {
		ctx.Logger.Errorln(err)
	}
}

// AnswerCbQuery answers the callback query with no text or alert.
func (ctx *MsgContext) AnswerCbQuery() { ctx.answerCallbackQuery("", "", false) }

// AnswerCbQueryText answers the callback query with a text notification.
func (ctx *MsgContext) AnswerCbQueryText(text string) { ctx.answerCallbackQuery("", text, false) }

// AnswerCbQueryAlert answers the callback query with a user-visible alert.
func (ctx *MsgContext) AnswerCbQueryAlert(text string) { ctx.answerCallbackQuery("", text, true) }

// AnswerCbQueryUrl answers the callback query with a URL redirect.
func (ctx *MsgContext) AnswerCbQueryUrl(u string) { ctx.answerCallbackQuery(u, "", false) }

// SendAction sends a chat action (typing, uploading_photo, etc.) to indicate bot activity.
func (ctx *MsgContext) SendAction(action tgapi.ChatActionType) {
	if ctx.Msg == nil {
		ctx.Logger.Errorln("Can't send action without chat message context")
		return
	}
	params := tgapi.SendChatAction{
		ChatID: ctx.Msg.Chat.ID, Action: action,
	}
	if ctx.Msg.MessageThreadID > 0 {
		params.MessageThreadID = ctx.Msg.MessageThreadID
	}
	_, err := ctx.Api.SendChatActionWithContext(ctx.Context(), params)
	if err != nil {
		ctx.Logger.Errorln(err)
	}
}

// Internal helper that formats, sends, and logs an error.
func (ctx *MsgContext) error(err error) {
	if err == nil {
		return
	}
	ctx.Logger.Errorln(err)
	if IsInternalError(err) {
		return
	}
	text := fmt.Sprintf(ctx.errorTemplate, err.Error())

	if ctx.CallbackQueryId != "" {
		ctx.answerCallbackQuery("", text, false)
	} else {
		ctx.answer(text, nil, tgapi.ParseNone)
	}
}

// Error is an alias for error().
func (ctx *MsgContext) Error(err error) { ctx.error(err) }

func (ctx *MsgContext) newDraft(parseMode tgapi.ParseMode) *Draft {
	if ctx.Msg == nil {
		ctx.Logger.Errorln(ErrMessageContextNil)
		return nil
	}
	if ctx.Api == nil {
		ctx.Logger.Errorln(ErrAPIIsNil)
		return nil
	}
	if ctx.draftProvider == nil {
		ctx.Logger.Errorln(ErrDraftProviderNil)
		return nil
	}

	if ctx.Api.Limiter != nil {
		c, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
		defer cancel()
		if err := ctx.Api.Limiter.Wait(c, ctx.Msg.Chat.ID); err != nil {
			ctx.Logger.Errorln(err)
			return nil
		}
	}

	draft := ctx.draftProvider.NewDraft(parseMode).SetChat(ctx.Msg.Chat.ID, ctx.Msg.MessageThreadID)
	return draft
}

// NewDraft creates a new message draft associated with the current chat.
// Uses the API limiter to avoid rate limiting.
func (ctx *MsgContext) NewDraft() *Draft {
	return ctx.newDraft(tgapi.ParseNone)
}

// NewDraftMarkdown creates a new message draft associated with the current chat,
// with Markdown V2 parse mode enabled.
// Uses the API limiter to avoid rate limiting.
func (ctx *MsgContext) NewDraftMarkdown() *Draft {
	return ctx.newDraft(tgapi.ParseMDV2)
}

// Translate looks up a key in the current user's language.
// Falls back to the bot's default language if user's language is unknown or unsupported.
func (ctx *MsgContext) Translate(key string) string {
	if ctx.From == nil {
		return key
	}
	lang := Val(ctx.From.LanguageCode, ctx.l10n.GetFallbackLanguage())
	return ctx.l10n.Translate(lang, key)
}

// NewInlineKeyboard creates a new keyboard builder with the context's payload
// encoding type and the specified maximum number of buttons per row.
func (ctx *MsgContext) NewInlineKeyboard(maxRow int) *InlineKeyboard {
	return NewInlineKeyboard(ctx.payloadType, maxRow)
}

func bindPositional(args []string, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return ErrBindArgsTargetNotPointer
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrBindArgsTargetNotStruct
	}

	t := v.Type()
	fields := make([]int, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}
		fields = append(fields, i)
	}

	argIndex := 0
	for fieldPos, fieldIndex := range fields {
		field := v.Field(fieldIndex)
		fieldType := t.Field(fieldIndex)

		if argIndex >= len(args) {
			// Leave trailing fields at their zero values when arguments run out.
			break
		}

		isLastBindableField := fieldPos == len(fields)-1

		raw := args[argIndex]
		if isLastBindableField && field.Kind() == reflect.String {
			raw = strings.Join(args[argIndex:], " ")
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(raw)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: field %s: %v", ErrBindArgsConversion, fieldType.Name, err)
			}
			field.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(raw, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: field %s: %v", ErrBindArgsConversion, fieldType.Name, err)
			}
			field.SetUint(n)
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return fmt.Errorf("%w: field %s: %v", ErrBindArgsConversion, fieldType.Name, err)
			}
			field.SetFloat(f)
		case reflect.Bool:
			b, err := strconv.ParseBool(raw)
			if err != nil {
				return fmt.Errorf("%w: field %s: %v", ErrBindArgsConversion, fieldType.Name, err)
			}
			field.SetBool(b)
		default:
			return fmt.Errorf("%w: field %s: %s", ErrBindArgsUnsupportedFieldType, fieldType.Name, field.Kind())
		}

		if isLastBindableField && field.Kind() == reflect.String {
			break
		}
		argIndex++
	}

	return nil
}

// BindArgs binds positional command arguments from ctx.Args into dst.
//
// Exported struct fields are filled in declaration order. When fewer arguments
// are provided than fields, the remaining fields keep their zero values. If the
// final bindable field is a string, it receives the remaining arguments joined
// with spaces.
func (ctx *MsgContext) BindArgs(dst any) error {
	return bindPositional(ctx.Args, dst)
}

// Context returns the request-scoped context associated with the current update.
func (ctx *MsgContext) Context() context.Context {
	if ctx.ctx == nil {
		return context.Background()
	}
	return ctx.ctx
}

func (ctx *MsgContext) emitPolicyChecked(event PolicyCheckedEvent) {
	if ctx == nil || ctx.observer == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			if ctx.Logger != nil {
				ctx.Logger.Errorln(fmt.Sprintf("panic in observer policy event: %v", r))
				return
			}
			log.Printf("panic in observer policy event: %v", r)
		}
	}()
	ctx.observer.OnPolicyChecked(ctx.Context(), event)
}

// EnterScene enters the named scene at its configured entry step.
func (ctx *MsgContext) EnterScene(name string) error {
	if ctx.sceneRuntime == nil {
		return ErrSceneRuntimeNil
	}

	scene, ok := ctx.sceneRuntime.findScene(name)
	if !ok {
		return ErrSceneNotFound
	}

	key, ok := ctx.sceneRuntime.buildSceneKey(scene.Scope, ctx)
	if !ok {
		return ErrCantFindSession
	}
	if scene.Entry == "" {
		return ErrSceneEntryNotSet
	}
	if _, ok := scene.Steps[scene.Entry]; !ok {
		return ErrSceneStepNotFound
	}

	session := SceneSession{
		Scene: scene.Name,
		Step:  scene.Entry,
	}

	return ctx.sceneRuntime.setSession(key, session)
}

// EnterSceneStep enters the named scene at a specific step.
func (ctx *MsgContext) EnterSceneStep(name, step string) error {
	if ctx.sceneRuntime == nil {
		return ErrSceneRuntimeNil
	}

	scene, ok := ctx.sceneRuntime.findScene(name)
	if !ok {
		return ErrSceneNotFound
	}
	if _, ok := scene.Steps[step]; !ok {
		return ErrSceneStepNotFound
	}

	key, ok := ctx.sceneRuntime.buildSceneKey(scene.Scope, ctx)
	if !ok {
		return ErrCantFindSession
	}

	session := SceneSession{
		Scene: scene.Name,
		Step:  step,
	}

	return ctx.sceneRuntime.setSession(key, session)
}

// ExitScene leaves the currently active scene for this context.
func (ctx *MsgContext) ExitScene() error {
	if ctx.sceneRuntime == nil {
		return ErrSceneRuntimeNil
	}

	_, session, err := ctx.sceneRuntime.findSceneSession(ctx)
	if err != nil {
		return err
	}
	if session.Scene == "" {
		return ErrNotInScene
	}

	scene, ok := ctx.sceneRuntime.findScene(session.Scene)
	if !ok {
		return ErrSceneNotFound
	}

	key, ok := ctx.sceneRuntime.buildSceneKey(scene.Scope, ctx)
	if !ok {
		return ErrCantFindSession
	}

	return ctx.sceneRuntime.deleteSession(key)
}
