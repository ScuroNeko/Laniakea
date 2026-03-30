package laniakea

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"git.scuroneko.dev/scuroneko/extypes"
	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/laniakea/utils"
	"git.scuroneko.dev/scuroneko/slog"
	"github.com/alitto/pond/v2"
)

// AppData is the generic shared application data type injected into bots,
// plugins, and handlers.
//
// Use it for long-lived shared dependencies such as database handles, service
// containers, API clients, or immutable configuration snapshots.
//
// Example:
//
//	type MyDB struct { ... }
//	myDB := &MyDB{}
//	bot, err := NewBot[*MyDB](opts)
//	if err != nil {
//		return err
//	}
//	bot.SetAppData(myDB)
//
// Use NoData if no shared application data is needed.
type AppData any

// NoData is a placeholder type for bots that do not use shared application
// data.
//
// Use Bot[NoData] to indicate no shared dependency injection is required.
type NoData struct{ AppData }

// AppDataLogger builds a slog.LoggerWriter from injected application data.
//
// Use it when shared application data exposes a log sink or adapter that should
// receive framework logs.
type AppDataLogger[T AppData] func(data T) slog.LoggerWriter

// BotPayloadType defines the serialization format for callback data payloads.
type BotPayloadType string

var (
	// BotPayloadBase64 encodes callback data as a Base64 string.
	BotPayloadBase64 BotPayloadType = "base64"
	// BotPayloadJson encodes callback data as a JSON string.
	BotPayloadJson BotPayloadType = "json"
)

