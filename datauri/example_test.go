package datauri_test

import (
	"fmt"

	"github.com/octohelm/x/datauri"
)

func ExampleParse() {
	d, err := datauri.Parse("data:text/plain;charset=utf-8;base64,aGVsbG8=")
	if err != nil {
		panic(err)
	}

	fmt.Println(d.MediaType)
	fmt.Println(d.Params["charset"])
	fmt.Println(string(d.Data))
	// Output:
	// text/plain
	// utf-8
	// hello
}

func ExampleDataURI_Encoded() {
	d := &datauri.DataURI{
		MediaType: "text/plain",
		Data:      []byte("hello world"),
	}

	fmt.Println(d.Encoded(false))
	// Output:
	// data:text/plain,hello%20world
}
