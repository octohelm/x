package cmp_test

import (
	"fmt"

	"github.com/octohelm/x/cmp"
)

func ExampleEq() {
	fmt.Println(cmp.Eq(3)(3) == nil)
	fmt.Println(cmp.Eq(3)(4) == nil)
	// Output:
	// true
	// false
}

func ExampleLen() {
	values := []string{"a", "b", "c"}

	fmt.Println(cmp.Len[[]string](3)(values) == nil)
	fmt.Println(cmp.Len[[]string](cmp.Gt(3))(values) == nil)
	// Output:
	// true
	// false
}
