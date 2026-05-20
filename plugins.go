package laniakea

import (
	"errors"

	"git.scuroneko.dev/scuroneko/extypes"
	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/laniakea/utils"
	"git.scuroneko.dev/scuroneko/sneklog/v2"
)

// Plugin represents a collection of commands and payloads (e.g., callback handlers),
// with shared middleware and configuration.
//
// A Plugin is intended to be fully configured before it is passed to Bot.AddPlugins.
// After registration, treat the plugin as committed and do not mutate it further.
// Post-registration changes through the original *Plugin are not a supported API.
type Plugin[T AppData] struct {
	name        string                       // Name of the plugin (e.g., "admin", "user")
	commands    map[string]*Command[T]       // Registered commands (triggered by message)
	payloads    map[string]*Command[T]       // Registered payloads (triggered by callback data)
	scenes      map[string]*Scene[T]         // Optional scenes for multi-step interactions
	middlewares extypes.Slice[Middleware[T]] // Shared middlewares for all commands/payloads
	skipAutoCmd bool                         // If true, all commands in this plugin are excluded from auto-help
	logger      *sneklog.Logger
	loggerOwned bool // true when the logger was created by the bot during registration; only owned loggers are closed by Close

	messageFallback CommandExecutor[T]
	handlers        map[tgapi.UpdateType]CommandExecutor[T]

	onClose func() error
}

// NewPlugin creates a new Plugin with the given name.
func NewPlugin[T AppData](name string) *Plugin[T] {
	return &Plugin[T]{
		name:        name,
		commands:    make(map[string]*Command[T]),
		payloads:    make(map[string]*Command[T]),
		middlewares: make(extypes.Slice[Middleware[T]], 0),
		scenes:      make(map[string]*Scene[T]),
		skipAutoCmd: false,
		logger:      nil,
		handlers:    make(map[tgapi.UpdateType]CommandExecutor[T]),
	}
}

// AddCommand registers a command in the plugin.
func (p *Plugin[T]) AddCommand(command *Command[T]) *Plugin[T] {
	if command == nil {
		if p.logger != nil {
			p.logger.Warnln("trying to add nil command")
		}
		return p
	}
	if _, exists := p.commands[command.command]; exists && p.logger != nil {
		p.logger.Warnf("command '%s' already registered in plugin '%s'; overwriting", command.command, p.name)
	}
	p.commands[command.command] = command
	return p
}

// Command creates and immediately adds a new command to the plugin.
// Returns the created command for further configuration.
func (p *Plugin[T]) Command(command string, exec CommandExecutor[T], args ...CommandArg) *Command[T] {
	cmd := NewCommand(command, exec, args...)
	p.AddCommand(cmd)
	return cmd
}

// AddPayload registers a payload (e.g., callback query data) in the plugin.
// Payloads are triggered by inline button callback_data, not by message text.
func (p *Plugin[T]) AddPayload(command *Command[T]) *Plugin[T] {
	if command == nil {
		if p.logger != nil {
			p.logger.Warnln("trying to add nil command")
		}
		return p
	}
	if _, exists := p.payloads[command.command]; exists && p.logger != nil {
		p.logger.Warnf("payload '%s' is already registered in plugin '%s'; overwriting", command.command, p.name)
	}
	p.payloads[command.command] = command
	return p
}

// Payload creates and immediately adds a new payload command to the plugin.
// Returns the created payload command for further configuration.
func (p *Plugin[T]) Payload(command string, exec CommandExecutor[T], args ...CommandArg) *Command[T] {
	cmd := NewCommand(command, exec, args...)
	p.AddPayload(cmd)
	return cmd
}

// Scene creates, registers, and returns a new scene owned by the plugin.
func (p *Plugin[T]) Scene(name string) *Scene[T] {
	scene := NewScene[T](name)
	scene.setPluginName(p.name)
	p.AddScene(scene)
	return scene
}

