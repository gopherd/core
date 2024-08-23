// Package term provides terminal-related utilities.
package term

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

var (
	isSupportsAnsi = func() bool {
		return runtime.GOOS != "windows" || os.Getenv("TERM") != ""
	}()
	isSupports256Colors = func() bool {
		term := os.Getenv("TERM")
		if strings.Contains(term, "256color") {
			return true
		}
		if strings.Contains(term, "xterm") || strings.Contains(term, "screen") || strings.Contains(term, "tmux") || strings.Contains(term, "rxvt") {
			return true
		}
		colorterm := os.Getenv("COLORTERM")
		if strings.Contains(strings.ToLower(colorterm), "truecolor") || strings.Contains(strings.ToLower(colorterm), "24bit") {
			return true
		}
		return false
	}()
)

// IsTerminal reports whether w is a terminal.
func IsTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// IsSupportsAnsi reports whether the terminal supports ANSI escape codes.
func IsSupportsAnsi() bool {
	return isSupportsAnsi
}

// IsSupports256Colors reports whether the terminal supports 256 colors.
func IsSupports256Colors() bool {
	return isSupports256Colors
}

func ColorizeWriter(w io.Writer, c Color) io.Writer {
	return &colorizeWriter{w: w, color: c, isTerminal: IsTerminal(w)}
}

type colorizeWriter struct {
	w          io.Writer
	color      Color
	isTerminal bool
}

func (w *colorizeWriter) Write(p []byte) (n int, err error) {
	if w.isTerminal && isSupportsAnsi && (!w.color.Is256() || isSupports256Colors) {
		return w.w.Write([]byte(w.color.Format(string(p))))
	}
	return w.w.Write(p)
}

// Reset is the ANSI escape code to reset all attributes.
const Reset = "\033[0m"

const (
	None    = Color("-")
	Red     = Color("\033[31m")
	Green   = Color("\033[32m")
	Yellow  = Color("\033[33m")
	Blue    = Color("\033[34m")
	Magenta = Color("\033[35m")
	Cyan    = Color("\033[36m")
	White   = Color("\033[37m")

	BrightRed     = Color("\033[91m")
	BrightGreen   = Color("\033[92m")
	BrightYellow  = Color("\033[93m")
	BrightBlue    = Color("\033[94m")
	BrightMagenta = Color("\033[95m")
	BrightCyan    = Color("\033[96m")
	BrightWhite   = Color("\033[97m")

	// 256 colors
	Turquoise  = Color("\033[38;5;80m")
	Orange     = Color("\033[38;5;214m")
	Pink       = Color("\033[38;5;200m")
	Violet     = Color("\033[38;5;135m")
	LightGreen = Color("\033[38;5;119m")
	LightBlue  = Color("\033[38;5;123m")
)

// Color represents a terminal color.
type Color string

// Is256 reports whether the color is a 256 color.
func (c Color) Is256() bool {
	return strings.HasPrefix(string(c), "\033[38;5;")
}

// Colorize returns a colorized string.
func (c Color) Colorize(s string) colorizedString {
	return colorizedString{value: s, color: c}
}

// Format formats the string s with the color c.
func (c Color) Format(s string) string {
	if c == "" || c == None {
		return s
	}
	return string(c) + s + Reset
}

// Fprint formats using the default formats for its operands and writes to w.
func Fprint(w io.Writer, a ...any) (n int, err error) {
	isTerminal := IsTerminal(w)
	if isTerminal && isSupports256Colors {
		return fmt.Fprint(w, a...)
	}
	return fmt.Fprint(w, removeColors(isTerminal, a)...)
}

// Fprintf formats according to a format specifier and writes to w.
func Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	isTerminal := IsTerminal(w)
	if isTerminal && isSupports256Colors {
		return fmt.Fprintf(w, format, a...)
	}
	return fmt.Fprintf(w, format, removeColors(isTerminal, a)...)
}

// Fprintln formats using the default formats for its operands and writes to w.
func Fprintln(w io.Writer, a ...any) (n int, err error) {
	isTerminal := IsTerminal(w)
	if isTerminal && isSupports256Colors {
		return fmt.Fprintln(w, a...)
	}
	return fmt.Fprintln(w, removeColors(isTerminal, a)...)
}

type colorizedString struct {
	value string
	color Color
}

func (s colorizedString) String() string {
	return s.color.Format(s.value)
}

func getColorizedString(a any) *colorizedString {
	if s, ok := a.(colorizedString); ok {
		return &s
	}
	if s, ok := a.(*colorizedString); ok && s != nil {
		return s
	}
	return nil
}

func removeColors(isTerminal bool, a []any) []any {
	for _, arg := range a {
		s := getColorizedString(arg)
		if s == nil || (isTerminal && isSupportsAnsi && !s.color.Is256()) {
			continue
		}
		args := make([]any, len(a))
		for i := range a {
			s := getColorizedString(arg)
			if s != nil && (!isTerminal || !isSupportsAnsi || s.color.Is256()) {
				args[i] = s.value
			} else {
				args[i] = a[i]
			}
		}
		return args
	}
	return a
}
