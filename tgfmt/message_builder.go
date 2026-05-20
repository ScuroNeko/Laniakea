package tgfmt

import (
	"strings"
	"time"

	"git.scuroneko.dev/scuroneko/extypes"
	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

// MessageBuilder builds Telegram message text with explicit message entities.
// MessageBuilder is not safe for concurrent use.
type MessageBuilder struct {
	str      string
	offset   int
	entities extypes.Slice[tgapi.MessageEntity]

	entries extypes.Slice[*MessageBuilderEntry]
	isDirty bool
}

// NewMessageBuilder returns an empty MessageBuilder.
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		entities: make([]tgapi.MessageEntity, 0),
		entries:  make(extypes.Slice[*MessageBuilderEntry], 0),
		isDirty:  false,
	}
}

// String returns the built message text.
func (b *MessageBuilder) String() string {
	if b.isDirty {
		b.update()
	}
	return b.str
}

// Entities returns a copy of the built message entities.
func (b *MessageBuilder) Entities() []tgapi.MessageEntity {
	if b.isDirty {
		b.update()
	}
	return append([]tgapi.MessageEntity(nil), b.entities...)
}

// Build returns the built message text and a copy of its entities.
func (b *MessageBuilder) Build() (string, []tgapi.MessageEntity) {
	if b.isDirty {
		b.update()
	}
	return b.str, append([]tgapi.MessageEntity(nil), b.entities...)
}

// Reset clears the builder and keeps it ready for reuse.
func (b *MessageBuilder) Reset() {
	b.str = ""
	b.offset = 0
	b.entities = b.entities[:0]
	b.entries = b.entries[:0]
	b.isDirty = false
}

func (b *MessageBuilder) update() *MessageBuilder {
	b.offset = 0

	var textLen int
	var entitiesLen int
	for _, e := range b.entries {
		textLen += len(e.text)
		entitiesLen += len(e.entities)
	}

	b.entities = make(extypes.Slice[tgapi.MessageEntity], 0, entitiesLen)

	var sb strings.Builder
	sb.Grow(textLen)

	for _, e := range b.entries {
		sb.WriteString(e.text)

		for _, entity := range e.entities {
			entity.Offset += b.offset
			b.entities = append(b.entities, entity)
		}

		b.offset += e.length
	}

	b.str = sb.String()
	b.isDirty = false
	return b
}

func (b *MessageBuilder) markDirty() {
	b.isDirty = true
}

// MessageBuilderEntry represents text appended to a MessageBuilder.
type MessageBuilderEntry struct {
	text   string
	length int

	b        *MessageBuilder
	entities extypes.Slice[tgapi.MessageEntity]
}

// Add appends plain text to the message and returns its entry for formatting.
func (b *MessageBuilder) Add(text string) *MessageBuilderEntry {
	e := &MessageBuilderEntry{
		b:        b,
		entities: make(extypes.Slice[tgapi.MessageEntity], 0),

		text:   text,
		length: telegramTextLen(text),
	}
	b.entries = b.entries.Push(e)
	b.markDirty()

	return e
}

// Mention marks the entry as a Telegram mention.
func (e *MessageBuilderEntry) Mention() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityMention,
		Offset: 0, Length: e.length,
	})
	return e
}

// Hashtag marks the entry as a Telegram hashtag.
func (e *MessageBuilderEntry) Hashtag() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityHashtag,
		Offset: 0, Length: e.length,
	})
	return e
}

// Cashtag marks the entry as a Telegram cashtag.
func (e *MessageBuilderEntry) Cashtag() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityCashtag,
		Offset: 0, Length: e.length,
	})
	return e
}

// BotCommand marks the entry as a Telegram bot command.
func (e *MessageBuilderEntry) BotCommand() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityBotCommand,
		Offset: 0, Length: e.length,
	})
	return e
}

// Email marks the entry as an email address.
func (e *MessageBuilderEntry) Email() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityEmail,
		Offset: 0, Length: e.length,
	})
	return e
}

