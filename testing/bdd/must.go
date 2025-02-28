package bdd

import "encoding"

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
