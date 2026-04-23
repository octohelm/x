// Package reflect 提供对标准库 reflect 的补充能力，聚焦零值判断、类型名称和结构标签解析等常见场景。
//
// 常见入口：
//   - 用 FullTypeName 获取包含包路径的完整类型名。
//   - 用 Deref 或 Indirect 递归去掉指针层级。
//   - 用 IsZero、IsEmptyValue 判断业务上的空值。
//   - 用 ParseStructTags 拆解结构体标签。
package reflect
