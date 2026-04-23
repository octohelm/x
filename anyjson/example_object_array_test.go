package anyjson_test

import (
	"fmt"

	"github.com/octohelm/x/anyjson"
)

func ExampleObject_KeyValues() {
	obj := &anyjson.Object{}
	obj.Set("name", anyjson.StringOf("octohelm"))
	obj.Set("count", anyjson.NumberOf(2))

	for key, value := range obj.KeyValues() {
		fmt.Printf("%s=%v\n", key, value.Value())
	}
	// Output:
	// name=octohelm
	// count=2
}

func ExampleArray_Index() {
	arr := &anyjson.Array{}
	arr.Append(anyjson.StringOf("go"))
	arr.Append(anyjson.StringOf("doc"))

	v, ok := arr.Index(1)
	fmt.Println(v.Value(), ok)
	// Output:
	// doc true
}
