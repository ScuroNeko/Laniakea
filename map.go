package laniakea

type HashMap[K comparable, V any] map[K]V

func NewMapFrom[K comparable, V any](m map[K]V) HashMap[K, V] {
	out := make(HashMap[K, V])
	for k, v := range m {
		out[k] = v
	}
	return out
}
func (m HashMap[K, V]) Len() int {
	return len(m)
}
func (m HashMap[K, V]) Add(key K, val V) HashMap[K, V] {
	m[key] = val
	return m
}
func (m HashMap[K, V]) Remove(key K) HashMap[K, V] {
	delete(m, key)
	return m
}
func (m HashMap[K, V]) Get(key K) V {
	return m[key]
}
func (m HashMap[K, V]) GetOrDefault(key K, def V) V {
	v, ok := m[key]
	if !ok {
		return def
	}
	return v
}
func (m HashMap[K, V]) Contains(key K) bool {
	_, ok := m[key]
	return ok
}
func (m HashMap[K, V]) Filter(f func(K, V) bool) HashMap[K, V] {
	out := make(HashMap[K, V])
	for k, v := range m {
		if !f(k, v) {
			continue
		}
		out[k] = v
	}
	return out
}
