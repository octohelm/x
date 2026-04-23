# `github.com/octohelm/x/testing/v2` API Map

用于在编写或审查测试时快速选择合适入口。

## 场景到 API

- 单个场景的一组断言：`Then(t, "摘要", ...)`
- 对已有值做检查：`Expect(actual, checkers...)`
- 先执行取值动作，再检查返回值：`ExpectMustValue(func() (V, error), checkers...)`
- 先执行动作，要求必须成功：`ExpectMust(func() error)`
- 先执行动作，再检查返回的 `error`：`ExpectDo(func() error, errorCheckers...)`
- 需要立即取值或取多个值，否则立刻失败：`Must`、`MustValue`、`MustValues`

## 常用值检查

- 深度相等：`Equal(expect)`
- 深度不等：`NotEqual(expect)`
- 自定义谓词：`Be(func(v V) error)`
- 与 `github.com/octohelm/x/cmp` 组合：`Be(cmp.Gt(...))`、`Be(cmp.Len(...))` 等

## 常用错误检查

- 错误链命中：`ErrorIs(expect)`
- 错误链不命中：`ErrorNotIs(expect)`
- 通过 `errors.As` 提取指定实例：`ErrorAs(expectPtr)`
- 不应提取到指定实例：`ErrorNotAs(expectPtr)`
- 按类型提取：`ErrorAsType[E]()`、`ErrorNotAsType[E]()`
- 错误文本匹配正则：`ErrorMatch(re)`、`ErrorNotMatch(re)`

## 快照检查

- 构造快照序列：`SnapshotOf(files...)`
- 从原始内容构造文件：`SnapshotFileFromRaw(name, raw)`
- 断言快照：`Expect(snapshot, MatchSnapshot(name))`

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
- 断言复杂条件时，优先复用 `github.com/octohelm/x/cmp` 与 `Be`，不要在测试体里散落 `if` 分支
- 需要新增特殊谓词时，优先把可复用逻辑实现到 `github.com/octohelm/x/cmp`
- 需要完整签名、泛型边界或未在此列出的能力时，优先执行 `go doc github.com/octohelm/x/testing/v2`、`go doc github.com/octohelm/x/cmp` 或对应完整包路径
