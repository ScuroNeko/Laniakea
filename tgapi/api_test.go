package tgapi

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestAPILeavesAcceptEncodingToHTTPTransport(t *testing.T) {
	var gotPath string
	var gotAcceptEncoding string

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotPath = req.URL.Path
			gotAcceptEncoding = req.Header.Get("Accept-Encoding")
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"Test"}}`)),
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

	user, err := api.GetMe()
	if err != nil {
		t.Fatalf("GetMe returned error: %v", err)
	}
	if user.FirstName != "Test" {
		t.Fatalf("unexpected first name: %q", user.FirstName)
	}
	if gotPath != "/bottoken/getMe" {
		t.Fatalf("unexpected request path: %s", gotPath)
	}
	if gotAcceptEncoding != "" {
		t.Fatalf("expected empty Accept-Encoding header, got %q", gotAcceptEncoding)
	}
}
