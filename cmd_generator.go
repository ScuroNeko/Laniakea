package laniakea

import (
	"fmt"
	"strings"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

func generateBotCommand(cmd Command) tgapi.BotCommand {
	desc := cmd.command
	if len(cmd.description) > 0 {
		desc = cmd.description
	}
	var descArgs []string
	for _, a := range cmd.args {
		if a.required {
			descArgs = append(descArgs, fmt.Sprintf("%s", a.text))
		} else {
			descArgs = append(descArgs, fmt.Sprintf("[%s]", a.text))
		}
	}
	desc = fmt.Sprintf("%s. Usage: /%s %s", desc, cmd.command, strings.Join(descArgs, " "))
	return tgapi.BotCommand{Command: cmd.command, Description: desc}
}

func generateBotCommandForPlugin(pl Plugin) []tgapi.BotCommand {
	commands := make([]tgapi.BotCommand, 0)
	for _, cmd := range pl.Commands {
		if cmd.skipAutoCmd {
			continue
		}
		commands = append(commands, generateBotCommand(cmd))
	}
	return commands
}

func (b *Bot) AutoGenerateCommands() error {
	_, err := b.api.DeleteMyCommands(tgapi.DeleteMyCommandsP{})
	if err != nil {
		return err
	}

	commands := make([]tgapi.BotCommand, 0)
	for _, pl := range b.plugins {
		commands = append(commands, generateBotCommandForPlugin(pl)...)
	}

	privateChatsScope := &tgapi.BotCommandScope{Type: tgapi.BotCommandScopePrivateType}
	groupChatsScope := &tgapi.BotCommandScope{Type: tgapi.BotCommandScopeGroupType}
	chatAdminsScope := &tgapi.BotCommandScope{Type: tgapi.BotCommandScopeAllChatAdministratorsType}
	_, err = b.api.SetMyCommands(tgapi.SetMyCommandsP{Commands: commands, Scope: privateChatsScope})
	if err != nil {
		return err
	}
	_, err = b.api.SetMyCommands(tgapi.SetMyCommandsP{Commands: commands, Scope: groupChatsScope})
	if err != nil {
		return err
	}
	_, err = b.api.SetMyCommands(tgapi.SetMyCommandsP{Commands: commands, Scope: chatAdminsScope})
	return err
}
