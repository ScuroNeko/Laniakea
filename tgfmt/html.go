package tgfmt

import (
	"strconv"
	"strings"
)

// HTML is an escaped Telegram HTML fragment.
//
// Methods on HTML compose formatting without escaping the fragment again.
type HTML string

// EscapeHTML escapes special characters for Telegram HTML parse mode.
func EscapeHTML(s string) HTML {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return HTML(s)
}

// Bold returns h wrapped as bold Telegram HTML text.
func (h HTML) Bold() HTML {
	return "<b>" + h + "</b>"
}

// Italic returns h wrapped as italic Telegram HTML text.
func (h HTML) Italic() HTML {
	return "<i>" + h + "</i>"
}

// Underline returns h wrapped as underlined Telegram HTML text.
func (h HTML) Underline() HTML {
	return "<u>" + h + "</u>"
}

// Strikethrough returns h wrapped as strikethrough Telegram HTML text.
func (h HTML) Strikethrough() HTML {
	return "<s>" + h + "</s>"
}

// Spoiler returns h wrapped as spoiler Telegram HTML text.
func (h HTML) Spoiler() HTML {
	return "<tg-spoiler>" + h + "</tg-spoiler>"
}

// Link returns h as a Telegram HTML text link.
func (h HTML) Link(url string) HTML {
	return `<a href="` + escapeHTMLAttr(url) + `">` + h + "</a>"
}

// Mention returns h as a Telegram HTML user mention.
func (h HTML) Mention(userID int64) HTML {
	return `<a href="tg://user?id=` + HTML(strconv.FormatInt(userID, 10)) + `">` + h + "</a>"
}

// Emoji returns h as a Telegram HTML custom emoji.
func (h HTML) Emoji(emojiID string) HTML {
	return `<tg-emoji emoji-id="` + escapeHTMLAttr(emojiID) + `">` + h + "</tg-emoji>"
}

// Time returns h as a Telegram HTML localized timestamp.
func (h HTML) Time(unix int64) HTML {
	return `<tg-time unix="` + HTML(strconv.FormatInt(unix, 10)) + `">` + h + "</tg-time>"
}

// TimeFormat returns h as a Telegram HTML localized timestamp with format.
func (h HTML) TimeFormat(unix int64, format string) HTML {
	return `<tg-time unix="` + HTML(strconv.FormatInt(unix, 10)) + `" format="` + escapeHTMLAttr(format) + `">` + h + "</tg-time>"
}

// InlineCode returns h wrapped as inline code Telegram HTML text.
func (h HTML) InlineCode() HTML {
	return "<code>" + h + "</code>"
}

// BlockCode returns h wrapped as a Telegram HTML code block.
func (h HTML) BlockCode() HTML {
	return "<pre>" + h + "</pre>"
}

// BlockCodeLanguage returns h wrapped as a Telegram HTML code block with language.
func (h HTML) BlockCodeLanguage(lang string) HTML {
	return `<pre><code class="language-` + escapeHTMLAttr(lang) + `">` + h + "</code></pre>"
}

// Quote returns h as a Telegram HTML blockquote.
func (h HTML) Quote() HTML {
	return "<blockquote>" + h + "</blockquote>"
}

// QuoteExpandable returns h as a Telegram HTML expandable blockquote.
func (h HTML) QuoteExpandable() HTML {
	return "<blockquote expandable>" + h + "</blockquote>"
}

func escapeHTMLAttr(s string) HTML {
	return EscapeHTML(s)
}
