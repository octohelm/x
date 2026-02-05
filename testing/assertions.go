package testing

import (
	"testing"

	"github.com/octohelm/x/testing/internal/deprecated"
)

func Expect[A any](t testing.TB, actual A, matchers ...Matcher[A]) {
	if t.Skipped() {
		return
	}

	for i := range matchers {
		deprecated.Expect(t, actual, matchers[i])
	}
}
