package laniakea

import "slices"

type Slice[T comparable] []T

func NewSliceFrom[T comparable](slice []T) Slice[T] {
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
func (s Slice[T]) Remove(el T) Slice[T] {
	index := slices.Index(s, el)
	if index == -1 {
		return s
	}
	return s.Pop(index)
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

type Set[T comparable] []T

func NewSetFrom[T comparable](slice []T) Set[T] {
	s := make(Set[T], 0)
	for _, v := range slice {
		if slices.Index(slice, v) >= 0 {
			continue
		}
		s = append(s, v)
	}
	return s
}
func (s Set[T]) Len() int {
	return len(s)
}
func (s Set[T]) Cap() int {
	return cap(s)
}
func (s Set[T]) Get(index int) T {
	return s[index]
}
func (s Set[T]) Last() T {
	return s[s.Len()-1]
}
func (s Set[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Set[T]) Add(v T) Set[T] {
	index := slices.Index(s, v)
	if index >= 0 {
		return s
	}
	return append(s, v)
}
func (s Set[T]) Pop(index int) Set[T] {
	if index == 0 {
		return s[1:]
	}
	out := make(Set[T], s.Len()-index)
	for i, e := range s {
		if i == index {
			continue
		}
		out[i] = e
	}
	return out
}
func (s Set[T]) Remove(el T) Set[T] {
	index := slices.Index(s, el)
	if index == -1 {
		return s
	}
	return s.Pop(index)
}
func (s Set[T]) ToSlice() Slice[T] {
	out := make(Slice[T], s.Len())
	copy(out, s)
	return out
}
func (s Set[T]) ToArray() []T {
	out := make([]T, len(s))
	copy(out, s)
	return out
}
func (s Set[T]) ToAnyArray() []any {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}
