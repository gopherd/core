package templates

import (
	"encoding/base64"
	"fmt"
	"html"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gopherd/core/text"
)

// Bool is a type alias for reflect.Value. It's used to mark a value as a boolean.
type Bool = reflect.Value

// Int is a type alias for reflect.Value. It's used to mark a value as an integer.
type Int = reflect.Value

// Number is a type alias for reflect.Value. It's used to mark a value as a number.
type Number = reflect.Value

// String is a type alias for reflect.Value. It's used to mark a value as a string.
type String = reflect.Value

// Slice is a type alias for reflect.Value. It's used to mark a value as a slice.
type Slice = reflect.Value

// SliceOrString is a type alias for reflect.Value. It's used to mark a value as a slice or a string.
type SliceOrString = reflect.Value

// Any is a type alias for reflect.Value. It's used to mark a value as any type.
type Any = reflect.Value

var null reflect.Value

// Func is a function that converts a reflect.Value to another reflect.Value.
type Func func(reflect.Value) (reflect.Value, error)

// Func2 is a function that converts a T and a reflect.Value to another reflect.Value.
type Func2[T any] func(T, reflect.Value) (reflect.Value, error)

// Func3 is a function that converts a T1, T2 and a reflect.Value to another reflect.Value.
type Func3[T1, T2 any] func(T1, T2, reflect.Value) (reflect.Value, error)

// Func4 is a function that converts a T1, T2, T3 and a reflect.Value to another reflect.Value.
type Func4[T1, T2, T3 any] func(T1, T2, T3, reflect.Value) (reflect.Value, error)

// FuncChain is a function chain that can be used to convert values.
//
// The function must accept 0 or 1 arguments.
//
// If the function accepts 0 arguments, it returns itself (the function).
// If the function accepts 1 argument and the argument is a FuncChain, it returns a new FuncChain
// that chains the two functions.
// Otherwise, it returns the result of calling the function with the argument.
type FuncChain func(...reflect.Value) (reflect.Value, error)

// Chain returns a FuncChain that chains the given function.
func Chain(f Func) FuncChain {
	var self FuncChain
	self = FuncChain(func(values ...reflect.Value) (reflect.Value, error) {
		return call(self, f, values...)
	})
	return self
}

// Chain2 returns a FuncChain that chains the given function.
func Chain2[T any](f Func2[T]) func(T, ...reflect.Value) (reflect.Value, error) {
	return func(x T, argument ...reflect.Value) (reflect.Value, error) {
		var self FuncChain
		self = FuncChain(func(arg ...reflect.Value) (reflect.Value, error) {
			return call(self, func(y reflect.Value) (reflect.Value, error) {
				return f(x, y)
			}, arg...)
		})
		return self(argument...)
	}
}

// Chain3 returns a FuncChain that chains the given function.
func Chain3[T1, T2 any](f Func3[T1, T2]) func(T1, T2, ...reflect.Value) (reflect.Value, error) {
	return func(x T1, y T2, argument ...reflect.Value) (reflect.Value, error) {
		var self FuncChain
		self = FuncChain(func(arg ...reflect.Value) (reflect.Value, error) {
			return call(self, func(z reflect.Value) (reflect.Value, error) {
				return f(x, y, z)
			}, arg...)
		})
		return self(argument...)
	}
}

// Chain4 returns a FuncChain that chains the given function.
func Chain4[T1, T2, T3 any](f Func4[T1, T2, T3]) func(T1, T2, T3, ...reflect.Value) (reflect.Value, error) {
	return func(x T1, y T2, z T3, argument ...reflect.Value) (reflect.Value, error) {
		var self FuncChain
		self = FuncChain(func(arg ...reflect.Value) (reflect.Value, error) {
			return call(self, func(w reflect.Value) (reflect.Value, error) {
				return f(x, y, z, w)
			}, arg...)
		})
		return self(argument...)
	}
}

