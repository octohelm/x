package snapshot

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/octohelm/x/anyjson"
)

func AsJSON(v any) ([]byte, error) {
	vv, err := anyjson.FromValue(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(anyjson.Sorted(vv), jsontext.WithIndent("  "))
}
