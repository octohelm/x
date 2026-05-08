---
name: testing-guideline
description: 使用 `github.com/octohelm/x/testing/v2` 编写、迁移或审查 Go 测试时使用，统一断言写法、错误检查和快照检查入口。
---

# Testing Guideline

按 `github.com/octohelm/x/testing/v2` 约定写 Go 测试。

## 写一个新测试

```go
import . "github.com/octohelm/x/testing/v2"

func TestXxx(t *testing.T) {
    Then(t, "返回值符合预期",
        Expect(actual, Equal(expect)),
    )
}
```

**场景选入口**：

| 场景 | 入口 |
|------|------|
| 已有值做检查 | `Expect(actual, checkers...)` |
| 执行动作再检查 | `ExpectDo(fn, ErrorIs(...))` |
| 执行动作必须成功 | `ExpectMust(fn)` |
| 取返回值 | `ExpectMustValue(fn, Equal(...))` |
| 立即取值否则失败 | `Must(fn)` / `MustValue(fn)` |
| 快照测试 | `SnapshotOf(...)` + `MatchSnapshot(name)` |
| 复杂谓词 | `Be(cmp.Gt(...))` / `Be(cmp.Len(...))` |

**关键约定**：
- 值比较用 `Equal` / `NotEqual`，通用谓词优先补到 `cmp` 包
- 新代码不引入旧 `testing` / `testing/bdd`
- 具体 API 签名以 `go doc github.com/octohelm/x/testing/v2` 为准

## 迁移旧测试

见 [references/migration.md](references/migration.md)——旧入口到新入口的对照替换。

## 选择原则和更多示例

见 [references/api-map.md](references/api-map.md)。
