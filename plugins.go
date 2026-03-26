package laniakea

import (
	"errors"
	"regexp"

	"git.nix13.pw/scuroneko/extypes"
	"git.nix13.pw/scuroneko/laniakea/tgapi"
	"git.nix13.pw/scuroneko/laniakea/utils"
	"git.nix13.pw/scuroneko/slog"
)

// CommandValueType defines the expected type of command argument.
type CommandValueType string

const (
	// CommandValueStringType expects any non-empty string.
	CommandValueStringType CommandValueType = "string"
	// CommandValueIntType expects a decimal integer (digits only).
	CommandValueIntType CommandValueType = "int"
	// CommandValueBoolType is reserved for future use (not implemented).
	CommandValueBoolType CommandValueType = "bool"
	// CommandValueAnyType accepts any input without validation.
	CommandValueAnyType CommandValueType = "any"
)

var (
	// CommandRegexInt matches one or more digits.
	CommandRegexInt = regexp.MustCompile(`^\d+$`)
	// CommandRegexString matches any non-empty string.
	CommandRegexString = regexp.MustCompile(`^.+$`)
	// CommandRegexBool matches true or false.
	CommandRegexBool = regexp.MustCompile(`^(true|false)$`)
)

// ErrCmdArgCountMismatch is returned when the number of provided arguments
// is less than the number of required arguments.
var ErrCmdArgCountMismatch = errors.New("command arg count mismatch")

// ErrCmdArgRegexpMismatch is returned when an argument fails regex validation.
var ErrCmdArgRegexpMismatch = errors.New("command arg regexp mismatch")

// CommandArg defines a single argument for a command, including type, regex,
// and whether it is required.
type CommandArg struct {
	valueType CommandValueType // Type of expected value
	text      string           // Human-readable description (not used in validation)
	regex     *regexp.Regexp   // Regex used to validate input
	required  bool             // Whether this argument must be provided
}

// NewCommandArg creates a new CommandArg with the given text and type.
// Uses a default regex based on the type (string or int).
// For CommandValueAnyType, no validation is performed.
func NewCommandArg(text string) CommandArg {
	return CommandArg{CommandValueAnyType, text, CommandRegexString, false}
}

// SetValueType sets expected value type and switches built-in validation regexp.
func (c CommandArg) SetValueType(t CommandValueType) CommandArg {
	regex := CommandRegexString
	switch t {
	case CommandValueIntType:
		regex = CommandRegexInt
	case CommandValueBoolType:
		regex = CommandRegexBool
	case CommandValueAnyType:
		regex = nil // Skip validation
	}
	c.valueType = t
	c.regex = regex
	return c
}

// SetRequired marks this argument as required.
// Returns the receiver for method chaining.
func (c CommandArg) SetRequired() CommandArg {
	c.required = true
	return c
}

// CommandExecutor is the function type that executes a command.
// It receives the message context and a database context (generic).
// Returning a non-nil error routes it through the bot's error handler.
type CommandExecutor[T DbContext] func(ctx *MsgContext, dbContext T) error

// Command represents a bot command with arguments, description, and executor.
// Can be registered in a Plugin and optionally skipped from auto-generation.
type Command[T DbContext] struct {
	command     string                       // The command trigger (e.g., "/start")
	description string                       // Human-readable description for help
	exec        CommandExecutor[T]           // Function to execute when command is triggered
	args        extypes.Slice[CommandArg]    // List of expected arguments
	middlewares extypes.Slice[Middleware[T]] // Optional middleware chain
	skipAutoCmd bool                         // If true, this command won't be auto-added to help menus
}

// NewCommand creates a new Command with the given executor, command string, and arguments.
// The command string should not include the leading slash (e.g., "start", not "/start").
func NewCommand[T any](exec CommandExecutor[T], command string, args ...CommandArg) *Command[T] {
	return &Command[T]{command, "", exec, args, make(extypes.Slice[Middleware[T]], 0), false}
}

// NewPayload creates a new Command with the given executor, command payload string, and arguments.
// The command string can contain any symbols, but it is recommended to use only "_", "-", ".", a-z, A-Z, and 0-9.
func NewPayload[T any](exec CommandExecutor[T], command string, args ...CommandArg) *Command[T] {
	return &Command[T]{command, "", exec, args, make(extypes.Slice[Middleware[T]], 0), false}
}

// Use adds a middleware to the command's execution chain.
// Middlewares are executed in the order they are added.
func (c *Command[T]) Use(m Middleware[T]) *Command[T] {
	c.middlewares = c.middlewares.Push(m)
	return c
}

// SetDescription sets the human-readable description of the command.
func (c *Command[T]) SetDescription(desc string) *Command[T] {
	c.description = desc
	return c
}

// SkipCommandAutoGen marks this command to be excluded from auto-generated help menus.
func (c *Command[T]) SkipCommandAutoGen() *Command[T] {
	c.skipAutoCmd = true
	return c
}

// Internal helper that validates provided command arguments.
func (c *Command[T]) validateArgs(args []string) error {
	for i := range c.args.Len() {
		if i >= len(args) && c.args.Get(i).required {
			return ErrCmdArgCountMismatch
		}
	}

	// Validate each argument against its regex
	for i, arg := range args {
		if i >= c.args.Len() {
			// Extra arguments beyond defined args are ignored
			break
		}
		cmdArg := c.args.Get(i)
		if cmdArg.regex == nil {
			continue // Skip validation for CommandValueAnyType
		}
		if !cmdArg.regex.MatchString(arg) {
			return ErrCmdArgRegexpMismatch
		}
	}
	return nil
}

