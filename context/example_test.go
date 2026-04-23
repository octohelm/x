package context_test

import (
	stdcontext "context"
	"fmt"

	xcontext "github.com/octohelm/x/context"
)

func ExampleNew() {
	slot := xcontext.New[string]()
	ctx := slot.Inject(stdcontext.Background(), "alice")

	fmt.Println(slot.From(ctx))
	// Output:
	// alice
}

func ExampleWithDefaults() {
	slot := xcontext.New(xcontext.WithDefaults("guest"))

	fmt.Println(slot.From(stdcontext.Background()))
	// Output:
	// guest
}
