package laniakea

// DictEntry represents a single localized entry with language-to-text mappings.
// Example: {"ru": "Привет", "en": "Hello"}.
type DictEntry map[string]string

// L10n is a localization manager that maps keys to language-specific strings.
type L10n struct {
	entries      map[string]DictEntry // Map of translation keys to language dictionaries
	fallbackLang string               // Language code to use when requested language is missing
}

// NewL10n creates a new L10n instance with the specified fallback language.
// The fallback language is used when a requested language is not available
// for a given key.
//
// Example: NewL10n("en") will return "Hello" for key "greeting" if "ru" is requested
// but no "ru" entry exists.
func NewL10n(fallbackLanguage string) *L10n {
	return &L10n{
		entries:      make(map[string]DictEntry),
		fallbackLang: fallbackLanguage,
	}
}

// AddDictEntry adds a new translation entry for the given key.
// The value must be a DictEntry mapping language codes (e.g., "en", "ru") to their translated strings.
//
// If a key already exists, it is overwritten.
//
// Returns the L10n instance for method chaining.
func (l *L10n) AddDictEntry(key string, value DictEntry) *L10n {
	l.entries[key] = value
	return l
}

// GetFallbackLanguage returns the currently configured fallback language code.
func (l *L10n) GetFallbackLanguage() string {
	return l.fallbackLang
}

// Translate retrieves the translation for the given key and language.
//
// Behavior:
//   - If the key exists and the language has a translation → returns the translation
//   - If the key exists but the language is missing → returns the fallback language's value
//   - If the key does not exist → returns the key string itself (as fallback)
//
// Example:
//
//	l.AddDictEntry("greeting", DictEntry{"en": "Hello", "ru": "Привет"})
//	l.Translate("en", "greeting") → "Hello"
//	l.Translate("es", "greeting") → "Hello" (fallback to "en")
//	l.Translate("en", "unknown") → "unknown" (key not found)
//
// This behavior ensures that missing translations do not break UI or logs —
// instead, the original key is displayed, making it easy to identify gaps.
func (l *L10n) Translate(lang, key string) string {
	entries, exists := l.entries[key]
	if !exists {
		return key // Return key as fallback when translation is missing
	}

	// Try requested language
	if translation, ok := entries[lang]; ok {
		return translation
	}

	// Fall back to configured fallback language
	if fallback, ok := entries[l.fallbackLang]; ok {
		return fallback
	}

	// If fallback language is also missing, return the key
	return key
}
