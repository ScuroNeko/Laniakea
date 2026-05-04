package tgfmt

import (
	"reflect"
	"testing"
	"time"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

func TestMessageBuilder_BuildPlainText(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("Hello")
	b.Add(", ")
	b.Add("world")

	text, entities := b.Build()

	if text != "Hello, world" {
		t.Fatalf("text = %q, want %q", text, "Hello, world")
	}

	if len(entities) != 0 {
		t.Fatalf("entities len = %d, want 0", len(entities))
	}
}

func TestMessageBuilder_EntityOffsetsAreUTF16(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("Hi ")
	b.Add("👋") // 2 UTF-16 code units
	b.Add(" ")
	b.Add("world").Bold()

	text, entities := b.Build()

	if text != "Hi 👋 world" {
		t.Fatalf("text = %q, want %q", text, "Hi 👋 world")
	}

	want := []tgapi.MessageEntity{
		{
			Type:   tgapi.MessageEntityBold,
			Offset: 6, // H i space = 3, 👋 = 2, space = 1
			Length: 5,
		},
	}

	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("entities = %#v, want %#v", entities, want)
	}
}

func TestMessageBuilder_EntityLengthIsUTF16(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("👋").Bold()

	text, entities := b.Build()

	if text != "👋" {
		t.Fatalf("text = %q, want %q", text, "👋")
	}

	want := []tgapi.MessageEntity{
		{
			Type:   tgapi.MessageEntityBold,
			Offset: 0,
			Length: 2,
		},
	}

	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("entities = %#v, want %#v", entities, want)
	}
}

func TestMessageBuilder_MultipleEntitiesOnSameEntry(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("hello").Bold().Italic()

	_, entities := b.Build()

	want := []tgapi.MessageEntity{
		{
			Type:   tgapi.MessageEntityBold,
			Offset: 0,
			Length: 5,
		},
		{
			Type:   tgapi.MessageEntityItalic,
			Offset: 0,
			Length: 5,
		},
	}

	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("entities = %#v, want %#v", entities, want)
	}
}

func TestMessageBuilder_DoesNotDuplicateAfterRepeatedReads(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("hello").Bold()

	text1 := b.String()
	entities1 := b.Entities()

	text2 := b.String()
	entities2 := b.Entities()

	if text1 != text2 {
		t.Fatalf("texts differ: %q != %q", text1, text2)
	}

	if !reflect.DeepEqual(entities1, entities2) {
		t.Fatalf("entities differ: %#v != %#v", entities1, entities2)
	}

	want := []tgapi.MessageEntity{
		{
			Type:   tgapi.MessageEntityBold,
			Offset: 0,
			Length: 5,
		},
	}

	if !reflect.DeepEqual(entities2, want) {
		t.Fatalf("entities = %#v, want %#v", entities2, want)
	}
}

func TestMessageBuilder_AddEntityAfterStringMarksDirty(t *testing.T) {
	b := NewMessageBuilder()

	entry := b.Add("hello")

	if got := b.String(); got != "hello" {
		t.Fatalf("String() = %q, want %q", got, "hello")
	}

	entry.Bold()

	entities := b.Entities()

	want := []tgapi.MessageEntity{
		{
			Type:   tgapi.MessageEntityBold,
			Offset: 0,
			Length: 5,
		},
	}

	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("entities = %#v, want %#v", entities, want)
	}
}

func TestMessageBuilder_EntitiesReturnsCopy(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("hello").Bold()

	entities1 := b.Entities()
	entities1[0].Offset = 999

	entities2 := b.Entities()

	if entities2[0].Offset != 0 {
		t.Fatalf("Entities() did not return copy: offset = %d, want 0", entities2[0].Offset)
	}
}

func TestMessageBuilder_BuildReturnsEntitiesCopy(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("hello").Bold()

	_, entities1 := b.Build()
	entities1[0].Offset = 999

	_, entities2 := b.Build()

	if entities2[0].Offset != 0 {
		t.Fatalf("Build() did not return entities copy: offset = %d, want 0", entities2[0].Offset)
	}
}

func TestMessageBuilder_Reset(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("hello").Bold()

	if got := b.String(); got != "hello" {
		t.Fatalf("String() before Reset = %q, want %q", got, "hello")
	}

	b.Reset()

	text, entities := b.Build()

	if text != "" {
		t.Fatalf("text after Reset = %q, want empty", text)
	}

	if len(entities) != 0 {
		t.Fatalf("entities len after Reset = %d, want 0", len(entities))
	}

	b.Add("world").Italic()

	text, entities = b.Build()

	if text != "world" {
		t.Fatalf("text after reuse = %q, want %q", text, "world")
	}

	want := []tgapi.MessageEntity{
		{
			Type:   tgapi.MessageEntityItalic,
			Offset: 0,
			Length: 5,
		},
	}

	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("entities after reuse = %#v, want %#v", entities, want)
	}
}

