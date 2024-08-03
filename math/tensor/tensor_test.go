package tensor_test

import (
	"math/rand"
	"testing"

	"github.com/gopherd/core/container/tuple"
	"github.com/gopherd/core/math/tensor"
)

func TestShape(t *testing.T) {
	const dim = 4
	var shape = make(tensor.Indices, dim)
	for i := range shape {
		shape[i] = 2*i + 3
	}
	const N = 1 << 20

	var size = tensor.SizeOf(shape)
	var index = make(tensor.Indices, dim)
	for i := 0; i < N; i++ {
		for j := 0; j < dim; j++ {
			index[j] = rand.Intn(2*j + 3)
		}
		var offset = tensor.OffsetOf(shape, index)
		var gotIndex = tensor.IndexOf(shape, offset, nil)
		if !tuple.Equal[int](index, gotIndex) {
			t.Fatalf("index mismatched: want %v, but got %v", index, gotIndex)
		}
		if offset+1 < size {
			var next = tensor.Next(shape, index)
			if offset+1 != tensor.OffsetOf(shape, next) {
				t.Fatalf("next offset mismatched: got next %v for %v", next, gotIndex)
			}
		}
	}
}
