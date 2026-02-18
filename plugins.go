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

type CommandExecutor func(ctx *MsgContext, dbContext *DatabaseContext)
type Command struct {
	command     string
	description string
	exec        CommandExecutor
	args        extypes.Slice[CommandArg]
	middlewares extypes.Slice[Middleware]
	skipAutoCmd bool
}

func NewCommand(exec CommandExecutor, command string, args ...CommandArg) *Command {
	return &Command{command, "", exec, args, make(extypes.Slice[Middleware], 0), false}
}
func (c *Command) Use(m Middleware) *Command {
	c.middlewares = c.middlewares.Push(m)
	return c
}
func (c *Command) SetDescription(desc string) *Command {
	c.description = desc
	return c
}
func (c *Command) SkipCommandAutoGen() *Command {
	c.skipAutoCmd = true
	return c
}
func (c *Command) validateArgs(args []string) error {
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

type Plugin struct {
	Name        string
	Commands    map[string]Command
	Payloads    map[string]Command
	Middlewares extypes.Slice[Middleware]
}

func NewPlugin(name string) *Plugin {
	return &Plugin{
		name, map[string]Command{},
		map[string]Command{}, extypes.Slice[Middleware]{},
	}
}

func (p *Plugin) AddCommand(command *Command) *Plugin {
	p.Commands[command.command] = *command
	return p
}
func (p *Plugin) NewCommand(exec CommandExecutor, command string, args ...CommandArg) *Command {
	return NewCommand(exec, command, args...)
}
func (p *Plugin) AddPayload(command *Command) *Plugin {
	p.Payloads[command.command] = *command
	return p
}
func (p *Plugin) AddMiddleware(middleware Middleware) *Plugin {
	p.Middlewares = p.Middlewares.Push(middleware)
	return p
}

func (p *Plugin) executeCmd(cmd string, ctx *MsgContext, dbContext *DatabaseContext) {
	command := p.Commands[cmd]
	if err := command.validateArgs(ctx.Args); err != nil {
		ctx.error(err)
		return
	}
	command.exec(ctx, dbContext)
}
func (p *Plugin) executePayload(payload string, ctx *MsgContext, dbContext *DatabaseContext) {
	pl := p.Payloads[payload]
	if err := pl.validateArgs(ctx.Args); err != nil {
		ctx.error(err)
		return
	}
	pl.exec(ctx, dbContext)
}
func (p *Plugin) executeMiddlewares(ctx *MsgContext, db *DatabaseContext) bool {
	for _, m := range p.Middlewares {
		if !m.Execute(ctx, db) {
			return false
		}
	}
	return true
}

type MiddlewareExecutor func(ctx *MsgContext, db *DatabaseContext) bool

// Middleware
// When async, returned value ignored
type Middleware struct {
	name     string
	executor MiddlewareExecutor
	order    int
	async    bool
}

func NewMiddleware(name string, executor MiddlewareExecutor) *Middleware {
	return &Middleware{name, executor, 0, false}
}
func (m *Middleware) SetOrder(order int) *Middleware {
	m.order = order
	return m
}
func (m *Middleware) SetAsync(async bool) *Middleware {
	m.async = async
	return m
}
func (m *Middleware) Execute(ctx *MsgContext, db *DatabaseContext) bool {
	if m.async {
		go m.executor(ctx, db)
		return true
	}
	return m.executor(ctx, db)
}
