// Package sync 提供对标准库 sync 的泛型补充封装。
//
// 当前主要覆盖类型安全的 Map、Pool，以及独立子包中的 singleflight 能力。
//
// 常见入口：
//   - 用 Map 替代 sync.Map，减少类型断言。
//   - 用 Pool 替代 sync.Pool，保留泛型返回值。
//   - 需要按 key 去重并发调用时，使用子包 github.com/octohelm/x/sync/singleflight。
package sync
