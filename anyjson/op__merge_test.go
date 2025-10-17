package anyjson_test

import (
	"testing"

	. "github.com/octohelm/x/anyjson"
	testingx "github.com/octohelm/x/testing"
)

func TestMerge(t *testing.T) {
	t.Run("normal merge", func(t *testing.T) {
		base := MustFromValue(Obj{
			"int": 1,
			"str": "string",
			"arr": List{
				"1", "2",
			},
			"boolRemoveWhenNotExists": true,
		}).(*Object)

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
			"arr": List{
				"2", "3",
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
			"arr": List{
				"2", "3",
			},
		}))
	})

	t.Run("nil as remover", func(t *testing.T) {
		base := MustFromValue(Obj{
			"int": 1,
			"str": "string",
		}).(*Object)

		patch := MustFromValue(Obj{
			"str": nil,
		}).(*Object)

		merged := Merge(base, patch, WithNullOp(NullAsRemover))

		testingx.Expect(t, merged.Value(), testingx.Equal[any](Obj{
			"int": 1,
		}))
	})

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
				Obj{
					"name":  "d",
					"value": "x",
				},
			},
			"withoutMergeKey": List{
				Obj{
					"value": "1",
				},
				Obj{
					"value": "2",
				},
			},
			"path": Obj{
				"to": Obj{
					"deep": Obj{
						"withoutMergeKey": List{
							Obj{
								"value": "1",
							},
						},
					},
				},
			},
		})

		patch := MustFromValue(Obj{
			"withMergeKey": List{
				Obj{
					"name":  "a",
					"value": "patched",
				},
				Obj{
					"name":  "c",
					"value": "new",
				},
				Obj{
					"name":   "d",
					"$patch": PatchOpDelete,
				},
			},
			"withoutMergeKey": List{
				Obj{
					"value": "3",
				},
			},
			"path": Obj{
				"to": Obj{
					"deep": Obj{
						"withoutMergeKey": List{
							Obj{
								"value": "1",
							},
						},
					},
				},
			},
		})

		merged := Merge(base, patch, WithArrayMergeKey("name"))

		testingx.Expect(t, merged.Value(), testingx.Equal[any](Obj{
			"withMergeKey": List{
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
			},
			"withoutMergeKey": List{
				Obj{
					"value": "3",
				},
			},
			"path": Obj{
				"to": Obj{
					"deep": Obj{
						"withoutMergeKey": List{
							Obj{
								"value": "1",
							},
						},
					},
				},
			},
		}))
	})
}
