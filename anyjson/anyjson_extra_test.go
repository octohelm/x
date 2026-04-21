package anyjson_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/octohelm/x/anyjson"
	. "github.com/octohelm/x/testing/v2"
)

func TestObjectAndArrayHelpers(t *testing.T) {
	t.Run("Object 保持插入顺序并支持更新删除", func(t *testing.T) {
		obj := &anyjson.Object{}

		firstInsert := obj.Set("b", anyjson.StringOf("first"))
		secondInsert := obj.Set("a", anyjson.NumberOf(1))
		updateInsert := obj.Set("b", anyjson.StringOf("updated"))
		deleted := obj.Delete("a")
		missingDeleted := obj.Delete("missing")

		keys := make([]string, 0)
		values := make([]any, 0)
		for k, v := range obj.KeyValues() {
			keys = append(keys, k)
			values = append(values, v.Value())
		}

		got, ok := obj.Get("b")
		_, missingOK := obj.Get("a")

		Then(t, "Object 应保持更新后的顺序与值",
			Expect(firstInsert, Equal(true)),
			Expect(secondInsert, Equal(true)),
			Expect(updateInsert, Equal(false)),
			Expect(deleted, Equal(true)),
			Expect(missingDeleted, Equal(false)),
			Expect(obj.Len(), Equal(1)),
			Expect(keys, Equal([]string{"b"})),
			Expect(values, Equal([]any{"updated"})),
			Expect(ok, Equal(true)),
			Expect(got.Value(), Equal(any("updated"))),
			Expect(missingOK, Equal(false)),
		)
	})

	t.Run("Array 支持索引、遍历和越界检查", func(t *testing.T) {
		arr := &anyjson.Array{}
		arr.Append(anyjson.StringOf("x"))
		arr.Append(anyjson.NumberOf(2))

		values := make([]any, 0)
		for v := range arr.Values() {
			values = append(values, v.Value())
		}

		indexed := make([]string, 0)
		for i, v := range arr.IndexedValues() {
			indexed = append(indexed, fmt.Sprintf("%d=%v", i, v.Value()))
		}

		item, ok := arr.Index(1)
		_, missingOK := arr.Index(3)

		Then(t, "Array 应保留顺序并正确处理越界",
			Expect(arr.Len(), Equal(2)),
			Expect(values, Equal([]any{"x", 2})),
			Expect(indexed, Equal([]string{"0=x", "1=2"})),
			Expect(ok, Equal(true)),
			Expect(item.Value(), Equal(any(2))),
			Expect(missingOK, Equal(false)),
		)
	})
}

func TestAsAndMustFromValue(t *testing.T) {
	t.Run("As 将 Valuer 反序列化到目标结构", func(t *testing.T) {
		type employee struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		valuer := MustValue(t, func() (anyjson.Valuer, error) {
			return anyjson.FromValue(map[string]any{
				"name": "octo",
				"age":  18,
			})
		})

		var out employee

		Then(t, "As 应先完成解码",
			ExpectDo(func() error {
				return anyjson.As(valuer, &out)
			}),
		)

		Then(t, "As 应成功写入目标结构",
			Expect(out, Equal(employee{Name: "octo", Age: 18})),
		)
	})

	t.Run("FromValue 遇到不支持的值返回错误，MustFromValue 会 panic", func(t *testing.T) {
		type unsupported struct {
			Ch chan int `json:"ch"`
		}

		_, err := anyjson.FromValue(unsupported{Ch: make(chan int)})

		Then(t, "FromValue 应返回错误",
			Expect(err, Be(func(actual error) error {
				if actual == nil {
					return &ErrNotEqual{Expect: "non-nil error", Got: actual}
				}
				return nil
			})),
		)

		panicMsg := capturePanic(func() {
			_ = anyjson.MustFromValue(unsupported{Ch: make(chan int)})
		})

		Then(t, "MustFromValue 应传播 panic",
			Expect(panicMsg, Be(func(actual string) error {
				if actual == "" {
					return &ErrNotEqual{Expect: "panic message", Got: actual}
				}
				return nil
			})),
		)
	})
}

func TestSortedAndTransform(t *testing.T) {
	t.Run("Sorted 递归稳定对象键顺序且保留数组顺序", func(t *testing.T) {
		obj := &anyjson.Object{}
		obj.Set("b", anyjson.StringOf("2"))
		obj.Set("a", anyjson.StringOf("1"))

		child := &anyjson.Object{}
		child.Set("d", anyjson.StringOf("4"))
		child.Set("c", anyjson.StringOf("3"))

		arr := &anyjson.Array{}
		arr.Append(child)
		arr.Append(anyjson.StringOf("tail"))
		obj.Set("arr", arr)

		sorted := anyjson.Sorted(obj)

		Then(t, "Sorted 应按键排序对象，但不改变数组元素顺序",
			Expect(sorted.String(), Equal(`{"a":"1","arr":[{"c":"3","d":"4"},"tail"],"b":"2"}`)),
		)
	})

	t.Run("Transform 仅对叶子值调用并支持删除节点", func(t *testing.T) {
		source := MustValue(t, func() (anyjson.Valuer, error) {
			return anyjson.FromValue(map[string]any{
				"name": "octo",
				"items": []any{
					"keep",
					"drop",
				},
				"meta": map[string]any{
					"enabled": true,
				},
			})
		})

		visited := make([]string, 0)

		transformed := anyjson.Transform(nil, source, func(v anyjson.Valuer, keyPath ...any) anyjson.Valuer {
			visited = append(visited, fmt.Sprint(keyPath...))

			if s, ok := v.Value().(string); ok {
				if s == "drop" {
					return nil
				}
				return anyjson.StringOf("prefix-" + s)
			}
			return v
		})

		Then(t, "Transform 应记录叶子路径并删除返回 nil 的节点",
			Expect(visited, Equal([]string{"items0", "items1", "metaenabled", "name"})),
			Expect(transformed.Value(), Equal[any](anyjson.Obj{
				"name": "prefix-octo",
				"items": anyjson.List{
					"prefix-keep",
				},
				"meta": anyjson.Obj{
					"enabled": true,
				},
			})),
		)
	})
}

