package laniakea

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
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

func TestAnswerRejectsEmptyMessage(t *testing.T) {
	ctx := &MsgContext{
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 42, Type: string(tgapi.ChatTypePrivate)}},
		Logger: slog.CreateLogger(),
	}

	if answer := ctx.Answer(""); answer != nil {
		t.Fatal("expected nil answer for empty message")
	}
}

func TestAnswerRejectsLongMessageWithoutSendingRequest(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatal("unexpected HTTP request")
			return nil, nil
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
		Api:    api,
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 42, Type: string(tgapi.ChatTypePrivate)}},
		Logger: slog.CreateLogger(),
	}

	if answer := ctx.Answer(strings.Repeat("a", maxMessageTextLen+1)); answer != nil {
		t.Fatal("expected nil answer for long message")
	}
}

func TestValidateMessageText(t *testing.T) {
	if err := validateMessageText(""); !errors.Is(err, ErrEmptyMessage) {
		t.Fatalf("expected ErrEmptyMessage, got %v", err)
	}
	if err := validateMessageText(strings.Repeat("a", maxMessageTextLen+1)); !errors.Is(err, ErrMessageTooLong) {
		t.Fatalf("expected ErrMessageTooLong, got %v", err)
	}
	if err := validateMessageText("ok"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCaptionText(t *testing.T) {
	if err := validateCaptionText(strings.Repeat("a", maxMessageCaptionLen+1)); !errors.Is(err, ErrCaptionTooLong) {
		t.Fatalf("expected ErrCaptionTooLong, got %v", err)
	}
	if err := validateCaptionText(""); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestSplitMessageTextPreservesContent(t *testing.T) {
	text := "alpha beta\n" + strings.Repeat("x", maxMessageTextLen) + " omega"

	parts := SplitMessageText(text)
	if len(parts) < 2 {
		t.Fatalf("expected multiple parts, got %d", len(parts))
	}

	for i, part := range parts {
		if got := len([]rune(part)); got > maxMessageTextLen {
			t.Fatalf("part %d exceeded limit: %d", i, got)
		}
	}

	if got := strings.Join(parts, ""); got != text {
		t.Fatalf("split/join mismatch: got %q want %q", got, text)
	}
}

func TestAnswerLongSplitsRequestsAndAttachesKeyboardToLastChunk(t *testing.T) {
	var requests []map[string]any

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("failed to read request body: %v", err)
			}
			var got map[string]any
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("failed to decode request body: %v", err)
			}
			requests = append(requests, got)
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
		Api:    api,
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 42, Type: string(tgapi.ChatTypePrivate)}},
		Logger: slog.CreateLogger(),
	}
	kb := NewInlineKeyboardJson(1).AddCallbackButton("A", "cmd")
	text := strings.Repeat("a", maxMessageTextLen) + " " + strings.Repeat("b", 32)

	messages := ctx.KeyboardLong(text, kb)
	if got := len(messages); got != 2 {
		t.Fatalf("expected 2 sent messages, got %d", got)
	}
	if got := len(requests); got != 2 {
		t.Fatalf("expected 2 requests, got %d", got)
	}
	if _, ok := requests[0]["reply_markup"]; ok {
		t.Fatal("did not expect keyboard on first chunk")
	}
	if _, ok := requests[1]["reply_markup"]; !ok {
		t.Fatal("expected keyboard on final chunk")
	}

	gotTexts := []string{requests[0]["text"].(string), requests[1]["text"].(string)}
	wantTexts := SplitMessageText(text)
	if !reflect.DeepEqual(gotTexts, wantTexts) {
		t.Fatalf("unexpected chunk texts: got %q want %q", gotTexts, wantTexts)
	}
}
