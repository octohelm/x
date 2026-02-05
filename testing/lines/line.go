package lines

import (
	"slices"
	"strings"
)

type Lines []string

type Differ interface {
	Lines() Lines
}

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
