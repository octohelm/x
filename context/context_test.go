package context

import (
	"context"
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestContext(t *testing.T) {
	c := New[string]()

	_, ok := c.MayFrom(context.Background())
	testingx.Expect(t, ok, testingx.BeFalse())

	ctx := c.Inject(context.Background(), "1")
	testingx.Expect(t, c.From(ctx), testingx.Be("1"))
}
