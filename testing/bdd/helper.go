package bdd

import (
	"iter"
)

func Cases[C any](cases ...C) iter.Seq[C] {
	return func(yield func(C) bool) {
		for _, c := range cases {
			yield(c)
		}
	}
}

func Build[T any](build func(v *T)) *T {
	v := new(T)
	build(v)
	return v
}
