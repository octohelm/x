package iter

import "iter"

// Action 将返回 error 的 yield 风格回调包装为 iter.Seq2。
func Action[V any](do func(yield func(*V) bool) error) iter.Seq2[*V, error] {
	return func(yield func(*V, error) bool) {
		if err := do(func(v *V) bool { return yield(v, nil) }); err != nil {
			yield(nil, err)
			return
		}
	}
}
