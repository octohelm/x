package snapshot

import (
	"bytes"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/octohelm/x/testing/internal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/txtar"
)

var updateSnapshots string

func init() {
	updateSnapshots = os.Getenv("UPDATE_SNAPSHOTS")
}

func NewSnapshot() *Snapshot {
	return &Snapshot{}
}

func Load(name string) *Snapshot {
	// testdata/__snapshots__/<name>.txtar

	filename := path.Join(
		"testdata",
		"__snapshots__",
		strings.ReplaceAll(cases.Lower(language.Und).String(name), " ", "_")+".txtar",
	)

	if strings.ToUpper(updateSnapshots) == "ALL" || strings.Contains(updateSnapshots, name) {
		return &Snapshot{
			filename: filename,
		}
	}

	snapshot, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Snapshot{
				filename: filename,
			}
		}
		panic(err)
	}

	t := txtar.Parse(snapshot)

	return &Snapshot{
		filename: filename,
		files:    t.Files,
	}
}

type File = txtar.File

func FileFromRaw(filename string, data []byte) File {
	return File{
		Name: filename,
		Data: data,
	}
}

func Files(files ...File) *Snapshot {
	return &Snapshot{files: files}
}

type Snapshot struct {
	filename string
	files    []txtar.File
}

func (s *Snapshot) Lines() internal.Lines {
	return internal.LinesFromBytes(txtar.Format(&txtar.Archive{
		Files: s.files,
	}))
}

func (s *Snapshot) Bytes() []byte {
	return txtar.Format(&txtar.Archive{
		Files: s.files,
	})
}

func (s *Snapshot) Add(file string, data []byte) {
	s.files = append(s.files, txtar.File{
		Name: file,
		Data: data,
	})
}

func (s Snapshot) With(file string, data []byte) *Snapshot {
	s.files = append(s.files, txtar.File{
		Name: file,
		Data: data,
	})
	return &s
}

func (s *Snapshot) Equal(a *Snapshot) bool {
	a.filename = s.filename
	if len(s.files) == 0 {
		return true
	}
	return bytes.Equal(s.Bytes(), a.Bytes())
}

func (s *Snapshot) PostMatched() error {
	if err := os.MkdirAll(filepath.Dir(s.filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(s.filename, s.Bytes(), 0o644)
}
