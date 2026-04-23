package anyjson_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/octohelm/x/anyjson"
)

func ExampleSorted() {
	v := anyjson.MustFromValue(anyjson.Obj{
		"b": 1,
		"a": anyjson.Obj{
			"d": 2,
			"c": 1,
		},
	})

	fmt.Println(anyjson.ToString(anyjson.Sorted(v)))
	// Output:
	// {"a":{"c":1,"d":2},"b":1}
}

func ExampleTransform() {
	v := anyjson.MustFromValue(anyjson.Obj{
		"name": "octohelm",
		"tags": anyjson.List{"go", "json"},
	})

	out := anyjson.Transform(context.Background(), v, func(v anyjson.Valuer, keyPath ...any) anyjson.Valuer {
		if s, ok := v.(*anyjson.String); ok && len(keyPath) > 0 {
			if key, ok := keyPath[len(keyPath)-1].(string); ok && key == "name" {
				return anyjson.StringOf(strings.ToUpper(strings.Trim(s.String(), `"`)))
			}
		}
		return v
	})

	fmt.Println(anyjson.ToString(out))
	// Output:
	// {"name":"OCTOHELM","tags":["go","json"]}
}
