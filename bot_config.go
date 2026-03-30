package laniakea

import (
	"fmt"
	"slices"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/slog"
)

// AddPrefixes adds one or more command prefixes (e.g., "/", "!").
// Must have at least one prefix before Run().
func (bot *Bot[T]) AddPrefixes(prefixes ...string) *Bot[T] {
	if !bot.configMutable("AddPrefixes") {
		return bot
	}
	bot.prefixes = append(bot.prefixes, prefixes...)
	return bot
}

// SetDraftProvider replaces the default DraftProvider with a custom one.
// Useful for using LinearDraftIdGenerator to persist draft IDs across restarts.
func (bot *Bot[T]) SetDraftProvider(p *DraftProvider) *Bot[T] {
	if !bot.configMutable("SetDraftProvider") {
		return bot
	}
	bot.draftProvider = p
	return bot
}

// GetDraftProvider returns the draft provider currently used by the bot.
func (bot *Bot[T]) GetDraftProvider() *DraftProvider {
	return bot.draftProvider
}

// SetSessionStore replaces the session store used for scene management.
func (bot *Bot[T]) SetSessionStore(store SessionStore) *Bot[T] {
	if !bot.configMutable("SetSessionStore") {
		return bot
	}
	if store == nil {
		bot.logger.Warn("SetSessionStore called with nil store; using default MemorySessionStore")
		return bot
	}
	bot.sessionStore = store
	return bot
}

// GetSessionStore returns the session store used for scene management.
func (bot *Bot[T]) GetSessionStore() SessionStore {
	return bot.sessionStore
}

// SetSceneScopePriority sets the lookup order for resolving active scene sessions.
func (bot *Bot[T]) SetSceneScopePriority(priority []SceneScope) *Bot[T] {
	if !bot.configMutable("SetSceneScopePriority") {
		return bot
	}
	newPriority := make([]SceneScope, 0, 3)
	for _, scope := range priority {
		if scope != SceneScopeUser && scope != SceneScopeChat && scope != SceneScopeUserChat {
			bot.logger.Warnln(fmt.Sprintf("invalid scene scope %v in priority list; ignoring", scope))
			continue
		}
		if slices.Index(newPriority, scope) >= 0 {
			bot.logger.Warnln(fmt.Sprintf("duplicate scope %v in scene scope priority; ignoring duplicates", scope))
			continue
		}
		newPriority = append(newPriority, scope)
	}
	if len(newPriority) == 0 || len(newPriority) > 3 {
		bot.logger.Warnln("scene scope priority must have 1 to 3 scopes; ignoring invalid input")
		return bot
	}
	bot.sceneScopePriority = append([]SceneScope(nil), newPriority...)
	return bot
}

// SetAppData injects shared application data into the bot.
//
// The data is accessible to commands, payload handlers, middleware, scenes,
// and runners through the generic type parameter T.
//
// For shared dependencies such as *sql.DB, prefer using a pointer type as T.
// Value-typed application data is supported, but the bot warns once because
// handlers receive T by value.
func (bot *Bot[T]) SetAppData(ctx T) *Bot[T] {
	if !bot.configMutable("SetAppData") {
		return bot
	}
	if !bot.warnedValueData && shouldWarnOnValueAppData[T]() && bot.logger != nil {
		bot.logger.Warnln("app data uses a value type; shared dependencies should usually use a pointer type as T")
		bot.warnedValueData = true
	}
	bot.appData = ctx
	bot.hasAppData = true
	return bot
}

// GetAppData returns the injected application data.
// If SetAppData was not called, it returns the zero value of T.
func (bot *Bot[T]) GetAppData() T { return bot.appData }

// SetUpdateTypes sets the list of update types the bot will request from Telegram.
// Overwrites any previously set types.
func (bot *Bot[T]) SetUpdateTypes(t ...tgapi.UpdateType) *Bot[T] {
	if !bot.configMutable("UpdateTypes") {
		return bot
	}
	bot.updateTypes = make([]tgapi.UpdateType, 0)
	bot.updateTypes = append(bot.updateTypes, t...)
	return bot
}

// AddUpdateType adds one or more update types to the list.
// Does not overwrite existing types.
func (bot *Bot[T]) AddUpdateType(t ...tgapi.UpdateType) *Bot[T] {
	if !bot.configMutable("AddUpdateType") {
		return bot
	}
	bot.updateTypes = append(bot.updateTypes, t...)
	return bot
}

// GetUpdateTypes returns the list of update types the bot is configured to receive.
func (bot *Bot[T]) GetUpdateTypes() []tgapi.UpdateType {
	return append([]tgapi.UpdateType(nil), bot.updateTypes...)
}

// SetPayloadType sets the default payload encoding type used for callback data.
// JSON stores payload as a string: `{"cmd":"command","args":[...]}`.
// Base64 stores the same JSON encoded as a Base64URL string.
// InlineKeyboard.SetPayloadType may override this value for an individual keyboard.
func (bot *Bot[T]) SetPayloadType(t BotPayloadType) *Bot[T] {
	if !bot.configMutable("SetPayloadType") {
		return bot
	}
	bot.payloadType = t
	return bot
}

// GetPayloadType returns the bot's default callback payload encoding type.
func (bot *Bot[T]) GetPayloadType() BotPayloadType { return bot.payloadType }

// SetStrictPayloadType enables or disables strict callback payload decoding.
// When enabled, callback payloads must match the bot's default payload type.
func (bot *Bot[T]) SetStrictPayloadType(strict bool) *Bot[T] {
	if !bot.configMutable("SetStrictPayloadType") {
		return bot
	}
	bot.strictPayloadType = strict
	return bot
}

// SetErrorTemplate sets the format string for error messages sent to users.
// Use "%s" to insert the error message.
// Example: "❌ Error: %s" → "❌ Error: Command not found".
func (bot *Bot[T]) SetErrorTemplate(s string) *Bot[T] {
	if !bot.configMutable("ErrorTemplate") {
		return bot
	}
	bot.errorTemplate = s
	return bot
}

// SetDebug enables or disables debug logging.
func (bot *Bot[T]) SetDebug(debug bool) *Bot[T] {
	bot.debug = debug
	level := slog.FATAL
	if debug {
		level = slog.DEBUG
	}

	bot.logger.Level(level)
	if bot.RequestLogger != nil {
		bot.RequestLogger.Level(level)
	}
	for _, p := range bot.plugins {
		if p.logger == nil {
			continue
		}
		p.logger.Level(level)
	}
	return bot
}

// SetL10n sets the localization (i18n) provider for the bot.
//
// The L10n instance must be pre-populated with translations.
// Translations are accessed via Bot.L10n(lang, key).
//
// Replaces any previously set L10n instance.
func (bot *Bot[T]) SetL10n(l *L10n) *Bot[T] {
	if !bot.configMutable("SetL10n") {
		return bot
	}
	if l == nil {
		bot.logger.Warn("SetL10n called with nil L10n; localization will be disabled")
		return bot
	}
	bot.l10n = l
	return bot
}
