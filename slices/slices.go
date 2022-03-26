package slices

func Map[E any, T any](list []E, m func(e E) T) []T {
	var out = make([]T, len(list))
	for i := range list {
		out[i] = m(list[i])
	}
	return out
}
