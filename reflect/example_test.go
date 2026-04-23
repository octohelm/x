package reflect_test

import (
	"fmt"
	stdreflect "reflect"

	reflectx "github.com/octohelm/x/reflect"
)

type sample struct{}

func ExampleFullTypeName() {
	fmt.Println(reflectx.FullTypeName(stdreflect.TypeFor[*sample]()))
	// Output:
	// *github.com/octohelm/x/reflect_test.sample
}

func ExampleDeref() {
	fmt.Println(reflectx.Deref(stdreflect.TypeFor[***sample]()).Name())
	// Output:
	// sample
}
