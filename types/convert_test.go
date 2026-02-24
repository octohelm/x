package types_test

import (
	"reflect"
	"testing"

	"github.com/octohelm/x/cmp"
	. "github.com/octohelm/x/testing/v2"
	. "github.com/octohelm/x/types"
)

func TestTypeFor(t *testing.T) {
	cases := []string{
		"string",
		"int",
		"map[int]int",
		"[]int",
		"[2]int",
		"error",

		"github.com/octohelm/x/types/testdata/typ.String",
		"github.com/octohelm/x/types/testdata/typ.AnyMap[int,string]",
	}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			Then(t, "type string should match",
				Expect(
					FromTType(TypeFor(c)).String(),
					Equal(c),
				),
			)
		})
	}
}

func Test_issue_for_chan(t *testing.T) {
	x := make(chan struct{})
	defer close(x)

	t.Run("GIVEN a set of channels", func(t *testing.T) {
		t.Run("WHEN reflect Both Chan", func(t *testing.T) {
			typ := NewTypesTypeFromReflectType(reflect.TypeFor[chan struct{}]())

			Then(t, "should be bidirectional",
				Expect(typ.String(),
					Be(cmp.Eq("chan struct{}"))),
			)
		})

		t.Run("WHEN reflect Send Chan", func(t *testing.T) {
			typ := NewTypesTypeFromReflectType(reflect.TypeFor[func(recv chan<- struct{}) <-chan struct{}]().In(0))

			Then(t, "should be send-only",
				Expect(typ.String(),
					Be(cmp.Eq("chan<- struct{}"))),
			)
		})

		t.Run("WHEN reflect Recv Chan", func(t *testing.T) {
			typ := NewTypesTypeFromReflectType(reflect.TypeFor[func(recv chan<- struct{}) <-chan struct{}]().Out(0))

			Then(t, "should be receive-only",
				Expect(typ.String(),
					Be(cmp.Eq("<-chan struct{}"))),
			)
		})
	})
}
