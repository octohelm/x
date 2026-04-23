package types_test

import (
	"fmt"
	"reflect"

	"github.com/octohelm/x/types"
)

type exampleUser struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age"`
}

func ExampleFullTypeName() {
	t := types.FromRType(reflect.TypeFor[*exampleUser]())
	fmt.Println(types.FullTypeName(t))
	// Output:
	// *github.com/octohelm/x/types_test.exampleUser
}

func ExampleEachField() {
	t := types.FromRType(reflect.TypeFor[exampleUser]())

	types.EachField(t, "json", func(field types.StructField, displayName string, omitempty bool) bool {
		fmt.Printf("%s %t\n", displayName, omitempty)
		return true
	})
	// Output:
	// name true
	// age false
}
