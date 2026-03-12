// Package laniakea provides a modular, extensible framework for building scalable
// Telegram bots with support for plugins, middleware, localization, draft messages,
// rate limiting, structured logging, and dependency injection.
//
// The framework is designed around a fluent API for configuration and separation of concerns:
//
//   - Plugins: Handle specific commands or events (e.g., /start, /help)
//   - Middleware: Intercept and modify updates before plugins run (auth, logging, validation)
//   - Runners: Background goroutines for cleanup, cron jobs, or monitoring
//   - DraftProvider: Safely build and resume multi-step messages
//   - L10n: Multi-language support via key-based translation
//   - RateLimiter: Enforces Telegram API limits to avoid bans
//   - Structured Logging: JSON stdout + optional file output with request-level tracing
//   - Dependency Injection: Inject custom database contexts (e.g., *gorm.DB, *sql.DB)
//
// Example usage:
//
//	bot := laniakea.NewBot[mydb.DBContext](laniakea.LoadOptsFromEnv()).
//		DatabaseContext(&myDB).
//		AddUpdateType(tgapi.UpdateTypeMessage).
//		AddPrefixes("/", "!").
//		AddPlugins(&startPlugin, &helpPlugin).
//		AddMiddleware(&authMiddleware, &logMiddleware).
//		AddRunner(&cleanupRunner).
//		AddL10n(l10n.New())
//
//	go bot.Run()
//	<-ctx.Done() // wait for shutdown signal
//
// All methods are thread-safe except direct field access. Use provided accessors
// (e.g., GetDBContext, SetUpdateOffset) for safe concurrent access.
package laniakea

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.nix13.pw/scuroneko/extypes"
	"git.nix13.pw/scuroneko/laniakea/tgapi"
	"git.nix13.pw/scuroneko/laniakea/utils"
	"git.nix13.pw/scuroneko/slog"
	"github.com/alitto/pond/v2"
)

// BotOpts holds configuration options for initializing a Bot.
//
// Values are loaded from environment variables via LoadOptsFromEnv().
// Use NewOpts() to create a zero-value struct and set fields manually.
type BotOpts struct {
	// Token is the Telegram bot token (required).
	Token string

	// UpdateTypes is a semicolon-separated list of update types to listen for.
	// Example: "message;edited_message;callback_query"
	// Defaults to empty (Telegram will return all types).
	UpdateTypes []string

	// Debug enables debug-level logging.
	Debug bool

	// ErrorTemplate is the format string used to wrap error messages sent to users.
	// Use "%s" to insert the actual error. Example: "❌ Error: %s"
	ErrorTemplate string

	// Prefixes is a list of command prefixes (e.g., ["/", "!"]).
	// Defaults to ["/"] if not set via environment.
	Prefixes []string

	// LoggerBasePath is the directory where log files are written.
	// Defaults to "./".
	LoggerBasePath string

	// UseRequestLogger enables detailed logging of all Telegram API requests.
	UseRequestLogger bool

	// WriteToFile enables writing logs to files (main.log and requests.log).
	WriteToFile bool

	// UseTestServer uses Telegram's test server (https://api.test.telegram.org).
	UseTestServer bool

	// APIUrl overrides the default Telegram API endpoint (useful for proxies or self-hosted).
	APIUrl string

	// RateLimit is the maximum number of API requests per second.
	// Telegram allows up to 30 req/s for most bots. Defaults to 30.
	RateLimit int

	// DropRLOverflow drops incoming updates when rate limit is exceeded instead of queuing.
	// Use this to prioritize responsiveness over reliability.
	DropRLOverflow bool
}

// NewOpts returns a new BotOpts with zero values.
func NewOpts() *BotOpts { return new(BotOpts) }

