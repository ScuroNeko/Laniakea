package laniakea

// DictEntry {key:{ru:123,en:123}}
type DictEntry map[string]string
type L10n struct {
	entries      map[string]DictEntry
	fallbackLang string
}

func NewL10n(fallbackLanguage string) *L10n {
	return &L10n{make(map[string]DictEntry), fallbackLanguage}
}
func (l *L10n) AddDictEntry(key string, value DictEntry) *L10n {
	l.entries[key] = value
	return l
}
func (l *L10n) GetFallbackLanguage() string {
	return l.fallbackLang
}
func (l *L10n) Translate(lang, key string) string {
	s, ok := l.entries[key]
	if !ok {
		return key
	}
	return s[lang]
}
