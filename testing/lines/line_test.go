package lines_test

import (
	"strings"
	"testing"

	"github.com/octohelm/x/testing/lines"
	. "github.com/octohelm/x/testing/v2"
)

func TestFromBytes(t *testing.T) {
	t.Run("空文本返回空行集合", func(t *testing.T) {
		ls := lines.FromBytes([]byte(""))

		Then(t, "空文本应返回 nil 行集合",
			Expect(ls, Equal(lines.Lines(nil))),
		)
	})

	t.Run("仅包含换行时保留空行", func(t *testing.T) {
		ls := lines.FromBytes([]byte("\n"))

		Then(t, "单个换行应保留为空字符串行",
			Expect(ls, Equal(lines.Lines{""})),
		)
	})

	t.Run("尾换行不额外产生空行", func(t *testing.T) {
		ls := lines.FromBytes([]byte("line1\nline2\n"))

		Then(t, "尾换行应只作为行结束符处理",
			Expect(ls, Equal(lines.Lines{
				"line1",
				"line2",
			})),
		)
	})

	t.Run("中间空行应被保留", func(t *testing.T) {
		ls := lines.FromBytes([]byte("line1\nline2\n\nline3"))

		Then(t, "应按顺序保留中间空行",
			Expect(ls, Equal(lines.Lines{
				"line1",
				"line2",
				"",
				"line3",
			})),
		)
	})
}

func TestDiff(t *testing.T) {
	t.Run("新增行时输出插入 diff", func(t *testing.T) {
		diff := lines.Diff(lines.Lines{"A", "B", "C"}, lines.Lines{"A", "B", "C", "D"})

		Then(t, "应输出新增行 patch",
			Expect(string(diff), Equal(`
@@ -3,0 +4,1 @@
+D
`)),
		)
	})

	t.Run("删除行时输出删除 diff", func(t *testing.T) {
		diff := lines.Diff(lines.Lines{"A", "B", "C"}, lines.Lines{"A", "C"})

		Then(t, "应输出删除行 patch",
			Expect(string(diff), Equal(`
@@ -2,1 +1,0 @@ 
-B
`)),
		)
	})

	t.Run("修改行时输出替换 diff", func(t *testing.T) {
		diff := lines.Diff(lines.Lines{"A", "B", "C"}, lines.Lines{"A", "B2", "C"})

		Then(t, "应输出删除与新增成对 patch",
			Expect(string(diff), Equal(`
@@ -2,1 +2,1 @@
-B
+B2
`)),
		)
	})

	t.Run("连续变更时按顺序输出 patch", func(t *testing.T) {
		diff := lines.Diff(
			lines.Lines{"keep", "remove", "change"},
			lines.Lines{"keep", "change2", "add", "add more"},
		)

		Then(t, "应稳定报告连续变更",
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

	t.Run("相同输入多次 diff 输出保持稳定", func(t *testing.T) {
		oldLines := lines.Lines{"A", "B", "C"}
		newLines := lines.Lines{"A", "B2", "C", "D"}

		first := string(lines.Diff(oldLines, newLines))
		second := string(lines.Diff(oldLines, newLines))

		Then(t, "重复计算应得到相同结果",
			Expect(second, Equal(first)),
		)
	})

	t.Run("相同文本不产生 diff", func(t *testing.T) {
		diff := lines.Diff(lines.Lines{"same"}, lines.Lines{"same"})

		Then(t, "无差异时应返回空文本",
			Expect(string(diff), Equal("")),
		)
	})
}
