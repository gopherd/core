package tensor

import (
	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/tuple"
)

// Shape represents shape of tensor
type Shape = tuple.Tuple[int]

// Index implements a 1-dim shape
type Index int

func (i Index) Len() int     { return 1 }
func (i Index) At(j int) int { return int(i) }

// Index implements a n-dim shape
type Indices []int

func (i Indices) Len() int     { return len(i) }
func (i Indices) At(j int) int { return i[j] }

// Tensor represents a tensor
type Tensor[T constraints.SignedReal] interface {
	Shape() tuple.Tuple[int]
	At(index tuple.Tuple[int]) T
	Sum() T
}

// Create creates a tensor by shape
func Create[T constraints.SignedReal](shape Shape) Tensor[T] {
	return ndarray[T]{
		shape: shape,
		data:  make(Vector[T], SizeOf(shape)),
	}
}

// ndarray implements Tensor with shape
type ndarray[T constraints.SignedReal] struct {
	shape Shape
	data  Vector[T]
}

// Shape implements Tensor Shape method
func (t ndarray[T]) Shape() Shape {
	return t.shape
}

// At implements Tensor At method
func (t ndarray[T]) At(index Shape) T {
	return t.data[OffsetOf(t.shape, index)]
}

// Sum implements Tensor Sum method
func (t ndarray[T]) Sum() T {
	return t.data.Sum()
}

// set updates value by index
func (t *ndarray[T]) set(index Shape, value T) {
	t.data[OffsetOf(t.shape, index)] = value
}

// OffsetOf calculates offset of index in shape
func OffsetOf(shape, index Shape) int {
	var off int
	var prev = 1
	for i, n := 0, index.Len(); i < n; i++ {
		off += index.At(i) * prev
		prev *= shape.At(i)
	}
	return off
}

// IndexOf calculates index of offset in shape and returns indices
func IndexOf(shape Shape, offset int, indices Indices) Indices {
	if len(indices) != shape.Len() {
		indices = make(Indices, shape.Len())
	}
	var size = SizeOf(shape)
	for i := len(indices) - 1; i >= 0; i-- {
		if offset == 0 {
			indices[i] = 0
		} else {
			size /= shape.At(i)
			var value = offset / size
			indices[i] = value
			offset -= value * size
		}
	}
	return indices
}

// SizeOf calculates size of shape
func SizeOf(shape Shape) int {
	if shape.Len() == 0 {
		return 1
	}
	var size = 1
	for i, n := 0, shape.Len(); i < n; i++ {
		size *= shape.At(i)
	}
	return size
}

// Next retrieves next index
func Next(shape Shape, index Indices) Indices {
	for i, x := range index {
		var up = shape.At(i)
		if x+1 > up {
			continue
		}
		x++
		if x < up {
			index[i] = x
			return index
		}
		var j = i + 1
		for ; j < len(index); j++ {
			if index[j]+1 < shape.At(j) {
				break
			}
		}
		if j == len(index) {
			return nil
		}
		index[j]++
		for k := i; k < j; k++ {
			index[k] = 0
		}
		return index
	}
	return nil
}

// Product computes tensor product: a⊗b
func Product[T constraints.SignedReal](a, b Tensor[T]) Tensor[T] {
	var ashape, bshape = a.Shape(), b.Shape()
	return productedTensor[T]{
		a:     a,
		b:     b,
		m:     ashape.Len(),
		n:     bshape.Len(),
		shape: tuple.Concat(ashape, bshape),
	}
}

type productedTensor[T constraints.SignedReal] struct {
	a, b  Tensor[T]
	m, n  int
	shape Shape
}

// Shape implements Tensor Shape method
func (t productedTensor[T]) Shape() Shape {
	return t.shape
}

// At implements Tensor At method
func (t productedTensor[T]) At(index Shape) T {
	var i, j = tuple.Slice(index, 0, t.m), tuple.Slice(index, t.m, t.m+t.n)
	return t.a.At(i) * t.b.At(j)
}

// Sum implements Tensor Sum method
func (t productedTensor[T]) Sum() T {
	return t.a.Sum() * t.b.Sum()
}

// Dot computes dot product: a‧b
func Dot[T constraints.SignedReal](a, b Tensor[T]) Tensor[T] {
	var ashape = a.Shape()
	var bshape = b.Shape()
	var alen, blen = ashape.Len(), bshape.Len()
	if (alen == 0) != (blen == 0) {
		panic("tensor.dot: size mismatched")
	} else if alen == 0 {
		return Scalar(a.Sum() * b.Sum())
	}
	var n = ashape.At(ashape.Len() - 1)
	if n != bshape.At(0) {
		panic("tensor.dot: size mismatched")
	}

	var shape = tuple.Concat(
		tuple.Slice(ashape, 0, alen),
		tuple.Slice(bshape, 1, blen),
	)

	// shape.Len == alen + blen - 2
	if shape.Len() == 0 {
		var sum T
		for i := Index(0); i < Index(n); i++ {
			sum += a.At(i) * b.At(i)
		}
		return Scalar(sum)
	}

	// c = a‧b
	var c = ndarray[T]{
		shape: shape,
		data:  make(Vector[T], SizeOf(shape)),
	}
	var indices = make(Indices, shape.Len())
	var aindices = make(Indices, alen)
	var bindices = make(Indices, blen)
	for len(indices) > 0 {
		indices = Next(shape, indices)
		copy(aindices[:alen-1], indices[:alen-1])
		copy(bindices[1:], indices[alen-1:])
		var sum T
		for i := 0; i < n; i++ {
			aindices[alen-1] = i
			bindices[0] = i
			sum += a.At(aindices) * b.At(bindices)
			c.set(indices, sum)
		}
	}
	return c
}

// Reshape reshapes the tensor
func Reshape[T constraints.SignedReal](tensor Tensor[T], shape Shape) Tensor[T] {
	var c = ndarray[T]{
		shape: shape,
		data:  make(Vector[T], SizeOf(shape)),
	}
	var tshape = tensor.Shape()
	var tsize = SizeOf(tshape)
	var indices = make(Indices, tshape.Len())
	for i := range c.data {
		if i >= tsize {
			break
		}
		c.data[i] = tensor.At(indices)
		indices = Next(tshape, indices)
		if len(indices) == 0 {
			break
		}
	}
	return c
}
