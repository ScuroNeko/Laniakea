package laniakea

import (
	"regexp"

	"git.nix13.pw/scuroneko/extypes"
)

type CommandExecutor func(ctx *MsgContext, dbContext *DatabaseContext)

const (
	CommandValueStringType CommandValueType = "string"
	CommandValueIntType    CommandValueType = "int"
	CommandValueBoolType   CommandValueType = "bool"
	CommandValueAnyType    CommandValueType = "any"
)

var (
	CommandRegexInt    = regexp.MustCompile("\\d+")
	CommandRegexString = regexp.MustCompile("\\.+")
)

type CommandValueType string
type CommandArg struct {
	valueType CommandValueType
	text      string
	regex     *regexp.Regexp
}

func NewCommandArg(text string, valueType CommandValueType) CommandArg {
	regex := CommandRegexString
	switch valueType {
	case CommandValueIntType:
		regex = CommandRegexInt
	}
	return CommandArg{valueType, text, regex}
}

type Command struct {
	command     string
	exec        CommandExecutor
	args        []CommandArg
	middlewares []Middleware
}

func NewCommand(exec CommandExecutor, command string, args ...CommandArg) *Command {
	return &Command{command, exec, args, make([]Middleware, 0)}
}

type Plugin struct {
	Name        string
	Commands    map[string]Command
	Payloads    map[string]Command
	Middlewares extypes.Slice[PluginMiddleware]
}

func NewPlugin(name string) *Plugin {
	return &Plugin{
		Name:        name,
		Commands:    map[string]Command{},
		Payloads:    map[string]Command{},
		Middlewares: extypes.Slice[PluginMiddleware]{},
	}
}

func (p *Plugin) AddCommand(command Command) *Plugin {
	p.Commands[command.command] = command
	return p
}

func (p *Plugin) AddPayload(command Command) *Plugin {
	p.Payloads[command.command] = command
	return p
}

func (p *Plugin) AddMiddleware(middleware PluginMiddleware) *Plugin {
	p.Middlewares = p.Middlewares.Push(middleware)
	return p
}

func (p *Plugin) Execute(cmd string, ctx *MsgContext, dbContext *DatabaseContext) {
	command := p.Commands[cmd]
	if !command.validateArgs(ctx.Args) {
		return
	}
	command.exec(ctx, dbContext)
}
func (p *Plugin) ExecutePayload(payload string, ctx *MsgContext, dbContext *DatabaseContext) {
	pl := p.Payloads[payload]
	if !pl.validateArgs(ctx.Args) {
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
func (c *Command) validateArgs(args []string) bool {
	if len(args) != len(c.args) {
		return false
	}

	for i, arg := range c.args {
		if arg.regex == nil {
			continue
		}
		if !arg.regex.MatchString(args[i]) {
			return false
		}
	}
	return true
}

type Middleware struct {
	name  string
	exec  CommandExecutor
	order int
	async bool
}

func NewMiddleware(name string, executor CommandExecutor) *Middleware {
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
func (m *Middleware) Execute(ctx *MsgContext, db *DatabaseContext) {
	if m.async {
		go m.exec(ctx, db)
	} else {
		m.Execute(ctx, db)
	}
}

type PluginMiddlewareExecutor func(ctx *MsgContext, db *DatabaseContext) bool

// PluginMiddleware
// When async, returned value ignored
type PluginMiddleware struct {
	executor PluginMiddlewareExecutor
	order    int
	async    bool
}

func NewPluginMiddleware(executor PluginMiddlewareExecutor) *PluginMiddleware {
	return &PluginMiddleware{
		executor: executor,
		order:    0,
		async:    false,
	}
}
func (m *PluginMiddleware) SetOrder(order int) *PluginMiddleware {
	m.order = order
	return m
}
func (m *PluginMiddleware) SetAsync(async bool) *PluginMiddleware {
	m.async = async
	return m
}
func (m *PluginMiddleware) Execute(ctx *MsgContext, db *DatabaseContext) bool {
	if m.async {
		go m.executor(ctx, db)
		return true
	}
	return m.executor(ctx, db)
}
