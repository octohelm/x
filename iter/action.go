package iter

import "iter"

func Action[V any](do func(yield func(*V) bool) error) iter.Seq2[*V, error] {
	return func(yield func(*V, error) bool) {
		if err := do(func(v *V) bool { return yield(v, nil) }); err != nil {
			yield(nil, err)
			return
		}
	}
}
