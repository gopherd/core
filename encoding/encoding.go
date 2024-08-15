// Package encoding provides interfaces, utilities, and common functions for
// encoding and decoding data in various formats.
//
// This package defines common types for encoders and decoders, as well as
// utility functions that can be used across different encoding schemes.
// It serves as a foundation for building more specific encoding/decoding
// functionalities while providing a consistent interface.
package encoding

import (
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
