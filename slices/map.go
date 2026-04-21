package slices

// Map 将 list 中的每个元素映射为目标类型，并按原顺序返回结果。
func Map[E any, T any](list []E, m func(e E) T) []T {
	out := make([]T, len(list))
	for i := range list {
		out[i] = m(list[i])
	}
	return out
}
