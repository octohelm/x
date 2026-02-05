package snapshot

import (
	"bytes"
	"fmt"

	"github.com/octohelm/x/testing/lines"
)

func Diff(src *Snapshot, dst *Snapshot) ([]byte, bool) {
	srcFiles := map[string]*File{}
	for _, f := range src.files {
		srcFiles[f.Name] = f
	}

	buf := bytes.NewBuffer(nil)

	for _, f := range dst.files {
		if srcFile, ok := srcFiles[f.Name]; ok {
			diffed := lines.Diff(lines.FromBytes(srcFile.Data), lines.FromBytes(f.Data))

			if len(diffed) > 0 {
				_, _ = fmt.Fprintf(buf, "M -- %s --\n", f.Name)
				_, _ = buf.Write(diffed)
			}

			delete(srcFiles, f.Name)
		}
	}

	if len(srcFiles) > 0 {
		for _, f := range srcFiles {
			_, _ = fmt.Fprintf(buf, "D -- %s --\n", f.Name)
		}
	}

	return buf.Bytes(), buf.Len() > 0
}
