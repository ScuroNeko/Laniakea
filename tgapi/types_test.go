package tgapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUpdateUnmarshalSetsType(t *testing.T) {
	tests := []struct {
		name string
		body string
		want UpdateType
	}{
		{
			name: "deleted business messages",
			body: `{
				"update_id": 1,
				"deleted_business_messages": {
					"business_connection_id": "conn",
					"chat": {"id": 42, "type": "private"},
					"message_ids": [3, 5]
				}
			}`,
			want: UpdateTypeDeletedBusinessMessages,
		},
		{
			name: "callback query",
			body: `{
				"update_id": 2,
				"callback_query": {
					"id": "cb",
					"from": {"id": 1, "is_bot": false, "first_name": "Test"},
					"chat_instance": "instance",
					"data": "payload"
				}
			}`,
			want: UpdateTypeCallbackQuery,
		},
		{
			name: "chat boost",
			body: `{
				"update_id": 3,
				"chat_boost": {
					"chat": {"id": -1001, "type": "supergroup", "title": "Boosted"},
					"boost": {
						"boost_id": "boost-1",
						"add_date": 1735689600,
						"expiration_date": 1738291600,
						"source": {
							"source": "premium",
							"user": {"id": 1, "is_bot": false, "first_name": "Test"}
						}
					}
				}
			}`,
			want: UpdateTypeChatBoost,
		},
		{
			name: "unknown",
			body: `{"update_id":4}`,
			want: UpdateTypeUnknown,
		},
		{
			name: "managed bot",
			body: `{
				"update_id": 5,
				"managed_bot": {
					"user": {"id": 11, "is_bot": false, "first_name": "Manager"},
					"bot": {"id": 12, "is_bot": true, "first_name": "Worker"}
				}
			}`,
			want: UpdateTypeManagedBot,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var update Update
			if err := json.Unmarshal([]byte(tt.body), &update); err != nil {
				t.Fatalf("Unmarshal returned error: %v", err)
			}
			if update.Type != tt.want {
				t.Fatalf("unexpected update type: got %q want %q", update.Type, tt.want)
			}
			if tt.want == UpdateTypeChatBoost && update.ChatBoost.Boost.BoostID != "boost-1" {
				t.Fatalf("unexpected boost id: got %q want %q", update.ChatBoost.Boost.BoostID, "boost-1")
			}
			if tt.want == UpdateTypeManagedBot && update.ManagedBot.Bot.ID != 12 {
				t.Fatalf("unexpected managed bot id: got %d want %d", update.ManagedBot.Bot.ID, 12)
			}
		})
	}
}

func TestPollUnmarshalSupportsBotAPI96Fields(t *testing.T) {
	var poll Poll

	body := `{
		"id": "poll-1",
		"question": "Pick winners",
		"question_entities": [],
		"options": [],
		"total_voter_count": 2,
		"is_closed": false,
		"is_anonymous": false,
		"type": "quiz",
		"allows_multiple_answers": true,
		"allows_revoting": true,
		"correct_option_ids": [1, 3]
	}`

	if err := json.Unmarshal([]byte(body), &poll); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if !poll.AllowsRevoting {
		t.Fatal("expected allows_revoting to be decoded")
	}
	if len(poll.CorrectOptionIDs) != 2 || poll.CorrectOptionIDs[0] != 1 || poll.CorrectOptionIDs[1] != 3 {
		t.Fatalf("unexpected correct option ids: %#v", poll.CorrectOptionIDs)
	}
}

func TestUpdateMarshalOmitsSyntheticTypeField(t *testing.T) {
	update := Update{
		UpdateID: 1,
		Type:     UpdateTypeCallbackQuery,
		CallbackQuery: &CallbackQuery{
			ID:           "cb",
			From:         User{ID: 1, FirstName: "Test"},
			ChatInstance: "instance",
			Data:         "payload",
		},
	}

	data, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	got := string(data)
	if strings.Contains(got, `"type"`) {
		t.Fatalf("unexpected synthetic type field, got %s", got)
	}
}

func TestUpdateShippingQueryIsNilWhenAbsent(t *testing.T) {
	var update Update
	if err := json.Unmarshal([]byte(`{"update_id":1}`), &update); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if update.ShippingQuery != nil {
		t.Fatalf("expected ShippingQuery to be nil, got %+v", update.ShippingQuery)
	}
	if update.Type != UpdateTypeUnknown {
		t.Fatalf("expected UpdateTypeUnknown, got %q", update.Type)
	}
}

func TestMaybeInaccessibleMessageUnmarshalAccessibleMessage(t *testing.T) {
	var wrapper MaybeInaccessibleMessage

	body := `{
		"message_id": 10,
		"date": 1700000000,
		"chat": {"id": 42, "type": "private"},
		"text": "hello"
	}`

	if err := json.Unmarshal([]byte(body), &wrapper); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if !wrapper.IsAccessible() {
		t.Fatal("expected accessible message payload")
	}
	if wrapper.IsInaccessible() {
		t.Fatal("expected inaccessible payload to be empty")
	}
	if wrapper.Message() == nil || wrapper.Message().Text != "hello" {
		t.Fatalf("unexpected accessible payload: %#v", wrapper.Message())
	}
	if wrapper.MessageID() != 10 {
		t.Fatalf("unexpected message id: got %d want %d", wrapper.MessageID(), 10)
	}
	if wrapper.Chat() == nil || wrapper.Chat().ID != 42 {
		t.Fatalf("unexpected chat payload: %#v", wrapper.Chat())
	}
}

func TestMaybeInaccessibleMessageUnmarshalInaccessibleMessage(t *testing.T) {
	var wrapper MaybeInaccessibleMessage

	body := `{
		"message_id": 7,
		"date": 0,
		"chat": {"id": -1001, "type": "supergroup"}
	}`

	if err := json.Unmarshal([]byte(body), &wrapper); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if wrapper.IsAccessible() {
		t.Fatal("expected accessible payload to be empty")
	}
	if !wrapper.IsInaccessible() {
		t.Fatal("expected inaccessible message payload")
	}
	if wrapper.InaccessibleMessage() == nil || wrapper.InaccessibleMessage().MessageID != 7 {
		t.Fatalf("unexpected inaccessible payload: %#v", wrapper.InaccessibleMessage())
	}
	if wrapper.MessageID() != 7 {
		t.Fatalf("unexpected message id: got %d want %d", wrapper.MessageID(), 7)
	}
	if wrapper.Chat() == nil || wrapper.Chat().ID != -1001 {
		t.Fatalf("unexpected chat payload: %#v", wrapper.Chat())
	}
}
