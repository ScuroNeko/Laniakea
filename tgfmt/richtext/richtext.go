package richtext

import (
	"encoding/json"
	"fmt"

	"git.scuroneko.dev/scuroneko/laniakea/tgapi"
)

// RichText — узел дерева форматированного текста: строка, массив или
// один из типизированных объектов ниже.
type RichText interface {
	isRichText()
}

// ---------------------------------------------------------------------------
// Базовые формы: строка и массив
// ---------------------------------------------------------------------------

type String string

func (String) isRichText() {}

type Array []RichText

func (Array) isRichText() {}

// ---------------------------------------------------------------------------
// Узлы только с полем text. Их 9; различает только тег.
// bold italic underline strikethrough spoiler subscript superscript marked code
// ---------------------------------------------------------------------------

// Wrap покрывает все «чистые» оборачивающие узлы одним типом.
type Wrap struct {
	Tag  string // "bold", "italic", ...
	Text RichText
}

func (Wrap) isRichText() {}

func (w Wrap) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string   `json:"type"`
		Text RichText `json:"text"`
	}{w.Tag, w.Text})
}

var wrapTags = map[string]bool{
	"bold": true, "italic": true, "underline": true,
	"strikethrough": true, "spoiler": true, "subscript": true,
	"superscript": true, "marked": true, "code": true,
}

// Удобные конструкторы для wrap-узлов.
func Bold(t RichText) Wrap          { return Wrap{"bold", t} }
func Italic(t RichText) Wrap        { return Wrap{"italic", t} }
func Underline(t RichText) Wrap     { return Wrap{"underline", t} }
func Strikethrough(t RichText) Wrap { return Wrap{"strikethrough", t} }
func Spoiler(t RichText) Wrap       { return Wrap{"spoiler", t} }
func Subscript(t RichText) Wrap     { return Wrap{"subscript", t} }
func Superscript(t RichText) Wrap   { return Wrap{"superscript", t} }
func Marked(t RichText) Wrap        { return Wrap{"marked", t} }
func Code(t RichText) Wrap          { return Wrap{"code", t} }

// ---------------------------------------------------------------------------
// Узлы с text + одно строковое доп. поле.
// ---------------------------------------------------------------------------

type URL struct {
	Text RichText
	URL  string
}

func (URL) isRichText() {}
func (v URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string   `json:"type"`
		Text RichText `json:"text"`
		URL  string   `json:"url"`
	}{"url", v.Text, v.URL})
}

type EmailAddress struct {
	Text         RichText
	EmailAddress string
}

func (EmailAddress) isRichText() {}
func (v EmailAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type         string   `json:"type"`
		Text         RichText `json:"text"`
		EmailAddress string   `json:"email_address"`
	}{"email_address", v.Text, v.EmailAddress})
}

type PhoneNumber struct {
	Text        RichText
	PhoneNumber string
}

func (PhoneNumber) isRichText() {}
func (v PhoneNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type        string   `json:"type"`
		Text        RichText `json:"text"`
		PhoneNumber string   `json:"phone_number"`
	}{"phone_number", v.Text, v.PhoneNumber})
}

type BankCardNumber struct {
	Text           RichText
	BankCardNumber string
}

func (BankCardNumber) isRichText() {}
func (v BankCardNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type           string   `json:"type"`
		Text           RichText `json:"text"`
		BankCardNumber string   `json:"bank_card_number"`
	}{"bank_card_number", v.Text, v.BankCardNumber})
}

type Mention struct {
	Text     RichText
	Username string
}

func (Mention) isRichText() {}
func (v Mention) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string   `json:"type"`
		Text     RichText `json:"text"`
		Username string   `json:"username"`
	}{"mention", v.Text, v.Username})
}

type Hashtag struct {
	Text    RichText
	Hashtag string
}

func (Hashtag) isRichText() {}
func (v Hashtag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string   `json:"type"`
		Text    RichText `json:"text"`
		Hashtag string   `json:"hashtag"`
	}{"hashtag", v.Text, v.Hashtag})
}

