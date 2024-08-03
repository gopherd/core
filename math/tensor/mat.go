package tensor

import (
	"bytes"
	"fmt"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/tuple"
	"github.com/gopherd/core/math/mathutil"
	"github.com/gopherd/core/operator"
)

// Matrix represents a MxN matrix
type Matrix[T constraints.SignedReal] struct {
	m, n     int       // m rows x n columns
	elements Vector[T] // a(i,j) = elements[i+j*m]

	transposed bool
}

// ZeroMxN creates a zero MxN matrix
func ZeroMxN[T constraints.SignedReal](m, n int) Matrix[T] {
	return Matrix[T]{
		m:        m,
		n:        n,
		elements: Repeat[T](0, m*n),
	}
}

// OneMxN creates a MxN matrix which every element is 1
func OneMxN[T constraints.SignedReal](m, n int) Matrix[T] {
	return Matrix[T]{
		m:        m,
		n:        n,
		elements: Repeat[T](1, m*n),
	}
}

// IdentityN creates a NxN identity matrix
func IdentityN[T constraints.SignedReal](n int) Matrix[T] {
	var m = ZeroMxN[T](n, n)
	for i := 0; i < n; i++ {
		m.Set(i, i, 1)
	}
	return m
}

// String converts matrix as a string
func (mat Matrix[T]) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	var m, n = mat.m, mat.n
	if mat.transposed {
		m, n = n, m
	}
	for i := 0; i < m; i++ {
		if i > 0 {
			buf.WriteByte(';')
		}
		buf.WriteByte('(')
		for j := 0; j < n; j++ {
			if j > 0 {
				buf.WriteByte(',')
			}
			if mat.transposed {
				fmt.Fprint(&buf, mat.elements[j+i*m])
			} else {
				fmt.Fprint(&buf, mat.elements[i+j*m])
			}
		}
		buf.WriteByte(')')
	}
	buf.WriteByte(']')
	return buf.String()
}

// Shape implements Tensor Shape method
func (mat Matrix[T]) Shape() Shape {
	return tuple.T2(mat.Rows(), mat.Columns())
}

// At implements Tensor At method
func (mat Matrix[T]) At(index Shape) T {
	return mat.Get(index.At(0), index.At(1))
}

// Sum implements Tensor Sum method
func (mat Matrix[T]) Sum() T {
	return mat.elements.Sum()
}

//----------------------------------------------
// basic functions

func (mat Matrix[T]) Rows() int {
	return operator.Ternary(mat.transposed, mat.n, mat.m)
}

func (mat Matrix[T]) Columns() int {
	return operator.Ternary(mat.transposed, mat.m, mat.n)
}

func (mat Matrix[T]) Get(i, j int) T {
	if mat.transposed {
		return mat.elements[j+i*mat.n]
	} else {
		return mat.elements[i+j*mat.m]
	}
}

func (mat *Matrix[T]) Set(i, j int, x T) {
	if mat.transposed {
		mat.elements[j+i*mat.n] = x
	} else {
		mat.elements[i+j*mat.m] = x
	}
}

func (mat Matrix[T]) Elements() []T {
	return mat.elements
}

//----------------------------------------------
// operator functions

// Transpose returns a transposed matrix
func (mat Matrix[T]) Transpose() Matrix[T] {
	mat.transposed = !mat.transposed
	return mat
}

func (mat Matrix[T]) unaryElementwise(x T, op func(T, T) T) Matrix[T] {
	for i := 0; i < mat.m; i++ {
		for j := 0; j < mat.n; j++ {
			mat.elements[i+j*mat.m] = op(mat.elements[i+j*mat.m], x)
		}
	}
	return mat
}

func (mat Matrix[T]) binaryElementwise(other Matrix[T], op func(T, T) T) Matrix[T] {
	var m, n = mat.Rows(), mat.Columns()
	if m != other.Rows() || n != other.Columns() {
		panic("matrix.add: size mismatched")
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			mat.Set(i, j, op(mat.Get(i, j), other.Get(i, j)))
		}
	}
	return mat
}

func (mat Matrix[T]) Add(other Matrix[T]) Matrix[T] {
	return mat.binaryElementwise(other, mathutil.Add[T])
}

func (mat Matrix[T]) Sub(other Matrix[T]) Matrix[T] {
	return mat.binaryElementwise(other, mathutil.Sub[T])
}

func (mat Matrix[T]) Mul(other Matrix[T]) Matrix[T] {
	return mat.binaryElementwise(other, mathutil.Mul[T])
}

func (mat Matrix[T]) Div(other Matrix[T]) Matrix[T] {
	return mat.binaryElementwise(other, mathutil.Div[T])
}

func (mat Matrix[T]) Scale(x T) Matrix[T] {
	return mat.unaryElementwise(x, mathutil.Mul[T])
}

func (mat Matrix[T]) Normalize() Matrix[T] {
	return mat.Scale(1 / mat.Norm())
}

func (mat Matrix[T]) Dot(other Matrix[T]) Matrix[T] {
	if mat.Columns() != other.Rows() {
		panic("matrix.dot: size mismatched")
	}
	var m, n, p = mat.Rows(), mat.Columns(), other.Columns()
	var result = Matrix[T]{
		m:        m,
		n:        p,
		elements: make(Vector[T], m*p),
	}
	for i := 0; i < m; i++ {
		for j := 0; j < p; j++ {
			index := i + j*m
			for k := 0; k < n; k++ {
				result.elements[index] += mat.elements[i+k*m] * other.elements[k+j*n]
			}
		}
	}
	return result
}

func (mat Matrix[T]) DotVec(vec Vector[T]) Vector[T] {
	if mat.Columns() != vec.Dim() {
		panic("matrix.dotVec: size mismatched")
	}
	var m, n = mat.Rows(), mat.Columns()
	var result = make(Vector[T], m)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			result[i] += mat.Get(i, j) * vec.Get(j)
		}
	}
	return result
}

func (mat Matrix[T]) Invert() Matrix[T] {
	panic("TODO: matrix.invert not implemented")
}

//----------------------------------------------
// measure functions

func (mat Matrix[T]) Determinant() T {
	// TODO: computes determinant of matrix
	return 0
}

func (mat Matrix[T]) SquaredLength() T {
	return mat.elements.SquaredLength()
}

func (mat Matrix[T]) Norm() T {
	return mat.elements.Norm()
}

func (mat Matrix[T]) Normp(p T) T {
	return mat.elements.Normp(p)
}
