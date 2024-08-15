//go:build go1.23

package iters_test

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"testing"

	"github.com/gopherd/core/container/iters"
)

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
			for k, v := range iters.EnumerateKV(tt.input) {
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

			for range iters.EnumerateKV(tt.input) {
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
			input:    iters.EnumerateKV(map[int]string{}),
			expected: 0,
		},
		{
			name:     "non-empty map",
			input:    iters.EnumerateKV(map[int]string{1: "a", 2: "b", 3: "c"}),
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
			input:    iters.EnumerateKV(map[string]int{}),
			expected: 0,
		},
		{
			name:     "non-empty map",
			input:    iters.EnumerateKV(map[string]int{"a": 1, "b": 2, "c": 3}),
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
			input:    iters.EnumerateKV(map[int]string{1: "a", 2: "b", 3: "c"}),
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
			input:    iters.EnumerateKV(map[string]int{"a": 1, "b": 2, "c": 3}),
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
			input: iters.EnumerateKV(map[string]int{"a": 1, "bb": 2, "ccc": 3}),
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
			input: iters.EnumerateKV(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}),
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
			input:    iters.EnumerateKV(map[string]int{"a": 1, "b": 2, "c": 3}),
			f:        func(k string, v int) string { return fmt.Sprintf("%s:%d", k, v) },
			expected: []string{"a:1", "b:2", "c:3"},
		},
		{
			name:     "empty map",
			input:    iters.EnumerateKV(map[string]int{}),
			f:        func(k string, v int) string { return fmt.Sprintf("%s:%d", k, v) },
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]string, 0)
			for v := range iters.MapKV(tt.input, tt.f) {
				result = append(result, v)
			}

			sort.Strings(result)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MapKV() = %v, want %v", result, tt.expected)
			}

			for range iters.MapKV(tt.input, tt.f) {
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
			input:    iters.EnumerateKV(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}),
			f:        func(k string, v int) bool { return v > 2 },
			expected: map[string]int{"c": 3, "d": 4},
		},
		{
			name:     "no matches",
			input:    iters.EnumerateKV(map[string]int{"a": 1, "b": 2}),
			f:        func(k string, v int) bool { return v > 5 },
			expected: map[string]int{},
		},
		{
			name:     "all match",
			input:    iters.EnumerateKV(map[string]int{"a": 1, "b": 2}),
			f:        func(k string, v int) bool { return v > 0 },
			expected: map[string]int{"a": 1, "b": 2},
		},
		{
			name:     "empty map",
			input:    iters.EnumerateKV(map[string]int{}),
			f:        func(k string, v int) bool { return true },
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]int)
			for k, v := range iters.FilterKV(tt.input, tt.f) {
				result[k] = v
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FilterKV() = %v, want %v", result, tt.expected)
			}

			for range iters.FilterKV(tt.input, tt.f) {
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
	for k, v := range iters.EnumerateKV(m) {
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
