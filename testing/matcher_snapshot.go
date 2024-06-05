package testing

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/txtar"
)

var updateSnapshots string

func init() {
	updateSnapshots = os.Getenv("UPDATE_SNAPSHOTS")
}

func WithSnapshotUpdate() func(m *snapshotMatcher) {
	return func(m *snapshotMatcher) {
		m.update = true
	}
}

func WithWorkDir(wd string) func(m *snapshotMatcher) {
	return func(m *snapshotMatcher) {
		m.wd = wd
	}
}

func MatchSnapshot(name string, optionFuncs ...func(m *snapshotMatcher)) Matcher[*txtar.Archive] {
	// testdata/__snapshots__/<name>.txtar

	snapshotFilename := fmt.Sprintf("testdata/__snapshots__/%s.txtar", name)

	snapshot, _ := os.ReadFile(snapshotFilename)

	m := &snapshotMatcher{
		filename: snapshotFilename,
		expected: snapshot,
		update:   updateSnapshots == "ALL" || strings.Contains(updateSnapshots, name),
	}

	for _, fn := range optionFuncs {
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

func (s *snapshotMatcher) Name() string {
	return "MatchSnapshot"
}

func (s *snapshotMatcher) Negative() bool {
	return false
}

func (s *snapshotMatcher) FormatExpected() string {
	return string(s.expected)
}

func (s *snapshotMatcher) FormatActual(a *txtar.Archive) string {
	return string(txtar.Format(a))
}

func (s *snapshotMatcher) Match(a *txtar.Archive) bool {
	data := txtar.Format(a)
	if s.update || len(s.expected) == 0 {
		_ = s.commitSnapshots(data)

		return true
	}
	return bytes.Equal(data, s.expected)
}

func (s *snapshotMatcher) commitSnapshots(data []byte) error {
	filename := s.filename
	if s.wd != "" {
		filename = filepath.Join(s.wd, filename)
	}
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0o644)
}

var _ ExpectedFormatter = &snapshotMatcher{}
