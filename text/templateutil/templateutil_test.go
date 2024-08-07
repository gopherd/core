package templateutil

import (
	"fmt"
	"html/template"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestDefaultTemplate(t *testing.T) {
	// Verify that the template includes default functions
	for name := range DefaultFuncs() {
		tmpl := DefaultTemplate("test")
		if tmpl == nil {
			t.Fatal("DefaultTemplate returned nil")
		}
		if tmpl.Name() != "test" {
			t.Errorf("Expected template name 'test', got '%s'", tmpl.Name())
		}
		testTemplate := fmt.Sprintf(`{{if %s}}{{end}}`, name)
		_, err := tmpl.Parse(testTemplate)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("function %s not defined", name)) {
				t.Errorf("DefaultTemplate missing function %s", name)
			} else {
				t.Errorf("Error parsing template with function %s: %v", name, err)
			}
		}
	}
}

func TestContainerFuncs(t *testing.T) {
	tests := []struct {
		name     string
		fn       string
		args     []any
		expected any
	}{
		{"len", "len", []any{[]any{1, 2, 3}}, 3},
		{"list", "list", []any{1, 2, 3}, []any{1, 2, 3}},
		{"bools", "bools", []any{true, false, true}, []bool{true, false, true}},
		{"strings", "strings", []any{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"ints", "ints", []any{1, 2, 3}, []int{1, 2, 3}},
		{"int8s", "int8s", []any{int8(1), int8(2), int8(3)}, []int8{1, 2, 3}},
		{"int16s", "int16s", []any{int16(1), int16(2), int16(3)}, []int16{1, 2, 3}},
		{"int32s", "int32s", []any{int32(1), int32(2), int32(3)}, []int32{1, 2, 3}},
		{"int64s", "int64s", []any{int64(1), int64(2), int64(3)}, []int64{1, 2, 3}},
		{"uints", "uints", []any{uint(1), uint(2), uint(3)}, []uint{1, 2, 3}},
		{"uint8s", "uint8s", []any{uint8(1), uint8(2), uint8(3)}, []uint8{1, 2, 3}},
		{"uint16s", "uint16s", []any{uint16(1), uint16(2), uint16(3)}, []uint16{1, 2, 3}},
		{"uint32s", "uint32s", []any{uint32(1), uint32(2), uint32(3)}, []uint32{1, 2, 3}},
		{"uint64s", "uint64s", []any{uint64(1), uint64(2), uint64(3)}, []uint64{1, 2, 3}},
		{"float32s", "float32s", []any{float32(1.1), float32(2.2), float32(3.3)}, []float32{1.1, 2.2, 3.3}},
		{"float64s", "float64s", []any{1.1, 2.2, 3.3}, []float64{1.1, 2.2, 3.3}},
	}

	funcs := DefaultFuncs()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := reflect.ValueOf(funcs[tt.fn])
			result := fn.Call(convToValues(tt.args))
			if !reflect.DeepEqual(result[0].Interface(), tt.expected) {
				t.Errorf("%s() = %v, want %v", tt.name, result[0].Interface(), tt.expected)
			}
		})
	}
}

func TestStringsFuncs(t *testing.T) {
	tests := []struct {
		name     string
		fn       string
		args     []any
		expected any
	}{
		{"contains", "contains", []any{"hello", "ll"}, true},
		{"contains false", "contains", []any{"hello", "xx"}, false},
		{"count", "count", []any{"hello", "l"}, 2},
		{"hasPrefix", "hasPrefix", []any{"hello", "he"}, true},
		{"hasPrefix false", "hasPrefix", []any{"hello", "x"}, false},
		{"hasSuffix", "hasSuffix", []any{"hello", "lo"}, true},
		{"hasSuffix false", "hasSuffix", []any{"hello", "x"}, false},
		{"index", "index", []any{"hello", "l"}, 2},
		{"index not found", "index", []any{"hello", "x"}, -1},
		{"join", "join", []any{",", "a", "b", "c"}, "a,b,c"},
		{"lastIndex", "lastIndex", []any{"hello", "l"}, 3},
		{"lastIndex not found", "lastIndex", []any{"hello", "x"}, -1},
		{"repeat", "repeat", []any{"a", 3}, "aaa"},
		{"replace", "replace", []any{"hello", "l", "x", 1}, "hexlo"},
		{"replaceAll", "replaceAll", []any{"hello", "l", "x"}, "hexxo"},
		{"split", "split", []any{"a,b,c", ","}, []string{"a", "b", "c"}},
		{"toLower", "toLower", []any{"HeLLo"}, "hello"},
		{"toUpper", "toUpper", []any{"HeLLo"}, "HELLO"},
		{"toValidUTF8", "toValidUTF8", []any{"hello\xffworld", ""}, "helloworld"},
		{"trim", "trim", []any{" hello ", " "}, "hello"},
		{"trimLeft", "trimLeft", []any{" hello ", " "}, "hello "},
		{"trimRight", "trimRight", []any{" hello ", " "}, " hello"},
		{"trimPrefix", "trimPrefix", []any{"hello world", "hello "}, "world"},
		{"trimSuffix", "trimSuffix", []any{"hello world", " world"}, "hello"},
		{"trimSpace", "trimSpace", []any{" hello world "}, "hello world"},
	}

	funcs := DefaultFuncs()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := reflect.ValueOf(funcs[tt.fn])
			result := fn.Call(convToValues(tt.args))
			if !reflect.DeepEqual(result[0].Interface(), tt.expected) {
				t.Errorf("%s() = %v, want %v", tt.name, result[0].Interface(), tt.expected)
			}
		})
	}
}

