// Package flags provides custom flag types and utilities.
package flags

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

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
	parts := strings.SplitN(s, "=", 2)
	if len(parts) == 1 {
		k = parts[0]
	} else if len(parts) == 2 {
		k, v = parts[0], parts[1]
	} else {
		return fmt.Errorf("invalid format: %q, expected key=value", s)
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
	values := strings.Split(v, ",")
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		*s = append(*s, v)
	}
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

func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			// Quote anything except a backslash that would need quoting in a
			// JSON string, as well as space and '='
			if b != '\\' && (b == ' ' || b == '=' || !safeSet[b]) {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return true
		}
		i += size
	}
	return false
}

var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      false,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

type options struct {
	nameColor term.Color
}

func (o options) isEmpty() bool {
	return o.nameColor == "" || o.nameColor == term.None
}

// Option is an option for flag types.
type Option func(*options)

// NameColor sets the color of command names.
func NameColor(c term.Color) Option {
	return func(opts *options) {
		opts.nameColor = c
	}
}

// UsageFunc is a function that formats usage text.
type UsageFunc func(usage string) string

// UseUsage returns a UsageFunc that formats usage text with colorized command names.
func UseUsage(w io.Writer, opts ...Option) UsageFunc {
	o := options{
		nameColor: term.Turquoise,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return func(usage string) string {
		return formatUsage(w, usage, o)
	}
}

// formatUsage returns a usage string with colorized command names.
func formatUsage(w io.Writer, usage string, opt options) string {
	if opt.isEmpty() {
		return usage
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
