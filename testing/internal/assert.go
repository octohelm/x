package internal

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type CanPostMatched interface {
	PostMatched() error
}

func Expect[A any](t testing.TB, actual A, m Matcher[A]) {
	ok := m.Match(actual)
	if m.Negative() {
		if !ok {
			if x, ok := any(actual).(CanPostMatched); ok {
				if err := x.PostMatched(); err != nil {
					t.Helper()
					t.Fatal(err)
				}
			}
			return
		}
		t.Helper()
		t.Fatal("\n" + failureMessage(actual, m))
		return
	}
	if ok {
		if x, ok := any(actual).(CanPostMatched); ok {
			if err := x.PostMatched(); err != nil {
				t.Helper()
				t.Fatal(err)
			}
		}
		return
	}
	t.Helper()
	t.Fatal("\n" + failureMessage(actual, m))
}

func failureMessage[A any](actual A, m Matcher[A]) string {
	var v any = actual

	if n, ok := m.(MatcherWithActualNormalizer[A]); ok {
		v = n.NormalizeActual(actual)
	}

	if m.Negative() {
		return fmt.Sprintf("should not %s, but got\n%v", m.Action(), maybeDiff(v, m))
	}

	return fmt.Sprintf("should %s, but got\n%v", m.Action(), maybeDiff(v, m))
}

func maybeDiff(actual any, m any) any {
	if f, ok := m.(MatcherWithNormalizedExpected); ok {
		return cmp.Diff(f.NormalizedExpected(), actual)
	}
	return actual
}
