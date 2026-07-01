package richtext

import (
	"encoding/json"
	"fmt"
)

// RichBlock — блок в структурированном rich-сообщении.
type RichBlock interface {
	isRichBlock()
}

// RichMessage — корневой тип структурированного сообщения (Bot API 10.1).
type RichMessage struct {
	Blocks []RichBlock
}

func (m RichMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []RichBlock `json:"blocks"`
	}{m.Blocks})
}

func (m *RichMessage) UnmarshalJSON(data []byte) error {
	msg, err := UnmarshalMessage(data)
	if err != nil {
		return err
	}
	*m = msg
	return nil
}

// ---------------------------------------------------------------------------
// Вспомогательные типы
// ---------------------------------------------------------------------------

// RichBlockListItem — один элемент списка (ordered/unordered).
type RichBlockListItem struct {
	Blocks []RichBlock
}

func (i RichBlockListItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []RichBlock `json:"blocks"`
	}{i.Blocks})
}

func (i *RichBlockListItem) UnmarshalJSON(data []byte) error {
	item, err := unmarshalListItem(data)
	if err != nil {
		return err
	}
	*i = item
	return nil
}

// RichBlockTableCell — ячейка таблицы.
type RichBlockTableCell struct {
	Content    []RichBlock
	ColumnSpan int
	RowSpan    int
}

func (c RichBlockTableCell) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Content    []RichBlock `json:"content"`
		ColumnSpan int         `json:"column_span,omitempty"`
		RowSpan    int         `json:"row_span,omitempty"`
	}{c.Content, c.ColumnSpan, c.RowSpan})
}

func (c *RichBlockTableCell) UnmarshalJSON(data []byte) error {
	cell, err := unmarshalTableCell(data)
	if err != nil {
		return err
	}
	*c = cell
	return nil
}

// ---------------------------------------------------------------------------
// BlockWrap: чистые текстовые блоки — paragraph, section_heading, footer, thinking.
// ---------------------------------------------------------------------------

// BlockWrap покрывает все блоки, у которых есть только поле text.
type BlockWrap struct {
	Tag  string
	Text RichText
}

func (BlockWrap) isRichBlock() {}

func (b BlockWrap) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string   `json:"type"`
		Text RichText `json:"text"`
	}{b.Tag, b.Text})
}

var blockWrapTags = map[string]bool{
	"paragraph": true, "section_heading": true,
	"footer": true, "thinking": true,
}

func Paragraph(t RichText) BlockWrap      { return BlockWrap{"paragraph", t} }
func SectionHeading(t RichText) BlockWrap { return BlockWrap{"section_heading", t} }
func Footer(t RichText) BlockWrap         { return BlockWrap{"footer", t} }
func Thinking(t RichText) BlockWrap       { return BlockWrap{"thinking", t} }

// ---------------------------------------------------------------------------
// Блок с text + language
// ---------------------------------------------------------------------------

type BlockPreformatted struct {
	Text     RichText
	Language string
}

func (BlockPreformatted) isRichBlock() {}
func (b BlockPreformatted) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string   `json:"type"`
		Text     RichText `json:"text"`
		Language string   `json:"language,omitempty"`
	}{"preformatted", b.Text, b.Language})
}

// ---------------------------------------------------------------------------
// Блоки с text + caption
// ---------------------------------------------------------------------------

type BlockBlockQuotation struct {
	Text    RichText
	Caption RichText
}

func (BlockBlockQuotation) isRichBlock() {}
func (b BlockBlockQuotation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string   `json:"type"`
		Text    RichText `json:"text"`
		Caption RichText `json:"caption,omitempty"`
	}{"block_quotation", b.Text, b.Caption})
}

type BlockPullQuotation struct {
	Text    RichText
	Caption RichText
}

func (BlockPullQuotation) isRichBlock() {}
func (b BlockPullQuotation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string   `json:"type"`
		Text    RichText `json:"text"`
		Caption RichText `json:"caption,omitempty"`
	}{"pull_quotation", b.Text, b.Caption})
}

// ---------------------------------------------------------------------------
// Список
// ---------------------------------------------------------------------------

type BlockList struct {
	Items   []RichBlockListItem
	Ordered bool
}

func (BlockList) isRichBlock() {}
func (b BlockList) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string              `json:"type"`
		Items   []RichBlockListItem `json:"items"`
		Ordered bool                `json:"ordered"`
	}{"list", b.Items, b.Ordered})
}

