package reflect_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/octohelm/x/ptr"
	reflectx "github.com/octohelm/x/reflect"
	testingx "github.com/octohelm/x/testing"
)

type Bytes []byte

type Uint8 uint8

type NotBytes []Uint8

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
		testingx.Expect(t, reflectx.IsBytes([]byte("")), testingx.BeTrue())
	})
	t.Run("Not Bytes", func(t *testing.T) {
		testingx.Expect(t, reflectx.IsBytes(NotBytes("")), testingx.BeFalse())
	})
	t.Run("Named Bytes", func(t *testing.T) {
		testingx.Expect(t, reflectx.IsBytes(Bytes("")), testingx.BeTrue())
	})
	t.Run("Others", func(t *testing.T) {
		testingx.Expect(t, reflectx.IsBytes(""), testingx.BeFalse())
		testingx.Expect(t, reflectx.IsBytes(true), testingx.BeFalse())
	})
}

func TestFullTypeName(t *testing.T) {
	testingx.Expect(t, reflectx.FullTypeName(reflect.TypeOf(ptr.Ptr(1))), testingx.Equal("*int"))
	testingx.Expect(t, reflectx.FullTypeName(reflect.PtrTo(reflect.TypeOf(1))), testingx.Equal("*int"))
	testingx.Expect(t, reflectx.FullTypeName(reflect.PtrTo(reflect.TypeOf(time.Now()))), testingx.Equal("*time.Time"))
	testingx.Expect(t, reflectx.FullTypeName(reflect.PtrTo(reflect.TypeOf(struct {
		Name string
	}{}))), testingx.Equal("*struct { Name string }"))
}

func TestIndirectType(t *testing.T) {
	testingx.Expect(t, reflect.TypeOf(1), testingx.Equal(reflectx.Deref(reflect.TypeOf(ptr.Ptr(1)))))
	testingx.Expect(t, reflect.TypeOf(1), testingx.Equal(reflectx.Deref(reflect.PtrTo(reflect.TypeOf(1)))))

	tpe := reflect.TypeOf(1)
	for i := 0; i < 10; i++ {
		tpe = reflect.PtrTo(tpe)
	}
	testingx.Expect(t, reflect.TypeOf(1), testingx.Equal(reflectx.Deref(tpe)))
}
