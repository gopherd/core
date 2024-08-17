package pair_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gopherd/core/container/pair"
)

func TestPair(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		tests := []struct {
			name     string
			first    any
			second   any
			expected pair.Pair[any, any]
		}{
			{
				name:     "int and string",
				first:    42,
				second:   "hello",
				expected: pair.Pair[any, any]{First: 42, Second: "hello"},
			},
			{
				name:     "float and bool",
				first:    3.14,
				second:   true,
				expected: pair.Pair[any, any]{First: 3.14, Second: true},
			},
			{
				name:     "string and nil",
				first:    "test",
				second:   nil,
				expected: pair.Pair[any, any]{First: "test", Second: nil},
			},
			{
				name:     "nil and nil",
				first:    nil,
				second:   nil,
				expected: pair.Pair[any, any]{First: nil, Second: nil},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := pair.New(tt.first, tt.second)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("New() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("TypeSafety", func(t *testing.T) {
		intStringPair := pair.New(10, "ten")
		if intStringPair.First != 10 || intStringPair.Second != "ten" {
			t.Errorf("Type safety failed for int-string pair")
		}

		floatBoolPair := pair.New(3.14, true)
		if floatBoolPair.First != 3.14 || floatBoolPair.Second != true {
			t.Errorf("Type safety failed for float-bool pair")
		}
	})

	t.Run("ZeroValues", func(t *testing.T) {
		zeroPair := pair.New[int, string](0, "")
		if zeroPair.First != 0 || zeroPair.Second != "" {
			t.Errorf("Zero value test failed: got (%v, %v), want (0, '')", zeroPair.First, zeroPair.Second)
		}
	})

	t.Run("Mutability", func(t *testing.T) {
		p := pair.New(1, "one")
		p.First = 2
		p.Second = "two"
		if p.First != 2 || p.Second != "two" {
			t.Errorf("Mutability test failed: got (%v, %v), want (2, 'two')", p.First, p.Second)
		}
	})

	t.Run("DifferentTypes", func(t *testing.T) {
		type custom struct {
			value int
		}
		p := pair.New([]int{1, 2, 3}, custom{value: 42})
		if !reflect.DeepEqual(p.First, []int{1, 2, 3}) || p.Second.value != 42 {
			t.Errorf("Different types test failed: got (%v, %v), want ([1 2 3], {42})", p.First, p.Second)
		}
	})
}

func TestString(t *testing.T) {
	var tests = []struct {
		name     string
		p        pair.Pair[any, any]
		expected string
	}{
		{
			name:     "IntString",
			p:        pair.New[any, any](42, "hello"),
			expected: "(42,hello)",
		},
		{
			name:     "FloatBool",
			p:        pair.New[any, any](3.14, "world"),
			expected: "(3.14,world)",
		},
		{
			name:     "StringNil",
			p:        pair.New[any, any](42, ""),
			expected: "(42,)",
		},
		{
			name:     "NilNil",
			p:        pair.New[any, any](nil, nil),
			expected: "(<nil>,<nil>)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.p.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	var tests = []struct {
		name     string
		p1       pair.Pair[int, string]
		p2       pair.Pair[int, string]
		expected int
	}{
		{
			name:     "Equal",
			p1:       pair.New(1, "one"),
			p2:       pair.New(1, "one"),
			expected: 0,
		},
		{
			name:     "FirstLess",
			p1:       pair.New(1, "one"),
			p2:       pair.New(2, "two"),
			expected: -1,
		},
		{
			name:     "FirstGreater",
			p1:       pair.New(2, "two"),
			p2:       pair.New(1, "one"),
			expected: 1,
		},
		{
			name:     "SecondLess",
			p1:       pair.New(1, "one"),
			p2:       pair.New(1, "two"),
			expected: -1,
		},
		{
			name:     "SecondGreater",
			p1:       pair.New(1, "two"),
			p2:       pair.New(1, "one"),
			expected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pair.Compare(tt.p1, tt.p2)
			if result != tt.expected {
				t.Errorf("Compare() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCompareFirst(t *testing.T) {
	var tests = []struct {
		name     string
		p1       pair.Pair[int, string]
		p2       pair.Pair[int, string]
		expected int
	}{
		{
			name:     "Equal",
			p1:       pair.New(1, "one"),
			p2:       pair.New(1, "one"),
			expected: 0,
		},
		{
			name:     "FirstLess",
			p1:       pair.New(1, "one"),
			p2:       pair.New(2, "two"),
			expected: -1,
		},
		{
			name:     "FirstGreater",
			p1:       pair.New(2, "two"),
			p2:       pair.New(1, "one"),
			expected: 1,
		},
		{
			name:     "SecondLess",
			p1:       pair.New(1, "one"),
			p2:       pair.New(1, "two"),
			expected: 0,
		},
		{
			name:     "SecondGreater",
			p1:       pair.New(1, "two"),
			p2:       pair.New(1, "one"),
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pair.CompareFirst(tt.p1, tt.p2)
			if result != tt.expected {
				t.Errorf("CompareFirst() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCompareSecond(t *testing.T) {
	var tests = []struct {
		name     string
		p1       pair.Pair[int, string]
		p2       pair.Pair[int, string]
		expected int
	}{
		{
			name:     "Equal",
			p1:       pair.New(1, "one"),
			p2:       pair.New(1, "one"),
			expected: 0,
		},
		{
			name:     "FirstLess",
			p1:       pair.New(1, "one"),
			p2:       pair.New(2, "one"),
			expected: 0,
		},
		{
			name:     "FirstGreater",
			p1:       pair.New(2, "one"),
			p2:       pair.New(1, "one"),
			expected: 0,
		},
		{
			name:     "SecondLess",
			p1:       pair.New(1, "one"),
			p2:       pair.New(1, "two"),
			expected: -1,
		},
		{
			name:     "SecondGreater",
			p1:       pair.New(1, "two"),
			p2:       pair.New(1, "one"),
			expected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pair.CompareSecond(tt.p1, tt.p2)
			if result != tt.expected {
				t.Errorf("CompareSecond() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func ExampleNew() {
	p := pair.New(10, "ten")
	fmt.Printf("First: %v, Second: %v\n", p.First, p.Second)
	// Output: First: 10, Second: ten
}
