package snapshot

import (
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/octohelm/x/anyjson"
)

// AsJSON 将值稳定化后编码为格式化 JSON。
func AsJSON(v any) ([]byte, error) {
	vv, err := anyjson.FromValue(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(anyjson.Sorted(vv), jsontext.WithIndent("  "))
}
