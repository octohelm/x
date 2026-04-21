package internal_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"regexp"
	"strings"
	stdtesting "testing"

	"github.com/octohelm/x/cmp"
	"github.com/octohelm/x/testing/internal"
	. "github.com/octohelm/x/testing/v2"
)

type fakeTB struct {
	output bytes.Buffer
	errors []any
	fatals []any
}

func (f *fakeTB) Context() context.Context { return context.Background() }
func (f *fakeTB) TempDir() string          { return "" }
func (f *fakeTB) Output() io.Writer        { return &f.output }
func (f *fakeTB) Chdir(string)             {}
func (f *fakeTB) Setenv(string, string)    {}
func (f *fakeTB) Cleanup(func())           {}
func (f *fakeTB) Helper()                  {}
func (f *fakeTB) Skip(args ...any)         {}
func (f *fakeTB) Skipped() bool            { return false }
func (f *fakeTB) Error(args ...any)        { f.errors = append(f.errors, args...) }
func (f *fakeTB) Fatal(args ...any)        { f.fatals = append(f.fatals, args...) }

type plainTB struct {
	output bytes.Buffer
}

func (f *plainTB) Context() context.Context { return context.Background() }
func (f *plainTB) TempDir() string          { return "" }
func (f *plainTB) Output() io.Writer        { return &f.output }
func (f *plainTB) Chdir(string)             {}
func (f *plainTB) Setenv(string, string)    {}
func (f *plainTB) Cleanup(func())           {}
func (f *plainTB) Helper()                  {}
func (f *plainTB) Skip(args ...any)         {}

func TestReporter(t *stdtesting.T) {
	t.Run("Helper should preserve value", func(t *stdtesting.T) {
		reporter := &internal.Reporter{}

		Then(t, "Helper should return the original pointer",
			Expect(internal.Helper(0, reporter), Equal(reporter)),
		)
	})

	t.Run("Error and Fatal", func(t *stdtesting.T) {
		reporter := internal.Helper(0, &internal.Reporter{})
		tb := &fakeTB{}
		sentinel := errors.New("line one\nline two")

		reporter.Error(tb, sentinel)
		reporter.Fatal(tb, errors.New("fatal message"))
		reporter.Error(tb, nil)
		reporter.Fatal(tb, nil)

		Then(t, "nil error should be ignored",
			Expect(len(tb.errors), Equal(1)),
			Expect(len(tb.fatals), Equal(1)),
		)

		errValue, ok := tb.errors[0].(error)
		fatalValue, fatalOK := tb.fatals[0].(error)

		Then(t, "wrapped error should retain source and indentation",
			Expect(ok, Be(cmp.True())),
			Expect(fatalOK, Be(cmp.True())),
			Expect(errors.Is(errValue, sentinel), Be(cmp.True())),
			Expect(regexp.MustCompile(`core_test\.go:\d+:`).MatchString(errValue.Error()), Be(cmp.True())),
			Expect(strings.Contains(errValue.Error(), "    line one"), Be(cmp.True())),
			Expect(strings.Contains(errValue.Error(), "    line two"), Be(cmp.True())),
			Expect(strings.Contains(fatalValue.Error(), "fatal message"), Be(cmp.True())),
		)
	})

	t.Run("Reporter should allow TB without controller", func(t *stdtesting.T) {
		reporter := &internal.Reporter{}
		tb := &plainTB{}

		reporter.Error(tb, errors.New("ignored"))
		reporter.Fatal(tb, errors.New("ignored"))
	})
}

func TestFormatErrorMessage(t *stdtesting.T) {
	t.Run("regexp", func(t *stdtesting.T) {
		Then(t, "regexp message should describe match mode",
			Expect(strings.Contains(internal.FormatErrorMessage(regexp.MustCompile(`^a+$`), "bbb", false), "match:"), Be(cmp.True())),
			Expect(strings.Contains(internal.FormatErrorMessage(regexp.MustCompile(`^a+$`), "aaa", true), "not match:"), Be(cmp.True())),
		)
	})

	t.Run("multi line strings and bytes", func(t *stdtesting.T) {
		stringDiff := internal.FormatErrorMessage("a\nb\n", "a\nc\n", false)
		bytesDiff := internal.FormatErrorMessage([]byte("a\nb\n"), []byte("a\nc\n"), false)

		Then(t, "diff output should be used for multi line payloads",
			Expect(strings.Contains(stringDiff, "not match:"), Be(cmp.True())),
			Expect(strings.Contains(bytesDiff, "not match:"), Be(cmp.True())),
		)
	})

	t.Run("struct diff", func(t *stdtesting.T) {
		type sample struct {
			Name string
		}

		diff := internal.FormatErrorMessage(sample{Name: "x"}, sample{Name: "y"}, false)
		neg := internal.FormatErrorMessage(sample{Name: "x"}, sample{Name: "x"}, true)

		Then(t, "cmp reporter should describe field path",
			Expect(strings.Contains(diff, "Name:"), Be(cmp.True())),
			Expect(strings.Contains(diff, "expect:"), Be(cmp.True())),
			Expect(strings.Contains(neg, "not expect:"), Be(cmp.True())),
		)
	})
}
