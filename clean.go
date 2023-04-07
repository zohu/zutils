package zutils

import "reflect"

// Clean
// @Description: 清空对象
// @param v
func Clean[T any](v T) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}
