package testing

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/onsi/gomega/format"
)

func Expect[A any](t *testing.T, actual A, assertions ...Assertion[A]) {
	t.Helper()
	for i := range assertions {
		assertions[i].Check(t, actual)
	}
}

func Should[A any, E any](match func(actual A, expected E) bool, expected E) Assertion[A] {
	return &assertion[A, E]{
		expected: expected,
		match:    match,
	}
}

func ShouldNot[A any, E any](match func(actual A, expected E) bool, expected E) Assertion[A] {
	return &assertion[A, E]{
		match:    match,
		expected: expected,
		negative: true,
	}
}

type Assertion[A any] interface {
	Check(t *testing.T, actual A)
}

type assertion[A any, E any] struct {
	match    func(actual A, expected E) bool
	expected E
	negative bool
}

func (s *assertion[A, E]) Check(t *testing.T, actual A) {
	t.Helper()

	ok := s.match(actual, s.expected)
	if s.negative {
		if !ok {
			return
		}
		s.fail(t, actual)
		return
	}
	if ok {
		return
	}
	s.fail(t, actual)
	return
}

func (s *assertion[A, E]) fail(t *testing.T, actual A) {
	name := s.matchName()
	t.Helper()
	t.Fatalf("\n" + s.failureMessage(actual, name))
}

func (s *assertion[A, E]) matchName() string {
	pc := reflect.ValueOf(s.match).Pointer()
	f := runtime.FuncForPC(pc)
	name := f.Name()
	parts := strings.Split(name, ".")
	if len(parts) == 2 {
		return strings.ToLower(parts[1])
	}
	file, line := f.FileLine(pc)
	return fmt.Sprintf("match (%s:%d)", file, line)
}

func (s *assertion[A, E]) failureMessage(actual A, name string) string {
	if s.negative {
		return format.MessageWithDiff(
			fmt.Sprintf("%v", actual),
			fmt.Sprintf("Should not %s", name),
			fmt.Sprintf("%v", s.expected),
		)
	}
	return format.MessageWithDiff(
		fmt.Sprintf("%v", actual),
		fmt.Sprintf("Should %s", name),
		fmt.Sprintf("%v", s.expected),
	)
}
