// Package encoding provides interfaces, utilities, and common functions for
// encoding and decoding data in various formats.
//
// This package defines common types for encoders and decoders, as well as
// utility functions that can be used across different encoding schemes.
// It serves as a foundation for building more specific encoding/decoding
// functionalities while providing a consistent interface.
package encoding

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Encoder is a function type that encodes a value into bytes.
type Encoder func(any) ([]byte, error)

// Decoder is a function type that decodes bytes into a provided value.
type Decoder func([]byte, any) error

// Transform decodes the input data using the provided decoder,
// then re-encodes it using the provided encoder.
// It returns the encoded bytes and any error encountered during the process.
//
// The decoder should populate the provided value with the decoded data.
// The encoder should take the decoded value and produce the encoded bytes.
//
// If an error occurs during decoding or encoding, Transform returns nil for the bytes
// and the error describing the failure.
func Transform(data []byte, decoder Decoder, encoder Encoder) ([]byte, error) {
	var v any
	if err := decoder(data, &v); err != nil {
		return nil, fmt.Errorf("decoding error: %w", err)
	}

	encodedValue, err := encoder(v)
	if err != nil {
		return nil, fmt.Errorf("encoding error: %w", err)
	}

	return encodedValue, nil
}

// GetPosition returns the line and column number of the given offset in the data.
func GetPosition(data []byte, offset int) (line, column int) {
	line = 1
	column = 1

	for i := 0; i < offset && i < len(data); i++ {
		if data[i] == '\n' {
			line++
			column = 1
		} else if data[i] == '\r' {
			if i+1 < len(data) && data[i+1] == '\n' {
				continue
			}
			line++
			column = 1
		} else {
			column++
		}
	}

	return line, column
}

type SourceError struct {
	Filename string
	Line     int
	Column   int
	Offset   int
	Context  string
	Err      error
}

func (e *SourceError) Error() string {
	return fmt.Sprintf("%s:%d:%d(%s): %v", e.Filename, e.Line, e.Column, e.Context, e.Err)
}

func (e *SourceError) Unwrap() error {
	return e.Err
}

func GetJSONSourceError(filename string, data []byte, err error) error {
	if err == nil {
		return err
	}

	var offset int
	switch e := err.(type) {
	case *json.SyntaxError:
		offset = int(e.Offset)
	case *json.UnmarshalTypeError:
		offset = int(e.Offset)
	default:
		return err
	}
	if offset <= 0 {
		return err
	}

	const maxContext = 64
	line, column := GetPosition(data, offset)
	begin := bytes.LastIndexByte(data[:offset], '\n') + 1
	context := string(data[begin:offset])
	if offset-begin > maxContext {
		begin = offset - maxContext
		context = "..." + string(data[begin:offset])
	}
	return &SourceError{
		Filename: filename,
		Line:     line,
		Column:   column,
		Offset:   offset,
		Context:  context,
		Err:      err,
	}
}
