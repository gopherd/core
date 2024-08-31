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
	None      = Color("-")
	Bold      = Color("\033[1m")
	Dim       = Color("\033[2m")
	Italic    = Color("\033[3m")
	Underline = Color("\033[4m")
	Blink     = Color("\033[5m")
	Reverse   = Color("\033[7m")
	Hidden    = Color("\033[8m")

	// Basic colors
	Black   = Color("\033[30m")
	Red     = Color("\033[31m")
	Green   = Color("\033[32m")
	Yellow  = Color("\033[33m")
	Blue    = Color("\033[34m")
	Magenta = Color("\033[35m")
	Cyan    = Color("\033[36m")
	White   = Color("\033[37m")

	// Bright colors
	BrightBlack   = Color("\033[90m")
	BrightRed     = Color("\033[91m")
	BrightGreen   = Color("\033[92m")
	BrightYellow  = Color("\033[93m")
	BrightBlue    = Color("\033[94m")
	BrightMagenta = Color("\033[95m")
	BrightCyan    = Color("\033[96m")
	BrightWhite   = Color("\033[97m")

	// Background colors
	BgBlack   = Color("\033[40m")
	BgRed     = Color("\033[41m")
	BgGreen   = Color("\033[42m")
	BgYellow  = Color("\033[43m")
	BgBlue    = Color("\033[44m")
	BgMagenta = Color("\033[45m")
	BgCyan    = Color("\033[46m")
	BgWhite   = Color("\033[47m")

	// Bright background colors
	BgBrightBlack   = Color("\033[100m")
	BgBrightRed     = Color("\033[101m")
	BgBrightGreen   = Color("\033[102m")
	BgBrightYellow  = Color("\033[103m")
	BgBrightBlue    = Color("\033[104m")
	BgBrightMagenta = Color("\033[105m")
	BgBrightCyan    = Color("\033[106m")
	BgBrightWhite   = Color("\033[107m")

	// 256-color mode
	Turquoise      = Color("\033[38;5;80m")
	Orange         = Color("\033[38;5;214m")
	Pink           = Color("\033[38;5;200m")
	Violet         = Color("\033[38;5;135m")
	LightGreen     = Color("\033[38;5;119m")
	LightBlue      = Color("\033[38;5;123m")
	DeepPink       = Color("\033[38;5;198m")
	LightSeaGreen  = Color("\033[38;5;37m")
	MediumPurple   = Color("\033[38;5;141m")
	DarkOrange     = Color("\033[38;5;208m")
	SteelBlue      = Color("\033[38;5;67m")
	IndianRed      = Color("\033[38;5;167m")
	Chartreuse     = Color("\033[38;5;118m")
	MediumOrchid   = Color("\033[38;5;134m")
	DodgerBlue     = Color("\033[38;5;33m")
	Crimson        = Color("\033[38;5;160m")
	MediumSeaGreen = Color("\033[38;5;48m")
	Gold           = Color("\033[38;5;220m")
)

// RGB color function
func RGB(r, g, b int) Color {
	return Color(fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b))
}

// Background RGB color function
func BgRGB(r, g, b int) Color {
	return Color(fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b))
}

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
