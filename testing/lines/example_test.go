package lines_test

import (
	"fmt"

	"github.com/octohelm/x/testing/lines"
)

func ExampleFromBytes() {
	fmt.Println(lines.FromBytes([]byte("a\nb\n")))
	// Output:
	// [a b]
}

func ExampleDiff() {
	diff := lines.Diff(
		lines.Lines{"a", "b"},
		lines.Lines{"a", "c"},
	)

	fmt.Print(string(diff))
	// Output:
	//
	// @@ -2,1 +2,1 @@
	// -b
	// +c
}
