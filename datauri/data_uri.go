package datauri

import (
	"encoding"
	"encoding/base64"
	"errors"
	"mime"
	"net/url"
	"strings"
)

// ErrInvalidDataURI 表示输入不符合 data URI 语法。
var ErrInvalidDataURI = errors.New("invalid data uri")

// Parse 解析 data URI 文本，或将无前缀输入按 base64 数据体处理。
func Parse(dataURI string) (*DataURI, error) {
	withoutPrefix := ""
	if strings.HasPrefix(dataURI, "data:") {
		withoutPrefix = dataURI[len("data:"):]
	} else {
		withoutPrefix = ";base64," + dataURI
	}

	parts := strings.SplitN(withoutPrefix, ",", 2)
	if len(parts) != 2 {
		return nil, ErrInvalidDataURI
	}

	meta, raw := parts[0], parts[1]
	isBase64 := strings.Contains(meta, ";base64")
	if isBase64 {
		meta = strings.Replace(meta, ";base64", "", 1)
	}

	mediaType, params, err := mime.ParseMediaType(meta)
	if err != nil {
		mediaType = ""
	}

	d := &DataURI{
		MediaType: mediaType,
		Params:    params,
	}

	if isBase64 {
		data, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			return nil, errors.Join(ErrInvalidDataURI, err)
		}
		d.Data = data

		return d, nil
	}

	s, err := url.PathUnescape(raw)
	if err != nil {
		return nil, errors.Join(ErrInvalidDataURI, err)
	}

	d.Data = []byte(s)

	return d, nil
}

// DataURI 表示一个解析后的 data URI。
type DataURI struct {
	// MediaType 是 data URI 中声明的媒体类型，例如 text/plain 或 image/png。
	MediaType string
	// Params 保存媒体类型后的附加参数，例如 charset。
	Params map[string]string
	// Data 是解码后的原始数据体。
	Data []byte
}

// IsZero 判断 data URI 是否不含数据体。
func (d DataURI) IsZero() bool {
	return len(d.Data) == 0
}

// OpenAPISchemaFormat 返回 data URI 的 OpenAPI format 名称。
func (DataURI) OpenAPISchemaFormat() string {
	return "data-uri"
}

// Encoded 按指定是否使用 base64 返回 data URI 文本。
func (d *DataURI) Encoded(base64Encoded bool) string {
	s := &strings.Builder{}
	s.WriteString("data:")
	s.WriteString(strings.ReplaceAll(mime.FormatMediaType(d.MediaType, d.Params), "; ", ";"))

	if base64Encoded {
		s.WriteString(";base64,")
		s.WriteString(base64.StdEncoding.EncodeToString(d.Data))
	} else {
		s.WriteString(",")
		s.WriteString(url.PathEscape(string(d.Data)))
	}
	return s.String()
}

// String 返回 base64 形式的 data URI 文本。
func (d DataURI) String() string {
	return d.Encoded(true)
}

var _ encoding.TextMarshaler = (*DataURI)(nil)

// MarshalText 将 data URI 编码为文本。
func (d DataURI) MarshalText() ([]byte, error) {
	return []byte(d.Encoded(true)), nil
}

var _ encoding.TextUnmarshaler = (*DataURI)(nil)

// UnmarshalText 从文本解析 data URI。
func (d *DataURI) UnmarshalText(text []byte) error {
	dd, err := Parse(string(text))
	if err != nil {
		return err
	}
	*d = *dd
	return nil
}
