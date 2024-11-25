package anyjson

import "slices"

func Sorted[T Valuer](v T) T {
	switch x := any(v).(type) {
	case *Object:
		keys := slices.Sorted(func(yield func(string2 string) bool) {
			for k, _ := range x.KeyValues() {
				if !yield(k) {
					return
				}
			}
		})

		o := &Object{}

		for _, k := range keys {
			if propValue, ok := x.Get(k); ok {
				o.Set(k, Sorted(propValue))
			}
		}

		return any(o).(T)
	case *Array:
		arr := &Array{
			items: make([]Valuer, x.Len()),
		}

		for i := range arr.items {
			arr.items[i] = Sorted(x.items[i])
		}

		return any(arr).(T)
	}

	return v
}