// ---------------------------------------------------------------------------
// Контейнеры с items []RichBlock + caption
// ---------------------------------------------------------------------------

type BlockCollage struct {
	Items   []RichBlock
	Caption RichText
}

func (BlockCollage) isRichBlock() {}
func (b BlockCollage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string      `json:"type"`
		Items   []RichBlock `json:"items"`
		Caption RichText    `json:"caption,omitempty"`
	}{"collage", b.Items, b.Caption})
}

type BlockSlideshow struct {
	Items   []RichBlock
	Caption RichText
}

func (BlockSlideshow) isRichBlock() {}
func (b BlockSlideshow) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string      `json:"type"`
		Items   []RichBlock `json:"items"`
		Caption RichText    `json:"caption,omitempty"`
	}{"slideshow", b.Items, b.Caption})
}

// ---------------------------------------------------------------------------
// Details — раскрывающийся блок
// ---------------------------------------------------------------------------

type BlockDetails struct {
	Title  RichText
	Blocks []RichBlock
	Open   bool
}

func (BlockDetails) isRichBlock() {}
func (b BlockDetails) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string      `json:"type"`
		Title  RichText    `json:"title"`
		Blocks []RichBlock `json:"blocks"`
		Open   bool        `json:"open"`
	}{"details", b.Title, b.Blocks, b.Open})
}

// ---------------------------------------------------------------------------
// Таблица
// ---------------------------------------------------------------------------

type BlockTable struct {
	Title    RichText
	Rows     [][]RichBlockTableCell
	Bordered bool
	Striped  bool
}

func (BlockTable) isRichBlock() {}
func (b BlockTable) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string                 `json:"type"`
		Title    RichText               `json:"title,omitempty"`
		Rows     [][]RichBlockTableCell `json:"rows"`
		Bordered bool                   `json:"bordered"`
		Striped  bool                   `json:"striped"`
	}{"table", b.Title, b.Rows, b.Bordered, b.Striped})
}

// ---------------------------------------------------------------------------
// Карта
// ---------------------------------------------------------------------------

type BlockMap struct {
	Latitude  float64
	Longitude float64
	Zoom      int
	Width     int
	Height    int
	Caption   RichText
}

func (BlockMap) isRichBlock() {}
func (b BlockMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string   `json:"type"`
		Latitude  float64  `json:"latitude"`
		Longitude float64  `json:"longitude"`
		Zoom      int      `json:"zoom"`
		Width     int      `json:"width"`
		Height    int      `json:"height"`
		Caption   RichText `json:"caption,omitempty"`
	}{"map", b.Latitude, b.Longitude, b.Zoom, b.Width, b.Height, b.Caption})
}

// ---------------------------------------------------------------------------
// Медиа-блоки (file_id + caption)
// ---------------------------------------------------------------------------

type BlockPhoto struct {
	FileID  string
	Caption RichText
	URL     string
}

func (BlockPhoto) isRichBlock() {}
func (b BlockPhoto) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string   `json:"type"`
		FileID  string   `json:"file_id"`
		Caption RichText `json:"caption,omitempty"`
		URL     string   `json:"url,omitempty"`
	}{"photo", b.FileID, b.Caption, b.URL})
}

type BlockVideo struct {
	FileID   string
	Caption  RichText
	Autoplay bool
	Loop     bool
}

func (BlockVideo) isRichBlock() {}
func (b BlockVideo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string   `json:"type"`
		FileID   string   `json:"file_id"`
		Caption  RichText `json:"caption,omitempty"`
		Autoplay bool     `json:"autoplay"`
		Loop     bool     `json:"loop"`
	}{"video", b.FileID, b.Caption, b.Autoplay, b.Loop})
}

type BlockAudio struct {
	FileID  string
	Caption RichText
}

func (BlockAudio) isRichBlock() {}
func (b BlockAudio) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string   `json:"type"`
		FileID  string   `json:"file_id"`
		Caption RichText `json:"caption,omitempty"`
	}{"audio", b.FileID, b.Caption})
}

type BlockAnimation struct {
	FileID  string
	Caption RichText
}

func (BlockAnimation) isRichBlock() {}
func (b BlockAnimation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string   `json:"type"`
		FileID  string   `json:"file_id"`
		Caption RichText `json:"caption,omitempty"`
	}{"animation", b.FileID, b.Caption})
}

