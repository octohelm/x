package internal

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/octohelm/x/testing/lines"
)

func FormatErrorMessage(expectOrNot any, actual any, neg bool) string {
	switch x := expectOrNot.(type) {
	case *regexp.Regexp:
		if neg {
			return fmt.Sprintf(`
	not match: %v
	got:       %v
`, expectOrNot, actual)
		}
		return fmt.Sprintf(`
match: %v
got:   %v
`, expectOrNot, actual)
	case string:
		if strings.Index(x, "\n") > -1 {
			if actualStr, ok := actual.(string); ok {
				return fmt.Sprintf(`
not match:
%s
`, lines.Diff(lines.FromBytes([]byte(x)), lines.FromBytes([]byte(actualStr))))
			}
		}
	case []byte:
		if bytes.Index(x, []byte("\n")) > -1 {
			if actualBytes, ok := actual.([]byte); ok {
				return fmt.Sprintf(`
not match:
%s
`, lines.Diff(lines.FromBytes(x), lines.FromBytes(actualBytes)))
			}
		}
	}

	r := &diffReporter{
		neg: neg,
		ret: &bytes.Buffer{},
	}
	_ = cmp.Diff(
		root{Value: expectOrNot},
		root{Value: actual},
		cmp.Reporter(r),
	)
	return r.ret.String()
}

type root struct {
	Value any
}

const prefix = "{internal.root}.Value."

type diffReporter struct {
	neg    bool
	ret    *bytes.Buffer
	path   cmp.Path
	lx, ly int
}

func (r *diffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *diffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *diffReporter) Report(rs cmp.Result) {
	if r.neg {
		if rs.Equal() {
			vx, vy := r.path.Last().Values()
			_, _ = fmt.Fprintf(r.ret, `
%s: 
	not expect: %+v
	got:        %+v
`, r.path.GoString()[len(prefix):], vx, vy)
			return
		}
	}

	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		_, _ = fmt.Fprintf(r.ret, `
%s: 
	expect: %+v
	got:    %+v
`, r.path.GoString()[len(prefix):], vx, vy)
	}
}
