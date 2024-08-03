package tensor_test

import (
	"testing"

	"github.com/gopherd/core/math/tensor"
)

func TestMatrix(t *testing.T) {
	type T = float32
	var mat1 = tensor.ZeroMxN[T](2, 3)
	var mat2 = tensor.ZeroMxN[T](3, 2)

	mat1.Set(0, 0, 1)
	mat1.Set(0, 1, 2)
	mat1.Set(0, 2, 3)
	mat1.Set(1, 0, 4)
	mat1.Set(1, 1, 5)
	mat1.Set(1, 2, 6)

	mat2.Set(0, 0, 1)
	mat2.Set(0, 1, 2)
	mat2.Set(1, 0, 3)
	mat2.Set(1, 1, 4)
	mat2.Set(2, 0, 5)
	mat2.Set(2, 1, 6)

	var mat3 = mat1.Dot(mat2)
	t.Logf("mat1 %v dot mat2 %v: %v", mat1, mat2, mat3)
}
