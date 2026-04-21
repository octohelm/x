package ptr

// Ptr 返回 v 的指针副本。
//
// Deprecated: Go 1.26 起可直接使用内建的 new(T) 能力完成同类用途。
//
//go:fix inline
func Ptr[T any](v T) *T {
	return new(v)
}