// AddScene registers a multi-step scene in the plugin.
func (p *Plugin[T]) AddScene(scene *Scene[T]) *Plugin[T] {
	if scene == nil {
		return p
	}
	scene.pluginName = p.name
	if _, exists := p.scenes[scene.name]; exists && p.logger != nil {
		p.logger.Warnf("scene '%s'да already registered in plugin '%s'; overwriting", scene.name, p.name)
	}
	p.scenes[scene.name] = scene
	return p
}

// CommandGroup configures and registers a prefixed command group.
func (p *Plugin[T]) CommandGroup(prefix string, groupFunc func(group *CommandGroup[T])) *Plugin[T] {
	if groupFunc == nil {
		return p
	}
	group := NewCommandGroup[T](prefix)
	groupFunc(group)
	if len(group.commands) == 0 {
		return p
	}
	for _, cmd := range group.Build() {
		p.AddCommand(cmd)
	}
	return p
}

// AddCommandGroup registers every command built by group.
func (p *Plugin[T]) AddCommandGroup(group *CommandGroup[T]) *Plugin[T] {
	if group == nil {
		return p
	}
	if len(group.commands) == 0 {
		return p
	}
	for _, cmd := range group.Build() {
		p.AddCommand(cmd)
	}
	return p
}

// UsePolicy registers a Policy as plugin middleware for all plugin handlers.
func (p *Plugin[T]) UsePolicy(name string, policy Policy[T]) *Plugin[T] {
	mw := RequirePolicy(name, policy)
	return p.AddMiddleware(mw)
}

// AddUpdateHandler registers a handler for a non-command update type.
// Message, channel post, and callback query updates stay on the command/payload flow.
func (p *Plugin[T]) AddUpdateHandler(t tgapi.UpdateType, handler CommandExecutor[T]) *Plugin[T] {
	switch t {
	case tgapi.UpdateTypeMessage, tgapi.UpdateTypeChannelPost, tgapi.UpdateTypeCallbackQuery:
		if p.logger == nil {
			logger := utils.CreateLogger(p.name, utils.GetLoggerLevel(), utils.LogFormatText, nil)
			logger.Warnf("%s can't be registered through AddUpdateHandler. Use AddPayload/Payload or AddCommand/Command", t)
			_ = logger.Close()
			return p
		}
		p.logger.Warnf("%s can't be registered through AddUpdateHandler. Use AddPayload/Payload or AddCommand/Command", t)
		return p
	}
	p.handlers[t] = handler
	return p
}

// AddMiddleware adds a middleware to the plugin's global middleware chain.
// Middlewares are executed before any command or payload.
func (p *Plugin[T]) AddMiddleware(middleware Middleware[T]) *Plugin[T] {
	p.middlewares = p.middlewares.Push(middleware)
	return p
}

// SkipCommandAutoGen marks the entire plugin to be excluded from auto-generated help menus.
func (p *Plugin[T]) SkipCommandAutoGen() *Plugin[T] {
	p.skipAutoCmd = true
	return p
}

// SetLogger sets the logger used for this plugin's handlers.
//
// Call this before Bot.AddPlugins. If the plugin is already registered, changing
// the original *Plugin does not update the Bot's internal copy.
func (p *Plugin[T]) SetLogger(l *sneklog.Logger) *Plugin[T] {
	p.logger = l
	return p
}

// RemoveLogger clears the custom logger for this plugin.
//
// Call this before Bot.AddPlugins. If the plugin is already registered, changing
// the original *Plugin does not update the Bot's internal copy.
func (p *Plugin[T]) RemoveLogger() *Plugin[T] {
	p.logger = nil
	return p
}

// SetOnClose registers a callback invoked from Plugin.Close after the plugin
// logger is closed.
//
// Call this before Bot.AddPlugins. If the plugin is already registered, changing
// the original *Plugin does not update the Bot's internal copy.
func (p *Plugin[T]) SetOnClose(f func() error) *Plugin[T] {
	p.onClose = f
	return p
}

// SetMessageFallback registers a fallback handler for messages that do not
// match a command.
func (p *Plugin[T]) SetMessageFallback(handler CommandExecutor[T]) *Plugin[T] {
	p.messageFallback = handler
	return p
}

