package anyjson_test

import (
	"testing"

	"github.com/octohelm/x/anyjson"
	. "github.com/octohelm/x/testing/v2"
)

type (
	Object = anyjson.Object
	Obj    = anyjson.Obj
	List   = anyjson.List
)

func TestMerge(t *testing.T) {
	t.Run("normal merge", func(t *testing.T) {
		base := MustValue(t, func() (*Object, error) {
			v, err := anyjson.FromValue(Obj{
				"int":                     1,
				"str":                     "string",
				"arr":                     List{"1", "2"},
				"boolRemoveWhenNotExists": true,
			})
			return v.(*Object), err
		})

		patch := MustValue(t, func() (*Object, error) {
			v, err := anyjson.FromValue(Obj{
				"str":    "changed",
				"extra":  true,
				"ignore": nil,
				"deap": Obj{
					"a": Obj{
						"b": Obj{
							"ignore": nil,
						},
					},
				},
				"arr": List{"2", "3"},
			})
			return v.(*Object), err
		})

		merged := anyjson.Merge(base, patch)

		Then(t, "should deep merge and override values",
			Expect(merged.Value(), Equal[any](Obj{
				"int":   1,
				"str":   "changed",
				"extra": true,
				"deap": Obj{
					"a": Obj{
						"b": Obj{},
					},
				},
				"arr": List{"2", "3"},
			})),
		)
	})

	t.Run("nil as remover", func(t *testing.T) {
		base := MustValue(t, func() (*Object, error) {
			v, err := anyjson.FromValue(Obj{"int": 1, "str": "string"})
			return v.(*Object), err
		})

		patch := MustValue(t, func() (*Object, error) {
			v, err := anyjson.FromValue(Obj{"str": nil})
			return v.(*Object), err
		})

		merged := anyjson.Merge(base, patch, anyjson.WithNullOp(anyjson.NullAsRemover))

		Then(t, "str field should be removed",
			Expect(merged.Value(), Equal[any](Obj{"int": 1})),
		)
	})

	t.Run("array object merge", func(t *testing.T) {
		base := MustValue(t, func() (anyjson.Valuer, error) {
			return anyjson.FromValue(Obj{
				"withMergeKey": List{
					Obj{"name": "a", "value": "x"},
					Obj{"name": "b", "value": "x"},
					Obj{"name": "d", "value": "x"},
				},
				"withoutMergeKey": List{
					Obj{"value": "1"},
					Obj{"value": "2"},
				},
				"path": Obj{
					"to": Obj{
						"deep": Obj{
							"withoutMergeKey": List{
								Obj{"value": "1"},
							},
						},
					},
				},
			})
		})

		patch := MustValue(t, func() (anyjson.Valuer, error) {
			return anyjson.FromValue(Obj{
				"withMergeKey": List{
					Obj{"name": "a", "value": "patched"},
					Obj{"name": "c", "value": "new"},
					Obj{"name": "d", "$patch": anyjson.PatchOpDelete},
				},
				"withoutMergeKey": List{
					Obj{"value": "3"},
				},
				"path": Obj{
					"to": Obj{
						"deep": Obj{
							"withoutMergeKey": List{
								Obj{"value": "1"},
							},
						},
					},
				},
			})
		})

		merged := anyjson.Merge(base, patch, anyjson.WithArrayMergeKey("name"))

		Then(t, "should merge array by key or replace entirely",
			Expect(merged.Value(), Equal[any](Obj{
				"withMergeKey": List{
					Obj{"name": "a", "value": "patched"},
					Obj{"name": "b", "value": "x"},
					Obj{"name": "c", "value": "new"},
				},
				"withoutMergeKey": List{
					Obj{"value": "3"},
				},
				"path": Obj{
					"to": Obj{
						"deep": Obj{
							"withoutMergeKey": List{
								Obj{"value": "1"},
							},
						},
					},
				},
			})),
		)
	})
}
