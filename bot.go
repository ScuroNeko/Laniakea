package laniakea

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"git.nix13.pw/scuroneko/extypes"
	"git.nix13.pw/scuroneko/laniakea/tgapi"
	"git.nix13.pw/scuroneko/slog"
	"github.com/alitto/pond/v2"
	"golang.org/x/time/rate"
)

type BotOpts struct {
	Token       string
	UpdateTypes []string

	Debug         bool
	ErrorTemplate string
	Prefixes      []string

	LoggerBasePath   string
	UseRequestLogger bool
	WriteToFile      bool

	UseTestServer bool
	APIUrl        string

	RateLimit      int
	DropRLOverflow bool
}

func NewOpts() *BotOpts { return new(BotOpts) }
func LoadOptsFromEnv() *BotOpts {
	rateLimit := 30
	if rl := os.Getenv("RATE_LIMIT"); rl != "" {
		rateLimit, _ = strconv.Atoi(rl)
	}

	return &BotOpts{
		Token:       os.Getenv("TG_TOKEN"),
		UpdateTypes: strings.Split(os.Getenv("UPDATE_TYPES"), ";"),

		Debug:         os.Getenv("DEBUG") == "true",
		ErrorTemplate: os.Getenv("ERROR_TEMPLATE"),
		Prefixes:      LoadPrefixesFromEnv(),

		UseRequestLogger: os.Getenv("USE_REQ_LOG") == "true",
		WriteToFile:      os.Getenv("WRITE_TO_FILE") == "true",

		UseTestServer: os.Getenv("USE_TEST_SERVER") == "true",
		APIUrl:        os.Getenv("API_URL"),

		RateLimit:      rateLimit,
		DropRLOverflow: os.Getenv("DROP_RL_OVERFLOW") == "true",
	}
}
func LoadPrefixesFromEnv() []string {
	prefixesS, exists := os.LookupEnv("PREFIXES")
	if !exists {
		return []string{"/"}
	}
	return strings.Split(prefixesS, ";")
}

type DbContext interface{}
type NoDB struct{ DbContext }
type Bot[T DbContext] struct {
	token         string
	debug         bool
	errorTemplate string

	logger        *slog.Logger
	RequestLogger *slog.Logger
	extraLoggers  extypes.Slice[*slog.Logger]

	plugins     []Plugin[T]
	middlewares []Middleware[T]
	prefixes    []string
	runners     []Runner[T]

	api       *tgapi.API
	uploader  *tgapi.Uploader
	dbContext *T
	l10n      *L10n

	updateOffsetMu sync.Mutex
	updateOffset   int
	updateTypes    []tgapi.UpdateType
	updateQueue    chan *tgapi.Update
}

func NewBot[T any](opts *BotOpts) *Bot[T] {
	updateQueue := make(chan *tgapi.Update, 512)

	var limiter *rate.Limiter
	if opts.RateLimit > 0 {
		limiter = rate.NewLimiter(rate.Limit(opts.RateLimit), opts.RateLimit)
	}

	apiOpts := tgapi.NewAPIOpts(opts.Token).SetAPIUrl(opts.APIUrl).UseTestServer(opts.UseTestServer).SetLimiter(limiter)
	api := tgapi.NewAPI(apiOpts)

	uploader := tgapi.NewUploader(api)

	bot := &Bot[T]{
		updateOffset:  0,
		errorTemplate: "%s",
		updateQueue:   updateQueue,
		api:           api,
		uploader:      uploader,
		debug:         opts.Debug,
		prefixes:      opts.Prefixes,
		token:         opts.Token,
		plugins:       make([]Plugin[T], 0),
		updateTypes:   make([]tgapi.UpdateType, 0),
		runners:       make([]Runner[T], 0),
		extraLoggers:  make([]*slog.Logger, 0),
		l10n:          &L10n{},
	}
	bot.extraLoggers = bot.extraLoggers.Push(api.GetLogger()).Push(uploader.GetLogger())

	if len(opts.ErrorTemplate) > 0 {
		bot.errorTemplate = opts.ErrorTemplate
	}
	if len(opts.LoggerBasePath) == 0 {
		opts.LoggerBasePath = "./"
	}
	bot.initLoggers(opts)

	u, err := api.GetMe()
	if err != nil {
		_ = bot.Close()
		bot.logger.Fatal(err)
	}
	bot.logger.Infof("Authorized as %s\n", u.FirstName)

	return bot
}
func (bot *Bot[T]) Close() error {
	if err := bot.uploader.Close(); err != nil {
		bot.logger.Errorln(err)
	}
	if err := bot.api.CloseApi(); err != nil {
		bot.logger.Errorln(err)
	}
	if err := bot.RequestLogger.Close(); err != nil {
		bot.logger.Errorln(err)
	}
	if err := bot.logger.Close(); err != nil {
		return err
	}
	return nil
}
func (bot *Bot[T]) initLoggers(opts *BotOpts) {
	level := slog.FATAL
	if opts.Debug {
		level = slog.DEBUG
	}

	bot.logger = slog.CreateLogger().Level(level).Prefix("BOT")
	bot.logger.AddWriter(bot.logger.CreateJsonStdoutWriter())
	if opts.WriteToFile {
		path := fmt.Sprintf("%s/main.log", strings.TrimRight(opts.LoggerBasePath, "/"))
		fileWriter, err := bot.logger.CreateTextFileWriter(path)
		if err != nil {
			bot.logger.Fatal(err)
		}
		bot.logger.AddWriter(fileWriter)
	}

	if opts.UseRequestLogger {
		bot.RequestLogger = slog.CreateLogger().Level(level).Prefix("REQUESTS")
		bot.RequestLogger.AddWriter(bot.RequestLogger.CreateJsonStdoutWriter())
		if opts.WriteToFile {
			path := fmt.Sprintf("%s/requests.log", strings.TrimRight(opts.LoggerBasePath, "/"))
			fileWriter, err := bot.RequestLogger.CreateTextFileWriter(path)
			if err != nil {
				bot.logger.Fatal(err)
			}
			bot.RequestLogger.AddWriter(fileWriter)
		}
	}
}

