package laniakea

import (
	"errors"
	"fmt"
	"regexp"

	"git.scuroneko.dev/scuroneko/extypes"
)

// CommandValueType defines the expected type of command argument.
type CommandValueType string

const (
	// CommandValueStringType expects any non-empty string.
	CommandValueStringType CommandValueType = "string"
	// CommandValueIntType expects a decimal integer (digits only).
	CommandValueIntType CommandValueType = "int"
	// CommandValueBoolType expects a exact "true" or "false".
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
// It receives the message context and injected application data.
// Returning a non-nil error routes it through the bot's error handler.
type CommandExecutor[T AppData] func(ctx *MsgContext, dbContext T) error

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
	prefix    string
	separator string

	middlewares extypes.Slice[Middleware[T]]
	commands    extypes.Slice[*Command[T]]
}

// NewCommandGroup creates a command group that prefixes every added command.
func NewCommandGroup[T any](prefix string) *CommandGroup[T] {
	return &CommandGroup[T]{
		prefix: prefix, separator: "",

		middlewares: make([]Middleware[T], 0),
		commands:    make([]*Command[T], 0),
	}
}

// SetSeparator sets the text inserted between the group prefix and command name.
func (g *CommandGroup[T]) SetSeparator(separator string) *CommandGroup[T] {
	g.separator = separator
	return g
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
	newCmd.command = fmt.Sprintf("%s%s%s", g.prefix, g.separator, cmd.command)
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
