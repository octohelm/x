package anyjson_test

import (
	"fmt"

	"github.com/octohelm/x/anyjson"
)

func ExampleFromValue() {
	v, err := anyjson.FromValue(map[string]any{
		"name": "octohelm",
		"tags": []any{"go", "json"},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(anyjson.ToString(v))
	// Output:
	// {"name":"octohelm","tags":["go","json"]}
}

func ExampleMerge() {
	base := anyjson.MustFromValue(anyjson.Obj{
		"name":  "octohelm",
		"count": 1,
	})
	patch := anyjson.MustFromValue(anyjson.Obj{
		"count": 2,
		"lang":  "go",
	})

	merged := anyjson.Merge(base, patch)
	fmt.Println(anyjson.ToString(merged))
	// Output:
	// {"count":2,"name":"octohelm","lang":"go"}
}

func ExampleDiff() {
	template := map[string]any{
		"name":  "octohelm",
		"count": 1,
	}
	live := map[string]any{
		"name":  "octohelm",
		"count": 2,
	}

	diff, err := anyjson.Diff(&template, &live)
	if err != nil {
		panic(err)
	}

	fmt.Println(anyjson.ToString(diff))
	// Output:
	// {"count":2}
}
