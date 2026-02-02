package snapshot

import (
	"github.com/octohelm/x/testing/internal"
)

type Option = func(m *snapshotMatcher)

func WithWorkDir(wd string) Option {
	return func(m *snapshotMatcher) {
		m.wd = wd
	}
}

func Match(name string, options ...Option) internal.Matcher[*Snapshot] {
	m := &snapshotMatcher{
		expected: Load(name),
	}

	for _, fn := range options {
		fn(m)
	}

	return m
}

type snapshotMatcher struct {
	wd       string
	expected *Snapshot
}

func (snapshotMatcher) Action() string {
	return "match snapshot"
}

func (s *snapshotMatcher) Negative() bool {
	return false
}

func (s *snapshotMatcher) Match(a *Snapshot) bool {
	return s.expected.Equal(a)
}
