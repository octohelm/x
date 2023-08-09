package testing_test

import (
	"strings"
	"testing"

	"slices"

	. "github.com/octohelm/x/testing"
)

func Test(t *testing.T) {
	var (
		ContainsStringItem = MatcherWith(slices.Contains[[]string, string], "ContainsStringItem")
		HaveStringSuffix   = MatcherWith(strings.HasSuffix, "HaveStringSuffix")
	)

	t.Run("Matchers", func(t *testing.T) {
		t.Run("Should check", func(t *testing.T) {
			Expect(t, "1",
				Equal("1"),
				Not(Equal("2")),
			)
		})

		t.Run("Should equal", func(t *testing.T) {
			Expect(t, map[string]string{"1": "1"},
				Equal(map[string]string{"1": "1"}),
			)
		})

		t.Run("Should Be", func(t *testing.T) {
			Expect(t, error(nil),
				Be[error](nil),
			)
		})

		t.Run("Should Contains and Have Len", func(t *testing.T) {
			Expect(t, "x1x",
				HaveStringSuffix("x"),
				HaveLen[string](3),
			)
		})

		t.Run("Should Have HaveLen", func(t *testing.T) {
			Expect(t, []string{"1", "2"},
				HaveLen[[]string](2),
			)
		})

		t.Run("Should Have slices.Contains", func(t *testing.T) {
			Expect(t, []string{"1", "2"},
				ContainsStringItem("2"),
			)
		})
	})

	t.Run("Should get project root", func(t *testing.T) {
		pr := ProjectRoot()
		Expect(t, pr,
			HaveStringSuffix("/x"),
		)
	})
}
