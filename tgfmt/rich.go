package tgfmt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

// Rich is an inline fragment of rich-message HTML (Bot API 10.1).
// Raw text enters through NewRich, which escapes it; fragments compose as-is.
type Rich string

// RichBlock is a block-level fragment of rich-message HTML. Block
// constructors accept only Rich arguments, so invalid nesting (a block
// inside inline content) does not compile.
type RichBlock string

// RichItem is any rich-message fragment: Rich or RichBlock. Both are valid
// at the top level of a message — Telegram merges adjacent inline content
// into paragraphs.
type RichItem interface{ richItem() string }

func (r Rich) richItem() string      { return string(r) }
func (r RichBlock) richItem() string { return string(r) }

// RichHTML concatenates fragments into the final rich-message HTML string.
func RichHTML(items ...RichItem) string {
	var b strings.Builder
	for _, item := range items {
		b.WriteString(item.richItem())
	}
	return b.String()
}

// RichMessage builds a ready-to-send InputRichMessage from fragments.
// SkipEntityDetection is enabled so the server does not add auto-detected
// entities; set IsRTL on the result if needed.
func RichMessage(items ...RichItem) tgapi.InputRichMessage {
	return tgapi.InputRichMessage{
		HTML:                RichHTML(items...),
		SkipEntityDetection: true,
	}
}

func richJoin(items ...Rich) Rich {
	var out Rich
	for _, item := range items {
		out += item
	}
	return out
}
func richJoinSep(sep Rich, items ...Rich) Rich {
	var out Rich
	for i, item := range items {
		if i > 0 {
			out += sep
		}
		out += item
	}
	return out
}
func richBlocksJoin(items ...RichBlock) RichBlock {
	var out RichBlock
	for _, item := range items {
		out += item
	}
	return out
}

// openTag omits the space when there are no attributes.
func openTag(name string, attrs []string) string {
	if len(attrs) == 0 {
		return "<" + name + ">"
	}
	return "<" + name + " " + strings.Join(attrs, " ") + ">"
}

func cite(credit Rich) Rich {
	if credit == "" {
		return ""
	}
	return "<cite>" + credit + "</cite>"
}

// NewRich escapes raw text and returns it as an inline fragment.
func NewRich(text string) Rich { return Rich(escapeHTML(text)) }

// Bold wraps the fragment in <b>.
func (r Rich) Bold() Rich { return "<b>" + r + "</b>" }

// Italic wraps the fragment in <i>.
func (r Rich) Italic() Rich { return "<i>" + r + "</i>" }

// Underline wraps the fragment in <u>.
func (r Rich) Underline() Rich { return "<u>" + r + "</u>" }

// Strike wraps the fragment in <s>.
func (r Rich) Strike() Rich { return "<s>" + r + "</s>" }

// Code wraps the fragment in <code>.
func (r Rich) Code() Rich { return "<code>" + r + "</code>" }

// Mark wraps the fragment in <mark>.
func (r Rich) Mark() Rich { return "<mark>" + r + "</mark>" }

// Sub wraps the fragment in <sub>.
func (r Rich) Sub() Rich { return "<sub>" + r + "</sub>" }

// Sup wraps the fragment in <sup>.
func (r Rich) Sup() Rich { return "<sup>" + r + "</sup>" }

// Spoiler wraps the fragment in <tg-spoiler>.
func (r Rich) Spoiler() Rich { return "<tg-spoiler>" + r + "</tg-spoiler>" }

// Link wraps the fragment in a hyperlink to url.
func (r Rich) Link(url string) Rich { return `<a href="` + escapeRich(url) + `">` + r + "</a>" }

// Email wraps the fragment in a mailto: link.
func (r Rich) Email(email string) Rich { return r.Link("mailto:" + email) }

// Phone wraps the fragment in a tel: link.
func (r Rich) Phone(phone string) Rich { return r.Link("tel:" + phone) }

// Mention wraps the fragment in an inline user mention link.
func (r Rich) Mention(userID int64) Rich {
	return r.Link("tg://user?id=" + strconv.FormatInt(userID, 10))
}

// Anchor marks the fragment as a named anchor (<a name>).
func (r Rich) Anchor(name string) Rich { return `<a name="` + escapeRich(name) + `">` + r + "</a>" }

// AnchorLink wraps the fragment in an in-document link to a named anchor
// or reference; the server resolves which one by the target name.
func (r Rich) AnchorLink(anchor string) Rich { return r.Link("#" + anchor) }

