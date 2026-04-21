# x

[![GoDoc Widget](https://godoc.org/github.com/octohelm/x?status.svg)](https://godoc.org/github.com/octohelm/x)

`x` 是一个面向 Go 的增强标准库集合，围绕“补 Go 标准库之间的组合空缺”组织了一组可按需引入的公共包。

仓库按包目录直接组织，适合按需引入单个能力，而不是作为一个需要整体初始化的框架。

## 内容总览

### 基础能力

- [`cmp`](cmp)：通用比较谓词与测试中可复用的条件判断
- [`container/list`](container/list)：泛型链表
- [`context`](context)：带类型辅助的上下文取值封装
- [`datauri`](datauri)：`data:` URI 的编解码
- [`iter`](iter)：轻量迭代辅助
- [`reflect`](reflect)：反射补充工具
- [`slices`](slices)：切片辅助函数
- [`sync`](sync)：泛型并发容器与并发辅助
- [`sync/singleflight`](sync/singleflight)：单飞控制扩展

### 类型与编码

- [`anyjson`](anyjson)：JSON value、patch、merge、diff 等相关能力
- [`encoding`](encoding)：文本编解码辅助
- [`types`](types)：`reflect.Type` 与 `go/types.Type` 之间的统一抽象与转换

### 日志

- [`logr`](logr)：仓库内统一日志接口
- [`logr/slog`](logr/slog)：基于标准库 `log/slog` 的适配实现

### 测试

- [`testing`](testing)：兼容旧用法的测试辅助入口
- [`testing/v2`](testing/v2)：当前推荐使用的测试断言入口
- [`testing/lines`](testing/lines)：文本行与 diff 辅助
- [`testing/snapshot`](testing/snapshot)：快照测试能力

## 相关文档

- [`AGENTS.md`](AGENTS.md)：仓库内协同约束、改动边界与接管条件
- [`go.mod`](go.mod)：模块定义与当前 Go / 依赖约束
- [`LICENSE`](LICENSE)：许可信息
