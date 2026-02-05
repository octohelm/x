package anyjson_test

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/octohelm/x/anyjson"
	. "github.com/octohelm/x/testing/v2"
)

type Employee struct {
	Name    string `json:"name"`
	Salary  int    `json:"salary"`
	Married bool   `json:"married"`
	Age     int    `json:"age,omitempty"`
}

func TestFrom(t *testing.T) {
	t.Run("GIVEN a struct with employees", func(t *testing.T) {
		x := struct {
			Employees []Employee `json:"employees"`
		}{
			Employees: []Employee{
				{
					Name:    "octo",
					Salary:  56000,
					Married: false,
				},
			},
		}

		// 使用 MustValue 表达：转换必须成功，转换后的 obj 是后续验证的基础
		obj := MustValue(t, func() (anyjson.Valuer, error) {
			return anyjson.FromValue(x)
		})

		Then(t, "it should be converted to Object/List hierarchy",
			Expect(obj.Value(), Equal[any](anyjson.Obj{
				"employees": anyjson.List{
					anyjson.Obj{
						"name":    "octo",
						"salary":  56000,
						"married": false,
					},
				},
			})),
		)
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("WHEN unmarshal raw JSON to Object", func(t *testing.T) {
		raw := []byte(`
{  
    "employees": [
       {  
          "name":      "octo",   
          "salary":     56000,   
          "married":    false  
       }  
    ] 
}
`)
		var obj anyjson.Object
		Then(t, "unmarshal should success",
			ExpectMust(func() error {
				return json.Unmarshal(raw, &obj)
			}),
		)

		Then(t, "the content should match expected Obj/List structure",
			Expect(obj.Value(), Equal[any](anyjson.Obj{
				"employees": anyjson.List{
					anyjson.Obj{
						"name":    "octo",
						"salary":  56000,
						"married": false,
					},
				},
			})),
		)
	})
}