// Phone marks the entry as a phone number.
func (e *MessageBuilderEntry) Phone() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityPhoneNumber,
		Offset: 0, Length: e.length,
	})
	return e
}

// Bold marks the entry as bold text.
func (e *MessageBuilderEntry) Bold() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityBold,
		Offset: 0, Length: e.length,
	})
	return e
}

// Italic marks the entry as italic text.
func (e *MessageBuilderEntry) Italic() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityItalic,
		Offset: 0, Length: e.length,
	})
	return e
}

// Underline marks the entry as underlined text.
func (e *MessageBuilderEntry) Underline() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityUnderline,
		Offset: 0, Length: e.length,
	})
	return e
}

// Strikethrough marks the entry as strikethrough text.
func (e *MessageBuilderEntry) Strikethrough() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityStrike,
		Offset: 0, Length: e.length,
	})
	return e
}

// Spoiler marks the entry as spoiler text.
func (e *MessageBuilderEntry) Spoiler() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntitySpoiler,
		Offset: 0, Length: e.length,
	})
	return e
}

// Quote marks the entry as a blockquote.
func (e *MessageBuilderEntry) Quote() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityBlockquote,
		Offset: 0, Length: e.length,
	})
	return e
}

// ExpandableQuote marks the entry as an expandable blockquote.
func (e *MessageBuilderEntry) ExpandableQuote() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityExpandableBlockquote,
		Offset: 0, Length: e.length,
	})
	return e
}

// InlineCode marks the entry as inline code.
func (e *MessageBuilderEntry) InlineCode() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityCode,
		Offset: 0, Length: e.length,
	})
	return e
}

// CodeBlock marks the entry as a preformatted code block.
func (e *MessageBuilderEntry) CodeBlock() *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityPre,
		Offset: 0, Length: e.length,
	})
	return e
}

// CodeBlockWithLanguage marks the entry as a preformatted code block with language.
func (e *MessageBuilderEntry) CodeBlockWithLanguage(lang string) *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityPre,
		Offset: 0, Length: e.length, Language: lang,
	})
	return e
}

// Link marks the entry as a text link.
func (e *MessageBuilderEntry) Link(url string) *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityTextLink,
		Offset: 0, Length: e.length, URL: url,
	})
	return e
}

// TextMention marks the entry as a mention of user.
func (e *MessageBuilderEntry) TextMention(user *tgapi.User) *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityTextMention,
		Offset: 0, Length: e.length, User: user,
	})
	return e
}

// CustomEmoji marks the entry as a custom emoji.
func (e *MessageBuilderEntry) CustomEmoji(emojiID string) *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityCustomEmoji,
		Offset: 0, Length: e.length, CustomEmojiID: emojiID,
	})
	return e
}

// DateTime marks the entry as a localized timestamp.
func (e *MessageBuilderEntry) DateTime(time time.Time) *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityDateTime,
		Offset: 0, Length: e.length, UnixTime: time.Unix(),
	})
	return e
}

// DateTimeFormat marks the entry as a localized timestamp with format.
func (e *MessageBuilderEntry) DateTimeFormat(time time.Time, format string) *MessageBuilderEntry {
	e.addEntity(tgapi.MessageEntity{
		Type:   tgapi.MessageEntityDateTime,
		Offset: 0, Length: e.length,
		UnixTime: time.Unix(), DateTimeFormat: format,
	})
	return e
}

func telegramTextLen(text string) int {
	n := 0
	for _, r := range text {
		if r <= 0xFFFF {
			n++
		} else {
			n += 2
		}
	}
	return n
}

func (e *MessageBuilderEntry) addEntity(entity tgapi.MessageEntity) {
	if entity.Length <= 0 {
		return
	}
	e.entities = append(e.entities, entity)
	if e.b != nil {
		e.b.markDirty()
	}
}
