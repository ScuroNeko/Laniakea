package tgapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetFileByLinkUsesConfiguredAPIURL(t *testing.T) {
	var gotPath string

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotPath = req.URL.Path
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("payload")),
			}, nil
		}),
	}

	api := NewAPI(
		NewAPIOpts("token").
			SetAPIUrl("https://example.test").
			SetHTTPClient(client),
	)
	defer func() {
		if err := api.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()

	data, err := api.GetFileByLink("files/report.txt")
	if err != nil {
		t.Fatalf("GetFileByLink returned error: %v", err)
	}
	if string(data) != "payload" {
		t.Fatalf("unexpected payload: %q", string(data))
	}
	if gotPath != "/file/bottoken/files/report.txt" {
		t.Fatalf("unexpected request path: %s", gotPath)
	}
}

func TestOpenFileByLinkStreamsResponseBody(t *testing.T) {
	api := NewAPI(
		NewAPIOpts("token").
			SetAPIUrl("https://example.test").
			SetHTTPClient(&http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("streamed payload")),
					}, nil
				}),
			}),
	)
	defer func() {
		if err := api.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()

	body, err := api.OpenFileByLink("files/report.txt")
	if err != nil {
		t.Fatalf("OpenFileByLink returned error: %v", err)
	}
	defer func() {
		if err := body.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()

	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	if string(data) != "streamed payload" {
		t.Fatalf("unexpected payload: %q", string(data))
	}
}

func TestGetFileByLinkReturnsHTTPStatusError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("missing\n")),
			}, nil
		}),
	}

	api := NewAPI(
		NewAPIOpts("token").
			SetAPIUrl("https://example.test").
			SetHTTPClient(client),
	)
	defer func() {
		if err := api.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()

	_, err := api.GetFileByLink("files/report.txt")
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestGetUpdatesOmitsAllowedUpdatesWhenEmpty(t *testing.T) {
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
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":[]}`)),
			}, nil
		}),
	}

	api := NewAPI(
		NewAPIOpts("token").
			SetAPIUrl("https://example.test").
			SetHTTPClient(client),
	)
	defer func() {
		if err := api.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()

	updates, err := api.GetUpdates(UpdateParams{})
	if err != nil {
		t.Fatalf("GetUpdates returned error: %v", err)
	}
	if len(updates) != 0 {
		t.Fatalf("expected no updates, got %d", len(updates))
	}
	if _, exists := gotBody["allowed_updates"]; exists {
		t.Fatalf("expected allowed_updates to be omitted, got %v", gotBody["allowed_updates"])
	}
}
