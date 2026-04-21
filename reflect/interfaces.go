package reflect

// ZeroChecker 表示类型可自行定义零值判定规则。
type ZeroChecker interface {
	IsZero() bool
}
