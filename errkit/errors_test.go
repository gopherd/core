package errkit

import (
	"errors"
	"fmt"
	"testing"
)

func TestList_Len(t *testing.T) {
	tests := []struct {
		name     string
		list     List
		expected int
	}{
		{"Empty list", List{}, 0},
		{"List with one error", List{errors: multiError{errors.New("error1")}}, 1},
		{"List with multiple errors", List{errors: multiError{errors.New("error1"), errors.New("error2"), errors.New("error3")}}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.Len(); got != tt.expected {
				t.Errorf("List.Len() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestList_Append(t *testing.T) {
	tests := []struct {
		name     string
		list     *List
		err      error
		expected int
	}{
		{"Append to empty list", &List{}, errors.New("error1"), 1},
		{"Append to non-empty list", &List{errors: multiError{errors.New("existing")}}, errors.New("error2"), 2},
		{"Append nil error", &List{}, nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.list.Append(tt.err)
			if got := tt.list.Len(); got != tt.expected {
				t.Errorf("After List.Append(), Len() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestList_First(t *testing.T) {
	tests := []struct {
		name     string
		list     *List
		expected error
	}{
		{"Empty list", &List{}, nil},
		{"List with one error", &List{errors: multiError{errors.New("error1")}}, errors.New("error1")},
		{"List with multiple errors", &List{errors: multiError{errors.New("error1"), errors.New("error2")}}, errors.New("error1")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.list.First()
			if (got == nil) != (tt.expected == nil) {
				t.Errorf("List.First() = %v, want %v", got, tt.expected)
			}
			if got != nil && got.Error() != tt.expected.Error() {
				t.Errorf("List.First() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestList_Last(t *testing.T) {
	tests := []struct {
		name     string
		list     *List
		expected error
	}{
		{"Empty list", &List{}, nil},
		{"List with one error", &List{errors: multiError{errors.New("error1")}}, errors.New("error1")},
		{"List with multiple errors", &List{errors: multiError{errors.New("error1"), errors.New("error2")}}, errors.New("error2")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.list.Last()
			if (got == nil) != (tt.expected == nil) {
				t.Errorf("List.Last() = %v, want %v", got, tt.expected)
			}
			if got != nil && got.Error() != tt.expected.Error() {
				t.Errorf("List.Last() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestList_Err(t *testing.T) {
	tests := []struct {
		name     string
		list     *List
		expected string
	}{
		{"Empty list", &List{}, ""},
		{"List with one error", &List{errors: multiError{errors.New("error1")}}, "error1"},
		{"List with multiple errors", &List{errors: multiError{errors.New("error1"), errors.New("error2")}}, "error1 | error2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.list.Err()
			if (got == nil) != (tt.expected == "") {
				t.Errorf("List.Err() = %v, want %v", got, tt.expected)
			}
			if got != nil && got.Error() != tt.expected {
				t.Errorf("List.Err() = %v, want %v", got.Error(), tt.expected)
			}
		})
	}
}

func TestMultiError_Error(t *testing.T) {
	tests := []struct {
		name     string
		me       multiError
		expected string
	}{
		{"Empty multiError", multiError{}, ""},
		{"Single error", multiError{errors.New("error1")}, "error1"},
		{"Multiple errors", multiError{errors.New("error1"), errors.New("error2"), errors.New("error3")}, "error1 | error2 | error3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.me.Error(); got != tt.expected {
				t.Errorf("multiError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMultiError_Is(t *testing.T) {
	err1 := errors.New("error1")
	err2 := errors.New("error2")
	err3 := fmt.Errorf("wrapped: %w", err2)

	tests := []struct {
		name     string
		me       multiError
		target   error
		expected bool
	}{
		{"Empty multiError", multiError{}, err1, false},
		{"Target in multiError", multiError{err1, err2}, err1, true},
		{"Target not in multiError", multiError{err1, err2}, errors.New("error3"), false},
		{"Wrapped error in multiError", multiError{err1, err3}, err2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.me.Is(tt.target); got != tt.expected {
				t.Errorf("multiError.Is() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// customError now correctly implements the error interface
type customError struct{ msg string }

func (e *customError) Error() string {
	return e.msg
}

func TestMultiError_As(t *testing.T) {
	err1 := errors.New("error1")
	err2 := &customError{msg: "custom error"}
	err3 := fmt.Errorf("wrapped: %w", err2)

	tests := []struct {
		name     string
		me       multiError
		expected bool
	}{
		{"Empty multiError", multiError{}, false},
		{"Target type in multiError", multiError{err1, err2}, true},
		{"Target type not in multiError", multiError{err1}, false},
		{"Wrapped target type in multiError", multiError{err1, err3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target *customError
			got := tt.me.As(&target)
			if got != tt.expected {
				t.Errorf("multiError.As() = %v, want %v", got, tt.expected)
			}
			if got {
				if target == nil {
					t.Errorf("multiError.As() succeeded but target is nil")
				} else if target.msg != "custom error" {
					t.Errorf("multiError.As() target has unexpected message: got %q, want %q", target.msg, "custom error")
				}
			}
		})
	}
}

// Test to ensure multiError.As works with error interface
func TestMultiError_AsErrorInterface(t *testing.T) {
	err1 := errors.New("error1")
	err2 := &customError{msg: "custom error"}
	me := multiError{err1, err2}

	var target error
	if !me.As(&target) {
		t.Errorf("multiError.As() = false, want true")
	}
	if target == nil {
		t.Errorf("multiError.As() succeeded but target is nil")
	}
}
