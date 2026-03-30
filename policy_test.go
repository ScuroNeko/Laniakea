package laniakea

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/slog"
)

func TestRequirePolicyStopsExecutionOnDeniedPolicy(t *testing.T) {
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

	mw := RequirePolicy[NoData]("deny", func(ctx *MsgContext, data NoData) error {
		return AsUserError(errors.New("blocked"))
	})

	if mw.Execute(ctx, NoData{}) {
		t.Fatal("expected denied policy middleware to stop execution")
	}
	if requests != 1 {
		t.Fatalf("expected one user-facing error reply, got %d requests", requests)
	}
	if got := gotBody["text"]; got != "Error: blocked" {
		t.Fatalf("unexpected policy error reply text: %v", got)
	}
}

func TestRequirePrivateChatAllowsPrivateChat(t *testing.T) {
	ctx := &MsgContext{
		Msg: &tgapi.Message{
			Chat: &tgapi.Chat{ID: 42, Type: tgapi.ChatTypePrivate},
		},
		Logger: slog.CreateLogger(),
	}

	if err := RequirePrivateChat[NoData]()(ctx, NoData{}); err != nil {
		t.Fatalf("RequirePrivateChat returned error: %v", err)
	}
}

func TestRequirePrivateChatDeniesNonPrivateChat(t *testing.T) {
	ctx := &MsgContext{
		Msg: &tgapi.Message{
			Chat: &tgapi.Chat{ID: -100, Type: tgapi.ChatTypeSupergroup},
		},
		Logger: slog.CreateLogger(),
	}

	err := RequirePrivateChat[NoData]()(ctx, NoData{})
	if err == nil {
		t.Fatal("expected RequirePrivateChat to deny non-private chats")
	}
	if !IsUserError(err) {
		t.Fatalf("expected user-visible deny error, got %v", err)
	}
}

func TestRequireChatAdminUsesNormalizedIDs(t *testing.T) {
	var sawGetChatMember bool
	var gotBody map[string]any

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.Path, "getChatMember") {
				t.Fatalf("unexpected API method: %s", req.URL.Path)
			}
			sawGetChatMember = true
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
				Body: io.NopCloser(strings.NewReader(
					`{"ok":true,"result":{"status":"administrator","user":{"id":55,"is_bot":false,"first_name":"tester"}}}`,
				)),
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
		ChatID: -2001,
		FromID: 55,
		Logger: slog.CreateLogger(),
	}

	if err := RequireChatAdmin[NoData]()(ctx, NoData{}); err != nil {
		t.Fatalf("RequireChatAdmin returned error: %v", err)
	}
	if !sawGetChatMember {
		t.Fatal("expected GetChatMember to be called")
	}
	if got := gotBody["chat_id"]; got != float64(-2001) {
		t.Fatalf("unexpected chat_id in request: %v", got)
	}
	if got := gotBody["user_id"]; got != float64(55) {
		t.Fatalf("unexpected user_id in request: %v", got)
	}
}

func TestAllPoliciesReturnsFirstError(t *testing.T) {
	want := AsUserError(errors.New("blocked"))
	policy := AllPolicies[NoData](
		func(ctx *MsgContext, data NoData) error { return nil },
		func(ctx *MsgContext, data NoData) error { return want },
		func(ctx *MsgContext, data NoData) error {
			t.Fatal("unexpected evaluation after first failure")
			return nil
		},
	)

	err := policy(&MsgContext{Logger: slog.CreateLogger()}, NoData{})
	if !errors.Is(err, want) {
		t.Fatalf("expected first policy error, got %v", err)
	}
}

func TestAnyPolicyAllowsLaterSuccessAfterInternalError(t *testing.T) {
	policy := AnyPolicy[NoData](
		func(ctx *MsgContext, data NoData) error { return AsInternalError(errors.New("temporary")) },
		func(ctx *MsgContext, data NoData) error { return nil },
	)

	if err := policy(&MsgContext{Logger: slog.CreateLogger()}, NoData{}); err != nil {
		t.Fatalf("expected later success to allow access, got %v", err)
	}
}

func TestAnyPolicyReturnsInternalErrorWhenNonePass(t *testing.T) {
	internal := AsInternalError(errors.New("temporary"))
	policy := AnyPolicy[NoData](
		func(ctx *MsgContext, data NoData) error { return AsUserError(errors.New("denied")) },
		func(ctx *MsgContext, data NoData) error { return internal },
	)

	err := policy(&MsgContext{Logger: slog.CreateLogger()}, NoData{})
	if !errors.Is(err, internal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestAnyPolicyReturnsFirstDenyWhenNoPolicyPasses(t *testing.T) {
	first := AsUserError(errors.New("first deny"))
	policy := AnyPolicy[NoData](
		func(ctx *MsgContext, data NoData) error { return first },
		func(ctx *MsgContext, data NoData) error { return AsUserError(errors.New("second deny")) },
	)

	err := policy(&MsgContext{Logger: slog.CreateLogger()}, NoData{})
	if !errors.Is(err, first) {
		t.Fatalf("expected first deny error, got %v", err)
	}
}

func TestNotPolicyInvertsUserDenyButPreservesInternalErrors(t *testing.T) {
	inverted := NotPolicy[NoData](func(ctx *MsgContext, data NoData) error {
		return AsUserError(errors.New("denied"))
	})
	if err := inverted(&MsgContext{Logger: slog.CreateLogger()}, NoData{}); err != nil {
		t.Fatalf("expected inverted deny to succeed, got %v", err)
	}

	internal := AsInternalError(errors.New("temporary"))
	preserve := NotPolicy[NoData](func(ctx *MsgContext, data NoData) error {
		return internal
	})
	err := preserve(&MsgContext{Logger: slog.CreateLogger()}, NoData{})
	if !errors.Is(err, internal) {
		t.Fatalf("expected internal error to be preserved, got %v", err)
	}
}
