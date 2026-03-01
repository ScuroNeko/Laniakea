package laniakea

import (
	"fmt"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
	"git.nix13.pw/scuroneko/laniakea/utils"
	"git.nix13.pw/scuroneko/slog"
)

type MsgContext struct {
	Api *tgapi.API

	Msg             *tgapi.Message
	Update          tgapi.Update
	From            *tgapi.User
	CallbackMsgId   int
	CallbackQueryId string
	FromID          int
	Prefix          string
	Text            string
	Args            []string

	errorTemplate string
	botLogger     *slog.Logger
	l10n          *L10n
	draftProvider *DraftProvider
}

type AnswerMessage struct {
	MessageID int
	Text      string
	IsMedia   bool
	ctx       *MsgContext
}

func (ctx *MsgContext) edit(messageId int, text string, keyboard *InlineKeyboard) *AnswerMessage {
	params := tgapi.EditMessageTextP{
		MessageID: messageId,
		ChatID:    ctx.Msg.Chat.ID,
		Text:      text,
		ParseMode: tgapi.ParseMD,
	}
	if keyboard != nil {
		params.ReplyMarkup = keyboard.Get()
	}
	msg, _, err := ctx.Api.EditMessageText(params)
	if err != nil {
		ctx.botLogger.Errorln(err)
		return nil
	}
	return &AnswerMessage{
		MessageID: msg.MessageID, ctx: ctx, Text: text, IsMedia: false,
	}
}
func (m *AnswerMessage) Edit(text string) *AnswerMessage {
	return m.ctx.edit(m.MessageID, text, nil)
}
func (ctx *MsgContext) EditCallback(text string, keyboard *InlineKeyboard) *AnswerMessage {
	if ctx.CallbackMsgId == 0 {
		ctx.botLogger.Errorln("Can't edit non-callback update message")
		return nil
	}

	return ctx.edit(ctx.CallbackMsgId, text, keyboard)
}
func (ctx *MsgContext) EditCallbackf(format string, keyboard *InlineKeyboard, args ...any) *AnswerMessage {
	return ctx.EditCallback(fmt.Sprintf(format, args...), keyboard)
}

func (ctx *MsgContext) editPhotoText(messageId int, text string, kb *InlineKeyboard) *AnswerMessage {
	params := tgapi.EditMessageCaptionP{
		ChatID:    ctx.Msg.Chat.ID,
		MessageID: messageId,
		Caption:   text,
		ParseMode: tgapi.ParseMD,
	}
	if kb != nil {
		params.ReplyMarkup = kb.Get()
	}

	msg, _, err := ctx.Api.EditMessageCaption(params)
	if err != nil {
		ctx.botLogger.Errorln(err)
	}
	return &AnswerMessage{
		MessageID: msg.MessageID, ctx: ctx, Text: text, IsMedia: true,
	}
}
func (m *AnswerMessage) EditCaption(text string) *AnswerMessage {
	if m.MessageID == 0 {
		m.ctx.botLogger.Errorln("Can't edit caption message, message id is zero")
		return m
	}
	return m.ctx.editPhotoText(m.MessageID, text, nil)
}
func (m *AnswerMessage) EditCaptionKeyboard(text string, kb *InlineKeyboard) *AnswerMessage {
	return m.ctx.editPhotoText(m.MessageID, text, kb)
}

func (ctx *MsgContext) answer(text string, keyboard *InlineKeyboard) *AnswerMessage {
	params := tgapi.SendMessageP{
		ChatID:    ctx.Msg.Chat.ID,
		Text:      text,
		ParseMode: tgapi.ParseMD,
	}
	if keyboard != nil {
		params.ReplyMarkup = keyboard.Get()
	}
	if ctx.Msg.MessageThreadID > 0 {
		params.MessageThreadID = ctx.Msg.MessageThreadID
	}

	msg, err := ctx.Api.SendMessage(params)
	if err != nil {
		ctx.botLogger.Errorln(err)
		return nil
	}
	return &AnswerMessage{
		MessageID: msg.MessageID, ctx: ctx, IsMedia: false, Text: text,
	}
}
func (ctx *MsgContext) Answer(text string) *AnswerMessage {
	return ctx.answer(text, nil)
}
func (ctx *MsgContext) Answerf(template string, args ...any) *AnswerMessage {
	return ctx.answer(fmt.Sprintf(template, args...), nil)
}
func (ctx *MsgContext) Keyboard(text string, kb *InlineKeyboard) *AnswerMessage {
	return ctx.answer(text, kb)
}

