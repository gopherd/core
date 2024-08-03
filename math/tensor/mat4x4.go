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

// Matrix4 represents a 4x4 matrix
type Matrix4[T constraints.SignedReal] [4 * 4]T

// Zero4x4 creates a zero 4x4 matrix
func Zero4x4[T constraints.SignedReal]() Matrix4[T] {
	return Matrix4[T]{}
}

// One4x4 creates a 4x4 matrix which every element is 1
func One4x4[T constraints.SignedReal]() Matrix4[T] {
	return Matrix4[T]{
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}
}

// Identity4 creates a 4x4 identity matrix
func Identity4[T constraints.SignedReal]() Matrix4[T] {
	return Matrix4[T]{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// String converts matrix as a string
func (mat Matrix4[T]) String() string {
	const dim = 4
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i := 0; i < dim; i++ {
		if i > 0 {
			buf.WriteByte(';')
		}
		buf.WriteByte('(')
		for j := 0; j < dim; j++ {
			if j > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprint(&buf, mat[i+j*dim])
		}
		buf.WriteByte(')')
	}
	buf.WriteByte('}')
	return buf.String()
}

var shape4x4 = tuple.T2(4, 4)

// Shape implements Tensor Shape method
func (mat Matrix4[T]) Shape() Shape {
	return shape4x4
}

// At implements Tensor At method
func (mat Matrix4[T]) At(index Shape) T {
	return mat.Get(index.At(0), index.At(1))
}

// Sum implements Tensor Sum method
func (mat Matrix4[T]) Sum() T {
	var result T
	for i := range mat {
		result += mat[i]
	}
	return result
}

//----------------------------------------------
// basic functions

func (mat Matrix4[T]) Get(i, j int) T {
	return mat[i+j*4]
}

func (mat *Matrix4[T]) Set(i, j int, x T) {
	mat[i+j*4] = x
}

func (mat Matrix4[T]) Elements() []T {
	return mat[:]
}

func (mat *Matrix4[T]) SetElements(n11, n12, n13, n14, n21, n22, n23, n24, n31, n32, n33, n34, n41, n42, n43, n44 T) *Matrix4[T] {
	(*mat)[0], (*mat)[4], (*mat)[8], (*mat)[12] = n11, n12, n13, n14
	(*mat)[1], (*mat)[5], (*mat)[9], (*mat)[13] = n21, n22, n23, n24
	(*mat)[2], (*mat)[6], (*mat)[10], (*mat)[14] = n31, n32, n33, n34
	(*mat)[3], (*mat)[7], (*mat)[11], (*mat)[15] = n41, n42, n43, n44
	return mat
}

//----------------------------------------------
// operator functions

func (mat Matrix4[T]) Transpose() Matrix4[T] {
	const dim = 4
	for i := 0; i < dim-1; i++ {
		for j := i + 1; j < dim; j++ {
			mat[i+j*dim], mat[j+i*dim] = mat[j+i*dim], mat[i+j*dim]
		}
	}
	return mat
}

func (mat Matrix4[T]) Add(other Matrix4[T]) Matrix4[T] {
	for i := range mat {
		mat[i] += other[i]
	}
	return mat
}

func (mat Matrix4[T]) Sub(other Matrix4[T]) Matrix4[T] {
	for i := range mat {
		mat[i] -= other[i]
	}
	return mat
}

func (mat Matrix4[T]) Mul(other Matrix4[T]) Matrix4[T] {
	for i := range mat {
		mat[i] *= other[i]
	}
	return mat
}

func (mat Matrix4[T]) Div(other Matrix4[T]) Matrix4[T] {
	for i := range mat {
		mat[i] /= other[i]
	}
	return mat
}

func (mat Matrix4[T]) Scale(v T) Matrix4[T] {
	for i := range mat {
		mat[i] *= v
	}
	return mat
}

func (mat Matrix4[T]) Normalize() Matrix4[T] {
	return mat.Scale(1 / mat.Norm())
}

func (mat Matrix4[T]) Dot(other Matrix4[T]) Matrix4[T] {
	const dim = 4
	var result Matrix4[T]
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

func (mat Matrix4[T]) DotVec2(vec Vector2[T]) Vector3[T] {
	return mat.DotVec4(vec.Vec4()).Vec3()
}

func (mat Matrix4[T]) DotVec3(vec Vector3[T]) Vector3[T] {
	return mat.DotVec4(vec.Vec4()).Vec3()
}

func (mat Matrix4[T]) DotVec4(vec Vector4[T]) Vector4[T] {
	const dim = 4
	var result Vector4[T]
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			result[i] += mat[i+j*dim] * vec[j]
		}
	}
	return result
}

// based on http://www.euclideanspace.com/maths/algebra/matrix/functions/inverse/fourD/index.htm
func (mat Matrix4[T]) Invert() Matrix4[T] {
	var n11, n21, n31, n41 = mat[0], mat[1], mat[2], mat[3]
	var n12, n22, n32, n42 = mat[4], mat[5], mat[6], mat[7]
	var n13, n23, n33, n43 = mat[8], mat[9], mat[10], mat[11]
	var n14, n24, n34, n44 = mat[12], mat[13], mat[14], mat[15]

	var t11 = n23*n34*n42 - n24*n33*n42 + n24*n32*n43 - n22*n34*n43 - n23*n32*n44 + n22*n33*n44
	var t12 = n14*n33*n42 - n13*n34*n42 - n14*n32*n43 + n12*n34*n43 + n13*n32*n44 - n12*n33*n44
	var t13 = n13*n24*n42 - n14*n23*n42 + n14*n22*n43 - n12*n24*n43 - n13*n22*n44 + n12*n23*n44
	var t14 = n14*n23*n32 - n13*n24*n32 - n14*n22*n33 + n12*n24*n33 + n13*n22*n34 - n12*n23*n34

	var det = n11*t11 + n21*t12 + n31*t13 + n41*t14

	if det == 0 {
		return Matrix4[T]{}
	}

	var detInv = 1 / det

	mat[0] = t11 * detInv
	mat[1] = (n24*n33*n41 - n23*n34*n41 - n24*n31*n43 + n21*n34*n43 + n23*n31*n44 - n21*n33*n44) * detInv
	mat[2] = (n22*n34*n41 - n24*n32*n41 + n24*n31*n42 - n21*n34*n42 - n22*n31*n44 + n21*n32*n44) * detInv
	mat[3] = (n23*n32*n41 - n22*n33*n41 - n23*n31*n42 + n21*n33*n42 + n22*n31*n43 - n21*n32*n43) * detInv

	mat[4] = t12 * detInv
	mat[5] = (n13*n34*n41 - n14*n33*n41 + n14*n31*n43 - n11*n34*n43 - n13*n31*n44 + n11*n33*n44) * detInv
	mat[6] = (n14*n32*n41 - n12*n34*n41 - n14*n31*n42 + n11*n34*n42 + n12*n31*n44 - n11*n32*n44) * detInv
	mat[7] = (n12*n33*n41 - n13*n32*n41 + n13*n31*n42 - n11*n33*n42 - n12*n31*n43 + n11*n32*n43) * detInv

	mat[8] = t13 * detInv
	mat[9] = (n14*n23*n41 - n13*n24*n41 - n14*n21*n43 + n11*n24*n43 + n13*n21*n44 - n11*n23*n44) * detInv
	mat[10] = (n12*n24*n41 - n14*n22*n41 + n14*n21*n42 - n11*n24*n42 - n12*n21*n44 + n11*n22*n44) * detInv
	mat[11] = (n13*n22*n41 - n12*n23*n41 - n13*n21*n42 + n11*n23*n42 + n12*n21*n43 - n11*n22*n43) * detInv

	mat[12] = t14 * detInv
	mat[13] = (n13*n24*n31 - n14*n23*n31 + n14*n21*n33 - n11*n24*n33 - n13*n21*n34 + n11*n23*n34) * detInv
	mat[14] = (n14*n22*n31 - n12*n24*n31 - n14*n21*n32 + n11*n24*n32 + n12*n21*n34 - n11*n22*n34) * detInv
	mat[15] = (n12*n23*n31 - n13*n22*n31 + n13*n21*n32 - n11*n23*n32 - n12*n21*n33 + n11*n22*n33) * detInv

	return mat
}

//----------------------------------------------
// measure functions

func (mat Matrix4[T]) Determinant() T {
	var n11, n12, n13, n14 = mat[0], mat[4], mat[8], mat[12]
	var n21, n22, n23, n24 = mat[1], mat[5], mat[9], mat[13]
	var n31, n32, n33, n34 = mat[2], mat[6], mat[10], mat[14]
	var n41, n42, n43, n44 = mat[3], mat[7], mat[11], mat[15]

	return n41*(+n14*n23*n32-n13*n24*n32-n14*n22*n33+n12*n24*n33+n13*n22*n34-n12*n23*n34) +
		n42*(+n11*n23*n34-n11*n24*n33+n14*n21*n33-n13*n21*n34+n13*n24*n31-n14*n23*n31) +
		n43*(+n11*n24*n32-n11*n22*n34-n14*n21*n32+n12*n21*n34+n14*n22*n31-n12*n24*n31) +
		n44*(-n13*n22*n31-n11*n23*n32+n11*n22*n33+n13*n21*n32-n12*n21*n33+n12*n23*n31)
}

func (mat Matrix4[T]) SquaredLength() T {
	return mat.Mul(mat).Sum()
}

func (mat Matrix4[T]) Norm() T {
	return T(math.Sqrt(float64(mat.SquaredLength())))
}

func (mat Matrix4[T]) Normp(p T) T {
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

//----------------------------------------------
// geometry functions

func (mat Matrix4[T]) GetPosition() Vector3[T] {
	return Vec3(mat[12], mat[13], mat[14])
}

func (mat *Matrix4[T]) SetPosition(pos Vector3[T]) *Matrix4[T] {
	(*mat)[12] = pos.X()
	(*mat)[13] = pos.Y()
	(*mat)[14] = pos.Z()
	return mat
}

// http://www.euclideanspace.com/maths/geometry/rotations/conversions/matrixToQuaternion/index.htm
func (mat Matrix4[T]) GetQuaternion() Vector4[T] {
	// assumes the upper 3x3 of m is a pure rotation matrix (i.e, unscaled)
	var m11, m12, m13 = float64(mat[0]), float64(mat[4]), float64(mat[8])
	var m21, m22, m23 = float64(mat[1]), float64(mat[5]), float64(mat[9])
	var m31, m32, m33 = float64(mat[2]), float64(mat[6]), float64(mat[10])
	var trace = m11 + m22 + m33
	var x, y, z, w float64

	if trace > 0 {
		var s = 0.5 / math.Sqrt(trace+1.0)
		w = 0.25 / s
		x = (m32 - m23) * s
		y = (m13 - m31) * s
		z = (m21 - m12) * s
	} else if m11 > m22 && m11 > m33 {
		var s = 2.0 * math.Sqrt(1.0+m11-m22-m33)
		w = (m32 - m23) / s
		x = 0.25 * s
		y = (m12 + m21) / s
		z = (m13 + m31) / s
	} else if m22 > m33 {
		var s = 2.0 * math.Sqrt(1.0+m22-m11-m33)
		w = (m13 - m31) / s
		x = (m12 + m21) / s
		y = 0.25 * s
		z = (m23 + m32) / s
	} else {
		var s = 2.0 * math.Sqrt(1.0+m33-m11-m22)
		w = (m21 - m12) / s
		x = (m13 + m31) / s
		y = (m23 + m32) / s
		z = 0.25 * s
	}

	return Vec4(T(x), T(y), T(z), T(w))
}

func (mat *Matrix4[T]) MakeIdentity() *Matrix4[T] {
	return mat.SetElements(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
}

func (mat *Matrix4[T]) MakeZero() *Matrix4[T] {
	return mat.SetElements(
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	)
}

func (mat *Matrix4[T]) MakeFromMatrix3(m Matrix3[T]) *Matrix4[T] {
	return mat.SetElements(
		m[0], m[3], m[6], 0,
		m[1], m[4], m[7], 0,
		m[2], m[5], m[8], 0,
		0, 0, 0, 1,
	)
}

func (mat *Matrix4[T]) MakeTranslation(vec Vector3[T]) *Matrix4[T] {
	return mat.SetElements(
		1, 0, 0, vec.X(),
		0, 1, 0, vec.Y(),
		0, 0, 1, vec.Z(),
		0, 0, 0, 1,
	)
}

func (mat *Matrix4[T]) MakeRotationX(theta T) *Matrix4[T] {
	var s0, c0 = math.Sincos(float64(theta))
	var s, c = T(s0), T(c0)
	return mat.SetElements(
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
		0, 0, 0, 1,
	)
}

func (mat *Matrix4[T]) MakeRotationY(theta T) *Matrix4[T] {
	var s0, c0 = math.Sincos(float64(theta))
	var s, c = T(s0), T(c0)
	return mat.SetElements(
		c, 0, s, 0,
		0, 1, 0, 0,
		-s, 0, c, 0,
		0, 0, 0, 1,
	)

}

func (mat *Matrix4[T]) MakeRotationZ(theta T) *Matrix4[T] {
	var s0, c0 = math.Sincos(float64(theta))
	var s, c = T(s0), T(c0)
	return mat.SetElements(
		c, -s, 0, 0,
		s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
}

// Based on http://www.gamedev.net/reference/articles/article1199.asp
func (mat *Matrix4[T]) MakeRotationAxis(axis Vector3[T], angle T) *Matrix4[T] {
	var s0, c0 = math.Sincos(float64(angle))
	var s, c = T(s0), T(c0)
	var t = 1 - c
	var x, y, z = axis.X(), axis.Y(), axis.Z()
	var tx, ty = t * x, t * y
	return mat.SetElements(
		tx*x+c, tx*y-s*z, tx*z+s*y, 0,
		tx*y+s*z, ty*y+c, ty*z-s*x, 0,
		tx*z-s*y, ty*z+s*x, t*z*z+c, 0,
		0, 0, 0, 1,
	)
}

func (mat *Matrix4[T]) MakeScale(vec Vector3[T]) *Matrix4[T] {
	return mat.SetElements(
		vec.X(), 0, 0, 0,
		0, vec.Y(), 0, 0,
		0, 0, vec.Z(), 0,
		0, 0, 0, 1,
	)
}

func (mat *Matrix4[T]) MakeShear(xy, xz, yx, yz, zx, zy T) *Matrix4[T] {
	return mat.SetElements(
		1, yx, zx, 0,
		xy, 1, zy, 0,
		xz, yz, 1, 0,
		0, 0, 0, 1,
	)
}

func (mat *Matrix4[T]) Compose(position Vector3[T], quaternion Vector4[T], scale Vector3[T]) *Matrix4[T] {
	var x, y, z, w = quaternion.X(), quaternion.Y(), quaternion.Z(), quaternion.W()
	var x2, y2, z2 = x + x, y + y, z + z
	var xx, xy, xz = x * x2, x * y2, x * z2
	var yy, yz, zz = y * y2, y * z2, z * z2
	var wx, wy, wz = w * x2, w * y2, w * z2
	var sx, sy, sz = scale.X(), scale.Y(), scale.Z()

	(*mat)[0] = (1 - (yy + zz)) * sx
	(*mat)[1] = (xy + wz) * sx
	(*mat)[2] = (xz - wy) * sx
	(*mat)[3] = 0

	(*mat)[4] = (xy - wz) * sy
	(*mat)[5] = (1 - (xx + zz)) * sy
	(*mat)[6] = (yz + wx) * sy
	(*mat)[7] = 0

	(*mat)[8] = (xz + wy) * sz
	(*mat)[9] = (yz - wx) * sz
	(*mat)[10] = (1 - (xx + yy)) * sz
	(*mat)[11] = 0

	(*mat)[12] = position.X()
	(*mat)[13] = position.Y()
	(*mat)[14] = position.Z()
	(*mat)[15] = 1

	return mat
}

func (mat Matrix4[T]) Decompose() (position Vector3[T], quaternion Vector4[T], scale Vector3[T]) {
	var sx = Vec3(mat[0], mat[1], mat[2]).Norm()
	var sy = Vec3(mat[4], mat[5], mat[6]).Norm()
	var sz = Vec3(mat[8], mat[9], mat[10]).Norm()

	// if determine is negative, we need to invert one scale
	var det = mat.Determinant()
	if det < 0 {
		sx = -sx
	}

	position.SetElements(mat[12], mat[13], mat[14])

	var invSX = 1 / sx
	var invSY = 1 / sy
	var invSZ = 1 / sz

	mat[0] *= invSX
	mat[1] *= invSX
	mat[2] *= invSX

	mat[4] *= invSY
	mat[5] *= invSY
	mat[6] *= invSY

	mat[8] *= invSZ
	mat[9] *= invSZ
	mat[10] *= invSZ

	quaternion = mat.GetQuaternion()
	scale.SetElements(sx, sy, sz)

	return
}

func (mat *Matrix4[T]) MakePerspective(left, right, top, bottom, near, far T) *Matrix4[T] {
	var x = 2 * near / (right - left)
	var y = 2 * near / (top - bottom)
	var a = (right + left) / (right - left)
	var b = (top + bottom) / (top - bottom)
	var c = -(far + near) / (far - near)
	var d = -2 * far * near / (far - near)

	return mat.SetElements(
		x, 0, a, 0,
		0, y, b, 0,
		0, 0, c, d,
		0, 0, -1, 0,
	)
}

func (mat *Matrix4[T]) MakeOrthographic(left, right, top, bottom, near, far T) *Matrix4[T] {
	var w = 1.0 / (right - left)
	var h = 1.0 / (top - bottom)
	var p = 1.0 / (far - near)
	var x = (right + left) * w
	var y = (top + bottom) * h
	var z = (far + near) * p

	return mat.SetElements(
		2*w, 0, 0, -x,
		0, 2*h, 0, -y,
		0, 0, -2*p, -z,
		0, 0, 0, 1,
	)
}