type Cashtag struct {
	Text    RichText
	Cashtag string
}

func (Cashtag) isRichText() {}
func (v Cashtag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string   `json:"type"`
		Text    RichText `json:"text"`
		Cashtag string   `json:"cashtag"`
	}{"cashtag", v.Text, v.Cashtag})
}

type BotCommand struct {
	Text       RichText
	BotCommand string
}

func (BotCommand) isRichText() {}
func (v BotCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string   `json:"type"`
		Text       RichText `json:"text"`
		BotCommand string   `json:"bot_command"`
	}{"bot_command", v.Text, v.BotCommand})
}

type AnchorLink struct {
	Text       RichText
	AnchorName string
}

func (AnchorLink) isRichText() {}
func (v AnchorLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string   `json:"type"`
		Text       RichText `json:"text"`
		AnchorName string   `json:"anchor_name"`
	}{"anchor_link", v.Text, v.AnchorName})
}

type Reference struct {
	Text RichText
	Name string
}

func (Reference) isRichText() {}
func (v Reference) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string   `json:"type"`
		Text RichText `json:"text"`
		Name string   `json:"name"`
	}{"reference", v.Text, v.Name})
}

type ReferenceLink struct {
	Text          RichText
	ReferenceName string
}

func (ReferenceLink) isRichText() {}
func (v ReferenceLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type          string   `json:"type"`
		Text          RichText `json:"text"`
		ReferenceName string   `json:"reference_name"`
	}{"reference_link", v.Text, v.ReferenceName})
}

// ---------------------------------------------------------------------------
// Узлы с text + несколько/нестроковых полей.
// ---------------------------------------------------------------------------

type DateTime struct {
	Text           RichText
	UnixTime       int64
	DateTimeFormat string
}

func (DateTime) isRichText() {}
func (v DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type           string   `json:"type"`
		Text           RichText `json:"text"`
		UnixTime       int64    `json:"unix_time"`
		DateTimeFormat string   `json:"date_time_format"`
	}{"date_time", v.Text, v.UnixTime, v.DateTimeFormat})
}

type TextMention struct {
	Text RichText
	User tgapi.User
}

func (TextMention) isRichText() {}
func (v TextMention) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string     `json:"type"`
		Text RichText   `json:"text"`
		User tgapi.User `json:"user"`
	}{"text_mention", v.Text, v.User})
}

// ---------------------------------------------------------------------------
// ЛИСТЬЯ: без поля text.
// ---------------------------------------------------------------------------

type CustomEmoji struct {
	CustomEmojiID   string
	AlternativeText string
}

func (CustomEmoji) isRichText() {}
func (v CustomEmoji) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type            string `json:"type"`
		CustomEmojiID   string `json:"custom_emoji_id"`
		AlternativeText string `json:"alternative_text"`
	}{"custom_emoji", v.CustomEmojiID, v.AlternativeText})
}

type MathematicalExpression struct {
	Expression string
}

func (MathematicalExpression) isRichText() {}
func (v MathematicalExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string `json:"type"`
		Expression string `json:"expression"`
	}{"mathematical_expression", v.Expression})
}

type Anchor struct {
	Name string
}

func (Anchor) isRichText() {}
func (v Anchor) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{"anchor", v.Name})
}

// ---------------------------------------------------------------------------
// Разбор JSON -> RichText
// ---------------------------------------------------------------------------