// call calls the given function with the given arguments in the function chain.
// If no arguments are given, it returns the function chain itself.
// If one argument is given and it is a FuncChain, it returns a new FuncChain that
// chains the two functions.
// Otherwise, it returns the result of calling the function with the argument.
func call(fc FuncChain, f Func, argument ...reflect.Value) (reflect.Value, error) {
	if len(argument) == 0 {
		// Return self if no arguments are given.
		return reflect.ValueOf(fc), nil
	}
	if len(argument) > 1 {
		return reflect.Value{}, fmt.Errorf("expected 0 or 1 arguments, got %d", len(argument))
	}
	v := argument[0]
	if v.Kind() != reflect.Func {
		return f(v)
	}
	if !v.CanInterface() {
		return reflect.Value{}, fmt.Errorf("function is not exported")
	}
	c, ok := v.Interface().(FuncChain)
	if !ok {
		return reflect.Value{}, fmt.Errorf("expected function chain, got %s", v.Type())
	}
	return reflect.ValueOf(c.then(f)), nil
}

// then returns a new function chain that chains the given function with the current function chain.
func (f FuncChain) then(next Func) FuncChain {
	return Chain(func(v reflect.Value) (reflect.Value, error) {
		v, err := f(v)
		if err != nil {
			return reflect.Value{}, err
		}
		return next(v)
	})
}

// Map maps the given value to a slice of values using the given function chain.
func Map(f FuncChain, v Slice) (Slice, error) {
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return null, fmt.Errorf("map: expected slice or array, got %s", v.Type())
	}
	if v.Len() == 0 {
		return reflect.MakeSlice(reflect.SliceOf(v.Type().Elem()), 0, 0), nil
	}
	var result reflect.Value
	for i := 0; i < v.Len(); i++ {
		r, err := f(v.Index(i))
		if err != nil {
			return null, err
		}
		if i == 0 {
			result = reflect.MakeSlice(reflect.SliceOf(r.Type()), v.Len(), v.Len())
		}
		result.Index(i).Set(r)
	}
	return result, nil
}

// noError returns a function that calls the given function and returns the result and nil.
func noError[T, U any](f func(T) U) func(T) (U, error) {
	return func(s T) (U, error) {
		return f(s), nil
	}
}

// stringFunc converts a function that takes a string and returns a string to a funtion
// that takes a reflect.Value and returns a reflect.Value.
func stringFunc(name string, f func(string) (string, error)) Func {
	return func(v reflect.Value) (reflect.Value, error) {
		s, ok := asString(v)
		if !ok {
			return reflect.Value{}, fmt.Errorf("%s: expected string, got %s", name, v.Type())
		}
		r, err := f(s)
		return reflect.ValueOf(r), err
	}
}

// Funcs is a map of utility functions for use in templates
var Funcs = map[string]any{
	// _ is a no-op function that returns an empty string.
	// It's useful to place a newline in the template.
	"_": func() string { return "" },

	// map maps a list of values using the given function and returns a list of results.
	"map": Chain2(Map),

	// String functions

	"quote":       Chain(stringFunc("quote", noError(strconv.Quote))),
	"unquote":     Chain(stringFunc("unquote", strconv.Unquote)),
	"capitalize":  Chain(stringFunc("capitalize", noError(capitalize))),
	"lower":       Chain(stringFunc("lower", noError(strings.ToLower))),
	"upper":       Chain(stringFunc("upper", noError(strings.ToUpper))),
	"replace":     Chain3(replace),
	"replaceN":    Chain4(replaceN),
	"trim":        Chain(stringFunc("trim", noError(strings.TrimSpace))),
	"trimPrefix":  Chain2(trimPrefix),
	"hasPrefix":   Chain2(hasPrefix),
	"trimSuffix":  Chain2(trimSuffix),
	"hasSuffix":   Chain2(hasSuffix),
	"split":       Chain2(split),
	"join":        Chain2(join),
	"striptags":   Chain(stringFunc("striptags", striptags)),
	"substr":      Chain3(substr),
	"repeat":      Chain2(repeat),
	"camelCase":   Chain(stringFunc("camelCase", noError(camelCase))),
	"pascalCase":  Chain(stringFunc("pascalCase", noError(pascalCase))),
	"snakeCase":   Chain(stringFunc("snakeCase", noError(snakeCase))),
	"kebabCase":   Chain(stringFunc("kebabCase", noError(kebabCase))),
	"truncate":    Chain3(truncate),
	"wordwrap":    Chain2(wordwrap),
	"center":      Chain2(center),
	"matchRegex":  Chain2(matchRegex),
	"html":        Chain(stringFunc("html", noError(html.EscapeString))),
	"urlquery":    Chain(stringFunc("urlquery", noError(url.QueryEscape))),
	"urlUnescape": Chain(stringFunc("urlUnescape", url.QueryUnescape)),

	// Encoding functions

	"b64enc": Chain(stringFunc("b64enc", noError(b64enc))),
	"b64dec": Chain(stringFunc("b64dec", b64dec)),

	// List functions

	"list":     list,
	"first":    Chain(first),
	"last":     Chain(last),
	"reverse":  Chain(reverse),
	"sort":     Chain(sortSlice),
	"uniq":     Chain(uniq),
	"includes": Chain2(includes),

	// Math functions

	"add":   Chain2(add),
	"sub":   Chain2(sub),
	"mul":   Chain2(mul),
	"quo":   Chain2(quo),
	"rem":   Chain2(rem),
	"mod":   Chain2(mod),
	"ceil":  Chain(ceil),
	"floor": Chain(floor),
	"round": Chain2(round),
	"min":   minFunc,
	"max":   maxFunc,

	// Type conversion functions

	"int":    Chain(toInt64),
	"float":  Chain(toFloat64),
	"string": Chain(toString),
	"bool":   Chain(toBool),

	// Date functions

	"now":       time.Now,
	"parseTime": parseTime,
}

