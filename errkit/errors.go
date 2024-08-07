package errkit

import (
	"errors"
	"fmt"
)

// exitError represents an error that causes the service to exit.
type exitError struct {
	code int
}

func (e *exitError) Error() string {
	return fmt.Sprintf("exit with code %d", e.code)
}

// NewExitError creates a new exit error with the given exit code.
func NewExitError(code int) error {
	return &exitError{code: code}
}

// ExitCode returns the exit code of the error if it is an exit error.
func ExitCode(err error) (int, bool) {
	var exit *exitError
	if !errors.As(err, &exit) {
		return 0, false
	}
	return exit.code, true
}
