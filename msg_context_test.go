package laniakea

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"git.nix13.pw/scuroneko/laniakea/tgapi"
	"git.nix13.pw/scuroneko/slog"
)

func TestAnswerPhotoIncludesDirectMessagesTopicID(t *testing.T) {
	var gotBody map[string]any

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("failed to read request body: %v", err)
			}
			if err := json.Unmarshal(body, &gotBody); err != nil {
				t.Fatalf("failed to decode request body: %v", err)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":{"message_id":9,"date":1}}`)),
			}, nil
		}),
	}

	api := tgapi.NewAPI(
		tgapi.NewAPIOpts("token").
			SetAPIUrl("https://example.test").
			SetHTTPClient(client),
	)
	defer func() {
		if err := api.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()

	ctx := &MsgContext{
		Api: api,
		Msg: &tgapi.Message{
			Chat:               &tgapi.Chat{ID: 42, Type: string(tgapi.ChatTypePrivate)},
			DirectMessageTopic: &tgapi.DirectMessageTopic{TopicID: 77},
		},
		Logger: slog.CreateLogger(),
	}

	answer := ctx.AnswerPhoto("photo-id", "caption")
	if answer == nil {
		t.Fatal("expected answer message")
	}
	if answer.MessageID != 9 {
		t.Fatalf("unexpected message id: %d", answer.MessageID)
	}
	if got := gotBody["direct_messages_topic_id"]; got != float64(77) {
		t.Fatalf("unexpected direct_messages_topic_id: %v", got)
	}
}
