package iter_test

import (
	"fmt"

	iterx "github.com/octohelm/x/iter"
)

func ExampleAction() {
	seq := iterx.Action(func(yield func(*int) bool) error {
		for _, v := range []int{1, 2, 3} {
			v := v
			if !yield(&v) {
				return nil
			}
		}
		return nil
	})

	for v, err := range seq {
		fmt.Println(*v, err == nil)
	}
	// Output:
	// 1 true
	// 2 true
	// 3 true
}
