package tensor

import (
	"bytes"
	"fmt"
	"math"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/slices"
	"github.com/gopherd/core/container/tuple"
	"github.com/gopherd/core/math/mathutil"
	"github.com/gopherd/core/operator"
)

// Vector implements n-dim vector
type Vector[T constraints.SignedReal] []T

// Vec creates a vector by elements
func Vec[T constraints.SignedReal](elements ...T) Vector[T] {
	return Vector[T](elements)
}

// Repeat creates a n-dim vector which has same value x
func Repeat[T constraints.SignedReal](x T, n int) Vector[T] {
	var vec = make(Vector[T], n)
	if x != 0 {
		for i := range vec {
			vec[i] = x
		}
	}
	return vec
}

// Range creates a vector [start, end)
func Range[T constraints.SignedReal](start, end T) Vector[T] {
	if end <= start {
		return nil
	}
	var vec = make(Vector[T], 0, int(end-start))
	for start < end {
		vec = append(vec, start)
		start++
	}
	return vec
}

// RangeN creates a vector [0, n)
func RangeN[T constraints.SignedReal](n T) Vector[T] {
	return Range(0, n)
}

// Linspace creates a vector {x1=from, ..., xn=to}
func Linspace[T constraints.SignedReal](from, to T, n int) Vector[T] {
	if n < 1 {
		return nil
	}
	if n == 1 {
		return Vec(from)
	}
	var interval = (to - from) / T(n-1)
	var vec = make(Vector[T], n)
	for i := 0; i < n; i++ {
		if i+1 == n {
			vec[i] = to
		} else {
			vec[i] = from + T(i)*interval
		}
	}
	return vec
}

// String converts vector as a string
func (vec Vector[T]) String() string {
	var buf bytes.Buffer
	buf.WriteByte('(')
	for i := range vec {
		if i > 0 {
			buf.WriteByte(',')
		}
		fmt.Fprint(&buf, vec[i])
	}
	buf.WriteByte(')')
	return buf.String()
}

// Shape implements Tensor Shape method
func (vec Vector[T]) Shape() Shape {
	return tuple.T1(len(vec))
}

// At implements Tensor At method
func (vec Vector[T]) At(index Shape) T {
	return vec[index.At(0)]
}

// Sum implements Tensor Sum method
func (vec Vector[T]) Sum() T {
	var sum T
	for i := range vec {
		sum += vec[i]
	}
	return sum
}

//----------------------------------------------
// basic functions

func (vec Vector[T]) Dim() int       { return len(vec) }
func (vec Vector[T]) Get(i int) T    { return operator.Ternary(i < len(vec), vec[i], 0) }
func (vec Vector[T]) Set(i int, v T) { vec[i] = v }

func (vec Vector[T]) Elements() []T {
	return []T(vec)
}

//----------------------------------------------
// operator functions

func (vec Vector[T]) Add(other Vector[T]) Vector[T] {
	var min, max = mathutil.Minmax(vec.Dim(), other.Dim())
	var out = make(Vector[T], max)
	for i := 0; i < min; i++ {
		out[i] = vec[i] + other[i]
	}
	return out
}

func (vec Vector[T]) Sub(other Vector[T]) Vector[T] {
	var min, max = mathutil.Minmax(vec.Dim(), other.Dim())
	var out = make(Vector[T], max)
	for i := 0; i < min; i++ {
		out[i] = vec[i] - other[i]
	}
	return out
}

func (vec Vector[T]) Mul(other Vector[T]) Vector[T] {
	var min, max = mathutil.Minmax(vec.Dim(), other.Dim())
	var out = make(Vector[T], max)
	for i := 0; i < min; i++ {
		out[i] = vec[i] * other[i]
	}
	return out
}

func (vec Vector[T]) Div(other Vector[T]) Vector[T] {
	var min, max = mathutil.Minmax(vec.Dim(), other.Dim())
	var out = make(Vector[T], max)
	for i := 0; i < max; i++ {
		out[i] = vec[i] / operator.Ternary(i < min, other[i], 1)
	}
	return out
}

//----------------------------------------------
// measure functions

func (vec Vector[T]) Dot(other Vector[T]) T {
	var sum T
	for i := range vec {
		if i >= len(other) {
			break
		}
		sum += vec[i] * other[i]
	}
	return sum
}

// SquaredLength computes ║vec║²
func (vec Vector[T]) SquaredLength() T {
	var sum T
	for i := range vec {
		sum += vec[i] * vec[i]
	}
	return sum
}

// Norm computes 2-norm
func (vec Vector[T]) Norm() T {
	return T(math.Sqrt(float64(vec.SquaredLength())))
}

// Normp computes p-norm
func (vec Vector[T]) Normp(p T) T {
	switch p {
	case 0:
		return T(len(vec)) - slices.SumFunc(vec, mathutil.IsZero[T])
	case 1:
		return slices.SumFunc(vec, mathutil.Abs[T])
	case 2:
		return vec.Norm()
	default:
		var sum float64
		for i := range vec {
			sum += math.Pow(float64(vec[i]), float64(p))
		}
		return T(math.Pow(sum, 1.0/float64(p)))
	}
}
