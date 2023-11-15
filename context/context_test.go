package context

import (
	"context"
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestContext(t *testing.T) {
	c := New[string]()

	ctx := c.Inject(context.Background(), "1")
	testingx.Expect(t, c.From(ctx), testingx.Be("1"))
}
