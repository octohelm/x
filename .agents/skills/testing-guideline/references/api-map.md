# `github.com/octohelm/x/testing/v2` 示例与选择原则

具体 API 签名、参数和完整列表始终以 `go doc` 为准：
- `go doc github.com/octohelm/x/testing/v2`
- `go doc github.com/octohelm/x/cmp`

## 推荐写法

```go
Then(t, "返回值符合预期",
    Expect(actual,
        Equal(expect),
    ),
)

Then(t, "错误链包含目标错误",
    ExpectDo(
        func() error { return doSomething() },
        ErrorIs(os.ErrNotExist),
    ),
)
```

## 选择原则

- 已经有值：优先 `Expect`
- 先执行动作：优先 `ExpectDo` / `ExpectMust`
- 需要在准备阶段立刻拿结果：优先 `MustValue` / `MustValues`
- 断言复杂条件时，复用 `github.com/octohelm/x/cmp` 与 `Be`，不要在测试体内写 `if` 分支
- 需要新增通用谓词时，优先实现到 `github.com/octohelm/x/cmp`，通过 `Be(cmp.Eq(...))` 复用
