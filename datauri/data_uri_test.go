package datauri

import (
	"strings"
	"testing"

	. "github.com/octohelm/x/testing/v2"
)

func TestDataURI(t *testing.T) {
	cases := []struct {
		URI     string
		DataURI DataURI
	}{
		{
			URI: "data:,A%20brief%20note",
			DataURI: DataURI{
				Data: []byte("A brief note"),
			},
		},
		{
			URI: `data:text/plain;charset=utf-8;filename="file x",A%20brief%20note`,
			DataURI: DataURI{
				MediaType: "text/plain",
				Params: map[string]string{
					"charset":  "utf-8",
					"filename": "file x",
				},
				Data: []byte("A brief note"),
			},
		},
	}

	for _, c := range cases {
		t.Run("GIVEN "+c.URI, func(t *testing.T) {
			t.Run("WHEN parse", func(t *testing.T) {
				// 使用 MustValue，因为 dataURI 是后续断言的前提
				dataURI := MustValue(t, func() (*DataURI, error) {
					return Parse(c.URI)
				})

				Then(t, "success",
					Expect(*dataURI, Equal(c.DataURI)),
				)
			})

			t.Run("WHEN encoded", func(t *testing.T) {
				isBase64 := strings.Contains(c.URI, ";base64,")
				uri := c.DataURI.Encoded(isBase64)

				Then(t, "success",
					Expect(uri, Equal(c.URI)),
				)
			})
		})
	}

	t.Run("WHEN parse without data proto", func(t *testing.T) {
		dataURI := MustValue(t, func() (*DataURI, error) {
			return Parse("QSBicmllZiBub3Rl")
		})

		Then(t, "success",
			Expect(string(dataURI.Data), Equal(`A brief note`)),
		)
	})
}
