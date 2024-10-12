package templates

import (
	"encoding/base64"
	"fmt"
	"html"
	"math"
	"math/big"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gopherd/core/stringutil"
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

// stringFunc converts a function that takes a string to a funtion
// that takes a reflect.Value.
func stringFunc[R any](name string, f func(string) (R, error)) Func {
	return func(v reflect.Value) (reflect.Value, error) {
		s, ok := asString(v)
		if !ok {
			return null, fmt.Errorf("%s: expected string, got %s", name, v.Type())
		}
		r, err := f(s)
		if r, ok := any(r).(reflect.Value); ok {
			return r, err
		}
		return reflect.ValueOf(r), err
	}
}

// stringFunc2 converts a function that takes a T and a string to a funtion
// that takes a T and a reflect.Value.
func stringFunc2[T, R any](name string, f func(T, string) (R, error)) Func2[T] {
	return func(t T, v reflect.Value) (reflect.Value, error) {
		s, ok := asString(v)
		if !ok {
			return reflect.Value{}, fmt.Errorf("%s: expected string, got %s", name, v.Type())
		}
		r, err := f(t, s)
		if r, ok := any(r).(reflect.Value); ok {
			return r, err
		}
		return reflect.ValueOf(r), err
	}
}

// Funcs is a map of utility functions for use in templates
var Funcs = template.FuncMap{
	// @api(_) is a no-op function that returns an empty string.
	// It's useful to place a newline in the template.
	//
	// Example:
	// ```tmpl
	// {{- if true -}}
	// {{_}} ok
	// {{- end}}
	// ```
	//
	// Output:
	// ```
	// 	ok
	// ```
	"_": func(...any) string { return "" },

	// String functions

	// @api(Strings/linespace) adds a newline after a string if it is not empty and does not end with a newline.
	//
	// Example:
	// ```tmpl
	// {{linespace "hello"}}
	// {{linespace ""}}
	// ```
	//
	// Output:
	// ```
	// hello
	//
	// ```
	"linespace": Chain(stringFunc("linespace", noError(linespace))),

	// @api(Strings/fields) splits a string into fields separated by whitespace.
	//
	// Example:
	// ```tmpl
	// {{fields "hello world  !"}}
	// ```
	//
	// Output:
	// ```
	// [hello world !]
	// ```
	"fields": Chain(stringFunc("fields", noError(strings.Fields))),

	// @api(Strings/quote) returns a double-quoted string.
	//
	// Example:
	// ```tmpl
	// {{print "hello"}}
	// {{quote "hello"}}
	// ```
	//
	// Output:
	// ```
	// hello
	// "hello"
	// ```
	"quote": Chain(stringFunc("quote", noError(strconv.Quote))),

	// @api(Strings/unquote) returns an unquoted string.
	//
	// Example:
	// ```tmpl
	// {{unquote "\"hello\""}}
	// ```
	//
	// Output:
	// ```
	// hello
	// ```
	"unquote": Chain(stringFunc("unquote", strconv.Unquote)),

	// @api(Strings/capitalize) capitalizes the first character of a string.
	//
	// Example:
	// ```tmpl
	// {{capitalize "hello"}}
	// ```
	//
	// Output:
	// ```
	// Hello
	// ```
	"capitalize": Chain(stringFunc("capitalize", noError(stringutil.Capitalize))),

	// @api(Strings/lower) converts a string to lowercase.
	//
	// Example:
	// ```tmpl
	// {{lower "HELLO"}}
	// ```
	//
	// Output:
	// ```
	// hello
	// ```
	"lower": Chain(stringFunc("lower", noError(strings.ToLower))),

	// @api(Strings/upper) converts a string to uppercase.
	//
	// Example:
	// ```tmpl
	// {{upper "hello"}}
	// ```
	//
	// Output:
	// ```
	// HELLO
	// ```
	"upper": Chain(stringFunc("upper", noError(strings.ToUpper))),

	// @api(Strings/replace) replaces all occurrences of a substring with another substring.
	//
	// - **Parameters**: (_old_: string, _new_: string, _target_: string)
	//
	// Example:
	// ```tmpl
	// {{replace "o" "0" "hello world"}}
	// ```
	//
	// Output:
	// ```
	// hell0 w0rld
	// ```
	"replace": Chain3(replace),

	// @api(Strings/replaceN) replaces the first n occurrences of a substring with another substring.
	//
	// - **Parameters**: (_old_: string, _new_: string, _n_: int, _target_: string)
	//
	// Example:
	// ```tmpl
	// {{replaceN "o" "0" 1 "hello world"}}
	// ```
	//
	// Output:
	// ```
	// hell0 world
	// ```
	"replaceN": Chain4(replaceN),

	// @api(Strings/trim) removes leading and trailing whitespace from a string.
	//
	// Example:
	// ```tmpl
	// {{trim "  hello  "}}
	// ```
	//
	// Output:
	// ```
	// hello
	// ```
	"trim": Chain(stringFunc("trim", noError(strings.TrimSpace))),

	// @api(Strings/trimPrefix) removes a prefix from a string if it exists.
	//
	// - **Parameters**: (_prefix_: string, _target_: string)
	//
	// Example:
	// ```tmpl
	// {{trimPrefix "Hello, " "Hello, World!"}}
	// ```
	//
	// Output:
	// ```
	// World!
	// ```
	"trimPrefix": Chain2(trimPrefix),

	// @api(Strings/hasPrefix) checks if a string starts with a given prefix.
	//
	// - **Parameters**: (_prefix_: string, _target_: string)
	// - **Returns**: bool
	//
	// Example:
	// ```tmpl
	// {{hasPrefix "Hello" "Hello, World!"}}
	// ```
	//
	// Output:
	// ```
	// true
	// ```
	"hasPrefix": Chain2(hasPrefix),

	// @api(Strings/trimSuffix) removes a suffix from a string if it exists.
	//
	// - **Parameters**: (_suffix_: string, _target_: string)
	//
	// Example:
	// ```tmpl
	// {{trimSuffix ", World!" "Hello, World!"}}
	// ```
	//
	// Output:
	// ```
	// Hello
	// ```
	"trimSuffix": Chain2(trimSuffix),

	// @api(Strings/hasSuffix) checks if a string ends with a given suffix.
	//
	// - **Parameters**: (_suffix_: string, _target_: string)
	// - **Returns**: bool
	//
	// Example:
	// ```tmpl
	// {{hasSuffix "World!" "Hello, World!"}}
	// ```
	//
	// Output:
	// ```
	// true
	// ```
	"hasSuffix": Chain2(hasSuffix),

	// @api(Strings/split) splits a string by a separator.
	//
	// - **Parameters**: (_separator_: string, _target_: string)
	// - **Returns**: slice of strings
	//
	// Example:
	// ```tmpl
	// {{split "," "apple,banana,cherry"}}
	// ```
	//
	// Output:
	// ```
	// [apple banana cherry]
	// ```
	"split": Chain2(split),

	// @api(Strings/join) joins a slice of strings with a separator.
	//
	// - **Parameters**: (_separator_: string, _values_: slice of strings)
	// - **Returns**: string
	//
	// Example:
	// ```tmpl
	// {{join "-" (list "apple" "banana" "cherry")}}
	// ```
	//
	// Output:
	// ```
	// apple-banana-cherry
	// ```
	"join": Chain2(join),

	// @api(Strings/striptags) removes HTML tags from a string.
	//
	// Example:
	// ```tmpl
	// {{striptags "<p>Hello <b>World</b>!</p>"}}
	// ```
	//
	// Output:
	// ```
	// Hello World!
	// ```
	"striptags": Chain(stringFunc("striptags", striptags)),

	// @api(Strings/substr) extracts a substring from a string.
	//
	// - **Parameters**: (_start_: int, _length_: int, _target_: string)
	//
	// Example:
	// ```tmpl
	// {{substr 0 5 "Hello, World!"}}
	// ```
	//
	// Output:
	// ```
	// Hello
	// ```
	"substr": Chain3(substr),

	// @api(Strings/repeat) repeats a string a specified number of times.
	//
	// - **Parameters**: (_count_: int, _target_: string)
	//
	// Example:
	// ```tmpl
	// {{repeat 3 "abc"}}
	// ```
	//
	// Output:
	// ```
	// abcabcabc
	// ```
	"repeat": Chain2(repeat),

	// @api(Strings/camelCase) converts a string to camelCase.
	//
	// Example:
	// ```tmpl
	// {{camelCase "hello world"}}
	// ```
	//
	// Output:
	// ```
	// helloWorld
	// ```
	"camelCase": Chain(stringFunc("camelCase", noError(stringutil.CamelCase))),

	// @api(Strings/pascalCase) converts a string to PascalCase.
	//
	// Example:
	// ```tmpl
	// {{pascalCase "hello world"}}
	// ```
	//
	// Output:
	// ```
	// HelloWorld
	// ```
	"pascalCase": Chain(stringFunc("pascalCase", noError(stringutil.PascalCase))),

	// @api(Strings/snakeCase) converts a string to snake_case.
	//
	// Example:
	// ```tmpl
	// {{snakeCase "helloWorld"}}
	// ```
	//
	// Output:
	// ```
	// hello_world
	// ```
	"snakeCase": Chain(stringFunc("snakeCase", noError(stringutil.SnakeCase))),

	// @api(Strings/kebabCase) converts a string to kebab-case.
	//
	// Example:
	// ```tmpl
	// {{kebabCase "helloWorld"}}
	// ```
	//
	// Output:
	// ```
	// hello-world
	// ```
	"kebabCase": Chain(stringFunc("kebabCase", noError(stringutil.KebabCase))),

	// @api(Strings/truncate) truncates a string to a specified length and adds a suffix if truncated.
	//
	// - **Parameters**: (_length_: int, _suffix_: string, _target_: string)
	//
	// Example:
	// ```tmpl
	// {{truncate 10 "..." "This is a long sentence."}}
	// ```
	//
	// Output:
	// ```
	// This is a...
	// ```
	"truncate": Chain3(truncate),

	// @api(Strings/wordwrap) wraps words in a string to a specified width.
	//
	// - **Parameters**: (_width_: int, _target_: string)
	//
	// Example:
	// ```tmpl
	// {{wordwrap 10 "This is a long sentence that needs wrapping."}}
	// ```
	//
	// Output:
	// ```
	// This is a
	// long
	// sentence
	// that needs
	// wrapping.
	// ```
	"wordwrap": Chain2(wordwrap),

	// @api(Strings/center) centers a string in a field of a given width.
	//
	// - **Parameters**: (_width_: int, _target_: string)
	//
	// Example:
	// ```tmpl
	// {{center 20 "Hello"}}
	// ```
	//
	// Output:
	// ```
	// "       Hello        "
	// ```
	"center": Chain2(center),

	// @api(Strings/matchRegex) checks if a string matches a regular expression.
	//
	// - **Parameters**: (_pattern_: string, _target_: string)
	// - **Returns**: bool
	//
	// Example:
	// ```tmpl
	// {{matchRegex "^[a-z]+$" "hello"}}
	// ```
	//
	// Output:
	// ```
	// true
	// ```
	"matchRegex": Chain2(matchRegex),

	// @api(Strings/html) escapes special characters in a string for use in HTML.
	//
	// Example:
	// ```tmpl
	// {{html "<script>alert('XSS')</script>"}}
	// ```
	//
	// Output:
	// ```
	// &lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;
	// ```
	"html": Chain(stringFunc("html", noError(html.EscapeString))),

	// @api(Strings/urlEscape) escapes a string for use in a URL query.
	//
	// Example:
	// ```tmpl
	// {{urlEscape "hello world"}}
	// ```
	//
	// Output:
	// ```
	// hello+world
	// ```
	"urlEscape": Chain(stringFunc("urlEscape", noError(url.QueryEscape))),

	// @api(Strings/urlUnescape) unescapes a URL query string.
	//
	// Example:
	// ```tmpl
	// {{urlUnescape "hello+world"}}
	// ```
	//
	// Output:
	// ```
	// hello world
	// ```
	"urlUnescape": Chain(stringFunc("urlUnescape", url.QueryUnescape)),

	// Encoding functions

	// @api(Encoding/b64enc) encodes a string to base64.
	//
	// Example:
	// ```tmpl
	// {{b64enc "Hello, World!"}}
	// ```
	//
	// Output:
	// ```
	// SGVsbG8sIFdvcmxkIQ==
	// ```
	"b64enc": Chain(stringFunc("b64enc", noError(b64enc))),

	// @api(Encoding/b64dec) decodes a base64 encoded string.
	//
	// Example:
	// ```tmpl
	// {{b64dec "SGVsbG8sIFdvcmxkIQ=="}}
	// ```
	//
	// Output:
	// ```
	// Hello, World!
	// ```
	"b64dec": Chain(stringFunc("b64dec", b64dec)),

	// Container functions

	"list": list,

	// @api(Container/dict) creates a [Dictionary](#Container/Dictionary) from the given key/value pairs.
	//
	// - **Parameters**: (_dict_or_pairs_: ...any)
	//
	// dict supports keys of any comparable type and values of any type.
	// If an odd number of arguments is provided, it returns an error.
	// If the first argument is already a map, it extends that map with the following key/value pairs.
	//
	// Example:
	// ```tmpl
	// {{$user := dict "name" "Alice" "age" 30}}
	// {{dict "user" ($user) "active" true}}
	// ```
	//
	// Output:
	// ```
	// map[user:map[name:Alice age:30] active:true]
	// ```
	"dict": dict,

	// @api(Container/map) maps a list of values using the given function and returns a list of results.
	//
	// - **Parameters**: (_fn_: function, _list_: slice)
	//
	// Example:
	// ```tmpl
	// {{list 1 2 3 | map (add 1)}}
	// {{list "a" "b" "c" | map (upper | replace "A" "X")}}
	// {{"math/rand.Int63, *io.Reader, *io.Writer" | split "," | map (trim | split "." | first | trimPrefix "*") | sort | uniq}}
	// ```
	//
	// Output:
	// ```
	// [2 3 4]
	// [X B C]
	// [io math/rand]
	// ```
	"map": Chain2(Map),

	// @api(Container/first) returns the first element of a list or string.
	//
	// Example:
	// ```tmpl
	// {{first (list 1 2 3)}}
	// {{first "hello"}}
	// ```
	//
	// Output:
	// ```
	// 1
	// h
	// ```
	"first": Chain(first),

	// @api(Container/last) returns the last element of a list or string.
	//
	// Example:
	// ```tmpl
	// {{last (list 1 2 3)}}
	// {{last "hello"}}
	// ```
	//
	// Output:
	// ```
	// 3
	// o
	// ```
	"last": Chain(last),

	// @api(Container/reverse) reverses a list or string.
	//
	// Example:
	// ```tmpl
	// {{reverse (list 1 2 3)}}
	// {{reverse "hello"}}
	// ```
	//
	// Output:
	// ```
	// [3 2 1]
	// olleh
	// ```
	"reverse": Chain(reverse),

	// @api(Container/sort) sorts a list of numbers or strings.
	//
	// Example:
	// ```tmpl
	// {{sort (list 3 1 4 1 5 9)}}
	// {{sort (list "banana" "apple" "cherry")}}
	// ```
	//
	// Output:
	// ```
	// [1 1 3 4 5 9]
	// [apple banana cherry]
	// ```
	"sort": Chain(sortSlice),

	// @api(Container/uniq) removes duplicate elements from a list.
	//
	// Example:
	// ```tmpl
	// {{uniq (list 1 2 2 3 3 3)}}
	// ```
	//
	// Output:
	// ```
	// [1 2 3]
	// ```
	"uniq": Chain(uniq),

	// @api(Container/includes) checks if an item is present in a list, map, or string.
	//
	// - **Parameters**: (_item_: any, _collection_: slice | map | string)
	// - **Returns**: bool
	//
	// Example:
	// ```tmpl
	// {{includes 2 (list 1 2 3)}}
	// {{includes "world" "hello world"}}
	// ```
	//
	// Output:
	// ```
	// true
	// true
	// ```
	"includes": Chain2(includes),

	// Math functions

	// @api(Math/add) adds two numbers.
	//
	// - **Parameters**: (_a_: number, _b_: number)
	//
	// Example:
	// ```tmpl
	// {{add 2 3}}
	// ```
	//
	// Output:
	// ```
	// 5
	// ```
	"add": Chain2(add),

	// @api(Math/sub) subtracts the second number from the first.
	//
	// - **Parameters**: (_a_: number, _b_: number)
	//
	// Example:
	// ```tmpl
	// {{sub 5 3}}
	// ```
	//
	// Output:
	// ```
	// 2
	// ```
	"sub": Chain2(sub),

	// @api(Math/mul) multiplies two numbers.
	//
	// - **Parameters**: (_a_: number, _b_: number)
	//
	// Example:
	// ```tmpl
	// {{mul 2 3}}
	// ```
	//
	// Output:
	// ```
	// 6
	// ```
	"mul": Chain2(mul),

	// @api(Math/quo) divides the first number by the second.
	//
	// - **Parameters**: (_a_: number, _b_: number)
	//
	// Example:
	// ```tmpl
	// {{quo 6 3}}
	// ```
	//
	// Output:
	// ```
	// 2
	// ```
	"quo": Chain2(quo),

	// @api(Math/rem) returns the remainder of dividing the first number by the second.
	//
	// - **Parameters**: (_a_: number, _b_: number)
	//
	// Example:
	// ```tmpl
	// {{rem 7 3}}
	// ```
	//
	// Output:
	// ```
	// 1
	// ```
	"rem": Chain2(rem),

	// @api(Math/mod) returns the modulus of dividing the first number by the second.
	//
	// - **Parameters**: (_a_: number, _b_: number)
	//
	// Example:
	// ```tmpl
	// {{mod -7 3}}
	// ```
	//
	// Output:
	// ```
	// 2
	// ```
	"mod": Chain2(mod),

	// @api(Math/ceil) returns the least integer value greater than or equal to the input.
	//
	// Example:
	// ```tmpl
	// {{ceil 3.14}}
	// ```
	//
	// Output:
	// ```
	// 4
	// ```
	"ceil": Chain(ceil),

	// @api(Math/floor) returns the greatest integer value less than or equal to the input.
	//
	// Example:
	// ```tmpl
	// {{floor 3.14}}
	// ```
	//
	// Output:
	// ```
	// 3
	// ```
	"floor": Chain(floor),

	// @api(Math/round) rounds a number to a specified number of decimal places.
	//
	// - **Parameters**: (_precision_: integer, _value_: number)
	//
	// Example:
	// ```tmpl
	// {{round 2 3.14159}}
	// ```
	//
	// Output:
	// ```
	// 3.14
	// ```
	"round": Chain2(round),

	// @api(Math/min) returns the minimum of a list of numbers.
	//
	// - **Parameters**: numbers (variadic)
	//
	// Example:
	// ```tmpl
	// {{min 3 1 4 1 5 9}}
	// ```
	//
	// Output:
	// ```
	// 1
	// ```
	"min": minFunc,

	// @api(Math/max) returns the maximum of a list of numbers.
	//
	// - **Parameters**: numbers (variadic)
	//
	// Example:
	// ```tmpl
	// {{max 3 1 4 1 5 9}}
	// ```
	//
	// Output:
	// ```
	// 9
	// ```
	"max": maxFunc,

	// Type conversion functions

	// @api(Convert/int) converts a value to an integer.
	//
	// Example:
	// ```tmpl
	// {{int "42"}}
	// {{int 3.14}}
	// ```
	//
	// Output:
	// ```
	// 42
	// 3
	// ```
	"int": Chain(toInt64),

	// @api(Convert/float) converts a value to a float.
	//
	// Example:
	// ```tmpl
	// {{float "3.14"}}
	// {{float 42}}
	// ```
	//
	// Output:
	// ```
	// 3.14
	// 42
	// ```
	"float": Chain(toFloat64),

	// @api(Convert/string) converts a value to a string.
	//
	// Example:
	// ```tmpl
	// {{string 42}}
	// {{string true}}
	// ```
	//
	// Output:
	// ```
	// 42
	// true
	// ```
	"string": Chain(toString),

	// @api(Convert/bool) converts a value to a boolean.
	//
	// Example:
	// ```tmpl
	// {{bool 1}}
	// {{bool "false"}}
	// ```
	//
	// Output:
	// ```
	// true
	// false
	// ```
	"bool": Chain(toBool),

	// Date functions

	// @api(Date/now) returns the current time.
	//
	// Example:
	// ```tmpl
	// {{now}}
	// ```
	//
	// Output:
	// ```
	// 2024-09-12 15:04:05.999999999 +0000 UTC
	// ```
	"now": time.Now,

	// @api(Date/parseTime) parses a time string using the specified layout.
	//
	// - **Parameters**: (_layout_: string, _value_: string)
	//
	// Example:
	// ```tmpl
	// {{parseTime "2006-01-02" "2024-09-12"}}
	// ```
	//
	// Output:
	// ```
	// 2024-09-12 00:00:00 +0000 UTC
	// ```
	"parseTime": parseTime,

	// OS functions

	// @api(OS/joinPath) joins path elements into a single path.
	//
	// - **Parameters**: elements (variadic)
	//
	// Example:
	// ```tmpl
	// {{joinPath "path" "to" "file.txt"}}
	// ```
	// Output:
	// ```
	// path/to/file.txt
	// ```
	"joinPath": joinPath,

	// @api(OS/splitPath) splits a path into its elements.
	//
	// Example:
	// ```tmpl
	// {{splitPath "path/to/file.txt"}}
	// ```
	// Output:
	// ```
	// [path to file.txt]
	// ```
	"splitPath": Chain(splitPath),

	// @api(OS/absPath) returns the absolute path of a file or directory.
	//
	// Example:
	// ```tmpl
	// {{absPath "file.txt"}}
	// ```
	// Output:
	// ```
	// /path/to/file.txt
	// ```
	"absPath": Chain(stringFunc("absPath", absPath)),

	// @api(OS/relPath) returns the relative path between two paths.
	//
	// - **Parameters**: (_base_: string, _target_: string)
	//
	// Example:
	// ```tmpl
	// {{relPath "/path/to" "/path/to/file.txt"}}
	// ```
	// Output:
	// ```
	// file.txt
	// ```
	"relPath": Chain2(stringFunc2("relPath", relPath)),

	// @api(OS/cleanPath) returns the cleaned path.
	//
	// Example:
	// ```tmpl
	// {{cleanPath "path/to/../file.txt"}}
	// ```
	// Output:
	// ```
	// path/file.txt
	// ```
	"cleanPath": Chain(stringFunc("cleanPath", noError(cleanPath))),

	// @api(OS/basename) returns the last element of a path.
	//
	// Example:
	// ```tmpl
	// {{basename "path/to/file.txt"}}
	// ```
	// Output:
	// ```
	// file.txt
	// ```
	"basename": Chain(stringFunc("basename", noError(basename))),

	// @api(OS/dirname) returns the directory of a path.
	//
	// Example:
	// ```tmpl
	// {{dirname "path/to/file.txt"}}
	// ```
	// Output:
	// ```
	// path/to
	// ```
	"dirname": Chain(stringFunc("dirname", noError(dirname))),

	// @api(OS/extname) returns the extension of a path.
	//
	// Example:
	// ```tmpl
	// {{extname "path/to/file.txt"}}
	// ```
	// Output:
	// ```
	// .txt
	// ```
	"extname": Chain(stringFunc("extname", noError(extname))),

	// @api(OS/removeExt) removes the extension from a path.
	//
	// Example:
	// ```tmpl
	// {{removeExt "path/to/file.txt"}}
	// ```
	// Output:
	// ```
	// path/to/file
	// ```
	"removeExt": Chain(stringFunc("removeExt", noError(removeExt))),

	// @api(OS/isAbs) reports whether a path is absolute.
	//
	// Example:
	// ```tmpl
	// {{isAbs "/path/to/file.txt"}}
	// ```
	// Output:
	// ```
	// true
	// ```
	"isAbs": Chain(stringFunc("isAbs", noError(isAbs))),

	// @api(OS/glob) returns the names of all files matching a pattern.
	//
	// Example:
	// ```tmpl
	// {{glob "/path/to/*.txt"}}
	// ```
	// Output:
	// ```
	// [/path/to/file1.txt /path/to/file2.txt]
	// ```
	"glob": Chain(stringFunc("glob", glob)),

	// @api(OS/matchPath) reports whether a path matches a pattern.
	//
	// - **Parameters**: (_pattern_: string, _path_: string)
	//
	// Example:
	// ```tmpl
	// {{matchPath "/path/to/*.txt" "/path/to/file.txt"}}
	// ```
	//
	// Output:
	// ```
	// true
	// ```
	"matchPath": Chain2(stringFunc2("matchPath", matchPath)),
}

// String functions

func linespace(s string) string {
	if s == "" || s[len(s)-1] == '\n' {
		return s
	}
	return s + "\n"
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

func hasPrefix(prefix string, v String) (Bool, error) {
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

func hasSuffix(suffix string, v String) (Bool, error) {
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

// Container functions

func dict(args ...any) (Dictionary, error) {
	if len(args) == 0 {
		return make(Dictionary), nil
	}
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("dict: odd number of arguments: %d", len(args))
	}
	m := make(Dictionary, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		m[args[i]] = args[i+1]
	}
	return m, nil
}

func list(values ...any) (Vector, error) {
	return values, nil
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

func isInteger(v Any) bool {
	return isInt(v) || isUint(v)
}

func isNumber(v Any) bool {
	return isInt(v) || isUint(v) || isFloat(v)
}

func isBool(v Any) bool {
	return v.Kind() == reflect.Bool
}

func isString(v Any) bool {
	return v.Kind() == reflect.String
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
	case reflect.Invalid:
		return reflect.ValueOf(false), nil
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
	case reflect.Map, reflect.Pointer:
		return reflect.ValueOf(v.IsNil()), nil
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

// Os functions

func joinPath(parts ...string) string {
	return filepath.Join(parts...)
}

func splitPath(path String) (Slice, error) {
	p, ok := asString(path)
	if !ok {
		return null, fmt.Errorf("splitPath: expected string as the argument, got %s", path.Type())
	}
	return reflect.ValueOf(strings.Split(p, string(filepath.Separator))), nil
}

func absPath(path string) (string, error) {
	return filepath.Abs(path)
}

func relPath(basepath, targpath string) (string, error) {
	return filepath.Rel(basepath, targpath)
}

func cleanPath(path string) string {
	return filepath.Clean(path)
}

func basename(path string) string {
	return filepath.Base(path)
}

func dirname(path string) string {
	return filepath.Dir(path)
}

func extname(path string) string {
	return filepath.Ext(path)
}

func removeExt(path string) string {
	if ext := filepath.Ext(path); ext != "" {
		return strings.TrimSuffix(path, ext)
	}
	return path
}

func isAbs(path string) bool {
	return filepath.IsAbs(path)
}

func glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func matchPath(pattern, path string) (bool, error) {
	return filepath.Match(pattern, path)
}
