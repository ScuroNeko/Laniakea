package tgfmt

import (
	"strconv"
	"strings"
)

// MarkdownV2 is an escaped Telegram MarkdownV2 fragment.
//
// Methods on MarkdownV2 compose formatting without escaping the fragment again.
type MarkdownV2 string

// EscapeMarkdownV2 escapes special characters for Telegram MarkdownV2.
// https://core.telegram.org/bots/api#markdownv2-style
func EscapeMarkdownV2(s string) MarkdownV2 {
	symbols := []string{"\\", "_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, symbol := range symbols {
		s = strings.ReplaceAll(s, symbol, "\\"+symbol)
	}
	return MarkdownV2(s)
}

// Bold returns s wrapped as bold Telegram MarkdownV2 text.
func (s MarkdownV2) Bold() MarkdownV2 {
	return "*" + s + "*"
}

// Italic returns s wrapped as italic Telegram MarkdownV2 text.
func (s MarkdownV2) Italic() MarkdownV2 {
	return "_" + s + "_"
}

// Underline returns s wrapped as underlined Telegram MarkdownV2 text.
func (s MarkdownV2) Underline() MarkdownV2 {
	return "__" + s + "__"
}

// Strikethrough returns s wrapped as strikethrough Telegram MarkdownV2 text.
func (s MarkdownV2) Strikethrough() MarkdownV2 {
	return "~" + s + "~"
}

// Spoiler returns s wrapped as spoiler Telegram MarkdownV2 text.
func (s MarkdownV2) Spoiler() MarkdownV2 {
	return "||" + s + "||"
}

// Link returns s as a Telegram MarkdownV2 text link.
func (s MarkdownV2) Link(url string) MarkdownV2 {
	return "[" + s + "](" + escapeMarkdownV2LinkDestination(url) + ")"
}

// Mention returns s as a Telegram MarkdownV2 user mention.
func (s MarkdownV2) Mention(userID uint64) MarkdownV2 {
	return "[" + s + "](tg://user?id=" + MarkdownV2(strconv.FormatUint(userID, 10)) + ")"
}

// Emoji returns s as a Telegram MarkdownV2 custom emoji.
func (s MarkdownV2) Emoji(emojiID string) MarkdownV2 {
	return "[" + s + "](tg://emoji?id=" + escapeMarkdownV2LinkDestination(emojiID) + ")"
}

// Time returns s as a Telegram MarkdownV2 localized timestamp.
func (s MarkdownV2) Time(unix uint64) MarkdownV2 {
	return "![" + s + "](tg://time?unix=" + MarkdownV2(strconv.FormatUint(unix, 10)) + ")"
}

// TimeFormat returns s as a Telegram MarkdownV2 localized timestamp with format.
func (s MarkdownV2) TimeFormat(unix uint64, format string) MarkdownV2 {
	dest := "tg://time?unix=" + strconv.FormatUint(unix, 10) + "&format=" + format
	return "![" + s + "](" + escapeMarkdownV2LinkDestination(dest) + ")"
}

// InlineCode returns s wrapped as inline code Telegram MarkdownV2 text.
func (s MarkdownV2) InlineCode() MarkdownV2 {
	return "`" + s + "`"
}

// BlockCode returns s wrapped as a Telegram MarkdownV2 code block.
func (s MarkdownV2) BlockCode() MarkdownV2 {
	return "```\n" + s + "\n```"
}

// BlockCodeLanguage returns s wrapped as a Telegram MarkdownV2 code block with language.
func (s MarkdownV2) BlockCodeLanguage(lang string) MarkdownV2 {
	return "```" + MarkdownV2(lang) + "\n" + s + "\n```"
}

// Quote returns s as a Telegram MarkdownV2 blockquote.
func (s MarkdownV2) Quote() MarkdownV2 {
	return MarkdownV2(">" + strings.ReplaceAll(string(s), "\n", "\n>"))
}

// QuoteExpandable returns s as a Telegram MarkdownV2 expandable blockquote.
func (s MarkdownV2) QuoteExpandable() MarkdownV2 {
	return "**>" + s
}

func escapeMarkdownV2LinkDestination(s string) MarkdownV2 {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ")", "\\)")
	return MarkdownV2(s)
}
