package mathutil_test

import (
	"math"
	"testing"

	"github.com/gopherd/core/math/mathutil"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"positive", 5.0, 5.0},
		{"negative", -3.0, 3.0},
		{"zero", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.Abs(tt.x); got != tt.want {
				t.Errorf("Abs(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestPredict(t *testing.T) {
	tests := []struct {
		name string
		ok   bool
		want int
	}{
		{"true", true, 1},
		{"false", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.Predict[int](tt.ok); got != tt.want {
				t.Errorf("Predict(%v) = %v, want %v", tt.ok, got, tt.want)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		min  float64
		max  float64
		want float64
	}{
		{"within range", 5.0, 0.0, 10.0, 5.0},
		{"below min", -1.0, 0.0, 10.0, 0.0},
		{"above max", 11.0, 0.0, 10.0, 10.0},
		{"equal to min", 0.0, 0.0, 10.0, 0.0},
		{"equal to max", 10.0, 0.0, 10.0, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.Clamp(tt.x, tt.min, tt.max); got != tt.want {
				t.Errorf("Clamp(%v, %v, %v) = %v, want %v", tt.x, tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestEuclideanModulo(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		y    float64
		want float64
	}{
		{"positive x, positive y", 7.0, 3.0, 1.0},
		{"negative x, positive y", -7.0, 3.0, 2.0},
		{"positive x, negative y", 7.0, -3.0, -2.0},
		{"negative x, negative y", -7.0, -3.0, -1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.EuclideanModulo(tt.x, tt.y); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("EuclideanModulo(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestMapLinear(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		a1   float64
		a2   float64
		b1   float64
		b2   float64
		want float64
	}{
		{"middle of range", 5.0, 0.0, 10.0, 0.0, 100.0, 50.0},
		{"start of range", 0.0, 0.0, 10.0, 0.0, 100.0, 0.0},
		{"end of range", 10.0, 0.0, 10.0, 0.0, 100.0, 100.0},
		{"reverse range", 5.0, 0.0, 10.0, 100.0, 0.0, 50.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.MapLinear(tt.x, tt.a1, tt.a2, tt.b1, tt.b2); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("MapLinear(%v, %v, %v, %v, %v) = %v, want %v", tt.x, tt.a1, tt.a2, tt.b1, tt.b2, got, tt.want)
			}
		})
	}
}

func TestLerp(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		y    float64
		t    float64
		want float64
	}{
		{"midpoint", 0.0, 10.0, 0.5, 5.0},
		{"start", 0.0, 10.0, 0.0, 0.0},
		{"end", 0.0, 10.0, 1.0, 10.0},
		{"quarter", 0.0, 10.0, 0.25, 2.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.Lerp(tt.x, tt.y, tt.t); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Lerp(%v, %v, %v) = %v, want %v", tt.x, tt.y, tt.t, got, tt.want)
			}
		})
	}
}

func TestInverseLerp(t *testing.T) {
	tests := []struct {
		name  string
		x     float64
		y     float64
		value float64
		want  float64
	}{
		{"midpoint", 0.0, 10.0, 5.0, 0.5},
		{"start", 0.0, 10.0, 0.0, 0.0},
		{"end", 0.0, 10.0, 10.0, 1.0},
		{"quarter", 0.0, 10.0, 2.5, 0.25},
		{"equal x and y", 5.0, 5.0, 5.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.InverseLerp(tt.x, tt.y, tt.value); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("InverseLerp(%v, %v, %v) = %v, want %v", tt.x, tt.y, tt.value, got, tt.want)
			}
		})
	}
}

func TestDamp(t *testing.T) {
	tests := []struct {
		name   string
		x      float64
		y      float64
		lambda float64
		dt     float64
		want   float64
	}{
		{"slow damping", 0.0, 10.0, 0.1, 1.0, 0.9516},
		{"fast damping", 0.0, 10.0, 1.0, 1.0, 6.3212},
		{"no time passed", 0.0, 10.0, 1.0, 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.Damp(tt.x, tt.y, tt.lambda, tt.dt); math.Abs(got-tt.want) > 1e-4 {
				t.Errorf("Damp(%v, %v, %v, %v) = %v, want %v", tt.x, tt.y, tt.lambda, tt.dt, got, tt.want)
			}
		})
	}
}

func TestPingPong(t *testing.T) {
	tests := []struct {
		name   string
		x      float64
		length float64
		want   float64
	}{
		{"within range", 1.5, 2.0, 1.5},
		{"over range", 2.5, 2.0, 1.5},
		{"double range", 4.5, 2.0, 0.5},
		{"negative input", -0.5, 2.0, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.PingPong(tt.x, tt.length); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("PingPong(%v, %v) = %v, want %v", tt.x, tt.length, got, tt.want)
			}
		})
	}
}

