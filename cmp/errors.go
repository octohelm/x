package cmp

import "fmt"

type ErrCondition struct {
	Op     string
	Expect any
	Actual any
}

func (e *ErrCondition) Error() string {
	return fmt.Sprintf("should be %s %v (got %v)", e.Op, e.Expect, e.Actual)
}

type ErrState struct {
	State  string
	Actual any
}

func (e *ErrState) Error() string {
	return fmt.Sprintf("should be %s (got %v)", e.State, e.Actual)
}

type ErrCheck struct {
	Topic  string
	Err    error
	Actual any
}

func (e *ErrCheck) Error() string {
	return fmt.Sprintf("%s check failed: %v", e.Topic, e.Err)
}
