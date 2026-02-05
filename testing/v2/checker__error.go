package v2

import (
	"errors"
	"regexp"

	"github.com/octohelm/x/testing/internal"
)

func ErrorMatch(re *regexp.Regexp) ValueChecker[error] {
	return internal.Helper(1, &beChecker[error]{
		be: func(actual error) error {
			if re.MatchString(actual.Error()) {
				return nil
			}
			return &ErrNotEqual{
				Expect: re,
				Got:    actual,
			}
		},
	})
}

func ErrorNotMatch(re *regexp.Regexp) ValueChecker[error] {
	return internal.Helper(1, &beChecker[error]{
		be: func(actual error) error {
			if !re.MatchString(actual.Error()) {
				return nil
			}
			return &ErrEqual{
				NotExpect: re,
				Got:       actual,
			}
		},
	})
}

func ErrorAs[E error](expect *E) ValueChecker[error] {
	return internal.Helper(1, &beChecker[error]{
		be: func(actual error) error {
			if errors.As(actual, expect) {
				return nil
			}
			return &ErrNotEqual{Expect: expect, Got: actual}
		},
	})
}

func ErrorNotAs[E error](expect *E) ValueChecker[error] {
	return internal.Helper(1, &beChecker[error]{
		be: func(actual error) error {
			if !errors.As(actual, expect) {
				return nil
			}
			return &ErrEqual{NotExpect: expect, Got: actual}
		},
	})
}

func ErrorNotIs(expect error) ValueChecker[error] {
	return internal.Helper(1, &beChecker[error]{
		be: func(actual error) error {
			if !errors.Is(actual, expect) {
				return nil
			}
			return &ErrEqual{NotExpect: expect, Got: actual}
		},
	})
}

func ErrorIs(expect error) ValueChecker[error] {
	return internal.Helper(1, &beChecker[error]{
		be: func(actual error) error {
			if errors.Is(actual, expect) {
				return nil
			}
			return &ErrNotEqual{Expect: expect, Got: actual}
		},
	})
}
