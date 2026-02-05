package lines

import (
	"bytes"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

func Diff(oldLines Lines, newLines Lines) []byte {
	r := &diffReporter{
		ret: &bytes.Buffer{},
	}
	_ = cmp.Diff(oldLines, newLines, cmp.Reporter(r))
	return r.ret.Bytes()
}

type diffReporter struct {
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
	last := r.path.Last()

	vx, vy := last.Values()

	switch x := r.path.Last().(type) {
	case cmp.SliceIndex:
		ix, iy := x.SplitKeys()
		if ix > -1 {
			r.lx = ix + 1
		}
		if iy > -1 {
			r.ly = iy + 1
		}

		if rs.Equal() {
			return
		}

		if ix == -1 && iy > -1 {
			_, _ = fmt.Fprintf(r.ret, `
@@ -%d,0 +%d,1 @@
+%v
`, r.lx, r.ly, vy)
		} else if ix > -1 && iy == -1 {
			_, _ = fmt.Fprintf(r.ret, `
@@ -%d,1 +%d,0 @@ 
-%v
`, r.lx, r.ly, vx)
		} else if ix == iy {
			_, _ = fmt.Fprintf(r.ret, `
@@ -%d,1 +%d,1 @@
-%v
+%v
`, r.lx, r.ly, vx, vy)
		}
	}
}
