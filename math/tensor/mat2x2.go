package tensor

import (
	"bytes"
	"fmt"
	"math"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/slices"
	"github.com/gopherd/core/container/tuple"
	"github.com/gopherd/core/math/mathutil"
)

// Matrix2 represents a 2x2 matrix
type Matrix2[T constraints.SignedReal] [2 * 2]T

// Zero2x2 creates a zero 2x2 matrix
func Zero2x2[T constraints.SignedReal]() Matrix2[T] {
	return Matrix2[T]{}
}

// One2x2 creates a 2x2 matrix which every element is 1
func One2x2[T constraints.SignedReal]() Matrix2[T] {
	return Matrix2[T]{
		1, 1,
		1, 1,
	}
}

// Identity2 creates a 2x2 identity matrix
func Identity2[T constraints.SignedReal]() Matrix2[T] {
	return Matrix2[T]{
		1, 0,
		0, 1,
	}
}

// String converts matrix as a string
func (mat Matrix2[T]) String() string {
	const n = 2
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteByte(';')
		}
		buf.WriteByte('(')
		for j := 0; j < n; j++ {
			if j > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprint(&buf, mat[i+j*n])
		}
		buf.WriteByte(')')
	}
	buf.WriteByte('}')
	return buf.String()
}

var shape2x2 = tuple.T2(2, 2)

// Shape implements Tensor Shape method
func (mat Matrix2[T]) Shape() Shape {
	return shape2x2
}

// At implements Tensor At method
func (mat Matrix2[T]) At(index Shape) T {
	return mat.Get(index.At(0), index.At(1))
}

// Sum implements Tensor Sum method
func (mat Matrix2[T]) Sum() T {
	var result T
	for i := range mat {
		result += mat[i]
	}
	return result
}

//----------------------------------------------
// basic functions

func (mat Matrix2[T]) Get(i, j int) T {
	return mat[i+j*2]
}

func (mat *Matrix2[T]) Set(i, j int, x T) {
	mat[i+j*2] = x
}

func (mat *Matrix2[T]) SetElements(n11, n12, n21, n22 T) *Matrix2[T] {
	(*mat)[0], (*mat)[2] = n11, n12
	(*mat)[1], (*mat)[3] = n21, n22
	return mat
}

//----------------------------------------------
// operator functions

func (mat Matrix2[T]) Transpose() Matrix2[T] {
	const dim = 2
	for i := 0; i < dim-1; i++ {
		for j := i + 1; j < dim; j++ {
			mat[i+j*dim], mat[j+i*dim] = mat[j+i*dim], mat[i+j*dim]
		}
	}
	return mat
}

func (mat Matrix2[T]) Add(other Matrix2[T]) Matrix2[T] {
	for i := range mat {
		mat[i] += other[i]
	}
	return mat
}

func (mat Matrix2[T]) Sub(other Matrix2[T]) Matrix2[T] {
	for i := range mat {
		mat[i] -= other[i]
	}
	return mat
}

func (mat Matrix2[T]) Mul(other Matrix2[T]) Matrix2[T] {
	for i := range mat {
		mat[i] *= other[i]
	}
	return mat
}

func (mat Matrix2[T]) Div(other Matrix2[T]) Matrix2[T] {
	for i := range mat {
		mat[i] /= other[i]
	}
	return mat
}

func (mat Matrix2[T]) Scale(v T) Matrix2[T] {
	for i := range mat {
		mat[i] *= v
	}
	return mat
}

func (mat Matrix2[T]) Normalize() Matrix2[T] {
	return mat.Scale(1 / mat.Norm())
}

func (mat Matrix2[T]) Dot(other Matrix2[T]) Matrix2[T] {
	const dim = 2
	var result Matrix2[T]
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			index := i + j*dim
			for k := 0; k < dim; k++ {
				result[index] += mat[i+k*dim] * other[k+j*dim]
			}
		}
	}
	return result
}

func (mat Matrix2[T]) DotVec2(vec Vector2[T]) Vector2[T] {
	const dim = 2
	var result Vector2[T]
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			result[i] += mat[i+j*dim] * vec[j]
		}
	}
	return result
}

func (mat Matrix2[T]) Invert() Matrix2[T] {
	var n11, n21 = mat[0], mat[1]
	var n12, n22 = mat[2], mat[3]
	var det = n11*n22 - n12*n21
	if det == 0 {
		return Matrix2[T]{}
	}
	var detInv = 1 / det
	mat[0] = n11 * detInv
	mat[1] = -n21 * detInv
	mat[2] = -n12 * detInv
	mat[3] = n22 * detInv
	return mat
}

//----------------------------------------------
// measure functions

func (mat Matrix2[T]) Determaint() T {
	return mat[0]*mat[3] - mat[1]*mat[2]
}

func (mat Matrix2[T]) SquaredLength() T {
	return mat.Mul(mat).Sum()
}

func (mat Matrix2[T]) Norm() T {
	return T(math.Sqrt(float64(mat.SquaredLength())))
}

func (mat Matrix2[T]) Normp(p T) T {
	switch p {
	case 0:
		return T(len(mat)) - slices.SumFunc(mat[:], mathutil.IsZero[T])
	case 1:
		return slices.SumFunc(mat[:], mathutil.Abs[T])
	case 2:
		return mat.Norm()
	default:
		var sum float64
		for _, v := range mat {
			sum += math.Pow(float64(v), float64(p))
		}
		return T(math.Pow(sum, 1/float64(p)))
	}
}

//-------------------------------------------------
// geometry functions

func (mat *Matrix2[T]) MakeIdentity() *Matrix2[T] {
	return mat.SetElements(
		1, 0,
		0, 1,
	)
}

func (mat *Matrix2[T]) MakeZero() *Matrix2[T] {
	return mat.SetElements(
		0, 0,
		0, 0,
	)
}

func (mat *Matrix2[T]) MakeRotation(theta T) *Matrix2[T] {
	var s0, c0 = math.Sincos(float64(theta))
	var s, c = T(s0), T(c0)
	return mat.SetElements(
		c, -s,
		s, c,
	)
}

func (mat *Matrix2[T]) MakeScale(vec Vector2[T]) *Matrix2[T] {
	return mat.SetElements(
		vec.X(), 0,
		0, vec.Y(),
	)
}