func (bot *Bot[T]) GetUpdateOffset() int {
	bot.updateOffsetMu.Lock()
	defer bot.updateOffsetMu.Unlock()
	return bot.updateOffset
}
func (bot *Bot[T]) SetUpdateOffset(offset int) {
	bot.updateOffsetMu.Lock()
	defer bot.updateOffsetMu.Unlock()
	bot.updateOffset = offset
}
func (bot *Bot[T]) GetUpdateTypes() []tgapi.UpdateType { return bot.updateTypes }
func (bot *Bot[T]) GetLogger() *slog.Logger            { return bot.logger }
func (bot *Bot[T]) GetDBContext() *T                   { return bot.dbContext }
func (bot *Bot[T]) L10n(lang, key string) string       { return bot.l10n.Translate(lang, key) }

type DbLogger[T DbContext] func(db *T) slog.LoggerWriter

func (bot *Bot[T]) AddDatabaseLoggerWriter(writer DbLogger[T]) *Bot[T] {
	w := writer(bot.dbContext)
	bot.logger.AddWriter(w)
	if bot.RequestLogger != nil {
		bot.RequestLogger.AddWriter(w)
	}
	for _, l := range bot.extraLoggers {
		l.AddWriter(w)
	}
	return bot
}

func (bot *Bot[T]) DatabaseContext(ctx *T) *Bot[T] {
	bot.dbContext = ctx
	return bot
}
func (bot *Bot[T]) UpdateTypes(t ...tgapi.UpdateType) *Bot[T] {
	bot.updateTypes = make([]tgapi.UpdateType, 0)
	bot.updateTypes = append(bot.updateTypes, t...)
	return bot
}
func (bot *Bot[T]) AddUpdateType(t ...tgapi.UpdateType) *Bot[T] {
	bot.updateTypes = append(bot.updateTypes, t...)
	return bot
}
func (bot *Bot[T]) AddPrefixes(prefixes ...string) *Bot[T] {
	bot.prefixes = append(bot.prefixes, prefixes...)
	return bot
}
func (bot *Bot[T]) ErrorTemplate(s string) *Bot[T] {
	bot.errorTemplate = s
	return bot
}
func (bot *Bot[T]) Debug(debug bool) *Bot[T] {
	bot.debug = debug
	return bot
}
func (bot *Bot[T]) AddPlugins(plugin ...*Plugin[T]) *Bot[T] {
	for _, p := range plugin {
		bot.plugins = append(bot.plugins, *p)
		bot.logger.Debugln(fmt.Sprintf("plugins with name \"%s\" registered", p.name))
	}
	return bot
}
func (bot *Bot[T]) AddMiddleware(middleware ...Middleware[T]) *Bot[T] {
	bot.middlewares = append(bot.middlewares, middleware...)
	for _, m := range middleware {
		bot.logger.Debugln(fmt.Sprintf("middleware with name \"%s\" registered", m.name))
	}

	sort.Slice(bot.middlewares, func(i, j int) bool {
		first := bot.middlewares[i]
		second := bot.middlewares[j]
		if first.order == second.order {
			return first.name < second.name
		}
		return first.order < second.order
	})

	return bot
}
func (bot *Bot[T]) AddRunner(runner Runner[T]) *Bot[T] {
	bot.runners = append(bot.runners, runner)
	bot.logger.Debugln(fmt.Sprintf("runner with name \"%s\" registered", runner.name))
	return bot
}
func (bot *Bot[T]) AddL10n(l *L10n) *Bot[T] {
	bot.l10n = l
	return bot
}

func (bot *Bot[T]) enqueueUpdate(u *tgapi.Update) error {
	select {
	case bot.updateQueue <- u:
		return nil
	default:
		return extypes.QueueFullErr
	}
}
func (bot *Bot[T]) RunWithContext(ctx context.Context) {
	if len(bot.prefixes) == 0 {
		bot.logger.Fatalln("no prefixes defined")
		return
	}

	if len(bot.plugins) == 0 {
		bot.logger.Fatalln("no plugins defined")
		return
	}

	bot.ExecRunners()

	bot.logger.Infoln("Bot running. Press CTRL+C to exit.")
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				updates, err := bot.Updates()
				if err != nil {
					bot.logger.Errorln(err)
					continue
				}

				for _, u := range updates {
					select {
					case bot.updateQueue <- new(u):
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	pool := pond.NewPool(16)
	for update := range bot.updateQueue {
		update := update
		log.Println(update)
		pool.Submit(func() {
			bot.handle(update)
		})
	}
}
func (bot *Bot[T]) Run() {
	bot.RunWithContext(context.Background())
}
