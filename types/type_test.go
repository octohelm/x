package types_test

import (
	"encoding"
	"fmt"
	"go/types"
	"reflect"
	"testing"
	"unsafe"

	"github.com/octohelm/x/ptr"
	testingx "github.com/octohelm/x/testing"
	. "github.com/octohelm/x/types"
	"github.com/octohelm/x/types/testdata/typ"
	typ2 "github.com/octohelm/x/types/testdata/typ/typ"
)

func TestType(t *testing.T) {
	fn := func(a, b string) bool {
		return true
	}

	values := []any{
		typ.AnyStruct[string]{
			Name: "x",
		},
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

		typ.String(""),
		"",
		typ.Bool(true),
		true,
		typ.Int(0),
		ptr.Ptr(1),
		int(0),
		typ.Int8(0),
		int8(0),
		typ.Int16(0),
		int16(0),
		typ.Int32(0),
		int32(0),
		typ.Int64(0),
		int64(0),
		typ.Uint(0),
		uint(0),
		typ.Uintptr(0),
		uintptr(0),
		typ.Uint8(0),
		uint8(0),
		typ.Uint16(0),
		uint16(0),
		typ.Uint32(0),
		uint32(0),
		typ.Uint64(0),
		uint64(0),
		typ.Float32(0),
		float32(0),
		typ.Float64(0),
		float64(0),
		typ.Complex64(0),
		complex64(0),
		typ.Complex128(0),
		complex128(0),
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
			Struct struct {
				B string
			}
		}{},
	}

	for i := range values {
		check(t, values[i])
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
		testingx.Expect(t, rt.String(), testingx.Equal(tt.String()))
		testingx.Expect(t, rt.Kind().String(), testingx.Equal(tt.Kind().String()))
		testingx.Expect(t, rt.Name(), testingx.Equal(tt.Name()))
		testingx.Expect(t, rt.PkgPath(), testingx.Equal(tt.PkgPath()))
		testingx.Expect(t, rt.Comparable(), testingx.Equal(tt.Comparable()))
		testingx.Expect(t, rt.AssignableTo(FromRType(reflect.TypeOf(""))), testingx.Equal(tt.AssignableTo(FromTType(types.Typ[types.String]))))
		testingx.Expect(t, rt.ConvertibleTo(FromRType(reflect.TypeOf(""))), testingx.Equal(tt.ConvertibleTo(FromTType(types.Typ[types.String]))))

		testingx.Expect(t, rt.NumMethod(), testingx.Equal(tt.NumMethod()))

		for i := 0; i < rt.NumMethod(); i++ {
			rMethod := rt.Method(i)
			tMethod, ok := tt.MethodByName(rMethod.Name())

			t.Run("M "+rMethod.Name()+" | "+rMethod.Type().String(), func(t *testing.T) {
				testingx.Expect(t, ok, testingx.BeTrue())
				testingx.Expect(t, rMethod.Name(), testingx.Equal(tMethod.Name()))
				testingx.Expect(t, rMethod.PkgPath(), testingx.Equal(tMethod.PkgPath()))
				testingx.Expect(t, rMethod.Type().String(), testingx.Equal(tMethod.Type().String()))
			})
		}

		{
			_, rOk := rt.MethodByName("String")
			_, tOk := tt.MethodByName("String")
			testingx.Expect(t, rOk, testingx.Equal(tOk))
		}

		{
			rReplacer, rIs := EncodingTextMarshalerTypeReplacer(rt)
			tReplacer, tIs := EncodingTextMarshalerTypeReplacer(tt)
			testingx.Expect(t, rIs, testingx.Equal(tIs))
			testingx.Expect(t, rReplacer.String(), testingx.Equal(tReplacer.String()))
		}

		if rt.Kind() == reflect.Array {
			testingx.Expect(t, rt.Len(), testingx.Equal(tt.Len()))
		}

		if rt.Kind() == reflect.Map {
			testingx.Expect(t, FullTypeName(rt.Key()), testingx.Equal(FullTypeName(tt.Key())))
		}

		if rt.Kind() == reflect.Array || rt.Kind() == reflect.Slice || rt.Kind() == reflect.Map {
			testingx.Expect(t, FullTypeName(rt.Elem()), testingx.Equal(FullTypeName(tt.Elem())))
		}

		if rt.Kind() == reflect.Struct {
			testingx.Expect(t, rt.NumField(), testingx.Equal(tt.NumField()))

			for i := 0; i < rt.NumField(); i++ {
				rsf := rt.Field(i)
				tsf := tt.Field(i)

				t.Run("F "+rsf.Name(), func(t *testing.T) {
					testingx.Expect(t, rsf.Anonymous(), testingx.Equal(tsf.Anonymous()))
					testingx.Expect(t, rsf.Tag(), testingx.Equal(tsf.Tag()))
					testingx.Expect(t, rsf.Name(), testingx.Equal(tsf.Name()))
					testingx.Expect(t, rsf.PkgPath(), testingx.Equal(tsf.PkgPath()))
					testingx.Expect(t, FullTypeName(rsf.Type()), testingx.Equal(FullTypeName(tsf.Type())))

					if rsf.Type().Kind() == reflect.Struct {
						elmT := rsf.Type()

						for i := 0; i < elmT.NumField(); i++ {
							rsf := elmT.Field(i)
							tsf := elmT.Field(i)

							testingx.Expect(t, rsf.Anonymous(), testingx.Equal(tsf.Anonymous()))
							testingx.Expect(t, rsf.Tag(), testingx.Equal(tsf.Tag()))
							testingx.Expect(t, rsf.Name(), testingx.Equal(tsf.Name()))
							testingx.Expect(t, rsf.PkgPath(), testingx.Equal(tsf.PkgPath()))
							testingx.Expect(t, FullTypeName(rsf.Type()), testingx.Equal(FullTypeName(tsf.Type())))
						}
					}
				})
			}

			if rt.NumField() > 0 {
				{
					rsf, _ := rt.FieldByName("A")
					tsf, _ := tt.FieldByName("A")

					testingx.Expect(t, rsf.Anonymous(), testingx.Equal(tsf.Anonymous()))
					testingx.Expect(t, rsf.Tag(), testingx.Equal(tsf.Tag()))
					testingx.Expect(t, rsf.Name(), testingx.Equal(tsf.Name()))
					testingx.Expect(t, rsf.PkgPath(), testingx.Equal(tsf.PkgPath()))
					testingx.Expect(t, FullTypeName(rsf.Type()), testingx.Equal(FullTypeName(tsf.Type())))

					{
						_, ok := rt.FieldByName("_")
						testingx.Expect(t, ok, testingx.BeFalse())
					}
					{
						_, ok := tt.FieldByName("_")
						testingx.Expect(t, ok, testingx.BeFalse())
					}
				}

				{
					rsf, _ := rt.FieldByNameFunc(func(s string) bool {
						return s == "A"
					})
					tsf, _ := tt.FieldByNameFunc(func(s string) bool {
						return s == "A"
					})

					testingx.Expect(t, rsf.Anonymous(), testingx.Equal(tsf.Anonymous()))
					testingx.Expect(t, rsf.Tag(), testingx.Equal(tsf.Tag()))
					testingx.Expect(t, rsf.Name(), testingx.Equal(tsf.Name()))
					testingx.Expect(t, rsf.PkgPath(), testingx.Equal(tsf.PkgPath()))
					testingx.Expect(t, FullTypeName(rsf.Type()), testingx.Equal(FullTypeName(tsf.Type())))

					{
						_, ok := rt.FieldByNameFunc(func(s string) bool {
							return false
						})
						testingx.Expect(t, ok, testingx.BeFalse())
					}
					{
						_, ok := tt.FieldByNameFunc(func(s string) bool {
							return false
						})
						testingx.Expect(t, ok, testingx.BeFalse())
					}
				}
			}
		}

		if rt.Kind() == reflect.Func {
			testingx.Expect(t, rt.NumIn(), testingx.Equal(tt.NumIn()))
			testingx.Expect(t, rt.NumOut(), testingx.Equal(tt.NumOut()))

			for i := 0; i < rt.NumIn(); i++ {
				rParam := rt.In(i)
				tParam := tt.In(i)
				testingx.Expect(t, rParam.String(), testingx.Equal(tParam.String()))
			}

			for i := 0; i < rt.NumOut(); i++ {
				rResult := rt.Out(i)
				tResult := tt.Out(i)
				testingx.Expect(t, rResult.String(), testingx.Equal(tResult.String()))
			}
		}

		if rt.Kind() == reflect.Ptr {
			rt = Deref(rt).(*RType)
			tt = Deref(tt).(*TType)

			testingx.Expect(t, rt.String(), testingx.Equal(tt.String()))
		}
	})
}

func TestTryNew(t *testing.T) {
	{
		_, ok := TryNew(FromRType(reflect.TypeOf(typ.Struct{})))
		testingx.Expect(t, ok, testingx.BeTrue())
	}
	{
		_, ok := TryNew(FromTType(NewTypesTypeFromReflectType(reflect.TypeOf(typ.Struct{}))))
		testingx.Expect(t, ok, testingx.BeFalse())
	}
}

func TestEachField(t *testing.T) {
	expect := []string{
		"a", "b", "bool", "c", "Part2",
	}

	{
		rtype := FromRType(reflect.TypeOf(typ.Struct{}))
		names := make([]string, 0)
		EachField(rtype, "json", func(field StructField, fieldDisplayName string, omitempty bool) bool {
			names = append(names, fieldDisplayName)
			return true
		})
		testingx.Expect(t, expect, testingx.Equal(names))
	}

	{
		ttype := FromTType(NewTypesTypeFromReflectType(reflect.TypeOf(typ.Struct{})))
		names := make([]string, 0)
		EachField(ttype, "json", func(field StructField, fieldDisplayName string, omitempty bool) bool {
			names = append(names, fieldDisplayName)
			return true
		})
		testingx.Expect(t, expect, testingx.Equal(names))
	}
}
