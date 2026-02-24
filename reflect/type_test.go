package reflect_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/octohelm/x/cmp"
	reflectx "github.com/octohelm/x/reflect"
	. "github.com/octohelm/x/testing/v2"
)

type (
	Bytes    []byte
	Uint8    uint8
	NotBytes []Uint8
)

func BenchmarkIsBytes(b *testing.B) {
	b.Run("Raw Bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflectx.IsBytes([]byte(""))
		}
	})
	b.Run("Named Bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflectx.IsBytes(Bytes(""))
		}
	})
}

func TestIsBytes(t *testing.T) {
	t.Run("Raw Bytes", func(t *testing.T) {
		Then(t, "[]byte should be identified as bytes",
			Expect(reflectx.IsBytes([]byte("")), Be(cmp.True())),
		)
	})
	t.Run("Not Bytes", func(t *testing.T) {
		Then(t, "[]Uint8 (named uint8) should not be identified as bytes",
			Expect(reflectx.IsBytes(NotBytes("")), Be(cmp.False())),
		)
	})
	t.Run("Named Bytes", func(t *testing.T) {
		Then(t, "Named []byte should still be identified as bytes",
			Expect(reflectx.IsBytes(Bytes("")), Be(cmp.True())),
		)
	})
	t.Run("Others", func(t *testing.T) {
		Then(t, "non-slice types should be false",
			Expect(reflectx.IsBytes(""), Be(cmp.False())),
			Expect(reflectx.IsBytes(true), Be(cmp.False())),
		)
	})
}

func TestFullTypeName(t *testing.T) {
	t.Run("GIVEN various types", func(t *testing.T) {
		Then(t, "should return correct full name string",
			Expect(
				reflectx.FullTypeName(reflect.TypeFor[*int]()),
				Equal("*int"),
			),
			Expect(
				reflectx.FullTypeName(reflect.PointerTo(reflect.TypeFor[int]())),
				Equal("*int"),
			),
			Expect(
				reflectx.FullTypeName(reflect.PointerTo(reflect.TypeOf(time.Now()))),
				Equal("*time.Time"),
			),
			Expect(
				reflectx.FullTypeName(reflect.PointerTo(reflect.TypeFor[struct{ Name string }]())),
				Equal("*struct { Name string }"),
			),
		)
	})
}

func TestIndirectType(t *testing.T) {
	t.Run("GIVEN a pointer type", func(t *testing.T) {
		expected := reflect.TypeFor[int]()

		Then(t, "Deref should return the underlying element type",
			Expect(reflectx.Deref(reflect.TypeFor[*int]()), Equal(expected)),
			Expect(reflectx.Deref(reflect.PointerTo(reflect.TypeFor[int]())), Equal(expected)),
		)
	})

	t.Run("WHEN having deep nested pointers", func(t *testing.T) {
		tpe := reflect.TypeFor[int]()
		for range 10 {
			tpe = reflect.PointerTo(tpe)
		}

		Then(t, "Deref should recursively unwrap all levels",
			Expect(reflectx.Deref(tpe), Equal(reflect.TypeFor[int]())),
		)
	})
}
