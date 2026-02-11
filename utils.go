package laniakea

func Ptr[T any](v T) *T { return &v }

func Val[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}
