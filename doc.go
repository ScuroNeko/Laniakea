/*
Package laniakea provides a modular, extensible framework for building scalable Telegram bots.

It offers a fluent API for configuration and separates concerns through several core concepts:

  - Bot: The central instance managing API communication, update processing, logging,
    rate limiting, and dependency injection. Created via NewBot[T].

  - Plugins: Organize commands and payloads into reusable units.
    A plugin can have multiple commands and shared middlewares.

  - Commands: Named bot commands with descriptions, argument validation, and
    execution logic. Automatically registrable across different chat scopes.

  - Middleware: Functions that intercept and modify updates before they reach plugins.
    Useful for authentication, logging, validation, etc. Return false to stop processing.

  - MsgContext: Provides access to the incoming update and convenient methods for
    responding, editing, deleting, and translating messages. Includes built-in rate limiting
    and error handling. ⚠️ MarkdownV2 methods require manual escaping via EscapeMarkdownV2().

  - InlineKeyboard: A fluent builder for constructing inline keyboards with styled buttons,
    icons, URLs, and structured callback data (JSON or Base64).

  - DraftProvider: Manages ephemeral, multi-step message drafts with automatic ID generation
    (random or linear). Drafts can be built incrementally and flushed atomically.

  - L10n: Simple key-based localization system with fallback language support.

  - Runners: Background goroutines for periodic tasks or one‑off initialization,
    with configurable timeouts and async execution.

  - RateLimiting & Logging: Built‑in rate limiter (respects Telegram's retry_after)
    and structured logging (JSON stdout + optional file output) with request‑level tracing.

  - Dependency Injection: Pass any custom database context (e.g., *sql.DB) to all handlers
    via the type parameter T in Bot[T].

Example usage:

	bot := laniakea.NewBot[mydb.DBContext](laniakea.LoadOptsFromEnv()).
	    DatabaseContext(&myDB).
	    AddUpdateType(tgapi.UpdateTypeMessage).
	    AddPrefixes("/", "!").
	    AddPlugins(&startPlugin, &helpPlugin).
	    AddMiddleware(&authMiddleware, &logMiddleware).
	    AddRunner(&cleanupRunner).
	    AddL10n(l10n.New())

	bot.Run()

All public methods are safe for concurrent use unless stated otherwise.
Direct field access is not recommended; use provided accessors (e.g., GetDBContext, SetUpdateOffset).
*/
package laniakea
