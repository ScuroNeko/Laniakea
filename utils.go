package laniakea

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MapToStruct unsafe function
func MapToStruct(m map[string]any, s any) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, s)
	return err
}

// SliceToStruct unsafe function
func SliceToStruct(sl []any, s any) error {
	data, err := json.Marshal(sl)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, s)
	return err
}

// AnyToStruct unsafe function
func AnyToStruct(src, dest any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, dest)
	return err
}

func MapToJson(m map[string]any) (string, error) {
	data, err := json.Marshal(m)
	return string(data), err
}

func StructToMap(s any) (map[string]any, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	err = json.Unmarshal(data, &m)
	return m, err
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func EscapeMarkdown(s string) string {
	s = strings.ReplaceAll(s, "_", `\_`)
	s = strings.ReplaceAll(s, "*", `\*`)
	s = strings.ReplaceAll(s, "[", `\[`)
	return strings.ReplaceAll(s, "`", "\\`")
}

func EscapeMarkdownV2(s string) string {
	symbols := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, symbol := range symbols {
		s = strings.ReplaceAll(s, symbol, fmt.Sprintf("\\%s", symbol))
	}
	return s
}
