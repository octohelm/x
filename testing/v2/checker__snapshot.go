package v2

import (
	"fmt"
	"go/token"
	"iter"
	"slices"

	"github.com/octohelm/x/testing/internal"
	"github.com/octohelm/x/testing/snapshot"
)

func SnapshotFileFromRaw(filename string, raw []byte) *SnapshotFile {
	return snapshot.FileFromRaw(filename, raw)
}

type Snapshot = iter.Seq[*SnapshotFile]

func SnapshotOf(files ...*SnapshotFile) Snapshot {
	return func(yield func(*SnapshotFile) bool) {
		for _, file := range files {
			if !yield(file) {
				return
			}
		}
	}
}

type SnapshotFile = snapshot.File

type ErrNotMatch = snapshot.ErrNotMatch

func MatchSnapshot[Name comparable](name Name) ValueChecker[Snapshot] {
	return internal.Helper(1, &snapshotMatcher{
		ctx: &snapshot.Context{
			Name: fmt.Sprintf("%v", name),
		},
	})
}

type snapshotMatcher struct {
	Reporter

	ctx      *snapshot.Context
	snapshot *snapshot.Snapshot
	pos      token.Position
}

func (r *snapshotMatcher) Check(t TB, actual Snapshot) {
	t.Helper()

	if s, err := r.ctx.Load(); err != nil {
		r.Fatal(t, err)
	} else {
		r.snapshot = s
	}

	a := snapshot.FromFiles(slices.Collect(actual)...)

	if r.snapshot.IsZero() {
		if err := a.Commit(r.ctx); err != nil {
			r.Error(t, err)
			return
		}
		return
	}

	if raw, changed := snapshot.Diff(r.snapshot, a); changed {
		r.Fatal(t, &snapshot.ErrNotMatch{
			Name:   r.ctx.Name,
			Diffed: raw,
		})
	}

	if err := a.Commit(r.ctx); err != nil {
		r.Error(t, err)
		return
	}
}
