package laniakea

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"git.nix13.pw/scuroneko/extypes"
	"git.nix13.pw/scuroneko/slog"
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

	logger        *slog.Logger
	requestLogger *slog.Logger

	plugins     []Plugin
	middlewares []Middleware
	prefixes    []string
	runners     []Runner

	dbContext *DatabaseContext
	api       *Api

	dbWriterRequested extypes.Slice[*slog.Logger]

	updateOffset int
	updateTypes  []string
	updateQueue  *extypes.Queue[*Update]
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

func LoadPrefixesFromEnv() []string {
	prefixesS, exists := os.LookupEnv("PREFIXES")
	if !exists {
		return []string{"!"}
	}
	return strings.Split(prefixesS, ";")
}
func NewBot(settings *BotSettings) *Bot {
	updateQueue := extypes.CreateQueue[*Update](256)
	api := NewAPI(settings.Token)
	bot := &Bot{
		updateOffset: 0, plugins: make([]Plugin, 0), debug: settings.Debug, errorTemplate: "%s",
		prefixes: settings.Prefixes, updateTypes: make([]string, 0), runners: make([]Runner, 0),
		updateQueue: updateQueue, api: api, dbWriterRequested: make([]*slog.Logger, 0),
		token: settings.Token,
	}
	bot.dbWriterRequested = bot.dbWriterRequested.Push(api.logger)

	if len(settings.ErrorTemplate) > 0 {
		bot.errorTemplate = settings.ErrorTemplate
	}
	if len(settings.LoggerBasePath) == 0 {
		settings.LoggerBasePath = "./"
	}

	level := slog.FATAL
	if settings.Debug {
		level = slog.DEBUG
	}

	bot.logger = slog.CreateLogger().Level(level).Prefix("BOT")
	bot.logger.AddWriter(bot.logger.CreateJsonStdoutWriter())
	if settings.WriteToFile {
		path := fmt.Sprintf("%s/main.log", strings.TrimRight(settings.LoggerBasePath, "/"))
		fileWriter, err := bot.logger.CreateTextFileWriter(path)
		if err != nil {
			bot.logger.Fatal(err)
		}
		bot.logger.AddWriter(fileWriter)
	}

	if settings.UseRequestLogger {
		bot.requestLogger = slog.CreateLogger().Level(level).Prefix("REQUESTS")
		bot.requestLogger.AddWriter(bot.requestLogger.CreateJsonStdoutWriter())
		if settings.WriteToFile {
			path := fmt.Sprintf("%s/requests.log", strings.TrimRight(settings.LoggerBasePath, "/"))
			fileWriter, err := bot.requestLogger.CreateTextFileWriter(path)
			if err != nil {
				bot.logger.Fatal(err)
			}
			bot.requestLogger.AddWriter(fileWriter)
		}
	}

	u, err := api.GetMe()
	if err != nil {
		bot.logger.Fatal(err)
	}
	bot.logger.Infof("Authorized as %s\n", u.FirstName)

	return bot
}

func (b *Bot) Close() {
	err := b.logger.Close()
	if err != nil {
		log.Println(err)
	}
	err = b.requestLogger.Close()
	if err != nil {
		log.Println(err)
	}
}

type DatabaseContext struct {
	PostgresSQL *sqlx.DB
	MongoDB     *mongo.Client
	Redis       *redis.Client
}

func (b *Bot) AddDatabaseLogger(writer func(db *DatabaseContext) slog.LoggerWriter) *Bot {
	w := writer(b.dbContext)
	b.logger.AddWriter(w)
	if b.requestLogger != nil {
		b.requestLogger.AddWriter(w)
	}
	for _, l := range b.dbWriterRequested {
		l.AddWriter(w)
	}
	return b
}

func (b *Bot) DatabaseContext(ctx *DatabaseContext) *Bot {
	b.dbContext = ctx
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
func (b *Bot) ErrorTemplate(s string) *Bot {
	b.errorTemplate = s
	return b
}
func (b *Bot) Debug(debug bool) *Bot {
	b.debug = debug
	return b
}
func (b *Bot) AddPlugins(plugin ...Plugin) *Bot {
	b.plugins = append(b.plugins, plugin...)
	for _, p := range plugin {
		b.logger.Debugln(fmt.Sprintf("plugins with name \"%s\" registered", p.Name))
	}
	return b
}
func (b *Bot) AddMiddleware(middleware ...Middleware) *Bot {
	b.middlewares = append(b.middlewares, middleware...)
	for _, m := range middleware {
		b.logger.Debugln(fmt.Sprintf("middleware with name \"%s\" registered", m.Name))
	}

	sort.Slice(b.middlewares, func(i, j int) bool {
		first := b.middlewares[i]
		second := b.middlewares[j]
		if first.Order == second.Order {
			return first.Name < second.Name
		}
		return first.Order < second.Order
	})

	return b
}
func (b *Bot) AddRunner(runner Runner) *Bot {
	b.runners = append(b.runners, runner)
	b.logger.Debugln(fmt.Sprintf("runner with name \"%s\" registered", runner.Name))
	return b
}
func (b *Bot) Logger() *slog.Logger {
	return b.logger
}
func (b *Bot) GetDBContext() *DatabaseContext {
	return b.dbContext
}

func (b *Bot) Run() {
	if len(b.prefixes) == 0 {
		b.logger.Fatalln("no prefixes defined")
		return
	}

	if len(b.plugins) == 0 {
		b.logger.Fatalln("no plugins defined")
		return
	}

	b.ExecRunners()

	b.logger.Infoln("Bot running. Press CTRL+C to exit.")
	go func() {
		for {
			_, err := b.Updates()
			if err != nil {
				b.logger.Errorln(err)
			}
		}
	}()

	for {
		if b.updateQueue.IsEmpty() {
			time.Sleep(time.Millisecond * 25)
			continue
		}

		u := b.updateQueue.Dequeue()
		if u == nil {
			b.logger.Errorln("update is nil")
			continue
		}

		b.handle(u)
	}
}
