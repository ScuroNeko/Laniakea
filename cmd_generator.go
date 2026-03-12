// Package laniakea provides a framework for building Telegram bots with plugin-based
// command registration and automatic command scope management.
//
// This module automatically generates and registers bot commands across different
// chat scopes (private, group, admin) based on plugin-defined commands.
//
// Commands are derived from Plugin and Command structs, with optional descriptions
// and argument formatting. Automatic registration avoids manual command setup and
// ensures consistency across chat types.
package laniakea

import (
	"errors"
	"fmt"
	"strings"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
)

// ErrTooManyCommands is returned when the total number of registered commands
// exceeds Telegram's limit of 100 bot commands per bot.
//
// Telegram Bot API enforces this limit strictly. If exceeded, SetMyCommands
// will fail with a 400 error. This error helps catch the issue early during
// bot initialization.
var ErrTooManyCommands = errors.New("too many commands. max 100")

// generateBotCommand converts a Command[T] into a tgapi.BotCommand with a
// formatted description that includes usage instructions.
//
// The description is built as:
//
//	"<original_description>. Usage: /<command> <arg1> [<arg2>] ..."
//
// Required arguments are shown as-is; optional arguments are wrapped in square brackets.
//
// Example:
//
//	Command{command: "start", description: "Start the bot", args: []Arg{{text: "name", required: false}}}
//	→ Description: "Start the bot. Usage: /start [name]"
func generateBotCommand[T any](cmd Command[T]) tgapi.BotCommand {
	desc := cmd.command
	if len(cmd.description) > 0 {
		desc = cmd.description
	}

	var descArgs []string
	for _, a := range cmd.args {
		if a.required {
			descArgs = append(descArgs, a.text)
		} else {
			descArgs = append(descArgs, fmt.Sprintf("[%s]", a.text))
		}
	}

	desc = fmt.Sprintf("%s. Usage: /%s %s", desc, cmd.command, strings.Join(descArgs, " "))
	return tgapi.BotCommand{Command: cmd.command, Description: desc}
}

// generateBotCommandForPlugin collects all non-skipped commands from a Plugin[T]
// and converts them into tgapi.BotCommand objects.
//
// Commands marked with skipAutoCmd = true are excluded from auto-registration.
// This allows plugins to opt out of automatic command generation (e.g., for
// internal or hidden commands).
func generateBotCommandForPlugin[T any](pl Plugin[T]) []tgapi.BotCommand {
	commands := make([]tgapi.BotCommand, 0)
	for _, cmd := range pl.commands {
		if cmd.skipAutoCmd {
			continue
		}
		commands = append(commands, generateBotCommand(cmd))
	}
	return commands
}

// AutoGenerateCommands registers all plugin-defined commands with Telegram's Bot API
// across three scopes:
//   - Private chats (users)
//   - Group chats
//   - Group administrators
//
// It first deletes existing commands to ensure a clean state, then sets the new
// set of commands for all scopes. This ensures consistency even if commands were
// previously modified manually via @BotFather.
//
// Returns ErrTooManyCommands if the total number of commands exceeds 100.
// Returns any API error from Telegram (e.g., network issues, invalid scope).
//
// Important: This method assumes the bot has been properly initialized and
// the API client is authenticated and ready.
//
// Usage:
//
//	err := bot.AutoGenerateCommands()
//	if err != nil {
//	    log.Fatal(err)
//	}
func (bot *Bot[T]) AutoGenerateCommands() error {
	// Clear existing commands to avoid duplication or stale entries
	_, err := bot.api.DeleteMyCommands(tgapi.DeleteMyCommandsP{})
	if err != nil {
		return fmt.Errorf("failed to delete existing commands: %w", err)
	}

	// Collect all non-skipped commands from all plugins
	commands := make([]tgapi.BotCommand, 0)
	for _, pl := range bot.plugins {
		if pl.skipAutoCmd {
			continue
		}
		commands = append(commands, generateBotCommandForPlugin(pl)...)
		bot.logger.Debugf("Registered %d commands from plugin %s", len(pl.commands), pl.name)
	}

	// Enforce Telegram's 100-command limit
	if len(commands) > 100 {
		return ErrTooManyCommands
	}

	// Register commands for each scope
	scopes := []*tgapi.BotCommandScope{
		{Type: tgapi.BotCommandScopePrivateType},
		{Type: tgapi.BotCommandScopeGroupType},
		{Type: tgapi.BotCommandScopeAllChatAdministratorsType},
	}

	for _, scope := range scopes {
		_, err = bot.api.SetMyCommands(tgapi.SetMyCommandsP{
			Commands: commands,
			Scope:    scope,
		})
		if err != nil {
			return fmt.Errorf("failed to set commands for scope %q: %w", scope.Type, err)
		}
	}

	return nil
}
