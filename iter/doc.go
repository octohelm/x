// Package iter 提供围绕标准 iter 包的轻量辅助封装。
//
// 当前主要入口是 Action，用于把 yield 风格且可能返回 error 的回调包装成 iter.Seq2。
package iter
