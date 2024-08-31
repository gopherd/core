// Package flags provides custom flag types and utilities.
package flags

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/gopherd/core/term"
)

// Map is a map of string key-value pairs that implements the flag.Value interface.
type Map map[string]string

// Get returns the value of the key.
func (m Map) Get(k string) string {
	if m == nil {
		return ""
	}
	return m[k]
}

// Contains reports whether the key is in the map.
func (m Map) Contains(k string) bool {
	if m == nil {
		return false
	}
	_, ok := m[k]
	return ok
}

// Lookup returns the value of the key and reports whether the key is in the map.
func (m Map) Lookup(k string) (string, bool) {
	if m == nil {
		return "", false
	}
	v, ok := m[k]
	return v, ok
}

// Set implements the flag.Value interface.
func (m *Map) Set(s string) error {
	if *m == nil {
		*m = make(Map)
	}
	var k, v string
	index := strings.Index(s, "=")
	if index < 0 {
		k = s
	} else {
		k, v = s[:index], s[index+1:]
	}
	if k == "" {
		return fmt.Errorf("empty key")
	}
	if _, dup := (*m)[k]; dup {
		return fmt.Errorf("already set: %q", k)
	}
	(*m)[k] = v
	return nil
}

// String implements the flag.Value interface.
func (m Map) String() string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	var sb strings.Builder
	for _, k := range keys {
		if sb.Len() > 0 {
			sb.WriteByte(',')
		}
		v := m[k]
		if needsQuoting(k) {
			sb.WriteString(strconv.Quote(k))
		} else {
			sb.WriteString(k)
		}
		if v != "" {
			sb.WriteByte('=')
			if needsQuoting(v) {
				sb.WriteString(strconv.Quote(v))
			} else {
				sb.WriteString(v)
			}
		}
	}
	return sb.String()
}

// IsAllValuesSet reports whether all values are set.
func (m Map) IsAllValuesSet() bool {
	for _, v := range m {
		if v == "" {
			return false
		}
	}
	return true
}

// Slice is a slice of strings that implements the flag.Value interface.
type Slice []string

// Set implements the flag.Value interface.
func (s *Slice) Set(v string) error {
	if v == "" {
		return fmt.Errorf("empty value")
	}
	*s = append(*s, v)
	return nil
}

// String implements the flag.Value interface.
func (s Slice) String() string {
	var sb strings.Builder
	for i, v := range s {
		if i > 0 {
			sb.WriteByte(',')
		}
		if needsQuoting(v) {
			sb.WriteString(strconv.Quote(v))
		} else {
			sb.WriteString(v)
		}
	}
	return sb.String()
}

// MapSlice is a map of string key-slice pairs that implements the flag.Value interface.
// It is used to parse multiple values for the same key.
type MapSlice map[string]Slice

// Get returns the slice of the key.
func (m MapSlice) Get(k string) Slice {
	if m == nil {
		return nil
	}
	return m[k]
}

// Contains reports whether the key is in the map.
func (m MapSlice) Contains(k string) bool {
	if m == nil {
		return false
	}
	_, ok := m[k]
	return ok
}

// Lookup returns the slice of the key and reports whether the key is in the map.
func (m MapSlice) Lookup(k string) (Slice, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m[k]
	return v, ok
}

// Set implements the flag.Value interface.
func (m *MapSlice) Set(s string) error {
	if *m == nil {
		*m = make(MapSlice)
	}
	var k, v string
	index := strings.Index(s, "=")
	if index <= 0 || index == len(s)-1 {
		return fmt.Errorf("invalid format: %q, expect key=value", s)
	}
	k, v = s[:index], s[index+1:]
	(*m)[k] = append((*m)[k], v)
	return nil
}

func (m MapSlice) String() string {
	if len(m) == 0 {
		return ""
	}
	var sb strings.Builder
	for k, vs := range m {
		if sb.Len() > 0 {
			sb.WriteByte(',')
		}
		if needsQuoting(k) {
			sb.WriteString(strconv.Quote(k))
		} else {
			sb.WriteString(k)
		}
		sb.WriteByte('=')
		sb.WriteString(vs.String())
	}
	return sb.String()
}

type options struct {
	nameColor term.Color
	newline   bool
}

// Option is an option for flag types.
type Option func(*options)

// NameColor sets the color of command names.
func NameColor(c term.Color) Option {
	return func(opts *options) {
		opts.nameColor = c
	}
}

// Newline adds a newline after the usage text.
func Newline() Option {
	return func(opts *options) {
		opts.newline = true
	}
}

// UsageFunc is a function that formats usage text.
type UsageFunc func(usage string) string

// UseUsage returns a UsageFunc that formats usage text with colorized command names.
func UseUsage(w io.Writer, opts ...Option) UsageFunc {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return func(usage string) string {
		return formatUsage(w, usage, o)
	}
}

// formatUsage returns a usage string with colorized command names.
func formatUsage(w io.Writer, usage string, opt options) string {
	if opt.newline {
		usage += "\n"
	}
	if !term.IsTerminal(w) || !term.IsSupportsAnsi() || !term.IsSupports256Colors() {
		return usage
	}
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					return usage[:i+1] + opt.nameColor.Format(usage[i+1:j]) + usage[j:]
				}
			}
			break
		}
	}
	return usage
}
