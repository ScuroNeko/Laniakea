package tgmd

import (
	"strconv"
	"strings"

	"git.scuroneko.dev/scuroneko/laniakea"
)

// Helpers in this file generate Telegram MarkdownV2.
// All user-provided text is escaped.

// TODO Markdown v2 escaping. GoDoc and tests

// WithBold returns s wrapped as bold Telegram Markdown text.
func WithBold(s string) string {
	return "*" + s + "*"
}

// WithItalic returns s wrapped as italic Telegram Markdown text.
func WithItalic(s string) string {
	return "_" + s + "_"
}

func WithUnderline(s string) string {
	return "__" + s + "__"
}
func WithStrikethrough(s string) string {
	return "~" + s + "~"
}
func WithSpoiler(s string) string {
	return "||" + s + "||"
}

// WithLink returns a Telegram Markdown link for text and URL.
func WithLink(text, url string) string {
	return "[" + text + "](" + url + ")"
}

func WithMention(text string, userID uint64) string {
	return "[" + text + "](tg://user?id=" + strconv.FormatUint(userID, 10) + ")"
}
func WithEmoji(text, emojiID string) string {
	return "[" + text + "](tg://emoji?id=" + emojiID + ")"
}

func WithTime(text string, unix uint64) string {
	return "![" + text + "](tg://time?unix=" + strconv.FormatUint(unix, 10) + ")"
}
func WithTimeFormat(text string, unix uint64, format string) string {
	return "![" + text + "](tg://time?unix=" +
		strconv.FormatUint(unix, 10) +
		"&format=" + format + ")"
}

// WithInlineCode returns s wrapped as inline code Telegram Markdown text.
func WithInlineCode(s string) string {
	return "`" + s + "`"
}
func WithBlockCode(s string) string {
	return "```\n" + s + "\n```"
}
func WithBlockCodeLanguage(s, lang string) string {
	return "```" + lang + "\n" + s + "\n```"
}
func WithQuote(s string) string {
	return ">" + strings.ReplaceAll(laniakea.EscapeMarkdownV2(s), "\n", "\n>")
}
func WithQuoteExpandable(s string) string {
	return "**>" + s
}
