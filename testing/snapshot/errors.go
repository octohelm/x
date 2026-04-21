package snapshot

import (
	"fmt"
)

// ErrNotMatch 表示实际快照与期望快照不一致。
type ErrNotMatch struct {
	Name   string
	Diffed []byte
}

func (e *ErrNotMatch) Error() string {
	return fmt.Sprintf(`
Snapshot(%s) failed. To update, run: UPDATE_SNAPSHOTS=%s go test
%s
`, e.Name, e.Name, e.Diffed)
}