var (
	// ErrNoPrefixes reports that the bot was started without any command prefixes.
	ErrNoPrefixes = errors.New("no prefixes defined")
	// ErrNoPlugins reports that the bot was started without any registered plugins.
	ErrNoPlugins = errors.New("no plugins defined")
	// ErrBotAlreadyRun reports that Run or RunWithContext was called more than once.
	ErrBotAlreadyRun = errors.New("bot can only be run once")

	// ErrTokenRequired reports that BotOpts.Token was empty.
	ErrTokenRequired = errors.New("token required")
	// ErrOptsIsNil reports that NewBot was called with a nil BotOpts pointer.
	ErrOptsIsNil = errors.New("opts is nil")
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
// Runtime accessors are safe for concurrent use. Configure the bot before Run.
// A Bot is single-use: after Run or RunWithContext returns, create a new Bot for the next session.
type Bot[T AppData] struct {
	token             string
	debug             bool
	errorTemplate     string
	username          string
	payloadType       BotPayloadType
	strictPayloadType bool
	maxWorkers        int

	logger        *slog.Logger                // Main bot logger (JSON stdout + optional file)
	RequestLogger *slog.Logger                // Optional request-level API logging
	extraLoggers  extypes.Slice[*slog.Logger] // API, Uploader, and custom loggers

	plugins     []Plugin[T]     // Command/event handlers
	middlewares []Middleware[T] // Pre-processing filters (sorted by order)
	prefixes    []string        // Command prefixes (e.g., "/", "!")
	runners     []Runner[T]     // Background tasks (e.g., cleanup, cron)

	api           *tgapi.API      // Telegram API client
	uploader      *tgapi.Uploader // File uploader
	l10n          *L10n           // Localization manager
	draftProvider *DraftProvider  // Draft message builder

	appData         T // Injected application data
	hasAppData      bool
	warnedValueData bool

	sessionStore       SessionStore // Session store for scene management
	sceneScopePriority []SceneScope

	updateOffsetMu sync.Mutex
	updateOffset   int                // Last processed update ID
	updateTypes    []tgapi.UpdateType // Types of updates to fetch
	updateQueue    chan *tgapi.Update // Internal queue for processing updates
	runnerOnceWG   sync.WaitGroup     // Tracks one-time async runners
	runnerBgWG     sync.WaitGroup     // Tracks background async runners
	runStateMu     sync.Mutex
	running        bool
	ran            bool
}

func (bot *Bot[T]) configMutable(method string) bool {
	bot.runStateMu.Lock()
	defer bot.runStateMu.Unlock()
	if !bot.ran {
		return true
	}
	if bot.logger != nil {
		bot.logger.Warnln(fmt.Sprintf("%s called after bot configuration was frozen; ignoring", method))
	}
	return false
}

// NewBot creates and initializes a new Bot instance using the provided BotOpts.
//
// Automatically:
//   - Creates API and Uploader clients
//   - Initializes structured logging (JSON stdout + optional file)
//   - Fetches bot username via GetMe()
//   - Sets up DraftProvider with random IDs
//   - Adds API and Uploader loggers to extraLoggers
func NewBot[T any](opts *BotOpts) (*Bot[T], error) {
	if opts == nil {
		return nil, ErrOptsIsNil
	}
	if opts.Token == "" {
		return nil, ErrTokenRequired
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
		updateOffset:      0,
		errorTemplate:     "%s",
		payloadType:       BotPayloadBase64,
		strictPayloadType: opts.StrictPayloadType,
		maxWorkers:        workers,
		updateQueue:       updateQueue,
		api:               api,
		uploader:          uploader,
		debug:             opts.Debug,
		prefixes:          prefixes,
		token:             opts.Token,
		plugins:           make([]Plugin[T], 0),
		updateTypes:       append([]tgapi.UpdateType{}, opts.UpdateTypes...),
		runners:           make([]Runner[T], 0),
		extraLoggers:      make([]*slog.Logger, 0),
		l10n:              &L10n{},
		draftProvider:     NewRandomDraftProvider(api),

		sessionStore:       NewMemorySessionStore(),
		sceneScopePriority: []SceneScope{SceneScopeUserChat, SceneScopeChat, SceneScopeUser},
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
		return nil, err
	}
	bot.username = Val(u.Username, "")
	if bot.username == "" {
		bot.logger.Warn("Can't get bot username. Named command handlers won't work!")
	}
	bot.logger.Infoln(fmt.Sprintf("Authorized as %s (@%s)", u.FirstName, Val(u.Username, "unknown")))

	return bot, nil
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

// Internal logger setup for the bot and optional request logger.
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
			bot.logger.Errorln(err)
		} else {
			bot.logger = logger
		}
	}

	if opts.UseRequestLogger {
		bot.RequestLogger = utils.CreateLogger("REQUESTS", level)
		if opts.WriteToFile {
			path := fmt.Sprintf("%s/requests.log", strings.TrimRight(opts.LoggerBasePath, "/"))
			logger, err := utils.CreateFileLogger("REQUESTS", level, path)
			if err != nil {
				bot.logger.Errorln(err)
			} else {
				bot.RequestLogger = logger
			}
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
func (bot *Bot[T]) GetUpdateTypes() []tgapi.UpdateType {
	return append([]tgapi.UpdateType(nil), bot.updateTypes...)
}

// GetLogger returns the main bot logger.
func (bot *Bot[T]) GetLogger() *slog.Logger { return bot.logger }

// GetAppData returns the injected application data.
// If SetAppData was not called, it returns the zero value of T.
func (bot *Bot[T]) GetAppData() T { return bot.appData }

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
	if !bot.configMutable("SetDraftProvider") {
		return bot
	}
	bot.draftProvider = p
	return bot
}

// GetDraftProvider returns the draft provider currently used by the bot.
func (bot *Bot[T]) GetDraftProvider() *DraftProvider {
	return bot.draftProvider
}

// SetSessionStore replaces the session store used for scene management.
func (bot *Bot[T]) SetSessionStore(store SessionStore) *Bot[T] {
	if !bot.configMutable("SetSessionStore") {
		return bot
	}
	if store == nil {
		bot.logger.Warn("SetSessionStore called with nil store; using default MemorySessionStore")
		return bot
	}
	bot.sessionStore = store
	return bot
}

// GetSessionStore returns the session store used for scene management.
func (bot *Bot[T]) GetSessionStore() SessionStore {
	return bot.sessionStore
}

// SetSceneScopePriority sets the lookup order for resolving active scene sessions.
func (bot *Bot[T]) SetSceneScopePriority(priority []SceneScope) *Bot[T] {
	if !bot.configMutable("SetSceneScopePriority") {
		return bot
	}
	newPriority := make([]SceneScope, 0, 3)
	for _, scope := range priority {
		if scope != SceneScopeUser && scope != SceneScopeChat && scope != SceneScopeUserChat {
			bot.logger.Warnln(fmt.Sprintf("invalid scene scope %v in priority list; ignoring", scope))
			continue
		}
		if slices.Index(newPriority, scope) >= 0 {
			bot.logger.Warnln(fmt.Sprintf("duplicate scope %v in scene scope priority; ignoring duplicates", scope))
			continue
		}
		newPriority = append(newPriority, scope)
	}
	if len(newPriority) == 0 || len(newPriority) > 3 {
		bot.logger.Warnln("scene scope priority must have 1 to 3 scopes; ignoring invalid input")
		return bot
	}
	bot.sceneScopePriority = append([]SceneScope(nil), newPriority...)
	return bot
}

// SetAppData injects shared application data into the bot.
//
// The data is accessible to commands, payload handlers, middleware, scenes,
// and runners through the generic type parameter T.
//
// For shared dependencies such as *sql.DB, prefer using a pointer type as T.
// Value-typed application data is supported, but the bot warns once because
// handlers receive T by value.
func (bot *Bot[T]) SetAppData(ctx T) *Bot[T] {
	if !bot.configMutable("SetAppData") {
		return bot
	}
	if !bot.warnedValueData && shouldWarnOnValueAppData[T]() && bot.logger != nil {
		bot.logger.Warnln("app data uses a value type; shared dependencies should usually use a pointer type as T")
		bot.warnedValueData = true
	}
	bot.appData = ctx
	bot.hasAppData = true
	return bot
}

// SetUpdateTypes sets the list of update types the bot will request from Telegram.
// Overwrites any previously set types.
func (bot *Bot[T]) SetUpdateTypes(t ...tgapi.UpdateType) *Bot[T] {
	if !bot.configMutable("UpdateTypes") {
		return bot
	}
	bot.updateTypes = make([]tgapi.UpdateType, 0)
	bot.updateTypes = append(bot.updateTypes, t...)
	return bot
}

// SetPayloadType sets the default payload encoding type used for callback data.
// JSON stores payload as a string: `{"cmd":"command","args":[...]}`.
// Base64 stores the same JSON encoded as a Base64URL string.
// InlineKeyboard.SetPayloadType may override this value for an individual keyboard.
func (bot *Bot[T]) SetPayloadType(t BotPayloadType) *Bot[T] {
	if !bot.configMutable("SetPayloadType") {
		return bot
	}
	bot.payloadType = t
	return bot
}

// GetPayloadType returns the bot's default callback payload encoding type.
func (bot *Bot[T]) GetPayloadType() BotPayloadType { return bot.payloadType }

// SetStrictPayloadType enables or disables strict callback payload decoding.
// When enabled, callback payloads must match the bot's default payload type.
func (bot *Bot[T]) SetStrictPayloadType(strict bool) *Bot[T] {
	if !bot.configMutable("SetStrictPayloadType") {
		return bot
	}
	bot.strictPayloadType = strict
	return bot
}

// AddUpdateType adds one or more update types to the list.
// Does not overwrite existing types.
func (bot *Bot[T]) AddUpdateType(t ...tgapi.UpdateType) *Bot[T] {
	if !bot.configMutable("AddUpdateType") {
		return bot
	}
	bot.updateTypes = append(bot.updateTypes, t...)
	return bot
}

// AddPrefixes adds one or more command prefixes (e.g., "/", "!").
// Must have at least one prefix before Run().
func (bot *Bot[T]) AddPrefixes(prefixes ...string) *Bot[T] {
	if !bot.configMutable("AddPrefixes") {
		return bot
	}
	bot.prefixes = append(bot.prefixes, prefixes...)
	return bot
}

// SetErrorTemplate sets the format string for error messages sent to users.
// Use "%s" to insert the error message.
// Example: "❌ Error: %s" → "❌ Error: Command not found".
func (bot *Bot[T]) SetErrorTemplate(s string) *Bot[T] {
	if !bot.configMutable("ErrorTemplate") {
		return bot
	}
	bot.errorTemplate = s
	return bot
}

// SetDebug enables or disables debug logging.
func (bot *Bot[T]) SetDebug(debug bool) *Bot[T] {
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
	if !bot.configMutable("AddPlugins") {
		return bot
	}
	level := bot.GetLoggerLevel()
	for _, p := range plugin {
		if p == nil {
			if bot.logger != nil {
				bot.logger.Warn("nil plugin skipped")
			}
			continue
		}
		cloned := clonePlugin(p)
		if cloned.logger == nil {
			cloned.logger = utils.CreateLogger(cloned.name, level)
		}
		bot.plugins = append(bot.plugins, cloned)
		if bot.logger != nil {
			bot.logger.Debugln(fmt.Sprintf("plugins with name \"%s\" registered", cloned.name))
		}
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
//	bot.AddMiddleware(authMiddleware, rateLimitMiddleware)
//
// Middleware with an empty name are skipped with a warning.
func (bot *Bot[T]) AddMiddleware(middleware ...Middleware[T]) *Bot[T] {
	if !bot.configMutable("AddMiddleware") {
		return bot
	}
	for _, m := range middleware {
		if m.name == "" {
			bot.logger.Warnln("middleware must have a non-empty name")
			continue
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
//	bot.AddRunner(cleanupRunner)
//
// Runners with an empty name are skipped with a warning.
func (bot *Bot[T]) AddRunner(runner Runner[T]) *Bot[T] {
	if !bot.configMutable("AddRunner") {
		return bot
	}
	if runner.name == "" {
		bot.logger.Warnln("runner must have a non-empty name")
		return bot
	}
	bot.runners = append(bot.runners, runner)
	bot.logger.Debugln(fmt.Sprintf("runner with name \"%s\" registered", runner.name))
	return bot
}

// SetL10n sets the localization (i18n) provider for the bot.
//
// The L10n instance must be pre-populated with translations.
// Translations are accessed via Bot.L10n(lang, key).
//
// Example:
//
//	l10n := l10n.New()
//	l10n.Add("en", "hello", "Hello!")
//	l10n.Add("es", "hello", "¡Hola!")
//	bot.SetL10n(l10n)
//
// Replaces any previously set L10n instance.
func (bot *Bot[T]) SetL10n(l *L10n) *Bot[T] {
	if !bot.configMutable("SetL10n") {
		return bot
	}
	if l == nil {
		bot.logger.Warn("SetL10n called with nil L10n; localization will be disabled")
		return bot
	}
	bot.l10n = l
	return bot
}

// AddAppDataLoggerWriter adds an app-data-backed logger writer to all loggers.
//
// The writer will receive logs from:
//   - Main bot logger
//   - Request logger (if enabled)
//   - API and Uploader loggers
//   - Already registered plugin loggers
//
// Call this after AddPlugins if plugin loggers should also receive the writer.
// Plugins registered later do not automatically inherit previously added
// writers; call AddAppDataLoggerWriter again after adding them.
//
// Example:
//
//	bot.AddAppDataLoggerWriter(func(data *MyAppData) slog.LoggerWriter {
//	    return data.QueryLogger()
//	})
func (bot *Bot[T]) AddAppDataLoggerWriter(writer AppDataLogger[T]) *Bot[T] {
	if !bot.hasAppData {
		bot.logger.Warnln("app data is not set; skipping app-data logger writer")
		return bot
	}
	if isNilValue(bot.appData) {
		bot.logger.Warnln("app data is nil; skipping app-data logger writer")
		return bot
	}
	w := writer(bot.appData)
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
//	_ = bot.Close()
//
// A Bot is single-use. After RunWithContext returns, later calls return ErrBotAlreadyRun.
func (bot *Bot[T]) RunWithContext(ctx context.Context) error {
	if len(bot.prefixes) == 0 {
		return ErrNoPrefixes
	}

	if len(bot.plugins) == 0 {
		return ErrNoPlugins
	}
	if err := bot.beginRun(); err != nil {
		return err
	}
	defer bot.finishRun()

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
		retryDelay := time.Duration(0)
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
					retryDelay = nextPollRetryDelay(retryDelay)
					timer := time.NewTimer(retryDelay)
					select {
					case <-ctx.Done():
						if !timer.Stop() {
							<-timer.C
						}
						return
					case <-timer.C:
					}
					continue
				}
				retryDelay = 0

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
			bot.handle(ctx, u)
		})
	}
	pool.Stop() // Wait for all tasks to complete and stop the pool
	bot.runnerOnceWG.Wait()
	bot.runnerBgWG.Wait()
	return nil
}

// Run starts the bot using a background context.
//
// Equivalent to RunWithContext(context.Background()).
// Use this for simple bots where graceful shutdown is not required.
//
// For production use, prefer RunWithContext to handle SIGINT/SIGTERM gracefully.
func (bot *Bot[T]) Run() error {
	return bot.RunWithContext(context.Background())
}

func (bot *Bot[T]) beginRun() error {
	bot.runStateMu.Lock()
	defer bot.runStateMu.Unlock()
	if bot.running || bot.ran {
		return ErrBotAlreadyRun
	}
	bot.running = true
	bot.ran = true
	return nil
}

func (bot *Bot[T]) finishRun() {
	bot.runStateMu.Lock()
	bot.running = false
	bot.runStateMu.Unlock()
}

func nextPollRetryDelay(prev time.Duration) time.Duration {
	if prev <= 0 {
		return time.Second
	}
	next := prev * 2
	if next > 30*time.Second {
		return 30 * time.Second
	}
	return next
}

func isNilValue[T any](v T) bool {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return true
	}
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

func shouldWarnOnValueAppData[T any]() bool {
	t := reflect.TypeFor[T]()
	if t == reflect.TypeFor[NoData]() {
		return false
	}
	switch t.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		return false
	default:
		return true
	}
}

func clonePlugin[T AppData](p *Plugin[T]) Plugin[T] {
	cloned := Plugin[T]{
		name:        p.name,
		commands:    make(map[string]*Command[T], len(p.commands)),
		payloads:    make(map[string]*Command[T], len(p.payloads)),
		scenes:      make(map[string]*Scene[T], len(p.scenes)),
		middlewares: append(extypes.Slice[Middleware[T]](nil), p.middlewares...),
		skipAutoCmd: p.skipAutoCmd,
		logger:      p.logger,
		handlers:    make(map[tgapi.UpdateType]CommandExecutor[T]),
		onClose:     p.onClose,
	}

	for name, command := range p.commands {
		cloned.commands[name] = cloneCommand(command)
	}
	for name, command := range p.payloads {
		cloned.payloads[name] = cloneCommand(command)
	}
	for name, scene := range p.scenes {
		cloned.scenes[name] = cloneScene(scene)
	}
	maps.Copy(cloned.handlers, p.handlers)

	return cloned
}

func cloneCommand[T AppData](command *Command[T]) *Command[T] {
	if command == nil {
		return nil
	}

	cloned := *command
	cloned.args = append(extypes.Slice[CommandArg](nil), command.args...)
	cloned.middlewares = append(extypes.Slice[Middleware[T]](nil), command.middlewares...)
	return &cloned
}

func cloneScene[T AppData](scene *Scene[T]) *Scene[T] {
	if scene == nil {
		return nil
	}

	cloned := *scene
	cloned.steps = make(map[string]SceneHandler[T], len(scene.steps))
	cloned.commands = make(map[string]SceneHandler[T], len(scene.commands))

	for name, handler := range scene.steps {
		cloned.steps[name] = handler
	}
	for name, handler := range scene.commands {
		cloned.commands[name] = handler
	}

	return &cloned
}
