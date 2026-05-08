package laniakea

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"git.scuroneko.dev/scuroneko/extypes"
	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/laniakea/utils"
	"git.scuroneko.dev/scuroneko/sneklog/v2"
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

// AppDataLogger builds a sneklog.LoggerWriter from injected application data.
//
// Use it when shared application data exposes a log sink or adapter that should
// receive framework logs.
type AppDataLogger[T AppData] func(data T) sneklog.LoggerWriter

// BotPayloadType defines the serialization format for callback data payloads.
type BotPayloadType string

var (
	// BotPayloadBase64 encodes callback data as a Base64 string.
	BotPayloadBase64 BotPayloadType = "base64"
	// BotPayloadJSON encodes callback data as a JSON string.
	BotPayloadJSON BotPayloadType = "json"
)

var (
	// ErrNoPrefixes reports that the bot was started without any command prefixes.
	ErrNoPrefixes = errors.New("no prefixes defined")
	// ErrNoPlugins reports that the bot was started without any registered plugins.
	ErrNoPlugins = errors.New("no plugins defined")
	// ErrBotAlreadyRun reports that Run, RunWithContext, or RunWebhookWithContext was called more than once.
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
// Runtime accessors are safe for concurrent use. Configure the bot before Run,
// RunWithContext, or RunWebhookWithContext.
// A Bot is single-use: after Run, RunWithContext, or RunWebhookWithContext returns,
// create a new Bot for the next session.
type Bot[T AppData] struct {
	token             string
	debug             bool
	errorTemplate     string
	username          string
	payloadType       BotPayloadType
	strictPayloadType bool
	maxWorkers        int

	logFormat     utils.LogFormat
	logFormatter  *sneklog.Formatter
	logger        *sneklog.Logger // Main bot logger (JSON stdout + optional file)
	requestLogger *sneklog.Logger // Optional request-level API logging
	useReqLogger  bool
	webhookLogger *sneklog.Logger                // Webhook logger. Available only after Bot.RunWebhookWithContext.
	extraLoggers  extypes.Slice[*sneklog.Logger] // API, Uploader, and custom loggers

	plugins     []Plugin[T]     // Command/event handlers
	middlewares []Middleware[T] // Pre-processing filters (sorted by order)
	prefixes    []string        // Command prefixes (e.g., "/", "!")
	runners     []Runner[T]     // Background tasks (e.g., cleanup, cron)

	api           *tgapi.API      // Telegram API client
	uploader      *tgapi.Uploader // File uploader
	l10n          *L10n           // Localization manager
	draftProvider *DraftProvider  // Draft message builder
	observer      Observer        // Optional event observer for instrumentation

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

	limiter := utils.NewRateLimiter()
	limiter.SetGlobalRate(opts.RateLimit)

	apiOpts := tgapi.NewAPIOpts(opts.Token).
		SetAPIURL(opts.APIURL).
		UseTestServer(opts.UseTestServer).
		SetLimiter(limiter).
		SetDropRateLimitOverflow(opts.DropRateLimitOverflow).
		SetLogFormat(opts.LogFormat).
		SetLogFormatter(opts.LogFormatter)
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
		logFormat:         opts.LogFormat,
		logFormatter:      opts.LogFormatter,
		useReqLogger:      opts.UseRequestLogger,

		plugins:       make([]Plugin[T], 0),
		updateTypes:   append([]tgapi.UpdateType{}, opts.UpdateTypes...),
		runners:       make([]Runner[T], 0),
		extraLoggers:  make([]*sneklog.Logger, 0),
		l10n:          &L10n{},
		draftProvider: NewRandomDraftProvider(api),

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

	if opts.FileConfigVersion > 0 && opts.FileConfigVersion < ConfigVersion {
		bot.logger.Warnln(
			fmt.Sprintf(
				"Config file version %d is older than library version %d; please update your config file to access new features and avoid compatibility issues",
				opts.FileConfigVersion,
				ConfigVersion,
			),
		)
	}

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
	bot.logger.Debugln("Bot initialized with configuration:", fmt.Sprintf("%+v", opts))

	return bot, nil
}

// SetLogger replaces the main bot logger.
func (bot *Bot[T]) SetLogger(l *sneklog.Logger) *Bot[T] {
	bot.logger = l
	return bot
}

// SetRequestLogger replaces the request-level logger.
func (bot *Bot[T]) SetRequestLogger(l *sneklog.Logger) *Bot[T] {
	bot.requestLogger = l
	return bot
}

// SetWebhookLogger replaces the webhook logger.
func (bot *Bot[T]) SetWebhookLogger(l *sneklog.Logger) *Bot[T] {
	bot.webhookLogger = l
	return bot
}

func (bot *Bot[T]) GetAPI() *tgapi.API { return bot.api }

func (bot *Bot[T]) GetUploader() *tgapi.Uploader { return bot.uploader }

// Close gracefully shuts down bot-owned resources.
//
// Close shuts down, in order:
//   - Registered plugins via Plugin.Close
//   - Webhook logger (if initialized)
//   - Uploader (waits for pending uploads)
//   - API client internals
//   - RequestLogger (if enabled)
//   - Main logger
//
// RunWithContext and RunWebhookWithContext do not call Close automatically.
// The caller is responsible for invoking Close after runtime returns to release
// these resources.
//
// Close returns a joined error containing all shutdown failures, if any.
func (bot *Bot[T]) Close() error {
	var e []error
	logCloseErr := func(err error) {
		if err == nil {
			return
		}
		if bot.logger != nil {
			bot.logger.Errorln(err)
		}
		e = append(e, err)
	}

	for _, p := range bot.plugins {
		if err := p.Close(); err != nil {
			e = append(e, err)
		}
	}
	if bot.webhookLogger != nil {
		if err := bot.webhookLogger.Close(); err != nil {
			logCloseErr(err)
		}
		bot.webhookLogger = nil
	}
	if bot.uploader != nil {
		if err := bot.uploader.Close(); err != nil {
			logCloseErr(err)
		}
	}
	if bot.api != nil {
		if err := bot.api.Close(); err != nil {
			logCloseErr(err)
		}
	}
	if bot.requestLogger != nil {
		if err := bot.requestLogger.Close(); err != nil {
			logCloseErr(err)
		}
	}
	if bot.logger != nil {
		if err := bot.logger.Close(); err != nil {
			e = append(e, err)
		}
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

// GetLogger returns the main bot logger.
func (bot *Bot[T]) GetLogger() *sneklog.Logger { return bot.logger }

// GetRequestLogger returns the request-level logger, if configured.
func (bot *Bot[T]) GetRequestLogger() *sneklog.Logger { return bot.requestLogger }

// GetWebhookLogger returns the webhook logger, if configured.
func (bot *Bot[T]) GetWebhookLogger() *sneklog.Logger { return bot.webhookLogger }

// GetLoggerLevel returns the effective log level derived from the bot's debug
// flag.
func (bot *Bot[T]) GetLoggerLevel() sneklog.LogLevel {
	level := sneklog.FATAL
	if bot.debug {
		level = sneklog.DEBUG
	}
	return level
}

// L10n translates a key in the given language.
// Returns empty string if translation not found.
func (bot *Bot[T]) L10n(lang, key string) string {
	return bot.l10n.Translate(lang, key)
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
// If you are switching an existing deployment from webhook delivery to polling,
// delete the current webhook first with CloseWebhook or tgapi.DeleteWebhook.
// Telegram keeps webhook delivery active until the webhook is removed.
//
// RunWithContext does not close API, uploader, or logger resources on return.
// The caller must invoke Close after RunWithContext finishes.
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
	if !bot.useReqLogger && bot.requestLogger != nil {
		bot.logger.Warnln("Opts#UseRequestLogger is false, but Bot#requestLogger present. Remove Bot#SetRequestLogger or set Opts#UseRequestLogger to true!")
		err := bot.requestLogger.Close()
		if err != nil {
			bot.logger.Errorln(err)
		}
		bot.requestLogger = nil
	}
	if bot.webhookLogger != nil {
		bot.logger.Warnln("Bot#webhookLogger present. You shouldn't set this, if ran in Long Polling mode!")
		err := bot.webhookLogger.Close()
		if err != nil {
			bot.logger.Errorln(err)
		}
		bot.webhookLogger = nil
	}

	bot.ExecRunners(ctx)

	// Start update polling in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				bot.logger.Errorln(fmt.Sprintf("panic in update polling: %v", r))
			}
			close(bot.updateQueue)
		}()
		backoffDelay := time.Duration(0)
		retryCount := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				updates, err := bot.Updates(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
						return
					}
					retryDelay, ok := pollRetryAfterDelay(err)
					if ok {
						bot.logger.Warnln("getUpdates rate limited; retrying after", retryDelay)
						backoffDelay = 0
					} else {
						bot.logger.Errorln("failed to fetch updates:", err)
						backoffDelay = nextPollRetryDelay(backoffDelay)
						retryDelay = backoffDelay
					}
					retryCount++
					bot.safeEmitEvent(ctx, PollingRetryEvent{
						Attempt: retryCount,
						Delay:   retryDelay,
						Err:     err,
					})
					bot.safeEmitEvent(ctx, ErrorEvent{
						Plugin:      "bot",
						HandlerKind: HandlerPollingKind,
						HandlerName: "getUpdates",
						Err:         err,
						UserFacing:  false,
					})
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
				backoffDelay = 0
				retryCount = 0

				for _, update := range updates {
					if err := bot.enqueueUpdate(ctx, update); err != nil {
						return
					}
				}
			}
		}
	}()

	bot.logger.Infoln("Bot running. Press CTRL+C to exit.")
	// Start worker pool for concurrent update handling
	bot.startUpdateWorkers(ctx)

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
