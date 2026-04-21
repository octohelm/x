package v2

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
)

type fakeTB struct {
	errors  []any
	fatals  []any
	skipped bool
}

func (t *fakeTB) Context() context.Context { return context.Background() }
func (t *fakeTB) TempDir() string {
	dir, _ := os.MkdirTemp("", "v2-test-*")
	return dir
}
func (t *fakeTB) Output() io.Writer        { return io.Discard }
func (t *fakeTB) Chdir(dir string)         {}
func (t *fakeTB) Setenv(key, value string) {}
func (t *fakeTB) Cleanup(fn func())        { fn() }
func (t *fakeTB) Helper()                  {}
func (t *fakeTB) Skip(args ...any)         { t.skipped = true }
func (t *fakeTB) Skipped() bool            { return t.skipped }
func (t *fakeTB) Error(args ...any)        { t.errors = append(t.errors, args...) }
func (t *fakeTB) Fatal(args ...any)        { t.fatals = append(t.fatals, args...) }

type helperRecorder struct{ skipN int }

func (r *helperRecorder) skip(n int) { r.skipN = n }

func TestHelperAndPrepareFailures(t *testing.T) {
	t.Run("Helper 返回原对象", func(t *testing.T) {
		rec := &helperRecorder{}
		got := Helper(rec)

		Then(t, "Helper 应原样返回输入对象",
			Expect(got == rec, Equal(true)),
		)
	})

	t.Run("Must 系列在失败时走 Fatal", func(t *testing.T) {
		tb := &fakeTB{}

		Must(tb, func() error { return errors.New("must failed") })
		v := MustValue(tb, func() (int, error) { return 0, errors.New("must value failed") })
		a, b := MustValues(tb, func() (int, string, error) { return 0, "", errors.New("must values failed") })

		Then(t, "失败时应记录 fatal 且返回零值",
			Expect(len(tb.fatals), Equal(3)),
			Expect(v, Equal(0)),
			Expect(a, Equal(0)),
			Expect(b, Equal("")),
		)
	})
}

func TestInternalCheckers(t *testing.T) {
	t.Run("beChecker 失败时走 Fatal", func(t *testing.T) {
		tb := &fakeTB{}

		Be(func(v int) error { return errors.New("be failed") }).Check(tb, 1)

		Then(t, "失败应记录 fatal",
			Expect(len(tb.fatals), Equal(1)),
		)
	})

	t.Run("Expect 无 checker 时为空操作", func(t *testing.T) {
		tb := &fakeTB{}

		Expect(1).Check(tb)

		Then(t, "不应记录 error 或 fatal",
			Expect(len(tb.errors), Equal(0)),
			Expect(len(tb.fatals), Equal(0)),
		)
	})

	t.Run("ExpectMustValue do 失败时走 Fatal", func(t *testing.T) {
		tb := &fakeTB{}

		ExpectMustValue(func() (int, error) {
			return 0, errors.New("boom")
		}, Equal(1)).Check(tb)

		Then(t, "失败应记录 fatal",
			Expect(len(tb.fatals), Equal(1)),
		)
	})
}

func TestErrorHelpersAndMessages(t *testing.T) {
	t.Run("ErrNotEqual 与 ErrEqual 可输出错误文本", func(t *testing.T) {
		notEqual := (&ErrNotEqual{Expect: 1, Got: 2}).Error()
		equal := (&ErrEqual{NotExpect: 1, Got: 1}).Error()

		Then(t, "错误文本应非空",
			Expect(notEqual == "", Equal(false)),
			Expect(equal == "", Equal(false)),
		)
	})

	t.Run("ErrorNotAsType 覆盖正反路径", func(t *testing.T) {
		tb1 := &fakeTB{}
		tb2 := &fakeTB{}

		ErrorNotAsType[*os.PathError]().Check(tb1, errors.New("plain"))
		ErrorNotAsType[*os.PathError]().Check(tb2, &os.PathError{Op: "open", Path: "x", Err: errors.New("bad")})

		Then(t, "命中时应 fatal，未命中时通过",
			Expect(len(tb1.fatals), Equal(0)),
			Expect(len(tb2.fatals), Equal(1)),
		)
	})
}
