package constraints

// Signed is a constraint that permits any signed integer type.
// It includes all variations of signed integers: int, int8, int16, int32, int64.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
// It includes all variations of unsigned integers: uint, uint8, uint16, uint32, uint64, and uintptr.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type.
// This includes both signed and unsigned integers.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
// It includes both float32 and float64.
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
// It includes both complex64 and complex128.
type Complex interface {
	~complex64 | ~complex128
}

// Ordered is a constraint that permits any ordered type.
// This includes all types that support comparison operators (<, <=, >=, >),
// specifically integer, float, and string types.
type Ordered interface {
	Integer | Float | ~string
}

// SignedReal is a constraint that permits any signed integer and floating-point type.
// It includes all signed integers and both float32 and float64.
type SignedReal interface {
	Signed | Float
}

// Real is a constraint that permits any real number type.
// This includes all integer types (both signed and unsigned) and both float32 and float64.
type Real interface {
	Integer | Float
}

// SignedNumber is a constraint that permits any signed integer, floating-point, and complex numeric type.
// It includes all signed integers, both float32 and float64, and both complex64 and complex128.
type SignedNumber interface {
	SignedReal | Complex
}

// Number is a constraint that permits any number type.
// This includes all real and complex numeric types.
type Number interface {
	Real | Complex
}

// Field is a constraint that permits any number field type.
// This includes both floating-point and complex numeric types.
type Field interface {
	Float | Complex
}
