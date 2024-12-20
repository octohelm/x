package testing

import (
	"bytes"
	"fmt"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/octohelm/x/anyjson"
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

func NewSnapshot() *Snapshot {
	return &Snapshot{}
}

type Snapshot txtar.Archive

func (s Snapshot) With(file string, data []byte) *Snapshot {
	s.Files = append(s.Files, txtar.File{
		Name: file,
		Data: data,
	})

	return &s
}

func MustAsJSON(v any) []byte {
	raw, err := AsJSON(v)
	if err != nil {
		panic(err)
	}
	return raw
}

func AsJSON(v any) ([]byte, error) {
	vv, err := anyjson.FromValue(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(anyjson.Sorted(vv), jsontext.WithIndent("  "))
}

func MatchSnapshot(name string, optionFuncs ...func(m *snapshotMatcher)) Matcher[*Snapshot] {
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

func (s *snapshotMatcher) FormatActual(a *Snapshot) string {
	return string(txtar.Format((*txtar.Archive)(a)))
}

func (s *snapshotMatcher) Match(a *Snapshot) bool {
	data := txtar.Format((*txtar.Archive)(a))
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
