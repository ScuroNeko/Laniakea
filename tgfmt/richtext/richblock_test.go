package richtext

import (
	"encoding/json"
	"testing"
)

func roundtripBlock(t *testing.T, in RichBlock) {
	t.Helper()
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out, err := UnmarshalBlock(b)
	if err != nil {
		t.Fatalf("unmarshal %s: %v", b, err)
	}
	b2, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("remarshal: %v", err)
	}
	if string(b) != string(b2) {
		t.Fatalf("not stable:\n  %s\n  %s", b, b2)
	}
}

func TestBlockRoundtrip(t *testing.T) {
	cases := []RichBlock{
		// wrap-блоки
		Paragraph(String("Hello, world")),
		SectionHeading(Bold(String("Chapter 1"))),
		Footer(String("© 2024")),
		Thinking(String("Let me reason step by step.")),

		// preformatted
		BlockPreformatted{Text: String("fmt.Println(\"hi\")"), Language: "go"},
		BlockPreformatted{Text: String("no language")},

		// цитаты
		BlockBlockQuotation{Text: String("To be or not to be"), Caption: String("Shakespeare")},
		BlockPullQuotation{Text: String("Pull me"), Caption: nil},

		// список
		BlockList{
			Items: []RichBlockListItem{
				{Blocks: []RichBlock{Paragraph(String("item 1"))}},
				{Blocks: []RichBlock{Paragraph(String("item 2"))}},
			},
			Ordered: true,
		},
		BlockList{
			Items: []RichBlockListItem{
				{Blocks: []RichBlock{Paragraph(String("bullet"))}},
			},
			Ordered: false,
		},

		// коллаж и слайдшоу
		BlockCollage{
			Items:   []RichBlock{BlockPhoto{FileID: "abc123"}},
			Caption: String("A photo"),
		},
		BlockSlideshow{
			Items:   []RichBlock{BlockVideo{FileID: "vid1", Autoplay: true, Loop: false}},
			Caption: nil,
		},

		// details
		BlockDetails{
			Title:  String("Spoiler"),
			Blocks: []RichBlock{Paragraph(String("Hidden content"))},
			Open:   false,
		},
		BlockDetails{
			Title:  Bold(String("Open details")),
			Blocks: []RichBlock{BlockDivider{}, Paragraph(String("content"))},
			Open:   true,
		},

		// таблица
		BlockTable{
			Title: String("Results"),
			Rows: [][]RichBlockTableCell{
				{
					{Content: []RichBlock{Paragraph(String("Cell A1"))}},
					{Content: []RichBlock{Paragraph(String("Cell A2"))}, ColumnSpan: 2},
				},
				{
					{Content: []RichBlock{Paragraph(String("Cell B1"))}, RowSpan: 2},
					{Content: []RichBlock{Paragraph(String("Cell B2"))}},
				},
			},
			Bordered: true,
			Striped:  false,
		},

		// карта
		BlockMap{
			Latitude: 55.7558, Longitude: 37.6173,
			Zoom: 12, Width: 800, Height: 400,
			Caption: String("Moscow"),
		},

		// медиа
		BlockPhoto{FileID: "photo_file_id", Caption: String("A cat"), URL: "https://example.com/cat.jpg"},
		BlockPhoto{FileID: "bare_photo"},
		BlockVideo{FileID: "video_file_id", Caption: String("Demo"), Autoplay: true, Loop: true},
		BlockAudio{FileID: "audio_file_id", Caption: String("Podcast ep. 1")},
		BlockAnimation{FileID: "anim_file_id"},
		BlockVoiceNote{FileID: "voice_file_id"},

		// листья
		BlockDivider{},
		BlockMathematicalExpression{Expression: "E = mc^2"},
		BlockAnchor{Name: "section-2"},
	}
	for _, c := range cases {
		roundtripBlock(t, c)
	}
}

func TestRichMessageRoundtrip(t *testing.T) {
	msg := RichMessage{
		Blocks: []RichBlock{
			SectionHeading(String("Title")),
			Paragraph(Array{String("Some "), Bold(String("bold")), String(" text")}),
			BlockDivider{},
			BlockList{
				Items: []RichBlockListItem{
					{Blocks: []RichBlock{Paragraph(String("First"))}},
					{Blocks: []RichBlock{Paragraph(String("Second"))}},
				},
				Ordered: true,
			},
			BlockPhoto{FileID: "img1", Caption: String("Fig. 1")},
		},
	}

	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out RichMessage
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	b2, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("remarshal: %v", err)
	}
	if string(b) != string(b2) {
		t.Fatalf("not stable:\n  %s\n  %s", b, b2)
	}
}

func TestBlockDividerHasNoContent(t *testing.T) {
	b, _ := json.Marshal(BlockDivider{})
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if len(m) != 1 {
		t.Fatalf("divider must only have type field: %s", b)
	}
	if m["type"] != "divider" {
		t.Fatalf("unexpected type: %s", b)
	}
}

func TestBlockUnknownTypeWithTextIsForwardCompat(t *testing.T) {
	raw := []byte(`{"type":"future_tag","text":"hello"}`)
	b, err := UnmarshalBlock(raw)
	if err != nil {
		t.Fatalf("forward-compat failed: %v", err)
	}
	w, ok := b.(BlockWrap)
	if !ok || w.Tag != "future_tag" {
		t.Fatalf("expected BlockWrap{future_tag}, got %T", b)
	}
}

func TestBlockUnknownTypeWithoutTextIsError(t *testing.T) {
	raw := []byte(`{"type":"mystery_leaf","value":42}`)
	_, err := UnmarshalBlock(raw)
	if err == nil {
		t.Fatal("expected error for unknown type without text")
	}
}
