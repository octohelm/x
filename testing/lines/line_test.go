package lines_test

import (
	"strings"
	"testing"

	"github.com/octohelm/x/testing/lines"
	. "github.com/octohelm/x/testing/v2"
)

func TestDiff(t *testing.T) {
	t.Run("LinesFromBytes", func(t *testing.T) {
		t.Run("WHEN converting multiline data", func(t *testing.T) {
			data := []byte("line1\nline2\n\nline3")
			ls := lines.FromBytes(data)

			Then(t, "it should split into lines and trim trailing newline",
				Expect(ls, Equal(lines.Lines{
					"line1",
					"line2",
					"",
					"line3",
				})),
			)
		})
	})

	t.Run("Diff Logic", func(t *testing.T) {
		t.Run("GIVEN two sets of lines", func(t *testing.T) {
			oldLines := lines.Lines{"A", "B", "C"}

			t.Run("WHEN a line is added", func(t *testing.T) {
				newLines := lines.Lines{"A", "B", "C", "D"}
				diff := lines.Diff(oldLines, newLines)

				Then(t, "output should contain a plus line",
					Expect(string(diff), Equal(`
@@ -3,0 +4,1 @@
+D
`)),
				)
			})

			t.Run("WHEN a line is removed", func(t *testing.T) {
				newLines := lines.Lines{"A", "C"}
				diff := lines.Diff(oldLines, newLines)

				Then(t, "output should contain a minus line",
					Expect(string(diff), Equal(`
@@ -2,1 +1,0 @@ 
-B
`)),
				)
			})

			t.Run("WHEN a line is modified", func(t *testing.T) {
				newLines := lines.Lines{"A", "B2", "C"}
				diff := lines.Diff(oldLines, newLines)

				Then(t, "output should contain both minus and plus lines",
					Expect(string(diff), Equal(`
@@ -2,1 +2,1 @@
-B
+B2
`)),
				)
			})

			t.Run("WHEN multiple changes occur", func(t *testing.T) {
				oldLines := lines.Lines{"keep", "remove", "change"}
				newLines := lines.Lines{"keep", "change2", "add", "add more"}
				diff := lines.Diff(oldLines, newLines)

				Then(t, "it should report all differences in order",
					Expect(strings.TrimSpace(string(diff)),
						Equal(strings.TrimSpace(`
@@ -2,1 +2,1 @@
-remove
+change2

@@ -3,1 +3,1 @@
-change
+add

@@ -3,0 +4,1 @@
+add more
`))),
				)
			})
		})
	})
}
