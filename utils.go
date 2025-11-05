package laniakea

import "encoding/json"

func MapToStruct(m map[string]any, dst interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, dst)
	return err
}

func MapToJson(m map[string]any) (string, error) {
	data, err := json.Marshal(m)
	return string(data), err
}

func StructToMap(s interface{}) (map[string]any, error) {
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
