package datauri

import (
	"encoding/base64"
	"reflect"
	"regexp"
	"testing"

	"github.com/octohelm/x/cmp"
	. "github.com/octohelm/x/testing/v2"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name string
		uri  string
		want DataURI
	}{
		{
			name: "百分号编码文本",
			uri:  "data:,A%20brief%20note",
			want: DataURI{
				Data: []byte("A brief note"),
			},
		},
		{
			name: `带参数的文本`,
			uri:  `data:text/plain;charset=utf-8;filename="file x",A%20brief%20note`,
			want: DataURI{
				MediaType: "text/plain",
				Params: map[string]string{
					"charset":  "utf-8",
					"filename": "file x",
				},
				Data: []byte("A brief note"),
			},
		},
		{
			name: "省略 data 前缀时按 base64 处理",
			uri:  base64.StdEncoding.EncodeToString([]byte("A brief note")),
			want: DataURI{
				Data: []byte("A brief note"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			parsed := MustValue(t, func() (*DataURI, error) {
				return Parse(c.uri)
			})

			Then(t, "解析结果符合预期",
				Expect(*parsed, Equal(c.want)),
			)
		})
	}
}

func TestParseErrors(t *testing.T) {
	t.Run("缺少逗号分隔符", func(t *testing.T) {
		Then(t, "返回无效 data uri 错误",
			ExpectDo(func() error {
				_, err := Parse("data:text/plain;base64")
				return err
			}, ErrorIs(ErrInvalidDataURI)),
		)
	})

	t.Run("非法 base64", func(t *testing.T) {
		Then(t, "错误链保留 ErrInvalidDataURI",
			ExpectDo(func() error {
				_, err := Parse("data:text/plain;base64,%%%")
				return err
			}, ErrorIs(ErrInvalidDataURI)),
		)
	})

	t.Run("非法转义文本", func(t *testing.T) {
		Then(t, "错误链保留 ErrInvalidDataURI",
			ExpectDo(func() error {
				_, err := Parse("data:text/plain,%zz")
				return err
			}, ErrorIs(ErrInvalidDataURI)),
		)
	})
}

func TestEncodingRoundTrip(t *testing.T) {
	value := DataURI{
		MediaType: "text/plain",
		Params: map[string]string{
			"charset": "utf-8",
		},
		Data: []byte("hello world"),
	}

	t.Run("文本与 base64 编码可往返", func(t *testing.T) {
		base64Text := value.Encoded(true)
		plainText := value.Encoded(false)

		Then(t, "不同编码路径都可恢复为同一值",
			ExpectMustValue(func() (*DataURI, error) {
				return Parse(base64Text)
			}, Be(func(actual *DataURI) error {
				if reflect.DeepEqual(value, *actual) {
					return nil
				}
				return &ErrNotEqual{Expect: value, Got: *actual}
			})),
			ExpectMustValue(func() (*DataURI, error) {
				return Parse(plainText)
			}, Be(func(actual *DataURI) error {
				if reflect.DeepEqual(value, *actual) {
					return nil
				}
				return &ErrNotEqual{Expect: value, Got: *actual}
			})),
		)
	})

	t.Run("MarshalText 与 UnmarshalText", func(t *testing.T) {
		text := MustValue(t, value.MarshalText)

		var decoded DataURI

		Then(t, "文本解码应成功",
			ExpectDo(func() error {
				return decoded.UnmarshalText(text)
			}),
		)

		Then(t, "文本编解码保持一致",
			Expect(decoded, Equal(value)),
			Expect(string(text), Equal(value.String())),
		)
	})
}

func TestZeroAndFormat(t *testing.T) {
	t.Run("零值与格式", func(t *testing.T) {
		var zero DataURI

		Then(t, "零值无数据且保留 data-uri format",
			Expect(zero.IsZero(), Be(cmp.True())),
			Expect(zero.OpenAPISchemaFormat(), Equal("data-uri")),
			Expect(zero.String(), Be(func(actual string) error {
				if regexp.MustCompile(`^data:;base64,$`).MatchString(actual) {
					return nil
				}
				return &ErrNotEqual{Expect: "^data:;base64,$", Got: actual}
			})),
		)
	})

	t.Run("非法媒体类型当前会被忽略", func(t *testing.T) {
		parsed := MustValue(t, func() (*DataURI, error) {
			return Parse("data:%zz,abc")
		})

		Then(t, "解析继续进行并保留原始媒体类型文本",
			Expect(parsed.MediaType, Equal("%zz")),
			Expect(string(parsed.Data), Equal("abc")),
		)
	})
}
