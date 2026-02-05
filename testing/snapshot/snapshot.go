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

func NewSnapshot() *Snapshot {
	return &Snapshot{}
}

func FromFiles(files ...*File) *Snapshot {
	return &Snapshot{files: files}
}

type Snapshot struct {
	ctx   *Context
	files []*File
}

func (s *Snapshot) IsZero() bool {
	return s == nil || len(s.files) == 0
}

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

func (s *Snapshot) Files() iter.Seq[*File] {
	return slices.Values(s.files)
}

func (s *Snapshot) Lines() lines.Lines {
	return lines.FromBytes(s.Bytes())
}

func (s *Snapshot) Add(file string, data []byte) {
	s.files = append(s.files, &File{
		Name: file,
		Data: data,
	})
}

func (s Snapshot) With(file string, data []byte) *Snapshot {
	s.files = append(s.files, &File{
		Name: file,
		Data: data,
	})
	return &s
}

func (s *Snapshot) Commit(ctx *Context) error {
	if err := os.MkdirAll(filepath.Dir(ctx.Filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(ctx.Filename, s.Bytes(), 0o644)
}

func (s *Snapshot) Equal(a *Snapshot) bool {
	a.ctx = s.ctx
	if len(s.files) == 0 {
		return true
	}
	return bytes.Equal(s.Bytes(), a.Bytes())
}

func (s *Snapshot) PostMatched() error {
	return s.Commit(s.ctx)
}
