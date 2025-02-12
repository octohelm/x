package slices

func Filter[E any](list []E, filter func(e E) bool) []E {
	out := make([]E, 0, len(list))
	for _, e := range list {
		if filter(e) {
			out = append(out, e)
		}
	}
	return out
}
