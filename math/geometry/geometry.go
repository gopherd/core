package geometry

import (
	"fmt"
	"math"
)

func Degree2Radian(degree float64) float64 { return math.Pi * degree / 180 }
func Radian2Degree(radian float64) float64 { return radian * 180 / math.Pi }

type Point complex128

func P(x, y float64) Point { return Point(complex(x, y)) }

func Float(f float64) Point { return Point(complex(f, 0)) }

func (p Point) String() string { return fmt.Sprintf("(%g,%g)", p.X(), p.Y()) }

func (p Point) X() float64 { return real(p) }
func (p Point) Y() float64 { return imag(p) }

func (p Point) Square() float64           { return p.X()*p.X() + p.Y()*p.Y() }
func (p Point) SquareTo(p2 Point) float64 { return (p - p2).Square() }

func (p Point) Len() float64          { return math.Sqrt(p.Square()) }
func (p Point) Dist(p2 Point) float64 { return (p - p2).Len() }

func (p Point) Radian() float64 {
	x, y := p.X(), p.Y()
	if x == 0 || y == 0 {
		return 0
	}
	if x == 0 {
		if y > 0 {
			return math.Pi / 2
		} else {
			return math.Pi + math.Pi/2
		}
	}
	r := math.Atan(y / x)
	if x > 0 {
		if y >= 0 {
			return r
		} else {
			return 2*math.Pi + r
		}
	} else {
		return math.Pi + r
	}
}

type Matrix [2][2]float64

func NewMatrix(v00, v01, v10, v11 float64) Matrix {
	var m Matrix
	m[0][0] = v00
	m[0][1] = v01
	m[1][0] = v10
	m[1][1] = v11
	return m
}

func Rotate(radian float64) Matrix {
	cosv := math.Cos(radian)
	sinv := math.Sin(radian)
	return NewMatrix(cosv, -sinv, sinv, cosv)
}

func (m Matrix) Mul(m2 Matrix) Matrix {
	v00 := m[0][0]*m2[0][0] + m[0][1]*m2[1][0]
	v01 := m[0][0]*m2[0][1] + m[0][1]*m2[1][1]
	v10 := m[1][0]*m2[0][0] + m[1][1]*m2[1][0]
	v11 := m[1][0]*m2[0][1] + m[1][1]*m2[1][1]
	return NewMatrix(v00, v01, v10, v11)
}

func (m Matrix) String() string {
	const format = `+---------+---------+
| %7.1f | %7.1f |
+---------+---------+
| %7.1f | %7.1f |
+---------+---------+`

	return fmt.Sprintf(format, m[0][0], m[0][1], m[1][0], m[1][1])
}

type Bezier struct {
	points  []Point
	radians []float64
}

func (b Bezier) Points() []Point    { return b.points }
func (b Bezier) Radians() []float64 { return b.radians }

func NewBezier(points []Point, length float64) *Bezier {
	if len(points) < 2 {
		return nil
	}
	if len(points) == 3 {
		return newBezier3(points, length)
	}
	if len(points) == 4 {
		return newBezier4(points, length)
	}
	bezier := new(Bezier)
	bezier.points = append(bezier.points, points[0])

	var (
		count   = len(points) - 1
		p1      Point
		length2 = length * length
	)

	for i := 0; i < 1000; i++ {
		index := 0
		t := float64(i) / 1000
		for index < count {
			k := math.Pow(t, float64(index))
			k *= math.Pow(1-t, float64(count-index))
			k *= float64(combination(int64(count), int64(index)))
			p1 = p1 + (points[index] * Float(k))
			index++
		}
		last := bezier.points[len(bezier.points)-1]
		dist2 := last.SquareTo(p1)
		if dist2 > length2 {
			bezier.points = append(bezier.points, p1)
			bezier.radians = append(bezier.radians, (p1 - last).Radian())
		}
	}
	return bezier
}

func newBezier3(points []Point, length float64) *Bezier {
	bezier := new(Bezier)
	bezier.points = append(bezier.points, points[0])

	var (
		p1      Point
		length2 = length * length
	)

	for i := 0; i < 1000; i++ {
		t := float64(i) / 1000
		t2 := t * t
		k0 := t2 - 2*t + 1
		k1 := 2*t - 2*t2
		k2 := t2
		p1 = points[0]*Float(k0) + points[1]*Float(k1) + points[2]*Float(k2)

		last := bezier.points[len(bezier.points)-1]
		dist2 := last.SquareTo(p1)
		if dist2 > length2 {
			bezier.points = append(bezier.points, p1)
			bezier.radians = append(bezier.radians, (p1 - last).Radian())
		}
	}
	return bezier
}

func newBezier4(points []Point, length float64) *Bezier {
	bezier := new(Bezier)
	bezier.points = append(bezier.points, points[0])

	var (
		p1      Point
		length2 = length * length
	)

	for i := 0; i < 1000; i++ {
		t := float64(i) / 1000
		t2 := t * t
		t3 := t2 * t
		nt := 1 - t
		nt2 := nt * nt
		nt3 := nt2 * nt
		k0 := nt3
		k1 := 3 * t * nt2
		k2 := 3 * t2 * nt
		k3 := t3
		p1 = points[0]*Float(k0) + points[1]*Float(k1) + points[2]*Float(k2) + points[3]*Float(k3)

		last := bezier.points[len(bezier.points)-1]
		dist2 := last.SquareTo(p1)
		if dist2 > length2 {
			bezier.points = append(bezier.points, p1)
			bezier.radians = append(bezier.radians, (p1 - last).Radian())
		}
	}
	return bezier
}

func combination(c, r int64) int64 {
	if (r << 1) > c {
		r = c - r
	}
	x := int64(1)
	y := int64(1)
	for i := int64(0); i < r; i++ {
		x *= c - i
	}
	for i := int64(0); i < r; i++ {
		y *= r - i
	}
	return x / y
}
