package encoding

import (
	"testing"
)

func TestUnmarshalString(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		literalChar  byte
		allowNewline bool
		want         string
		wantErr      bool
	}{
		{"Empty input", []byte{}, 0, false, "", true},
		{"Single character", []byte{'"'}, 0, false, "", true},
		{"Valid quoted string", []byte(`"hello"`), 0, false, "hello", false},
		{"Valid literal string", []byte(`'hello'`), '\'', false, "hello", false},
		{"Quoted string with escapes", []byte(`"he\"llo"`), 0, false, `he"llo`, false},
		{"Literal string with newline allowed", []byte("'hello\nworld'"), '\'', true, "hello\nworld", false},
		{"Literal string with newline disallowed", []byte("'hello\nworld'"), '\'', false, "", true},
		{"Invalid UTF-8 in literal string", []byte{'\'', 0xFF, '\''}, '\'', false, "", true},
		{"Mismatched quotes in quoted string", []byte(`"hello`), 0, false, "", true},
		{"Mismatched quotes in literal string", []byte(`'hello`), '\'', false, "", true},
		{"Invalid string start", []byte(`hello`), 0, false, "", true},
		{"Quoted string with literal char", []byte(`"'hello'"`), '\'', false, "'hello'", false},
		{"Literal string with quote char", []byte(`'"hello"'`), '\'', false, `"hello"`, false},
		{"Literal char is quote", []byte(`"hello"`), '"', false, "hello", false},
		{"Whitespace before valid string", []byte("  'hello'"), '\'', false, "hello", false},
		{"Whitespace after valid string", []byte("'hello'  "), '\'', false, "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalString(tt.input, tt.literalChar, tt.allowNewline)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UnmarshalString() = %v, want %v", got, tt.want)
			}
		})
	}
}
