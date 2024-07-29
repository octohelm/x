package anyjson

import (
	"testing"

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
			base := MustFromValue(List{
				Obj{
					"name":  "a",
					"value": "x",
				},
				Obj{
					"name":  "b",
					"value": "x",
				},
			}).(*Array)

			target := MustFromValue(List{
				Obj{
					"name":  "a",
					"value": "patched",
				},
				Obj{
					"name":  "c",
					"value": "new",
				},
			}).(*Array)

			diffed, err := Diff(base, target)
			testingx.Expect(t, err, testingx.Be[error](nil))

			testingx.Expect(t, diffed.Value(), testingx.Equal[any](List{
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
			}))
		})
	})
}
