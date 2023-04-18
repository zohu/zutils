package types

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
type Float interface {
	~float32 | ~float64
}
type Complex interface {
	~complex64 | ~complex128
}
type Number interface {
	Int | Float | Complex
}
type NumStr interface {
	Number | ~string
}
type Primitive interface {
	NumStr | ~bool | ~uintptr
}
type Object interface {
	Primitive | any
}