// Ref wraps the fragment in a <tg-reference> to the named reference.
func (r Rich) Ref(ref string) Rich {
	return `<tg-reference name="` + escapeRich(ref) + `">` + r + "</tg-reference>"
}

// Emoji builds a custom emoji fragment with alt text as fallback.
func Emoji(emojiID, alt string) Rich {
	return `<tg-emoji emoji-id="` + escapeRich(emojiID) + `">` + escapeRich(alt) + `</tg-emoji>`
}

// Time marks the fragment as a <tg-time> bound to t.
func (r Rich) Time(t time.Time) Rich {
	return `<tg-time unix="` + Rich(strconv.FormatInt(t.Unix(), 10)) + `">` + r + "</tg-time>"
}

// TimeFormat marks the fragment as a <tg-time> with an explicit display format.
func (r Rich) TimeFormat(t time.Time, format string) Rich {
	return `<tg-time unix="` + Rich(strconv.FormatInt(t.Unix(), 10)) + `" format="` + escapeRich(format) + `">` + r + "</tg-time>"
}

// Math wraps the fragment in an inline <tg-math> expression.
func (r Rich) Math() Rich { return `<tg-math>` + r + `</tg-math>` }

// Br returns a line break fragment.
func Br() Rich { return "<br>" }

// H1 builds a level-1 heading block.
func H1(items ...Rich) RichBlock { return RichBlock("<h1>" + richJoin(items...) + "</h1>") }

// H2 builds a level-2 heading block.
func H2(items ...Rich) RichBlock { return RichBlock("<h2>" + richJoin(items...) + "</h2>") }

// H3 builds a level-3 heading block.
func H3(items ...Rich) RichBlock { return RichBlock("<h3>" + richJoin(items...) + "</h3>") }

// H4 builds a level-4 heading block.
func H4(items ...Rich) RichBlock { return RichBlock("<h4>" + richJoin(items...) + "</h4>") }

// H5 builds a level-5 heading block.
func H5(items ...Rich) RichBlock { return RichBlock("<h5>" + richJoin(items...) + "</h5>") }

// H6 builds a level-6 heading block.
func H6(items ...Rich) RichBlock { return RichBlock("<h6>" + richJoin(items...) + "</h6>") }

// P builds a paragraph block.
func P(items ...Rich) RichBlock { return RichBlock("<p>" + richJoin(items...) + "</p>") }

// Pre builds a preformatted code block.
func Pre(items ...Rich) RichBlock { return RichBlock("<pre>" + richJoin(items...) + "</pre>") }

// PreCode builds a preformatted code block tagged with a language.
func PreCode(lang string, items ...Rich) RichBlock {
	return RichBlock(`<pre><code class="language-` + escapeRich(lang) + `">` + richJoin(items...) + `</code></pre>`)
}

// Footer builds a footer block.
func Footer(items ...Rich) RichBlock { return RichBlock("<footer>" + richJoin(items...) + "</footer>") }

// Hr builds a divider block.
func Hr() RichBlock { return "<hr/>" }

// AnchorBlock builds a standalone named anchor between blocks: <a name></a>.
func AnchorBlock(name string) RichBlock {
	return RichBlock(`<a name="` + escapeHTML(name) + `"></a>`)
}

// LiItem is a list item under construction for Ul or Ol.
type LiItem struct {
	text     Rich
	value    int
	typ      string
	checkbox bool
	checked  bool
}

// Li builds a list item from inline fragments.
func Li(items ...Rich) LiItem { return LiItem{text: richJoin(items...)} }

// LiCheckbox builds a checkbox list item.
func LiCheckbox(checked bool, items ...Rich) LiItem {
	return LiItem{text: richJoin(items...), checkbox: true, checked: checked}
}

// SetValue sets the explicit ordinal of the item (like <li value>).
func (l LiItem) SetValue(val int) LiItem {
	l.value = val
	return l
}

// SetType sets the numbering type of the item: "1", "a", "A", "i", "I".
func (l LiItem) SetType(t string) LiItem {
	l.typ = t
	return l
}
func (l LiItem) build() Rich {
	if l.checkbox {
		input := Rich(`<input type="checkbox">`)
		if l.checked {
			input = `<input type="checkbox" checked>`
		}
		return "<li>" + input + l.text + "</li>"
	}
	attrs := make([]string, 0, 2)
	if l.value != 0 {
		attrs = append(attrs, `value="`+strconv.Itoa(l.value)+`"`)
	}
	if l.typ != "" {
		attrs = append(attrs, `type="`+escapeHTML(l.typ)+`"`)
	}
	return Rich(openTag("li", attrs)) + l.text + "</li>"
}
func joinLiItems(items []LiItem) Rich {
	var out Rich
	for _, item := range items {
		out += item.build()
	}
	return out
}