// LoadOptsFromEnv loads BotOpts from environment variables.
//
// Environment variables:
//   - TG_TOKEN: Bot token (required)
//   - UPDATE_TYPES: semicolon-separated update types (e.g., "message;callback_query")
//   - DEBUG: "true" to enable debug logging
//   - ERROR_TEMPLATE: format string for error messages (e.g., "❌ %s")
//   - PREFIXES: semicolon-separated prefixes (e.g., "/;!bot")
//   - LOGGER_BASE_PATH: directory for log files (default: "./")
//   - USE_REQ_LOG: "true" to enable request logging
//   - WRITE_TO_FILE: "true" to write logs to files
//   - USE_TEST_SERVER: "true" to use Telegram test server
//   - API_URL: custom API endpoint
//   - RATE_LIMIT: max requests per second (default: 30)
//   - DROP_RL_OVERFLOW: "true" to drop updates on rate limit overflow
//
// Returns a populated BotOpts. If TG_TOKEN is missing, behavior is undefined.
func LoadOptsFromEnv() *BotOpts {
	rateLimit := 30
	if rl := os.Getenv("RATE_LIMIT"); rl != "" {
		if n, err := strconv.Atoi(rl); err == nil {
			rateLimit = n
		}
	}

	return &BotOpts{
		Token:       os.Getenv("TG_TOKEN"),
		UpdateTypes: strings.Split(os.Getenv("UPDATE_TYPES"), ";"),

		Debug:         os.Getenv("DEBUG") == "true",
		ErrorTemplate: os.Getenv("ERROR_TEMPLATE"),
		Prefixes:      LoadPrefixesFromEnv(),

		LoggerBasePath:   os.Getenv("LOGGER_BASE_PATH"),
		UseRequestLogger: os.Getenv("USE_REQ_LOG") == "true",
		WriteToFile:      os.Getenv("WRITE_TO_FILE") == "true",

		UseTestServer: os.Getenv("USE_TEST_SERVER") == "true",
		APIUrl:        os.Getenv("API_URL"),

		RateLimit:      rateLimit,
		DropRLOverflow: os.Getenv("DROP_RL_OVERFLOW") == "true",
	}
}

// LoadPrefixesFromEnv returns the PREFIXES environment variable split by semicolon.
// Defaults to ["/"] if not set.
func LoadPrefixesFromEnv() []string {
	prefixesS, exists := os.LookupEnv("PREFIXES")
	if !exists {
		return []string{"/"}
	}
	return strings.Split(prefixesS, ";")
}

// DbContext is an interface representing the application's database context.
// It is injected into plugins and middleware via Bot.DatabaseContext().
//
// Example:
//
//	type MyDB struct { ... }
//	bot := NewBot[MyDB](opts).DatabaseContext(&myDB)
//
// Use NoDB if no database is needed.
type DbContext interface{}

// NoDB is a placeholder type for bots that do not use a database.
// Use Bot[NoDB] to indicate no dependency injection is required.
type NoDB struct{ DbContext }

// Bot is the core Telegram bot instance.
//
// Manages:
//   - API communication via tgapi
//   - Update processing pipeline (middleware → plugins)
//   - Background runners
//   - Logging and rate limiting
//   - Localization and draft message support
//
// All methods are safe for concurrent use. Direct field access is not recommended.
type Bot[T DbContext] struct {
	token         string
	debug         bool
	errorTemplate string
	username      string

	logger        *slog.Logger                // Main bot logger (JSON stdout + optional file)
	RequestLogger *slog.Logger                // Optional request-level API logging
	extraLoggers  extypes.Slice[*slog.Logger] // API, Uploader, and custom loggers

	plugins     []Plugin[T]     // Command/event handlers
	middlewares []Middleware[T] // Pre-processing filters (sorted by order)
	prefixes    []string        // Command prefixes (e.g., "/", "!")
	runners     []Runner[T]     // Background tasks (e.g., cleanup, cron)

	api           *tgapi.API      // Telegram API client
	uploader      *tgapi.Uploader // File uploader
	dbContext     *T              // Injected database context
	l10n          *L10n           // Localization manager
	draftProvider *DraftProvider  // Draft message builder

	updateOffsetMu sync.Mutex
	updateOffset   int                // Last processed update ID
	updateTypes    []tgapi.UpdateType // Types of updates to fetch
	updateQueue    chan *tgapi.Update // Internal queue for processing updates
}

