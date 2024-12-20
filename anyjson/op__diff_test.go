package anyjson_test

import (
	"testing"

	. "github.com/octohelm/x/anyjson"
	testingx "github.com/octohelm/x/testing"
)

func TestDiff(t *testing.T) {
	t.Run("normal diff", func(t *testing.T) {
		base := MustFromValue(Obj{
			"int_changed":      1,
			"str_not_changed":  "string",
			"list_not_changed": List{"1", "2"},
			"list_changed":     List{"2"},
		}).(*Object)

		target := MustFromValue(Obj{
			"int_changed":      2,
			"str_not_changed":  "string",
			"list_not_changed": List{"1", "2"},
			"list_changed":     List{"1"},
		}).(*Object)

		diffed, err := Diff(base, target)
		testingx.Expect(t, err, testingx.Be[error](nil))

		testingx.Expect(t, diffed.Value(), testingx.Equal[any](Obj{
			"int_changed":  2,
			"list_changed": List{"1"},
		}))
	})

	t.Run("array object diff", func(t *testing.T) {
		t.Run("array object merge", func(t *testing.T) {
			base := MustFromValue(Obj{
				"withMergeKey": List{
					Obj{
						"name":  "a",
						"value": "x",
					},
					Obj{
						"name":  "b",
						"value": "x",
					},
				},
				"withoutMergeKey": List{
					Obj{
						"value": "1",
					},
				},
			}).(*Object)

			target := MustFromValue(Obj{
				"withMergeKey": List{
					Obj{
						"name":  "a",
						"value": "patched",
					},
					Obj{
						"name":  "c",
						"value": "new",
					},
				},
				"withoutMergeKey": List{
					Obj{
						"value": "2",
					},
				},
			}).(*Object)

			diffed, err := Diff(base, target)
			testingx.Expect(t, err, testingx.Be[error](nil))

			testingx.Expect(t, diffed.Value(), testingx.Equal[any](Obj{
				"withMergeKey": List{
					Obj{
						"name":  "a",
						"value": "patched",
					},
					Obj{
						"$patch": "delete",
						"name":   "b",
					},
					Obj{
						"name":  "c",
						"value": "new",
					},
				},
				"withoutMergeKey": List{
					Obj{
						"value": "2",
					},
				},
			}))
		})
	})
}
