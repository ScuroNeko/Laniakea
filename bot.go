package laniakea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vinovest/sqlx"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ParseMode string

const (
	ParseMDV2 ParseMode = "MarkdownV2"
	ParseHTML ParseMode = "HTML"
	ParseMD   ParseMode = "Markdown"
)

type Bot struct {
	token         string
	debug         bool
	errorTemplate string

	logger        *Logger
	requestLogger *Logger

	plugins     []*Plugin
	middlewares []*Middleware
	prefixes    []string

	dbContext *DatabaseContext

	updateOffset int
	updateTypes  []string
	updateQueue  *Queue[*Update]
}

type BotSettings struct {
	Token            string
	Debug            bool
	ErrorTemplate    string
	Prefixes         []string
	UpdateTypes      []string
	LoggerBasePath   string
	UseRequestLogger bool
	WriteToFile      bool
}

func LoadSettingsFromEnv() *BotSettings {
	return &BotSettings{
		Token:            os.Getenv("TG_TOKEN"),
		Debug:            os.Getenv("DEBUG") == "true",
		ErrorTemplate:    os.Getenv("ERROR_TEMPLATE"),
		Prefixes:         LoadPrefixesFromEnv(),
		UpdateTypes:      strings.Split(os.Getenv("UPDATE_TYPES"), ";"),
		UseRequestLogger: os.Getenv("USE_REQ_LOG") == "true",
		WriteToFile:      os.Getenv("WRITE_TO_FILE") == "true",
	}
}

type MsgContext struct {
	Bot           *Bot
	Msg           *Message
	Update        *Update
	From          *User
	CallbackMsgId int
	FromID        int
	Prefix        string
	Text          string
	Args          []string
}

type DatabaseContext struct {
	PostgresSQL *sqlx.DB
	MongoDB     *mongo.Client
	Redis       *redis.Client
}

func NewBot(settings *BotSettings) *Bot {
	updateQueue := CreateQueue[*Update](256)
	bot := &Bot{
		updateOffset: 0, plugins: make([]*Plugin, 0), debug: settings.Debug, errorTemplate: "%s",
		prefixes: settings.Prefixes, updateTypes: make([]string, 0),
		updateQueue: updateQueue,
		token:       settings.Token,
	}

	if len(settings.ErrorTemplate) > 0 {
		bot.errorTemplate = settings.ErrorTemplate
	}

	if len(settings.LoggerBasePath) == 0 {
		settings.LoggerBasePath = "./"
	}
	level := FATAL
	if settings.Debug {
		level = DEBUG
	}
	bot.logger = CreateLogger().Level(level)
	if settings.WriteToFile {
		bot.logger = bot.logger.OpenFile(fmt.Sprintf("%s/main.log", strings.TrimRight(settings.LoggerBasePath, "/")))
	}
	bot.logger = bot.logger.PrintTraceback(true)
	if settings.UseRequestLogger {
		bot.requestLogger = CreateLogger().Level(level).Prefix("REQUESTS")
		if settings.WriteToFile {
			bot.requestLogger = bot.requestLogger.OpenFile(fmt.Sprintf("%s/requests.log", strings.TrimRight(settings.LoggerBasePath, "/")))
		}
	}

	return bot
}

func (b *Bot) Close() {
	err := b.logger.f.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("log closed")
	}
}

func (b *Bot) InitDatabaseContext(ctx *DatabaseContext) *Bot {
	b.dbContext = ctx
	return b
}
func (b *Bot) AddDatabaseLogger(writer func(db *DatabaseContext) LoggerWriter) *Bot {
	w := []LoggerWriter{writer(b.dbContext)}
	b.logger.AddWriters(w)
	if b.requestLogger != nil {
		b.requestLogger.AddWriters(w)
	}
	return b
}

func (b *Bot) UpdateTypes(t ...string) *Bot {
	b.updateTypes = make([]string, 0)
	b.updateTypes = append(b.updateTypes, t...)
	return b
}
func (b *Bot) AddUpdateType(t ...string) *Bot {
	b.updateTypes = append(b.updateTypes, t...)
	return b
}

func (b *Bot) AddPrefixes(prefixes ...string) *Bot {
	b.prefixes = append(b.prefixes, prefixes...)
	return b
}

func LoadPrefixesFromEnv() []string {
	prefixesS, exists := os.LookupEnv("PREFIXES")
	if !exists {
		return []string{"!"}
	}
	return strings.Split(prefixesS, ";")
}

func (b *Bot) ErrorTemplate(s string) *Bot {
	b.errorTemplate = s
	return b
}

func (b *Bot) Debug(debug bool) *Bot {
	b.debug = debug
	return b
}

