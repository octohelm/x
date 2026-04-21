package snapshot_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/octohelm/x/testing/lines"
	"github.com/octohelm/x/testing/snapshot"
	. "github.com/octohelm/x/testing/v2"
)

func TestFilesSeqAndFileFromRaw(t *testing.T) {
	t.Run("按 txtar 风格拆分多个文件", func(t *testing.T) {
		raw := []byte("-- a.txt --\nhello\r\n-- b.txt --\nworld")

		files := slices.Collect(snapshot.FilesSeq(raw))

		Then(t, "应按顺序解析文件名和内容",
			Expect(len(files), Equal(2)),
			Expect(files[0].Name, Equal("a.txt")),
			Expect(string(files[0].Data), Equal("hello\r\n")),
			Expect(files[1].Name, Equal("b.txt")),
			Expect(string(files[1].Data), Equal("world\n")),
		)
	})

	t.Run("FileFromRaw 创建单个文件对象", func(t *testing.T) {
		f := snapshot.FileFromRaw("x.txt", []byte("123"))

		Then(t, "应保留原始文件名与内容",
			Expect(f.Name, Equal("x.txt")),
			Expect(string(f.Data), Equal("123")),
		)
	})
}

func TestSnapshotBuilders(t *testing.T) {
	t.Run("NewSnapshot 与 Add", func(t *testing.T) {
		s := snapshot.NewSnapshot()
		s.Add("a.txt", []byte("hello"))

		Then(t, "新增文件后不再是零值快照",
			Expect(s.IsZero(), Equal(false)),
			Expect(string(s.Bytes()), Equal("-- a.txt --\nhello\n")),
			Expect(s.Lines(), Equal(lines.Lines{"-- a.txt --", "hello"})),
		)
	})

	t.Run("FromFiles 与 With", func(t *testing.T) {
		base := snapshot.FromFiles(
			snapshot.FileFromRaw("a.txt", []byte("1")),
		)
		extended := base.With("b.txt", []byte("2"))

		Then(t, "With 应返回附加文件后的副本",
			Expect(slices.Collect(base.Files())[0].Name, Equal("a.txt")),
			Expect(len(slices.Collect(base.Files())), Equal(1)),
			Expect(len(slices.Collect(extended.Files())), Equal(2)),
			Expect(string(extended.Bytes()), Equal("-- a.txt --\n1\n-- b.txt --\n2\n")),
		)
	})

	t.Run("空快照 Equal 行为", func(t *testing.T) {
		var zero snapshot.Snapshot
		other := snapshot.FromFiles(snapshot.FileFromRaw("a.txt", []byte("1")))

		Then(t, "零值快照当前会视作匹配任意快照",
			Expect(zero.Equal(other), Equal(true)),
		)
	})
}

func TestSnapshotDiff(t *testing.T) {
	t.Run("修改与删除文件生成 diff", func(t *testing.T) {
		src := snapshot.FromFiles(
			snapshot.FileFromRaw("a.txt", []byte("hello")),
			snapshot.FileFromRaw("b.txt", []byte("obsolete")),
		)
		dst := snapshot.FromFiles(
			snapshot.FileFromRaw("a.txt", []byte("world")),
		)

		raw, changed := snapshot.Diff(src, dst)

		Then(t, "应报告修改与删除",
			Expect(changed, Equal(true)),
			Expect(strings.Contains(string(raw), "M -- a.txt --"), Equal(true)),
			Expect(strings.Contains(string(raw), "D -- b.txt --"), Equal(true)),
		)
	})

	t.Run("完全相同时不生成 diff", func(t *testing.T) {
		src := snapshot.FromFiles(snapshot.FileFromRaw("a.txt", []byte("same")))
		dst := snapshot.FromFiles(snapshot.FileFromRaw("a.txt", []byte("same")))

		raw, changed := snapshot.Diff(src, dst)

		Then(t, "相同内容应返回空 diff",
			Expect(changed, Equal(false)),
			Expect(string(raw), Equal("")),
		)
	})
}

func TestContextLoadCommitAndPostMatched(t *testing.T) {
	t.Run("Load 在文件缺失时返回空快照并设置路径", func(t *testing.T) {
		wd := t.TempDir()
		t.Chdir(wd)

		ctx := &snapshot.Context{Name: "Case A"}
		s := MustValue(t, ctx.Load)

		Then(t, "缺失文件时返回空快照并设置标准路径",
			Expect(s.IsZero(), Equal(true)),
			Expect(ctx.Filename, Equal(filepath.Join("testdata", "__snapshots__", "case_a.txtar"))),
		)
	})

	t.Run("Commit 后 Load 可读回内容", func(t *testing.T) {
		wd := t.TempDir()
		t.Chdir(wd)

		ctx := &snapshot.Context{Name: "Case B", Filename: filepath.Join("testdata", "__snapshots__", "case_b.txtar")}
		s := snapshot.FromFiles(
			snapshot.FileFromRaw("a.txt", []byte("hello")),
			snapshot.FileFromRaw("b.txt", []byte("world")),
		)

		Then(t, "Commit 应成功写入文件",
			ExpectDo(func() error { return s.Commit(ctx) }),
		)

		loaded := MustValue(t, func() (*snapshot.Snapshot, error) {
			readCtx := &snapshot.Context{Name: "Case B"}
			return readCtx.Load()
		})

		Then(t, "Load 应按默认路径读取已提交内容",
			Expect(string(loaded.Bytes()), Equal("-- a.txt --\nhello\n-- b.txt --\nworld\n")),
		)
	})

	t.Run("PostMatched 复用上下文提交快照", func(t *testing.T) {
		wd := t.TempDir()
		t.Chdir(wd)

		ctx := &snapshot.Context{Name: "Case C"}
		s := MustValue(t, ctx.Load)
		s.Add("x.txt", []byte("payload"))

		Then(t, "PostMatched 应把内容写回 ctx 指定位置",
			ExpectDo(func() error { return s.PostMatched() }),
		)

		raw := MustValue(t, func() ([]byte, error) {
			return os.ReadFile(filepath.Join("testdata", "__snapshots__", "case_c.txtar"))
		})

		Then(t, "写回文件应包含新增快照内容",
			Expect(string(raw), Equal("-- x.txt --\npayload\n")),
		)
	})
}

func TestAsJSONAndErrNotMatch(t *testing.T) {
	t.Run("AsJSON 稳定输出格式化 JSON", func(t *testing.T) {
		raw := MustValue(t, func() ([]byte, error) {
			return snapshot.AsJSON(map[string]any{
				"b": 2,
				"a": map[string]any{
					"d": 4,
					"c": 3,
				},
			})
		})

		Then(t, "应按键排序并带缩进输出",
			Expect(string(raw), Equal("{\n  \"a\": {\n    \"c\": 3,\n    \"d\": 4\n  },\n  \"b\": 2\n}")),
		)
	})

	t.Run("ErrNotMatch 输出更新提示", func(t *testing.T) {
		err := (&snapshot.ErrNotMatch{Name: "demo", Diffed: []byte("DIFF")}).Error()

		Then(t, "错误文本应包含名称与更新提示",
			Expect(strings.Contains(err, "Snapshot(demo) failed"), Equal(true)),
			Expect(strings.Contains(err, "UPDATE_SNAPSHOTS=demo go test"), Equal(true)),
			Expect(strings.Contains(err, "DIFF"), Equal(true)),
		)
	})
}
