package laniakea

import (
	"os"
	"strconv"
	"strings"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

// BotOpts holds configuration options for initializing a Bot.
//
// Values are loaded from environment variables via LoadOptsFromEnv().
// Use NewOpts() to create a zero-value struct and set fields manually.
type BotOpts struct {
	// Token is the Telegram bot token (required).
	Token string

	// UpdateTypes is a list of update types to listen for.
	// Example: "["message", "edited_message", "callback_query"]"
	// Defaults to empty (Telegram will return all types).
	UpdateTypes []tgapi.UpdateType

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

	// MaxWorkers is the maximum number of concurrency running update handlers.
	MaxWorkers int
}

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

	stringUpdateTypes := strings.Split(os.Getenv("UPDATE_TYPES"), ";")
	updateTypes := make([]tgapi.UpdateType, len(stringUpdateTypes))
	for i, updateType := range stringUpdateTypes {
		updateTypes[i] = tgapi.UpdateType(updateType)
	}

	return &BotOpts{
		Token:       os.Getenv("TG_TOKEN"),
		UpdateTypes: updateTypes,

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

// SetToken sets the Telegram bot token (required).
func (opts *BotOpts) SetToken(token string) *BotOpts {
	opts.Token = token
	return opts
}

// SetUpdateTypes sets the list of update types to listen for.
// If empty (default), Telegram will return all update types.
// Example: opts.SetUpdateTypes("message", "callback_query")
func (opts *BotOpts) SetUpdateTypes(types ...tgapi.UpdateType) *BotOpts {
	opts.UpdateTypes = types
	return opts
}

// SetDebug enables or disables debug-level logging.
// Default is false.
func (opts *BotOpts) SetDebug(debug bool) *BotOpts {
	opts.Debug = debug
	return opts
}

// SetErrorTemplate sets the format string for error messages sent to users.
// Use "%s" to insert the actual error. Example: "❌ Error: %s"
// If not set, defaults to "%s".
func (opts *BotOpts) SetErrorTemplate(tpl string) *BotOpts {
	opts.ErrorTemplate = tpl
	return opts
}

// SetPrefixes sets the command prefixes (e.g., "/", "!").
// If not set via environment, defaults to ["/"].
func (opts *BotOpts) SetPrefixes(prefixes ...string) *BotOpts {
	opts.Prefixes = prefixes
	return opts
}

// SetLoggerBasePath sets the directory where log files are written.
// Defaults to "./".
func (opts *BotOpts) SetLoggerBasePath(path string) *BotOpts {
	opts.LoggerBasePath = path
	return opts
}

// SetUseRequestLogger enables detailed logging of all Telegram API requests.
// Default is false.
func (opts *BotOpts) SetUseRequestLogger(use bool) *BotOpts {
	opts.UseRequestLogger = use
	return opts
}

// SetWriteToFile enables writing logs to files (main.log and requests.log).
// Default is false.
func (opts *BotOpts) SetWriteToFile(write bool) *BotOpts {
	opts.WriteToFile = write
	return opts
}

// SetUseTestServer enables using Telegram's test server (https://api.telegram.org/bot<token>/test).
// Default is false.
func (opts *BotOpts) SetUseTestServer(use bool) *BotOpts {
	opts.UseTestServer = use
	return opts
}

// SetAPIUrl overrides the default Telegram API endpoint (useful for proxies or self-hosted).
// If not set, defaults to "https://api.telegram.org".
func (opts *BotOpts) SetAPIUrl(url string) *BotOpts {
	opts.APIUrl = url
	return opts
}

// SetRateLimit sets the maximum number of API requests per second.
// Telegram allows up to 30 req/s for most bots. Defaults to 30.
func (opts *BotOpts) SetRateLimit(limit int) *BotOpts {
	opts.RateLimit = limit
	return opts
}

// SetDropRLOverflow drops incoming updates when rate limit is exceeded instead of queuing.
// Use this to prioritize responsiveness over reliability. Default is false.
func (opts *BotOpts) SetDropRLOverflow(drop bool) *BotOpts {
	opts.DropRLOverflow = drop
	return opts
}

// SetMaxWorkers sets the maximum number of concurrent update handlers.
// Must be called before NewBot, as the value is captured during bot creation.
//
// The optimal value depends on your bot's workload:
//   - For I/O-bound handlers (e.g., database queries, external API calls), you may
//     need more workers, but be mindful of downstream service limits.
//   - For CPU-bound handlers, keep workers close to the number of CPU cores.
//
// Recommended starting points (adjust based on profiling and monitoring):
//   - Small to medium bots with fast handlers: 16–32
//   - Medium to large bots with fast handlers: 32–64
//   - Large bots with heavy I/O: 64–128 (ensure your infrastructure can handle it)
//
// The default is 32. Monitor queue length and processing latency to fine-tune.
func (opts *BotOpts) SetMaxWorkers(workers int) *BotOpts {
	opts.MaxWorkers = workers
	return opts
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