func TestConvFuncs(t *testing.T) {
	tests := []struct {
		name     string
		fn       string
		arg      any
		expected any
		err      bool
	}{
		{"float32", "float32", 1, float32(1), false},
		{"float32 error", "float32", "not a number", float32(0), true},
		{"float64", "float64", 1, float64(1), false},
		{"float64 error", "float64", "not a number", float64(0), true},
		{"int", "int", 1.5, 1, false},
		{"int error", "int", "not a number", 0, true},
		{"int8", "int8", 127, int8(127), false},
		{"int16", "int16", 32767, int16(32767), false},
		{"int32", "int32", 2147483647, int32(2147483647), false},
		{"int64", "int64", 9223372036854775807, int64(9223372036854775807), false},
		{"uint", "uint", 1, uint(1), false},
		{"uint negative", "uint", -1, uint(0), true},
		{"uint8", "uint8", 255, uint8(255), false},
		{"uint16", "uint16", 65535, uint16(65535), false},
		{"uint32", "uint32", 4294967295, uint32(4294967295), false},
		{"bool", "bool", 1, true, false},
		{"bool false", "bool", 0, false, false},
		{"string", "string", 123, "123", false},
		{"bytes", "bytes", "hello", []byte("hello"), false},
		{"bytes error", "bytes", 123, []byte(nil), true},
		{"runes", "runes", "hello", []rune("hello"), false},
		{"runes error", "runes", 123, []rune(nil), true},
		{"rune", "rune", "a", 'a', false},
		{"rune error", "rune", "", rune(0), true},
		{"byte", "byte", "a", byte('a'), false},
		{"byte error", "byte", "", byte(0), true},
	}

	funcs := DefaultFuncs()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := reflect.ValueOf(funcs[tt.fn])
			result := fn.Call([]reflect.Value{reflect.ValueOf(tt.arg)})
			if tt.err {
				if result[1].IsNil() {
					t.Errorf("%s() expected error, got nil", tt.name)
				}
			} else {
				if !result[1].IsNil() {
					t.Errorf("%s() unexpected error: %v", tt.name, result[1].Interface())
				}
				if !reflect.DeepEqual(result[0].Interface(), tt.expected) {
					t.Errorf("%s() = %v, want %v", tt.name, result[0].Interface(), tt.expected)
				}
			}
		})
	}
}

