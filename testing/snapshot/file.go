package snapshot

import (
	"bytes"
	"iter"
	"strings"
)

type File = struct {
	Name string // name of file ("foo/bar.txt")
	Data []byte // text content of file
}

func FileFromRaw(filename string, data []byte) *File {
	return &File{
		Name: filename,
		Data: data,
	}
}

func FilesSeq(data []byte) iter.Seq[*File] {
	return func(yield func(*File) bool) {
		var name string
		_, name, data = findFileMarker(data)
		for name != "" {
			f := &File{name, nil}
			f.Data, name, data = findFileMarker(data)
			if !yield(f) {
				return
			}
		}
	}
}

var (
	newlineMarker = []byte("\n-- ")
	marker        = []byte("-- ")
	markerEnd     = []byte(" --")
)

func findFileMarker(data []byte) (before []byte, name string, after []byte) {
	var i int
	for {
		if name, after = isMarker(data[i:]); name != "" {
			return data[:i], name, after
		}
		j := bytes.Index(data[i:], newlineMarker)
		if j < 0 {
			return fixNL(data), "", nil
		}
		i += j + 1 // positioned at start of new possible marker
	}
}

func isMarker(data []byte) (name string, after []byte) {
	if !bytes.HasPrefix(data, marker) {
		return "", nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		data, after = data[:i], data[i+1:]
		if data[i-1] == '\r' { // handle \r\n line ending
			data = data[:i-1]
		}
	}
	if !(bytes.HasSuffix(data, markerEnd) && len(data) >= len(marker)+len(markerEnd)) {
		return "", nil
	}
	return strings.TrimSpace(string(data[len(marker) : len(data)-len(markerEnd)])), after
}
