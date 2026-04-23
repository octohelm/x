package sync_test

import (
	"fmt"

	syncx "github.com/octohelm/x/sync"
)

func ExampleMap() {
	var m syncx.Map[string, int]
	m.Store("answer", 42)

	v, ok := m.Load("answer")
	fmt.Println(v, ok)
	// Output:
	// 42 true
}

func ExamplePool() {
	p := &syncx.Pool[[]int]{
		New: func() []int { return make([]int, 0, 2) },
	}

	v := p.Get()
	v = append(v, 1, 2)
	fmt.Println(v)
	p.Put(v[:0])

	fmt.Println(len(p.Get()))
	// Output:
	// [1 2]
	// 0
}