func TestSmoothStep(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		min  float64
		max  float64
		want float64
	}{
		{"below min", 0.0, 1.0, 2.0, 0.0},
		{"above max", 3.0, 1.0, 2.0, 1.0},
		{"midpoint", 1.5, 1.0, 2.0, 0.5},
		{"quarter", 1.25, 1.0, 2.0, 0.15625},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.SmoothStep(tt.x, tt.min, tt.max); math.Abs(got-tt.want) > 1e-5 {
				t.Errorf("SmoothStep(%v, %v, %v) = %v, want %v", tt.x, tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestSmoothStepFunc(t *testing.T) {
	customFunc := func(x float64) float64 {
		return x * x
	}

	tests := []struct {
		name string
		x    float64
		min  float64
		max  float64
		want float64
	}{
		{"below min", 0.0, 1.0, 2.0, 0.0},
		{"above max", 3.0, 1.0, 2.0, 1.0},
		{"midpoint", 1.5, 1.0, 2.0, 0.25},
		{"quarter", 1.25, 1.0, 2.0, 0.0625},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.SmoothStepFunc(tt.x, tt.min, tt.max, customFunc); math.Abs(got-tt.want) > 1e-5 {
				t.Errorf("SmoothStepFunc(%v, %v, %v, customFunc) = %v, want %v", tt.x, tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"zero", 0, false},
		{"one", 1, true},
		{"two", 2, true},
		{"four", 4, true},
		{"not power of two", 6, false},
		{"large power of two", 1024, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.IsPowerOfTwo(tt.value); got != tt.want {
				t.Errorf("IsPowerOfTwo(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestUpperPow2(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"zero", 0, 0},
		{"one", 1, 1},
		{"two", 2, 2},
		{"three", 3, 4},
		{"five", 5, 8},
		{"large number", 1000, 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.UpperPow2(tt.n); got != tt.want {
				t.Errorf("UpperPow2(%v) = %v, want %v",
					tt.n, got, tt.want)
			}
		})
	}
}

func TestDeg2Rad(t *testing.T) {
	tests := []struct {
		name string
		deg  float64
		want float64
	}{
		{"zero", 0, 0},
		{"45 degrees", 45, math.Pi / 4},
		{"90 degrees", 90, math.Pi / 2},
		{"180 degrees", 180, math.Pi},
		{"360 degrees", 360, 2 * math.Pi},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.Deg2Rad(tt.deg); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Deg2Rad(%v) = %v, want %v", tt.deg, got, tt.want)
			}
		})
	}
}

func TestRad2Deg(t *testing.T) {
	tests := []struct {
		name string
		rad  float64
		want float64
	}{
		{"zero", 0, 0},
		{"pi/4", math.Pi / 4, 45},
		{"pi/2", math.Pi / 2, 90},
		{"pi", math.Pi, 180},
		{"2pi", 2 * math.Pi, 360},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.Rad2Deg(tt.rad); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Rad2Deg(%v) = %v, want %v", tt.rad, got, tt.want)
			}
		})
	}
}

