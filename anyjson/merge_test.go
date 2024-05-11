package anyjson

import (
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestMerge(t *testing.T) {
	base := MustFromValue(Obj{
		"int": 1,
		"str": "string",
	}).(*Object)

	t.Run("normal merge", func(t *testing.T) {
		patch := MustFromValue(Obj{
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
		}).(*Object)

		merged := Merge(base, patch)

		testingx.Expect(t, merged.Value(), testingx.Equal[any](Obj{
			"int":   1,
			"str":   "changed",
			"extra": true,
			"deap": Obj{
				"a": Obj{
					"b": Obj{},
				},
			},
		}))
	})

	t.Run("nil remover", func(t *testing.T) {
		patch := MustFromValue(Obj{
			"str": nil,
		}).(*Object)

		merged := Merge(base, patch, WithNullOp(NullAsRemover))

		testingx.Expect(t, merged.Value(), testingx.Equal[any](Obj{
			"int": 1,
		}))
	})

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
		})

		patch := MustFromValue(List{
			Obj{
				"name":  "a",
				"value": "patched",
			},
			Obj{
				"name":  "c",
				"value": "new",
			},
		})

		merged := Merge(base, patch, WithArrayMergeKey("name"))

		testingx.Expect(t, merged.Value(), testingx.Equal[any](List{
			Obj{
				"name":  "a",
				"value": "patched",
			},
			Obj{
				"name":  "b",
				"value": "x",
			},
			Obj{
				"name":  "c",
				"value": "new",
			},
		}))
	})
}
