package testing

import (
	"os"
	"path/filepath"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/octohelm/x/anyjson"
)

func ProjectRoot() string {
	p, _ := os.Getwd()
	for {
		if p == "/" {
			break
		}
		if _, err := os.Stat(filepath.Join(p, "go.mod")); err == nil {
			return p
		}
		p = filepath.Dir(p)
	}
	return p
}

func MustAsJSON(v any) []byte {
	raw, err := AsJSON(v)
	if err != nil {
		panic(err)
	}
	return raw
}

func AsJSON(v any) ([]byte, error) {
	vv, err := anyjson.FromValue(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(anyjson.Sorted(vv), jsontext.WithIndent("  "))
}
