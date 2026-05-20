package laniakea

import (
	"errors"
	"regexp"

	"git.scuroneko.dev/scuroneko/extypes"
)

// CommandValueType defines the expected type of command argument.
type CommandValueType string

const (
	// CommandValueString expects any non-empty string.
	CommandValueString CommandValueType = "string"
	// CommandValueInt expects a decimal integer (digits only).
	CommandValueInt CommandValueType = "int"
	// CommandValueBool expects an exact "true" or "false".
	CommandValueBool CommandValueType = "bool"
	// CommandValueAny accepts any input without validation.
	CommandValueAny CommandValueType = "any"
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

var (
	errCommandNotFound = errors.New("command not found")
	errPayloadNotFound = errors.New("payload not found")
)

// CommandArg defines a single argument for a command, including type, regex,
// and whether it is required.
type CommandArg struct {
	valueType CommandValueType // Type of expected value
	text      string           // Human-readable description (not used in validation)
	regex     *regexp.Regexp   // Regex used to validate input
	required  bool             // Whether this argument must be provided
}

// NewCommandArg creates an optional argument without value validation.
func NewCommandArg(text string) CommandArg {
	return CommandArg{CommandValueAny, text, nil, false}
}

// SetValueType sets expected value type and switches built-in validation regexp.
func (c CommandArg) SetValueType(t CommandValueType) CommandArg {
	var regex *regexp.Regexp
	switch t {
	case CommandValueInt:
		regex = CommandRegexInt
	case CommandValueBool:
		regex = CommandRegexBool
	case CommandValueString:
		regex = CommandRegexString
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
// It receives the message context and injected application data.
// Returning a non-nil error routes it through the bot's error handler.
type CommandExecutor[T AppData] func(ctx *MessageContext, dbContext T) error

// Command represents a bot command with arguments, description, and executor.
// Can be registered in a Plugin and optionally skipped from auto-generation.
type Command[T AppData] struct {
	command     string                       // The command trigger (e.g., "/start")
	description string                       // Human-readable description for help
	exec        CommandExecutor[T]           // Function to execute when command is triggered
	args        extypes.Slice[CommandArg]    // List of expected arguments
	middlewares extypes.Slice[Middleware[T]] // Optional middleware chain
	skipAutoCmd bool                         // If true, this command won't be auto-added to help menus
}

// NewCommand creates a new Command with the given identifier, executor, and arguments.
//
// The identifier is used as the routing key for both /-prefixed commands and
// callback payloads — the difference is registration: pass the result to
// Plugin.AddCommand/Plugin.Command for message routing, or to
// Plugin.AddPayload/Plugin.Payload for callback_data routing.
//
// For /-commands the identifier must not include the leading slash
// (e.g. "start", not "/start") and should match [_a-z0-9]{1,32} to satisfy
// Telegram's BotCommand validation. Payload identifiers may use any bytes
// that fit Telegram's callback_data limit, though the configured payload
// encoding may impose its own restrictions.
func NewCommand[T any](command string, exec CommandExecutor[T], args ...CommandArg) *Command[T] {
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
			continue // Skip validation for CommandValueAny.
		}
		if !cmdArg.regex.MatchString(arg) {
			return ErrCmdArgRegexpMismatch
		}
	}
	return nil
}

func (c *Command[T]) clone() *Command[T] {
	if c == nil {
		return nil
	}

	cloned := *c
	cloned.args = append(extypes.Slice[CommandArg](nil), c.args...)
	cloned.middlewares = append(extypes.Slice[Middleware[T]](nil), c.middlewares...)
	return &cloned
}

// CommandGroup builds a set of commands with a shared name prefix and middleware.
type CommandGroup[T any] struct {
	prefix      string
	middlewares extypes.Slice[Middleware[T]]
	commands    extypes.Slice[*Command[T]]
}

// NewCommandGroup creates a command group that prefixes every added command.
func NewCommandGroup[T any](prefix string) *CommandGroup[T] {
	return &CommandGroup[T]{
		prefix: prefix,

		middlewares: make([]Middleware[T], 0),
		commands:    make([]*Command[T], 0),
	}
}

// Use adds middleware that runs before each command's own middleware.
func (g *CommandGroup[T]) Use(m Middleware[T]) *CommandGroup[T] {
	g.middlewares = append(g.middlewares, m)
	return g
}

// AddCommand adds a prefixed copy of cmd to the group.
func (g *CommandGroup[T]) AddCommand(cmd *Command[T]) *CommandGroup[T] {
	if cmd == nil {
		return g
	}
	newCmd := cmd.clone()
	newCmd.command = g.prefix + cmd.command
	g.commands = g.commands.Push(newCmd)
	return g
}

// Build returns command copies with group middleware prepended.
func (g *CommandGroup[T]) Build() []*Command[T] {
	commands := make([]*Command[T], 0)
	for _, cmd := range g.commands {
		cloned := cmd.clone()
		cloned.middlewares = append(
			append(extypes.Slice[Middleware[T]]{}, g.middlewares...),
			cloned.middlewares...,
		)
		commands = append(commands, cloned)
	}
	return commands
}
