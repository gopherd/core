package pointer

// TryDeref a pointer value, if the pointer is nil, return the zero value of the type.
func TryDeref[T any](v *T) T {
	var zero T
	if v == nil {
		return zero
	}
	return *v
}

// Deref a pointer value, if the pointer is nil, return the zero value of the type and false.
func Deref[T any](v *T, defaultValue T) T {
	if v == nil {
		return defaultValue
	}
	return *v
}

// Of returns a pointer to the value.
func Of[T any](v T) *T {
	return &v
}
