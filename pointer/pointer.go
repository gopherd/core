package pointer

// Dereference a pointer value, if the pointer is nil, return the zero value of the type.
func Deref[T any](v *T) T {
	var zero T
	if v == nil {
		return zero
	}
	return *v
}

// TryDeref a pointer value, if the pointer is nil, return the zero value of the type and false.
func TryDeref[T any](v *T) (T, bool) {
	var zero T
	if v == nil {
		return zero, false
	}
	return *v, true
}

// Of returns a pointer to the value.
func Of[T any](v T) *T {
	return &v
}
