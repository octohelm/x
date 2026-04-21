package reflect_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/octohelm/x/cmp"
	reflectx "github.com/octohelm/x/reflect"
	. "github.com/octohelm/x/testing/v2"
)

func TestIndirect(t *testing.T) {
	t.Run("GIVEN various pointer levels", func(t *testing.T) {
		Then(t, "Indirect should always return the underlying value",
			Expect(
				reflectx.Indirect(reflect.ValueOf(new(1))).Interface(),
				Equal(reflect.ValueOf(1).Interface()),
			),
			Expect(
				reflectx.Indirect(reflect.New(reflect.TypeFor[int]())).Interface(),
				Equal(reflect.ValueOf(0).Interface()),
			),
		)

		t.Run("WHEN deep nested pointers", func(t *testing.T) {
			rv := reflectx.New(reflect.PointerTo(reflect.PointerTo(reflect.PointerTo(reflect.TypeFor[string]()))))

			Then(t, "it should unwrap to the base value",
				Expect(reflectx.Indirect(rv).Interface(), Equal(reflect.ValueOf("").Interface())),
			)
		})
	})

	t.Run("WHEN the pointer is nil", func(t *testing.T) {
		rv := reflect.Zero(reflect.TypeFor[*string]())

		Then(t, "it should return an invalid value after dereference",
			Expect(reflectx.Indirect(rv).IsValid(), Be(cmp.False())),
		)
	})
}

type Zero string

func (Zero) IsZero() bool {
	return true
}

func TestNew(t *testing.T) {
	tpeZero := reflect.TypeFor[Zero]()

	t.Run("NewType", func(t *testing.T) {
		z := MustValue(t, func() (Zero, error) {
			v, ok := reflectx.New(tpeZero).Interface().(Zero)
			if !ok {
				return "", fmt.Errorf("not Zero type")
			}
			return v, nil
		})

		Then(t, "should be initialized zero value",
			Expect(z, Equal(Zero(""))),
		)
	})

	t.Run("NewPtrType", func(t *testing.T) {
		tpe := reflect.PointerTo(tpeZero)

		Then(t, "should create a pointer to zero value",
			Expect(
				func() *Zero { v, _ := reflectx.New(tpe).Interface().(*Zero); return v }(),
				Be(cmp.NotNil[*Zero]()),
			),
			Expect(
				func() Zero { v, _ := reflectx.New(tpe).Interface().(*Zero); return *v }(),
				Equal(Zero("")),
			),
		)
	})

	t.Run("NewPtrPtrType", func(t *testing.T) {
		tpe := reflect.PointerTo(reflect.PointerTo(tpeZero))

		Then(t, "should create nested pointers to zero value",
			Expect(
				func() Zero { v, _ := reflectx.New(tpe).Interface().(**Zero); return **v }(),
				Equal(Zero("")),
			),
		)
	})
}

type S struct {
	V any
}

func TestIsEmptyValue(t *testing.T) {
	emptyValues := []any{
		Zero(""),
		(*string)(nil),
		(any)(nil),
		(S{}).V,
		"",
		0,
		uint(0),
		float32(0),
		false,
		reflect.ValueOf(S{}).FieldByName("V"),
		nil,
	}

	for i, v := range emptyValues {
		t.Run(fmt.Sprintf("EmptyCase/%d", i), func(t *testing.T) {
			Then(t, fmt.Sprintf("value [%#v] should be empty", v),
				Expect(reflectx.IsEmptyValue(v), Be(cmp.True())),
			)

			if _, ok := v.(reflect.Value); !ok {
				rv := reflect.ValueOf(v)
				Then(t, "wrapped reflect.Value should also be empty",
					Expect(reflectx.IsEmptyValue(rv), Be(cmp.True())),
				)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	t.Run("反射容器零值判定", func(t *testing.T) {
		Then(t, "map、slice、string 和 invalid 值应按长度或有效性判断",
			Expect(reflectx.IsZero(reflect.ValueOf(map[string]int{})), Be(cmp.True())),
			Expect(reflectx.IsZero(reflect.ValueOf([]int{})), Be(cmp.True())),
			Expect(reflectx.IsZero(reflect.ValueOf("")), Be(cmp.True())),
			Expect(reflectx.IsZero(reflect.Value{}), Be(cmp.True())),
		)
	})

	t.Run("接口包装值会继续判断其底层元素", func(t *testing.T) {
		var zero any = ""
		var nonZero any = "x"

		Then(t, "接口中的零值和非零值应区分开",
			Expect(reflectx.IsZero(reflect.ValueOf(&zero).Elem()), Be(cmp.True())),
			Expect(reflectx.IsZero(reflect.ValueOf(&nonZero).Elem()), Be(cmp.False())),
		)
	})

	t.Run("nil map 与 nil slice 应视为零值", func(t *testing.T) {
		var m map[string]int
		var s []int

		Then(t, "nil 容器应被视为零值",
			Expect(reflectx.IsZero(m), Be(cmp.True())),
			Expect(reflectx.IsZero(s), Be(cmp.True())),
		)
	})

	t.Run("非零值应返回 false", func(t *testing.T) {
		Then(t, "非空集合和 true 应被识别为非零值",
			Expect(reflectx.IsZero([]int{1}), Be(cmp.False())),
			Expect(reflectx.IsZero(map[string]int{"x": 1}), Be(cmp.False())),
			Expect(reflectx.IsZero(true), Be(cmp.False())),
		)
	})
}
