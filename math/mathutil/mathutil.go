package mathutil

import (
	"math"

	"github.com/gopherd/core/constraints"
)

// Min returns mininum value
func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Max returns maxinum value
func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// Minmax returns ordered values
func Minmax[T constraints.Ordered](x, y T) (min, max T) {
	if x < y {
		return x, y
	}
	return y, x
}

// Abs returns abs of x
func Abs[T constraints.SignedReal](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Predict returns 1 if ok, otherwise 0
func Predict[T constraints.Number](ok bool) T {
	if ok {
		return 1
	}
	return 0
}

// Clamp clamps x into range [min, max]
func Clamp[T constraints.Ordered](x, min, max T) T {
	return Max(min, Min(max, x))
}

// EuclideanModulo computes euclidean modulo: x % y
func EuclideanModulo[T constraints.Float](x, y T) T {
	var x64, y64 = float64(x), float64(y)
	return T(math.Mod(math.Mod(x64, y64)+y64, y64))
}

// MapLinear mapping from range <a1, a2> to range <b1, b2>
func MapLinear[T constraints.Field](x, a1, a2, b1, b2 T) T {
	return b1 + (x-a1)*(b2-b1)/(a2-a1)
}

// https://en.wikipedia.org/wiki/Linear_interpolation
func Lerp[T constraints.Field](x, y, t T) T {
	return (1-t)*x + t*y
}

// https://www.gamedev.net/tutorials/programming/general-and-gameplay-programming/inverse-lerp-a-super-useful-yet-often-overlooked-function-r5230/
func InverseLerp[T constraints.Field](x, y, value T) T {
	if x == y {
		return 0
	}
	return (value - x) / (y - x)
}

// http://www.rorydriscoll.com/2016/03/07/frame-rate-independent-damping-using-lerp/
func Damp[T constraints.Float](x, y, lambda, dt T) T {
	return Lerp(x, y, 1-T(math.Exp(float64(-lambda*dt))))
}

// https://www.desmos.com/calculator/vcsjnyz7x4
func PingPong[T constraints.Float](x, length T) T {
	return length - Abs(EuclideanModulo(x, length*2)-length)
}

// http://en.wikipedia.org/wiki/Smoothstep
func SmoothStep[T constraints.Float](x, min, max T) T {
	if x <= min {
		return 0
	} else if x >= max {
		return 1
	}
	x = (x - min) / (max - min)
	return x * x * (3 - 2*x)
}

func SmoothStepFunc[T constraints.Float](x, min, max T, fn func(T) T) T {
	if x <= min {
		return 0
	} else if x >= max {
		return 1
	}
	x = (x - min) / (max - min)
	return fn(x)
}

func IsPowerOfTwo[T constraints.Integer](value T) bool {
	return (value&(value-1)) == 0 && value != 0
}

func CeilPowerOfTwo[T constraints.Integer](value T) T {
	return T(math.Pow(2, math.Ceil(math.Log(float64(value))/math.Ln2)))
}

func FloorPowerOfTwo[T constraints.Integer](value T) T {
	return T(math.Pow(2, math.Floor(math.Log(float64(value))/math.Ln2)))
}

const deg2Rad = math.Pi / 180
const rad2Deg = 180 / math.Pi

func Deg2Rad[T constraints.Float](deg T) T {
	return deg * deg2Rad
}

func Rad2Deg[T constraints.Float](rad T) T {
	return rad * rad2Deg
}

type UnaryFn[T constraints.Number] func(T) T

func (f UnaryFn[T]) Add(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) + f2(x)
	}
}

func (f UnaryFn[T]) Sub(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) - f2(x)
	}
}

func (f UnaryFn[T]) Mul(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) * f2(x)
	}
}

func (f UnaryFn[T]) Div(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) / f2(x)
	}
}

func Constant[T constraints.Number](c T) UnaryFn[T] {
	return func(x T) T { return c }
}

func KSigmoid[T constraints.Real](k T) UnaryFn[T] {
	return func(x T) T { return Sigmoid(k * x) }
}

func KSigmoidPrime[T constraints.Real](k T) UnaryFn[T] {
	return func(x T) T { return SigmoidPrime(k*x) * k }
}

func Scale[T constraints.Number](k T) UnaryFn[T] {
	return func(x T) T { return k * x }
}

func Offset[T constraints.Number](b T) UnaryFn[T] {
	return func(x T) T { return x + b }
}

func Affine[T constraints.Number](k, b T) UnaryFn[T] {
	return func(x T) T { return k*x + b }
}

func Power[T constraints.Real](p T) UnaryFn[T] {
	return func(x T) T { return T(math.Pow(float64(x), float64(p))) }
}

func Zero[T constraints.Number](x T) T {
	return 0
}

func One[T constraints.Number](x T) T {
	return 1
}

func Identity[T constraints.Number](x T) T {
	return x
}

func Square[T constraints.Number](x T) T {
	return x * x
}

func IsZero[T constraints.SignedReal](x T) T {
	if x == 0 {
		return 1
	}
	return 0
}

func Sign[T constraints.SignedReal](x T) T {
	if x == 0 {
		return 0
	}
	if x > 0 {
		return 1
	}
	return -1
}

func Sigmoid[T constraints.Real](x T) T {
	return T(1.0 / (1.0 + math.Exp(-float64(x))))
}

func SigmoidPrime[T constraints.Real](x T) T {
	x = Sigmoid(x)
	return x * (1 - x)
}

type BinaryFn[T constraints.Number] func(x, y T) T

func (f BinaryFn[T]) Add(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) + f2(x, y)
	}
}

func (f BinaryFn[T]) Sub(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) - f2(x, y)
	}
}

func (f BinaryFn[T]) Mul(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) * f2(x, y)
	}
}

func (f BinaryFn[T]) Div(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) / f2(x, y)
	}
}

func Add[T constraints.Number](x, y T) T { return x + y }
func Sub[T constraints.Number](x, y T) T { return x - y }
func Mul[T constraints.Number](x, y T) T { return x * y }
func Div[T constraints.Number](x, y T) T { return x / y }
func Pow[T constraints.Real](x, y T) T   { return T(math.Pow(float64(x), float64(y))) }