func TestUnaryFn(t *testing.T) {
	square := func(x float64) float64 { return x * x }
	double := func(x float64) float64 { return 2 * x }

	t.Run("Add", func(t *testing.T) {
		f := mathutil.UnaryFn[float64](square).Add(mathutil.UnaryFn[float64](double))
		if got := f(3); math.Abs(got-15) > 1e-9 {
			t.Errorf("square.Add(double)(3) = %v, want 15", got)
		}
	})

	t.Run("Sub", func(t *testing.T) {
		f := mathutil.UnaryFn[float64](square).Sub(mathutil.UnaryFn[float64](double))
		if got := f(3); math.Abs(got-3) > 1e-9 {
			t.Errorf("square.Sub(double)(3) = %v, want 3", got)
		}
	})

	t.Run("Mul", func(t *testing.T) {
		f := mathutil.UnaryFn[float64](square).Mul(mathutil.UnaryFn[float64](double))
		if got := f(3); math.Abs(got-54) > 1e-9 {
			t.Errorf("square.Mul(double)(3) = %v, want 54", got)
		}
	})

	t.Run("Div", func(t *testing.T) {
		f := mathutil.UnaryFn[float64](square).Div(mathutil.UnaryFn[float64](double))
		if got := f(3); math.Abs(got-1.5) > 1e-9 {
			t.Errorf("square.Div(double)(3) = %v, want 1.5", got)
		}
	})
}

func TestConstant(t *testing.T) {
	f := mathutil.Constant[float64](5)
	if got := f(10); got != 5 {
		t.Errorf("Constant(5)(10) = %v, want 5", got)
	}
}

func TestKSigmoid(t *testing.T) {
	f := mathutil.KSigmoid[float64](2)
	if got := f(0); math.Abs(got-0.5) > 1e-9 {
		t.Errorf("KSigmoid(2)(0) = %v, want 0.5", got)
	}
}

func TestKSigmoidPrime(t *testing.T) {
	f := mathutil.KSigmoidPrime[float64](2)
	if got := f(0); math.Abs(got-0.5) > 1e-9 {
		t.Errorf("KSigmoidPrime(2)(0) = %v, want 0.5", got)
	}
}

func TestScale(t *testing.T) {
	f := mathutil.Scale[float64](2)
	if got := f(3); got != 6 {
		t.Errorf("Scale(2)(3) = %v, want 6", got)
	}
}

func TestOffset(t *testing.T) {
	f := mathutil.Offset[float64](2)
	if got := f(3); got != 5 {
		t.Errorf("Offset(2)(3) = %v, want 5", got)
	}
}

func TestAffine(t *testing.T) {
	f := mathutil.Affine[float64](2, 1)
	if got := f(3); got != 7 {
		t.Errorf("Affine(2, 1)(3) = %v, want 7", got)
	}
}

func TestPower(t *testing.T) {
	f := mathutil.Power[float64](2)
	if got := f(3); got != 9 {
		t.Errorf("Power(2)(3) = %v, want 9", got)
	}
}

func TestZero(t *testing.T) {
	if got := mathutil.Zero[float64](5); got != 0 {
		t.Errorf("Zero(5) = %v, want 0", got)
	}
}

func TestOne(t *testing.T) {
	if got := mathutil.One[float64](5); got != 1 {
		t.Errorf("One(5) = %v, want 1", got)
	}
}

func TestIdentity(t *testing.T) {
	if got := mathutil.Identity[float64](5); got != 5 {
		t.Errorf("Identity(5) = %v, want 5", got)
	}
}

