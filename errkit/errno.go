package errkit

import (
	"errors"
	"fmt"

	"github.com/gopherd/core/constraints"
)

// Built-in error codes
const (
	EUnknown = -1 // Unknown error
	EOK      = 0  // No error
)

// Error is an interface that wraps the Error and Errno method.
type Error interface {
	error
	Errno() int
}

// errno is a struct that implements the Error interface.
type errno struct {
	no  int
	err error
}

// Errno returns the code of errno.
func (err errno) Errno() int {
	return err.no
}

// Error returns the error message of errno.
func (err errno) Error() string {
	return err.err.Error()
}

// Unwrap returns the wrapped error.
func (err errno) Unwrap() error {
	return err.err
}

// New wraps the error with code
func New[T constraints.Integer](code T, err error) Error {
	if err == nil {
		return nil
	}
	return errno{
		no:  int(code),
		err: err,
	}
}

// NewWithContext wraps an error with additional context information
func NewWithContext[T constraints.Integer](code T, err error, context string) Error {
	if err == nil {
		return nil
	}
	return errno{
		no:  int(code),
		err: fmt.Errorf("%s: %w", context, err),
	}
}

// Errno finds the first error in err's chain that contains errno.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// If err is nil, Errno returns EOK.
// If err does not contain errno, Errno returns EUnknown.
func Errno(err error) int {
	if err == nil {
		return EOK
	}
	for {
		if e, ok := err.(interface{ Errno() int }); ok {
			return e.Errno()
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
	}
	return EUnknown
}
