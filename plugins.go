package laniakea

import "log"

type CommandExecutor func(ctx *MsgContext, dbContext *DatabaseContext)

type PluginBuilder struct {
	name           string
	commands       map[string]*CommandExecutor
	payloads       map[string]*CommandExecutor
	updateListener *CommandExecutor
}

type Plugin struct {
	Name           string
	Commands       map[string]*CommandExecutor
	Payloads       map[string]*CommandExecutor
	UpdateListener *CommandExecutor
}

func NewPlugin(name string) *PluginBuilder {
	return &PluginBuilder{
		name:     name,
		commands: make(map[string]*CommandExecutor),
		payloads: make(map[string]*CommandExecutor),
	}
}

func (p *PluginBuilder) Command(f CommandExecutor, cmd ...string) *PluginBuilder {
	for _, c := range cmd {
		p.commands[c] = &f
	}
	return p
}

func (p *PluginBuilder) Payload(f CommandExecutor, payloads ...string) *PluginBuilder {
	for _, payload := range payloads {
		p.payloads[payload] = &f
	}
	return p
}

func (p *PluginBuilder) UpdateListener(listener CommandExecutor) *PluginBuilder {
	p.updateListener = &listener
	return p
}

func (p *PluginBuilder) Build() *Plugin {
	if len(p.commands) == 0 && len(p.payloads) == 0 {
		log.Println("no command or payloads")
	}
	plugin := &Plugin{
		Name:           p.name,
		Commands:       p.commands,
		Payloads:       p.payloads,
		UpdateListener: p.updateListener,
	}
	return plugin
}

func (p *Plugin) Execute(cmd string, ctx *MsgContext, dbContext *DatabaseContext) {
	(*p.Commands[cmd])(ctx, dbContext)
}

func (p *Plugin) ExecutePayload(payload string, ctx *MsgContext, dbContext *DatabaseContext) {
	(*p.Payloads[payload])(ctx, dbContext)
}

type Middleware struct {
	Name     string
	Executor CommandExecutor
	Order    int
	Async    bool
}
type MiddlewareBuilder struct {
	name     string
	executor CommandExecutor
	order    int
	async    bool
}

func NewMiddleware(name string) *MiddlewareBuilder {
	return &MiddlewareBuilder{name: name, async: false}
}
func (m *MiddlewareBuilder) SetName(name string) *MiddlewareBuilder {
	m.name = name
	return m
}
func (m *MiddlewareBuilder) SetExecutor(executor CommandExecutor) *MiddlewareBuilder {
	m.executor = executor
	return m
}
func (m *MiddlewareBuilder) SetOrder(order int) *MiddlewareBuilder {
	m.order = order
	return m
}
func (m *MiddlewareBuilder) SetAsync(async bool) *MiddlewareBuilder {
	m.async = async
	return m
}
func (m *MiddlewareBuilder) Build() Middleware {
	return Middleware{
		Name:     m.name,
		Executor: m.executor,
		Order:    m.order,
		Async:    m.async,
	}
}
func (m Middleware) Execute(ctx *MsgContext, db *DatabaseContext) {
	if m.Async {
		go m.Executor(ctx, db)
	} else {
		m.Execute(ctx, db)
	}
}
