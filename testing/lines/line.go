package lines

import (
	"slices"
	"strings"
)

// Lines 表示按行拆分后的文本。
type Lines []string

// Differ 表示可导出为 Lines 的对象。
type Differ interface {
	Lines() Lines
}

// FromBytes 将字节切片按行拆分为 Lines。
func FromBytes(data []byte) Lines {
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
