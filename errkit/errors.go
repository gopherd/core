package errkit

import (
	"bytes"
	"errors"
)

// List is a container for holding multiple errors.
//
// Usage:
//
//	var errList errkit.List
//	errList.Append(err1)
//	errList.Append(err2)
//	errList.Append(err3)
//	if err := errList.Err(); err != nil {
//		// handle error
//	}
type List struct {
	errors multiError
}

// Len returns the number of errors in the list.
func (l List) Len() int {
	return len(l.errors)
}

// Append adds err to the list if it is non-nil.
func (l *List) Append(err error) {
	if err != nil {
		l.errors = append(l.errors, err)
	}
}

// First returns the first error in the list, or nil if the list is empty.
func (l *List) First() error {
	if len(l.errors) == 0 {
		return nil
	}
	return l.errors[0]
}

// Last returns the last error in the list, or nil if the list is empty.
func (l *List) Last() error {
	if len(l.errors) == 0 {
		return nil
	}
	return l.errors[len(l.errors)-1]
}

// Err returns all errors as a single error, or nil if there are no errors.
func (l *List) Err() error {
	if len(l.errors) == 0 {
		return nil
	}
	return l.errors
}

type multiError []error

// Error returns a string representation of all errors in the list.
func (me multiError) Error() string {
	switch len(me) {
	case 0:
		return ""
	case 1:
		return me[0].Error()
	}
	var buf bytes.Buffer
	for i, e := range me {
		if i > 0 {
			buf.WriteString(" | ")
		}
		buf.WriteString(e.Error())
	}
	return buf.String()
}

// Is reports whether any error in the list matches the target error.
func (me multiError) Is(target error) bool {
	for _, err := range me {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// As finds the first error in the list that matches the target type.
func (me multiError) As(target interface{}) bool {
	for _, err := range me {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}
