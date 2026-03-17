package tgapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestReplyKeyboardMarkupMarshalsKeyboardButtons(t *testing.T) {
	markup := ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{{
			{
				Text:        "Create poll",
				RequestPoll: &KeyboardButtonPollType{Type: PollTypeQuiz},
			},
		}},
	}

	data, err := json.Marshal(markup)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	got := string(data)
	if !strings.Contains(got, `"keyboard":[[{"text":"Create poll","request_poll":{"type":"quiz"}}]]`) {
		t.Fatalf("unexpected reply keyboard JSON: %s", got)
	}
}

func TestChatActionUploadVideoNoteValue(t *testing.T) {
	if ChatActionUploadVideoNote != "upload_video_note" {
		t.Fatalf("unexpected chat action value: %q", ChatActionUploadVideoNote)
	}
	if ChatActionUploadVideoNone != ChatActionUploadVideoNote {
		t.Fatalf("expected deprecated alias to match upload_video_note, got %q", ChatActionUploadVideoNone)
	}
}
