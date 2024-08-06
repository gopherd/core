// Package mathutil provides various mathematical utility functions.
package mathutil

import (
	"math"

	"github.com/gopherd/core/constraints"
)

// Abs returns the absolute value of x.
func Abs[T constraints.SignedReal](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Predict returns 1 if ok is true, otherwise 0.
func Predict[T constraints.Integer | constraints.Float](ok bool) T {
	if ok {
		return 1
	}
	return 0
}

// Clamp restricts x to the range [min, max].
func Clamp[T constraints.Ordered](x, min, max T) T {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// EuclideanModulo computes the Euclidean modulo of x % y.
func EuclideanModulo[T constraints.Float](x, y T) T {
	return T(math.Mod(math.Mod(float64(x), float64(y))+float64(y), float64(y)))
}

// MapLinear performs linear mapping from range [a1, a2] to range [b1, b2].
func MapLinear[T constraints.Field](x, a1, a2, b1, b2 T) T {
	return b1 + (x-a1)*(b2-b1)/(a2-a1)
}

// Lerp performs linear interpolation between x and y based on t.
func Lerp[T constraints.Field](x, y, t T) T {
	return (1-t)*x + t*y
}

// InverseLerp calculates the inverse of linear interpolation.
func InverseLerp[T constraints.Field](x, y, value T) T {
	if x != y {
		return (value - x) / (y - x)
	}
	return 0
}

// Damp performs frame rate independent damping.
func Damp[T constraints.Float](x, y, lambda, dt T) T {
	return Lerp(x, y, 1-T(math.Exp(-float64(lambda*dt))))
}

// PingPong calculates a value that ping-pongs between 0 and length.
func PingPong[T constraints.Float](x, length T) T {
	return length - Abs(EuclideanModulo(x, length*2)-length)
}

// SmoothStep performs smooth interpolation between min and max.
func SmoothStep[T constraints.Float](x, min, max T) T {
	if x <= min {
		return 0
	}
	if x >= max {
		return 1
	}
	x = (x - min) / (max - min)
	return x * x * (3 - 2*x)
}

// SmoothStepFunc applies a custom function to the smoothstep interpolation.
func SmoothStepFunc[T constraints.Float](x, min, max T, fn func(T) T) T {
	if x <= min {
		return 0
	}
	if x >= max {
		return 1
	}
	x = (x - min) / (max - min)
	return fn(x)
}

// IsPowerOfTwo checks if the given value is a power of two.
func IsPowerOfTwo[T constraints.Integer](value T) bool {
	return value > 0 && (value&(value-1)) == 0
}

// UpperPow2 returns the smallest power of 2 greater than or equal to n.
func UpperPow2(n int) int {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}

// Deg2Rad converts degrees to radians.
func Deg2Rad[T constraints.Float](deg T) T {
	return deg * T(math.Pi) / 180
}

// Rad2Deg converts radians to degrees.
func Rad2Deg[T constraints.Float](rad T) T {
	return rad * 180 / T(math.Pi)
}

// UnaryFn represents a unary function.
type UnaryFn[T constraints.Number] func(T) T

// Add returns a new UnaryFn that adds the results of two UnaryFn.
func (f UnaryFn[T]) Add(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) + f2(x)
	}
}

// Sub returns a new UnaryFn that subtracts the result of f2 from f.
func (f UnaryFn[T]) Sub(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) - f2(x)
	}
}

// Mul returns a new UnaryFn that multiplies the results of two UnaryFn.
func (f UnaryFn[T]) Mul(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) * f2(x)
	}
}

// Div returns a new UnaryFn that divides the result of f by f2.
func (f UnaryFn[T]) Div(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) / f2(x)
	}
}

// Constant returns a UnaryFn that always returns c.
func Constant[T constraints.Number](c T) UnaryFn[T] {
	return func(T) T { return c }
}

// KSigmoid returns a UnaryFn that applies a sigmoid function with slope k.
func KSigmoid[T constraints.Real](k T) UnaryFn[T] {
	return func(x T) T { return Sigmoid(k * x) }
}

// KSigmoidPrime returns a UnaryFn that applies the derivative of a sigmoid function with slope k.
func KSigmoidPrime[T constraints.Real](k T) UnaryFn[T] {
	return func(x T) T { return SigmoidPrime(k*x) * k }
}

// Scale returns a UnaryFn that scales its input by k.
func Scale[T constraints.Number](k T) UnaryFn[T] {
	return func(x T) T { return k * x }
}

// Offset returns a UnaryFn that adds b to its input.
func Offset[T constraints.Number](b T) UnaryFn[T] {
	return func(x T) T { return x + b }
}

// Affine returns a UnaryFn that applies an affine transformation (kx + b).
func Affine[T constraints.Number](k, b T) UnaryFn[T] {
	return func(x T) T { return k*x + b }
}

// Power returns a UnaryFn that raises its input to the power of p.
func Power[T constraints.Real](p T) UnaryFn[T] {
	return func(x T) T { return T(math.Pow(float64(x), float64(p))) }
}

// Zero always returns 0.
func Zero[T constraints.Number](T) T {
	return 0
}

// One always returns 1.
func One[T constraints.Number](T) T {
	return 1
}

// Identity returns its input unchanged.
func Identity[T constraints.Number](x T) T {
	return x
}

// Square returns the square of its input.
func Square[T constraints.Number](x T) T {
	return x * x
}

// IsZero returns 1 if the input is zero, otherwise 0.
func IsZero[T constraints.SignedReal](x T) T {
	if x == 0 {
		return 1
	}
	return 0
}

// Sign returns the sign of the input (-1, 0, or 1).
func Sign[T constraints.SignedReal](x T) T {
	switch {
	case x < 0:
		return -1
	case x > 0:
		return 1
	default:
		return 0
	}
}

// Sigmoid applies the sigmoid function to the input.
func Sigmoid[T constraints.Real](x T) T {
	return T(1.0 / (1.0 + math.Exp(-float64(x))))
}

// SigmoidPrime applies the derivative of the sigmoid function to the input.
func SigmoidPrime[T constraints.Real](x T) T {
	sx := Sigmoid(x)
	return sx * (1 - sx)
}

// BinaryFn represents a binary function.
type BinaryFn[T constraints.Number] func(x, y T) T

// Add returns a new BinaryFn that adds the results of two BinaryFn.
func (f BinaryFn[T]) Add(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) + f2(x, y)
	}
}

// Sub returns a new BinaryFn that subtracts the result of f2 from f.
func (f BinaryFn[T]) Sub(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) - f2(x, y)
	}
}

// Mul returns a new BinaryFn that multiplies the results of two BinaryFn.
func (f BinaryFn[T]) Mul(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) * f2(x, y)
	}
}

// Div returns a new BinaryFn that divides the result of f by f2.
func (f BinaryFn[T]) Div(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) / f2(x, y)
	}
}

// Add returns the sum of x and y.
func Add[T constraints.Number](x, y T) T { return x + y }

// Sub returns the difference of x and y.
func Sub[T constraints.Number](x, y T) T { return x - y }

// Mul returns the product of x and y.
func Mul[T constraints.Number](x, y T) T { return x * y }

// Div returns the quotient of x and y.
func Div[T constraints.Number](x, y T) T { return x / y }

// Pow returns x raised to the power of y.
func Pow[T constraints.Real](x, y T) T { return T(math.Pow(float64(x), float64(y))) }

// ClampedLerp performs a linear interpolation and clamps the result.
func ClampedLerp[T constraints.Float](x, y, t, min, max T) T {
	return Clamp(Lerp(x, y, t), min, max)
}
