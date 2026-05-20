package tgfmt

import (
	"strconv"
	"strings"
)

// Markdown is an escaped legacy Telegram Markdown fragment.
//
// Deprecated: Use MarkdownV2 instead.
type Markdown string

// EscapeMarkdown escapes special characters for legacy Telegram Markdown.
//
// Deprecated: Use EscapeMarkdownV2 instead.
func EscapeMarkdown(s string) Markdown {
	s = strings.ReplaceAll(s, "_", `\_`)
	s = strings.ReplaceAll(s, "*", `\*`)
	s = strings.ReplaceAll(s, "[", `\[`)
	return Markdown(strings.ReplaceAll(s, "`", "\\`"))
}

// Bold returns s wrapped as bold legacy Telegram Markdown text.
func (s Markdown) Bold() Markdown {
	return "*" + s + "*"
}

// Italic returns s wrapped as italic legacy Telegram Markdown text.
func (s Markdown) Italic() Markdown {
	return "_" + s + "_"
}

// Link returns s as a legacy Telegram Markdown text link.
func (s Markdown) Link(url string) Markdown {
	return "[" + s + "](" + Markdown(url) + ")"
}

// Mention returns s as a legacy Telegram Markdown user mention.
func (s Markdown) Mention(userID int64) Markdown {
	return "[" + s + "](tg://user?id=" + Markdown(strconv.FormatInt(userID, 10)) + ")"
}

// InlineCode returns s wrapped as inline code legacy Telegram Markdown text.
func (s Markdown) InlineCode() Markdown {
	return "`" + s + "`"
}

// BlockCode returns s wrapped as a legacy Telegram Markdown code block.
func (s Markdown) BlockCode() Markdown {
	return "```\n" + s + "\n```"
}

// BlockCodeLanguage returns s wrapped as a legacy Telegram Markdown code block.
func (s Markdown) BlockCodeLanguage(lang string) Markdown {
	return "```" + Markdown(lang) + "\n" + s + "\n```"
}
