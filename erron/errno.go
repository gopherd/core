package erron

import (
	"encoding/json"
	"errors"

	"github.com/gopherd/core/constraints"
)

const (
	EUnknown = -1
	EOK      = 0

	// User-defined errno should be greater than zero
)

var errOK = errors.New("ok")

type errno struct {
	code int
	err  error
}

func (err errno) Errno() int {
	return err.code
}

func (err errno) Error() string {
	return err.err.Error()
}

func (err errno) Unwrap() error {
	return err.err
}

func (err errno) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Error       int    `json:"error"`
		Description string `json:"description,omitempty"`
	}{
		Error:       err.code,
		Description: err.err.Error(),
	})
}

func AsErrno(err error) error {
	if err == nil {
		return errno{
			code: EOK,
			err:  errOK,
		}
	}
	var origin = err
	for {
		if e, ok := err.(errno); ok {
			return e
		}
		if e, ok := err.(interface{ Errno() int }); ok {
			return errno{
				code: e.Errno(),
				err:  origin,
			}
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
	}
	return errno{
		code: EUnknown,
		err:  origin,
	}
}

// Errno wraps the error with code
func Errno[T constraints.Integer](code T, err error) error {
	if err == nil {
		return nil
	}
	return errno{
		code: int(code),
		err:  err,
	}
}

// Errnof returns an error that formats as the given text with code.
func Errnof[T constraints.Integer](code T, format string, args ...any) error {
	return errno{
		code: int(code),
		err:  New(format, args...),
	}
}

// GetErrno finds the first error in err's chain that contains errno.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
func GetErrno(err error) int {
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
