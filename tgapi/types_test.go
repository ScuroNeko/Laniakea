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
		})
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
