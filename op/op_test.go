package op_test

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/gopherd/core/op"
)

var (
	negzero = math.Copysign(0, -1)
)

func TestOr(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"zero_nonzero", 0, 1, 1},
		{"nonzero_zero", 1, 0, 1},
		{"nonzero_nonzero", 2, 3, 2},
		{"zero_zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := op.Or(tt.a, tt.b); got != tt.want {
				t.Errorf("Or(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}

	// Test with strings
	if got := op.Or("", "default"); got != "default" {
		t.Errorf("Or(\"\", \"default\") = %v, want \"default\"", got)
	}
	if got := op.Or("value", "default"); got != "value" {
		t.Errorf("Or(\"value\", \"default\") = %v, want \"value\"", got)
	}

	// Test with floats
	if got := op.Or(0.0, 1.0); got != 1.0 {
		t.Errorf("Or(0.0, 1.0) = %v, want 1.0", got)
	}
	if got := op.Or(negzero, 1.0); got != 1.0 {
		t.Errorf("Or(negzero, 1.0) = %v, want 1.0", got)
	}
}

func TestOrFunc(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    func() int
		want int
	}{
		{"zero_func", 0, func() int { return 1 }, 1},
		{"nonzero_func", 1, func() int { return 2 }, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := op.OrFunc(tt.a, tt.b); got != tt.want {
				t.Errorf("OrFunc(%v, func()) = %v, want %v", tt.a, got, tt.want)
			}
		})
	}
}

func TestSetOr(t *testing.T) {
	tests := []struct {
		name       string
		a, b, want int
	}{
		{"zero_nonzero", 0, 1, 1},
		{"nonzero_zero", 1, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.a
			if got := op.SetOr(&a, tt.b); got != tt.want || a != tt.want {
				t.Errorf("SetOr(&%v, %v) = %v, want %v, a = %v, want %v", tt.a, tt.b, got, tt.want, a, tt.want)
			}
		})
	}
}

func TestSetOrFunc(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    func() int
		want int
	}{
		{"zero_func", 0, func() int { return 1 }, 1},
		{"nonzero_func", 1, func() int { return 2 }, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.a
			if got := op.SetOrFunc(&a, tt.b); got != tt.want || a != tt.want {
				t.Errorf("SetOrFunc(&%v, func()) = %v, want %v, a = %v, want %v", tt.a, got, tt.want, a, tt.want)
			}
		})
	}
}

func TestIf(t *testing.T) {
	if got := op.If(true, 1, 2); got != 1 {
		t.Errorf("If(true, 1, 2) = %v, want 1", got)
	}
	if got := op.If(false, 1, 2); got != 2 {
		t.Errorf("If(false, 1, 2) = %v, want 2", got)
	}
	if got := op.If(true, "yes", "no"); got != "yes" {
		t.Errorf("If(true, \"yes\", \"no\") = %v, want \"yes\"", got)
	}
	if got := op.If(false, "yes", "no"); got != "no" {
		t.Errorf("If(false, \"yes\", \"no\") = %v, want \"no\"", got)
	}
}

func TestIfFunc(t *testing.T) {
	if got := op.IfFunc(true, 1, func() int { return 2 }); got != 1 {
		t.Errorf("IfFunc(true, 1, func()) = %v, want 1", got)
	}
	if got := op.IfFunc(false, 1, func() int { return 2 }); got != 2 {
		t.Errorf("IfFunc(false, 1, func()) = %v, want 2", got)
	}
}

func TestIfFunc2(t *testing.T) {
	if got := op.IfFunc2(true, func() int { return 1 }, func() int { return 2 }); got != 1 {
		t.Errorf("IfFunc2(true, func(), func()) = %v, want 1", got)
	}
	if got := op.IfFunc2(false, func() int { return 1 }, func() int { return 2 }); got != 2 {
		t.Errorf("IfFunc2(false, func(), func()) = %v, want 2", got)
	}
}

func TestBin(t *testing.T) {
	tests := []struct {
		name string
		x    interface{}
		want int
	}{
		{"zero_int", 0, 0},
		{"nonzero_int", 1, 1},
		{"zero_string", "", 0},
		{"nonzero_string", "hello", 1},
		{"zero_float", 0.0, 0},
		{"nonzero_float", 1.1, 1},
		{"negative_zero", negzero, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got int
			switch x := tt.x.(type) {
			case int:
				got = op.Bin(x)
			case string:
				got = op.Bin(x)
			case float64:
				got = op.Bin(x)
			}
			if got != tt.want {
				t.Errorf("Bin(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	if got := op.First(1, 2, 3); got != 1 {
		t.Errorf("First(1, 2, 3) = %v, want 1", got)
	}
	if got := op.First("a", "b", "c"); got != "a" {
		t.Errorf("First(\"a\", \"b\", \"c\") = %v, want \"a\"", got)
	}
}

func TestSecond(t *testing.T) {
	if got := op.Second(1, "two", 3.0); got != "two" {
		t.Errorf("Second(1, \"two\", 3.0) = %v, want \"two\"", got)
	}
}

func TestThird(t *testing.T) {
	if got := op.Third(1, "two", 3.0, 4); got != 3.0 {
		t.Errorf("Third(1, \"two\", 3.0, 4) = %v, want 3.0", got)
	}
}

func TestDeref(t *testing.T) {
	x := 42
	var nilPtr *int

	if got := op.Deref(&x); got != 42 {
		t.Errorf("Deref(&x) = %v, want 42", got)
	}
	if got := op.Deref(nilPtr); got != 0 {
		t.Errorf("Deref(nil) = %v, want 0", got)
	}
}

func TestDerefOr(t *testing.T) {
	x := 42
	var nilPtr *int

	if got := op.DerefOr(&x, 0); got != 42 {
		t.Errorf("DerefOr(&x, 0) = %v, want 42", got)
	}
	if got := op.DerefOr(nilPtr, 10); got != 10 {
		t.Errorf("DerefOr(nil, 10) = %v, want 10", got)
	}
}

func TestDerefOrFunc(t *testing.T) {
	x := 42
	var nilPtr *int

	if got := op.DerefOrFunc(&x, func() int { return 0 }); got != 42 {
		t.Errorf("DerefOrFunc(&x, func()) = %v, want 42", got)
	}
	if got := op.DerefOrFunc(nilPtr, func() int { return 10 }); got != 10 {
		t.Errorf("DerefOrFunc(nil, func()) = %v, want 10", got)
	}
}

func TestAddr(t *testing.T) {
	x := 42
	if got := op.Addr(x); *got != x {
		t.Errorf("Addr(%v) = %v, want address of %v", x, got, x)
	}
}

func TestMust(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Must did not panic")
		}
	}()
	op.Must(fmt.Errorf("error"))
}

func TestMustValue(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Must did not panic")
		}
	}()
	op.MustValue(42, fmt.Errorf("error"))
}

