package datauri

import (
	"strings"
	"testing"

	"github.com/octohelm/x/testing/bdd"
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
		bdd.FromT(t).Given(c.URI, func(b bdd.T) {
			b.When("parse", func(b bdd.T) {
				dataURI, err := Parse(c.URI)
				b.Then("success",
					bdd.NoError(err),
					bdd.Equal(c.DataURI, *dataURI),
				)
			})

			b.When("encoded", func(b bdd.T) {
				uri := c.DataURI.Encoded(strings.Contains(c.URI, ";base64,"))

				b.Then("success",
					bdd.Equal(c.URI, uri),
				)
			})
		})
	}

	bdd.FromT(t).When("parse without data proto", func(b bdd.T) {
		dataURI, err := Parse("QSBicmllZiBub3Rl")

		b.Then("success",
			bdd.NoError(err),
			bdd.Equal(string(dataURI.Data), `A brief note`),
		)
	})
}
