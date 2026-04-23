# x

`x` 是一个 Go 公共库集合，提供可复用的基础能力与测试辅助组件。

## 职责与边界

- 以库代码为主，面向复用，不在 root 承载业务项目入口
- 根目录负责聚合仓库级控制面与工具链入口
- 具体执行命令收敛到 `justfile` 和 `tool/` 下的对应工具链目录

## 目录索引

- `anyjson`：动态 JSON 值表示，以及 diff、merge、transform 等操作
- `cmp`：可组合断言谓词与结构化错误
- `container/list`：泛型双向链表实现
- `context`：类型化上下文槽位和值注入辅助
- `datauri`：data URI 的解析、编码与文本序列化
- `encoding`：围绕 `encoding.TextMarshaler` / `encoding.TextUnmarshaler` 的通用文本编解码辅助
- `iter`：标准库 `iter` 的轻量辅助封装
- `logr`：轻量日志抽象、上下文注入和日志级别辅助
- `logr/slog`：`logr` 与标准库 `log/slog` 的适配
- `ptr`：基础值到指针的便捷转换
- `reflect`：零值判断、类型名称和结构标签等反射辅助
- `slices`：切片相关的轻量泛型辅助函数
- `sync`：对标准库 `sync` 的泛型补充封装
- `sync/singleflight`：按键去重的并发调用抑制
- `testing`：旧版测试断言入口，已弃用，建议改用 `testing/v2`
- `testing/bdd`：旧版 BDD 测试辅助，已弃用，建议改用 `testing/v2`
- `testing/lines`：按行表示文本和生成行级 diff 的辅助
- `testing/snapshot`：基于 `txtar` 的测试快照装载、比较和更新
- `testing/v2`：当前推荐的测试断言 API
- `types`：桥接 `reflect.Type` 与 `go/types.Type` 的统一类型抽象
- `tool/go`：Go 工具链执行入口
- `.github/workflows`：仓库 CI 配置

## 继续阅读

- Go 模块定义见 `go.mod`
- 仓库级执行入口见 `justfile`
- 各包的导出能力与用法以源码和 Go doc 为准
