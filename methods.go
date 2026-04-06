package laniakea

import (
	"context"
	"encoding/json"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

// Updates fetches new updates from Telegram API using long polling.
// It respects the bot's current update offset and automatically advances it
// after successful retrieval. The method supports selective update types
// through AllowedUpdates and includes optional request logging.
//
// Parameters:
//   - ctx: request context used to cancel the in-flight long polling request
//
// Returns:
//   - []tgapi.Update: slice of received updates (empty if none available)
//   - error: any error encountered during the API call
//
// Behavior:
//  1. Uses the bot's current update offset (via GetUpdateOffset)
//  2. Requests updates with 30-second timeout
//  3. Filters updates by types specified in bot.GetUpdateTypes()
//  4. Logs raw update JSON if RequestLogger is configured
//  5. Automatically updates the offset to the last received update ID + 1
//  6. Returns all received updates (empty slice if none)
//
// Note: This is a blocking call that waits up to 30 seconds for new updates,
// unless ctx is canceled earlier. For non-blocking behavior, consider using
// webhooks instead.
//
// Example:
//
//	updates, err := bot.Updates(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, update := range updates {
//	    // process update
//	}
func (bot *Bot[T]) Updates(ctx context.Context) ([]tgapi.Update, error) {
	offset := bot.GetUpdateOffset()
	params := tgapi.UpdateParams{
		Offset:         new(offset),
		Timeout:        new(30),
		AllowedUpdates: bot.GetUpdateTypes(),
	}

	updates, err := bot.api.GetUpdatesWithContext(ctx, params)
	if err != nil {
		return nil, err
	}

	if bot.RequestLogger != nil {
		for _, u := range updates {
			j, err := json.Marshal(u)
			if err != nil {
				bot.GetLogger().Error(err)
			}
			bot.RequestLogger.Debugf("UPDATE %s\n", j)
		}
	}
	if len(updates) > 0 {
		bot.SetUpdateOffset(updates[len(updates)-1].UpdateID + 1)
	}
	return updates, err
}
