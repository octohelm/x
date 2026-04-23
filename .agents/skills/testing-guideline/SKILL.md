---
name: testing-guideline
description: 使用仓库内 `github.com/octohelm/x/testing/v2` 编写、迁移或审查 Go 测试时使用，帮助统一断言写法、错误检查和快照检查入口，并约束扩展优先落在 `github.com/octohelm/x/cmp`。
---

# Testing Guideline

用于在本仓库中统一使用 `github.com/octohelm/x/testing/v2` 编写测试。

## 何时使用

- 新增 Go 单测，且需要统一断言风格
- 将旧的 `github.com/octohelm/x/testing` / `github.com/octohelm/x/testing/bdd` 用法迁移到 `github.com/octohelm/x/testing/v2`
- 审查现有测试，判断断言、错误检查和快照检查是否符合仓库约定

## 非目标

- 不负责设计业务测试用例本身
- 不替代 `go test`、快照文件维护或仓库级测试策略
- 不继续扩展已弃用的 `testing` / `testing/bdd` 风格

## 关键约定

1. 优先使用 `Then(t, summary, ...)` 表达一个完整断言场景，`summary` 直接描述预期。
2. 值断言优先使用 `Expect(actual, ...)`，动作用 `ExpectMust`、`ExpectDo`、`ExpectMustValue`。
3. 通用比较优先使用 `Equal`、`NotEqual`；复杂谓词通过 `Be(...)` 组合 `github.com/octohelm/x/cmp` 包能力。
4. 需要立即取值并在失败时终止时，使用 `Must`、`MustValue`、`MustValues`。
5. 快照测试使用 `SnapshotOf(...)` 配合 `MatchSnapshot(name)`；不要手写分散的快照比对逻辑。
6. 新代码不要继续引入已弃用的 `github.com/octohelm/x/testing` 或 `github.com/octohelm/x/testing/bdd`。
7. 需要特殊断言扩展时，优先补充到 `github.com/octohelm/x/cmp`，再通过 `Be(...)` 复用，而不是直接在测试里散落一次性谓词。
8. 需要高阶或边界用法时，优先查 `go doc github.com/octohelm/x/testing/v2` 以及相关子包文档，不在 skill 中复制完整 API 手册。

## 工作方式

1. 先判断当前测试是值断言、错误断言、动作断言还是快照断言。
2. 先读 [`references/api-map.md`](references/api-map.md) 选择对应入口。
3. 若任务包含迁移旧测试，再读 [`references/migration.md`](references/migration.md)。
4. 若需要高阶 API、具体签名或边界行为，执行 `go doc github.com/octohelm/x/testing/v2`、`go doc github.com/octohelm/x/cmp` 或对应完整包路径。
5. 落地时保持测试描述准确中文，直接表达输入条件与预期。

## 完成标准

- 测试入口统一使用 `github.com/octohelm/x/testing/v2`
- 断言类型与 API 选择匹配
- 未继续引入旧测试包
- 需要扩展时优先复用或补充 `github.com/octohelm/x/cmp`
- 需要时已使用快照或错误链断言，而不是手写重复样板