// NewBot creates and initializes a new Bot instance using the provided BotOpts.
//
// Automatically:
//   - Creates API and Uploader clients
//   - Initializes structured logging (JSON stdout + optional file)
//   - Fetches bot username via GetMe()
//   - Sets up DraftProvider with random IDs
//   - Adds API and Uploader loggers to extraLoggers
//
// Panics if:
//   - Token is empty
//   - GetMe() fails (invalid token or network error)
func NewBot[T any](opts *BotOpts) *Bot[T] {
	if opts.Token == "" {
		panic("laniakea: BotOpts.Token is required")
	}

	updateQueue := make(chan *tgapi.Update, 512)

	//var limiter *utils.RateLimiter
	//if opts.RateLimit > 0 {
	//	limiter = utils.NewRateLimiter()
	//}
	limiter := utils.NewRateLimiter()

	apiOpts := tgapi.NewAPIOpts(opts.Token).
		SetAPIUrl(opts.APIUrl).
		UseTestServer(opts.UseTestServer).
		SetLimiter(limiter)
	api := tgapi.NewAPI(apiOpts)
	uploader := tgapi.NewUploader(api)

	prefixes := opts.Prefixes
	if len(prefixes) == 0 {
		prefixes = []string{"/"}
	}

	bot := &Bot[T]{
		updateOffset:  0,
		errorTemplate: "%s",
		updateQueue:   updateQueue,
		api:           api,
		uploader:      uploader,
		debug:         opts.Debug,
		prefixes:      prefixes,
		token:         opts.Token,
		plugins:       make([]Plugin[T], 0),
		updateTypes:   make([]tgapi.UpdateType, 0),
		runners:       make([]Runner[T], 0),
		extraLoggers:  make([]*slog.Logger, 0),
		l10n:          &L10n{},
		draftProvider: NewRandomDraftProvider(api),
	}

	// Add API and Uploader loggers to extraLoggers for unified output
	bot.extraLoggers = bot.extraLoggers.Push(api.GetLogger()).Push(uploader.GetLogger())

	if len(opts.ErrorTemplate) > 0 {
		bot.errorTemplate = opts.ErrorTemplate
	}
	if len(opts.LoggerBasePath) == 0 {
		opts.LoggerBasePath = "./"
	}
	bot.initLoggers(opts)

	// Fetch bot info to validate token and get username
	u, err := api.GetMe()
	if err != nil {
		_ = bot.Close()
		bot.logger.Fatal(err)
	}
	bot.username = Val(u.Username, "")
	if bot.username == "" {
		bot.logger.Warn("Can't get bot username. Named command handlers won't work!")
	}
	bot.logger.Infoln(fmt.Sprintf("Authorized as %s (@%s)", u.FirstName, Val(u.Username, "unknown")))

	return bot
}

// Close gracefully shuts down the bot.
//
// Closes:
//   - Uploader (waits for pending uploads)
//   - API client
//   - RequestLogger (if enabled)
//   - Main logger
//
// Returns the first error encountered, if any.
func (bot *Bot[T]) Close() error {
	if err := bot.uploader.Close(); err != nil {
		bot.logger.Errorln(err)
	}
	if err := bot.api.CloseApi(); err != nil {
		bot.logger.Errorln(err)
	}
	if bot.RequestLogger != nil {
		if err := bot.RequestLogger.Close(); err != nil {
			bot.logger.Errorln(err)
		}
	}
	if err := bot.logger.Close(); err != nil {
		return err
	}
	return nil
}

// initLoggers configures the main and optional request loggers.
//
// Uses DEBUG flag to set log level (DEBUG if true, FATAL otherwise).
// Writes to stdout in JSON format by default.
// If WriteToFile is true, writes to main.log and requests.log in LoggerBasePath.
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

// GetUpdateOffset returns the current update offset (thread-safe).
func (bot *Bot[T]) GetUpdateOffset() int {
	bot.updateOffsetMu.Lock()
	defer bot.updateOffsetMu.Unlock()
	return bot.updateOffset
}

// SetUpdateOffset sets the update offset for next GetUpdates call (thread-safe).
func (bot *Bot[T]) SetUpdateOffset(offset int) {
	bot.updateOffsetMu.Lock()
	defer bot.updateOffsetMu.Unlock()
	bot.updateOffset = offset
}

