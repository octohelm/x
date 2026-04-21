package v2

import "github.com/octohelm/x/testing/internal"

// Must 执行动作并在返回 error 时立即失败。
func Must(t TB, action func() error) {
	r := internal.Helper(1, &Reporter{})

	if err := action(); err != nil {
		r.Fatal(t, err)
	}
}

// MustValue 执行动作并返回其值；如果返回 error 则立即失败。
func MustValue[T any](t TB, action func() (T, error)) T {
	r := internal.Helper(1, &Reporter{})

	x, err := action()
	if err != nil {
		r.Fatal(t, err)
	}
	return x
}

// MustValues 执行动作并返回两个值；如果返回 error 则立即失败。
func MustValues[A any, B any](t TB, action func() (A, B, error)) (A, B) {
	r := internal.Helper(1, &Reporter{})

	a, b, err := action()
	if err != nil {
		r.Fatal(t, err)
	}
	return a, b
}
