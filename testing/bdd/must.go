package bdd

import "encoding"

func MustDo[T any](action func() (T, error)) T {
	x, err := action()
	if err != nil {
		panic(err)
	}
	return x
}

func MustDo2[A any, B any](action func() (A, B, error)) (A, B) {
	a, b, err := action()
	if err != nil {
		panic(err)
	}
	return a, b
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Must2[A any, B any](a A, b B, err error) (A, B) {
	if err != nil {
		panic(err)
	}
	return a, b
}

func MustText[T any](text string) T {
	t := new(T)
	if u, ok := any(t).(encoding.TextUnmarshaler); ok {
		if err := u.UnmarshalText([]byte(text)); err != nil {
			panic(err)
		}
	}
	return *t
}