func (b *Bot) AddPlugins(plugin ...*Plugin) *Bot {
	b.plugins = append(b.plugins, plugin...)
	for _, p := range plugin {
		b.logger.Debug(fmt.Sprintf("plugins with name \"%s\" was registered", p.Name))
	}
	return b
}

func (b *Bot) AddMiddleware(middleware ...*Middleware) *Bot {
	sort.Slice(middleware, func(a, b int) bool {
		first := middleware[a]
		second := middleware[b]
		if first.Order == second.Order {
			return first.Name < second.Name
		}
		return middleware[a].Order < middleware[b].Order
	})

	b.middlewares = append(b.middlewares, middleware...)
	for _, m := range middleware {
		b.logger.Debug(fmt.Sprintf("middleware with name \"%s\" was registered", m.Name))
	}
	return b
}

func (b *Bot) Run() {
	if len(b.prefixes) == 0 {
		b.logger.Fatal("no prefixes defined")
		return
	}

	if len(b.plugins) == 0 {
		b.logger.Fatal("no plugins defined")
		return
	}

	b.logger.Info("Bot running. Press CTRL+C to exit.")

	go func() {
		for {
			_, err := b.Updates()
			if err != nil {
				b.logger.Error(err)
			}
		}
	}()

	for {
		queue := b.updateQueue
		if queue.IsEmpty() {
			time.Sleep(time.Millisecond * 25)
			continue
		}

		u := queue.Dequeue()
		if u == nil {
			b.logger.Error("update is nil")
			continue
		}
		ctx := &MsgContext{
			Bot:    b,
			Update: u,
		}
		for _, middleware := range b.middlewares {
			middleware.Execute(ctx, b.dbContext)
		}

		for _, plugin := range b.plugins {
			if plugin.UpdateListener != nil {
				(*plugin.UpdateListener)(ctx, b.dbContext)
			}
		}

		if u.CallbackQuery != nil {
			b.handleCallback(u, ctx)
		} else {
			b.handleMessage(u, ctx)
		}
	}
}

// {"callback_query":{"chat_instance":"6202057960757700762","data":"aboba","from":{"first_name":"scuroneko","id":314834933,"is_bot":false,"language_code":"ru","username":"scuroneko"},"id":"1352205741990111553","message":{"chat":{"first_name":"scuroneko","id":314834933,"type":"private","username":"scuroneko"},"date":1734338107,"from":{"first_name":"Kurumi","id":7718900880,"is_bot":true,"username":"kurumi_game_bot"},"message_id":19,"reply_markup":{"inline_keyboard":[[{"callback_data":"aboba","text":"Test"},{"callback_data":"another","text":"Another"}]]},"text":"Aboba"}},"update_id":350979488}

func (b *Bot) handleMessage(update *Update, ctx *MsgContext) {
	var text string
	if update.Message == nil {
		return
	}
	if len(update.Message.Text) > 0 {
		text = update.Message.Text
	} else {
		text = update.Message.Caption
	}

	ctx.FromID = update.Message.From.ID
	ctx.From = update.Message.From
	ctx.Msg = update.Message
	text = strings.TrimSpace(text)
	prefix, hasPrefix := b.checkPrefixes(text)
	if !hasPrefix {
		return
	}
	ctx.Prefix = prefix

	text = strings.TrimSpace(text[len(prefix):])

	for _, plugin := range b.plugins {
		// Check every command
		for cmd := range plugin.Commands {
			if !strings.HasPrefix(text, cmd) {
				continue
			}

			ctx.Text = strings.TrimSpace(text[len(cmd):])
			ctx.Args = strings.Split(ctx.Text, " ")

			go plugin.Execute(cmd, ctx, b.dbContext)
		}
	}
}

func (b *Bot) handleCallback(update *Update, ctx *MsgContext) {
	data := new(CallbackData)
	err := json.Unmarshal([]byte(update.CallbackQuery.Data), data)
	if err != nil {
		b.logger.Error(err)
		return
	}

	ctx.FromID = update.CallbackQuery.From.ID
	ctx.From = update.CallbackQuery.From
	ctx.Msg = update.CallbackQuery.Message
	ctx.CallbackMsgId = update.CallbackQuery.Message.MessageID
	ctx.Args = data.Args

	for _, plugin := range b.plugins {
		_, ok := plugin.Payloads[data.Command]
		if !ok {
			continue
		}
		go plugin.ExecutePayload(data.Command, ctx, b.dbContext)
		break
	}
}

func (b *Bot) checkPrefixes(text string) (string, bool) {
	for _, prefix := range b.prefixes {
		if strings.HasPrefix(text, prefix) {
			return prefix, true
		}
	}
	return "", false
}