// Ul builds an unordered list block.
func Ul(items ...LiItem) RichBlock {
	return RichBlock(`<ul>` + joinLiItems(items) + `</ul>`)
}

// OlOpts holds the <ol> numbering attributes.
type OlOpts struct {
	Start    int
	Type     string
	Reversed bool
}

// Ol builds an ordered list block. Item labels are rendered by the server.
func Ol(opts OlOpts, items ...LiItem) RichBlock {
	attrs := make([]string, 0, 3)
	if opts.Start > 0 {
		attrs = append(attrs, `start="`+strconv.Itoa(opts.Start)+`"`)
	}
	if opts.Type != "" {
		attrs = append(attrs, `type="`+escapeHTML(opts.Type)+`"`)
	}
	if opts.Reversed {
		attrs = append(attrs, "reversed")
	}
	return RichBlock(openTag("ol", attrs) + string(joinLiItems(items)) + "</ol>")
}

// Blockquote builds a block quotation: lines are joined with <br> (as in the
// official HTML example) and credit renders as a trailing <cite>.
func Blockquote(credit Rich, lines ...Rich) RichBlock {
	return RichBlock(`<blockquote>` + richJoinSep(Br(), lines...) + cite(credit) + `</blockquote>`)
}

// Aside builds a pull quote (<aside>) with an optional <cite> credit.
func Aside(credit Rich, lines ...Rich) RichBlock {
	return RichBlock(`<aside>` + richJoinSep(Br(), lines...) + cite(credit) + `</aside>`)
}

// RichMedia is a media element for standalone blocks, collages, and
// slideshows. Media is sent by HTTP(S) URL only; file_id does not work in
// html mode.
type RichMedia struct {
	tag     string
	src     string
	spoiler bool
}

// Photo builds a photo element from an HTTP(S) URL.
func Photo(url string) RichMedia { return RichMedia{tag: "img", src: url} }

// Video builds a video or animation element from an HTTP(S) URL; the server
// distinguishes them by the URL extension.
func Video(url string) RichMedia { return RichMedia{tag: "video", src: url} }

// Audio builds an audio or voice-note element from an HTTP(S) URL; the
// server treats .ogg as a voice note.
func Audio(url string) RichMedia { return RichMedia{tag: "audio", src: url} }

// SetSpoiler hides the media behind a spoiler overlay.
func (m RichMedia) SetSpoiler() RichMedia {
	m.spoiler = true
	return m
}
func (m RichMedia) build() string {
	attrs := []string{`src="` + escapeHTML(m.src) + `"`}
	if m.spoiler {
		attrs = append(attrs, "tg-spoiler")
	}
	if m.tag == "img" {
		return "<img " + strings.Join(attrs, " ") + "/>"
	}
	return openTag(m.tag, attrs) + "</" + m.tag + ">"
}

// Block turns the media into a standalone block without a caption.
func (m RichMedia) Block() RichBlock { return RichBlock(m.build()) }

// Caption wraps the media in <figure> with a caption and optional credit.
func (m RichMedia) Caption(credit Rich, caption ...Rich) RichBlock {
	return RichBlock(`<figure>` + m.build() + string(figcaption(credit, caption...)) + `</figure>`)
}

func figcaption(credit Rich, caption ...Rich) Rich {
	text := richJoin(caption...) + cite(credit)
	if text == "" {
		return ""
	}
	return "<figcaption>" + text + "</figcaption>"
}

func richMediaJoin(items []RichMedia) string {
	var out string
	for _, item := range items {
		out += item.build()
	}
	return out
}

// Map builds a location map block.
func Map(lat, long float64, zoom int) RichBlock {
	latString := fmt.Sprintf("%.6f", lat)
	longString := fmt.Sprintf("%.6f", long)
	zoomString := strconv.Itoa(zoom)
	return RichBlock(`<tg-map lat="` + latString + `" long="` + longString + `" zoom="` + zoomString + `"/>`)
}

