package ptr_test

import (
	"testing"

	"github.com/octohelm/x/ptr"
	. "github.com/octohelm/x/testing/v2"
)

func TestPtr(t *testing.T) {
	t.Run("返回值副本的指针", func(t *testing.T) {
		value := 1
		p := ptr.Ptr(value)
		value = 2

		Then(t, "指针中的值不应受原变量后续修改影响",
			Expect(*p, Equal(1)),
		)
	})

	t.Run("支持结构体值", func(t *testing.T) {
		type user struct {
			Name string
		}

		p := ptr.Ptr(user{Name: "octo"})

		Then(t, "应返回对应类型的指针",
			Expect(p.Name, Equal("octo")),
		)
	})
}
