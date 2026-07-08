package tgapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEditMessageTextMarshalsInputRichMessage(t *testing.T) {
	params := EditMessageText{
		ChatID:    1,
		MessageID: 2,
		RichMessage: &InputRichMessage{
			HTML:                "<p>hi</p>",
			SkipEntityDetection: true,
		},
	}
	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	got := string(data)
	for _, want := range []string{`"rich_message":{"html":`, `"skip_entity_detection":true`} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %s in editMessageText JSON: %s", want, got)
		}
	}
	if strings.Contains(got, `"blocks"`) {
		t.Fatalf("rich_message must be an InputRichMessage, not a block tree: %s", got)
	}
	if strings.Contains(got, `"text"`) {
		t.Fatalf("empty text must be omitted when editing rich content: %s", got)
	}
}

func TestSendRichMessageDraftMarshal(t *testing.T) {
	params := SendRichMessageDraft{
		ChatID:      1,
		DraftID:     7,
		RichMessage: InputRichMessage{Markdown: "*hi*"},
	}
	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	got := string(data)
	for _, want := range []string{`"chat_id":1`, `"draft_id":7`, `"rich_message":{"markdown":"*hi*"}`} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %s in sendRichMessageDraft JSON: %s", want, got)
		}
	}
}

func TestInputRichMessageContentMarshal(t *testing.T) {
	content := InputRichMessageContent{
		RichMessage: InputRichMessage{HTML: "<p>hi</p>"},
	}
	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if got := string(data); !strings.Contains(got, `"rich_message":{"html":`) {
		t.Fatalf("unexpected InputRichMessageContent JSON: %s", got)
	}
}

func TestInputPollOptionMediaLinkMarshal(t *testing.T) {
	media := InputPollOptionMedia{Type: "link", URL: "https://example.com"}
	data, err := json.Marshal(media)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	got := string(data)
	if got != `{"type":"link","url":"https://example.com"}` {
		t.Fatalf("unexpected link media JSON: %s", got)
	}
}

func TestPollMediaUnmarshalLink(t *testing.T) {
	var media PollMedia
	if err := json.Unmarshal([]byte(`{"link":{"url":"https://example.com"}}`), &media); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if media.Link == nil || media.Link.URL != "https://example.com" {
		t.Fatalf("unexpected poll media link: %+v", media.Link)
	}
}

func TestChatJoinRequestUnmarshalQueryID(t *testing.T) {
	payload := `{"chat":{"id":1},"from":{"id":2,"first_name":"A"},"user_chat_id":2,"date":3,"query_id":"q42"}`
	var req ChatJoinRequest
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if req.QueryID == nil || *req.QueryID != "q42" {
		t.Fatalf("unexpected query_id: %+v", req.QueryID)
	}
}

func TestAnswerChatJoinRequestQueryResultValues(t *testing.T) {
	if JoinRequestApprove != "approve" || JoinRequestDecline != "decline" || JoinRequestQueue != "queue" {
		t.Fatalf("unexpected join request query result values: %q %q %q",
			JoinRequestApprove, JoinRequestDecline, JoinRequestQueue)
	}
}

func TestUserUnmarshalSupportsJoinRequestQueries(t *testing.T) {
	var user User
	if err := json.Unmarshal([]byte(`{"id":1,"first_name":"A","supports_join_request_queries":true}`), &user); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if user.SupportsJoinRequestQueries == nil || !*user.SupportsJoinRequestQueries {
		t.Fatalf("unexpected supports_join_request_queries: %+v", user.SupportsJoinRequestQueries)
	}
}
