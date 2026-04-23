package reflect_test

import (
	"fmt"

	reflectx "github.com/octohelm/x/reflect"
)

func ExampleParseStructTags() {
	tags := reflectx.ParseStructTags(`json:"name,omitempty" validate:"required"`)

	fmt.Println(tags["json"].Name())
	fmt.Println(tags["json"].HasFlag("omitempty"))
	fmt.Println(tags["validate"].Name())
	// Output:
	// name
	// true
	// required
}