type BlockVoiceNote struct {
	FileID  string
	Caption RichText
}

func (BlockVoiceNote) isRichBlock() {}
func (b BlockVoiceNote) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string   `json:"type"`
		FileID  string   `json:"file_id"`
		Caption RichText `json:"caption,omitempty"`
	}{"voice_note", b.FileID, b.Caption})
}

// ---------------------------------------------------------------------------
// Листья без вложенного контента
// ---------------------------------------------------------------------------

type BlockDivider struct{}

func (BlockDivider) isRichBlock() {}
func (b BlockDivider) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
	}{"divider"})
}

type BlockMathematicalExpression struct {
	Expression string
}

func (BlockMathematicalExpression) isRichBlock() {}
func (b BlockMathematicalExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string `json:"type"`
		Expression string `json:"expression"`
	}{"mathematical_expression", b.Expression})
}

type BlockAnchor struct {
	Name string
}

func (BlockAnchor) isRichBlock() {}
func (b BlockAnchor) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{"anchor", b.Name})
}

// ---------------------------------------------------------------------------
// Разбор JSON -> RichBlock
// ---------------------------------------------------------------------------

func UnmarshalBlock(data []byte) (RichBlock, error) {
	var head struct {
		Type    string          `json:"type"`
		Text    json.RawMessage `json:"text"`
		Caption json.RawMessage `json:"caption"`
		Title   json.RawMessage `json:"title"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return nil, fmt.Errorf("richblock: %w", err)
	}

	parseText := func(raw json.RawMessage) (RichText, error) {
		if len(raw) == 0 || string(raw) == "null" {
			return nil, nil
		}
		return Unmarshal(raw)
	}

	if blockWrapTags[head.Type] {
		text, err := parseText(head.Text)
		if err != nil {
			return nil, fmt.Errorf("richblock %q: text: %w", head.Type, err)
		}
		return BlockWrap{Tag: head.Type, Text: text}, nil
	}

	switch head.Type {
	case "preformatted":
		var v struct {
			Language string `json:"language"`
		}
		_ = json.Unmarshal(data, &v)
		text, _ := parseText(head.Text)
		return BlockPreformatted{text, v.Language}, nil

	case "block_quotation":
		text, _ := parseText(head.Text)
		caption, _ := parseText(head.Caption)
		return BlockBlockQuotation{text, caption}, nil

	case "pull_quotation":
		text, _ := parseText(head.Text)
		caption, _ := parseText(head.Caption)
		return BlockPullQuotation{text, caption}, nil

	case "list":
		var v struct {
			Items   []RichBlockListItem `json:"items"`
			Ordered bool                `json:"ordered"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return BlockList{v.Items, v.Ordered}, nil

	case "collage":
		var raw struct {
			Items json.RawMessage `json:"items"`
		}
		_ = json.Unmarshal(data, &raw)
		items, _ := unmarshalBlocks(raw.Items)
		caption, _ := parseText(head.Caption)
		return BlockCollage{items, caption}, nil

	case "slideshow":
		var raw struct {
			Items json.RawMessage `json:"items"`
		}
		_ = json.Unmarshal(data, &raw)
		items, _ := unmarshalBlocks(raw.Items)
		caption, _ := parseText(head.Caption)
		return BlockSlideshow{items, caption}, nil

	case "details":
		var raw struct {
			Blocks json.RawMessage `json:"blocks"`
			Open   bool            `json:"open"`
		}
		_ = json.Unmarshal(data, &raw)
		title, _ := parseText(head.Title)
		blocks, _ := unmarshalBlocks(raw.Blocks)
		return BlockDetails{title, blocks, raw.Open}, nil

	case "table":
		var raw struct {
			Title    json.RawMessage        `json:"title"`
			Rows     [][]RichBlockTableCell `json:"rows"`
			Bordered bool                   `json:"bordered"`
			Striped  bool                   `json:"striped"`
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, err
		}
		title, _ := parseText(raw.Title)
		return BlockTable{title, raw.Rows, raw.Bordered, raw.Striped}, nil

	case "map":
		var v struct {
			Latitude  float64         `json:"latitude"`
			Longitude float64         `json:"longitude"`
			Zoom      int             `json:"zoom"`
			Width     int             `json:"width"`
			Height    int             `json:"height"`
			Caption   json.RawMessage `json:"caption"`
		}
		_ = json.Unmarshal(data, &v)
		caption, _ := parseText(v.Caption)
		return BlockMap{v.Latitude, v.Longitude, v.Zoom, v.Width, v.Height, caption}, nil

	case "photo":
		var v struct {
			FileID  string          `json:"file_id"`
			Caption json.RawMessage `json:"caption"`
			URL     string          `json:"url"`
		}
		_ = json.Unmarshal(data, &v)
		caption, _ := parseText(v.Caption)
		return BlockPhoto{v.FileID, caption, v.URL}, nil

	case "video":
		var v struct {
			FileID   string          `json:"file_id"`
			Caption  json.RawMessage `json:"caption"`
			Autoplay bool            `json:"autoplay"`
			Loop     bool            `json:"loop"`
		}
		_ = json.Unmarshal(data, &v)
		caption, _ := parseText(v.Caption)
		return BlockVideo{v.FileID, caption, v.Autoplay, v.Loop}, nil

	case "audio":
		var v struct {
			FileID  string          `json:"file_id"`
			Caption json.RawMessage `json:"caption"`
		}
		_ = json.Unmarshal(data, &v)
		caption, _ := parseText(v.Caption)
		return BlockAudio{v.FileID, caption}, nil

	case "animation":
		var v struct {
			FileID  string          `json:"file_id"`
			Caption json.RawMessage `json:"caption"`
		}
		_ = json.Unmarshal(data, &v)
		caption, _ := parseText(v.Caption)
		return BlockAnimation{v.FileID, caption}, nil

	case "voice_note":
		var v struct {
			FileID  string          `json:"file_id"`
			Caption json.RawMessage `json:"caption"`
		}
		_ = json.Unmarshal(data, &v)
		caption, _ := parseText(v.Caption)
		return BlockVoiceNote{v.FileID, caption}, nil

	case "divider":
		return BlockDivider{}, nil

	case "mathematical_expression":
		var v struct {
			Expression string `json:"expression"`
		}
		_ = json.Unmarshal(data, &v)
		return BlockMathematicalExpression{v.Expression}, nil

	case "anchor":
		var v struct {
			Name string `json:"name"`
		}
		_ = json.Unmarshal(data, &v)
		return BlockAnchor{v.Name}, nil

	default:
		// forward-compat: неизвестный тип с text → BlockWrap, без text → ошибка.
		if text, err := parseText(head.Text); err == nil && text != nil {
			return BlockWrap{Tag: head.Type, Text: text}, nil
		}
		return nil, fmt.Errorf("richblock: unknown type %q", head.Type)
	}
}

