package encoding_test

import (
	"fmt"

	"github.com/octohelm/x/encoding"
)

func ExampleMarshalText() {
	data, err := encoding.MarshalText([]byte("hi"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	// Output:
	// aGk=
}

func ExampleUnmarshalText() {
	var text string
	if err := encoding.UnmarshalText(&text, []byte("hello")); err != nil {
		panic(err)
	}

	fmt.Println(text)
	// Output:
	// hello
}
