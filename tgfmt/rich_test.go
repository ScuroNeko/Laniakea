package tgfmt

import (
	"testing"
	"time"
)

func TestRichInline(t *testing.T) {
	cases := []struct {
		name string
		got  Rich
		want string
	}{
		{"bold", NewRich("bold text").Bold(), "<b>bold text</b>"},
		{"spoiler", NewRich("spoiler").Spoiler(), "<tg-spoiler>spoiler</tg-spoiler>"},
		{"escape", NewRich(`a<b> & "c"`), "a&lt;b&gt; &amp; &quot;c&quot;"},
		{"link", NewRich("inline URL").Link("https://t.me/"), `<a href="https://t.me/">inline URL</a>`},
		{"link attr escape", NewRich("x").Link(`https://e/?q="><b>`), `<a href="https://e/?q=&quot;&gt;&lt;b&gt;">x</a>`},
		{"mention", NewRich("user").Mention(123456789), `<a href="tg://user?id=123456789">user</a>`},
		{"anchor", Rich("").Anchor("chapter-1"), `<a name="chapter-1"></a>`},
		{"anchor link", NewRich("in-document link").AnchorLink("chapter-1"), `<a href="#chapter-1">in-document link</a>`},
		{"reference", NewRich("Referenced text").Ref("note-1"), `<tg-reference name="note-1">Referenced text</tg-reference>`},
		{"emoji", Emoji("5368324170671202286", "👍"), `<tg-emoji emoji-id="5368324170671202286">👍</tg-emoji>`},
		{"emoji alt escape", Emoji("1", `<x>`), `<tg-emoji emoji-id="1">&lt;x&gt;</tg-emoji>`},
		{"time format", NewRich("22:45 tomorrow").TimeFormat(time.Unix(1647531900, 0), "wDT"),
			`<tg-time unix="1647531900" format="wDT">22:45 tomorrow</tg-time>`},
		{"math", NewRich("x^2 + y^2").Math(), "<tg-math>x^2 + y^2</tg-math>"},
	}
	for _, c := range cases {
		if string(c.got) != c.want {
			t.Errorf("%s:\n got  %s\n want %s", c.name, c.got, c.want)
		}
	}
}

func TestRichBlocks(t *testing.T) {
	cases := []struct {
		name string
		got  RichBlock
		want string
	}{
		{"h1", H1(NewRich("Heading 1")), "<h1>Heading 1</h1>"},
		{"p", P(NewRich("Paragraph text")), "<p>Paragraph text</p>"},
		{"pre code", PreCode("python", NewRich("print('x')")),
			`<pre><code class="language-python">print('x')</code></pre>`},
		{"footer", Footer(NewRich("Footer text")), "<footer>Footer text</footer>"},
		{"hr", Hr(), "<hr/>"},
		{"anchor block", AnchorBlock("chapter-2"), `<a name="chapter-2"></a>`},
		{"math block", MathBlock("E = mc^2"), "<tg-math-block>E = mc^2</tg-math-block>"},
	}
	for _, c := range cases {
		if string(c.got) != c.want {
			t.Errorf("%s:\n got  %s\n want %s", c.name, c.got, c.want)
		}
	}
}

func TestRichLists(t *testing.T) {
	cases := []struct {
		name string
		got  RichBlock
		want string
	}{
		{"ul", Ul(Li(NewRich("unordered list item"))),
			"<ul><li>unordered list item</li></ul>"},
		{"ol plain", Ol(OlOpts{}, Li(NewRich("ordered list item"))),
			"<ol><li>ordered list item</li></ol>"},
		{"ol attrs", Ol(OlOpts{Start: 3, Type: "a", Reversed: true}, Li(NewRich("ordered list item"))),
			`<ol start="3" type="a" reversed><li>ordered list item</li></ol>`},
		{"li value type", Ol(OlOpts{}, Li(NewRich("item")).SetValue(7).SetType("i")),
			`<ol><li value="7" type="i">item</li></ol>`},
		{"checkboxes", Ul(
			LiCheckbox(true, NewRich("Checked checkbox")),
			LiCheckbox(false, NewRich("Unchecked checkbox")),
		), `<ul><li><input type="checkbox" checked>Checked checkbox</li><li><input type="checkbox">Unchecked checkbox</li></ul>`},
	}
	for _, c := range cases {
		if string(c.got) != c.want {
			t.Errorf("%s:\n got  %s\n want %s", c.name, c.got, c.want)
		}
	}
}

