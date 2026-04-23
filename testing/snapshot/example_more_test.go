package snapshot_test

import (
	"fmt"

	"github.com/octohelm/x/testing/snapshot"
)

func ExampleContext_Load() {
	ctx := &snapshot.Context{Name: "User Profile"}
	s, err := ctx.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println(ctx.Filename)
	fmt.Println(s.IsZero())
	// Output:
	// testdata/__snapshots__/user_profile.txtar
	// true
}

func ExampleFilesSeq() {
	raw := []byte("-- a.txt --\na\n-- b.txt --\nb\n")

	for file := range snapshot.FilesSeq(raw) {
		fmt.Printf("%s=%s\n", file.Name, string(file.Data))
	}
	// Output:
	// a.txt=a
	//
	// b.txt=b
}
