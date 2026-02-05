package anyjson_test

import (
	"testing"

	"github.com/octohelm/x/anyjson"
	. "github.com/octohelm/x/testing/v2"
)

func TestDiff(t *testing.T) {
	t.Run("normal diff", func(t *testing.T) {
		base := MustValue(t, func() (*Object, error) {
			v, err := anyjson.FromValue(Obj{
				"int_changed":                1,
				"str_not_changed":            "string",
				"list_not_changed":           List{"1", "2"},
				"list_changed":               List{"2"},
				"bool_not_exists_as_changed": true,
			})
			return v.(*Object), err
		})

		target := MustValue(t, func() (*Object, error) {
			v, err := anyjson.FromValue(Obj{
				"int_changed":      2,
				"str_not_changed":  "string",
				"list_not_changed": List{"1", "2"},
				"list_changed":     List{"1"},
			})
			return v.(*Object), err
		})

		diffed := MustValue(t, func() (anyjson.Valuer, error) {
			return anyjson.Diff(base, target)
		})

		Then(t, "diff should contain only changed fields",
			Expect(diffed.Value(), Equal[any](Obj{
				"int_changed":                2,
				"list_changed":               List{"1"},
				"bool_not_exists_as_changed": false,
			})),
		)
	})

	t.Run("array object diff", func(t *testing.T) {
		base := MustValue(t, func() (*Object, error) {
			v, err := anyjson.FromValue(Obj{
				"withMergeKey": List{
					Obj{"name": "a", "value": "x"},
					Obj{"name": "b", "value": "x"},
				},
				"withoutMergeKey": List{
					Obj{"value": "1"},
				},
			})
			return v.(*Object), err
		})

		target := MustValue(t, func() (*Object, error) {
			v, err := anyjson.FromValue(Obj{
				"withMergeKey": List{
					Obj{"name": "a", "value": "patched"},
					Obj{"name": "c", "value": "new"},
				},
				"withoutMergeKey": List{
					Obj{"value": "2"},
				},
			})
			return v.(*Object), err
		})

		diffed := MustValue(t, func() (anyjson.Valuer, error) {
			return anyjson.Diff(base, target)
		})

		Then(t, "should generate patch with delete and add operations",
			Expect(diffed.Value(), Equal[any](Obj{
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
			})),
		)
	})
}
