// Package encoding provides interfaces and utilities for encoding and decoding data.
package encoding

// Encoder is a function type that encodes a value into bytes.
type Encoder func(any) ([]byte, error)

// Decoder is a function type that decodes bytes into a provided value.
type Decoder func([]byte, any) error

// Transform decodes data using the provided decoder, then encodes it using the provided encoder.
func Transform(data []byte, encoder Encoder, decoder Decoder) ([]byte, error) {
	var v = make(map[string]any)
	if err := decoder(data, &v); err != nil {
		return nil, err
	}
	return encoder(v)
}
