package laniakea

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"git.nix13.pw/scuroneko/extypes"
	"git.nix13.pw/scuroneko/laniakea/tgapi"
	"git.nix13.pw/scuroneko/laniakea/utils"
	"git.nix13.pw/scuroneko/slog"
	"github.com/alitto/pond/v2"
)

// DbContext is an interface representing the application's database context.
// It is injected into plugins and middleware via Bot.DatabaseContext().
//
// Example:
//
//	type MyDB struct { ... }
//	bot := NewBot[MyDB](opts).DatabaseContext(&myDB)
//
// Use NoDB if no database is needed.
type DbContext any

// NoDB is a placeholder type for bots that do not use a database.
// Use Bot[NoDB] to indicate no dependency injection is required.
type NoDB struct{ DbContext }

// DbLogger is a function type that returns a slog.LoggerWriter for database logging.
// Used to inject database-specific log output (e.g., SQL queries, ORM events).
type DbLogger[T DbContext] func(db *T) slog.LoggerWriter

// BotPayloadType defines the serialization format for callback data payloads.
type BotPayloadType string

var (
	// BotPayloadBase64 encodes callback data as a Base64 string.
	BotPayloadBase64 BotPayloadType = "base64"
	// BotPayloadJson encodes callback data as a JSON string.
	BotPayloadJson BotPayloadType = "json"
)

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
	payloadType   BotPayloadType
	maxWorkers    int

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
	runnerOnceWG   sync.WaitGroup     // Tracks one-time async runners
	runnerBgWG     sync.WaitGroup     // Tracks background async runners
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
	limiter.SetGlobalRate(opts.RateLimit)

	apiOpts := tgapi.NewAPIOpts(opts.Token).
		SetAPIUrl(opts.APIUrl).
		UseTestServer(opts.UseTestServer).
		SetLimiter(limiter).
		SetLimiterDrop(opts.DropRLOverflow)
	api := tgapi.NewAPI(apiOpts)
	uploader := tgapi.NewUploader(api)

	prefixes := opts.Prefixes
	if len(prefixes) == 0 {
		prefixes = []string{"/"}
	}

	workers := 32
	if opts.MaxWorkers > 0 {
		workers = opts.MaxWorkers
	}

	bot := &Bot[T]{
		updateOffset:  0,
		errorTemplate: "%s",
		payloadType:   BotPayloadBase64,
		maxWorkers:    workers,
		updateQueue:   updateQueue,
		api:           api,
		uploader:      uploader,
		debug:         opts.Debug,
		prefixes:      prefixes,
		token:         opts.Token,
		plugins:       make([]Plugin[T], 0),
		updateTypes:   append([]tgapi.UpdateType{}, opts.UpdateTypes...),
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

// Close gracefully shuts down bot-owned resources.
//
// Close shuts down, in order:
//   - Registered plugins via Plugin.Close
//   - Uploader (waits for pending uploads)
//   - API client internals
//   - RequestLogger (if enabled)
//   - Main logger
//
// RunWithContext does not call Close automatically. The caller is responsible
// for invoking Close after RunWithContext returns to release these resources.
//
// Close returns a joined error containing all shutdown failures, if any.
func (bot *Bot[T]) Close() error {
	var e []error

	for _, p := range bot.plugins {
		if err := p.Close(); err != nil {
			e = append(e, err)
		}
	}
	if err := bot.uploader.Close(); err != nil {
		bot.logger.Errorln(err)
		e = append(e, err)
	}
	if err := bot.api.Close(); err != nil {
		bot.logger.Errorln(err)
		e = append(e, err)
	}
	if bot.RequestLogger != nil {
		if err := bot.RequestLogger.Close(); err != nil {
			bot.logger.Errorln(err)
			e = append(e, err)
		}
	}
	if err := bot.logger.Close(); err != nil {
		e = append(e, err)
	}
	return errors.Join(e...)
}

// CloseRemote sends Telegram Bot API "close" request for the current bot
// instance using ctx for cancellation and deadlines.
//
// This is separate from Bot.Close(), which only releases local resources.
func (bot *Bot[T]) CloseRemote(ctx context.Context) error {
	if _, err := bot.api.CloseRemoteWithContext(ctx); err != nil {
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

	bot.logger = utils.CreateLogger("BOT", level)
	if opts.WriteToFile {
		path := fmt.Sprintf("%s/main.log", strings.TrimRight(opts.LoggerBasePath, "/"))
		logger, err := utils.CreateFileLogger("BOT", level, path)
		if err != nil {
			bot.logger.Fatal(err)
		}
		bot.logger = logger
	}

	if opts.UseRequestLogger {
		bot.RequestLogger = utils.CreateLogger("REQUESTS", level)
		if opts.WriteToFile {
			path := fmt.Sprintf("%s/requests.log", strings.TrimRight(opts.LoggerBasePath, "/"))
			logger, err := utils.CreateFileLogger("REQUESTS", level, path)
			if err != nil {
				bot.logger.Fatal(err)
			}
			bot.RequestLogger = logger
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

// GetLoggerLevel returns the effective log level derived from the bot's debug
// flag.
func (bot *Bot[T]) GetLoggerLevel() slog.LogLevel {
	level := slog.FATAL
	if bot.debug {
		level = slog.DEBUG
	}
	return level
}

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

// SetPayloadType sets the payload encoding type used for callback data.
// JSON stores payload as a string: `{"cmd":"command","args":[...]}`.
// Base64 stores the same JSON encoded as a Base64URL string.
func (bot *Bot[T]) SetPayloadType(t BotPayloadType) *Bot[T] {
	bot.payloadType = t
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
// Example: "❌ Error: %s" → "❌ Error: Command not found".
func (bot *Bot[T]) ErrorTemplate(s string) *Bot[T] {
	bot.errorTemplate = s
	return bot
}

// Debug enables or disables debug logging.
func (bot *Bot[T]) Debug(debug bool) *Bot[T] {
	bot.debug = debug
	level := slog.FATAL
	if debug {
		level = slog.DEBUG
	}

	bot.logger.Level(level)
	if bot.RequestLogger != nil {
		bot.RequestLogger.Level(level)
	}
	for _, p := range bot.plugins {
		if p.logger == nil {
			continue
		}
		p.logger.Level(level)
	}
	return bot
}

// AddPlugins registers one or more plugins.
// Plugins are executed in registration order unless filtered by middleware.
//
// Registration is a commit point for plugin configuration. The Bot stores
// plugin metadata internally, so plugins must be fully configured before they
// are passed here. Post-registration mutation through the original *Plugin is
// not a supported API, even if some changes appear to work due to shared maps.
func (bot *Bot[T]) AddPlugins(plugin ...*Plugin[T]) *Bot[T] {
	level := bot.GetLoggerLevel()
	for _, p := range plugin {
		if p.logger == nil {
			logger := utils.CreateLogger(p.name, level)
			p.SetLogger(logger)
		}
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
		return bot
	}
	bot.l10n = l
	return bot
}

// AddDatabaseLoggerWriter adds a database logger writer to all loggers.
//
// The writer will receive logs from:
//   - Main bot logger
//   - Request logger (if enabled)
//   - API and Uploader loggers
//   - Already registered plugin loggers
//
// Call this after AddPlugins if plugin loggers should also receive the writer.
// Plugins registered later do not automatically inherit previously added
// database writers; call AddDatabaseLoggerWriter again after adding them.
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
	for _, p := range bot.plugins {
		if p.logger != nil {
			p.logger.AddWriter(w)
		}
	}
	return bot
}

// RunWithContext starts the bot with a given context for graceful shutdown.
//
// This is the main entry point for bot execution. It:
//   - Validates required configuration (prefixes, plugins)
//   - Starts all registered runners as background goroutines
//   - Begins polling for updates via Telegram's GetUpdates API
//   - Processes updates concurrently using a worker pool with size configurable via BotOpts.MaxWorkers
//
// The context controls graceful shutdown. When canceled, the bot:
//   - Stops polling for new updates
//   - Finishes processing currently queued updates
//   - Waits for registered runners to exit
//
// RunWithContext does not close API, uploader, or logger resources on return.
// The caller must invoke Close after RunWithContext finishes.
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	go bot.RunWithContext(ctx)
//	// ... later ...
//	cancel() // triggers graceful shutdown
//	_ = bot.Close(context.Background())
func (bot *Bot[T]) RunWithContext(ctx context.Context) {
	if len(bot.prefixes) == 0 {
		bot.logger.Fatalln("no prefixes defined")
		return
	}

	if len(bot.plugins) == 0 {
		bot.logger.Fatalln("no plugins defined")
		return
	}

	bot.ExecRunners(ctx)

	bot.logger.Infoln("Bot running. Press CTRL+C to exit.")

	// Start update polling in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				bot.logger.Errorln(fmt.Sprintf("panic in update polling: %v", r))
			}
			close(bot.updateQueue)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				updates, err := bot.Updates(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					bot.logger.Errorln("failed to fetch updates:", err)
					continue
				}

				for _, update := range updates {
					u := update // copy loop variable to avoid race condition
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
	pool := pond.NewPool(bot.maxWorkers)
	for update := range bot.updateQueue {
		u := update // capture loop variable
		pool.Submit(func() {
			bot.handle(u)
		})
	}
	pool.Stop() // Wait for all tasks to complete and stop the pool
	bot.runnerOnceWG.Wait()
	bot.runnerBgWG.Wait()
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
