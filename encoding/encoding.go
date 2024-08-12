// Package encoding provides interfaces and utilities for encoding and decoding data.
package encoding

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ErrInvalidString is returned when the input string is invalid.
var ErrInvalidString = errors.New("invalid string")

// Encoder is a function type that encodes a value into bytes.
type Encoder func(any) ([]byte, error)

// Decoder is a function type that decodes bytes into a provided value.
type Decoder func([]byte, any) error

// Marshaler is an interface for types that can marshal themselves into bytes.
type Marshaler interface {
	Marshal() ([]byte, error)
}

// Unmarshaler is an interface for types that can unmarshal bytes into themselves.
type Unmarshaler interface {
	Unmarshal([]byte) error
}

// UnmarshalString decodes a string from byte slice data.
// It supports both quoted strings and literal strings.
//
// Parameters:
//   - data: The byte slice containing the string to unmarshal.
//   - literalChar: The character used for literal strings. Use 0 for no literal strings.
//   - allowNewline: Whether newlines are allowed in literal strings.
//
// Returns:
//   - The unmarshaled string and nil error if successful.
//   - An empty string and an error if unmarshaling fails.
func UnmarshalString(data []byte, literalChar byte, allowNewline bool) (string, error) {
	data = bytes.TrimSpace(data)
	if len(data) < 2 {
		return "", errors.New("string too short")
	}

	switch data[0] {
	case '"':
		return unquoteString(data)
	case literalChar:
		if literalChar == 0 || literalChar == '"' {
			break
		}
		return extractLiteralString(data, literalChar, allowNewline)
	}

	return "", ErrInvalidString
}

func unquoteString(data []byte) (string, error) {
	if data[len(data)-1] != '"' {
		return "", errors.New("mismatched quotes")
	}
	return strconv.Unquote(string(data))
}

func extractLiteralString(data []byte, literalChar byte, allowNewline bool) (string, error) {
	if data[len(data)-1] != literalChar {
		return "", errors.New("mismatched quotes")
	}
	str := string(data[1 : len(data)-1])
	if !allowNewline && strings.ContainsRune(str, '\n') {
		return "", errors.New("newlines not allowed in literal string")
	}
	if !utf8.ValidString(str) {
		return "", errors.New("invalid UTF-8 in string")
	}
	return str, nil
}