func TestSquare(t *testing.T) {
	if got := mathutil.Square[float64](3); got != 9 {
		t.Errorf("Square(3) = %v, want 9", got)
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"zero", 0, 1},
		{"non-zero", 5, 0},
		{"negative", -5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.IsZero(tt.x); got != tt.want {
				t.Errorf("IsZero(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestSign(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"positive", 5, 1},
		{"negative", -5, -1},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.Sign(tt.x); got != tt.want {
				t.Errorf("Sign(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestSigmoid(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"zero", 0, 0.5},
		{"positive", 2, 0.8807970779778823},
		{"negative", -2, 0.11920292202211755},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.Sigmoid(tt.x); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Sigmoid(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestSigmoidPrime(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"zero", 0, 0.25},
		{"positive", 2, 0.1049935854035065},
		{"negative", -2, 0.1049935854035065},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.SigmoidPrime(tt.x); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("SigmoidPrime(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestBinaryFn(t *testing.T) {
	add := func(x, y float64) float64 { return x + y }
	mul := func(x, y float64) float64 { return x * y }

	t.Run("Add", func(t *testing.T) {
		f := mathutil.BinaryFn[float64](add).Add(mathutil.BinaryFn[float64](mul))
		if got := f(3, 4); math.Abs(got-19) > 1e-9 {
			t.Errorf("add.Add(mul)(3, 4) = %v, want 19", got)
		}
	})

	t.Run("Sub", func(t *testing.T) {
		f := mathutil.BinaryFn[float64](add).Sub(mathutil.BinaryFn[float64](mul))
		if got := f(3, 4); math.Abs(got-(-5)) > 1e-9 {
			t.Errorf("add.Sub(mul)(3, 4) = %v, want -5", got)
		}
	})

	t.Run("Mul", func(t *testing.T) {
		f := mathutil.BinaryFn[float64](add).Mul(mathutil.BinaryFn[float64](mul))
		if got := f(3, 4); math.Abs(got-84) > 1e-9 {
			t.Errorf("add.Mul(mul)(3, 4) = %v, want 84", got)
		}
	})

	t.Run("Div", func(t *testing.T) {
		f := mathutil.BinaryFn[float64](add).Div(mathutil.BinaryFn[float64](mul))
		if got := f(3, 4); math.Abs(got-0.5833333333333334) > 1e-9 {
			t.Errorf("add.Div(mul)(3, 4) = %v, want 0.4375", got)
		}
	})
}

func TestAdd(t *testing.T) {
	if got := mathutil.Add(3.0, 4.0); got != 7.0 {
		t.Errorf("Add(3.0, 4.0) = %v, want 7.0", got)
	}
}

func TestSub(t *testing.T) {
	if got := mathutil.Sub(3.0, 4.0); got != -1.0 {
		t.Errorf("Sub(3.0, 4.0) = %v, want -1.0", got)
	}
}

func TestMul(t *testing.T) {
	if got := mathutil.Mul(3.0, 4.0); got != 12.0 {
		t.Errorf("Mul(3.0, 4.0) = %v, want 12.0", got)
	}
}

func TestDiv(t *testing.T) {
	if got := mathutil.Div(3.0, 4.0); math.Abs(got-0.75) > 1e-9 {
		t.Errorf("Div(3.0, 4.0) = %v, want 0.75", got)
	}
}

func TestPow(t *testing.T) {
	if got := mathutil.Pow(2.0, 3.0); got != 8.0 {
		t.Errorf("Pow(2.0, 3.0) = %v, want 8.0", got)
	}
}

func TestClampedLerp(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		y    float64
		t    float64
		min  float64
		max  float64
		want float64
	}{
		{"within range", 0.0, 10.0, 0.5, 0.0, 10.0, 5.0},
		{"below min", 0.0, 10.0, -0.1, 0.0, 10.0, 0.0},
		{"above max", 0.0, 10.0, 1.1, 0.0, 10.0, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mathutil.ClampedLerp(tt.x, tt.y, tt.t, tt.min, tt.max); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("ClampedLerp(%v, %v, %v, %v, %v) = %v, want %v", tt.x, tt.y, tt.t, tt.min, tt.max, got, tt.want)
			}
		})
	}
}
