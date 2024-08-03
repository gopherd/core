package tensor

import (
	"fmt"
	"math"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/tuple"
	"github.com/gopherd/core/math/mathutil"
	"github.com/gopherd/core/operator"
)

// Vector3 implements 3d vector
type Vector3[T constraints.SignedReal] [3]T

// Vec3 creates a 3d vector by elements
func Vec3[T constraints.SignedReal](x, y, z T) Vector3[T] {
	return Vector3[T]{x, y, z}
}

// String converts vector as a string
func (vec Vector3[T]) String() string {
	return fmt.Sprintf("(%v,%v,%v)", vec[0], vec[1], vec[2])
}

// Shape implements Tensor Shape method
func (vec Vector3[T]) Shape() Shape { return tuple.T1(3) }

// At implements Tensor At method
func (vec Vector3[T]) At(index Shape) T {
	return vec.Get(index.At(0))
}

// Sum implements Tensor Sum method
func (vec Vector3[T]) Sum() T {
	return vec[0] + vec[1] + vec[2]
}

//----------------------------------------------
// basic functions

func (vec Vector3[T]) X() T { return vec[0] }
func (vec Vector3[T]) Y() T { return vec[1] }
func (vec Vector3[T]) Z() T { return vec[2] }

func (vec *Vector3[T]) SetX(x T) { vec[0] = x }
func (vec *Vector3[T]) SetY(y T) { vec[1] = y }
func (vec *Vector3[T]) SetZ(z T) { vec[2] = z }

func (vec Vector3[T]) Get(i int) T     { return operator.Ternary(i < len(vec), vec[i], 0) }
func (vec *Vector3[T]) Set(i int, v T) { vec[i] = v }

func (vec *Vector3[T]) SetElements(x, y, z T) {
	(*vec)[0], (*vec)[1], (*vec)[2] = x, y, z
}

func (vec *Vector3[T]) Copy(other Vector3[T]) {
	(*vec)[0], (*vec)[1], (*vec)[2] = other[0], other[1], other[2]
}

//----------------------------------------------
// operator functions

func (vec Vector3[T]) Add(other Vector3[T]) Vector3[T] {
	return Vec3(vec[0]+other[0], vec[1]+other[1], vec[2]+other[2])
}

func (vec Vector3[T]) Sub(other Vector3[T]) Vector3[T] {
	return Vec3(vec[0]-other[0], vec[1]-other[1], vec[2]-other[2])
}

func (vec Vector3[T]) Mul(other Vector3[T]) Vector3[T] {
	return Vec3(vec[0]*other[0], vec[1]*vec[1], vec[2]*other[2])
}

func (vec Vector3[T]) Div(other Vector3[T]) Vector3[T] {
	return Vec3(vec[0]/other[0], vec[1]/vec[1], vec[2]/other[2])
}

func (vec Vector3[T]) Scale(k T) Vector3[T] {
	return Vec3(vec[0]*k, vec[1]*k, vec[2]*k)
}

func (vec Vector3[T]) Normalize() Vector3[T] {
	return vec.Scale(vec.Norm())
}

func (vec Vector3[T]) Dot(other Vector3[T]) T {
	return vec[0]*other[0] + vec[1]*other[1] + vec[2]*other[2]
}

func (vec Vector3[T]) Cross(other Vector3[T]) Vector3[T] {
	var x1, y1, z1 = vec.X(), vec.Y(), vec.Z()
	var x2, y2, z2 = other.X(), other.Y(), other.Z()
	return Vec3(y1*z2-y2*z1, x2*z1-x1*z2, x1*y2-x2*y1)
}

//----------------------------------------------
// measure functions

func (vec Vector3[T]) SquaredLength() T {
	return vec.Dot(vec)
}

func (vec Vector3[T]) Norm() T {
	return T(math.Sqrt(float64(vec.SquaredLength())))
}

// Normp computes p-norm
func (vec Vector3[T]) Normp(p T) T {
	switch p {
	case 0:
		return 3 -
			mathutil.IsZero[T](vec[0]) -
			mathutil.IsZero[T](vec[1]) -
			mathutil.IsZero[T](vec[2])
	case 1:
		return mathutil.Abs(vec[0]) + mathutil.Abs(vec[1]) + mathutil.Abs(vec[2])
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

func (vec Vector3[T]) Vec2() Vector2[T] { return Vec2(vec[0], vec[1]) }
func (vec Vector3[T]) Vec4() Vector4[T] { return Vec4(vec[0], vec[1], vec[2], 1) }
