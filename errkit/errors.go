package errkit

import (
	"errors"
	"fmt"
	"strings"
)

// exitError represents an error that causes the service to exit.
type exitError struct {
	code    int
	message string
}

func (e *exitError) Error() string {
	return fmt.Sprintf("(exit %d) %s", e.code, e.message)
}

// NewExitError creates a new exit error with the given exit code.
func NewExitError(code int, messages ...string) error {
	return &exitError{code: code, message: strings.Join(messages, " ")}
}

// ExitCode returns the exit code of the error if it is an exit error.
func ExitCode(err error) (int, bool) {
	var exit *exitError
	if !errors.As(err, &exit) {
		return 0, false
	}
	return exit.code, true
}
