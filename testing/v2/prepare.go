package v2

import "github.com/octohelm/x/testing/internal"

func Must(t TB, action func() error) {
	r := internal.Helper(1, &Reporter{})

	if err := action(); err != nil {
		r.Fatal(t, err)
	}
}

func MustValue[T any](t TB, action func() (T, error)) T {
	r := internal.Helper(1, &Reporter{})

	x, err := action()
	if err != nil {
		r.Fatal(t, err)
	}
	return x
}

func MustValues[A any, B any](t TB, action func() (A, B, error)) (A, B) {
	r := internal.Helper(1, &Reporter{})

	a, b, err := action()
	if err != nil {
		r.Fatal(t, err)
	}
	return a, b
}