// String functions

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	return string(unicode.ToUpper(r[0])) + string(r[1:])
}

func replace(old, new string, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("replace: expected string as third argument, got %s", v.Type())
	}
	return reflect.ValueOf(strings.Replace(s, old, new, -1)), nil
}

func replaceN(old, new string, n int, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("replaceN: expected string as fourth argument, got %s", v.Type())
	}
	return reflect.ValueOf(strings.Replace(s, old, new, n)), nil
}

func trimPrefix(prefix string, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("trimPrefix: expected string as second argument, got %s", v.Type())
	}
	return reflect.ValueOf(strings.TrimPrefix(s, prefix)), nil
}

func hasPrefix(prefix string, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("hasPrefix: expected string as second argument, got %s", v.Type())
	}
	return reflect.ValueOf(strings.HasPrefix(s, prefix)), nil
}

func trimSuffix(suffix string, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("trimSuffix: expected string as second argument, got %s", v.Type())
	}
	return reflect.ValueOf(strings.TrimSuffix(s, suffix)), nil
}

func hasSuffix(suffix string, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("hasSuffix: expected string as second argument, got %s", v.Type())
	}
	return reflect.ValueOf(strings.HasSuffix(s, suffix)), nil
}

func split(sep string, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("split: expected string as second argument, got %s", v.Type())
	}
	return reflect.ValueOf(strings.Split(s, sep)), nil
}

func join(sep string, v String) (String, error) {
	kind := v.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return null, fmt.Errorf("join: expected slice or array as second argument, got %s", v.Type())
	}

	length := v.Len()
	parts := make([]string, length)

	for i := 0; i < length; i++ {
		parts[i] = fmt.Sprint(v.Index(i).Interface())
	}

	return reflect.ValueOf(strings.Join(parts, sep)), nil
}

func striptags(s string) (string, error) {
	return regexp.MustCompile("<[^>]*>").ReplaceAllString(s, ""), nil
}

func substr(start, length int, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("substr: expected string as third argument, got %s", v.Type())
	}
	if start < 0 {
		start = 0
	}
	if length < 0 {
		length = 0
	}
	end := start + length
	if end > len(s) {
		end = len(s)
	}
	if start > end {
		start = end
	}
	return reflect.ValueOf(s[start:end]), nil
}

func repeat(count int, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("repeat: expected string as second argument, got %s", v.Type())
	}
	if count <= 0 {
		return reflect.ValueOf(""), nil
	}
	return reflect.ValueOf(strings.Repeat(s, count)), nil
}

func camelCase(s string) string {
	var result strings.Builder
	capNext := false
	for i, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if i == 0 {
				result.WriteRune(unicode.ToLower(r))
			} else if capNext {
				result.WriteRune(unicode.ToUpper(r))
				capNext = false
			} else {
				result.WriteRune(r)
			}
		} else {
			capNext = true
		}
	}
	return result.String()
}

func pascalCase(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	capNext := true
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if capNext {
				result.WriteRune(unicode.ToUpper(r))
				capNext = false
			} else {
				result.WriteRune(r)
			}
		} else {
			capNext = true
		}
	}
	return result.String()
}

func snakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && (unicode.IsUpper(r) || unicode.IsNumber(r) && !unicode.IsNumber(rune(s[i-1]))) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func kebabCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && (unicode.IsUpper(r) || unicode.IsNumber(r) && !unicode.IsNumber(rune(s[i-1]))) {
			result.WriteRune('-')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func truncate(length int, suffix, v String) (String, error) {
	ss, ok := asString(suffix)
	if !ok {
		return null, fmt.Errorf("truncate: expected string as first argument, got %s", suffix.Type())
	}
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("truncate: expected string as second argument, got %s", v.Type())
	}
	if length <= 0 {
		return reflect.ValueOf(""), nil
	}
	if len(s) <= length {
		return reflect.ValueOf(s), nil
	}
	return reflect.ValueOf(s[:length-len(ss)] + ss), nil
}

func wordwrap(width int, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("wordwrap: expected string, got %s", v.Type())
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return reflect.ValueOf(s), nil
	}
	var lines []string
	var currentLine string
	for _, word := range words {
		if len(currentLine)+len(word) > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		} else {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return reflect.ValueOf(strings.Join(lines, "\n")), nil
}

func center(width int, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("center: expected string, got %s", v.Type())
	}
	if width <= len(s) {
		return reflect.ValueOf(s), nil
	}
	left := (width - len(s)) / 2
	right := width - len(s) - left
	return reflect.ValueOf(strings.Repeat(" ", left) + s + strings.Repeat(" ", right)), nil
}

func matchRegex(pattern string, v String) (String, error) {
	s, ok := asString(v)
	if !ok {
		return null, fmt.Errorf("matchRegex: expected string as second argument, got %s", v.Type())
	}
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return null, err
	}
	return reflect.ValueOf(matched), nil
}

// Encoding functions

func b64enc(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func b64dec(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// List functions

func list(values ...Any) (Slice, error) {
	if len(values) == 0 {
		return reflect.ValueOf([]string{}), nil
	}
	result := reflect.MakeSlice(reflect.SliceOf(values[0].Type()), len(values), len(values))
	for i, v := range values {
		result.Index(i).Set(v)
	}
	return result, nil
}

func first(v Slice) (Any, error) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return null, nil
		}
		return v.Index(0), nil
	default:
		if s, ok := asString(v); ok {
			if len(s) == 0 {
				return null, nil
			}
			return reflect.ValueOf(s[0]), nil
		}
		return null, fmt.Errorf("first: expected slice, array or string, got %s", v.Type())
	}
}

func last(v Slice) (Any, error) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return null, nil
		}
		return v.Index(v.Len() - 1), nil
	default:
		if s, ok := asString(v); ok {
			if len(s) == 0 {
				return null, nil
			}
			return reflect.ValueOf(s[len(s)-1]), nil
		}
		return null, fmt.Errorf("last: expected slice, array or string, got %s", v.Type())
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func reverse(v SliceOrString) (SliceOrString, error) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		length := v.Len()
		reversed := reflect.MakeSlice(v.Type(), length, length)
		for i := 0; i < length; i++ {
			reversed.Index(i).Set(v.Index(length - 1 - i))
		}
		return reversed, nil
	default:
		if s, ok := asString(v); ok {
			return reflect.ValueOf(reverseString(s)), nil
		}
		return null, fmt.Errorf("reverse: expected slice, array or string, got %s", v.Type())
	}
}

func sortSlice(v Slice) (Slice, error) {
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return null, fmt.Errorf("sort: expected slice or array, got %s", v.Type())
	}
	isInt := true
	isUint := true
	isNumber := true
	for i := 0; i < v.Len(); i++ {
		isInt = isInt && v.Index(i).CanInt()
		isUint = isUint && v.Index(i).CanUint()
		isNumber = isNumber && v.Index(i).CanFloat()
	}
	if isUint {
		sorted := make([]uint64, v.Len())
		for i := 0; i < v.Len(); i++ {
			sorted[i] = v.Index(i).Uint()
		}
		slices.Sort(sorted)
		return reflect.ValueOf(sorted), nil
	}
	if isInt {
		sorted := make([]int64, v.Len())
		for i := 0; i < v.Len(); i++ {
			sorted[i] = v.Index(i).Int()
		}
		slices.Sort(sorted)
		return reflect.ValueOf(sorted), nil
	}
	if isNumber {
		sorted := make([]float64, v.Len())
		for i := 0; i < v.Len(); i++ {
			sorted[i] = v.Index(i).Float()
		}
		slices.Sort(sorted)
		return reflect.ValueOf(sorted), nil
	}

	sorted := make([]string, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		if s, ok := asString(v.Index(i)); ok {
			sorted = append(sorted, s)
		} else {
			return null, fmt.Errorf("sort: expected slice of numbers or strings, got %s", v.Type())
		}
	}
	slices.Sort(sorted)
	return reflect.ValueOf(sorted), nil
}

