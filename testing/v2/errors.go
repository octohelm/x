package v2

import (
	"github.com/octohelm/x/testing/internal"
)

// ErrNotEqual 表示实际值与期望值不相等。
type ErrNotEqual struct {
	Expect any
	Got    any
}

func (e *ErrNotEqual) Error() string {
	return internal.FormatErrorMessage(e.Expect, e.Got, false)
}

// ErrEqual 表示实际值意外等于不期望值。
type ErrEqual struct {
	NotExpect any
	Got       any
}

func (e *ErrEqual) Error() string {
	return internal.FormatErrorMessage(e.NotExpect, e.Got, true)
}
