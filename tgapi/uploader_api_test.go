package tgapi

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

func TestUploaderEncodesJSONFieldsAndLeavesAcceptEncodingToHTTPTransport(t *testing.T) {
	var (
		gotPath           string
		gotAcceptEncoding string
		gotFields         map[string]string
		gotFileName       string
		gotFileData       []byte
		roundTripErr      error
	)

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotPath = req.URL.Path
			gotAcceptEncoding = req.Header.Get("Accept-Encoding")

			gotFields, gotFileName, gotFileData, roundTripErr = readMultipartRequest(req)
			if roundTripErr != nil {
				roundTripErr = fmt.Errorf("readMultipartRequest: %w", roundTripErr)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":{"message_id":5,"date":1}}`)),
			}, nil
		}),
	}

	api := NewAPI(
		NewAPIOpts("token").
			SetAPIURL("https://example.test").
			SetHTTPClient(client),
	)
	defer func() {
		if err := api.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()

	uploader := NewUploader(api)
	defer func() {
		if err := uploader.Close(); err != nil {
			t.Fatalf("Close returned error: %v", err)
		}
	}()

	msg, err := uploader.SendPhoto(
		UploadPhoto{
			ChatID: 42,
			CaptionEntities: []MessageEntity{{
				Type:   MessageEntityBold,
				Offset: 0,
				Length: 4,
			}},
			ReplyMarkup: &ReplyMarkup{
				InlineKeyboard: [][]InlineKeyboardButton{{
					{Text: "A", CallbackData: "b"},
				}},
			},
		},
		NewUploaderFile("photo.jpg", []byte("img")),
	)
	if err != nil {
		t.Fatalf("SendPhoto returned error: %v", err)
	}
	if msg.MessageID != 5 {
		t.Fatalf("unexpected message id: %d", msg.MessageID)
	}
	if roundTripErr != nil {
		t.Fatalf("multipart parse failed: %v", roundTripErr)
	}
	if gotPath != "/bottoken/sendPhoto" {
		t.Fatalf("unexpected request path: %s", gotPath)
	}
	if gotAcceptEncoding != "" {
		t.Fatalf("expected empty Accept-Encoding header, got %q", gotAcceptEncoding)
	}
	if got := gotFields["chat_id"]; got != "42" {
		t.Fatalf("chat_id mismatch: %q", got)
	}
	if got := gotFields["caption_entities"]; got != `[{"type":"bold","offset":0,"length":4}]` {
		t.Fatalf("caption_entities mismatch: %q", got)
	}
	if got := gotFields["reply_markup"]; got != `{"inline_keyboard":[[{"text":"A","callback_data":"b"}]]}` {
		t.Fatalf("reply_markup mismatch: %q", got)
	}
	if gotFileName != "photo.jpg" {
		t.Fatalf("unexpected file name: %q", gotFileName)
	}
	if string(gotFileData) != "img" {
		t.Fatalf("unexpected file content: %q", string(gotFileData))
	}
}

func TestNewUploaderFileDetectsFileTypeCaseInsensitively(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     UploaderFileType
	}{
		{name: "uppercase photo", filename: "PHOTO.JPG", want: UploaderPhotoType},
		{name: "uppercase voice", filename: "voice.OGG", want: UploaderVoiceType},
		{name: "unknown defaults to document", filename: "archive.BIN", want: UploaderDocumentType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := NewUploaderFile(tt.filename, []byte("x"))
			if file.field != tt.want {
				t.Fatalf("unexpected uploader field: got %q want %q", file.field, tt.want)
			}
		})
	}
}

func readMultipartRequest(req *http.Request) (map[string]string, string, []byte, error) {
	_, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return nil, "", nil, err
	}
	reader := multipart.NewReader(req.Body, params["boundary"])

	fields := make(map[string]string)
	var fileName string
	var fileData []byte
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return fields, fileName, fileData, nil
		}
		if err != nil {
			return nil, "", nil, err
		}

		data, err := io.ReadAll(part)
		if err != nil {
			return nil, "", nil, err
		}

		if part.FileName() != "" {
			fileName = part.FileName()
			fileData = data
			continue
		}
		fields[part.FormName()] = string(data)
	}
}
