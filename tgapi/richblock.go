package tgapi

import (
	"encoding/json"
	"fmt"
)

// RichBlock is a block in a structured rich message.
type RichBlock interface {
	isRichBlock()
}

// ---------------------------------------------------------------------------
// Helper types
// ---------------------------------------------------------------------------

// RichBlockCaption is the caption of a media block or container.
type RichBlockCaption struct {
	Text   RichText
	Credit RichText
}

func (c RichBlockCaption) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Text   RichText `json:"text"`
		Credit RichText `json:"credit,omitempty"`
	}{c.Text, c.Credit})
}

func (c *RichBlockCaption) UnmarshalJSON(data []byte) error {
	var raw struct {
		Text   json.RawMessage `json:"text"`
		Credit json.RawMessage `json:"credit"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	text, err := parseOptRichText(raw.Text)
	if err != nil {
		return err
	}
	credit, err := parseOptRichText(raw.Credit)
	if err != nil {
		return err
	}
	*c = RichBlockCaption{text, credit}
	return nil
}

// RichBlockListItem is a single list item. Label is the ready-to-display
// visible marker ("1.", "c.", "vii.", "•"): the server renders it itself
// when parsing html/markdown.
type RichBlockListItem struct {
	Label       string
	Blocks      []RichBlock
	HasCheckbox bool
	IsChecked   bool
	Value       int    // for ordered lists: numeric value of the marker
	Type        string // for ordered lists: "a", "A", "i", "I" or "1"
}

func (i RichBlockListItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Label       string      `json:"label"`
		Blocks      []RichBlock `json:"blocks"`
		HasCheckbox bool        `json:"has_checkbox,omitempty"`
		IsChecked   bool        `json:"is_checked,omitempty"`
		Value       int         `json:"value,omitempty"`
		Type        string      `json:"type,omitempty"`
	}{i.Label, i.Blocks, i.HasCheckbox, i.IsChecked, i.Value, i.Type})
}

func (i *RichBlockListItem) UnmarshalJSON(data []byte) error {
	var raw struct {
		Label       string          `json:"label"`
		Blocks      json.RawMessage `json:"blocks"`
		HasCheckbox bool            `json:"has_checkbox"`
		IsChecked   bool            `json:"is_checked"`
		Value       int             `json:"value"`
		Type        string          `json:"type"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	blocks, err := unmarshalRichBlocks(raw.Blocks)
	if err != nil {
		return err
	}
	*i = RichBlockListItem{raw.Label, blocks, raw.HasCheckbox, raw.IsChecked, raw.Value, raw.Type}
	return nil
}

// RichBlockTableCell is a table cell. An empty Text means an invisible cell.
type RichBlockTableCell struct {
	Text     RichText
	IsHeader bool
	Colspan  int
	Rowspan  int
	Align    string // "left", "center" or "right"
	VAlign   string // "top", "middle" or "bottom"
}

func (c RichBlockTableCell) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Text     RichText `json:"text,omitempty"`
		IsHeader bool     `json:"is_header,omitempty"`
		Colspan  int      `json:"colspan,omitempty"`
		Rowspan  int      `json:"rowspan,omitempty"`
		Align    string   `json:"align,omitempty"`
		VAlign   string   `json:"valign,omitempty"`
	}{c.Text, c.IsHeader, c.Colspan, c.Rowspan, c.Align, c.VAlign})
}

