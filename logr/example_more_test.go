package logr_test

import (
	"context"
	"fmt"

	"github.com/octohelm/x/logr"
)

func ExampleParseLevel() {
	level, err := logr.ParseLevel("warning")
	if err != nil {
		panic(err)
	}

	fmt.Println(level.String())
	// Output:
	// warning
}

func ExampleWithLogger() {
	ctx := logr.WithLogger(context.Background(), logr.Discard())

	_, ok := logr.LoggerFromContext(ctx)
	fmt.Println(ok)
	// Output:
	// true
}
