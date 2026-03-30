package laniakea

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/slog"
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
			Chat:               &tgapi.Chat{ID: 42, Type: tgapi.ChatTypePrivate},
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

func TestBindArgsBindsScalarFields(t *testing.T) {
	type input struct {
		ID     int
		Active bool
		Score  float64
		Name   string
	}

	ctx := &MsgContext{Args: []string{"42", "true", "3.5", "Ada", "Lovelace"}}
	var got input

	if err := ctx.BindArgs(&got); err != nil {
		t.Fatalf("BindArgs returned error: %v", err)
	}

	want := input{
		ID:     42,
		Active: true,
		Score:  3.5,
		Name:   "Ada Lovelace",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected bound value: got %#v want %#v", got, want)
	}
}

func TestBindArgsLeavesTrailingFieldsZeroWhenArgsRunOut(t *testing.T) {
	type input struct {
		ID     int
		Reason string
		Admin  bool
	}

	ctx := &MsgContext{Args: []string{"7"}}
	var got input

	if err := ctx.BindArgs(&got); err != nil {
		t.Fatalf("BindArgs returned error: %v", err)
	}

	if got.ID != 7 {
		t.Fatalf("unexpected ID: got %d want 7", got.ID)
	}
	if got.Reason != "" {
		t.Fatalf("expected zero-value Reason, got %q", got.Reason)
	}
	if got.Admin {
		t.Fatal("expected zero-value Admin")
	}
}

func TestBindArgsRejectsInvalidTargets(t *testing.T) {
	ctx := &MsgContext{Args: []string{"1"}}

	if err := ctx.BindArgs(nil); !errors.Is(err, ErrBindArgsTargetNotPointer) {
		t.Fatalf("expected ErrBindArgsTargetNotPointer for nil target, got %v", err)
	}

	var notStruct int
	if err := ctx.BindArgs(&notStruct); !errors.Is(err, ErrBindArgsTargetNotStruct) {
		t.Fatalf("expected ErrBindArgsTargetNotStruct for non-struct target, got %v", err)
	}
}

func TestBindArgsReportsConversionFailures(t *testing.T) {
	type input struct {
		ID int
	}

	ctx := &MsgContext{Args: []string{"oops"}}
	var got input

	err := ctx.BindArgs(&got)
	if err == nil {
		t.Fatal("expected BindArgs to fail")
	}
	if !errors.Is(err, ErrBindArgsConversion) {
		t.Fatalf("expected ErrBindArgsConversion, got %v", err)
	}
	if !strings.Contains(err.Error(), "field ID") {
		t.Fatalf("expected field name in error, got %v", err)
	}
}

func TestBindArgsRejectsUnsupportedFieldTypes(t *testing.T) {
	type input struct {
		Tags []string
	}

	ctx := &MsgContext{Args: []string{"tag"}}
	var got input

	err := ctx.BindArgs(&got)
	if err == nil {
		t.Fatal("expected BindArgs to fail")
	}
	if !errors.Is(err, ErrBindArgsUnsupportedFieldType) {
		t.Fatalf("expected ErrBindArgsUnsupportedFieldType, got %v", err)
	}
}

func TestErrorDefaultRemainsUserVisibleForMessageFlow(t *testing.T) {
	var requests int
	var gotBody map[string]any

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests++
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
		Api:           api,
		Msg:           &tgapi.Message{Chat: &tgapi.Chat{ID: 42, Type: tgapi.ChatTypePrivate}},
		Logger:        slog.CreateLogger(),
		errorTemplate: "Error: %s",
	}

	ctx.error(errors.New("boom"))

	if requests != 1 {
		t.Fatalf("expected one user-facing error reply, got %d requests", requests)
	}
	if got := gotBody["text"]; got != "Error: boom" {
		t.Fatalf("unexpected error reply text: %v", got)
	}
}

func TestErrorInternalSkipsUserReplyForMessageFlow(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatal("unexpected HTTP request for internal-only error")
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
		Api:           api,
		Msg:           &tgapi.Message{Chat: &tgapi.Chat{ID: 42, Type: tgapi.ChatTypePrivate}},
		Logger:        slog.CreateLogger(),
		errorTemplate: "Error: %s",
	}

	ctx.error(AsInternalError(errors.New("boom")))
}

func TestErrorInternalSkipsCallbackAnswer(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatal("unexpected callback answer request for internal-only error")
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
		Api:             api,
		Logger:          slog.CreateLogger(),
		errorTemplate:   "%s",
		CallbackQueryId: "cb-1",
	}

	ctx.error(AsInternalError(errors.New("boom")))
}

func TestErrorUserVisibleAnswersCallback(t *testing.T) {
	var requests int
	var gotBody map[string]any

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests++
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
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":true}`)),
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
		Api:             api,
		Logger:          slog.CreateLogger(),
		errorTemplate:   "Oops: %s",
		CallbackQueryId: "cb-1",
	}

	ctx.error(AsUserError(errors.New("boom")))

	if requests != 1 {
		t.Fatalf("expected one callback error answer, got %d requests", requests)
	}
	if got := gotBody["text"]; got != "Oops: boom" {
		t.Fatalf("unexpected callback error text: %v", got)
	}
}

func TestAnswerRejectsEmptyMessage(t *testing.T) {
	ctx := &MsgContext{
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 42, Type: tgapi.ChatTypePrivate}},
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
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 42, Type: tgapi.ChatTypePrivate}},
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
		Msg:    &tgapi.Message{Chat: &tgapi.Chat{ID: 42, Type: tgapi.ChatTypePrivate}},
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
