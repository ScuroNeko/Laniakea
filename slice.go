package laniakea

type Slice[T any] []T

func NewSliceFrom[T any](slice []T) Slice[T] {
	s := make(Slice[T], len(slice))
	copy(s[:], slice)
	return s
}
func (s Slice[T]) Len() int {
	return len(s)
}
func (s Slice[T]) Cap() int {
	return cap(s)
}
func (s Slice[T]) Get(index int) T {
	return s[index]
}
func (s Slice[T]) Last() T {
	return s.Get(s.Len() - 1)
}
func (s Slice[T]) Swap(i, j int) Slice[T] {
	s[i], s[j] = s[j], s[i]
	return s
}
func (s Slice[T]) Filter(f func(e T) bool) Slice[T] {
	out := make(Slice[T], 0)
	for _, v := range s {
		if f(v) {
			out = append(out, v)
		}
	}
	return out
}
func (s Slice[T]) Map(f func(e T) T) Slice[T] {
	out := make(Slice[T], s.Len())
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}
func (s Slice[T]) Pop(index int) Slice[T] {
	if index == 0 {
		return s[1:]
	}
	out := make(Slice[T], s.Len()-index)
	for i, e := range s {
		if i == index {
			continue
		}
		out[i] = e
	}
	return out
}
func (s Slice[T]) Push(e T) Slice[T] {
	return append(s, e)
}

func (s Slice[T]) ToArray() []T {
	out := make([]T, len(s))
	copy(out, s)
	return out
}
func (s Slice[T]) ToAnyArray() []any {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}
