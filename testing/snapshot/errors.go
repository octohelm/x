package snapshot

import (
	"fmt"
)

// ErrNotMatch 表示实际快照与期望快照不一致。
type ErrNotMatch struct {
	// Name 是快照名。
	Name string
	// Diffed 是期望与实际之间的差异文本。
	Diffed []byte
}

// Error 返回用于测试失败输出的可读错误信息。
func (e *ErrNotMatch) Error() string {
	return fmt.Sprintf(`
Snapshot(%s) failed. To update, run: UPDATE_SNAPSHOTS=%s go test
%s
`, e.Name, e.Name, e.Diffed)
}
