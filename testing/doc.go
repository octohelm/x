// Package testing 是旧版测试入口，已弃用。
//
// Deprecated: 请使用 github.com/octohelm/x/testing/v2。
//
// v2 提供向后兼容的入口，常见替换：
//   - Expect → testing/v2.Expect
//   - Equal → testing/v2.Equal
//   - Be → testing/v2.Be(cmp.Eq(...)) 或 testing/v2.Be(...)
//   - MustAsJSON → testing/v2.MustValue(..., Equal(...)) 或直接使用 go doc 查看对应能力
package testing