// Close releases plugin-owned resources such as its logger and optional
// OnClose callback.
//
// Only loggers created by the bot during registration are closed. A logger
// supplied via SetLogger remains the caller's responsibility — the framework
// never closes a logger it does not own.
func (p *Plugin[T]) Close() error {
	var e []error
	if p.logger != nil && p.loggerOwned {
		if err := p.logger.Close(); err != nil {
			e = append(e, err)
		}
	}
	if p.onClose != nil {
		if err := p.onClose(); err != nil {
			e = append(e, err)
		}
	}
	return errors.Join(e...)
}

func (p *Plugin[T]) executeCmd(cmd string, ctx *MessageContext, db T) error {
	command, exists := p.commands[cmd]
	if !exists {
		return AsInternalError(errCommandNotFound)
	}

	if err := command.validateArgs(ctx.Args); err != nil {
		return AsUserError(err)
	}

	// Run command-specific middlewares
	for _, m := range command.middlewares {
		if !m.Execute(ctx, db) {
			return AsInternalError(errors.New("middleware blocked call"))
		}
	}

	// Execute command
	return command.exec(ctx, db)
}

func (p *Plugin[T]) executePayload(payload string, ctx *MessageContext, db T) error {
	command, exists := p.payloads[payload]
	if !exists {
		return AsInternalError(errPayloadNotFound)
	}

	if err := command.validateArgs(ctx.Args); err != nil {
		return AsUserError(err)
	}

	// Run command-specific middlewares
	for _, m := range command.middlewares {
		if !m.Execute(ctx, db) {
			return AsInternalError(errors.New("middleware blocked call"))
		}
	}

	// Execute payload
	return command.exec(ctx, db)
}

func (p *Plugin[T]) executeMiddlewares(ctx *MessageContext, db T) bool {
	for _, m := range p.middlewares {
		if !m.Execute(ctx, db) {
			return false
		}
	}
	return true
}

// MiddlewareExecutor is the function type for middleware logic.
// Returns true to continue execution, false to block it.
// If async, return value is ignored.
type MiddlewareExecutor[T AppData] func(ctx *MessageContext, db T) bool

// Middleware represents a reusable execution interceptor.
// Can be synchronous (blocking) or asynchronous (non-blocking).
type Middleware[T AppData] struct {
	name     string                // Human-readable name for logging/debugging
	executor MiddlewareExecutor[T] // Function to execute
	order    int                   // Sort order for bot-level middleware ordering
	async    bool                  // If true, runs in goroutine and doesn't block
}

// NewMiddleware creates a new synchronous middleware.
func NewMiddleware[T AppData](name string, executor MiddlewareExecutor[T]) Middleware[T] {
	return Middleware[T]{name, executor, 0, false}
}

// SetOrder sets the execution order (currently ignored).
func (m Middleware[T]) SetOrder(order int) Middleware[T] {
	m.order = order
	return m
}

// SetAsync marks the middleware to run asynchronously.
// Execution continues regardless of its return value.
func (m Middleware[T]) SetAsync(async bool) Middleware[T] {
	m.async = async
	return m
}

// Execute runs the middleware.
// If async, runs in a goroutine and returns true immediately.
// Otherwise, returns the result of the executor.
//
// Async note: the goroutine receives a shallow copy of MessageContext, so
// scalar fields (FromID, ChatID, CallbackQueryID, ...) remain a stable
// snapshot. Pointer and slice fields (Msg, From, Chat, API, Logger, Args)
// continue to share storage with the synchronous flow. Async middleware
// must treat those fields as read-only — mutating them races the sync chain
// that mutates the same context concurrently.
func (m Middleware[T]) Execute(ctx *MessageContext, db T) bool {
	if m.async {
		ctxCopy := *ctx
		go func(ctx MessageContext) {
			m.executor(&ctx, db)
		}(ctxCopy)
		return true
	}
	return m.executor(ctx, db)
}
