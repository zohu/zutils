package zutils

func Ptr[T Primitive](s T) *T {
	return &s
}
