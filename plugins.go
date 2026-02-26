package laniakea

import (
	"errors"
	"regexp"

	"git.nix13.pw/scuroneko/extypes"
)

const (
	CommandValueStringType CommandValueType = "string"
	CommandValueIntType    CommandValueType = "int"
	CommandValueBoolType   CommandValueType = "bool"
	CommandValueAnyType    CommandValueType = "any"
)

var (
	CommandRegexInt    = regexp.MustCompile("\\d+")
	CommandRegexString = regexp.MustCompile(".+")
)

var (
	ErrCmdArgCountMismatch  = errors.New("command arg count mismatch")
	ErrCmdArgRegexpMismatch = errors.New("command arg regexp mismatch")
)

type CommandValueType string
type CommandArg struct {
	valueType CommandValueType
	text      string
	regex     *regexp.Regexp
	required  bool
}

func NewCommandArg(text string, valueType CommandValueType) *CommandArg {
	regex := CommandRegexString
	switch valueType {
	case CommandValueIntType:
		regex = CommandRegexInt
	}
	return &CommandArg{valueType, text, regex, false}
}
func (c *CommandArg) SetRequired() *CommandArg {
	c.required = true
	return c
}

type CommandExecutor[T DbContext] func(ctx *MsgContext, dbContext *T)

type Command[T DbContext] struct {
	command     string
	description string
	exec        CommandExecutor[T]
	args        extypes.Slice[CommandArg]
	middlewares extypes.Slice[Middleware[T]]
	skipAutoCmd bool
}

func NewCommand[T any](exec CommandExecutor[T], command string, args ...CommandArg) *Command[T] {
	return &Command[T]{command, "", exec, args, make(extypes.Slice[Middleware[T]], 0), false}
}
func (c *Command[T]) Use(m Middleware[T]) *Command[T] {
	c.middlewares = c.middlewares.Push(m)
	return c
}
func (c *Command[T]) SetDescription(desc string) *Command[T] {
	c.description = desc
	return c
}
func (c *Command[T]) SkipCommandAutoGen() *Command[T] {
	c.skipAutoCmd = true
	return c
}
func (c *Command[T]) validateArgs(args []string) error {
	cmdArgs := c.args.Filter(func(e CommandArg) bool { return !e.required })
	if len(args) < cmdArgs.Len() {
		return ErrCmdArgCountMismatch
	}

	for i, arg := range args {
		if i >= c.args.Len() {
			break
		}
		cmdArg := c.args.Get(i)
		if cmdArg.regex == nil {
			continue
		}
		if !cmdArg.regex.MatchString(arg) {
			return ErrCmdArgRegexpMismatch
		}
	}
	return nil
}

type Plugin[T DbContext] struct {
	name        string
	commands    map[string]Command[T]
	payloads    map[string]Command[T]
	middlewares extypes.Slice[Middleware[T]]
	skipAutoCmd bool
}

func NewPlugin[T DbContext](name string) *Plugin[T] {
	return &Plugin[T]{
		name, map[string]Command[T]{},
		map[string]Command[T]{}, extypes.Slice[Middleware[T]]{}, false,
	}
}

func (p *Plugin[T]) AddCommand(command *Command[T]) *Plugin[T] {
	p.commands[command.command] = *command
	return p
}
func (p *Plugin[T]) NewCommand(exec CommandExecutor[T], command string, args ...CommandArg) *Command[T] {
	return NewCommand(exec, command, args...)
}
func (p *Plugin[T]) AddPayload(command *Command[T]) *Plugin[T] {
	p.payloads[command.command] = *command
	return p
}
func (p *Plugin[T]) AddMiddleware(middleware Middleware[T]) *Plugin[T] {
	p.middlewares = p.middlewares.Push(middleware)
	return p
}
func (p *Plugin[T]) SkipCommandAutoGen() *Plugin[T] {
	p.skipAutoCmd = true
	return p
}

func (p *Plugin[T]) executeCmd(cmd string, ctx *MsgContext, dbContext *T) {
	command := p.commands[cmd]
	if err := command.validateArgs(ctx.Args); err != nil {
		ctx.error(err)
		return
	}
	command.exec(ctx, dbContext)
}
func (p *Plugin[T]) executePayload(payload string, ctx *MsgContext, dbContext *T) {
	pl := p.payloads[payload]
	if err := pl.validateArgs(ctx.Args); err != nil {
		ctx.error(err)
		return
	}
	pl.exec(ctx, dbContext)
}
func (p *Plugin[T]) executeMiddlewares(ctx *MsgContext, db *T) bool {
	for _, m := range p.middlewares {
		if !m.Execute(ctx, db) {
			return false
		}
	}
	return true
}

type MiddlewareExecutor[T DbContext] func(ctx *MsgContext, db *T) bool

// Middleware
// When async, returned value ignored
type Middleware[T DbContext] struct {
	name     string
	executor MiddlewareExecutor[T]
	order    int
	async    bool
}

func NewMiddleware[T DbContext](name string, executor MiddlewareExecutor[T]) *Middleware[T] {
	return &Middleware[T]{name, executor, 0, false}
}
func (m *Middleware[T]) SetOrder(order int) *Middleware[T] {
	m.order = order
	return m
}
func (m *Middleware[T]) SetAsync(async bool) *Middleware[T] {
	m.async = async
	return m
}
func (m *Middleware[T]) Execute(ctx *MsgContext, db *T) bool {
	if m.async {
		go m.executor(ctx, db)
		return true
	}
	return m.executor(ctx, db)
}
