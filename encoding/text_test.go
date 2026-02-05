package encoding

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/octohelm/x/cmp"
	"github.com/octohelm/x/ptr"
	"github.com/octohelm/x/slices"

	. "github.com/octohelm/x/testing/v2"
)

type Duration time.Duration

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

func (d *Duration) UnmarshalText(data []byte) error {
	dur, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

type (
	NamedString string
	NamedInt    int
)

var longBytes = strings.Join(slices.Map(make([]string, 1025), func(e string) string {
	return "1"
}), "")

var (
	v = struct {
		NamedString  NamedString
		NamedInt     NamedInt
		Duration     Duration
		PtrDuration  *Duration
		String       string
		PtrString    *string
		PtrPtrString **string
		Int          int
		PtrInt       *int
		Uint         uint
		PtrUint      *uint
		Float        float32
		PtrFloat     *float32
		Bool         bool
		PtrBool      *bool
		Bytes        []byte
		LongBytes    []byte
	}{}

	rv = reflect.ValueOf(&v).Elem()
	d  = Duration(2 * time.Second)
)

var cases = []struct {
	name   string
	v      any
	text   string
	expect any
}{
	{
		"Ptr String",
		rv.FieldByName("PtrString"),
		"string",
		ptr.Ptr("string"),
	},
	{
		"Ptr Ptr String",
		rv.FieldByName("PtrPtrString"),
		"string",
		func() **string {
			s := ptr.Ptr("string")
			return &s
		}(),
	},
	{
		"Ptr String raw value",
		&v.String,
		"ptr",
		ptr.Ptr("ptr"),
	},
	{
		"Named String",
		rv.FieldByName("NamedString"),
		"string",
		NamedString("string"),
	},
	{
		"Duration",
		rv.FieldByName("Duration"),
		"2s",
		Duration(2 * time.Second),
	},
	{
		"Ptr Duration",
		rv.FieldByName("PtrDuration"),
		"2s",
		&d,
	},
	{
		"Int",
		rv.FieldByName("Int"),
		"1",
		1,
	},
	{
		"Named Int",
		rv.FieldByName("NamedInt"),
		"11",
		NamedInt(11),
	},
	{
		"PtrInt",
		rv.FieldByName("PtrInt"),
		"1",
		ptr.Ptr(1),
	},
	{
		"Uint",
		rv.FieldByName("Uint"),
		"1",
		uint(1),
	},
	{
		"Int raw value",
		rv.FieldByName("Int").Addr().Interface(),
		"1",
		ptr.Ptr(1),
	},
	{
		"PtrUint",
		rv.FieldByName("PtrUint"),
		"1",
		ptr.Ptr[uint](1),
	},
	{
		"Float",
		rv.FieldByName("Float"),
		"1",
		float32(1),
	},
	{
		"PtrFloat",
		rv.FieldByName("PtrFloat"),
		"1.1",
		ptr.Ptr[float32](1.1),
	},
	{
		"Bool",
		rv.FieldByName("Bool"),
		"true",
		true,
	},
	{
		"PtrBool",
		rv.FieldByName("PtrBool"),
		"true",
		ptr.Ptr(true),
	},
	{
		"Bytes",
		rv.FieldByName("Bytes"),
		"MTEx",
		[]byte("111"),
	},
	{
		"Bytes direct",
		&v.Bytes,
		"MTEx",
		func() *[]byte {
			b := []byte("111")
			return &b
		}(),
	},
	{
		"LongBytes direct",
		&v.LongBytes,
		base64.StdEncoding.EncodeToString([]byte(longBytes)),
		func() *[]byte {
			b := []byte(longBytes)
			return &b
		}(),
	},
}

func BenchmarkPtrFloat(b *testing.B) {
	v.PtrFloat = ptr.Ptr[float32](1.1)
	// rv := reflect.ValueOf(v.PtrFloat).Elem()

	b.Run("append", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			//f := rv.Float()
			//_, _ = MarshalText(v.PtrFloat)
			d := make([]byte, 0)
			strconv.AppendFloat(d, float64(*v.PtrFloat), 'f', -1, 32)
		}

		// fmt.Println(string(d))
	})

	b.Run("format", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			//f := rv.Float()
			//_, _ = MarshalText(v.PtrFloat)
			_ = []byte(strconv.FormatFloat(float64(*v.PtrFloat), 'f', -1, 32))
		}
		// fmt.Println(string(d))
	})
}

func BenchmarkUnmarshalTextAndMarshalText(b *testing.B) {
	for i := range cases {
		c := cases[i]

		b.Run(fmt.Sprintf("UnmarshalText %s", c.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = UnmarshalText(c.v, []byte(c.text))
			}
		})

		b.Run(fmt.Sprintf("MarshalText %s", c.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = MarshalText(c.v)
			}
		})
	}
}

func TestUnmarshalTextAndMarshalText(t *testing.T) {
	for _, c := range cases {
		t.Run(fmt.Sprintf("UnmarshalText %s", c.name), func(t *testing.T) {
			Then(t, "success",
				ExpectMust(func() error {
					return UnmarshalText(c.v, []byte(c.text))
				}),
			)

			actual := c.v
			if rv, ok := c.v.(reflect.Value); ok {
				actual = rv.Interface()
			}

			Then(t, "value as expected",
				Expect(actual, Equal(c.expect)),
			)
		})
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("MarshalText %s", c.name), func(t *testing.T) {
			// 使用 MustValue 立即获取结果，因为它不仅是断言，还是下一步输入
			text := MustValue(t, func() ([]byte, error) {
				return MarshalText(c.v)
			})

			Then(t, "text as expected",
				Expect(string(text), Equal(c.text)),
			)
		})
	}

	t.Run("GIVEN reflect fields", func(t *testing.T) {
		v2 := struct {
			PtrString *string
			Slice     []string
		}{}
		rv2 := reflect.ValueOf(v2)

		t.Run("PtrString", func(t *testing.T) {
			Then(t, "success",
				ExpectMust(func() error {
					_, err := MarshalText(rv2.FieldByName("PtrString"))
					return err
				}),
			)
		})

		t.Run("Slice", func(t *testing.T) {
			Then(t, "should has error",
				ExpectDo(
					func() error {
						_, err := MarshalText(rv2.FieldByName("Slice"))
						return err
					},
					Be(cmp.NotNil[error]()),
				),
			)
		})
	})
}
