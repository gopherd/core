//go:build go1.23

package iters_test

import (
	"fmt"
	"iter"
	"math"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gopherd/core/container/iters"
)

func TestInfinite(t *testing.T) {
	tests := []struct {
		name     string
		f        func(int) int
		expected []int
	}{
		{
			name: "increment",
			f:    func(i int) int { return i + 1 },
			expected: []int{
				1, 2, 3, 4, 5,
			},
		},
		{
			name: "square",
			f:    func(i int) int { return i * i },
			expected: []int{
				0, 1, 4, 9, 16,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]int, 0, len(tt.expected))
			for i := range iters.Infinite() {
				result = append(result, tt.f(i))
				if len(result) == len(tt.expected) {
					break
				}
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Infinite() = %v, want %v", result, tt.expected)
			}
			for range iters.Infinite() {
				break
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		n        int
		expected []interface{}
	}{
		{"Repeat int", 5, 3, []interface{}{5, 5, 5}},
		{"Repeat string", "a", 4, []interface{}{"a", "a", "a", "a"}},
		{"Repeat zero times", 1, 0, []interface{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := iters.Repeat(tt.value, tt.n)
			result := slices.Collect(seq)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("At index %d: expected %v, got %v", i, tt.expected[i], v)
				}
			}
			for range iters.Repeat(tt.value, tt.n) {
				break
			}
		})
	}
}

func TestRepeatPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for negative n, but it didn't panic")
		}
	}()
	iters.Repeat(1, -1)
}

func TestEnumerate(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected [][2]interface{}
	}{
		{
			name:     "Normal slice",
			input:    []string{"a", "b", "c"},
			expected: [][2]interface{}{{0, "a"}, {1, "b"}, {2, "c"}},
		},
		{
			name:     "Empty slice",
			input:    []string{},
			expected: [][2]interface{}{},
		},
		{
			name:     "Single element slice",
			input:    []string{"x"},
			expected: [][2]interface{}{{0, "x"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := iters.Enumerate(tt.input)
			result := make([][2]interface{}, 0, len(tt.input))
			for i, v := range seq {
				result = append(result, [2]interface{}{i, v})
			}
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v[0] != tt.expected[i][0] || v[1] != tt.expected[i][1] {
					t.Errorf("At index %d: expected %v, got %v", i, tt.expected[i], v)
				}
			}
			for range iters.Enumerate(tt.input) {
				break
			}
		})
	}
}

func TestLoop(t *testing.T) {
	tests := []struct {
		name     string
		end      int
		expected []int
	}{
		{"Normal loop", 5, []int{0, 1, 2, 3, 4}},
		{"Zero loop", 0, []int{}},
		{"Single iteration", 1, []int{0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := iters.Loop(tt.end)
			result := make([]int, 0, tt.end)
			for v := range seq {
				result = append(result, v)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("At index %d: expected %d, got %d", i, tt.expected[i], v)
				}
			}
			for range iters.Loop(tt.end) {
				break
			}
		})
	}
}

func TestLoopPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for negative end, but it didn't panic")
		}
	}()
	iters.Loop(-1)
}

func TestRange(t *testing.T) {
	tests := []struct {
		name     string
		start    int
		end      int
		step     int
		expected []int
	}{
		{"Positive step", 1, 10, 2, []int{1, 3, 5, 7, 9}},
		{"Negative step", 10, 1, -2, []int{10, 8, 6, 4, 2}},
		{"Empty range", 5, 5, 1, []int{}},
		{"Single element", 5, 6, 1, []int{5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := iters.Range(tt.start, tt.end, tt.step)
			result := make([]int, 0)
			for v := range seq {
				result = append(result, v)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("At index %d: expected %d, got %d", i, tt.expected[i], v)
				}
			}
			for range iters.Range(tt.start, tt.end, tt.step) {
				break
			}
		})
	}
}

func TestRangePanic(t *testing.T) {
	testCases := []struct {
		name  string
		start int
		end   int
		step  int
	}{
		{"Zero step", 1, 10, 0},
		{"Positive step, end < start", 10, 1, 1},
		{"Negative step, start < end", 1, 10, -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for %s, but it didn't panic", tc.name)
				}
			}()
			iters.Range(tc.start, tc.end, tc.step)
		})
	}
}

func TestSteps(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		start    int
		steps    []int
		expected []int
	}{
		{"No steps", 5, 10, nil, []int{10, 11, 12, 13, 14}},
		{"Single step", 5, 1, []int{2}, []int{1, 3, 5, 7, 9}},
		{"Multiple steps", 6, 1, []int{2, 3, 4}, []int{1, 3, 6, 10, 12, 15}},
		{"Negative steps", 5, 20, []int{-1, -2, -3}, []int{20, 19, 17, 14, 13}},
		{"Zero elements", 0, 1, []int{1}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := iters.Steps(tt.n, tt.start, tt.steps...)
			result := make([]int, 0, tt.n)
			for v := range seq {
				result = append(result, v)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("At index %d: expected %d, got %d", i, tt.expected[i], v)
				}
			}
			for range iters.Steps(tt.n, tt.start, tt.steps...) {
				break
			}
		})
	}
}

func TestStepsPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for negative n, but it didn't panic")
		}
	}()
	iters.Steps(-1, 0)
}

func ExampleRepeat() {
	for v := range iters.Repeat("Go", 3) {
		fmt.Print(v, " ")
	}
	// Output: Go Go Go
}

func ExampleEnumerate() {
	fruits := []string{"apple", "banana", "cherry"}
	for i, v := range iters.Enumerate(fruits) {
		fmt.Printf("%d: %s\n", i, v)
	}
	// Output:
	// 0: apple
	// 1: banana
	// 2: cherry
}

func ExampleLoop() {
	for v := range iters.Loop(3) {
		fmt.Print(v, " ")
	}
	// Output: 0 1 2
}

func ExampleRange() {
	for v := range iters.Range(1, 10, 2) {
		fmt.Print(v, " ")
	}
	// Output: 1 3 5 7 9
}

func ExampleSteps() {
	for v := range iters.Steps(5, 1, 2, 3) {
		fmt.Print(v, " ")
	}
	// Output: 1 3 6 8 11
}

func TestGenericFunctions(t *testing.T) {
	t.Run("Loop with float64", func(t *testing.T) {
		sum := iters.Sum(iters.Loop(5.0))
		expected := 10.0
		if sum != expected {
			t.Errorf("Expected sum %f, got %f", expected, sum)
		}
		for range iters.Loop(5.0) {
			break
		}
	})

	t.Run("Range with int32", func(t *testing.T) {
		sum := iters.Sum(iters.Range(int32(1), int32(5), int32(1)))
		expected := int32(10)
		if sum != expected {
			t.Errorf("Expected sum %d, got %d", expected, sum)
		}
	})

	t.Run("Steps with float32", func(t *testing.T) {
		sum := iters.Sum(iters.Steps(4, float32(1.0), float32(0.5)))
		expected := float32(7.0) // 1 + 1.5 + 2 + 2.5
		if sum != expected {
			t.Errorf("Expected sum %f, got %f", expected, sum)
		}
	})
}

func TestEnumerateKV(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []struct {
			k string
			v int
		}
	}{
		{
			name:  "empty map",
			input: map[string]int{},
			expected: []struct {
				k string
				v int
			}{},
		},
		{
			name:  "non-empty map",
			input: map[string]int{"a": 1, "b": 2, "c": 3},
			expected: []struct {
				k string
				v int
			}{
				{"a", 1}, {"b", 2}, {"c", 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]struct {
				k string
				v int
			}, 0, len(tt.input))
			for k, v := range iters.EnumerateMap(tt.input) {
				result = append(result, struct {
					k string
					v int
				}{k, v})
			}

			// Sort the results and expected values for consistent comparison
			sort.Slice(result, func(i, j int) bool { return result[i].k < result[j].k })
			sort.Slice(tt.expected, func(i, j int) bool { return tt.expected[i].k < tt.expected[j].k })

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("EnumerateKV() = %v, want %v", result, tt.expected)
			}

			for range iters.EnumerateMap(tt.input) {
				break
			}
		})
	}
}

