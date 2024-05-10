package testing

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/txtar"
)

func MatchSnapshot(name string) Matcher[*txtar.Archive] {
	// testdata/__snapshots__/<name>.txtar

	snapshotFilename := fmt.Sprintf("testdata/__snapshots__/%s.txtar", name)

	snapshot, _ := os.ReadFile(snapshotFilename)

	return &snapshotMatcher{
		filename: snapshotFilename,
		expected: snapshot,
	}
}

type snapshotMatcher struct {
	filename string
	expected []byte
}

func (s snapshotMatcher) Name() string {
	return "MatchSnapshot"
}

func (s snapshotMatcher) Negative() bool {
	return false
}

func (s snapshotMatcher) FormatExpected() string {
	return string(s.expected)
}

func (s snapshotMatcher) FormatActual(a *txtar.Archive) string {
	return string(txtar.Format(a))
}

func (s snapshotMatcher) Match(a *txtar.Archive) bool {
	data := txtar.Format(a)
	if len(s.expected) == 0 {
		_ = os.MkdirAll(filepath.Dir(s.filename), os.ModePerm)
		_ = os.WriteFile(s.filename, data, 0o644)
		return true
	}
	return bytes.Equal(data, s.expected)
}

var _ ExpectedFormatter = &snapshotMatcher{}
