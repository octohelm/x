package list_test

import (
	"fmt"

	"github.com/octohelm/x/container/list"
)

func ExampleList() {
	l := list.New[string]()
	l.PushBack("a")
	l.PushBack("b")
	l.PushFront("z")

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	// Output:
	// z
	// a
	// b
}
