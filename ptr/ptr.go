package ptr

//go:fix inline
func Ptr[T any](v T) *T {
	return new(v)
}
