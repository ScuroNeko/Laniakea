package tgapi

import (
	"encoding/json"
	"fmt"
)

// Rich messages (Bot API 10.1), receive side: the RichText*/RichBlock* types
// mirror what the server sends in Message.rich_message, plus their parsers.
// These types intentionally have no constructors: sending goes only through
// InputRichMessage (html/markdown), and HTML generation lives in tgfmt
// (rich.go). Names follow the API objects; the exception is
// RichBlockQuotation (officially RichBlockBlockQuotation, the double Block
// is dropped).

// RichText is a node of the rich formatted text tree: a plain string, an
// array, or one of the typed objects below.
type RichText interface {
	isRichText()
}

// ---------------------------------------------------------------------------
// Base forms: string and array
// ---------------------------------------------------------------------------

// RichTextPlain is a plain text leaf.
type RichTextPlain string

func (RichTextPlain) isRichText() {}

// RichTextArray is a concatenation of rich text nodes.
type RichTextArray []RichText

func (RichTextArray) isRichText() {}

// ---------------------------------------------------------------------------
// Nodes with only a text field. There are 9; only the tag differs.
// bold italic underline strikethrough spoiler subscript superscript marked code
// ---------------------------------------------------------------------------

// RichTextWrap covers all "pure" wrapper nodes with a single type.
type RichTextWrap struct {
	Tag  string // "bold", "italic", ...
	Text RichText
}

func (RichTextWrap) isRichText() {}

func (w RichTextWrap) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string   `json:"type"`
		Text RichText `json:"text"`
	}{w.Tag, w.Text})
}

var richTextWrapTags = map[string]bool{
	"bold": true, "italic": true, "underline": true,
	"strikethrough": true, "spoiler": true, "subscript": true,
	"superscript": true, "marked": true, "code": true,
}

// ---------------------------------------------------------------------------
// Nodes with text + one extra string field.
// ---------------------------------------------------------------------------

// RichTextURL is rich text linking to a URL.
type RichTextURL struct {
	Text RichText
	URL  string
}

func (RichTextURL) isRichText() {}
func (v RichTextURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string   `json:"type"`
		Text RichText `json:"text"`
		URL  string   `json:"url"`
	}{"url", v.Text, v.URL})
}

// RichTextEmailAddress is rich text linking to an email address.
type RichTextEmailAddress struct {
	Text         RichText
	EmailAddress string
}

func (RichTextEmailAddress) isRichText() {}
func (v RichTextEmailAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type         string   `json:"type"`
		Text         RichText `json:"text"`
		EmailAddress string   `json:"email_address"`
	}{"email_address", v.Text, v.EmailAddress})
}

// RichTextPhoneNumber is rich text linking to a phone number.
type RichTextPhoneNumber struct {
	Text        RichText
	PhoneNumber string
}

func (RichTextPhoneNumber) isRichText() {}
func (v RichTextPhoneNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type        string   `json:"type"`
		Text        RichText `json:"text"`
		PhoneNumber string   `json:"phone_number"`
	}{"phone_number", v.Text, v.PhoneNumber})
}

// RichTextBankCardNumber is rich text marked as a bank card number.
type RichTextBankCardNumber struct {
	Text           RichText
	BankCardNumber string
}

func (RichTextBankCardNumber) isRichText() {}
func (v RichTextBankCardNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type           string   `json:"type"`
		Text           RichText `json:"text"`
		BankCardNumber string   `json:"bank_card_number"`
	}{"bank_card_number", v.Text, v.BankCardNumber})
}

// RichTextMention is rich text mentioning a user by username.
type RichTextMention struct {
	Text     RichText
	Username string
}

func (RichTextMention) isRichText() {}
func (v RichTextMention) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string   `json:"type"`
		Text     RichText `json:"text"`
		Username string   `json:"username"`
	}{"mention", v.Text, v.Username})
}

// RichTextHashtag is rich text marked as a hashtag.
type RichTextHashtag struct {
	Text    RichText
	Hashtag string
}

