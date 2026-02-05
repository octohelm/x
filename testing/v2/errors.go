package v2

import (
	"github.com/octohelm/x/testing/internal"
)

type ErrNotEqual struct {
	Expect any
	Got    any
}

func (e *ErrNotEqual) Error() string {
	return internal.FormatErrorMessage(e.Expect, e.Got, false)
}

type ErrEqual struct {
	NotExpect any
	Got       any
}

func (e *ErrEqual) Error() string {
	return internal.FormatErrorMessage(e.NotExpect, e.Got, true)
}
