package v2

import (
	"testing"

	"github.com/octohelm/x/testing/internal"
)

type (
	// TB 是测试运行上下文的最小接口。
	TB = internal.TB
	// TBController 是可报告失败或跳过状态的测试控制接口。
	TBController = internal.TBController
)

// Reporter 负责将检查失败包装为带位置信息的测试输出。
type Reporter = internal.Reporter

// Checker 表示一个可在测试上下文中执行的断言单元。
type Checker interface {
	Check(t TB)
}

// ValueChecker 表示针对某个实际值执行的断言。
type ValueChecker[V any] interface {
	Check(t TB, actual V)
}

// Then 以带摘要的子测试方式执行一组 Checker。
func Then(t *testing.T, summary string, checkers ...Checker) {
	if t.Skipped() {
		return
	}
	t.Helper()

	t.Run("THEN "+summary, func(t *testing.T) {
		if t.Skipped() {
			return
		}
		t.Helper()

		for _, c := range checkers {
			c.Check(t)
		}
	})
}