func (RichTextHashtag) isRichText() {}
func (v RichTextHashtag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string   `json:"type"`
		Text    RichText `json:"text"`
		Hashtag string   `json:"hashtag"`
	}{"hashtag", v.Text, v.Hashtag})
}

// RichTextCashtag is rich text marked as a cashtag.
type RichTextCashtag struct {
	Text    RichText
	Cashtag string
}

func (RichTextCashtag) isRichText() {}
func (v RichTextCashtag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string   `json:"type"`
		Text    RichText `json:"text"`
		Cashtag string   `json:"cashtag"`
	}{"cashtag", v.Text, v.Cashtag})
}

// RichTextBotCommand is rich text marked as a bot command.
type RichTextBotCommand struct {
	Text       RichText
	BotCommand string
}

func (RichTextBotCommand) isRichText() {}
func (v RichTextBotCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string   `json:"type"`
		Text       RichText `json:"text"`
		BotCommand string   `json:"bot_command"`
	}{"bot_command", v.Text, v.BotCommand})
}

// RichTextAnchorLink is rich text linking to a named anchor in the same message.
type RichTextAnchorLink struct {
	Text       RichText
	AnchorName string
}

func (RichTextAnchorLink) isRichText() {}
func (v RichTextAnchorLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string   `json:"type"`
		Text       RichText `json:"text"`
		AnchorName string   `json:"anchor_name"`
	}{"anchor_link", v.Text, v.AnchorName})
}

// RichTextReference is rich text marked as a named reference target.
type RichTextReference struct {
	Text RichText
	Name string
}

func (RichTextReference) isRichText() {}
func (v RichTextReference) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string   `json:"type"`
		Text RichText `json:"text"`
		Name string   `json:"name"`
	}{"reference", v.Text, v.Name})
}

// RichTextReferenceLink is rich text linking to a named reference.
type RichTextReferenceLink struct {
	Text          RichText
	ReferenceName string
}

func (RichTextReferenceLink) isRichText() {}
func (v RichTextReferenceLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type          string   `json:"type"`
		Text          RichText `json:"text"`
		ReferenceName string   `json:"reference_name"`
	}{"reference_link", v.Text, v.ReferenceName})
}

// ---------------------------------------------------------------------------
// Nodes with text + multiple/non-string fields.
// ---------------------------------------------------------------------------

// RichTextDateTime is rich text bound to a point in time with a display format.
type RichTextDateTime struct {
	Text           RichText
	UnixTime       int64
	DateTimeFormat string
}

func (RichTextDateTime) isRichText() {}
func (v RichTextDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type           string   `json:"type"`
		Text           RichText `json:"text"`
		UnixTime       int64    `json:"unix_time"`
		DateTimeFormat string   `json:"date_time_format"`
	}{"date_time", v.Text, v.UnixTime, v.DateTimeFormat})
}

// RichTextTextMention is rich text mentioning a user without a username.
type RichTextTextMention struct {
	Text RichText
	User User
}

func (RichTextTextMention) isRichText() {}
func (v RichTextTextMention) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string   `json:"type"`
		Text RichText `json:"text"`
		User User     `json:"user"`
	}{"text_mention", v.Text, v.User})
}

// ---------------------------------------------------------------------------
// LEAVES: no text field.
// ---------------------------------------------------------------------------

// RichTextCustomEmoji is a custom emoji leaf with alternative text.
type RichTextCustomEmoji struct {
	CustomEmojiID   string
	AlternativeText string
}

func (RichTextCustomEmoji) isRichText() {}
func (v RichTextCustomEmoji) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type            string `json:"type"`
		CustomEmojiID   string `json:"custom_emoji_id"`
		AlternativeText string `json:"alternative_text"`
	}{"custom_emoji", v.CustomEmojiID, v.AlternativeText})
}

// RichTextMathematicalExpression is an inline mathematical expression leaf.
type RichTextMathematicalExpression struct {
	Expression string
}

func (RichTextMathematicalExpression) isRichText() {}
func (v RichTextMathematicalExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string `json:"type"`
		Expression string `json:"expression"`
	}{"mathematical_expression", v.Expression})
}

