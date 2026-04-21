package v2

import (
	"fmt"
	"go/token"
	"iter"
	"slices"

	"github.com/octohelm/x/testing/internal"
	"github.com/octohelm/x/testing/snapshot"
)

// SnapshotFileFromRaw 根据文件名和原始内容创建快照文件。
func SnapshotFileFromRaw(filename string, raw []byte) *SnapshotFile {
	return snapshot.FileFromRaw(filename, raw)
}

// Snapshot 表示一组待校验的快照文件序列。
type Snapshot = iter.Seq[*SnapshotFile]

// SnapshotOf 将若干快照文件组装为 Snapshot。
func SnapshotOf(files ...*SnapshotFile) Snapshot {
	return func(yield func(*SnapshotFile) bool) {
		for _, file := range files {
			if !yield(file) {
				return
			}
		}
	}
}

// SnapshotFile 是快照文件的别名。
type SnapshotFile = snapshot.File

// ErrNotMatch 是快照不匹配错误的别名。
type ErrNotMatch = snapshot.ErrNotMatch

// MatchSnapshot 返回基于名称匹配测试快照的检查器。
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
