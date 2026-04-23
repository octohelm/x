package slices_test

import (
	"fmt"
	"strings"

	"github.com/octohelm/x/slices"
)

func ExampleMap() {
	upper := slices.Map([]string{"go", "doc"}, strings.ToUpper)
	fmt.Println(upper)
	// Output:
	// [GO DOC]
}

func ExampleFilter() {
	evens := slices.Filter([]int{1, 2, 3, 4, 5, 6}, func(v int) bool {
		return v%2 == 0
	})
	fmt.Println(evens)
	// Output:
	// [2 4 6]
}