func TestZip(t *testing.T) {
	tests := []struct {
		name     string
		seq1     iter.Seq[int]
		seq2     iter.Seq[string]
		expected []struct {
			v1 int
			v2 string
		}
	}{
		{
			name: "equal length sequences",
			seq1: iters.Of(1, 2, 3),
			seq2: iters.Of("a", "b", "c"),
			expected: []struct {
				v1 int
				v2 string
			}{{1, "a"}, {2, "b"}, {3, "c"}},
		},
		{
			name: "seq1 longer than seq2",
			seq1: iters.Of(1, 2, 3, 4),
			seq2: iters.Of("a", "b"),
			expected: []struct {
				v1 int
				v2 string
			}{{1, "a"}, {2, "b"}, {3, ""}, {4, ""}},
		},
		{
			name: "seq2 longer than seq1",
			seq1: iters.Of(1, 2),
			seq2: iters.Of("a", "b", "c", "d"),
			expected: []struct {
				v1 int
				v2 string
			}{{1, "a"}, {2, "b"}, {0, "c"}, {0, "d"}},
		},
		{
			name: "seq1 is empty",
			seq1: iters.Of[int](),
			seq2: iters.Of("a", "b", "c"),
			expected: []struct {
				v1 int
				v2 string
			}{{0, "a"}, {0, "b"}, {0, "c"}},
		},
		{
			name: "seq2 is empty",
			seq1: iters.Of[int](1, 2, 3),
			seq2: iters.Of[string](),
			expected: []struct {
				v1 int
				v2 string
			}{{1, ""}, {2, ""}, {3, ""}},
		},
		{
			name: "empty sequences",
			seq1: iters.Of[int](),
			seq2: iters.Of[string](),
			expected: []struct {
				v1 int
				v2 string
			}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]struct {
				v1 int
				v2 string
			}, 0)
			for v1, v2 := range iters.Zip(tt.seq1, tt.seq2) {
				result = append(result, struct {
					v1 int
					v2 string
				}{v1, v2})
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Zip() = %v, want %v", result, tt.expected)
			}

			for range iters.Zip(tt.seq1, tt.seq2) {
				break
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		expected int
	}{
		{"empty sequence", iters.Of[int](), 0},
		{"single element", iters.Of(5), 5},
		{"multiple elements", iters.Of(1, 2, 3, 4, 5), 15},
		{"negative numbers", iters.Of(-1, -2, -3), -6},
		{"mixed numbers", iters.Of(-1, 0, 1), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.Sum(tt.input); got != tt.expected {
				t.Errorf("Sum() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSumKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq2[int, string]
		expected int
	}{
		{
			name:     "empty map",
			input:    iters.EnumerateMap(map[int]string{}),
			expected: 0,
		},
		{
			name:     "non-empty map",
			input:    iters.EnumerateMap(map[int]string{1: "a", 2: "b", 3: "c"}),
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.SumKeys(tt.input); got != tt.expected {
				t.Errorf("SumKeys() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSumValues(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq2[string, int]
		expected int
	}{
		{
			name:     "empty map",
			input:    iters.EnumerateMap(map[string]int{}),
			expected: 0,
		},
		{
			name:     "non-empty map",
			input:    iters.EnumerateMap(map[string]int{"a": 1, "b": 2, "c": 3}),
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.SumValues(tt.input); got != tt.expected {
				t.Errorf("SumValues() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAccumulate(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		initial  int
		expected int
	}{
		{"empty sequence", iters.Of[int](), 10, 10},
		{"single element", iters.Of(5), 10, 15},
		{"multiple elements", iters.Of(1, 2, 3), 10, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.Accumulate(tt.input, tt.initial); got != tt.expected {
				t.Errorf("Accumulate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAccumulateFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[string]
		f        func(int, string) int
		initial  int
		expected int
	}{
		{
			name:  "count non-empty strings",
			input: iters.Of("a", "", "b", "c", ""),
			f: func(acc int, s string) int {
				if s != "" {
					return acc + 1
				}
				return acc
			},
			initial:  0,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.AccumulateFunc(tt.input, tt.f, tt.initial); got != tt.expected {
				t.Errorf("AccumulateFunc() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAccumulateKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq2[int, string]
		initial  int
		expected int
	}{
		{
			name:     "sum of keys",
			input:    iters.EnumerateMap(map[int]string{1: "a", 2: "b", 3: "c"}),
			initial:  10,
			expected: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.AccumulateKeys(tt.input, tt.initial); got != tt.expected {
				t.Errorf("AccumulateKeys() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAccumulateValues(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq2[string, int]
		initial  int
		expected int
	}{
		{
			name:     "sum of values",
			input:    iters.EnumerateMap(map[string]int{"a": 1, "b": 2, "c": 3}),
			initial:  10,
			expected: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.AccumulateValues(tt.input, tt.initial); got != tt.expected {
				t.Errorf("AccumulateValues() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAccumulateKeysFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq2[string, int]
		f        func(int, string) int
		initial  int
		expected int
	}{
		{
			name:  "count keys with length > 1",
			input: iters.EnumerateMap(map[string]int{"a": 1, "bb": 2, "ccc": 3}),
			f: func(acc int, k string) int {
				if len(k) > 1 {
					return acc + 1
				}
				return acc
			},
			initial:  0,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.AccumulateKeysFunc(tt.input, tt.f, tt.initial); got != tt.expected {
				t.Errorf("AccumulateKeysFunc() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAccumulateValuesFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq2[string, int]
		f        func(int, int) int
		initial  int
		expected int
	}{
		{
			name:  "sum of even values",
			input: iters.EnumerateMap(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}),
			f: func(acc, v int) int {
				if v%2 == 0 {
					return acc + v
				}
				return acc
			},
			initial:  0,
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.AccumulateValuesFunc(tt.input, tt.f, tt.initial); got != tt.expected {
				t.Errorf("AccumulateValuesFunc() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		expected []int
	}{
		{"no duplicates", iters.Of(1, 2, 3, 4), []int{1, 2, 3, 4}},
		{"with duplicates", iters.Of(1, 2, 2, 3, 3, 3, 4), []int{1, 2, 3, 4}},
		{"all duplicates", iters.Of(1, 1, 1, 1), []int{1}},
		{"empty sequence", iters.Of[int](), []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]int, 0)
			for v := range iters.Unique(tt.input) {
				result = append(result, v)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Unique() = %v, want %v", result, tt.expected)
			}

			for range iters.Unique(tt.input) {
				break
			}
		})
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		f        func(int) string
		expected []string
	}{
		{
			name:     "int to string",
			input:    iters.Of(1, 2, 3),
			f:        func(i int) string { return strconv.Itoa(i) },
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "empty sequence",
			input:    iters.Of[int](),
			f:        func(i int) string { return strconv.Itoa(i) },
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]string, 0)
			for v := range iters.Map(tt.input, tt.f) {
				result = append(result, v)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Map() = %v, want %v", result, tt.expected)
			}

			for range iters.Map(tt.input, tt.f) {
				break
			}
		})
	}
}

func TestMapKV(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq2[string, int]
		f        func(string, int) string
		expected []string
	}{
		{
			name:     "combine key and value",
			input:    iters.EnumerateMap(map[string]int{"a": 1, "b": 2, "c": 3}),
			f:        func(k string, v int) string { return fmt.Sprintf("%s:%d", k, v) },
			expected: []string{"a:1", "b:2", "c:3"},
		},
		{
			name:     "empty map",
			input:    iters.EnumerateMap(map[string]int{}),
			f:        func(k string, v int) string { return fmt.Sprintf("%s:%d", k, v) },
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]string, 0)
			for v := range iters.Map2(tt.input, tt.f) {
				result = append(result, v)
			}

			sort.Strings(result)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MapKV() = %v, want %v", result, tt.expected)
			}

			for range iters.Map2(tt.input, tt.f) {
				break
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		f        func(int) bool
		expected []int
	}{
		{
			name:     "even numbers",
			input:    iters.Of(1, 2, 3, 4, 5, 6),
			f:        func(i int) bool { return i%2 == 0 },
			expected: []int{2, 4, 6},
		},
		{
			name:     "no matches",
			input:    iters.Of(1, 3, 5),
			f:        func(i int) bool { return i%2 == 0 },
			expected: []int{},
		},
		{
			name:     "all match",
			input:    iters.Of(2, 4, 6),
			f:        func(i int) bool { return i%2 == 0 },
			expected: []int{2, 4, 6},
		},
		{
			name:     "empty sequence",
			input:    iters.Of[int](),
			f:        func(i int) bool { return i%2 == 0 },
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]int, 0)
			for v := range iters.Filter(tt.input, tt.f) {
				result = append(result, v)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Filter() = %v, want %v", result, tt.expected)
			}

			for range iters.Filter(tt.input, tt.f) {
				break
			}
		})
	}
}

func TestFilterKV(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq2[string, int]
		f        func(string, int) bool
		expected map[string]int
	}{
		{
			name:     "values greater than 2",
			input:    iters.EnumerateMap(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}),
			f:        func(k string, v int) bool { return v > 2 },
			expected: map[string]int{"c": 3, "d": 4},
		},
		{
			name:     "no matches",
			input:    iters.EnumerateMap(map[string]int{"a": 1, "b": 2}),
			f:        func(k string, v int) bool { return v > 5 },
			expected: map[string]int{},
		},
		{
			name:     "all match",
			input:    iters.EnumerateMap(map[string]int{"a": 1, "b": 2}),
			f:        func(k string, v int) bool { return v > 0 },
			expected: map[string]int{"a": 1, "b": 2},
		},
		{
			name:     "empty map",
			input:    iters.EnumerateMap(map[string]int{}),
			f:        func(k string, v int) bool { return true },
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]int)
			for k, v := range iters.Filter2(tt.input, tt.f) {
				result[k] = v
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FilterKV() = %v, want %v", result, tt.expected)
			}

			for range iters.Filter2(tt.input, tt.f) {
				break
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		target   int
		expected bool
	}{
		{"element present", iters.Of(1, 2, 3, 4, 5), 3, true},
		{"element not present", iters.Of(1, 2, 3, 4, 5), 6, false},
		{"empty sequence", iters.Of[int](), 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.Contains(tt.input, tt.target); got != tt.expected {
				t.Errorf("Contains() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContainsFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		f        func(int) bool
		expected bool
	}{
		{
			name:     "even number present",
			input:    iters.Of(1, 3, 5, 7, 8, 9),
			f:        func(i int) bool { return i%2 == 0 },
			expected: true,
		},
		{
			name:     "no even number present",
			input:    iters.Of(1, 3, 5, 7, 9),
			f:        func(i int) bool { return i%2 == 0 },
			expected: false,
		},
		{
			name:     "empty sequence",
			input:    iters.Of[int](),
			f:        func(i int) bool { return i%2 == 0 },
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.ContainsFunc(tt.input, tt.f); got != tt.expected {
				t.Errorf("ContainsFunc() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		expected int
	}{
		{"non-empty sequence", iters.Of(1, 2, 3, 4, 5), 5},
		{"empty sequence", iters.Of[int](), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.Count(tt.input); got != tt.expected {
				t.Errorf("Count() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCountFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		f        func(int) bool
		expected int
	}{
		{
			name:     "count even numbers",
			input:    iters.Of(1, 2, 3, 4, 5, 6),
			f:        func(i int) bool { return i%2 == 0 },
			expected: 3,
		},
		{
			name:     "no matches",
			input:    iters.Of(1, 3, 5),
			f:        func(i int) bool { return i%2 == 0 },
			expected: 0,
		},
		{
			name:     "all match",
			input:    iters.Of(2, 4, 6),
			f:        func(i int) bool { return i%2 == 0 },
			expected: 3,
		},
		{
			name:     "empty sequence",
			input:    iters.Of[int](),
			f:        func(i int) bool { return i%2 == 0 },
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iters.CountFunc(tt.input, tt.f); got != tt.expected {
				t.Errorf("CountFunc() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGroupBy(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		f        func(int) string
		expected map[string][]int
	}{
		{
			name:  "group by even/odd",
			input: iters.Of(1, 2, 3, 4, 5, 6),
			f: func(i int) string {
				if i%2 == 0 {
					return "even"
				} else {
					return "odd"
				}
			},
			expected: map[string][]int{
				"even": {2, 4, 6},
				"odd":  {1, 3, 5},
			},
		},
		{
			name:     "empty sequence",
			input:    iters.Of[int](),
			f:        func(i int) string { return "group" },
			expected: map[string][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string][]int)
			for k, v := range iters.GroupBy(tt.input, tt.f) {
				result[k] = v
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GroupBy() = %v, want %v", result, tt.expected)
			}

			for range iters.GroupBy(tt.input, tt.f) {
				break
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		expected []int
	}{
		{"non-empty sequence", iters.Of(3, 1, 4, 1, 5, 9, 2, 6, 5, 3), []int{1, 1, 2, 3, 3, 4, 5, 5, 6, 9}},
		{"empty sequence", iters.Of[int](), []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]int, 0)
			for v := range iters.Sort(tt.input) {
				result = append(result, v)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Sort() = %v, want %v", result, tt.expected)
			}

			for range iters.Sort(tt.input) {
				break
			}
		})
	}
}

// Example tests
func ExampleEnumerateKV() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range iters.EnumerateMap(m) {
		fmt.Printf("%s: %d\n", k, v)
	}
	// Unordered output:
	// a: 1
	// b: 2
	// c: 3
}

func ExampleZip() {
	s1 := iters.Of(1, 2, 3)
	s2 := iters.Of("a", "b", "c")
	for v1, v2 := range iters.Zip(s1, s2) {
		fmt.Printf("%d: %s\n", v1, v2)
	}
	// Output:
	// 1: a
	// 2: b
	// 3: c
}

func ExampleUnique() {
	s := iters.Of(1, 2, 2, 3, 3, 3, 4)
	for v := range iters.Unique(s) {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4
}

func ExampleFilter() {
	s := iters.Of(1, 2, 3, 4, 5, 6)
	evenNumbers := iters.Filter(s, func(i int) bool { return i%2 == 0 })
	for v := range evenNumbers {
		fmt.Printf("%d ", v)
	}
	// Output: 2 4 6
}

func ExampleGroupBy() {
	s := iters.Of(1, 2, 3, 4, 5, 6)
	groups := iters.GroupBy(s, func(i int) string {
		if i%2 == 0 {
			return "even"
		}
		return "odd"
	})
	for k, v := range groups {
		fmt.Printf("%s: %v\n", k, v)
	}
	// Unordered output:
	// even: [2 4 6]
	// odd: [1 3 5]
}

func TestMin(t *testing.T) {
	tests := []struct {
		name string
		seq  []int
		want int
	}{
		{"SingleElement", []int{5}, 5},
		{"MultipleElements", []int{3, 1, 4, 1, 5, 9}, 1},
		{"NegativeNumbers", []int{-3, -1, -4, -1, -5, -9}, -9},
		{"MixedNumbers", []int{-3, 0, 4, -1, 5, -9}, -9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := iters.Min(iters.Of(tt.seq...))
			if got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin_EmptySequence(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Min() did not panic on empty sequence")
		}
	}()

	emptySeq := func(yield func(int) bool) {}
	iters.Min(emptySeq)
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		seq  []int
		want int
	}{
		{"SingleElement", []int{5}, 5},
		{"MultipleElements", []int{3, 1, 4, 1, 5, 9}, 9},
		{"NegativeNumbers", []int{-3, -1, -4, -1, -5, -9}, -1},
		{"MixedNumbers", []int{-3, 0, 4, -1, 5, -9}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := iters.Max(iters.Of(tt.seq...))
			if got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax_EmptySequence(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Max() did not panic on empty sequence")
		}
	}()

	emptySeq := func(yield func(int) bool) {}
	iters.Max(emptySeq)
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		name    string
		seq     []int
		wantMin int
		wantMax int
	}{
		{"SingleElement", []int{5}, 5, 5},
		{"MultipleElements", []int{3, 1, 4, 1, 5, 9}, 1, 9},
		{"NegativeNumbers", []int{-3, -1, -4, -1, -5, -9}, -9, -1},
		{"MixedNumbers", []int{-3, 0, 4, -1, 5, -9}, -9, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := iters.MinMax(iters.Of(tt.seq...))
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("MinMax() = (%v, %v), want (%v, %v)", gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestMinMax_EmptySequence(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MinMax() did not panic on empty sequence")
		}
	}()

	emptySeq := func(yield func(int) bool) {}
	iters.MinMax(emptySeq)
}

func TestMinKey(t *testing.T) {
	tests := []struct {
		name string
		seq  map[int]string
		want int
	}{
		{"SingleElement", map[int]string{5: "five"}, 5},
		{"MultipleElements", map[int]string{3: "three", 1: "one", 4: "four", 5: "five", 9: "nine"}, 1},
		{"NegativeNumbers", map[int]string{-3: "minus three", -1: "minus one", -4: "minus four", -5: "minus five"}, -5},
		{"MixedNumbers", map[int]string{-3: "minus three", 0: "zero", 4: "four", -1: "minus one", 5: "five"}, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := iters.MinKey(iters.EnumerateMap(tt.seq))
			if got != tt.want {
				t.Errorf("MinKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxKey(t *testing.T) {
	tests := []struct {
		name string
		seq  map[int]string
		want int
	}{
		{"SingleElement", map[int]string{5: "five"}, 5},
		{"MultipleElements", map[int]string{3: "three", 1: "one", 4: "four", 5: "five", 9: "nine"}, 9},
		{"NegativeNumbers", map[int]string{-3: "minus three", -1: "minus one", -4: "minus four", -5: "minus five"}, -1},
		{"MixedNumbers", map[int]string{-3: "minus three", 0: "zero", 4: "four", -1: "minus one", 5: "five"}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := iters.MaxKey(iters.EnumerateMap(tt.seq))
			if got != tt.want {
				t.Errorf("MaxKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinMaxKey(t *testing.T) {
	tests := []struct {
		name    string
		seq     map[int]string
		wantMin int
		wantMax int
	}{
		{"SingleElement", map[int]string{5: "five"}, 5, 5},
		{"MultipleElements", map[int]string{3: "three", 1: "one", 4: "four", 5: "five", 9: "nine"}, 1, 9},
		{"NegativeNumbers", map[int]string{-3: "minus three", -1: "minus one", -4: "minus four", -5: "minus five"}, -5, -1},
		{"MixedNumbers", map[int]string{-3: "minus three", 0: "zero", 4: "four", -1: "minus one", 5: "five"}, -3, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := iters.MinMaxKey(iters.EnumerateMap(tt.seq))
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("MinMaxKey() = (%v, %v), want (%v, %v)", gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestMinValue(t *testing.T) {
	tests := []struct {
		name string
		seq  map[string]int
		want int
	}{
		{"SingleElement", map[string]int{"five": 5}, 5},
		{"MultipleElements", map[string]int{"three": 3, "one": 1, "four": 4, "five": 5, "nine": 9}, 1},
		{"NegativeNumbers", map[string]int{"minus three": -3, "minus one": -1, "minus four": -4, "minus five": -5}, -5},
		{"MixedNumbers", map[string]int{"minus three": -3, "zero": 0, "four": 4, "minus one": -1, "five": 5}, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := iters.MinValue(iters.EnumerateMap(tt.seq))
			if got != tt.want {
				t.Errorf("MinValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxValue(t *testing.T) {
	tests := []struct {
		name string
		seq  map[string]int
		want int
	}{
		{"SingleElement", map[string]int{"five": 5}, 5},
		{"MultipleElements", map[string]int{"three": 3, "one": 1, "four": 4, "five": 5, "nine": 9}, 9},
		{"NegativeNumbers", map[string]int{"minus three": -3, "minus one": -1, "minus four": -4, "minus five": -5}, -1},
		{"MixedNumbers", map[string]int{"minus three": -3, "zero": 0, "four": 4, "minus one": -1, "five": 5}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := iters.MaxValue(iters.EnumerateMap(tt.seq))
			if got != tt.want {
				t.Errorf("MaxValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinMaxValue(t *testing.T) {
	tests := []struct {
		name    string
		seq     map[string]int
		wantMin int
		wantMax int
	}{
		{"SingleElement", map[string]int{"five": 5}, 5, 5},
		{"MultipleElements", map[string]int{"three": 3, "one": 1, "four": 4, "five": 5, "nine": 9}, 1, 9},
		{"NegativeNumbers", map[string]int{"minus three": -3, "minus one": -1, "minus four": -4, "minus five": -5}, -5, -1},
		{"MixedNumbers", map[string]int{"minus three": -3, "zero": 0, "four": 4, "minus one": -1, "five": 5}, -3, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := iters.MinMaxValue(iters.EnumerateMap(tt.seq))
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("MinMaxValue() = (%v, %v), want (%v, %v)", gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		n     int
		want  [][]int
	}{
		{"EvenSplit", []int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}},
		{"OddSplit", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2, 3}, {4, 5}}},
		{"SingleChunk", []int{1, 2, 3, 4}, 1, [][]int{{1, 2, 3, 4}}},
		{"MoreChunksThanElements", []int{1, 2, 3}, 4, [][]int{{1}, {2}, {3}}},
		{"EmptySlice", []int{}, 3, [][]int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := [][]int{}
			for chunk := range iters.Split(tt.slice, tt.n) {
				got = append(got, chunk)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() = %v, want %v", got, tt.want)
			}
			for range iters.Split(tt.slice, tt.n) {
				break
			}
		})
	}
}

func TestSplit_InvalidN(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Split() did not panic on invalid n")
		}
	}()

	iters.Split([]int{1, 2, 3}, 0)
}

// Example tests

func ExampleMin() {
	min := iters.Min(iters.Of(3, 1, 4, 1, 5, 9))
	fmt.Println(min)
	// Output: 1
}

func ExampleMax() {
	max := iters.Max(iters.Of(3, 1, 4, 1, 5, 9))
	fmt.Println(max)
	// Output: 9
}

func ExampleMinMax() {
	min, max := iters.MinMax(iters.Of(3, 1, 4, 1, 5, 9))
	fmt.Printf("Min: %d, Max: %d\n", min, max)
	// Output: Min: 1, Max: 9
}

func ExampleMinKey() {
	data := map[int]string{3: "three", 1: "one", 4: "four", 5: "five"}
	minKey := iters.MinKey(iters.EnumerateMap(data))
	fmt.Println(minKey)
	// Output: 1
}

func ExampleMaxKey() {
	data := map[int]string{3: "three", 1: "one", 4: "four", 5: "five"}
	maxKey := iters.MaxKey(iters.EnumerateMap(data))
	fmt.Println(maxKey)
	// Output: 5
}

func ExampleMinValue() {
	data := map[string]int{"three": 3, "one": 1, "four": 4, "five": 5}
	minValue := iters.MinValue(iters.EnumerateMap(data))
	fmt.Println(minValue)
	// Output: 1
}

func ExampleMaxValue() {
	data := map[string]int{"three": 3, "one": 1, "four": 4, "five": 5}
	maxValue := iters.MaxValue(iters.EnumerateMap(data))
	fmt.Println(maxValue)
	// Output: 5
}

func ExampleSplit() {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	chunks := iters.Split(slice, 3)
	for chunk := range chunks {
		fmt.Println(chunk)
	}
	// Output:
	// [1 2 3 4]
	// [5 6 7]
	// [8 9 10]
}

func TestMinMaxWithInfinity(t *testing.T) {
	values := []float64{math.Inf(1), math.Inf(-1), 0, 1, -1}
	min, max := iters.MinMax(iters.Of(values...))
	if min != math.Inf(-1) {
		t.Errorf("Min should be negative infinity, got %v", min)
	}
	if max != math.Inf(1) {
		t.Errorf("Max should be positive infinity, got %v", max)
	}
}

func TestSplitWithLargeN(t *testing.T) {
	slice := []int{1, 2, 3}
	chunks := iters.Split(slice, 1000000)
	count := 0
	for chunk := range chunks {
		count++
		if len(chunk) != 1 {
			t.Errorf("Expected chunk size 1, got %d", len(chunk))
		}
	}
	if count != 3 {
		t.Errorf("Expected 3 chunks, got %d", count)
	}
}

func TestMinMaxWithAllEqualValues(t *testing.T) {
	seq := func(yield func(int) bool) {
		for i := 0; i < 1000; i++ {
			if !yield(42) {
				return
			}
		}
	}
	min, max := iters.MinMax(seq)
	if min != 42 || max != 42 {
		t.Errorf("Expected min and max to be 42, got min=%d, max=%d", min, max)
	}
}

func TestSplitWithEmptySlice(t *testing.T) {
	slice := []int{}
	chunks := iters.Split(slice, 5)
	count := 0
	for range chunks {
		count++
	}
	if count != 0 {
		t.Errorf("Expected 0 chunks for empty slice, got %d", count)
	}
}

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		seq  map[string]int
		want []string
	}{
		{"EmptyMap", map[string]int{}, []string{}},
		{"SingleElement", map[string]int{"a": 1}, []string{"a"}},
		{"MultipleElements", map[string]int{"a": 1, "b": 2, "c": 3}, []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]string, 0, len(tt.seq))
			for k := range iters.Keys(iters.EnumerateMap(tt.seq)) {
				got = append(got, k)
			}

			// Sort both slices for comparison
			sort.Strings(got)
			sort.Strings(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}

			for range iters.Keys(iters.EnumerateMap(tt.seq)) {
				break
			}
		})
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		name string
		seq  map[string]int
		want []int
	}{
		{"EmptyMap", map[string]int{}, []int{}},
		{"SingleElement", map[string]int{"a": 1}, []int{1}},
		{"MultipleElements", map[string]int{"a": 1, "b": 2, "c": 3}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]int, 0, len(tt.seq))
			for v := range iters.Values(iters.EnumerateMap(tt.seq)) {
				got = append(got, v)
			}

			// Sort both slices for comparison
			sort.Ints(got)
			sort.Ints(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}

			for range iters.Values(iters.EnumerateMap(tt.seq)) {
				break
			}
		})
	}
}

func ExampleKeys() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := make([]string, 0, 3)
	for k := range iters.Keys(iters.EnumerateMap(m)) {
		result = append(result, k)
	}

	// Sort the result for consistent output
	sort.Strings(result)
	fmt.Println(result)
	// Output: [a b c]
}

func ExampleValues() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := make([]int, 0, 3)
	for v := range iters.Values(iters.EnumerateMap(m)) {
		result = append(result, v)
	}

	// Sort the result for consistent output
	sort.Ints(result)
	fmt.Println(result)
	// Output: [1 2 3]
}

func TestKeysWithLargeMap(t *testing.T) {
	largeMap := make(map[int]string)
	for i := 0; i < 1000000; i++ {
		largeMap[i] = fmt.Sprintf("value%d", i)
	}

	count := 0
	for range iters.Keys(iters.EnumerateMap(largeMap)) {
		count++
	}

	if count != 1000000 {
		t.Errorf("Expected 1000000 keys, got %d", count)
	}
}

func TestValuesWithLargeMap(t *testing.T) {
	largeMap := make(map[int]string)
	for i := 0; i < 1000000; i++ {
		largeMap[i] = fmt.Sprintf("value%d", i)
	}

	count := 0
	for range iters.Values(iters.EnumerateMap(largeMap)) {
		count++
	}

	if count != 1000000 {
		t.Errorf("Expected 1000000 values, got %d", count)
	}
}

func TestUniqueFunc(t *testing.T) {
	tests := []struct {
		name string
		seq  []int
		eq   func(int, int) bool
		want []int
	}{
		{
			name: "StandardEquality",
			seq:  []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4},
			eq:   func(a, b int) bool { return a == b },
			want: []int{1, 2, 3, 4},
		},
		{
			name: "EmptySequence",
			seq:  []int{},
			eq:   func(a, b int) bool { return a == b },
			want: []int{},
		},
		{
			name: "AllUnique",
			seq:  []int{1, 2, 3, 4, 5},
			eq:   func(a, b int) bool { return a == b },
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "AllSame",
			seq:  []int{1, 1, 1, 1, 1},
			eq:   func(a, b int) bool { return a == b },
			want: []int{1},
		},
		{
			name: "CustomEquality",
			seq:  []int{1, 11, 2, 22, 3, 33},
			eq:   func(a, b int) bool { return a%10 == b%10 },
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uniqueSeq := iters.UniqueFunc(iters.Of(tt.seq...), tt.eq)
			got := make([]int, 0)
			for v := range uniqueSeq {
				got = append(got, v)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueFunc() = %v, want %v", got, tt.want)
			}

			for range iters.UniqueFunc(iters.Of(tt.seq...), tt.eq) {
				break
			}
		})
	}
}

func TestUniqueFuncWithStrings(t *testing.T) {
	words := []string{"hello", "HELLO", "world", "WORLD"}
	eq := func(a, b string) bool {
		return strings.ToLower(a) == strings.ToLower(b)
	}

	uniqueSeq := iters.UniqueFunc(iters.Of(words...), eq)
	got := make([]string, 0)
	for v := range uniqueSeq {
		got = append(got, v)
	}

	want := []string{"hello", "world"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("UniqueFunc() with case-insensitive comparison = %v, want %v", got, want)
	}
}

func ExampleUniqueFunc() {
	numbers := []int{1, 1, 2, 3, 3, 4, 4, 4, 5}
	eq := func(a, b int) bool { return a == b }

	uniqueSeq := iters.UniqueFunc(iters.Of(numbers...), eq)
	for v := range uniqueSeq {
		fmt.Print(v, " ")
	}
	// Output: 1 2 3 4 5
}

func TestUniqueFuncWithCustomType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Bob", 25},
		{"Charlie", 35},
		{"Charlie", 26},
	}

	eq := func(a, b Person) bool {
		return a.Name == b.Name && a.Age == b.Age
	}

	uniqueSeq := iters.UniqueFunc(iters.Of(people...), eq)
	got := make([]Person, 0)
	for v := range uniqueSeq {
		got = append(got, v)
	}

	want := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
		{"Charlie", 26},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("UniqueFunc() with custom type = %v, want %v", got, want)
	}
}
