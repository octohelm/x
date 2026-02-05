package reflect_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/octohelm/x/cmp"
	"github.com/octohelm/x/ptr"
	reflectx "github.com/octohelm/x/reflect"
	. "github.com/octohelm/x/testing/v2"
)

func TestIndirect(t *testing.T) {
	t.Run("GIVEN various pointer levels", func(t *testing.T) {
		Then(t, "Indirect should always return the underlying value",
			Expect(
				reflectx.Indirect(reflect.ValueOf(ptr.Ptr(1))).Interface(),
				Equal(reflect.ValueOf(1).Interface()),
			),
			Expect(
				reflectx.Indirect(reflect.New(reflect.TypeOf(0))).Interface(),
				Equal(reflect.ValueOf(0).Interface()),
			),
		)

		t.Run("WHEN deep nested pointers", func(t *testing.T) {
			rv := reflectx.New(reflect.PointerTo(reflect.PointerTo(reflect.PointerTo(reflect.TypeOf("")))))

			Then(t, "it should unwrap to the base value",
				Expect(reflectx.Indirect(rv).Interface(), Equal(reflect.ValueOf("").Interface())),
			)
		})
	})
}

type Zero string

func (Zero) IsZero() bool {
	return true
}

func TestNew(t *testing.T) {
	tpeZero := reflect.TypeOf(Zero(""))

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
