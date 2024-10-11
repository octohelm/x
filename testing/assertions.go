package testing

import (
	"fmt"
	"testing"

	"github.com/onsi/gomega/format"
)

func Expect[A any](t testing.TB, actual A, matchers ...Matcher[A]) {
	t.Helper()
	for i := range matchers {
		assert(t, actual, matchers[i])
	}
}

func assert[A any](t testing.TB, actual A, m Matcher[A]) {
	ok := m.Match(actual)
	if m.Negative() {
		if !ok {
			return
		}
		t.Helper()
		t.Fatal("\n" + failureMessage(actual, m))
		return
	}
	if ok {
		return
	}
	t.Helper()
	t.Fatal("\n" + failureMessage(actual, m))
}

func failureMessage[A any](actual A, m Matcher[A]) string {
	if m.Negative() {
		if f, ok := m.(ExpectedFormatter); ok {
			return format.MessageWithDiff(
				m.FormatActual(actual),
				fmt.Sprintf("Should not %s", m.Name()),
				f.FormatExpected(),
			)
		}

		return format.Message(actual, fmt.Sprintf("Should not %s", m.Name()))
	}

	if f, ok := m.(ExpectedFormatter); ok {
		return format.MessageWithDiff(
			m.FormatActual(actual),
			fmt.Sprintf("Should %s", m.Name()),
			f.FormatExpected(),
		)
	}

	return format.Message(actual, fmt.Sprintf("Should %s", m.Name()))
}

func init() {
	format.TruncateThreshold = 200
}
