package errkit

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		err      error
		expected Error
	}{
		{"nil error", 1, nil, nil},
		{"non-nil error", 2, errors.New("test error"), errno{no: 2, err: errors.New("test error")}},
		{"zero code", 0, errors.New("zero code"), errno{no: 0, err: errors.New("zero code")}},
		{"negative code", -1, errors.New("negative code"), errno{no: -1, err: errors.New("negative code")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.code, tt.err)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("New(%d, %v) = %v, want nil", tt.code, tt.err, result)
				}
			} else {
				if result == nil {
					t.Errorf("New(%d, %v) = nil, want %v", tt.code, tt.err, tt.expected)
				} else if result.Errno() != tt.expected.Errno() || result.Error() != tt.expected.Error() {
					t.Errorf("New(%d, %v) = %v, want %v", tt.code, tt.err, result, tt.expected)
				}
			}
		})
	}
}

func TestNewWithContext(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		err      error
		context  string
		expected Error
	}{
		{"nil error", 1, nil, "context", nil},
		{"non-nil error", 2, errors.New("test error"), "context", errno{no: 2, err: fmt.Errorf("context: %w", errors.New("test error"))}},
		{"empty context", 3, errors.New("test error"), "", errno{no: 3, err: fmt.Errorf(": %w", errors.New("test error"))}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewWithContext(tt.code, tt.err, tt.context)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("NewWithContext(%d, %v, %q) = %v, want nil", tt.code, tt.err, tt.context, result)
				}
			} else {
				if result == nil {
					t.Errorf("NewWithContext(%d, %v, %q) = nil, want %v", tt.code, tt.err, tt.context, tt.expected)
				} else if result.Errno() != tt.expected.Errno() || result.Error() != tt.expected.Error() {
					t.Errorf("NewWithContext(%d, %v, %q) = %v, want %v", tt.code, tt.err, tt.context, result, tt.expected)
				}
			}
		})
	}
}

func TestErrno(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"nil error", nil, EOK},
		{"custom error", New(42, errors.New("custom error")), 42},
		{"wrapped custom error", fmt.Errorf("wrapped: %w", New(42, errors.New("custom error"))), 42},
		{"non-errno error", errors.New("regular error"), EUnknown},
		{"deeply wrapped custom error", fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", New(42, errors.New("custom error")))), 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Errno(tt.err)
			if result != tt.expected {
				t.Errorf("Errno(%v) = %d, want %d", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     int
		expected bool
	}{
		{"nil error", nil, EOK, true},
		{"matching code", New(42, errors.New("custom error")), 42, true},
		{"non-matching code", New(42, errors.New("custom error")), 43, false},
		{"wrapped matching code", fmt.Errorf("wrapped: %w", New(42, errors.New("custom error"))), 42, true},
		{"non-errno error", errors.New("regular error"), EUnknown, true},
		{"EOK with nil error", nil, EOK, true},
		{"EOK with non-nil error", errors.New("some error"), EOK, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Is(tt.err, tt.code)
			if result != tt.expected {
				t.Errorf("Is(%v, %d) = %v, want %v", tt.err, tt.code, result, tt.expected)
			}
		})
	}
}

func TestErrnoMethods(t *testing.T) {
	err := New(42, errors.New("test error"))

	t.Run("Errno", func(t *testing.T) {
		if err.Errno() != 42 {
			t.Errorf("err.Errno() = %d, want 42", err.Errno())
		}
	})

	t.Run("Error", func(t *testing.T) {
		if err.Error() != "test error" {
			t.Errorf("err.Error() = %q, want \"test error\"", err.Error())
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil || unwrapped.Error() != "test error" {
			t.Errorf("errors.Unwrap(err) = %v, want error with message \"test error\"", unwrapped)
		}
	})
}
