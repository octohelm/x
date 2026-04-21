package v2

import (
	"errors"
	"regexp"

	"github.com/octohelm/x/testing/internal"
)

// ErrorMatch 返回要求 error 文本匹配正则的检查器。
func ErrorMatch(re *regexp.Regexp) ValueChecker[error] {
	return internal.Helper(1, &beChecker[error]{
		be: func(actual error) error {
			if actual != nil && re.MatchString(actual.Error()) {
				return nil
			}
			return &ErrNotEqual{
				Expect: re,
				Got:    actual,
			}
		},
	})
}

// ErrorNotMatch 返回要求 error 文本不匹配正则的检查器。
func ErrorNotMatch(re *regexp.Regexp) ValueChecker[error] {
	return internal.Helper(1, &beChecker[error]{
		be: func(actual error) error {
			if actual == nil || !re.MatchString(actual.Error()) {
				return nil
			}
			return &ErrEqual{
				NotExpect: re,
				Got:       actual,
			}
		},
	})
}

// ErrorAsType 返回要求 error 链中可提取到指定类型的检查器。
func ErrorAsType[E error]() ValueChecker[error] {
	return internal.Helper(1, &beChecker[error]{
		be: func(actual error) error {
			if _, ok := errors.AsType[E](actual); ok {
				return nil
			}
			return &ErrNotEqual{Expect: new(E), Got: actual}
		},
	})
}

// ErrorNotAsType 返回要求 error 链中不可提取到指定类型的检查器。
func ErrorNotAsType[E error]() ValueChecker[error] {
	return internal.Helper(1, &beChecker[error]{
		be: func(actual error) error {
			if _, ok := errors.AsType[E](actual); !ok {
				return nil
			}
			return &ErrEqual{NotExpect: new(E), Got: actual}
		},
	})
}

// ErrorAs 返回要求 errors.As 可提取到 expect 的检查器。
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

// ErrorNotAs 返回要求 errors.As 不可提取到 expect 的检查器。
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

// ErrorNotIs 返回要求 errors.Is 不命中 expect 的检查器。
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

// ErrorIs 返回要求 errors.Is 命中 expect 的检查器。
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
