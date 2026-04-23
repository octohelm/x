package singleflight

import "fmt"

func ExampleGroupValue() {
	var g GroupValue[string, int]

	v, err, shared := g.Do("answer", func() (int, error) {
		return 42, nil
	})

	fmt.Println(v, err == nil, shared)
	// Output:
	// 42 true false
}

func ExampleGroupValue_DoChan() {
	var g GroupValue[string, int]

	ch := g.DoChan("answer", func() (int, error) {
		return 42, nil
	})

	result := <-ch
	fmt.Println(result.Val, result.Err == nil, result.Shared)
	// Output:
	// 42 true false
}
