package tgapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUpdateDeletedBusinessMessagesUnmarshalSetsAlias(t *testing.T) {
	var update Update
	err := json.Unmarshal([]byte(`{
		"update_id": 1,
		"deleted_business_messages": {
			"business_connection_id": "conn",
			"chat": {"id": 42, "type": "private"},
			"message_ids": [3, 5]
		}
	}`), &update)
	if err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if update.DeletedBusinessMessages == nil {
		t.Fatal("expected DeletedBusinessMessages to be populated")
	}
	if update.DeletedBusinessMessage == nil {
		t.Fatal("expected deprecated DeletedBusinessMessage alias to be populated")
	}
	if update.DeletedBusinessMessages != update.DeletedBusinessMessage {
		t.Fatal("expected deleted business message fields to share the same payload")
	}
	if got := update.DeletedBusinessMessages.MessageIDs; len(got) != 2 || got[0] != 3 || got[1] != 5 {
		t.Fatalf("unexpected message ids: %v", got)
	}
}

func TestUpdateMarshalUsesCanonicalDeletedBusinessMessagesField(t *testing.T) {
	update := Update{
		UpdateID: 1,
		DeletedBusinessMessage: &BusinessMessagesDeleted{
			BusinessConnectionID: "conn",
			Chat:                 Chat{ID: 42, Type: string(ChatTypePrivate)},
			MessageIDs:           []int{7},
		},
	}

	data, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	got := string(data)
	if !strings.Contains(got, `"deleted_business_messages"`) {
		t.Fatalf("expected canonical deleted_business_messages field, got %s", got)
	}
	if strings.Contains(got, `"deleted_business_message"`) {
		t.Fatalf("unexpected singular deleted_business_message field, got %s", got)
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
}
