//go:build go1.23

package iters_test

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
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

func TestEnumerateMap(t *testing.T) {
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
