package reflect_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/octohelm/x/ptr"
	testingx "github.com/octohelm/x/testing"

	. "github.com/octohelm/x/reflect"
)

func TestIndirect(t *testing.T) {
	testingx.Expect(t, reflect.ValueOf(1).Interface(), testingx.Equal(Indirect(reflect.ValueOf(ptr.Ptr(1))).Interface()))
	testingx.Expect(t, reflect.ValueOf(0).Interface(), testingx.Equal(Indirect(reflect.New(reflect.TypeOf(0))).Interface()))

	rv := New(reflect.PointerTo(reflect.PointerTo(reflect.PointerTo(reflect.TypeOf("")))))
	testingx.Expect(t, reflect.ValueOf("").Interface(), testingx.Equal(Indirect(rv).Interface()))
}

type Zero string

func (Zero) IsZero() bool {
	return true
}

func BenchmarkNew(b *testing.B) {
	tpe := reflect.PointerTo(reflect.TypeOf(Zero("")))

	for i := 0; i < b.N; i++ {
		_ = New(tpe)
	}
}

func BenchmarkIndirect(b *testing.B) {
	x := New(reflect.PointerTo(reflect.TypeOf(Zero(""))))

	for i := 0; i < b.N; i++ {
		_ = Indirect(x)
	}
}

func TestNew(t *testing.T) {
	t.Run("NewType", func(t *testing.T) {
		tpe := reflect.TypeOf(Zero(""))
		z, ok := New(tpe).Interface().(Zero)
		testingx.Expect(t, ok, testingx.BeTrue())
		testingx.Expect(t, z, testingx.Equal(Zero("")))
	})

	t.Run("NewPtrType", func(t *testing.T) {
		tpe := reflect.PointerTo(reflect.TypeOf(Zero("")))
		z, ok := New(tpe).Interface().(*Zero)
		testingx.Expect(t, ok, testingx.BeTrue())
		testingx.Expect(t, *z, testingx.Equal(Zero("")))
	})

	t.Run("NewPtrPtrType", func(t *testing.T) {
		tpe := reflect.PointerTo(reflect.PointerTo(reflect.TypeOf(Zero(""))))
		z, ok := New(tpe).Interface().(**Zero)
		testingx.Expect(t, ok, testingx.BeTrue())
		testingx.Expect(t, **z, testingx.Equal(Zero("")))
	})
}

type S struct {
	V any
}

var emptyValues = []any{
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

var nonEmptyValues = []any{
	Zero("11111111111"),
	ptr.Ptr("12322"),
}

func BenchmarkIsEmptyValue(b *testing.B) {
	for i, v := range append(emptyValues, nonEmptyValues...) {
		b.Run(fmt.Sprintf("%d: %#v", i, v), func(b *testing.B) {
			IsEmptyValue(v)
		})

		if _, ok := v.(reflect.Value); !ok {
			rv := reflect.ValueOf(v)
			b.Run(fmt.Sprintf("%d: reflect.Value(%#v)", i, v), func(b *testing.B) {
				IsEmptyValue(rv)
			})
		}
	}
}

func TestIsEmptyValue(t *testing.T) {
	for i, v := range emptyValues {
		t.Run(fmt.Sprintf("%d: %#v", i, v), func(t *testing.T) {
			testingx.Expect(t, IsEmptyValue(v), testingx.BeTrue())
		})

		if _, ok := v.(reflect.Value); !ok {
			rv := reflect.ValueOf(v)

			t.Run(fmt.Sprintf("%d: reflect.Value(%#v)", i, v), func(t *testing.T) {
				testingx.Expect(t, IsEmptyValue(rv), testingx.BeTrue())
			})
		}
	}
}
