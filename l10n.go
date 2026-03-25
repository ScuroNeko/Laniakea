package laniakea

import "sync"

// DictEntry maps language codes to translated strings.
type DictEntry map[string]string

// L10n stores translations with a configurable fallback language and is safe for concurrent use.
type L10n struct {
	mu           sync.RWMutex
	entries      map[string]DictEntry
	fallbackLang string
}

// NewL10n creates a localization store with the given fallback language.
func NewL10n(fallbackLanguage string) *L10n {
	return &L10n{
		entries:      make(map[string]DictEntry),
		fallbackLang: fallbackLanguage,
	}
}

// AddDictEntry stores translations for key.
func (l *L10n) AddDictEntry(key string, value DictEntry) *L10n {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.entries == nil {
		l.entries = make(map[string]DictEntry)
	}
	l.entries[key] = cloneDictEntry(value)
	return l
}

// GetFallbackLanguage returns the currently configured fallback language code.
func (l *L10n) GetFallbackLanguage() string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.fallbackLang
}

// Translate returns the translation for key in lang, falling back to the configured language or the key itself.
func (l *L10n) Translate(lang, key string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	entries, exists := l.entries[key]
	if !exists {
		return key
	}

	if translation, ok := entries[lang]; ok {
		return translation
	}

	if fallback, ok := entries[l.fallbackLang]; ok {
		return fallback
	}

	return key
}

func cloneDictEntry(src DictEntry) DictEntry {
	if src == nil {
		return nil
	}
	cloned := make(DictEntry, len(src))
	for lang, text := range src {
		cloned[lang] = text
	}
	return cloned
}
