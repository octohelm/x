// Package v2 提供当前推荐使用的测试断言 API，覆盖值检查、错误检查和快照检查等能力。
//
// 常见入口：
//   - 用 Then 组织一个具名断言场景。
//   - 用 Expect 检查已有值。
//   - 用 ExpectDo、ExpectMust、ExpectMustValue 检查会返回 error 的动作或取值过程。
//   - 用 Equal、NotEqual、Be 和 ErrorIs / ErrorAs / ErrorMatch 等 ValueChecker 组合具体断言。
//   - 用 SnapshotOf 和 MatchSnapshot 组织快照检查。
//
// 需要更细粒度的比较逻辑时，优先与 github.com/octohelm/x/cmp 组合使用。
package v2
