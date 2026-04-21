package types_test

import (
	"go/types"
	"reflect"
	"testing"

	"github.com/octohelm/x/cmp"
	. "github.com/octohelm/x/testing/v2"
	. "github.com/octohelm/x/types"
	"github.com/octohelm/x/types/testdata/typ"
)

func TestTypeHelpers(t *testing.T) {
	t.Run("FieldDisplayName", func(t *testing.T) {
		cases := []struct {
			name        string
			tag         reflect.StructTag
			defaultName string
			displayName string
			omitempty   bool
			keepNested  bool
		}{
			{name: "missing tag", tag: `xml:"x"`, defaultName: "Name", displayName: "Name", keepNested: false},
			{name: "empty tag", tag: `json:""`, defaultName: "Name", displayName: "Name", keepNested: true},
			{name: "leading comma", tag: `json:",omitempty"`, defaultName: "Name", displayName: "Name", omitempty: true, keepNested: true},
			{name: "renamed", tag: `json:"alias"`, defaultName: "Name", displayName: "alias", keepNested: true},
			{name: "renamed with omitempty", tag: `json:"alias,omitempty"`, defaultName: "Name", displayName: "alias", omitempty: true, keepNested: true},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				name, omitempty, keepNested := FieldDisplayName(c.tag, "json", c.defaultName)

				Then(t, "should parse tag flags correctly",
					Expect(name, Equal(c.displayName)),
					Expect(omitempty, Equal(c.omitempty)),
					Expect(keepNested, Equal(c.keepNested)),
				)
			})
		}
	})

	t.Run("Deref and FullTypeName", func(t *testing.T) {
		ptrType := PtrTo(PtrTo(FromRType(reflect.TypeFor[typ.Struct]())))

		Then(t, "Deref should unwrap pointer chain",
			Expect(Deref(ptrType).String(), Equal("github.com/octohelm/x/types/testdata/typ.Struct")),
		)
		Then(t, "FullTypeName should keep pointer prefix",
			Expect(FullTypeName(ptrType), Equal("**github.com/octohelm/x/types/testdata/typ.Struct")),
		)
		Then(t, "nil type should render nil", Expect(FullTypeName(nil), Equal("nil")))
	})

	t.Run("TypeByName and PtrTo", func(t *testing.T) {
		rptr := PtrTo(FromRType(reflect.TypeFor[typ.String]()))
		tptr := PtrTo(FromTType(types.Typ[types.String]))

		Then(t, "TypeByName should support builtin lookup",
			Expect(FromTType(TypeByName("", "string")).String(), Equal("string")),
		)
		Then(t, "NewPackage should allow empty import path",
			Expect(NewPackage(""), Be(cmp.Nil[*types.Package]())),
		)
		Then(t, "PtrTo should preserve concrete type family",
			Expect(FullTypeName(rptr), Equal("*github.com/octohelm/x/types/testdata/typ.String")),
			Expect(tptr.Kind(), Equal(reflect.Pointer)),
			Expect(tptr.Elem().String(), Equal("string")),
		)
	})

	t.Run("EncodingTextMarshalerTypeReplacer", func(t *testing.T) {
		replacedRType, okForEnum := EncodingTextMarshalerTypeReplacer(FromRType(reflect.TypeFor[typ.Enum]()))
		replacedTType, okForNamed := EncodingTextMarshalerTypeReplacer(FromTType(TypeFor("github.com/octohelm/x/types/testdata/typ.Enum")))
		_, okForPlain := EncodingTextMarshalerTypeReplacer(FromRType(reflect.TypeFor[int]()))

		Then(t, "types implementing TextMarshaler should be replaced by string",
			Expect(okForEnum, Be(cmp.True())),
			Expect(okForNamed, Be(cmp.True())),
			Expect(replacedRType.String(), Equal("string")),
			Expect(replacedTType.String(), Equal("string")),
		)
		Then(t, "plain types should not be replaced",
			Expect(okForPlain, Be(cmp.False())),
		)
	})
}
