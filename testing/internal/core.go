package internal

import (
	"context"
	"fmt"
	"go/token"
	"io"
	"path/filepath"
	"runtime"
	"strings"
)

type TB interface {
	Context() context.Context

	TempDir() string
	Output() io.Writer
	Chdir(dir string)
	Setenv(key, value string)

	Cleanup(func())

	Helper()
	Skip(args ...any)
}

type TBController interface {
	Skipped() bool
	Error(args ...any)
	Fatal(args ...any)
}

type Reporter struct {
	pos token.Position
}

func Helper[X any](skip int, x X) X {
	skip = skip + 1

	m := any(x)

	if v, ok := m.(interface{ skip(int) }); ok {
		v.skip(skip)
	}

	return x
}

func (r *Reporter) skip(skip int) {
	_, file, line, _ := runtime.Caller(skip + 1)

	r.pos.Filename = filepath.Base(file)
	r.pos.Line = line
}

func (r *Reporter) Error(t TB, err error) {
	if err == nil {
		return
	}

	t.Helper()
	if tbc, ok := t.(TBController); ok {
		tbc.Error(&errWithPos{err: err, pos: r.pos})
	}
}

func (r *Reporter) Fatal(t TB, err error) {
	if err == nil {
		return
	}

	t.Helper()
	if tbc, ok := t.(TBController); ok {
		tbc.Fatal(&errWithPos{err: err, pos: r.pos})
	}
}

type errWithPos struct {
	err error
	pos token.Position
}

func (e *errWithPos) Unwrap() error {
	return e.err
}

func (e *errWithPos) Error() string {
	return fmt.Sprintf(`
%s:
%s
`, e.pos, indent(e.err.Error()))
}

func indent(lines string) string {
	s := &strings.Builder{}
	for line := range strings.Lines(strings.TrimSpace(lines)) {
		s.WriteString("    ")
		s.WriteString(line)
	}
	return s.String()
}
