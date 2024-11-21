package anyjson

import (
	"testing"

	"github.com/go-json-experiment/json"
	testingx "github.com/octohelm/x/testing"
)

type Employee struct {
	Name    string `json:"name"`
	Salary  int    `json:"salary"`
	Married bool   `json:"married"`
	Age     int    `json:"age,omitempty"`
}

func TestFrom(t *testing.T) {
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

	obj, err := FromValue(x)
	testingx.Expect(t, err, testingx.BeNil[error]())

	testingx.Expect(t, obj.Value(), testingx.Equal[any](Obj{
		"employees": List{
			Obj{
				"name":    "octo",
				"salary":  56000,
				"married": false,
			},
		},
	}))
}

func TestUnmarshal(t *testing.T) {
	var obj Object

	err := json.Unmarshal([]byte(`
{  
    "employees": [
		{  
			"name":      "octo",   
			"salary":     56000,   
			"married":    false  
		}  
	] 
}
`), &obj)

	testingx.Expect(t, err, testingx.Be[error](nil))

	testingx.Expect(t, obj.Value(), testingx.Equal[any](Obj{
		"employees": List{
			Obj{
				"name":    "octo",
				"salary":  56000,
				"married": false,
			},
		},
	}))
}
