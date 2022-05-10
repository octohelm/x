package testing_test

import (
	"strings"
	"testing"

	"golang.org/x/exp/slices"

	. "github.com/go-courier/x/testing"
)

func Test(t *testing.T) {
	t.Run("Matchers", func(t *testing.T) {
		t.Run("Should check", func(t *testing.T) {
			Expect(t, "1",
				Should(Equal[string], "2"),
				ShouldNot(Equal[string], "2"),
				ShouldNot(strings.Contains, "x"),
			)
		})

		t.Run("Should equal", func(t *testing.T) {
			Expect(t, map[string]string{"1": "1"},
				Should(Equal[map[string]string], map[string]string{"1": "1"}),
			)
		})

		t.Run("Should Be", func(t *testing.T) {
			Expect(t, error(nil),
				Should(Be[error], nil),
			)
		})

		t.Run("Should Contains and Have Len", func(t *testing.T) {
			Expect(t, "x1x",
				Should(strings.HasPrefix, "x"),
				Should(HaveLen[string], 3),
			)
		})

		t.Run("Should Have HaveLen", func(t *testing.T) {
			Expect(t, []string{"1", "2"},
				Should(HaveLen[[]string], 2),
			)
		})

		t.Run("Should Have slices.Contains", func(t *testing.T) {
			Expect(t, []string{"1", "2"},
				Should(slices.Contains[string], "2"),
			)
		})
	})

	t.Run("Should get project root", func(t *testing.T) {
		pr := ProjectRoot()
		Expect(t, pr,
			Should(strings.HasSuffix, "github.com/go-courier/x"),
		)
	})
}
