package tensor

import (
	"math"
	"strings"

	"github.com/gopherd/core/constraints"
	"github.com/gopherd/core/math/mathutil"
)

type EulerRotationOrder int8

const (
	EulerRotationOrderXYZ EulerRotationOrder = iota
	EulerRotationOrderXZY
	EulerRotationOrderYXZ
	EulerRotationOrderYZX
	EulerRotationOrderZXY
	EulerRotationOrderZYX
)

func (order EulerRotationOrder) String() string {
	switch order {
	case EulerRotationOrderXYZ:
		return "XYZ"
	case EulerRotationOrderXZY:
		return "XZY"
	case EulerRotationOrderYXZ:
		return "YXZ"
	case EulerRotationOrderYZX:
		return "YZX"
	case EulerRotationOrderZXY:
		return "ZXY"
	case EulerRotationOrderZYX:
		return "ZYX"
	default:
		return "UNK"
	}
}

type Euler[T constraints.Float] struct {
	Vector3[T]
	Order EulerRotationOrder // XYZ, XZY, YXZ, YZX, ZXY, ZYX
}

func (euler Euler[T]) String() string {
	var sb strings.Builder
	sb.WriteString("{rotation:")
	sb.WriteString(euler.Vector3.String())
	sb.WriteString(",order:")
	sb.WriteString(euler.Order.String())
	sb.WriteByte('}')
	return sb.String()
}

func (euler *Euler[T]) SetFromRotationMatrix(mat Matrix4[T]) *Euler[T] {
	const one = 0.9999999
	// assumes the upper 3x3 of m is a pure rotation matrix (i.e, unscaled)
	var m11, m12, m13 = float64(mat[0]), float64(mat[4]), float64(mat[8])
	var m21, m22, m23 = float64(mat[1]), float64(mat[5]), float64(mat[9])
	var m31, m32, m33 = float64(mat[2]), float64(mat[6]), float64(mat[10])
	var x, y, z float64
	switch euler.Order {
	case EulerRotationOrderYXZ:
		x = math.Asin(-mathutil.Clamp(m23, -1, 1))
		if mathutil.Abs(m23) < one {
			y = math.Atan2(m13, m33)
			z = math.Atan2(m21, m22)
		} else {
			y = math.Atan2(-m31, m11)
			z = 0
		}
	case EulerRotationOrderZXY:
		x = math.Asin(mathutil.Clamp(m32, -1, 1))
		if math.Abs(m32) < one {
			y = math.Atan2(-m31, m33)
			z = math.Atan2(-m12, m22)
		} else {
			y = 0
			z = math.Atan2(m21, m11)
		}
	case EulerRotationOrderZYX:
		y = math.Asin(-mathutil.Clamp(m31, -1, 1))
		if math.Abs(m31) < one {
			x = math.Atan2(m32, m33)
			z = math.Atan2(m21, m11)
		} else {
			x = 0
			z = math.Atan2(-m12, m22)
		}
	case EulerRotationOrderYZX:
		z = math.Asin(mathutil.Clamp(m21, -1, 1))
		if math.Abs(m21) < one {
			x = math.Atan2(-m23, m22)
			y = math.Atan2(-m31, m11)
		} else {
			x = 0
			y = math.Atan2(m13, m33)
		}
	case EulerRotationOrderXZY:
		z = math.Asin(-mathutil.Clamp(m12, -1, 1))
		if math.Abs(m12) < one {
			x = math.Atan2(m32, m22)
			y = math.Atan2(m13, m11)
		} else {
			x = math.Atan2(-m23, m33)
			y = 0
		}
	default:
		y = math.Asin(mathutil.Clamp(m13, -1, 1))
		if math.Abs(m13) < one {
			x = math.Atan2(-m23, m33)
			z = math.Atan2(-m12, m11)
		} else {
			x = math.Atan2(m32, m22)
			z = 0
		}
	}
	euler.SetElements(T(x), T(y), T(z))
	return euler
}
