package snapshot

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
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

type Option = func(m *snapshotMatcher)

func WithSnapshotUpdate() Option {
	return func(m *snapshotMatcher) {
		m.update = true
	}
}

func WithWorkDir(wd string) Option {
	return func(m *snapshotMatcher) {
		m.wd = wd
	}
}

func NewSnapshot() *Snapshot {
	return &Snapshot{}
}

type Snapshot txtar.Archive

func (s *Snapshot) Add(file string, data []byte) {
	s.Files = append(s.Files, txtar.File{
		Name: file,
		Data: data,
	})
}

func (s Snapshot) With(file string, data []byte) *Snapshot {
	s.Files = append(s.Files, txtar.File{
		Name: file,
		Data: data,
	})
	return &s
}

func Match(name string, options ...Option) internal.Matcher[*Snapshot] {
	// testdata/__snapshots__/<name>.txtar

	snapshotFilename := fmt.Sprintf("testdata/__snapshots__/%s.txtar",
		strings.ReplaceAll(cases.Lower(language.Und).String(name), " ", "_"),
	)

	snapshot, _ := os.ReadFile(snapshotFilename)

	m := &snapshotMatcher{
		filename: snapshotFilename,
		expected: snapshot,
		update:   strings.ToUpper(updateSnapshots) == "ALL" || strings.Contains(updateSnapshots, name),
	}

	for _, fn := range options {
		fn(m)
	}

	return m
}

type snapshotMatcher struct {
	wd       string
	filename string
	expected []byte
	update   bool
}

func (snapshotMatcher) Action() string {
	return "match snapshot"
}

func (s *snapshotMatcher) Negative() bool {
	return false
}

var _ internal.MatcherWithNormalizedExpected = &snapshotMatcher{}

func (s *snapshotMatcher) Match(a *Snapshot) bool {
	data := txtar.Format((*txtar.Archive)(a))
	if s.update || len(s.expected) == 0 {
		if err := s.commitSnapshots(data); err != nil {
			panic(err)
		}
		return true
	}
	return bytes.Equal(data, s.expected)
}

func (s *snapshotMatcher) commitSnapshots(data []byte) error {
	filename := s.filename
	if s.wd != "" {
		filename = path.Join(s.wd, filename)
	}
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0o644)
}

func (s *snapshotMatcher) NormalizedExpected() any {
	return LinesFromBytes(s.expected)
}

func (s *snapshotMatcher) NormalizeActual(a *Snapshot) any {
	return LinesFromBytes(txtar.Format((*txtar.Archive)(a)))
}

func LinesFromBytes(data []byte) Lines {
	return slices.Collect(func(yield func(line string) bool) {
		for line := range strings.Lines(string(data)) {
			if len(line) > 0 {
				if line[len(line)-1] == '\n' {
					line = line[:len(line)-1]
				}
			}

			if !yield(line) {
				return
			}
		}
	})
}

type Lines []string