func Unmarshal(data []byte) (RichText, error) {
	// 1. строка
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return String(s), nil
	}
	// 2. массив
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err == nil {
		arr := make(Array, len(raw))
		for i, it := range raw {
			rt, err := Unmarshal(it)
			if err != nil {
				return nil, err
			}
			arr[i] = rt
		}
		return arr, nil
	}
	// 3. объект -> смотрим type, попутно вытаскиваем сырой text
	var head struct {
		Type string          `json:"type"`
		Text json.RawMessage `json:"text"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return nil, fmt.Errorf("richtext: not a string, array or object: %w", err)
	}

	// Рекурсивно разбираем вложенный text, если он есть.
	var inner RichText
	if len(head.Text) > 0 {
		var err error
		if inner, err = Unmarshal(head.Text); err != nil {
			return nil, fmt.Errorf("richtext %q: bad text: %w", head.Type, err)
		}
	}

	if wrapTags[head.Type] {
		return Wrap{Tag: head.Type, Text: inner}, nil
	}

	switch head.Type {
	case "url":
		var v struct {
			URL string `json:"url"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return URL{inner, v.URL}, nil
	case "email_address":
		var v struct {
			V string `json:"email_address"`
		}
		_ = json.Unmarshal(data, &v)
		return EmailAddress{inner, v.V}, nil
	case "phone_number":
		var v struct {
			V string `json:"phone_number"`
		}
		_ = json.Unmarshal(data, &v)
		return PhoneNumber{inner, v.V}, nil
	case "bank_card_number":
		var v struct {
			V string `json:"bank_card_number"`
		}
		_ = json.Unmarshal(data, &v)
		return BankCardNumber{inner, v.V}, nil
	case "mention":
		var v struct {
			V string `json:"username"`
		}
		_ = json.Unmarshal(data, &v)
		return Mention{inner, v.V}, nil
	case "hashtag":
		var v struct {
			V string `json:"hashtag"`
		}
		_ = json.Unmarshal(data, &v)
		return Hashtag{inner, v.V}, nil
	case "cashtag":
		var v struct {
			V string `json:"cashtag"`
		}
		_ = json.Unmarshal(data, &v)
		return Cashtag{inner, v.V}, nil
	case "bot_command":
		var v struct {
			V string `json:"bot_command"`
		}
		_ = json.Unmarshal(data, &v)
		return BotCommand{inner, v.V}, nil
	case "anchor_link":
		var v struct {
			V string `json:"anchor_name"`
		}
		_ = json.Unmarshal(data, &v)
		return AnchorLink{inner, v.V}, nil
	case "reference":
		var v struct {
			V string `json:"name"`
		}
		_ = json.Unmarshal(data, &v)
		return Reference{inner, v.V}, nil
	case "reference_link":
		var v struct {
			V string `json:"reference_name"`
		}
		_ = json.Unmarshal(data, &v)
		return ReferenceLink{inner, v.V}, nil
	case "date_time":
		var v struct {
			UnixTime       int64  `json:"unix_time"`
			DateTimeFormat string `json:"date_time_format"`
		}
		_ = json.Unmarshal(data, &v)
		return DateTime{inner, v.UnixTime, v.DateTimeFormat}, nil
	case "text_mention":
		var v struct {
			User tgapi.User `json:"user"`
		}
		_ = json.Unmarshal(data, &v)
		return TextMention{inner, v.User}, nil

	// --- листья без text ---
	case "custom_emoji":
		var v struct {
			ID  string `json:"custom_emoji_id"`
			Alt string `json:"alternative_text"`
		}
		_ = json.Unmarshal(data, &v)
		return CustomEmoji{v.ID, v.Alt}, nil
	case "mathematical_expression":
		var v struct {
			Expression string `json:"expression"`
		}
		_ = json.Unmarshal(data, &v)
		return MathematicalExpression{v.Expression}, nil
	case "anchor":
		var v struct {
			Name string `json:"name"`
		}
		_ = json.Unmarshal(data, &v)
		return Anchor{v.Name}, nil

	default:
		// forward-compat: неизвестный тег с полем text сохраняем как Wrap,
		// без text — как ошибку (нельзя угадать форму).
		if inner != nil {
			return Wrap{Tag: head.Type, Text: inner}, nil
		}
		return nil, fmt.Errorf("richtext: unknown type %q", head.Type)
	}
}