func TestRichQuotes(t *testing.T) {
	got := Blockquote(NewRich("The Author"),
		NewRich("Block quotation started"),
		NewRich("Block quotation continued"),
		NewRich("The last line of the block quotation"),
	)
	want := "<blockquote>Block quotation started<br>Block quotation continued<br>" +
		"The last line of the block quotation<cite>The Author</cite></blockquote>"
	if string(got) != want {
		t.Errorf("blockquote:\n got  %s\n want %s", got, want)
	}

	got = Aside(NewRich("The Author"), NewRich("Pull quote"))
	want = "<aside>Pull quote<cite>The Author</cite></aside>"
	if string(got) != want {
		t.Errorf("aside:\n got  %s\n want %s", got, want)
	}

	got = Blockquote("", NewRich("no credit"))
	want = "<blockquote>no credit</blockquote>"
	if string(got) != want {
		t.Errorf("blockquote without credit:\n got  %s\n want %s", got, want)
	}
}

func TestRichMedia(t *testing.T) {
	cases := []struct {
		name string
		got  RichBlock
		want string
	}{
		{"photo", Photo("https://telegram.org/example/photo.jpg").Block(),
			`<img src="https://telegram.org/example/photo.jpg"/>`},
		{"video", Video("https://telegram.org/example/video.mp4").Block(),
			`<video src="https://telegram.org/example/video.mp4"></video>`},
		{"audio", Audio("https://telegram.org/example/audio.mp3").Block(),
			`<audio src="https://telegram.org/example/audio.mp3"></audio>`},
		{"photo spoiler caption", Photo("https://telegram.org/example/photo.jpg").SetSpoiler().
			Caption(NewRich("Photo credit"), NewRich("Photo caption")),
			`<figure><img src="https://telegram.org/example/photo.jpg" tg-spoiler/>` +
				`<figcaption>Photo caption<cite>Photo credit</cite></figcaption></figure>`},
		{"video caption no credit", Video("https://telegram.org/example/video.mp4").
			Caption("", NewRich("Video caption")),
			`<figure><video src="https://telegram.org/example/video.mp4"></video>` +
				`<figcaption>Video caption</figcaption></figure>`},
		{"map", Map(41.9, 12.5, 14), `<tg-map lat="41.900000" long="12.500000" zoom="14"/>`},
		{"map caption", MapCaption(41.9, 12.5, 14, "", NewRich("Map caption")),
			`<figure><tg-map lat="41.900000" long="12.500000" zoom="14"/><figcaption>Map caption</figcaption></figure>`},
		{"collage", Collage(
			Photo("https://telegram.org/example/photo.jpg"),
			Video("https://telegram.org/example/video.mp4"),
		), `<tg-collage><img src="https://telegram.org/example/photo.jpg"/>` +
			`<video src="https://telegram.org/example/video.mp4"></video></tg-collage>`},
		{"slideshow caption", SlideshowCaption("", NewRich("Slideshow caption"),
			Photo("https://telegram.org/example/photo.jpg"),
		), `<tg-slideshow><img src="https://telegram.org/example/photo.jpg"/>` +
			`<figcaption>Slideshow caption</figcaption></tg-slideshow>`},
	}
	for _, c := range cases {
		if string(c.got) != c.want {
			t.Errorf("%s:\n got  %s\n want %s", c.name, c.got, c.want)
		}
	}
}

func TestRichTable(t *testing.T) {
	got := Table(false, false, "",
		Row(true, Cell(NewRich("Header 1")), Cell(NewRich("Header 2"))),
		Row(false, Cell(NewRich("Value 1")), Cell(NewRich("Value 2"))),
	)
	want := "<table><tr><th>Header 1</th><th>Header 2</th></tr>" +
		"<tr><td>Value 1</td><td>Value 2</td></tr></table>"
	if string(got) != want {
		t.Errorf("plain table:\n got  %s\n want %s", got, want)
	}

	got = Table(true, true, NewRich("Table caption"),
		Row(false,
			Cell(NewRich("Value")).SetSpan(2, 2).SetAlign("left"),
			Cell(NewRich("Value2")).SetAlign("center"),
		),
	)
	want = `<table bordered striped><caption>Table caption</caption>` +
		`<tr><td colspan="2" rowspan="2" align="left">Value</td><td align="center">Value2</td></tr></table>`
	if string(got) != want {
		t.Errorf("table attrs:\n got  %s\n want %s", got, want)
	}
}

func TestRichDetails(t *testing.T) {
	got := Details(true, NewRich("Title"), P(NewRich("Content")))
	want := "<details open><summary>Title</summary><p>Content</p></details>"
	if string(got) != want {
		t.Errorf("details:\n got  %s\n want %s", got, want)
	}
}