func TestMathFuncs(t *testing.T) {
	tests := []struct {
		name     string
		fn       string
		args     []any
		expected any
		err      bool
	}{
		{"sum ints", "sum", []any{1, 2, 3}, int64(6), false},
		{"sum floats", "sum", []any{1.1, 2.2, 3.3}, 6.6, false},
		{"sum mixed", "sum", []any{1, 2.5, 3}, 6.5, false},
		{"sum strings", "sum", []any{"a", "b", "c"}, "abc", false},
		{"sum error", "sum", []any{1, "a", 2}, nil, true},
		{"add ints", "add", []any{1, 2}, int64(3), false},
		{"add floats", "add", []any{1.1, 2.2}, 3.3, false},
		{"add mixed", "add", []any{1, 2.5}, 3.5, false},
		{"add strings", "add", []any{"a", "b"}, "ab", false},
		{"add error", "add", []any{1, "a"}, nil, true},
		{"sub ints", "sub", []any{3, 2}, int64(1), false},
		{"sub floats", "sub", []any{3.3, 2.2}, 1.1, false},
		{"sub mixed", "sub", []any{3, 2.5}, 0.5, false},
		{"sub error", "sub", []any{3, "a"}, nil, true},
		{"mul ints", "mul", []any{2, 3}, int64(6), false},
		{"mul floats", "mul", []any{2.2, 3.3}, 7.26, false},
		{"mul mixed", "mul", []any{2, 3.5}, 7.0, false},
		{"mul error", "mul", []any{2, "a"}, nil, true},
		{"div ints", "div", []any{6, 3}, int64(2), false},
		{"div floats", "div", []any{6.6, 3.3}, 2.0, false},
		{"div mixed", "div", []any{7, 2.0}, 3.5, false},
		{"div by zero", "div", []any{6, 0}, nil, true},
		{"div error", "div", []any{6, "a"}, nil, true},
		{"mod ints", "mod", []any{7, 3}, int64(1), false},
		{"mod floats", "mod", []any{7.5, 3.2}, 1.1, false},
		{"mod mixed", "mod", []any{7, 3.0}, 1.0, false},
		{"mod by zero", "mod", []any{7, 0}, nil, true},
		{"mod error", "mod", []any{7, "a"}, nil, true},
		{"pow ints", "pow", []any{2, 3}, int64(8), false},
		{"pow floats", "pow", []any{2.0, 3.0}, 8.0, false},
		{"pow mixed", "pow", []any{2, 3.0}, 8.0, false},
		{"pow negative", "pow", []any{2, -2}, 0.25, false},
		{"pow error", "pow", []any{2, "a"}, nil, true},
	}

	funcs := DefaultFuncs()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := reflect.ValueOf(funcs[tt.fn])
			result := fn.Call(convToValues(tt.args))
			if tt.err {
				if result[1].IsNil() {
					t.Errorf("%s() expected error, got nil", tt.name)
				}
			} else {
				if !result[1].IsNil() {
					t.Errorf("%s() unexpected error: %v", tt.name, result[1].Interface())
				}
				if !reflect.DeepEqual(result[0].Interface(), tt.expected) {
					if f1, err := toFloat64(result[0].Interface()); err == nil {
						if f2, err := toFloat64(tt.expected); err == nil {
							if math.Abs(f1-f2) < 1e-6 {
								return
							}
						}
					}
					t.Errorf("%s() = %v, want %v", tt.name, result[0].Interface(), tt.expected)
				}
			}
		})
	}
}

func TestWithFuncs(t *testing.T) {
	baseFunc := template.FuncMap{
		"base": func() string { return "base" },
	}

	additionalFunc1 := template.FuncMap{
		"add1": func() string { return "add1" },
	}

	additionalFunc2 := template.FuncMap{
		"add2": func() string { return "add2" },
		"base": func() string { return "overwritten" },
	}

	result := withFuncs(baseFunc, additionalFunc1, additionalFunc2)

	expectedFuncs := []string{"base", "add1", "add2"}
	for _, funcName := range expectedFuncs {
		if _, ok := result[funcName]; !ok {
			t.Errorf("Expected function %s not found in result", funcName)
		}
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 functions, got %d", len(result))
	}

	if result["base"].(func() string)() != "overwritten" {
		t.Errorf("Expected 'base' function to be overwritten")
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
		err      bool
	}{
		{"int", 42, 42.0, false},
		{"int8", int8(8), 8.0, false},
		{"int16", int16(16), 16.0, false},
		{"int32", int32(32), 32.0, false},
		{"int64", int64(64), 64.0, false},
		{"uint", uint(42), 42.0, false},
		{"uint8", uint8(8), 8.0, false},
		{"uint16", uint16(16), 16.0, false},
		{"uint32", uint32(32), 32.0, false},
		{"uint64", uint64(64), 64.0, false},
		{"float32", float32(3.14), 3.14, false},
		{"float64", 3.14, 3.14, false},
		{"string", "not a number", 0, true},
		{"bool", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toFloat64(tt.input)
			if tt.err {
				if err == nil {
					t.Errorf("toFloat64(%v) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("toFloat64(%v) unexpected error: %v", tt.input, err)
				}
				if math.Abs(result-tt.expected) > 1e-6 {
					t.Errorf("toFloat64(%v) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int64
		err      bool
	}{
		{"int", 42, 42, false},
		{"int8", int8(8), 8, false},
		{"int16", int16(16), 16, false},
		{"int32", int32(32), 32, false},
		{"int64", int64(64), 64, false},
		{"uint", uint(42), 42, false},
		{"uint8", uint8(8), 8, false},
		{"uint16", uint16(16), 16, false},
		{"uint32", uint32(32), 32, false},
		{"uint64", uint64(64), 64, false},
		{"uint64 overflow", uint64(math.MaxUint64), 0, true},
		{"float32", float32(3.14), 0, true},
		{"float64", 3.14, 0, true},
		{"string", "not a number", 0, true},
		{"bool", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toInt64(tt.input, true)
			if tt.err {
				if err == nil {
					t.Errorf("toInt64(%v) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("toInt64(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("toInt64(%v) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// Helper function to convert []any to []reflect.Value
func convToValues(args []any) []reflect.Value {
	vals := make([]reflect.Value, len(args))
	for i, arg := range args {
		vals[i] = reflect.ValueOf(arg)
	}
	return vals
}
