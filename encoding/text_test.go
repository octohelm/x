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

type failText struct{}

func (failText) MarshalText() ([]byte, error) {
	return nil, fmt.Errorf("marshal fail")
}

func (*failText) UnmarshalText(data []byte) error {
	return fmt.Errorf("unmarshal fail")
}

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
		new("string"),
	},
	{
		"Ptr Ptr String",
		rv.FieldByName("PtrPtrString"),
		"string",
		func() **string {
			s := new("string")
			return &s
		}(),
	},
	{
		"Ptr String raw value",
		&v.String,
		"ptr",
		new("ptr"),
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
		new(1),
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
		new(1),
	},
	{
		"PtrUint",
		rv.FieldByName("PtrUint"),
		"1",
		new(uint(1)),
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
		new(float32(1.1)),
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
		new(true),
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
	v.PtrFloat = new(float32(1.1))
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

	t.Run("边界与错误路径", func(t *testing.T) {
		t.Run("MarshalText 处理 nil", func(t *testing.T) {
			var nilPtr *int

			rawNil := MustValue(t, func() ([]byte, error) {
				return MarshalText(nil)
			})
			rawPtrNil := MustValue(t, func() ([]byte, error) {
				return MarshalText(nilPtr)
			})

			Then(t, "nil 与 nil 指针都返回 nil 文本",
				Expect(rawNil, Equal([]byte(nil))),
				Expect(rawPtrNil, Equal([]byte(nil))),
			)
		})

		t.Run("MarshalText 不支持的类型返回错误", func(t *testing.T) {
			Then(t, "slice of string 返回 unsupported type",
				ExpectDo(func() error {
					_, err := MarshalText([]string{"x"})
					return err
				}, Be(cmp.NotNil[error]())),
			)
		})

		t.Run("MarshalText 直接标量与自定义编码分支", func(t *testing.T) {
			for _, c := range []struct {
				name string
				v    any
				text string
			}{
				{"int8", int8(8), "8"},
				{"int16", int16(16), "16"},
				{"int32", int32(32), "32"},
				{"int64", int64(64), "64"},
				{"uint8", uint8(8), "8"},
				{"uint16", uint16(16), "16"},
				{"uint32", uint32(32), "32"},
				{"uint64", uint64(64), "64"},
				{"float64", float64(1.25), "1.25"},
				{"bool false", false, "false"},
				{"string", "text", "text"},
			} {
				t.Run(c.name, func(t *testing.T) {
					raw := MustValue(t, func() ([]byte, error) {
						return MarshalText(c.v)
					})

					Then(t, "应命中对应直接类型编码分支",
						Expect(string(raw), Equal(c.text)),
					)
				})
			}

			Then(t, "自定义 TextMarshaler 错误应透传",
				ExpectDo(func() error {
					_, err := MarshalText(failText{})
					return err
				}, Be(cmp.NotNil[error]())),
			)
		})

		t.Run("UnmarshalText 非法输入返回错误", func(t *testing.T) {
			var i int
			var b bool
			var data []byte

			Then(t, "非法数值、布尔和 base64 应失败",
				ExpectDo(func() error { return UnmarshalText(&i, []byte("x")) }, Be(cmp.NotNil[error]())),
				ExpectDo(func() error { return UnmarshalText(&b, []byte("not-bool")) }, Be(cmp.NotNil[error]())),
				ExpectDo(func() error { return UnmarshalText(&data, []byte("%%%")) }, Be(cmp.NotNil[error]())),
			)
		})

		t.Run("UnmarshalText 直接指针分支与自定义解码错误", func(t *testing.T) {
			var (
				i8  int8
				i16 int16
				i32 int32
				i64 int64
				u8  uint8
				u16 uint16
				u32 uint32
				u64 uint64
				f64 float64
			)

			Then(t, "直接指针类型应完成解码",
				ExpectDo(func() error { return UnmarshalText(&i8, []byte("8")) }),
				ExpectDo(func() error { return UnmarshalText(&i16, []byte("16")) }),
				ExpectDo(func() error { return UnmarshalText(&i32, []byte("32")) }),
				ExpectDo(func() error { return UnmarshalText(&i64, []byte("64")) }),
				ExpectDo(func() error { return UnmarshalText(&u8, []byte("8")) }),
				ExpectDo(func() error { return UnmarshalText(&u16, []byte("16")) }),
				ExpectDo(func() error { return UnmarshalText(&u32, []byte("32")) }),
				ExpectDo(func() error { return UnmarshalText(&u64, []byte("64")) }),
				ExpectDo(func() error { return UnmarshalText(&f64, []byte("1.25")) }),
			)

			Then(t, "解码结果应正确写入",
				Expect(i8, Equal(int8(8))),
				Expect(i16, Equal(int16(16))),
				Expect(i32, Equal(int32(32))),
				Expect(i64, Equal(int64(64))),
				Expect(u8, Equal(uint8(8))),
				Expect(u16, Equal(uint16(16))),
				Expect(u32, Equal(uint32(32))),
				Expect(u64, Equal(uint64(64))),
				Expect(f64, Equal(float64(1.25))),
			)

			target := &failText{}
			Then(t, "自定义 TextUnmarshaler 错误在直接接口和 reflect.Value 路径都应透传",
				ExpectDo(func() error { return UnmarshalText(target, []byte("x")) }, Be(cmp.NotNil[error]())),
				ExpectDo(func() error { return UnmarshalText(reflect.ValueOf(target), []byte("x")) }, Be(cmp.NotNil[error]())),
			)
		})

		t.Run("UnmarshalText 处理 reflect.Value 非指针输入", func(t *testing.T) {
			value := 0
			rv := reflect.ValueOf(value)
			panicMsg := capturePanic(func() {
				_ = UnmarshalText(rv, []byte("1"))
			})

			Then(t, "不可寻址的 reflect.Value 当前会 panic",
				Expect(panicMsg == "", Equal(false)),
			)
		})

		t.Run("UnmarshalText 处理多级 nil 指针 reflect.Value", func(t *testing.T) {
			var target **string

			Then(t, "应自动初始化指针链并写入值",
				ExpectDo(func() error {
					return UnmarshalText(reflect.ValueOf(&target).Elem(), []byte("hello"))
				}),
			)

			Then(t, "初始化后的值应可读取",
				Expect(**target, Equal("hello")),
			)
		})

		t.Run("UnmarshalText nil 输入当前会 panic", func(t *testing.T) {
			panicMsg := capturePanic(func() {
				_ = UnmarshalText(nil, []byte("1"))
			})

			Then(t, "nil 输入当前不是返回 error 而是 panic",
				Expect(panicMsg == "", Equal(false)),
			)
		})
	})
}

func capturePanic(do func()) (msg string) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Sprint(r)
		}
	}()
	do()
	return ""
}