// MapCaption builds a map block wrapped in <figure> with a caption.
func MapCaption(lat, long float64, zoom int, credit Rich, caption ...Rich) RichBlock {
	return `<figure>` + Map(lat, long, zoom) + RichBlock(figcaption(credit, caption...)) + `</figure>`
}

// Collage builds a media collage block.
func Collage(items ...RichMedia) RichBlock {
	return RichBlock(`<tg-collage>` + richMediaJoin(items) + `</tg-collage>`)
}

// CollageCaption builds a media collage block with a caption.
func CollageCaption(credit Rich, caption Rich, items ...RichMedia) RichBlock {
	return RichBlock(`<tg-collage>` + richMediaJoin(items) + string(figcaption(credit, caption)) + `</tg-collage>`)
}

// Slideshow builds a media slideshow block.
func Slideshow(items ...RichMedia) RichBlock {
	return RichBlock(`<tg-slideshow>` + richMediaJoin(items) + `</tg-slideshow>`)
}

// SlideshowCaption builds a media slideshow block with a caption.
func SlideshowCaption(credit Rich, caption Rich, items ...RichMedia) RichBlock {
	return RichBlock(`<tg-slideshow>` + richMediaJoin(items) + string(figcaption(credit, caption)) + `</tg-slideshow>`)
}

// RichCell is a table cell under construction for Row.
type RichCell struct {
	text    Rich
	colspan int
	rowspan int
	align   string
	valign  string
}

// Cell builds a table cell from inline fragments.
func Cell(items ...Rich) RichCell {
	return RichCell{text: richJoin(items...)}
}

// SetSpan sets colspan and rowspan; zero leaves the attribute out.
func (r RichCell) SetSpan(col, row int) RichCell {
	r.colspan = col
	r.rowspan = row
	return r
}

// SetAlign sets horizontal alignment: "left", "center", or "right".
func (r RichCell) SetAlign(align string) RichCell {
	r.align = align
	return r
}

// SetVAlign sets vertical alignment: "top", "middle", or "bottom".
func (r RichCell) SetVAlign(align string) RichCell {
	r.valign = align
	return r
}
func (r RichCell) build(isHeader bool) string {
	attrs := make([]string, 0, 4)
	if r.colspan > 0 {
		attrs = append(attrs, `colspan="`+strconv.Itoa(r.colspan)+`"`)
	}
	if r.rowspan > 0 {
		attrs = append(attrs, `rowspan="`+strconv.Itoa(r.rowspan)+`"`)
	}
	if r.align != "" {
		attrs = append(attrs, `align="`+escapeHTML(r.align)+`"`)
	}
	if r.valign != "" {
		attrs = append(attrs, `valign="`+escapeHTML(r.valign)+`"`)
	}
	tag := "td"
	if isHeader {
		tag = "th"
	}
	return openTag(tag, attrs) + string(r.text) + "</" + tag + ">"
}

// RichRow is a table row under construction for Table.
type RichRow struct {
	cells    []RichCell
	isHeader bool
}

// Row builds a table row; isHeader renders every cell as <th>.
func Row(isHeader bool, cells ...RichCell) RichRow {
	return RichRow{cells: cells, isHeader: isHeader}
}
func (r RichRow) build() string {
	var buildCells string
	for _, cell := range r.cells {
		buildCells += cell.build(r.isHeader)
	}
	return `<tr>` + buildCells + `</tr>`
}

// Table builds a table block; an empty caption is omitted.
func Table(bordered, striped bool, caption Rich, rows ...RichRow) RichBlock {
	attrs := make([]string, 0, 2)
	if bordered {
		attrs = append(attrs, "bordered")
	}
	if striped {
		attrs = append(attrs, "striped")
	}
	out := openTag("table", attrs)
	if caption != "" {
		out += `<caption>` + string(caption) + `</caption>`
	}
	for _, row := range rows {
		out += row.build()
	}
	out += `</table>`
	return RichBlock(out)
}

// Details builds an expandable block with an inline summary.
func Details(isOpen bool, summary Rich, blocks ...RichBlock) RichBlock {
	tag := `<details>`
	if isOpen {
		tag = `<details open>`
	}
	return RichBlock(tag) + `<summary>` + RichBlock(summary) + `</summary>` + richBlocksJoin(blocks...) + `</details>`
}

// MathBlock builds a block-level mathematical expression.
func MathBlock(expr string) RichBlock {
	return RichBlock(`<tg-math-block>` + escapeHTML(expr) + `</tg-math-block>`)
}
