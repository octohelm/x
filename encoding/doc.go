// Package encoding 提供围绕 encoding.TextMarshaler 和 encoding.TextUnmarshaler 的通用文本编解码辅助。
//
// 常见入口：
//   - 用 MarshalText 统一把常见标量、[]byte 或实现了 encoding.TextMarshaler 的值编码为文本。
//   - 用 UnmarshalText 把文本解码回目标值，支持指针和 reflect.Value。
package encoding
