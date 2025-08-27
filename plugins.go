package laniakea

type CommandExecutor func(ctx *MsgContext)

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
	if len(p.commands) == 0 {
		return nil
	}
	plugin := &Plugin{
		Name:           p.name,
		Commands:       p.commands,
		Payloads:       p.payloads,
		UpdateListener: p.updateListener,
	}
	return plugin
}

func (p *Plugin) Execute(cmd string, ctx *MsgContext) {
	(*p.Commands[cmd])(ctx)
}

func (p *Plugin) ExecutePayload(payload string, ctx *MsgContext) {
	(*p.Payloads[payload])(ctx)
}
