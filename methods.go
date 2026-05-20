package laniakea

import (
	"context"
	"encoding/json"
	"iter"

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
//  2. Requests updates with the timeout configured via PollTimeout
//  3. Filters updates by types specified in bot.GetUpdateTypes()
//  4. Logs raw update JSON if RequestLogger is configured
//  5. Automatically updates the offset to the last received update ID + 1
//  6. Returns all received updates (empty slice if none)
//
// Note: This is a blocking call that waits up to the configured PollTimeout
// for new updates, unless ctx is canceled earlier. For non-blocking behavior,
// consider using webhooks instead.
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
	timeout := bot.pollTimeout
	params := tgapi.UpdateParams{
		Offset:         new(offset),
		Timeout:        new(timeout),
		AllowedUpdates: bot.GetUpdateTypes(),
	}

	zero := make([]tgapi.Update, 0)
	updates, err := bot.api.GetUpdatesWithContext(ctx, params)
	if err != nil {
		return zero, err
	}

	if bot.requestLogger != nil {
		for _, u := range updates {
			j, err := json.Marshal(u)
			if err != nil {
				bot.GetLogger().Error(err)
			}
			bot.requestLogger.Debugf("UPDATE %s\n", j)
		}
	}
	if len(updates) > 0 {
		bot.SetUpdateOffset(updates[len(updates)-1].UpdateID + 1)
	}
	if updates == nil {
		return zero, nil
	}
	return updates, nil
}

// UpdatesIter fetches updates once and yields each update in order.
//
// If fetching updates fails, the iterator yields the error once with a zero
// update and then stops.
func (bot *Bot[T]) UpdatesIter(ctx context.Context) iter.Seq2[tgapi.Update, error] {
	return func(yield func(tgapi.Update, error) bool) {
		updates, err := bot.Updates(ctx)
		if err != nil {
			yield(tgapi.Update{}, err)
			return
		}
		for _, u := range updates {
			if !yield(u, nil) {
				return
			}
		}
	}
}