func TestMustValue_no_panic(t *testing.T) {
	if got := op.MustValue(42, nil); got != 42 {
		t.Errorf("MustValue(42, nil) = %v, want 42", got)
	}
}

func TestReverseCompare(t *testing.T) {
	cmp := func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	}
	revCmp := op.ReverseCompare(cmp)

	if got := revCmp(1, 2); got != 1 {
		t.Errorf("ReverseCompare(cmp)(1, 2) = %v, want 1", got)
	}
	if got := revCmp(2, 1); got != -1 {
		t.Errorf("ReverseCompare(cmp)(2, 1) = %v, want -1", got)
	}
	if got := revCmp(1, 1); got != 0 {
		t.Errorf("ReverseCompare(cmp)(1, 1) = %v, want 0", got)
	}
}

func TestZero(t *testing.T) {
	if zero := op.Zero[bool](); zero != false {
		t.Errorf("Zero[bool]() = %v, want false", zero)
	}
	if zero := op.Zero[int](); zero != 0 {
		t.Errorf("Zero[int]() = %v, want 0", zero)
	}
	if zero := op.Zero[string](); zero != "" {
		t.Errorf("Zero[string]() = %v, want \"\"", zero)
	}
	if zero := op.Zero[float64](); zero != 0.0 {
		t.Errorf("Zero[float64]() = %v, want 0.0", zero)
	}
	if zero := op.Zero[error](); zero != nil {
		t.Errorf("Zero[error]() = %v, want nil", zero)
	}
}

func TestIdentity(t *testing.T) {
	x := 42
	if got := op.Identity(x)(); got != x {
		t.Errorf("Identity(%v)() = %v, want %v", x, got, x)
	}
	e := fmt.Errorf("error")
	if got := op.Identity(e)(); got != e {
		t.Errorf("Identity(%v)() = %v, want %v", e, got, e)
	}
}

func TestResult(t *testing.T) {
	var err = fmt.Errorf("error")
	tests := []struct {
		name string
		a    any
		b    error
		want any
	}{
		{"nil", 1, nil, 1},
		{"error", 1, err, err},
		{"nil_nil", nil, nil, nil},
		{"error_nil", nil, err, err},
		{"nil_error", err, nil, err},
		{"nil_string", "string", nil, "string"},
		{"error_string", "string", err, err},
		{"nil_float", 1.0, nil, 1.0},
		{"error_float", 1.0, err, err},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := op.Result(tt.a, tt.b); got != tt.want {
				t.Errorf("Result(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func ExampleOr() {
	fmt.Println(op.Or(0, 1))
	fmt.Println(op.Or(2, 3))
	fmt.Println(op.Or("", "default"))
	fmt.Println(op.Or("value", "default"))
	// Output:
	// 1
	// 2
	// default
	// value
}

func ExampleIf() {
	condition := true
	fmt.Println(op.If(condition, "It's true", "It's false"))
	condition = false
	fmt.Println(op.If(condition, "It's true", "It's false"))
	// Output:
	// It's true
	// It's false
}

func ExampleOr_sort() {
	type Order struct {
		Product  string
		Customer string
		Price    float64
	}
	orders := []Order{
		{"foo", "alice", 1.00},
		{"bar", "bob", 3.00},
		{"baz", "carol", 4.00},
		{"foo", "alice", 2.00},
		{"bar", "carol", 1.00},
		{"foo", "bob", 4.00},
	}
	// Sort by customer first, product second, and last by higher price
	slices.SortFunc(orders, func(a, b Order) int {
		customerCmp := op.If(a.Customer < b.Customer, -1, op.If(a.Customer > b.Customer, 1, 0))
		productCmp := op.If(a.Product < b.Product, -1, op.If(a.Product > b.Product, 1, 0))
		priceCmp := op.If(b.Price < a.Price, -1, op.If(b.Price > a.Price, 1, 0))

		return op.Or(op.Or(customerCmp, productCmp), priceCmp)
	})
	for _, order := range orders {
		fmt.Printf("%s %s %.2f\n", order.Customer, order.Product, order.Price)
	}
	// Output:
	// alice foo 2.00
	// alice foo 1.00
	// bob bar 3.00
	// bob foo 4.00
	// carol bar 1.00
	// carol baz 4.00
}
