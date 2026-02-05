package snapshot

import (
	"fmt"
)

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
