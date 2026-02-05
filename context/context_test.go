package context

import (
	"context"
	"testing"

	"github.com/octohelm/x/cmp"
	. "github.com/octohelm/x/testing/v2"
)

func TestContext(t *testing.T) {
	c := New[string]()

	t.Run("GIVEN a new context container", func(t *testing.T) {
		t.Run("WHEN fetching from empty context", func(t *testing.T) {
			_, ok := c.MayFrom(context.Background())

			Then(t, "should not be found",
				Expect(ok, Be(cmp.False())),
			)
		})

		t.Run("WHEN fetching after injected", func(t *testing.T) {
			ctx := c.Inject(context.Background(), "1")

			Then(t, "value should be extracted correctly",
				Expect(c.From(ctx), Equal("1")),
			)
		})
	})
}
