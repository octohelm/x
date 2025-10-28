package datauri

import (
	"encoding"
	"encoding/base64"
	"errors"
	"mime"
	"net/url"
	"strings"
)

var (
	ErrInvalidDataURI = errors.New("invalid data uri")
)

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

type DataURI struct {
	MediaType string
	Params    map[string]string
	Data      []byte
}

func (d DataURI) IsZero() bool {
	return len(d.Data) == 0
}

func (DataURI) OpenAPISchemaFormat() string {
	return "data-uri"
}

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

func (d DataURI) String() string {
	return d.Encoded(true)
}

var _ encoding.TextMarshaler = (*DataURI)(nil)

func (d DataURI) MarshalText() ([]byte, error) {
	return []byte(d.Encoded(true)), nil
}

var _ encoding.TextUnmarshaler = (*DataURI)(nil)

func (d *DataURI) UnmarshalText(text []byte) error {
	dd, err := Parse(string(text))
	if err != nil {
		return err
	}
	*d = *dd
	return nil
}
