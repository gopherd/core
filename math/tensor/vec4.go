package tensor

import (
	"fmt"
	"math"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/tuple"
	"github.com/gopherd/core/math/mathutil"
	"github.com/gopherd/core/operator"
)

// Vector4 implements 4d vector
type Vector4[T constraints.SignedReal] [4]T

// Vec4 creates a 4d vector by elements
func Vec4[T constraints.SignedReal](x, y, z, w T) Vector4[T] {
	return Vector4[T]{x, y, z, w}
}

// String converts vector as a string
func (vec Vector4[T]) String() string {
	return fmt.Sprintf("(%v,%v,%v,%v)", vec[0], vec[1], vec[2], vec[3])
}

// Shape implements Tensor Shape method
func (vec Vector4[T]) Shape() Shape { return tuple.T1(4) }

// At implements Tensor At method
func (vec Vector4[T]) At(index Shape) T {
	return vec.Get(index.At(0))
}

// Sum implements Tensor Sum method
func (vec Vector4[T]) Sum() T {
	return vec[0] + vec[1] + vec[2] + vec[3]
}

//----------------------------------------------
// basic functions

func (vec Vector4[T]) X() T { return vec[0] }
func (vec Vector4[T]) Y() T { return vec[1] }
func (vec Vector4[T]) Z() T { return vec[2] }
func (vec Vector4[T]) W() T { return vec[3] }

func (vec *Vector4[T]) SetX(x T) { vec[0] = x }
func (vec *Vector4[T]) SetY(y T) { vec[1] = y }
func (vec *Vector4[T]) SetZ(z T) { vec[2] = z }
func (vec *Vector4[T]) SetW(w T) { vec[3] = w }

func (vec Vector4[T]) Get(i int) T     { return operator.Ternary(i < len(vec), vec[i], 0) }
func (vec *Vector4[T]) Set(i int, v T) { vec[i] = v }

func (vec *Vector4[T]) SetElements(x, y, z, w T) {
	(*vec)[0], (*vec)[1], (*vec)[2], (*vec)[3] = x, y, z, w
}

func (vec *Vector4[T]) Copy(other Vector4[T]) {
	(*vec)[0], (*vec)[1], (*vec)[2], (*vec)[3] = other[0], other[1], other[2], other[3]
}

//----------------------------------------------
// operator functions

func (vec Vector4[T]) Add(other Vector4[T]) Vector4[T] {
	return Vec4(vec[0]+other[0], vec[1]+other[1], vec[2]+other[2], vec[3]+other[3])
}

func (vec Vector4[T]) Sub(other Vector4[T]) Vector4[T] {
	return Vec4(vec[0]-other[0], vec[1]-other[1], vec[2]-other[2], vec[3]-other[3])
}

func (vec Vector4[T]) Div(other Vector4[T]) Vector4[T] {
	return Vec4(vec[0]/other[0], vec[1]/vec[1], vec[2]/other[2], vec[3]/other[3])
}

func (vec Vector4[T]) Mul(other Vector4[T]) Vector4[T] {
	return Vec4(vec[0]*other[0], vec[1]*vec[1], vec[2]*other[2], vec[3]*other[3])
}

func (vec Vector4[T]) Scale(k T) Vector4[T] {
	return Vec4(vec[0]*k, vec[1]*k, vec[2]*k, vec[3]*k)
}

func (vec Vector4[T]) Normalize() Vector4[T] {
	return vec.Scale(vec.Norm())
}

func (vec Vector4[T]) Dot(other Vector4[T]) T {
	return vec[0]*other[0] + vec[1]*other[1] + vec[2]*other[2] + vec[3]*other[3]
}

//----------------------------------------------
// measure functions

func (vec Vector4[T]) SquaredLength() T {
	return vec.Dot(vec)
}

func (vec Vector4[T]) Norm() T {
	return T(math.Sqrt(float64(vec.SquaredLength())))
}

// Normp computes p-norm
func (vec Vector4[T]) Normp(p T) T {
	switch p {
	case 0:
		return 4 -
			mathutil.IsZero[T](vec[0]) -
			mathutil.IsZero[T](vec[1]) -
			mathutil.IsZero[T](vec[2]) -
			mathutil.IsZero[T](vec[3])
	case 1:
		return mathutil.Abs(vec[0]) +
			mathutil.Abs(vec[1]) +
			mathutil.Abs(vec[2]) +
			mathutil.Abs(vec[3])
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

//----------------------------------------------
// special functions

func (vec Vector4[T]) Vec2() Vector2[T] { return Vec2(vec[0], vec[1]) }
func (vec Vector4[T]) Vec3() Vector3[T] { return Vec3(vec[0], vec[1], vec[2]) }
