package laniakea

import (
	"context"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
	"time"

	"git.scuroneko.dev/scuroneko/extypes"
	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/laniakea/utils"

	"git.scuroneko.dev/scuroneko/sneklog/v2"
	"github.com/alitto/pond/v2"
)

func (bot *Bot[T]) addTokenReplacer(loggers ...*sneklog.Logger) {
	if bot.token == "" {
		return
	}
	for _, logger := range loggers {
		if logger == nil {
			continue
		}
		logger.AddReplacer(bot.token, "<TOKEN>")
	}
}

func appendUniqueLogger(loggers []*sneklog.Logger, logger *sneklog.Logger) []*sneklog.Logger {
	if logger == nil {
		return loggers
	}
	if slices.Contains(loggers, logger) {
		return loggers
	}
	return append(loggers, logger)
}

func (bot *Bot[T]) managedExtraLoggers() []*sneklog.Logger {
	loggers := append([]*sneklog.Logger(nil), bot.extraLoggers...)
	if bot.api != nil {
		loggers = appendUniqueLogger(loggers, bot.api.GetLogger())
	}
	if bot.uploader != nil {
		loggers = appendUniqueLogger(loggers, bot.uploader.GetLogger())
	}
	return loggers
}

func (bot *Bot[T]) enqueueUpdate(ctx context.Context, update tgapi.Update) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case bot.updateQueue <- new(update):
		return nil
	}
}

func (bot *Bot[T]) startUpdateWorkers(ctx context.Context) {
	pool := pond.NewPool(bot.maxWorkers)
	for update := range bot.updateQueue {
		u := update // capture loop variable
		pool.Submit(func() {
			bot.handle(ctx, u)
		})
	}
	pool.Stop() // Wait for all tasks to complete and stop the pool
}

func (bot *Bot[T]) initLoggers(opts *BotOpts) {
	level := sneklog.FATAL
	if opts.Debug {
		level = sneklog.DEBUG
	}

	format, formatter := opts.LogFormat, opts.LogFormatter
	if bot.logger == nil {
		bot.logger = utils.CreateLogger("BOT", level, format, formatter)
		if opts.WriteToFile {
			path := fmt.Sprintf("%s/main.log", strings.TrimRight(opts.LoggerBasePath, "/"))
			logger, err := utils.CreateFileLogger("BOT", level, path, format, formatter)
			if err != nil {
				bot.logger.Errorln(err)
			} else {
				bot.logger = logger
			}
		}
	}

	if opts.UseRequestLogger && bot.requestLogger == nil {
		bot.requestLogger = utils.CreateLogger("REQUESTS", level, format, formatter)
		if opts.WriteToFile {
			path := fmt.Sprintf("%s/requests.log", strings.TrimRight(opts.LoggerBasePath, "/"))
			logger, err := utils.CreateFileLogger("REQUESTS", level, path, format, formatter)
			if err != nil {
				bot.logger.Errorln(err)
			} else {
				bot.requestLogger = logger
			}
		}
	}

	bot.addTokenReplacer(bot.logger, bot.requestLogger)
	bot.addTokenReplacer(bot.managedExtraLoggers()...)
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
		name:            p.name,
		commands:        make(map[string]*Command[T], len(p.commands)),
		payloads:        make(map[string]*Command[T], len(p.payloads)),
		scenes:          make(map[string]*Scene[T], len(p.scenes)),
		middlewares:     append(extypes.Slice[Middleware[T]](nil), p.middlewares...),
		skipAutoCmd:     p.skipAutoCmd,
		logger:          p.logger,
		messageFallback: p.messageFallback,
		handlers:        make(map[tgapi.UpdateType]CommandExecutor[T]),
		onClose:         p.onClose,
	}

	for name, command := range p.commands {
		cloned.commands[name] = command.clone()
	}
	for name, command := range p.payloads {
		cloned.payloads[name] = command.clone()
	}
	for name, scene := range p.scenes {
		cloned.scenes[name] = scene.clone()
	}
	maps.Copy(cloned.handlers, p.handlers)

	return cloned
}