// GetUpdateTypes returns the list of update types the bot is configured to receive.
func (bot *Bot[T]) GetUpdateTypes() []tgapi.UpdateType { return bot.updateTypes }

// GetLogger returns the main bot logger.
func (bot *Bot[T]) GetLogger() *slog.Logger { return bot.logger }

// GetDBContext returns the injected database context.
// Returns nil if not set via DatabaseContext().
func (bot *Bot[T]) GetDBContext() *T { return bot.dbContext }

// L10n translates a key in the given language.
// Returns empty string if translation not found.
func (bot *Bot[T]) L10n(lang, key string) string {
	return bot.l10n.Translate(lang, key)
}

// SetDraftProvider replaces the default DraftProvider with a custom one.
// Useful for using LinearDraftIdGenerator to persist draft IDs across restarts.
func (bot *Bot[T]) SetDraftProvider(p *DraftProvider) *Bot[T] {
	bot.draftProvider = p
	return bot
}

// DbLogger is a function type that returns a slog.LoggerWriter for database logging.
// Used to inject database-specific log output (e.g., SQL queries, ORM events).
type DbLogger[T DbContext] func(db *T) slog.LoggerWriter

// AddDatabaseLoggerWriter adds a database logger writer to all loggers.
//
// The writer will receive logs from:
//   - Main bot logger
//   - Request logger (if enabled)
//   - API and Uploader loggers
//
// Example:
//
//	bot.AddDatabaseLoggerWriter(func(db *MyDB) slog.LoggerWriter {
//	    return db.QueryLogger()
//	})
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

// DatabaseContext injects a database context into the bot.
// This context is accessible to plugins and middleware via GetDBContext().
func (bot *Bot[T]) DatabaseContext(ctx *T) *Bot[T] {
	bot.dbContext = ctx
	return bot
}

// UpdateTypes sets the list of update types the bot will request from Telegram.
// Overwrites any previously set types.
func (bot *Bot[T]) UpdateTypes(t ...tgapi.UpdateType) *Bot[T] {
	bot.updateTypes = make([]tgapi.UpdateType, 0)
	bot.updateTypes = append(bot.updateTypes, t...)
	return bot
}

// AddUpdateType adds one or more update types to the list.
// Does not overwrite existing types.
func (bot *Bot[T]) AddUpdateType(t ...tgapi.UpdateType) *Bot[T] {
	bot.updateTypes = append(bot.updateTypes, t...)
	return bot
}

// AddPrefixes adds one or more command prefixes (e.g., "/", "!").
// Must have at least one prefix before Run().
func (bot *Bot[T]) AddPrefixes(prefixes ...string) *Bot[T] {
	bot.prefixes = append(bot.prefixes, prefixes...)
	return bot
}

// ErrorTemplate sets the format string for error messages sent to users.
// Use "%s" to insert the error message.
// Example: "❌ Error: %s" → "❌ Error: Command not found"
func (bot *Bot[T]) ErrorTemplate(s string) *Bot[T] {
	bot.errorTemplate = s
	return bot
}

// Debug enables or disables debug logging.
func (bot *Bot[T]) Debug(debug bool) *Bot[T] {
	bot.debug = debug
	return bot
}

// AddPlugins registers one or more plugins.
// Plugins are executed in registration order unless filtered by middleware.
func (bot *Bot[T]) AddPlugins(plugin ...*Plugin[T]) *Bot[T] {
	for _, p := range plugin {
		bot.plugins = append(bot.plugins, *p)
		bot.logger.Debugln(fmt.Sprintf("plugins with name \"%s\" registered", p.name))
	}
	return bot
}