// RichTextAnchor is a named anchor leaf that anchor links can point to.
type RichTextAnchor struct {
	Name string
}

func (RichTextAnchor) isRichText() {}
func (v RichTextAnchor) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{"anchor", v.Name})
}

// ---------------------------------------------------------------------------
// JSON -> RichText parsing
// ---------------------------------------------------------------------------

// UnmarshalRichText parses a RichText tree from JSON: a string, an array, or
// a typed object. Unknown object types that carry a text field are preserved
// as RichTextWrap for forward compatibility.
func UnmarshalRichText(data []byte) (RichText, error) {
	// 1. string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return RichTextPlain(s), nil
	}
	// 2. array
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err == nil {
		arr := make(RichTextArray, len(raw))
		for i, it := range raw {
			rt, err := UnmarshalRichText(it)
			if err != nil {
				return nil, err
			}
			arr[i] = rt
		}
		return arr, nil
	}
	// 3. object -> dispatch on type, grabbing the raw text along the way
	var head struct {
		Type string          `json:"type"`
		Text json.RawMessage `json:"text"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return nil, fmt.Errorf("richtext: not a string, array or object: %w", err)
	}

	// Recursively parse the nested text, if any.
	var inner RichText
	if len(head.Text) > 0 {
		var err error
		if inner, err = UnmarshalRichText(head.Text); err != nil {
			return nil, fmt.Errorf("richtext %q: bad text: %w", head.Type, err)
		}
	}

	if richTextWrapTags[head.Type] {
		return RichTextWrap{Tag: head.Type, Text: inner}, nil
	}

	switch head.Type {
	case "url":
		var v struct {
			URL string `json:"url"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}
		return RichTextURL{inner, v.URL}, nil
	case "email_address":
		var v struct {
			V string `json:"email_address"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextEmailAddress{inner, v.V}, nil
	case "phone_number":
		var v struct {
			V string `json:"phone_number"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextPhoneNumber{inner, v.V}, nil
	case "bank_card_number":
		var v struct {
			V string `json:"bank_card_number"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextBankCardNumber{inner, v.V}, nil
	case "mention":
		var v struct {
			V string `json:"username"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextMention{inner, v.V}, nil
	case "hashtag":
		var v struct {
			V string `json:"hashtag"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextHashtag{inner, v.V}, nil
	case "cashtag":
		var v struct {
			V string `json:"cashtag"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextCashtag{inner, v.V}, nil
	case "bot_command":
		var v struct {
			V string `json:"bot_command"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextBotCommand{inner, v.V}, nil
	case "anchor_link":
		var v struct {
			V string `json:"anchor_name"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextAnchorLink{inner, v.V}, nil
	case "reference":
		var v struct {
			V string `json:"name"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextReference{inner, v.V}, nil
	case "reference_link":
		var v struct {
			V string `json:"reference_name"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextReferenceLink{inner, v.V}, nil
	case "date_time":
		var v struct {
			UnixTime       int64  `json:"unix_time"`
			DateTimeFormat string `json:"date_time_format"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextDateTime{inner, v.UnixTime, v.DateTimeFormat}, nil
	case "text_mention":
		var v struct {
			User User `json:"user"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextTextMention{inner, v.User}, nil

	// --- leaves without text ---
	case "custom_emoji":
		var v struct {
			ID  string `json:"custom_emoji_id"`
			Alt string `json:"alternative_text"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextCustomEmoji{v.ID, v.Alt}, nil
	case "mathematical_expression":
		var v struct {
			Expression string `json:"expression"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextMathematicalExpression{v.Expression}, nil
	case "anchor":
		var v struct {
			Name string `json:"name"`
		}
		_ = json.Unmarshal(data, &v)
		return RichTextAnchor{v.Name}, nil

	default:
		// forward-compat: keep an unknown tag with a text field as
		// RichTextWrap; without text it is an error (the shape cannot be guessed).
		if inner != nil {
			return RichTextWrap{Tag: head.Type, Text: inner}, nil
		}
		return nil, fmt.Errorf("richtext: unknown type %q", head.Type)
	}
}
