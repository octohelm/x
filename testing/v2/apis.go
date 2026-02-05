package v2

import (
	"testing"

	"github.com/octohelm/x/testing/internal"
)

type (
	TB           = internal.TB
	TBController = internal.TBController
)

type Reporter = internal.Reporter

type Checker interface {
	Check(t TB)
}

type ValueChecker[V any] interface {
	Check(t TB, actual V)
}

func Then(t *testing.T, summary string, checkers ...Checker) {
	if t.Skipped() {
		return
	}
	t.Helper()

	t.Run("THEN "+summary, func(t *testing.T) {
		if t.Skipped() {
			return
		}
		t.Helper()

		for _, c := range checkers {
			c.Check(t)
		}
	})
}
