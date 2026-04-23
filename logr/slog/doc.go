// Package slog 提供 logr.Logger 与标准库 log/slog 之间的适配实现。
//
// 常见入口：
//   - 用 Logger 把标准库 *slog.Logger 适配为 github.com/octohelm/x/logr.Logger。
//   - 用 Default 创建带最小级别过滤的 *slog.Logger。
//   - 用 EnableLevel 配置 Default 返回的 logger 最低输出级别。
package slog
