# Migration Notes

用于把旧测试风格收敛到 `github.com/octohelm/x/testing/v2`。

## 迁移方向

- `github.com/octohelm/x/testing`：旧入口，已弃用
- `github.com/octohelm/x/testing/bdd`：旧 BDD 入口，已弃用
- `github.com/octohelm/x/testing/v2`：当前推荐入口

## 迁移原则

1. 优先保持原测试语义不变，只替换断言入口。
2. 一个子场景对应一个 `Then(...)`，不要把多个互不相干的预期揉在同一段里。
3. 错误检查优先改成 `ErrorIs` / `ErrorAs` / `ErrorMatch`，不要继续手写 `if err == nil` 或字符串散比。
4. 可复用的值比较优先改成 `Equal`、`NotEqual`、`Be(cmp...)`；若缺少通用谓词，优先补到 `github.com/octohelm/x/cmp`。
5. 快照测试优先统一到 `MatchSnapshot`。

## 常见替换

- 手写 `if err != nil { t.Fatal(err) }`
  改为准备阶段 `Must(...)` / `MustValue(...)`
  或断言阶段 `ExpectMust(...)`

- 手写 `if !reflect.DeepEqual(...) { t.Fatal(...) }`
  改为 `Expect(actual, Equal(expect))`

- 手写 `errors.Is` / `errors.As` 分支
  改为 `ExpectDo(..., ErrorIs(...))`、`ExpectDo(..., ErrorAs(...))`

- 旧 BDD 风格嵌套断言
  改为 `Then(t, "预期描述", ...)`

## 审查点

- 是否仍引入旧测试包
- 是否把准备阶段失败和断言阶段失败混在一起
- `summary` 是否直接表达业务预期
- 是否存在可以复用 `github.com/octohelm/x/cmp` / `github.com/octohelm/x/testing/snapshot` 却仍手写的重复逻辑
- 是否需要补查 `go doc` 才能确认高阶用法，却被硬编码进局部猜测
