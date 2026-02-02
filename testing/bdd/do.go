package bdd

func Do(t TB, action func() error, args ...any) {
	if x, ok := t.(WithHelper); ok {
		x.Helper()
	}

	if err := action(); err != nil {
		if len(args) > 0 {
			t.Fatal(append([]any{err}, args...)...)
		} else {
			t.Fatal(err)
		}
	}
}

func DoValue[T any](t TB, action func() (T, error), args ...any) T {
	if x, ok := t.(WithHelper); ok {
		x.Helper()
	}

	x, err := action()
	if err != nil {
		if len(args) > 0 {
			t.Fatal(append([]any{err}, args...)...)
		} else {
			t.Fatal(err)
		}
	}
	return x
}

func DoValues[A any, B any](t TB, action func() (A, B, error), args ...any) (A, B) {
	if x, ok := t.(WithHelper); ok {
		x.Helper()
	}

	a, b, err := action()
	if err != nil {
		if len(args) > 0 {
			t.Fatal(append([]any{err}, args...)...)
		} else {
			t.Fatal(err)
		}
	}
	return a, b
}