// AddMiddleware registers one or more middleware handlers.
//
// Middleware are executed in order of increasing .order value before plugins.
// If two middleware have the same order, they are sorted lexicographically by name.
//
// Middleware can:
//   - Modify or reject updates before they reach plugins
//   - Inject context (e.g., user auth state, rate limit status)
//   - Log, validate, or transform incoming data
//
// Example:
//
//	bot.AddMiddleware(&authMiddleware, &rateLimitMiddleware)
//
// Panics if any middleware has a nil name.
func (bot *Bot[T]) AddMiddleware(middleware ...Middleware[T]) *Bot[T] {
	for _, m := range middleware {
		if m.name == "" {
			panic("laniakea: middleware must have a non-empty name")
		}
		bot.middlewares = append(bot.middlewares, m)
		bot.logger.Debugln(fmt.Sprintf("middleware with name \"%s\" registered", m.name))
	}

	// Stable sort by order (ascending), then by name (lexicographic)
	sort.Slice(bot.middlewares, func(i, j int) bool {
		first := bot.middlewares[i]
		second := bot.middlewares[j]
		if first.order != second.order {
			return first.order < second.order
		}
		return first.name < second.name
	})

	return bot
}

// AddRunner registers a background runner to execute concurrently with the bot.
//
// Runners are goroutines that run independently of update processing.
// Common use cases:
//   - Periodic cleanup (e.g., expiring drafts, clearing temp files)
//   - Metrics collection or health checks
//   - Scheduled tasks (e.g., daily announcements)
//
// Runners are started immediately after Bot.Run() is called.
//
// Example:
//
//	bot.AddRunner(&cleanupRunner)
//
// Panics if runner has a nil name.
func (bot *Bot[T]) AddRunner(runner Runner[T]) *Bot[T] {
	if runner.name == "" {
		panic("laniakea: runner must have a non-empty name")
	}
	bot.runners = append(bot.runners, runner)
	bot.logger.Debugln(fmt.Sprintf("runner with name \"%s\" registered", runner.name))
	return bot
}

// AddL10n sets the localization (i18n) provider for the bot.
//
// The L10n instance must be pre-populated with translations.
// Translations are accessed via Bot.L10n(lang, key).
//
// Example:
//
//	l10n := l10n.New()
//	l10n.Add("en", "hello", "Hello!")
//	l10n.Add("es", "hello", "¡Hola!")
//	bot.AddL10n(l10n)
//
// Replaces any previously set L10n instance.
func (bot *Bot[T]) AddL10n(l *L10n) *Bot[T] {
	if l == nil {
		bot.logger.Warn("AddL10n called with nil L10n; localization will be disabled")
	}
	bot.l10n = l
	return bot
}

// enqueueUpdate attempts to add an update to the internal processing queue.
//
// Returns extypes.QueueFullErr if the queue is full and the update cannot be enqueued.
// This is non-blocking and used to implement rate-limiting behavior.
//
// When DropRLOverflow is enabled, this error is ignored and the update is dropped.
// Otherwise, the update is retried via the main update loop.
func (bot *Bot[T]) enqueueUpdate(u *tgapi.Update) error {
	select {
	case bot.updateQueue <- u:
		return nil
	default:
		return extypes.QueueFullErr
	}
}

// RunWithContext starts the bot with a given context for graceful shutdown.
//
// This is the main entry point for bot execution. It:
//   - Validates required configuration (prefixes, plugins)
//   - Starts all registered runners as background goroutines
//   - Begins polling for updates via Telegram's GetUpdates API
//   - Processes updates concurrently using a worker pool (16 goroutines)
//
// The context controls graceful shutdown. When canceled, the bot:
//   - Stops polling for new updates
//   - Finishes processing currently queued updates
//   - Closes all resources (API, uploader, loggers)
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	go bot.RunWithContext(ctx)
//	// ... later ...
//	cancel() // triggers graceful shutdown
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

	// Start update polling in a goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				updates, err := bot.Updates()
				if err != nil {
					bot.logger.Errorln("failed to fetch updates:", err)
					time.Sleep(2 * time.Second) // exponential backoff
					continue
				}

				for _, u := range updates {
					select {
					case bot.updateQueue <- &u:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	// Start worker pool for concurrent update handling
	pool := pond.NewPool(16)
	for update := range bot.updateQueue {
		update := update // capture loop variable
		pool.Submit(func() {
			bot.handle(update)
		})
	}
}

// Run starts the bot using a background context.
//
// Equivalent to RunWithContext(context.Background()).
// Use this for simple bots where graceful shutdown is not required.
//
// For production use, prefer RunWithContext to handle SIGINT/SIGTERM gracefully.
func (bot *Bot[T]) Run() {
	bot.RunWithContext(context.Background())
}