// Plugin represents a collection of commands and payloads (e.g., callback handlers),
// with shared middleware and configuration.
//
// A Plugin is intended to be fully configured before it is passed to Bot.AddPlugins.
// After registration, treat the plugin as committed and do not mutate it further.
// Post-registration changes through the original *Plugin are not a supported API.
type Plugin[T DbContext] struct {
	name        string                       // Name of the plugin (e.g., "admin", "user")
	commands    map[string]*Command[T]       // Registered commands (triggered by message)
	payloads    map[string]*Command[T]       // Registered payloads (triggered by callback data)
	middlewares extypes.Slice[Middleware[T]] // Shared middlewares for all commands/payloads
	skipAutoCmd bool                         // If true, all commands in this plugin are excluded from auto-help
	logger      *slog.Logger

	handlers map[tgapi.UpdateType]CommandExecutor[T]

	onClose func() error
}

// NewPlugin creates a new Plugin with the given name.
func NewPlugin[T DbContext](name string) *Plugin[T] {
	return &Plugin[T]{
		name:        name,
		commands:    make(map[string]*Command[T]),
		payloads:    make(map[string]*Command[T]),
		middlewares: make(extypes.Slice[Middleware[T]], 0),
		skipAutoCmd: false,
		logger:      nil,
		handlers:    make(map[tgapi.UpdateType]CommandExecutor[T]),
	}
}

// AddCommand registers a command in the plugin.
// The command's .command field is used as the key.
func (p *Plugin[T]) AddCommand(command *Command[T]) *Plugin[T] {
	p.commands[command.command] = command
	return p
}

// NewCommand creates and immediately adds a new command to the plugin.
// Returns the created command for further configuration.
func (p *Plugin[T]) NewCommand(exec CommandExecutor[T], command string, args ...CommandArg) *Command[T] {
	cmd := NewCommand(exec, command, args...)
	p.AddCommand(cmd)
	return cmd
}

// AddPayload registers a payload (e.g., callback query data) in the plugin.
// Payloads are triggered by inline button callback_data, not by message text.
func (p *Plugin[T]) AddPayload(command *Command[T]) *Plugin[T] {
	p.payloads[command.command] = command
	return p
}

// NewPayload creates and immediately adds a new payload command to the plugin.
// Returns the created payload command for further configuration.
func (p *Plugin[T]) NewPayload(exec CommandExecutor[T], command string, args ...CommandArg) *Command[T] {
	cmd := NewPayload(exec, command, args...)
	p.AddPayload(cmd)
	return cmd
}

// AddUpdateHandler registers a handler for a non-command update type.
// Message, channel post, and callback query updates stay on the command/payload flow.
func (p *Plugin[T]) AddUpdateHandler(t tgapi.UpdateType, handler CommandExecutor[T]) *Plugin[T] {
	switch t {
	case tgapi.UpdateTypeMessage, tgapi.UpdateTypeChannelPost, tgapi.UpdateTypeCallbackQuery:
		if p.logger == nil {
			logger := utils.CreateLogger(p.name, utils.GetLoggerLevel())
			logger.Warnf("%s can't be registred through AddUpdateHandler. Use AddPayload/NewPayload or AddCommand/NewCommand", t)
			_ = logger.Close()
			return p
		}
		p.logger.Warnf("%s can't be registred through AddUpdateHandler. Use AddPayload/NewPayload or AddCommand/NewCommand", t)
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
func (p *Plugin[T]) SetLogger(l *slog.Logger) *Plugin[T] {
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

// Close releases plugin-owned resources such as its logger and optional
// OnClose callback.
func (p *Plugin[T]) Close() error {
	var e []error
	if p.logger != nil {
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

// Internal helper that validates and executes a command handler.
func (p *Plugin[T]) executeCmd(cmd string, ctx *MsgContext, db T) {
	command, exists := p.commands[cmd]
	if !exists {
		ctx.error(errors.New("command not found"))
		return
	}

	if err := command.validateArgs(ctx.Args); err != nil {
		ctx.error(err)
		return
	}

	// Run command-specific middlewares
	for _, m := range command.middlewares {
		if !m.Execute(ctx, db) {
			return
		}
	}

	// Execute command
	if err := command.exec(ctx, db); err != nil {
		ctx.error(err)
	}
}

// Internal helper that validates and executes a payload handler.
func (p *Plugin[T]) executePayload(payload string, ctx *MsgContext, db T) {
	command, exists := p.payloads[payload]
	if !exists {
		ctx.error(errors.New("payload not found"))
		return
	}

	if err := command.validateArgs(ctx.Args); err != nil {
		ctx.error(err)
		return
	}

	// Run command-specific middlewares
	for _, m := range command.middlewares {
		if !m.Execute(ctx, db) {
			return
		}
	}

	// Execute payload
	if err := command.exec(ctx, db); err != nil {
		ctx.error(err)
	}
}

// Internal helper that runs plugin middlewares in order.
func (p *Plugin[T]) executeMiddlewares(ctx *MsgContext, db T) bool {
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
type MiddlewareExecutor[T DbContext] func(ctx *MsgContext, db T) bool

// Middleware represents a reusable execution interceptor.
// Can be synchronous (blocking) or asynchronous (non-blocking).
type Middleware[T DbContext] struct {
	name     string                // Human-readable name for logging/debugging
	executor MiddlewareExecutor[T] // Function to execute
	order    int                   // Optional sort order (not used yet)
	async    bool                  // If true, runs in goroutine and doesn't block
}

// NewMiddleware creates a new synchronous middleware.
func NewMiddleware[T DbContext](name string, executor MiddlewareExecutor[T]) Middleware[T] {
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
func (m Middleware[T]) Execute(ctx *MsgContext, db T) bool {
	if m.async {
		ctx := *ctx // copy context to avoid race condition
		go func(ctx MsgContext) {
			m.executor(&ctx, db)
		}(ctx)
		return true
	}
	return m.executor(ctx, db)
}
