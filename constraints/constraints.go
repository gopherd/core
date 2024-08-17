// Package constraints provides type constraints for generic programming in Go.
//
// It defines a set of interface types that can be used as type constraints
// in generic functions and types. These constraints cover various numeric
// types and their combinations, allowing for more precise and flexible
// generic programming with numbers in Go.
package constraints

// Signed is a constraint that permits any signed integer type.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
type Complex interface {
	~complex64 | ~complex128
}

// SignedReal is a constraint that permits any signed real number type.
type SignedReal interface {
	Signed | Float
}

// Real is a constraint that permits any real number type.
type Real interface {
	Integer | Float
}

// SignedNumber is a constraint that permits any signed numeric type.
type SignedNumber interface {
	SignedReal | Complex
}

// Number is a constraint that permits any numeric type.
type Number interface {
	Real | Complex
}

// Field is a constraint that permits any number field type.
type Field interface {
	Float | Complex
}

// Addable is a constraint that permits any number or string type.
type Addable interface {
	Number | string
}