func (ctx *MsgContext) answerPhoto(photoId, text string, kb *InlineKeyboard) *AnswerMessage {
	params := tgapi.SendPhotoP{
		ChatID:    ctx.Msg.Chat.ID,
		Caption:   text,
		ParseMode: tgapi.ParseMD,
		Photo:     photoId,
	}
	if kb != nil {
		params.ReplyMarkup = kb.Get()
	}
	if ctx.Msg.MessageThreadID > 0 {
		params.MessageThreadID = ctx.Msg.MessageThreadID
	}

	msg, err := ctx.Api.SendPhoto(params)
	if err != nil {
		ctx.botLogger.Errorln(err)
		return &AnswerMessage{
			ctx: ctx, Text: text, IsMedia: true,
		}
	}
	return &AnswerMessage{
		MessageID: msg.MessageID, ctx: ctx, Text: text, IsMedia: true,
	}
}
func (ctx *MsgContext) AnswerPhoto(photoId, text string) *AnswerMessage {
	return ctx.answerPhoto(photoId, text, nil)
}
func (ctx *MsgContext) AnswerPhotoKeyboard(photoId, text string, kb *InlineKeyboard) *AnswerMessage {
	return ctx.answerPhoto(photoId, text, kb)
}

func (ctx *MsgContext) delete(messageId int) {
	_, err := ctx.Api.DeleteMessage(tgapi.DeleteMessageP{
		ChatID:    ctx.Msg.Chat.ID,
		MessageID: messageId,
	})
	if err != nil {
		ctx.botLogger.Errorln(err)
	}
}
func (m *AnswerMessage) Delete()        { m.ctx.delete(m.MessageID) }
func (ctx *MsgContext) CallbackDelete() { ctx.delete(ctx.CallbackMsgId) }

func (ctx *MsgContext) answerCallbackQuery(url, text string, showAlert bool) {
	if len(ctx.CallbackQueryId) == 0 {
		return
	}
	_, err := ctx.Api.AnswerCallbackQuery(tgapi.AnswerCallbackQueryP{
		CallbackQueryID: ctx.CallbackQueryId,
		Text:            text, ShowAlert: showAlert, URL: url,
	})
	if err != nil {
		ctx.botLogger.Errorln(err)
	}
}
func (ctx *MsgContext) AnswerCbQuery()                 { ctx.answerCallbackQuery("", "", false) }
func (ctx *MsgContext) AnswerCbQueryText(text string)  { ctx.answerCallbackQuery("", text, false) }
func (ctx *MsgContext) AnswerCbQueryAlert(text string) { ctx.answerCallbackQuery("", text, true) }
func (ctx *MsgContext) AnswerCbQueryUrl(u string)      { ctx.answerCallbackQuery(u, "", false) }

func (ctx *MsgContext) SendAction(action tgapi.ChatActionType) {
	params := tgapi.SendChatActionP{
		ChatID: ctx.Msg.Chat.ID, Action: action,
	}
	if ctx.Msg.MessageThreadID > 0 {
		params.MessageThreadID = ctx.Msg.MessageThreadID
	}
	_, err := ctx.Api.SendChatAction(params)
	if err != nil {
		ctx.botLogger.Errorln(err)
	}
}

func (ctx *MsgContext) error(err error) {
	text := fmt.Sprintf(ctx.errorTemplate, utils.EscapeMarkdown(err.Error()))

	if ctx.CallbackQueryId != "" {
		ctx.answerCallbackQuery("", text, false)
	} else {
		ctx.answer(text, nil)
	}
	ctx.botLogger.Errorln(err)
}
func (ctx *MsgContext) Error(err error) { ctx.error(err) }

func (ctx *MsgContext) NewDraft() *Draft { return ctx.draftProvider.NewDraft() }
func (ctx *MsgContext) Translate(key string) string {
	if ctx.From == nil {
		return key
	}
	lang := Val(ctx.From.LanguageCode, ctx.l10n.GetFallbackLanguage())
	return ctx.l10n.Translate(lang, key)
}
