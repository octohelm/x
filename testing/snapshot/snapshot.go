package snapshot

import (
	"bytes"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"slices"

	"github.com/octohelm/x/testing/lines"
)

var updateSnapshots string

func init() {
	updateSnapshots = os.Getenv("UPDATE_SNAPSHOTS")
}

// NewSnapshot 创建一个空快照。
func NewSnapshot() *Snapshot {
	return &Snapshot{}
}

// FromFiles 通过一组文件构造快照。
func FromFiles(files ...*File) *Snapshot {
	return &Snapshot{files: files}
}

// Snapshot 表示一组可比较、可提交的快照文件。
type Snapshot struct {
	ctx   *Context
	files []*File
}

// IsZero 判断快照是否为空。
func (s *Snapshot) IsZero() bool {
	return s == nil || len(s.files) == 0
}

// Bytes 将快照编码为 txtar 风格文本。
func (s *Snapshot) Bytes() []byte {
	var buf bytes.Buffer
	for _, f := range s.files {
		_, _ = fmt.Fprintf(&buf, "-- %s --\n", f.Name)
		buf.Write(fixNL(f.Data))
	}
	return buf.Bytes()
}

func fixNL(data []byte) []byte {
	if len(data) == 0 || data[len(data)-1] == '\n' {
		return data
	}
	d := make([]byte, len(data)+1)
	copy(d, data)
	d[len(data)] = '\n'
	return d
}

// Files 返回快照文件的迭代序列。
func (s *Snapshot) Files() iter.Seq[*File] {
	return slices.Values(s.files)
}

// Lines 返回快照整体文本的按行表示。
func (s *Snapshot) Lines() lines.Lines {
	return lines.FromBytes(s.Bytes())
}

// Add 向快照追加一个文件。
func (s *Snapshot) Add(file string, data []byte) {
	s.files = append(s.files, &File{
		Name: file,
		Data: data,
	})
}

// With 返回附加了新文件的快照副本。
func (s Snapshot) With(file string, data []byte) *Snapshot {
	s.files = append(s.files, &File{
		Name: file,
		Data: data,
	})
	return &s
}

// Commit 将快照写回 Context 指定位置。
func (s *Snapshot) Commit(ctx *Context) error {
	if err := os.MkdirAll(filepath.Dir(ctx.Filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(ctx.Filename, s.Bytes(), 0o644)
}

// Equal 判断两个快照编码后的内容是否一致。
func (s *Snapshot) Equal(a *Snapshot) bool {
	a.ctx = s.ctx
	if len(s.files) == 0 {
		return true
	}
	return bytes.Equal(s.Bytes(), a.Bytes())
}

// PostMatched 在匹配完成后持久化快照。
func (s *Snapshot) PostMatched() error {
	return s.Commit(s.ctx)
}
