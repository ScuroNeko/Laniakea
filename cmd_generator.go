package laniakea

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

// CmdRegexp matches command names allowed for Telegram command registration.
var CmdRegexp = regexp.MustCompile("^[_a-z0-9]{1,32}$")

// ErrTooManyCommands is returned when the total number of registered commands
// exceeds Telegram's limit of 100 bot commands per bot.
//
// Telegram Bot API enforces this limit strictly. If exceeded, SetMyCommands
// will fail with a 400 error. This error helps catch the issue early during
// bot initialization.
var ErrTooManyCommands = errors.New("too many commands. max 100")

// Internal helper to build a BotCommand description with generated usage text.
func generateBotCommand[T any](cmd *Command[T]) tgapi.BotCommand {
	desc := ""
	if len(cmd.description) > 0 {
		desc = cmd.description
	}

	var descArgs []string
	for _, a := range cmd.args {
		if a.required {
			descArgs = append(descArgs, fmt.Sprintf("<%s>", a.text))
		} else {
			descArgs = append(descArgs, fmt.Sprintf("[%s]", a.text))
		}
	}

	usage := fmt.Sprintf("Usage: /%s %s", cmd.command, strings.Join(descArgs, " "))
	if desc != "" {
		desc = fmt.Sprintf("%s. %s", desc, usage)
		return tgapi.BotCommand{Command: cmd.command, Description: desc}
	}
	return tgapi.BotCommand{Command: cmd.command, Description: usage}
}

// Internal helper to validate Telegram command names.
func checkCmdRegex(cmd string) bool { return CmdRegexp.MatchString(cmd) }

// Internal helper to collect non-skipped, valid commands from one plugin.
func gatherCommandsForPlugin[T any](pl Plugin[T]) []tgapi.BotCommand {
	commands := make([]tgapi.BotCommand, 0)
	names := make([]string, 0, len(pl.commands))
	for name := range pl.commands {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		cmd := pl.commands[name]
		if cmd.skipAutoCmd {
			continue
		}
		if !checkCmdRegex(cmd.command) {
			continue
		}
		commands = append(commands, generateBotCommand(cmd))
	}
	return commands
}

// Internal helper to collect all auto-generated commands from registered plugins.
func gatherCommands[T any](bot *Bot[T]) []tgapi.BotCommand {
	commands := make([]tgapi.BotCommand, 0)
	for _, pl := range bot.plugins {
		if pl.skipAutoCmd {
			continue
		}
		commands = append(commands, gatherCommandsForPlugin(pl)...)
		bot.logger.Debugln(fmt.Sprintf("Registered %d commands from plugin %s", len(pl.commands), pl.name))
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
	commands := gatherCommands(bot)
	if len(commands) > 100 {
		return ErrTooManyCommands
	}

	// Clear existing commands to avoid duplication or stale entries
	_, err := bot.api.DeleteMyCommands(tgapi.DeleteMyCommandsP{})
	if err != nil {
		return fmt.Errorf("failed to delete existing commands: %w", err)
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

// AutoGenerateCommandsForScope registers all plugin-defined commands with Telegram's Bot API
// for the specified command scope. It first deletes any existing commands in that scope
// to ensure a clean state, then sets the new set of commands.
//
// The scope parameter defines where the commands should be available (e.g., private chats,
// group chats, chat administrators). See tgapi.BotCommandScope and its predefined types.
//
// Returns ErrTooManyCommands if the total number of commands exceeds 100.
// Returns any API error from Telegram (e.g., network issues, invalid scope).
//
// Usage:
//
//	privateScope := &tgapi.BotCommandScope{Type: tgapi.BotCommandScopePrivateType}
//	if err := bot.AutoGenerateCommandsForScope(privateScope); err != nil {
//	    log.Fatal(err)
//	}
func (bot *Bot[T]) AutoGenerateCommandsForScope(scope *tgapi.BotCommandScope) error {
	commands := gatherCommands(bot)
	if len(commands) > 100 {
		return ErrTooManyCommands
	}

	_, err := bot.api.DeleteMyCommands(tgapi.DeleteMyCommandsP{Scope: scope})
	if err != nil {
		return fmt.Errorf("failed to delete existing commands: %w", err)
	}

	_, err = bot.api.SetMyCommands(tgapi.SetMyCommandsP{
		Commands: commands,
		Scope:    scope,
	})
	if err != nil {
		return fmt.Errorf("failed to set commands for scope %q: %w", scope.Type, err)
	}
	return nil
}