func uniq(v Slice) (Slice, error) {
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return null, fmt.Errorf("uniq: expected slice or array, got %s", v.Type())
	}

	length := v.Len()
	seen := make(map[any]bool)
	uniqueSlice := reflect.MakeSlice(v.Type(), 0, length)

	for i := 0; i < length; i++ {
		elem := v.Index(i)
		if !seen[elem.Interface()] {
			seen[elem.Interface()] = true
			uniqueSlice = reflect.Append(uniqueSlice, elem)
		}
	}

	return uniqueSlice, nil
}

func includes(item Any, collection Slice) (Bool, error) {
	switch collection.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < collection.Len(); i++ {
			if reflect.DeepEqual(item.Interface(), collection.Index(i).Interface()) {
				return reflect.ValueOf(true), nil
			}
		}
		return reflect.ValueOf(false), nil
	case reflect.Map:
		return reflect.ValueOf(collection.MapIndex(item).IsValid()), nil
	default:
		if s, ok := asString(collection); ok {
			if i, ok := asString(item); ok {
				return reflect.ValueOf(text.ContainsWord(s, i)), nil
			}
		}
		return null, fmt.Errorf("includes: expected slice, array, map or string as second argument, got %s", collection.Type())
	}
}

// Math functions

func add(a, b Number) (Number, error) {
	return numericBinaryOp(a, b, func(x, y *big.Float) *big.Float {
		return new(big.Float).Add(x, y)
	})
}

func sub(a, b Number) (Number, error) {
	return numericBinaryOp(a, b, func(x, y *big.Float) *big.Float {
		return new(big.Float).Sub(x, y)
	})
}

func mul(a, b Number) (Number, error) {
	return numericBinaryOp(a, b, func(x, y *big.Float) *big.Float {
		return new(big.Float).Mul(x, y)
	})
}

func quo(a, b Number) (Number, error) {
	return numericBinaryOp(a, b, func(x, y *big.Float) *big.Float {
		if y.Sign() == 0 {
			return nil // Division by zero
		}
		return new(big.Float).Quo(x, y)
	})
}

func rem(a, b Number) (Number, error) {
	return numericBinaryOp(a, b, func(x, y *big.Float) *big.Float {
		if y.Sign() == 0 {
			return nil // Division by zero
		}
		return new(big.Float).Quo(x, y).SetMode(big.ToZero)
	})
}

func mod(a, b Number) (Number, error) {
	return numericBinaryOp(a, b, func(x, y *big.Float) *big.Float {
		if y.Sign() == 0 {
			return nil // Division by zero
		}
		q := new(big.Float).Quo(x, y)
		q.SetMode(big.ToZero)
		q.SetPrec(0)
		return new(big.Float).Sub(x, new(big.Float).Mul(q, y))
	})
}

func minFunc(x Number, y ...Number) (Number, error) {
	minVal := x
	for _, arg := range y {
		result, err := numericCompare(minVal, arg)
		if err != nil {
			return null, err
		}
		if result > 0 {
			minVal = arg
		}
	}

	return minVal, nil
}

func maxFunc(x Number, y ...Number) (Number, error) {
	maxVal := x
	for _, arg := range y {
		result, err := numericCompare(maxVal, arg)
		if err != nil {
			return null, err
		}
		if result < 0 {
			maxVal = arg
		}
	}

	return maxVal, nil
}

func ceil(x Number) (Number, error) {
	f, err := toFloat64(x)
	if err != nil {
		return null, err
	}
	return reflect.ValueOf(math.Ceil(f.Float())), nil
}

func floor(x Number) (Number, error) {
	f, err := toFloat64(x)
	if err != nil {
		return null, err
	}
	return reflect.ValueOf(math.Floor(f.Float())), nil
}

