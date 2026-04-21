package cmp

import (
	"fmt"

	"github.com/go-json-experiment/json/jsontext"
)

// ErrCondition 表示条件比较失败。
type ErrCondition struct {
	Op     string
	Expect any
	Actual any
}

func (e *ErrCondition) Error() string {
	return fmt.Sprintf("should be %s %v (got %v)", e.Op, e.Expect, e.Actual)
}

// ErrState 表示状态类检查失败，例如 nil 或 zero。
type ErrState struct {
	State  string
	Actual any
}

func (e *ErrState) Error() string {
	return fmt.Sprintf("should be %s (got %v)", e.State, e.Actual)
}

// ErrCheck 表示带主题和路径信息的嵌套检查失败。
type ErrCheck struct {
	Topic   string
	Err     error
	Actual  any
	Pointer jsontext.Pointer
}

func (e *ErrCheck) Error() string {
	if len(e.Pointer) > 0 {
		return fmt.Sprintf("%s: %s check failed: %v", e.Pointer, e.Topic, e.Err)
	}
	return fmt.Sprintf("%s check failed: %v", e.Topic, e.Err)
}
