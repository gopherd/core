package erron

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strconv"
)

// New returns an error that formats as the given text.
func New(text string, args ...any) error {
	if len(args) > 0 {
		return fmt.Errorf(text, args...)
	}
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// throwedError wraps the error with source code position information.
type throwedError struct {
	pos string
	err error
}

func (e *throwedError) Error() string {
	return e.pos + ": " + e.err.Error()
}

func (e *throwedError) Unwrap() error {
	return e.err
}

func getCaller(depth int) string {
	_, file, line, ok := runtime.Caller(depth + 1)
	if !ok {
		file = "???"
		line = 0
	}
	return file + ":" + strconv.Itoa(line)
}

// Throw wraps the error with source code position information.
func Throw(err error) error {
	if err == nil {
		return nil
	}
	unwrapped := err
	for {
		if _, ok := unwrapped.(*throwedError); ok {
			return err
		}
		if unwrapped = errors.Unwrap(unwrapped); unwrapped == nil {
			break
		}
	}
	return &throwedError{
		pos: getCaller(1),
		err: err,
	}
}

// Throwf returns an error that formats as the given text with source code position information.
func Throwf(format string, args ...any) error {
	return &throwedError{
		pos: getCaller(1),
		err: fmt.Errorf(format, args...),
	}
}

// Try executes the fn, and then recovers any panics as an error.
// If no panics, return the result of fn.
func Try(fn func() error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			if err, ok = e.(error); !ok {
				if s, ok := e.(string); ok {
					err = New(s)
				} else {
					err = fmt.Errorf("%v", e)
				}
			}
			err = &throwedError{
				pos: getCaller(2),
				err: err,
			}
		}
	}()
	return fn()
}

// Errors is a container for holds errors
type Errors struct {
	errors errorList
}

func (errs Errors) Len() int {
	return len(errs.errors)
}

// Append appends err to list if e is non-nil
func (errs *Errors) Append(err error) {
	if err != nil {
		errs.errors = append(errs.errors, err)
	}
}

// First gets the first error of list
func (errs *Errors) First() error {
	if len(errs.errors) == 0 {
		return nil
	}
	return errs.errors[0]
}

// Last gets the last error of list
func (errs *Errors) Last() error {
	if len(errs.errors) == 0 {
		return nil
	}
	return errs.errors[len(errs.errors)-1]
}

// All merges all errors as an error, nil returned if no errors.
func (errs *Errors) All() error {
	if len(errs.errors) == 0 {
		return nil
	}
	return errs.errors
}

type errorList []error

func (errs errorList) Error() string {
	if len(errs) == 1 {
		return errs[0].Error()
	}
	var buf bytes.Buffer
	for i, e := range errs {
		if i > 0 {
			buf.WriteByte('|')
		}
		buf.WriteString(e.Error())
	}
	return buf.String()
}

// Is reports whether all of the errors are err
func (errs errorList) Is(err error) bool {
	for _, x := range errs {
		if !errors.Is(x, err) {
			return false
		}
	}
	return true
}
