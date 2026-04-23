package snapshot_test

import (
	"fmt"

	"github.com/octohelm/x/testing/snapshot"
)

func ExampleNewSnapshot() {
	s := snapshot.NewSnapshot()
	s.Add("hello.txt", []byte("hello"))

	fmt.Print(string(s.Bytes()))
	// Output:
	// -- hello.txt --
	// hello
}

func ExampleFromFiles() {
	s := snapshot.FromFiles(
		snapshot.FileFromRaw("a.txt", []byte("a")),
		snapshot.FileFromRaw("b.txt", []byte("b")),
	)

	fmt.Println(s.Equal(s))
	// Output:
	// true
}
