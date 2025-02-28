package types

import (
	"reflect"
	"testing"

	testingx "github.com/octohelm/x/testing"

	. "github.com/onsi/gomega"
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

	for i := range cases {
		c := cases[i]
		NewWithT(t).Expect(FromTType(TypeFor(c)).String()).To(Equal(c))
	}
}

func Test_issue_for_chan(t *testing.T) {
	x := make(chan struct{})
	defer close(x)

	c := func(recv chan<- struct{}) <-chan struct{} {
		return x
	}

	t.Run("Both Chan", func(t *testing.T) {
		typ := NewTypesTypeFromReflectType(reflect.TypeOf(x))
		testingx.Expect(t, typ.String(), testingx.Be("chan struct{}"))
	})

	t.Run("Send Chan", func(t *testing.T) {
		typ := NewTypesTypeFromReflectType(reflect.TypeOf(c).In(0))
		testingx.Expect(t, typ.String(), testingx.Be("chan<- struct{}"))
	})

	t.Run("Recv Chan", func(t *testing.T) {
		typ := NewTypesTypeFromReflectType(reflect.TypeOf(c).Out(0))
		testingx.Expect(t, typ.String(), testingx.Be("<-chan struct{}"))
	})
}
