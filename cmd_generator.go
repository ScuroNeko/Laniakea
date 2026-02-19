package laniakea

import (
	"errors"
	"fmt"
	"strings"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

func generateBotCommand[T any](cmd Command[T]) tgapi.BotCommand {
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

func generateBotCommandForPlugin[T any](pl Plugin[T]) []tgapi.BotCommand {
	commands := make([]tgapi.BotCommand, 0)
	for _, cmd := range pl.Commands {
		if cmd.skipAutoCmd {
			continue
		}
		commands = append(commands, generateBotCommand(cmd))
	}
	return commands
}

var ErrTooManyCommands = errors.New("too many commands. max 100")

func (bot *Bot[T]) AutoGenerateCommands() error {
	_, err := bot.api.DeleteMyCommands(tgapi.DeleteMyCommandsP{})
	if err != nil {
		return err
	}

	commands := make([]tgapi.BotCommand, 0)
	for _, pl := range bot.plugins {
		commands = append(commands, generateBotCommandForPlugin(pl)...)
	}
	if len(commands) > 100 {
		return ErrTooManyCommands
	}

	privateChatsScope := &tgapi.BotCommandScope{Type: tgapi.BotCommandScopePrivateType}
	groupChatsScope := &tgapi.BotCommandScope{Type: tgapi.BotCommandScopeGroupType}
	chatAdminsScope := &tgapi.BotCommandScope{Type: tgapi.BotCommandScopeAllChatAdministratorsType}
	_, err = bot.api.SetMyCommands(tgapi.SetMyCommandsP{Commands: commands, Scope: privateChatsScope})
	if err != nil {
		return err
	}
	_, err = bot.api.SetMyCommands(tgapi.SetMyCommandsP{Commands: commands, Scope: groupChatsScope})
	if err != nil {
		return err
	}
	_, err = bot.api.SetMyCommands(tgapi.SetMyCommandsP{Commands: commands, Scope: chatAdminsScope})
	return err
}
