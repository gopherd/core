package mathutil_test

import (
	"testing"
	"unsafe"
)

func TestXXX(t *testing.T) {
	println("sizeof(float32)", unsafe.Sizeof(float32(0)))
	println("sizeof(float64)", unsafe.Sizeof(float64(0)))
	println("sizeof(complex64)", unsafe.Sizeof(complex64(0)))
	println("sizeof(complex128)", unsafe.Sizeof(complex128(0)))
}
