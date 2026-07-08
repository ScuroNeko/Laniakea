package tgapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func roundtripRichBlock(t *testing.T, in RichBlock) {
	t.Helper()
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out, err := UnmarshalRichBlock(b)
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

func par(s string) RichBlockWrap { return RichBlockWrap{"paragraph", RichTextPlain(s)} }

func TestRichBlockRoundtrip(t *testing.T) {
	cases := []RichBlock{
		// wrap blocks
		par("Hello, world"),
		RichBlockWrap{"footer", RichTextPlain("© 2024")},
		RichBlockWrap{"thinking", RichTextPlain("Let me reason step by step.")},

		// heading
		RichBlockSectionHeading{RichTextWrap{"bold", RichTextPlain("Chapter 1")}, 1},
		RichBlockSectionHeading{RichTextPlain("smallest"), 6},

		// preformatted
		RichBlockPreformatted{RichTextPlain(`fmt.Println("hi")`), "go"},
		RichBlockPreformatted{Text: RichTextPlain("no language")},

		// quotations
		RichBlockQuotation{[]RichBlock{par("To be or not to be")}, RichTextPlain("Shakespeare")},
		RichBlockQuotation{Blocks: []RichBlock{par("anonymous"), par("second block")}},
		RichBlockPullQuotation{Text: RichTextPlain("Pull me")},
		RichBlockPullQuotation{RichTextPlain("Wisdom"), RichTextWrap{"italic", RichTextPlain("someone")}},

		// list: label is the ready-made marker, numbering lives on the items
		RichBlockList{
			Items: []RichBlockListItem{
				{Label: "c.", Blocks: []RichBlock{par("item 3")}, Value: 3, Type: "a"},
				{Label: "vii.", Blocks: []RichBlock{par("item 7")}, Value: 7, Type: "i"},
			},
		},
		RichBlockList{
			Items: []RichBlockListItem{
				{Label: "•", Blocks: []RichBlock{par("todo")}, HasCheckbox: true},
				{Label: "•", Blocks: []RichBlock{par("done")}, HasCheckbox: true, IsChecked: true},
			},
		},

		// collage and slideshow
		RichBlockCollage{
			Blocks:  []RichBlock{RichBlockPhoto{Photo: []PhotoSize{{FileID: "abc123", Width: 100, Height: 100}}}},
			Caption: &RichBlockCaption{Text: RichTextPlain("A photo")},
		},
		RichBlockSlideshow{
			Blocks: []RichBlock{
				RichBlockVideo{Video: Video{FileID: "vid1", Width: 640, Height: 480, Duration: 10}},
			},
		},

		// details
		RichBlockDetails{
			Summary: RichTextPlain("Spoiler"),
			Blocks:  []RichBlock{par("Hidden content")},
		},
		RichBlockDetails{
			Summary: RichTextWrap{"bold", RichTextPlain("Open details")},
			Blocks:  []RichBlock{RichBlockDivider{}, par("content")},
			IsOpen:  true,
		},

		// table: text cells, headers, spans, alignment, invisible cell
		RichBlockTable{
			Cells: [][]RichBlockTableCell{
				{
					{Text: RichTextPlain("Name"), IsHeader: true, Align: "center"},
					{Text: RichTextPlain("Score"), IsHeader: true, VAlign: "middle"},
				},
				{
					{Text: RichTextPlain("Alice"), Colspan: 2},
				},
				{
					{}, // invisible cell
					{Text: RichTextPlain("42"), Rowspan: 2},
				},
			},
			IsBordered: true,
			Caption:    RichTextPlain("Results"),
		},

		// map
		RichBlockMap{
			Location: Location{Latitude: 55.7558, Longitude: 37.6173},
			Zoom:     13, Width: 800, Height: 400,
			Caption: &RichBlockCaption{Text: RichTextPlain("Moscow"), Credit: RichTextPlain("OpenStreetMap")},
		},

		// media
		RichBlockPhoto{
			Photo:      []PhotoSize{{FileID: "p1", Width: 1280, Height: 720}},
			HasSpoiler: true,
			Caption:    &RichBlockCaption{Text: RichTextPlain("A cat"), Credit: RichTextWrap{"italic", RichTextPlain("photographer")}},
		},
		RichBlockVideo{Video: Video{FileID: "v1", Width: 1920, Height: 1080, Duration: 30}, HasSpoiler: true},
		RichBlockAudio{
			Audio:   Audio{FileID: "a1", Duration: 60},
			Caption: &RichBlockCaption{Text: RichTextPlain("Podcast ep. 1")},
		},
		RichBlockAnimation{Animation: Animation{FileID: "g1", Width: 320, Height: 240, Duration: 2}},
		RichBlockVoiceNote{VoiceNote: Voice{FileID: "vn1", Duration: 5}},

		// leaves
		RichBlockDivider{},
		RichBlockMathematicalExpression{Expression: "E = mc^2"},
		RichBlockAnchor{Name: "section-2"},
	}
	for _, c := range cases {
		roundtripRichBlock(t, c)
	}
}

func TestRichMessageRoundtrip(t *testing.T) {
	for _, msg := range []RichMessage{
		{
			Blocks: []RichBlock{
				RichBlockSectionHeading{RichTextPlain("Title"), 1},
				RichBlockWrap{"paragraph", RichTextArray{RichTextPlain("Some "), RichTextWrap{"bold", RichTextPlain("bold")}, RichTextPlain(" text")}},
				RichBlockDivider{},
				RichBlockList{
					Items: []RichBlockListItem{
						{Label: "1.", Blocks: []RichBlock{par("First")}, Value: 1, Type: "1"},
						{Label: "2.", Blocks: []RichBlock{par("Second")}, Value: 2, Type: "1"},
					},
				},
				RichBlockPhoto{
					Photo:   []PhotoSize{{FileID: "img1", Width: 10, Height: 10}},
					Caption: &RichBlockCaption{Text: RichTextPlain("Fig. 1")},
				},
			},
		},
		{
			Blocks: []RichBlock{par("שלום")},
			IsRTL:  true,
		},
	} {
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
}

func TestRichBlockTags(t *testing.T) {
	// tags per spec: heading, pre, blockquote, pullquote
	cases := map[string]RichBlock{
		"heading":    RichBlockSectionHeading{RichTextPlain("h"), 2},
		"pre":        RichBlockPreformatted{Text: RichTextPlain("x")},
		"blockquote": RichBlockQuotation{Blocks: []RichBlock{par("q")}},
		"pullquote":  RichBlockPullQuotation{Text: RichTextPlain("p")},
	}
	for want, block := range cases {
		b, _ := json.Marshal(block)
		var m map[string]any
		_ = json.Unmarshal(b, &m)
		if m["type"] != want {
			t.Fatalf("expected type %q, got %s", want, b)
		}
	}
}

func TestRichBlockOptionalFieldsOmitted(t *testing.T) {
	// nil credit/caption and false flags must not appear in the JSON
	for _, c := range []struct {
		block RichBlock
		bad   []string
	}{
		{RichBlockQuotation{Blocks: []RichBlock{par("q")}}, []string{"credit"}},
		{RichBlockPullQuotation{Text: RichTextPlain("p")}, []string{"credit"}},
		{RichBlockPhoto{Photo: []PhotoSize{{FileID: "p"}}}, []string{"caption", "has_spoiler"}},
		{RichBlockTable{Cells: [][]RichBlockTableCell{}}, []string{"caption", "is_bordered", "is_striped"}},
		{RichBlockDetails{Summary: RichTextPlain("s")}, []string{"is_open"}},
	} {
		b, _ := json.Marshal(c.block)
		for _, key := range c.bad {
			if strings.Contains(string(b), `"`+key+`"`) {
				t.Fatalf("%T: %q must be omitted: %s", c.block, key, b)
			}
		}
	}
	// same for RichMessage.is_rtl
	b, _ := json.Marshal(RichMessage{Blocks: []RichBlock{par("x")}})
	if strings.Contains(string(b), "is_rtl") {
		t.Fatalf("is_rtl must be omitted: %s", b)
	}
}

func TestRichBlockDividerHasNoContent(t *testing.T) {
	b, _ := json.Marshal(RichBlockDivider{})
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if len(m) != 1 {
		t.Fatalf("divider must only have type field: %s", b)
	}
	if m["type"] != "divider" {
		t.Fatalf("unexpected type: %s", b)
	}
}

func TestRichBlockUnknownTypeWithTextIsForwardCompat(t *testing.T) {
	raw := []byte(`{"type":"future_tag","text":"hello"}`)
	b, err := UnmarshalRichBlock(raw)
	if err != nil {
		t.Fatalf("forward-compat failed: %v", err)
	}
	w, ok := b.(RichBlockWrap)
	if !ok || w.Tag != "future_tag" {
		t.Fatalf("expected RichBlockWrap{future_tag}, got %T", b)
	}
}

func TestRichBlockUnknownTypeWithoutTextIsError(t *testing.T) {
	raw := []byte(`{"type":"mystery_leaf","value":42}`)
	_, err := UnmarshalRichBlock(raw)
	if err == nil {
		t.Fatal("expected error for unknown type without text")
	}
}