type AnswerMessage struct {
	MessageID int
	Text      string
	IsMedia   bool
	Keyboard  *InlineKeyboard
	ctx       *MsgContext
}

func (ctx *MsgContext) edit(messageId int, text string, keyboard *InlineKeyboard) *AnswerMessage {
	params := &EditMessageTextP{
		MessageID: messageId,
		ChatID:    ctx.Msg.Chat.ID,
		Text:      text,
		ParseMode: ParseMD,
	}
	if keyboard != nil {
		params.ReplyMarkup = keyboard.Get()
	}
	msg, err := ctx.Bot.EditMessageText(params)
	if err != nil {
		ctx.Bot.logger.Error(err)
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
		ctx.Bot.logger.Error("Can't edit non-callback update message")
		return nil
	}

	return ctx.edit(ctx.CallbackMsgId, text, keyboard)
}
func (ctx *MsgContext) EditCallbackf(format string, keyboard *InlineKeyboard, args ...any) *AnswerMessage {
	return ctx.EditCallback(fmt.Sprintf(format, args...), keyboard)
}

func (ctx *MsgContext) answer(text string, keyboard *InlineKeyboard) *AnswerMessage {
	params := &SendMessageP{
		ChatID:    ctx.Msg.Chat.ID,
		Text:      text,
		ParseMode: ParseMD,
	}
	if keyboard != nil {
		params.ReplyMarkup = keyboard.Get()
	}

	msg, err := ctx.Bot.SendMessage(params)
	if err != nil {
		ctx.Bot.logger.Error(err)
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

func (ctx *MsgContext) AnswerPhoto(photoId string, text string) *AnswerMessage {
	_, err := ctx.Bot.SendPhoto(&SendPhotoP{
		ChatID:    ctx.Msg.Chat.ID,
		Caption:   text,
		Photo:     photoId,
		ParseMode: ParseMD,
	})
	if err != nil {
		ctx.Bot.logger.Error(err)
	}
	return &AnswerMessage{
		MessageID: ctx.Msg.MessageID, ctx: ctx, IsMedia: true, Text: text,
	}
}

func (ctx *MsgContext) delete(messageId int) {
	_, err := ctx.Bot.DeleteMessage(&DeleteMessageP{
		ChatID:    ctx.Msg.Chat.ID,
		MessageID: messageId,
	})
	if err != nil {
		ctx.Bot.logger.Error(err)
	}
}
func (m *AnswerMessage) Delete() {
	m.ctx.delete(m.MessageID)
}
func (ctx *MsgContext) CallbackDelete() {
	ctx.delete(ctx.CallbackMsgId)
}

func (ctx *MsgContext) Error(err error) {
	_, sendErr := ctx.Bot.SendMessage(&SendMessageP{
		ChatID: ctx.Msg.Chat.ID,
		Text:   fmt.Sprintf(ctx.Bot.errorTemplate, EscapeMarkdown(err.Error())),
	})
	ctx.Bot.logger.Error(err)

	if sendErr != nil {
		ctx.Bot.logger.Error(sendErr)
	}
}

func (b *Bot) Logger() *Logger {
	return b.logger
}

type ApiResponse struct {
	Ok          bool           `json:"ok"`
	Result      map[string]any `json:"result,omitempty"`
	Description string         `json:"description,omitempty"`
	ErrorCode   int            `json:"error_code,omitempty"`
}

type ApiResponseA struct {
	Ok          bool   `json:"ok"`
	Result      []any  `json:"result,omitempty"`
	Description string `json:"description,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
}

// request is a low-level call to api.
func (b *Bot) request(methodName string, params any) (map[string]interface{}, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(params)
	if err != nil {
		return nil, err
	}

	if b.debug && b.requestLogger != nil {
		b.requestLogger.Debug(strings.ReplaceAll(fmt.Sprintf(
			"POST https://api.telegram.org/bot%s/%s %s",
			"<TOKEN>",
			methodName,
			buf.String(),
		), "\n", ""))
	}
	r, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.token, methodName), "application/json", &buf)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	b.requestLogger.Debug(fmt.Sprintf("RES %s %s", methodName, string(data)))
	response := new(ApiResponse)

	var result map[string]any

	err = json.Unmarshal(data, &response)
	if err != nil {
		responseArray := new(ApiResponseA)
		err = json.Unmarshal(data, responseArray)
		if err != nil {
			return nil, err
		}
		result = map[string]interface{}{
			"data": responseArray.Result,
		}
	} else {
		result = response.Result
	}
	if !response.Ok {
		return nil, fmt.Errorf("[%d] %s", response.ErrorCode, response.Description)
	}
	return result, err
}
