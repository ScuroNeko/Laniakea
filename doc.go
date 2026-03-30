/*
Package laniakea provides a modular, extensible framework for building scalable Telegram bots.

Core concepts:

  - Bot manages Telegram API access, update processing, logging, rate limiting, and dependency injection.
  - Plugins group commands, payloads, and non-command update handlers behind shared middleware.
  - MsgContext provides access to the current update and reply/edit/delete helpers.
  - InlineKeyboard builds callback-driven keyboards and structured payloads.
  - DraftProvider accumulates multi-step replies before sending them.
  - L10n stores key-based translations with fallback behavior.
  - Runners execute startup or background tasks alongside the polling loop.

Example usage:

	bot, err := laniakea.NewBot[*mydb.AppData](laniakea.LoadOptsFromEnv())
	if err != nil {
	    return err
	}
	bot.SetAppData(myDB).
	    AddUpdateType(tgapi.UpdateTypeMessage).
	    AddPrefixes("/", "!").
	    AddPlugins(&startPlugin, &helpPlugin).
	    AddMiddleware(authMiddleware, logMiddleware).
	    AddRunner(cleanupRunner).
	    SetL10n(l10n.New())

	return bot.Run()

Configure bots, plugins, and localization before starting Run or RunWithContext.
Runtime accessors are safe for concurrent use unless stated otherwise.
*/
package laniakea
