package utils_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"testing"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
	"git.scuroneko.dev/scuroneko/laniakea/utils"
)

type multipartEncodeParams struct {
	ChatID          int64                  `json:"chat_id"`
	MessageThreadID *int                   `json:"message_thread_id,omitempty"`
	ReplyMarkup     *tgapi.ReplyMarkup     `json:"reply_markup,omitempty"`
	CaptionEntities []tgapi.MessageEntity  `json:"caption_entities,omitempty"`
	ReplyParameters *tgapi.ReplyParameters `json:"reply_parameters,omitempty"`
}

func TestEncodeMultipartJSONFields(t *testing.T) {
	params := multipartEncodeParams{
		ChatID:          42,
		MessageThreadID: new(7),
		ReplyMarkup: &tgapi.ReplyMarkup{
			InlineKeyboard: [][]tgapi.InlineKeyboardButton{{
				{Text: "A", CallbackData: "b"},
			}},
		},
		CaptionEntities: []tgapi.MessageEntity{{
			Type:   tgapi.MessageEntityBold,
			Offset: 0,
			Length: 4,
		}},
	}

	body := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(body)
	if err := utils.Encode(writer, params); err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close returned error: %v", err)
	}

	got := readMultipartFields(t, body.Bytes(), writer.Boundary())
	if got["chat_id"] != "42" {
		t.Fatalf("chat_id mismatch: %q", got["chat_id"])
	}
	if got["message_thread_id"] != "7" {
		t.Fatalf("message_thread_id mismatch: %q", got["message_thread_id"])
	}
	if got["reply_markup"] != `{"inline_keyboard":[[{"text":"A","callback_data":"b"}]]}` {
		t.Fatalf("reply_markup mismatch: %q", got["reply_markup"])
	}
	if got["caption_entities"] != `[{"type":"bold","offset":0,"length":4}]` {
		t.Fatalf("caption_entities mismatch: %q", got["caption_entities"])
	}
	if _, ok := got["reply_parameters"]; ok {
		t.Fatalf("reply_parameters should be omitted when nil")
	}
}

func readMultipartFields(t *testing.T, body []byte, boundary string) map[string]string {
	t.Helper()

	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	fields := make(map[string]string)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return fields
		}
		if err != nil {
			t.Fatalf("NextPart returned error: %v", err)
		}

		data, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("ReadAll returned error: %v", err)
		}
		fields[part.FormName()] = string(data)
	}
}