func TestMessageBuilder_EmptyEntryDoesNotCreateEntity(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("").Bold()
	b.Add("x")

	text, entities := b.Build()

	if text != "x" {
		t.Fatalf("text = %q, want %q", text, "x")
	}

	if len(entities) != 0 {
		t.Fatalf("entities len = %d, want 0: %#v", len(entities), entities)
	}
}

func TestMessageBuilder_Link(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("OpenAI").Link("https://openai.com")

	_, entities := b.Build()

	want := []tgapi.MessageEntity{
		{
			Type:   tgapi.MessageEntityTextLink,
			Offset: 0,
			Length: 6,
			URL:    "https://openai.com",
		},
	}

	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("entities = %#v, want %#v", entities, want)
	}
}

func TestMessageBuilder_CodeBlockWithLanguage(t *testing.T) {
	b := NewMessageBuilder()

	b.Add("fmt.Println(\"hi\")").CodeBlockWithLanguage("go")

	_, entities := b.Build()

	want := []tgapi.MessageEntity{
		{
			Type:     tgapi.MessageEntityPre,
			Offset:   0,
			Length:   17,
			Language: "go",
		},
	}

	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("entities = %#v, want %#v", entities, want)
	}
}

func TestMessageBuilder_DateTimeFormat(t *testing.T) {
	b := NewMessageBuilder()

	ts := time.Unix(1772323200, 0)

	b.Add("date").DateTimeFormat(ts, "MMMM d, yyyy")

	_, entities := b.Build()

	want := []tgapi.MessageEntity{
		{
			Type:           tgapi.MessageEntityDateTime,
			Offset:         0,
			Length:         4,
			UnixTime:       1772323200,
			DateTimeFormat: "MMMM d, yyyy",
		},
	}

	if !reflect.DeepEqual(entities, want) {
		t.Fatalf("entities = %#v, want %#v", entities, want)
	}
}

func TestTelegramTextLen(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{
			name: "ascii",
			text: "hello",
			want: 5,
		},
		{
			name: "cyrillic",
			text: "привет",
			want: 6,
		},
		{
			name: "emoji",
			text: "👋",
			want: 2,
		},
		{
			name: "mixed",
			text: "a👋b",
			want: 4,
		},
		{
			name: "zwj sequence",
			text: "👨‍👩‍👧‍👦",
			want: 11,
		},
		{
			name: "flag",
			text: "🇫🇮",
			want: 4,
		},
		{
			name: "variation selector",
			text: "❤️",
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := telegramTextLen(tt.text)
			if got != tt.want {
				t.Fatalf("telegramTextLen(%q) = %d, want %d", tt.text, got, tt.want)
			}
		})
	}
}

func TestMessageBuilder_SimpleEntityTypes(t *testing.T) {
	tests := []struct {
		name string
		add  func(*MessageBuilderEntry)
		want tgapi.MessageEntityType
	}{
		{"mention", func(e *MessageBuilderEntry) { e.Mention() }, tgapi.MessageEntityMention},
		{"hashtag", func(e *MessageBuilderEntry) { e.Hashtag() }, tgapi.MessageEntityHashtag},
		{"cashtag", func(e *MessageBuilderEntry) { e.Cashtag() }, tgapi.MessageEntityCashtag},
		{"bot command", func(e *MessageBuilderEntry) { e.BotCommand() }, tgapi.MessageEntityBotCommand},
		{"email", func(e *MessageBuilderEntry) { e.Email() }, tgapi.MessageEntityEmail},
		{"phone", func(e *MessageBuilderEntry) { e.Phone() }, tgapi.MessageEntityPhoneNumber},
		{"bold", func(e *MessageBuilderEntry) { e.Bold() }, tgapi.MessageEntityBold},
		{"italic", func(e *MessageBuilderEntry) { e.Italic() }, tgapi.MessageEntityItalic},
		{"underline", func(e *MessageBuilderEntry) { e.Underline() }, tgapi.MessageEntityUnderline},
		{"strikethrough", func(e *MessageBuilderEntry) { e.Strikethrough() }, tgapi.MessageEntityStrike},
		{"spoiler", func(e *MessageBuilderEntry) { e.Spoiler() }, tgapi.MessageEntitySpoiler},
		{"quote", func(e *MessageBuilderEntry) { e.Quote() }, tgapi.MessageEntityBlockquote},
		{"expandable quote", func(e *MessageBuilderEntry) { e.ExpandableQuote() }, tgapi.MessageEntityExpandableBlockquote},
		{"inline code", func(e *MessageBuilderEntry) { e.InlineCode() }, tgapi.MessageEntityCode},
		{"code block", func(e *MessageBuilderEntry) { e.CodeBlock() }, tgapi.MessageEntityPre},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewMessageBuilder()

			e := b.Add("hello")
			tt.add(e)

			_, entities := b.Build()

			want := []tgapi.MessageEntity{
				{
					Type:   tt.want,
					Offset: 0,
					Length: 5,
				},
			}

			if !reflect.DeepEqual(entities, want) {
				t.Fatalf("entities = %#v, want %#v", entities, want)
			}
		})
	}
}
