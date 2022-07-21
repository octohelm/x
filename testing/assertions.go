package testing

import (
	"fmt"
	"testing"

	"github.com/onsi/gomega/format"
)

func Expect[A any](t testing.TB, actual A, matcheres ...Matcher[A]) {
	t.Helper()
	for i := range matcheres {
		assert(t, actual, matcheres[i])
	}
}

func assert[A any](t testing.TB, actual A, m Matcher[A]) {
	ok := m.Match(actual)
	if m.Negative() {
		if !ok {
			return
		}
		t.Helper()
		t.Fatalf("\n" + failureMessage(actual, m))
		return
	}
	if ok {
		return
	}
	t.Helper()
	t.Fatalf("\n" + failureMessage(actual, m))
}

func failureMessage[A any](actual A, m Matcher[A]) string {
	if m.Negative() {
		return format.MessageWithDiff(
			fmt.Sprintf("%v", actual),
			fmt.Sprintf("Should not %s", m.Name()),
			fmt.Sprintf("%v", m.Expected()),
		)
	}
	return format.MessageWithDiff(
		fmt.Sprintf("%v", actual),
		fmt.Sprintf("Should %s", m.Name()),
		fmt.Sprintf("%v", m.Expected()),
	)
}

func init() {
	format.TruncateThreshold = 200
}