func (c *RichBlockTableCell) UnmarshalJSON(data []byte) error {
	var raw struct {
		Text     json.RawMessage `json:"text"`
		IsHeader bool            `json:"is_header"`
		Colspan  int             `json:"colspan"`
		Rowspan  int             `json:"rowspan"`
		Align    string          `json:"align"`
		VAlign   string          `json:"valign"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	text, err := parseOptRichText(raw.Text)
	if err != nil {
		return err
	}
	*c = RichBlockTableCell{text, raw.IsHeader, raw.Colspan, raw.Rowspan, raw.Align, raw.VAlign}
	return nil
}

// ---------------------------------------------------------------------------
// RichBlockWrap: pure text blocks — paragraph, footer, thinking.
// ---------------------------------------------------------------------------

// RichBlockWrap covers all blocks that have only a text field.
type RichBlockWrap struct {
	Tag  string
	Text RichText
}

func (RichBlockWrap) isRichBlock() {}

func (b RichBlockWrap) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string   `json:"type"`
		Text RichText `json:"text"`
	}{b.Tag, b.Text})
}

var richBlockWrapTags = map[string]bool{
	"paragraph": true, "footer": true, "thinking": true,
}

// ---------------------------------------------------------------------------
// Section heading: text + size
// ---------------------------------------------------------------------------

// RichBlockSectionHeading is a section heading block.
type RichBlockSectionHeading struct {
	Text RichText
	Size int // 1-6, 1 is the largest
}

func (RichBlockSectionHeading) isRichBlock() {}
func (b RichBlockSectionHeading) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string   `json:"type"`
		Text RichText `json:"text"`
		Size int      `json:"size"`
	}{"heading", b.Text, b.Size})
}

// ---------------------------------------------------------------------------
// Block with text + language
// ---------------------------------------------------------------------------

// RichBlockPreformatted is a preformatted code block.
type RichBlockPreformatted struct {
	Text     RichText
	Language string
}

func (RichBlockPreformatted) isRichBlock() {}
func (b RichBlockPreformatted) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string   `json:"type"`
		Text     RichText `json:"text"`
		Language string   `json:"language,omitempty"`
	}{"pre", b.Text, b.Language})
}

// ---------------------------------------------------------------------------
// Quotations
// ---------------------------------------------------------------------------

// RichBlockQuotation is a block quotation with block-level content
// (officially RichBlockBlockQuotation).
type RichBlockQuotation struct {
	Blocks []RichBlock
	Credit RichText
}

func (RichBlockQuotation) isRichBlock() {}
func (b RichBlockQuotation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string      `json:"type"`
		Blocks []RichBlock `json:"blocks"`
		Credit RichText    `json:"credit,omitempty"`
	}{"blockquote", b.Blocks, b.Credit})
}

// RichBlockPullQuotation is a pull quotation with inline content.
type RichBlockPullQuotation struct {
	Text   RichText
	Credit RichText
}

func (RichBlockPullQuotation) isRichBlock() {}
func (b RichBlockPullQuotation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string   `json:"type"`
		Text   RichText `json:"text"`
		Credit RichText `json:"credit,omitempty"`
	}{"pullquote", b.Text, b.Credit})
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// RichBlockList is a list block.
type RichBlockList struct {
	Items []RichBlockListItem
}

func (RichBlockList) isRichBlock() {}
func (b RichBlockList) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string              `json:"type"`
		Items []RichBlockListItem `json:"items"`
	}{"list", b.Items})
}

// ---------------------------------------------------------------------------
// Containers with blocks []RichBlock + caption
// ---------------------------------------------------------------------------

// RichBlockCollage is a collage of media blocks.
type RichBlockCollage struct {
	Blocks  []RichBlock
	Caption *RichBlockCaption
}

func (RichBlockCollage) isRichBlock() {}
func (b RichBlockCollage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string            `json:"type"`
		Blocks  []RichBlock       `json:"blocks"`
		Caption *RichBlockCaption `json:"caption,omitempty"`
	}{"collage", b.Blocks, b.Caption})
}

// RichBlockSlideshow is a slideshow of media blocks.
type RichBlockSlideshow struct {
	Blocks  []RichBlock
	Caption *RichBlockCaption
}

func (RichBlockSlideshow) isRichBlock() {}
func (b RichBlockSlideshow) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string            `json:"type"`
		Blocks  []RichBlock       `json:"blocks"`
		Caption *RichBlockCaption `json:"caption,omitempty"`
	}{"slideshow", b.Blocks, b.Caption})
}

// ---------------------------------------------------------------------------
// Details — expandable block
// ---------------------------------------------------------------------------

// RichBlockDetails is an expandable block with an inline summary.
type RichBlockDetails struct {
	Summary RichText
	Blocks  []RichBlock
	IsOpen  bool
}

func (RichBlockDetails) isRichBlock() {}
func (b RichBlockDetails) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string      `json:"type"`
		Summary RichText    `json:"summary"`
		Blocks  []RichBlock `json:"blocks"`
		IsOpen  bool        `json:"is_open,omitempty"`
	}{"details", b.Summary, b.Blocks, b.IsOpen})
}

// ---------------------------------------------------------------------------
// Table
// ---------------------------------------------------------------------------

// RichBlockTable is a table block.
type RichBlockTable struct {
	Cells      [][]RichBlockTableCell
	IsBordered bool
	IsStriped  bool
	Caption    RichText
}

func (RichBlockTable) isRichBlock() {}
func (b RichBlockTable) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string                 `json:"type"`
		Cells      [][]RichBlockTableCell `json:"cells"`
		IsBordered bool                   `json:"is_bordered,omitempty"`
		IsStriped  bool                   `json:"is_striped,omitempty"`
		Caption    RichText               `json:"caption,omitempty"`
	}{"table", b.Cells, b.IsBordered, b.IsStriped, b.Caption})
}

// ---------------------------------------------------------------------------
// Map
// ---------------------------------------------------------------------------

// RichBlockMap is a location map block.
type RichBlockMap struct {
	Location Location
	Zoom     int // 13-20
	Width    int
	Height   int
	Caption  *RichBlockCaption
}

func (RichBlockMap) isRichBlock() {}
func (b RichBlockMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string            `json:"type"`
		Location Location          `json:"location"`
		Zoom     int               `json:"zoom"`
		Width    int               `json:"width"`
		Height   int               `json:"height"`
		Caption  *RichBlockCaption `json:"caption,omitempty"`
	}{"map", b.Location, b.Zoom, b.Width, b.Height, b.Caption})
}

// ---------------------------------------------------------------------------
// Media blocks
// ---------------------------------------------------------------------------

// RichBlockPhoto is a photo block.
type RichBlockPhoto struct {
	Photo      []PhotoSize
	HasSpoiler bool
	Caption    *RichBlockCaption
}

func (RichBlockPhoto) isRichBlock() {}
func (b RichBlockPhoto) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string            `json:"type"`
		Photo      []PhotoSize       `json:"photo"`
		HasSpoiler bool              `json:"has_spoiler,omitempty"`
		Caption    *RichBlockCaption `json:"caption,omitempty"`
	}{"photo", b.Photo, b.HasSpoiler, b.Caption})
}

// RichBlockVideo is a video block.
type RichBlockVideo struct {
	Video      Video
	HasSpoiler bool
	Caption    *RichBlockCaption
}

func (RichBlockVideo) isRichBlock() {}
func (b RichBlockVideo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string            `json:"type"`
		Video      Video             `json:"video"`
		HasSpoiler bool              `json:"has_spoiler,omitempty"`
		Caption    *RichBlockCaption `json:"caption,omitempty"`
	}{"video", b.Video, b.HasSpoiler, b.Caption})
}

// RichBlockAudio is an audio block.
type RichBlockAudio struct {
	Audio   Audio
	Caption *RichBlockCaption
}

func (RichBlockAudio) isRichBlock() {}
func (b RichBlockAudio) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string            `json:"type"`
		Audio   Audio             `json:"audio"`
		Caption *RichBlockCaption `json:"caption,omitempty"`
	}{"audio", b.Audio, b.Caption})
}

// RichBlockAnimation is an animation block.
type RichBlockAnimation struct {
	Animation  Animation
	HasSpoiler bool
	Caption    *RichBlockCaption
}

func (RichBlockAnimation) isRichBlock() {}
func (b RichBlockAnimation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string            `json:"type"`
		Animation  Animation         `json:"animation"`
		HasSpoiler bool              `json:"has_spoiler,omitempty"`
		Caption    *RichBlockCaption `json:"caption,omitempty"`
	}{"animation", b.Animation, b.HasSpoiler, b.Caption})
}

// RichBlockVoiceNote is a voice note block.
type RichBlockVoiceNote struct {
	VoiceNote Voice
	Caption   *RichBlockCaption
}

func (RichBlockVoiceNote) isRichBlock() {}
func (b RichBlockVoiceNote) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string            `json:"type"`
		VoiceNote Voice             `json:"voice_note"`
		Caption   *RichBlockCaption `json:"caption,omitempty"`
	}{"voice_note", b.VoiceNote, b.Caption})
}

// ---------------------------------------------------------------------------
// Leaves without nested content
// ---------------------------------------------------------------------------

// RichBlockDivider is a horizontal divider block.
type RichBlockDivider struct{}

func (RichBlockDivider) isRichBlock() {}
func (b RichBlockDivider) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
	}{"divider"})
}

// RichBlockMathematicalExpression is a block-level mathematical expression.
type RichBlockMathematicalExpression struct {
	Expression string
}

func (RichBlockMathematicalExpression) isRichBlock() {}
func (b RichBlockMathematicalExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string `json:"type"`
		Expression string `json:"expression"`
	}{"mathematical_expression", b.Expression})
}

// RichBlockAnchor is a named anchor block that anchor links can point to.
type RichBlockAnchor struct {
	Name string
}

func (RichBlockAnchor) isRichBlock() {}
func (b RichBlockAnchor) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{"anchor", b.Name})
}

// ---------------------------------------------------------------------------
// JSON -> RichBlock parsing
// ---------------------------------------------------------------------------

// UnmarshalRichBlock parses a single RichBlock from JSON, dispatching on the
// type tag. Unknown types that carry a text field are preserved as
// RichBlockWrap for forward compatibility.
func UnmarshalRichBlock(data []byte) (RichBlock, error) {
	var head struct {
		Type string          `json:"type"`
		Text json.RawMessage `json:"text"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return nil, fmt.Errorf("richblock: %w", err)
	}

	if richBlockWrapTags[head.Type] {
		text, err := parseOptRichText(head.Text)
		if err != nil {
			return nil, fmt.Errorf("richblock %q: text: %w", head.Type, err)
		}
		return RichBlockWrap{Tag: head.Type, Text: text}, nil
	}

	switch head.Type {
	case "heading":
		var v struct {
			Size int `json:"size"`
		}
		_ = json.Unmarshal(data, &v)
		text, _ := parseOptRichText(head.Text)
		return RichBlockSectionHeading{text, v.Size}, nil

	case "pre":
		var v struct {
			Language string `json:"language"`
		}
		_ = json.Unmarshal(data, &v)
		text, _ := parseOptRichText(head.Text)
		return RichBlockPreformatted{text, v.Language}, nil

	case "blockquote":
		var raw struct {
			Blocks json.RawMessage `json:"blocks"`
			Credit json.RawMessage `json:"credit"`
		}
		_ = json.Unmarshal(data, &raw)
		blocks, err := unmarshalRichBlocks(raw.Blocks)
		if err != nil {
			return nil, err
		}
		credit, _ := parseOptRichText(raw.Credit)
		return RichBlockQuotation{blocks, credit}, nil

	case "pullquote":
		var raw struct {
			Credit json.RawMessage `json:"credit"`
		}
		_ = json.Unmarshal(data, &raw)
		text, _ := parseOptRichText(head.Text)
		credit, _ := parseOptRichText(raw.Credit)
		return RichBlockPullQuotation{text, credit}, nil

	case "list":
		var v struct {
			Items []RichBlockListItem `json:"items"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return RichBlockList{v.Items}, nil

	case "collage":
		var raw struct {
			Blocks  json.RawMessage   `json:"blocks"`
			Caption *RichBlockCaption `json:"caption"`
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, err
		}
		blocks, err := unmarshalRichBlocks(raw.Blocks)
		if err != nil {
			return nil, err
		}
		return RichBlockCollage{blocks, raw.Caption}, nil

	case "slideshow":
		var raw struct {
			Blocks  json.RawMessage   `json:"blocks"`
			Caption *RichBlockCaption `json:"caption"`
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, err
		}
		blocks, err := unmarshalRichBlocks(raw.Blocks)
		if err != nil {
			return nil, err
		}
		return RichBlockSlideshow{blocks, raw.Caption}, nil

	case "details":
		var raw struct {
			Summary json.RawMessage `json:"summary"`
			Blocks  json.RawMessage `json:"blocks"`
			IsOpen  bool            `json:"is_open"`
		}
		_ = json.Unmarshal(data, &raw)
		summary, _ := parseOptRichText(raw.Summary)
		blocks, err := unmarshalRichBlocks(raw.Blocks)
		if err != nil {
			return nil, err
		}
		return RichBlockDetails{summary, blocks, raw.IsOpen}, nil

	case "table":
		var raw struct {
			Cells      [][]RichBlockTableCell `json:"cells"`
			IsBordered bool                   `json:"is_bordered"`
			IsStriped  bool                   `json:"is_striped"`
			Caption    json.RawMessage        `json:"caption"`
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, err
		}
		caption, _ := parseOptRichText(raw.Caption)
		return RichBlockTable{raw.Cells, raw.IsBordered, raw.IsStriped, caption}, nil

	case "map":
		var v struct {
			Location Location          `json:"location"`
			Zoom     int               `json:"zoom"`
			Width    int               `json:"width"`
			Height   int               `json:"height"`
			Caption  *RichBlockCaption `json:"caption"`
		}
		_ = json.Unmarshal(data, &v)
		return RichBlockMap{v.Location, v.Zoom, v.Width, v.Height, v.Caption}, nil

	case "photo":
		var v struct {
			Photo      []PhotoSize       `json:"photo"`
			HasSpoiler bool              `json:"has_spoiler"`
			Caption    *RichBlockCaption `json:"caption"`
		}
		_ = json.Unmarshal(data, &v)
		return RichBlockPhoto{v.Photo, v.HasSpoiler, v.Caption}, nil

	case "video":
		var v struct {
			Video      Video             `json:"video"`
			HasSpoiler bool              `json:"has_spoiler"`
			Caption    *RichBlockCaption `json:"caption"`
		}
		_ = json.Unmarshal(data, &v)
		return RichBlockVideo{v.Video, v.HasSpoiler, v.Caption}, nil

	case "audio":
		var v struct {
			Audio   Audio             `json:"audio"`
			Caption *RichBlockCaption `json:"caption"`
		}
		_ = json.Unmarshal(data, &v)
		return RichBlockAudio{v.Audio, v.Caption}, nil

	case "animation":
		var v struct {
			Animation  Animation         `json:"animation"`
			HasSpoiler bool              `json:"has_spoiler"`
			Caption    *RichBlockCaption `json:"caption"`
		}
		_ = json.Unmarshal(data, &v)
		return RichBlockAnimation{v.Animation, v.HasSpoiler, v.Caption}, nil

	case "voice_note":
		var v struct {
			VoiceNote Voice             `json:"voice_note"`
			Caption   *RichBlockCaption `json:"caption"`
		}
		_ = json.Unmarshal(data, &v)
		return RichBlockVoiceNote{v.VoiceNote, v.Caption}, nil

	case "divider":
		return RichBlockDivider{}, nil

	case "mathematical_expression":
		var v struct {
			Expression string `json:"expression"`
		}
		_ = json.Unmarshal(data, &v)
		return RichBlockMathematicalExpression{v.Expression}, nil

	case "anchor":
		var v struct {
			Name string `json:"name"`
		}
		_ = json.Unmarshal(data, &v)
		return RichBlockAnchor{v.Name}, nil

	default:
		// forward-compat: unknown type with text -> RichBlockWrap, without text -> error.
		if text, err := parseOptRichText(head.Text); err == nil && text != nil {
			return RichBlockWrap{Tag: head.Type, Text: text}, nil
		}
		return nil, fmt.Errorf("richblock: unknown type %q", head.Type)
	}
}

// UnmarshalRichMessage parses a root RichMessage from JSON.
func UnmarshalRichMessage(data []byte) (RichMessage, error) {
	var raw struct {
		Blocks json.RawMessage `json:"blocks"`
		IsRTL  bool            `json:"is_rtl"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return RichMessage{}, fmt.Errorf("richmessage: %w", err)
	}
	blocks, err := unmarshalRichBlocks(raw.Blocks)
	if err != nil {
		return RichMessage{}, err
	}
	return RichMessage{blocks, raw.IsRTL}, nil
}

// UnmarshalJSON parses the blocks through UnmarshalRichBlock: the Blocks
// field is interface-typed, so the standard unmarshaler cannot handle it.
func (m *RichMessage) UnmarshalJSON(data []byte) error {
	parsed, err := UnmarshalRichMessage(data)
	if err != nil {
		return err
	}
	*m = parsed
	return nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// parseOptRichText parses an optional RichText field: absent and null yield nil.
func parseOptRichText(raw json.RawMessage) (RichText, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	return UnmarshalRichText(raw)
}

func unmarshalRichBlocks(raw json.RawMessage) ([]RichBlock, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var raws []json.RawMessage
	if err := json.Unmarshal(raw, &raws); err != nil {
		return nil, err
	}
	blocks := make([]RichBlock, len(raws))
	for i, r := range raws {
		b, err := UnmarshalRichBlock(r)
		if err != nil {
			return nil, err
		}
		blocks[i] = b
	}
	return blocks, nil
}
