package types_test

import (
	"encoding"
	"fmt"
	"go/types"
	"reflect"
	"testing"
	"unsafe"

	"github.com/octohelm/x/cmp"
	"github.com/octohelm/x/ptr"
	. "github.com/octohelm/x/testing/v2"
	. "github.com/octohelm/x/types"
	"github.com/octohelm/x/types/testdata/typ"
	typ2 "github.com/octohelm/x/types/testdata/typ/typ"
)

func TestType(t *testing.T) {
	fn := func(a, b string) bool { return true }

	values := []any{
		typ.AnyStruct[string]{Name: "x"},
		typ.AnySlice[string]{},
		typ.AnyMap[int, string]{},
		typ.AnyMap[int, fmt.Stringer]{},
		typ.IntMap{},
		typ.DeepCompose{},
		func() *typ.Enum { v := typ.ENUM__ONE; return &v }(),
		typ.ENUM__ONE,
		reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem(),
		reflect.TypeOf((*typ.SomeMixInterface)(nil)).Elem(),
		unsafe.Pointer(t),
		make(typ.Chan),
		make(chan string, 100),
		typ.F,
		typ.Func(fn),
		fn,
		typ.String(""), "",
		typ.Bool(true), true,
		typ.Int(0), ptr.Ptr(1), int(0),
		typ.Int8(0), int8(0),
		typ.Int16(0), int16(0),
		typ.Int32(0), int32(0),
		typ.Int64(0), int64(0),
		typ.Uint(0), uint(0),
		typ.Uintptr(0), uintptr(0),
		typ.Uint8(0), uint8(0),
		typ.Uint16(0), uint16(0),
		typ.Uint32(0), uint32(0),
		typ.Uint64(0), uint64(0),
		typ.Float32(0), float32(0),
		typ.Float64(0), float64(0),
		typ.Complex64(0), complex64(0),
		typ.Complex128(0), complex128(0),
		typ.Array{},
		[1]string{},
		typ.Slice{},
		[]string{},
		typ.Map{},
		map[string]string{},
		typ.Struct{},
		struct{}{},
		struct {
			typ.Resource[typ.Int]
			typ.Part
			Part2  typ2.Part
			a      string
			A      string `json:"a"`
			Struct struct{ B string }
		}{},
	}

	for _, v := range values {
		check(t, v)
	}
}

func check(t *testing.T, v any) {
	rtype, ok := v.(reflect.Type)
	if !ok {
		rtype = reflect.TypeOf(v)
	}

	ttype := NewTypesTypeFromReflectType(rtype)
	rt := FromRType(rtype)
	tt := FromTType(ttype)

	t.Run(FullTypeName(rt), func(t *testing.T) {
		Then(t, "basic metadata should match",
			Expect(rt.String(), Equal(tt.String())),
			Expect(rt.Kind().String(), Equal(tt.Kind().String())),
			Expect(rt.Name(), Equal(tt.Name())),
			Expect(rt.PkgPath(), Equal(tt.PkgPath())),
			Expect(rt.Comparable(), Equal(tt.Comparable())),
		)

		Then(t, "assignability and convertibility should match",
			Expect(
				rt.AssignableTo(FromRType(reflect.TypeOf(""))),
				Equal(tt.AssignableTo(FromTType(types.Typ[types.String]))),
			),
			Expect(
				rt.ConvertibleTo(FromRType(reflect.TypeOf(""))),
				Equal(tt.ConvertibleTo(FromTType(types.Typ[types.String]))),
			),
		)

		t.Run("Methods", func(t *testing.T) {
			Then(t, "method count should match",
				Expect(rt.NumMethod(), Equal(tt.NumMethod())),
			)

			for i := 0; i < rt.NumMethod(); i++ {
				rM := rt.Method(i)
				tM, ok := tt.MethodByName(rM.Name())

				t.Run("M "+rM.Name(), func(t *testing.T) {
					Then(t, "method details should match",
						Expect(ok, Be(cmp.True())),
						Expect(rM.Name(), Equal(tM.Name())),
						Expect(rM.PkgPath(), Equal(tM.PkgPath())),
						Expect(rM.Type().String(), Equal(tM.Type().String())),
					)
				})
			}

			Then(t, "String() method presence should match",
				Expect(
					func() bool { _, ok := rt.MethodByName("String"); return ok }(),
					Equal(func() bool { _, ok := tt.MethodByName("String"); return ok }()),
				),
			)
		})

		if rt.Kind() == reflect.Struct {
			t.Run("Fields", func(t *testing.T) {
				Then(t, "field count should match",
					Expect(rt.NumField(), Equal(tt.NumField())),
				)

				for i := 0; i < rt.NumField(); i++ {
					rsf := rt.Field(i)
					tsf := tt.Field(i)

					t.Run("F "+rsf.Name(), func(t *testing.T) {
						Then(t, "field metadata should match",
							Expect(rsf.Anonymous(), Equal(tsf.Anonymous())),
							Expect(rsf.Tag(), Equal(tsf.Tag())),
							Expect(rsf.Name(), Equal(tsf.Name())),
							Expect(FullTypeName(rsf.Type()), Equal(FullTypeName(tsf.Type()))),
						)
					})
				}
			})
		}

		if rt.Kind() == reflect.Func {
			Then(t, "function signature should match",
				Expect(rt.NumIn(), Equal(tt.NumIn())),
				Expect(rt.NumOut(), Equal(tt.NumOut())),
			)
		}
	})
}

func TestTryNew(t *testing.T) {
	t.Run("TryNew behavior", func(t *testing.T) {
		Then(t, "RType should support TryNew",
			Expect(
				func() bool { _, ok := TryNew(FromRType(reflect.TypeOf(typ.Struct{}))); return ok }(),
				Be(cmp.True()),
			),
		)
		Then(t, "TType should not support TryNew",
			Expect(
				func() bool {
					_, ok := TryNew(FromTType(NewTypesTypeFromReflectType(reflect.TypeOf(typ.Struct{}))))
					return ok
				}(),
				Be(cmp.False()),
			),
		)
	})
}

func TestEachField(t *testing.T) {
	expect := []string{"a", "b", "bool", "c", "Part2"}

	t.Run("EachField parity", func(t *testing.T) {
		collect := func(tpe Type) []string {
			names := make([]string, 0)
			EachField(tpe, "json", func(field StructField, fieldDisplayName string, omitempty bool) bool {
				names = append(names, fieldDisplayName)
				return true
			})
			return names
		}

		Then(t, "RType and TType field names should match",
			Expect(collect(FromRType(reflect.TypeOf(typ.Struct{}))), Equal(expect)),
			Expect(collect(FromTType(NewTypesTypeFromReflectType(reflect.TypeOf(typ.Struct{})))), Equal(expect)),
		)
	})
}
