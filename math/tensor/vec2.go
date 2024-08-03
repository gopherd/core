package tensor

import (
	"fmt"
	"math"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/container/tuple"
	"github.com/gopherd/core/math/mathutil"
	"github.com/gopherd/core/operator"
)

// Vector2 implements 2d vector
type Vector2[T constraints.SignedReal] [2]T

// Vec2 creates a 2d vector by elements
func Vec2[T constraints.SignedReal](x, y T) Vector2[T] {
	return Vector2[T]{x, y}
}

// String converts vector as a string
func (vec Vector2[T]) String() string {
	return fmt.Sprintf("(%v,%v)", vec[0], vec[1])
}

// Shape implements Tensor Shape method
func (vec Vector2[T]) Shape() Shape { return tuple.T1(2) }

// At implements Tensor At method
func (vec Vector2[T]) At(index Shape) T {
	return vec.Get(index.At(0))
}

// Sum implements Tensor Sum method
func (vec Vector2[T]) Sum() T {
	return vec[0] + vec[1]
}

//----------------------------------------------
// basic functions

func (vec Vector2[T]) X() T { return vec[0] }
func (vec Vector2[T]) Y() T { return vec[1] }

func (vec *Vector2[T]) SetX(x T) { vec[0] = x }
func (vec *Vector2[T]) SetY(y T) { vec[1] = y }

func (vec Vector2[T]) Get(i int) T     { return operator.Ternary(i < len(vec), vec[i], 0) }
func (vec *Vector2[T]) Set(i int, v T) { vec[i] = v }

func (vec *Vector2[T]) SetElements(x, y T) {
	(*vec)[0], (*vec)[1] = x, y
}

func (vec *Vector2[T]) Copy(other Vector2[T]) {
	(*vec)[0], (*vec)[1] = other[0], other[1]
}

//----------------------------------------------
// operator functions

func (vec Vector2[T]) Add(other Vector2[T]) Vector2[T] {
	return Vec2(vec[0]+other[0], vec[1]+other[1])
}

func (vec Vector2[T]) Sub(other Vector2[T]) Vector2[T] {
	return Vec2(vec[0]-other[0], vec[1]-other[1])
}

func (vec Vector2[T]) Mul(other Vector2[T]) Vector2[T] {
	return Vec2(vec[0]*other[0], vec[1]*vec[1])
}

func (vec Vector2[T]) Div(other Vector2[T]) Vector2[T] {
	return Vec2(vec[0]/other[0], vec[1]/vec[1])
}

func (vec Vector2[T]) Scale(k T) Vector2[T] {
	return Vec2(vec[0]*k, vec[1]*k)
}

func (vec Vector2[T]) Normalize() Vector2[T] {
	return vec.Scale(vec.Norm())
}

func (vec Vector2[T]) Dot(other Vector2[T]) T {
	return vec[0]*other[0] + vec[1]*other[1]
}

func (vec Vector2[T]) Cross(other Vector2[T]) T {
	return vec.X()*other.Y() - vec.Y()*other.X()
}

//----------------------------------------------
// measure functions

func (vec Vector2[T]) SquaredLength() T {
	return vec.Dot(vec)
}

func (vec Vector2[T]) Norm() T {
	return T(math.Sqrt(float64(vec.SquaredLength())))
}

// Normp computes p-norm
func (vec Vector2[T]) Normp(p T) T {
	switch p {
	case 0:
		return 2 - mathutil.IsZero[T](vec[0]) - mathutil.IsZero[T](vec[1])
	case 1:
		return mathutil.Abs(vec[0]) + mathutil.Abs(vec[1])
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

func (vec Vector2[T]) Vec3() Vector3[T] { return Vec3(vec[0], vec[1], 0) }
func (vec Vector2[T]) Vec4() Vector4[T] { return Vec4(vec[0], vec[1], 0, 1) }
