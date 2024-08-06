package operator_test

import (
	"math"
	"testing"

	"github.com/gopherd/core/operator"
)

func TestOr(t *testing.T) {
	tests := []struct {
		name       string
		a, b, want int
	}{
		{"both non-zero", 1, 2, 1},
		{"a zero", 0, 2, 2},
		{"b zero", 1, 0, 1},
		{"both zero", 0, 0, 0},
		{"negative values", -1, -2, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := operator.Or(tt.a, tt.b); got != tt.want {
				t.Errorf("Or(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}

	// Test with string type
	if got := operator.Or("", "default"); got != "default" {
		t.Errorf(`Or("", "default") = %q, want "default"`, got)
	}
	if got := operator.Or("value", "default"); got != "value" {
		t.Errorf(`Or("value", "default") = %q, want "value"`, got)
	}
}

func TestOrFunc(t *testing.T) {
	counter := 0
	newFunc := func() int {
		counter++
		return 42
	}

	tests := []struct {
		name      string
		a         int
		wantCalls int
		want      int
	}{
		{"non-zero", 1, 0, 1},
		{"zero", 0, 1, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter = 0
			if got := operator.OrFunc(tt.a, newFunc); got != tt.want {
				t.Errorf("OrFunc(%v, newFunc) = %v, want %v", tt.a, got, tt.want)
			}
			if counter != tt.wantCalls {
				t.Errorf("newFunc called %d times, want %d", counter, tt.wantCalls)
			}
		})
	}
}

func TestTernary(t *testing.T) {
	tests := []struct {
		name       string
		condition  bool
		a, b, want int
	}{
		{"true condition", true, 1, 2, 1},
		{"false condition", false, 1, 2, 2},
		{"true with zero values", true, 0, 0, 0},
		{"false with zero values", false, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := operator.Ternary(tt.condition, tt.a, tt.b); got != tt.want {
				t.Errorf("Ternary(%v, %v, %v) = %v, want %v", tt.condition, tt.a, tt.b, got, tt.want)
			}
		})
	}

	// Test with string type
	if got := operator.Ternary(true, "yes", "no"); got != "yes" {
		t.Errorf(`Ternary(true, "yes", "no") = %q, want "yes"`, got)
	}
}

func TestTernaryFunc(t *testing.T) {
	aCounter, bCounter := 0, 0
	aFunc := func() int {
		aCounter++
		return 1
	}
	bFunc := func() int {
		bCounter++
		return 2
	}

	tests := []struct {
		name       string
		condition  bool
		wantResult int
		wantACalls int
		wantBCalls int
	}{
		{"true condition", true, 1, 1, 0},
		{"false condition", false, 2, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aCounter, bCounter = 0, 0
			if got := operator.TernaryFunc(tt.condition, aFunc, bFunc); got != tt.wantResult {
				t.Errorf("TernaryFunc(%v, aFunc, bFunc) = %v, want %v", tt.condition, got, tt.wantResult)
			}
			if aCounter != tt.wantACalls {
				t.Errorf("aFunc called %d times, want %d", aCounter, tt.wantACalls)
			}
			if bCounter != tt.wantBCalls {
				t.Errorf("bFunc called %d times, want %d", bCounter, tt.wantBCalls)
			}
		})
	}
}

func TestBool(t *testing.T) {
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
			if got := operator.Bool[int](tt.ok); got != tt.want {
				t.Errorf("Bool(%v) = %v, want %v", tt.ok, got, tt.want)
			}
		})
	}

	// Test with float64 type
	if got := operator.Bool[float64](true); got != 1.0 {
		t.Errorf("Bool[float64](true) = %v, want 1.0", got)
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name string
		x, y interface{}
		want bool
	}{
		{"equal ints", 1, 1, true},
		{"unequal ints", 1, 2, false},
		{"equal strings", "hello", "hello", true},
		{"unequal strings", "hello", "world", false},
		{"equal floats", 1.0, 1.0, true},
		{"unequal floats", 1.0, 1.1, false},
		{"float and int", 1.0, 1, false}, // Note: This will not be equal due to type mismatch
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := operator.Equal(tt.x, tt.y); got != tt.want {
				t.Errorf("Equal(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestLess(t *testing.T) {
	tests := []struct {
		name string
		x, y float64
		want bool
	}{
		{"less than", 1.0, 2.0, true},
		{"greater than", 2.0, 1.0, false},
		{"equal", 1.0, 1.0, false},
		{"x is NaN", math.NaN(), 1.0, true},
		{"y is NaN", 1.0, math.NaN(), false},
		{"both NaN", math.NaN(), math.NaN(), false},
		{"negative zero and zero", -0.0, 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := operator.Less(tt.x, tt.y); got != tt.want {
				t.Errorf("Less(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestGreater(t *testing.T) {
	tests := []struct {
		name string
		x, y float64
		want bool
	}{
		{"greater than", 2.0, 1.0, true},
		{"less than", 1.0, 2.0, false},
		{"equal", 1.0, 1.0, false},
		{"x is NaN", math.NaN(), 1.0, false},
		{"y is NaN", 1.0, math.NaN(), true},
		{"both NaN", math.NaN(), math.NaN(), false},
		{"negative zero and zero", -0.0, 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := operator.Greater(tt.x, tt.y); got != tt.want {
				t.Errorf("Greater(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestAsc(t *testing.T) {
	tests := []struct {
		name string
		x, y float64
		want int
	}{
		{"less than", 1.0, 2.0, -1},
		{"greater than", 2.0, 1.0, 1},
		{"equal", 1.0, 1.0, 0},
		{"x is NaN", math.NaN(), 1.0, -1},
		{"y is NaN", 1.0, math.NaN(), 1},
		{"both NaN", math.NaN(), math.NaN(), 0},
		{"negative zero and zero", -0.0, 0.0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := operator.Asc(tt.x, tt.y); got != tt.want {
				t.Errorf("Asc(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestDec(t *testing.T) {
	tests := []struct {
		name string
		x, y float64
		want int
	}{
		{"less than", 1.0, 2.0, 1},
		{"greater than", 2.0, 1.0, -1},
		{"equal", 1.0, 1.0, 0},
		{"x is NaN", math.NaN(), 1.0, 1},
		{"y is NaN", 1.0, math.NaN(), -1},
		{"both NaN", math.NaN(), math.NaN(), 0},
		{"negative zero and zero", -0.0, 0.0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := operator.Dec(tt.x, tt.y); got != tt.want {
				t.Errorf("Dec(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	if got := operator.First(1, 2, 3); got != 1 {
		t.Errorf("First(1, 2, 3) = %v, want 1", got)
	}
	if got := operator.First("a", "b", "c"); got != "a" {
		t.Errorf(`First("a", "b", "c") = %q, want "a"`, got)
	}
}

func TestSecond(t *testing.T) {
	if got := operator.Second(1, 2, 3); got != 2 {
		t.Errorf("Second(1, 2, 3) = %v, want 2", got)
	}
	if got := operator.Second("a", "b", "c"); got != "b" {
		t.Errorf(`Second("a", "b", "c") = %q, want "b"`, got)
	}
}

func TestThird(t *testing.T) {
	if got := operator.Third(1, 2, 3, 4); got != 3 {
		t.Errorf("Third(1, 2, 3, 4) = %v, want 3", got)
	}
	if got := operator.Third("a", "b", "c", "d"); got != "c" {
		t.Errorf(`Third("a", "b", "c", "d") = %q, want "c"`, got)
	}
}