func round(precision int, x Number) (Number, error) {
	f, err := toFloat64(x)
	if err != nil {
		return null, err
	}
	shift := math.Pow10(precision)
	return reflect.ValueOf(math.Round(f.Float()*shift) / shift), nil
}

// Helper functions for numeric operations

func numericBinaryOp(a, b Number, op func(*big.Float, *big.Float) *big.Float) (Number, error) {
	x, err := toBigFloat(a)
	if err != nil {
		return null, err
	}
	y, err := toBigFloat(b)
	if err != nil {
		return null, err
	}

	result := op(x, y)
	if result == nil {
		return null, fmt.Errorf("operation error (possibly division by zero)")
	}

	// Try to convert back to original type if possible
	switch {
	case isInt(a) && isInt(b):
		if i, acc := result.Int64(); acc == big.Exact {
			return reflect.ValueOf(i), nil
		}
	case isUint(a) && isUint(b):
		if u, acc := result.Uint64(); acc == big.Exact {
			return reflect.ValueOf(u), nil
		}
	case isFloat(a) || isFloat(b):
		f, _ := result.Float64()
		return reflect.ValueOf(f), nil
	}

	// If conversion is not possible, return as big.Float
	return reflect.ValueOf(result), nil
}

func numericCompare(a, b Number) (int, error) {
	x, err := toBigFloat(a)
	if err != nil {
		return 0, err
	}
	y, err := toBigFloat(b)
	if err != nil {
		return 0, err
	}
	return x.Cmp(y), nil
}

func toBigFloat(v Number) (*big.Float, error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return new(big.Float).SetInt64(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return new(big.Float).SetUint64(v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return new(big.Float).SetFloat64(v.Float()), nil
	default:
		return nil, fmt.Errorf("unsupported type for numeric operation: %s", v.Type())
	}
}

// Type checking functions

func isInt(v Any) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func isUint(v Any) bool {
	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func isFloat(v Any) bool {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// Type conversion functions

func toInt64(v Any) (Int, error) {
	if v.Kind() == reflect.Int64 {
		return v, nil
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(int64(v.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(int64(v.Float())), nil
	case reflect.String:
		i, err := strconv.ParseInt(v.String(), 10, 64)
		if err != nil {
			return null, err
		}
		return reflect.ValueOf(i), nil
	case reflect.Bool:
		if v.Bool() {
			return reflect.ValueOf(int64(1)), nil
		}
		return reflect.ValueOf(int64(0)), nil
	default:
		return null, fmt.Errorf("cannot convert %s to int", v.Type())
	}
}

func toFloat64(v Any) (Number, error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(float64(v.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(float64(v.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(v.Float()), nil
	case reflect.String:
		f, err := strconv.ParseFloat(v.String(), 64)
		if err != nil {
			return null, err
		}
		return reflect.ValueOf(f), nil
	case reflect.Bool:
		if v.Bool() {
			return reflect.ValueOf(float64(1)), nil
		}
		return reflect.ValueOf(float64(0)), nil
	default:
		return null, fmt.Errorf("cannot convert %s to float", v.Type())
	}
}

func toString(v Any) (String, error) {
	switch v.Kind() {
	case reflect.String:
		return v, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(strconv.FormatUint(v.Uint(), 10)), nil
	case reflect.Float32:
		return reflect.ValueOf(strconv.FormatFloat(v.Float(), 'f', -1, 32)), nil
	case reflect.Float64:
		return reflect.ValueOf(strconv.FormatFloat(v.Float(), 'f', -1, 64)), nil
	case reflect.Bool:
		return reflect.ValueOf(strconv.FormatBool(v.Bool())), nil
	default:
		return null, fmt.Errorf("cannot convert %s to string", v.Type())
	}
}

func toBool(v Any) (Bool, error) {
	if v.Kind() == reflect.Bool {
		return v, nil
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(v.Int() != 0), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(v.Uint() != 0), nil
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		return reflect.ValueOf(f != 0 && !math.IsNaN(f)), nil
	case reflect.String:
		x, err := strconv.ParseBool(v.String())
		if err != nil {
			return null, err
		}
		return reflect.ValueOf(x), nil
	default:
		return null, fmt.Errorf("cannot convert %s to bool", v.Type())
	}
}

func asString(v Any) (string, bool) {
	if v.Kind() == reflect.String {
		return v.String(), true
	}
	return "", false
}

// Date functions

func parseTime(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}
