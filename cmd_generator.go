package laniakea

import "git.nix13.pw/scuroneko/laniakea/tgapi"

func generateBotCommand(cmd Command) tgapi.BotCommand {
	return tgapi.BotCommand{
		Command: cmd.command, Description: cmd.command,
	}
}

func generateBotCommandForPlugin(pl Plugin) []tgapi.BotCommand {
	cmds := make([]tgapi.BotCommand, 0)
	for _, cmd := range pl.Commands {
		cmds = append(cmds, generateBotCommand(cmd))
	}
	return cmds
}

func (b *Bot) AutoGenerateCommands() error {
	commands := make([]tgapi.BotCommand, 0)
	for _, pl := range b.plugins {
		commands = append(commands, generateBotCommandForPlugin(pl)...)
	}
	_, err := b.api.SetMyCommands(tgapi.SetMyCommandsP{Commands: commands})
	return err
}
