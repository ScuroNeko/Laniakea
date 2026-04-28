package laniakea

import (
	"fmt"
	"sort"

	"git.scuroneko.dev/scuroneko/laniakea/utils"
)

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
			cloned.logger = utils.CreateLogger(cloned.name, level, bot.logFormat, bot.logFormatter)
		}
		bot.addTokenReplacer(cloned.logger)
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

// UsePolicy registers a Policy as a bot-level middleware.
func (bot *Bot[T]) UsePolicy(name string, policy Policy[T]) *Bot[T] {
	mw := RequirePolicy(name, policy)
	return bot.AddMiddleware(mw)
}

// AddRunner registers a background runner to execute concurrently with the bot.
//
// Runners are goroutines that run independently of update processing.
// Common use cases:
//   - Periodic cleanup (e.g., expiring drafts, clearing temp files)
//   - Metrics collection or health checks
//   - Scheduled tasks (e.g., daily announcements)
//
// Runners start from the bot runtime entry points, immediately after
// RunWithContext or RunWebHookWithContext begins.
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
//	bot.AddAppDataLoggerWriter(func(data *MyAppData) sneklog.LoggerWriter {
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
	bot.logger.AddWriters(w)
	if bot.requestLogger != nil {
		bot.requestLogger.AddWriters(w)
	}
	for _, l := range bot.managedExtraLoggers() {
		l.AddWriters(w)
	}
	for _, p := range bot.plugins {
		if p.logger != nil {
			p.logger.AddWriters(w)
		}
	}
	bot.addTokenReplacer(bot.logger, bot.requestLogger)
	bot.addTokenReplacer(bot.managedExtraLoggers()...)
	return bot
}