func TestPatchHelpers(t *testing.T) {
	t.Run("IsPatchObject 判断补丁对象类型", func(t *testing.T) {
		valid := &anyjson.Object{}
		valid.Set(anyjson.PatchKey, anyjson.StringOf(string(anyjson.PatchOpDelete)))

		invalid := &anyjson.Object{}
		invalid.Set(anyjson.PatchKey, anyjson.StringOf("unknown"))

		op, ok := anyjson.IsPatchObject(valid)
		invalidOp, invalidOK := anyjson.IsPatchObject(invalid)
		nilOp, nilOK := anyjson.IsPatchObject(nil)

		Then(t, "应只识别支持的补丁操作",
			Expect(op, Equal(anyjson.PatchOp("delete"))),
			Expect(ok, Equal(true)),
			Expect(invalidOp, Equal(anyjson.PatchOp(""))),
			Expect(invalidOK, Equal(false)),
			Expect(nilOp, Equal(anyjson.PatchOp(""))),
			Expect(nilOK, Equal(false)),
		)
	})
}

func TestScalarHelpersAndDecoders(t *testing.T) {
	t.Run("标量值可返回原生值与 JSON 文本", func(t *testing.T) {
		s := anyjson.StringOf("octo")
		n := anyjson.NumberOf(12)
		b := anyjson.BooleanOf(true)
		null := &anyjson.Null{}

		Then(t, "String 方法应返回 JSON 文本表示",
			Expect(s.Value(), Equal(any("octo"))),
			Expect(s.String(), Equal(`"octo"`)),
			Expect(n.Value(), Equal(any(12))),
			Expect(n.String(), Equal(`12`)),
			Expect(b.Value(), Equal(any(true))),
			Expect(b.String(), Equal(`true`)),
			Expect(null.Value(), Equal(any(nil))),
			Expect(null.String(), Equal(`null`)),
		)
	})

	t.Run("从 JSON decoder 读取各类标量", func(t *testing.T) {
		for _, c := range []struct {
			name string
			raw  string
			want any
		}{
			{name: "string", raw: `"octo"`, want: "octo"},
			{name: "bool-true", raw: `true`, want: true},
			{name: "bool-false", raw: `false`, want: false},
			{name: "number", raw: `12`, want: 12},
			{name: "null", raw: `null`, want: nil},
		} {
			t.Run(c.name, func(t *testing.T) {
				v := MustValue(t, func() (anyjson.Valuer, error) {
					return anyjson.FromJSONTextDecoder(jsontext.NewDecoder(bytes.NewBufferString(c.raw)))
				})

				Then(t, "decoder 应返回期望值",
					Expect(v.Value(), Equal(c.want)),
				)
			})
		}
	})

	t.Run("数组、对象与 null 的 JSON wrapper", func(t *testing.T) {
		var arr anyjson.Array
		var obj anyjson.Object
		var null anyjson.Null

		Then(t, "UnmarshalJSON wrapper 应成功解码",
			ExpectDo(func() error { return arr.UnmarshalJSON([]byte(`[1,"x"]`)) }),
			ExpectDo(func() error { return obj.UnmarshalJSON([]byte(`{"a":1,"b":"x"}`)) }),
			ExpectDo(func() error { return null.UnmarshalJSON([]byte(`null`)) }),
		)

		rawNull := MustValue(t, null.MarshalJSON)

		Then(t, "wrapper 解码后的值应符合预期",
			Expect(arr.String(), Equal(`[1,"x"]`)),
			Expect(obj.String(), Equal(`{"a":1,"b":"x"}`)),
			Expect(string(rawNull), Equal(`null`)),
		)
	})
}

func TestEmptyObjectAsNullOptions(t *testing.T) {
	t.Run("Merge 支持将空对象视为 null", func(t *testing.T) {
		base := MustValue(t, func() (anyjson.Valuer, error) {
			return anyjson.FromValue(anyjson.Obj{})
		})
		patch := MustValue(t, func() (anyjson.Valuer, error) {
			return anyjson.FromValue(anyjson.Obj{})
		})

		merged := anyjson.Merge(base, patch, anyjson.WithEmptyObjectAsNull())

		Then(t, "空对象 merge 后可转为 null",
			Expect(merged.Value(), Equal(any(nil))),
		)
	})

	t.Run("Diff 支持将空对象视为 null", func(t *testing.T) {
		template := anyjson.Obj{}
		live := anyjson.Obj{}

		diffed := MustValue(t, func() (anyjson.Valuer, error) {
			return anyjson.Diff(&template, &live, anyjson.WithEmptyObjectAsNull())
		})

		Then(t, "空对象 diff 后可转为 null",
			Expect(diffed.Value(), Equal(any(nil))),
		)
	})
}

func capturePanic(do func()) (msg string) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Sprint(r)
		}
	}()
	do()
	return ""
}
