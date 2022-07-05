package util

import (
	"reflect"
	"unsafe"
)

func StringToBytes(s string) (b []byte) {
	v := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bs := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bs.Data = v.Data
	bs.Len = v.Len
	bs.Cap = v.Len
	return
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
