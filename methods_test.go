package laniakea

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

func TestUpdatesIterYieldsFetchError(t *testing.T) {
	bot := newUpdatesIterTestBot(t, `{"ok":false,"error_code":500,"description":"boom"}`)

	var gotErr error
	var gotUpdates int
	bot.UpdatesIter(context.Background())(func(update tgapi.Update, err error) bool {
		gotUpdates++
		if update.UpdateID != 0 {
			t.Fatalf("expected zero update on error, got %d", update.UpdateID)
		}
		gotErr = err
		return true
	})

	if gotUpdates != 1 {
		t.Fatalf("expected one yielded error, got %d yields", gotUpdates)
	}
	if gotErr == nil {
		t.Fatal("expected fetch error")
	}
	if !strings.Contains(gotErr.Error(), "boom") {
		t.Fatalf("expected Telegram error description, got %v", gotErr)
	}
}

func TestUpdatesIterStopsWhenYieldReturnsFalse(t *testing.T) {
	bot := newUpdatesIterTestBot(t, `{"ok":true,"result":[{"update_id":11},{"update_id":12}]}`)

	var gotIDs []int
	bot.UpdatesIter(context.Background())(func(update tgapi.Update, err error) bool {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		gotIDs = append(gotIDs, update.UpdateID)
		return false
	})

	if len(gotIDs) != 1 || gotIDs[0] != 11 {
		t.Fatalf("expected only first update, got %v", gotIDs)
	}
}

func newUpdatesIterTestBot(t *testing.T, response string) *Bot[NoData] {
	t.Helper()

	api := tgapi.NewAPI(
		tgapi.NewAPIOpts("token").
			SetAPIURL("https://example.test").
			SetHTTPClient(&http.Client{
				Transport: pollingRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(response)),
					}, nil
				}),
			}),
	)
	t.Cleanup(func() {
		if err := api.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	})

	return &Bot[NoData]{api: api}
}