// UnmarshalMessage разбирает корневой RichMessage из JSON.
func UnmarshalMessage(data []byte) (RichMessage, error) {
	var raw struct {
		Blocks json.RawMessage `json:"blocks"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return RichMessage{}, fmt.Errorf("richmessage: %w", err)
	}
	blocks, err := unmarshalBlocks(raw.Blocks)
	if err != nil {
		return RichMessage{}, err
	}
	return RichMessage{blocks}, nil
}

// ---------------------------------------------------------------------------
// Внутренние хелперы
// ---------------------------------------------------------------------------

func unmarshalBlocks(raw json.RawMessage) ([]RichBlock, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var raws []json.RawMessage
	if err := json.Unmarshal(raw, &raws); err != nil {
		return nil, err
	}
	blocks := make([]RichBlock, len(raws))
	for i, r := range raws {
		b, err := UnmarshalBlock(r)
		if err != nil {
			return nil, err
		}
		blocks[i] = b
	}
	return blocks, nil
}

func unmarshalListItem(data []byte) (RichBlockListItem, error) {
	var raw struct {
		Blocks json.RawMessage `json:"blocks"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return RichBlockListItem{}, err
	}
	blocks, err := unmarshalBlocks(raw.Blocks)
	if err != nil {
		return RichBlockListItem{}, err
	}
	return RichBlockListItem{blocks}, nil
}

func unmarshalTableCell(data []byte) (RichBlockTableCell, error) {
	var raw struct {
		Content    json.RawMessage `json:"content"`
		ColumnSpan int             `json:"column_span"`
		RowSpan    int             `json:"row_span"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return RichBlockTableCell{}, err
	}
	content, err := unmarshalBlocks(raw.Content)
	if err != nil {
		return RichBlockTableCell{}, err
	}
	return RichBlockTableCell{content, raw.ColumnSpan, raw.RowSpan}, nil
}
