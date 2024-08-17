//go:build go1.23

package iters_test

import (
	"fmt"
	"iter"
	"math"
	"testing"

	"github.com/gopherd/core/container/iters"
	"github.com/gopherd/core/container/pair"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		expected int
	}{
		{"empty sequence", iters.List[int](), 0},
		{"single element", iters.List(5), 5},
		{"multiple elements", iters.List(1, 2, 3, 4, 5), 15},
		{"negative numbers", iters.List(-1, -2, -3), -6},
		{"mixed numbers", iters.List(-1, 0, 1), 0},
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
			input:    iters.Enumerate2(map[int]string{}),
			expected: 0,
		},
		{
			name:     "non-empty map",
			input:    iters.Enumerate2(map[int]string{1: "a", 2: "b", 3: "c"}),
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
			input:    iters.Enumerate2(map[string]int{}),
			expected: 0,
		},
		{
			name:     "non-empty map",
			input:    iters.Enumerate2(map[string]int{"a": 1, "b": 2, "c": 3}),
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
		{"empty sequence", iters.List[int](), 10, 10},
		{"single element", iters.List(5), 10, 15},
		{"multiple elements", iters.List(1, 2, 3), 10, 16},
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
			input: iters.List("a", "", "b", "c", ""),
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
			input:    iters.Enumerate2(map[int]string{1: "a", 2: "b", 3: "c"}),
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
			input:    iters.Enumerate2(map[string]int{"a": 1, "b": 2, "c": 3}),
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
			input: iters.Enumerate2(map[string]int{"a": 1, "bb": 2, "ccc": 3}),
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
			input: iters.Enumerate2(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}),
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

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		target   int
		expected bool
	}{
		{"element present", iters.List(1, 2, 3, 4, 5), 3, true},
		{"element not present", iters.List(1, 2, 3, 4, 5), 6, false},
		{"empty sequence", iters.List[int](), 1, false},
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
			input:    iters.List(1, 3, 5, 7, 8, 9),
			f:        func(i int) bool { return i%2 == 0 },
			expected: true,
		},
		{
			name:     "no even number present",
			input:    iters.List(1, 3, 5, 7, 9),
			f:        func(i int) bool { return i%2 == 0 },
			expected: false,
		},
		{
			name:     "empty sequence",
			input:    iters.List[int](),
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
		{"non-empty sequence", iters.List(1, 2, 3, 4, 5), 5},
		{"empty sequence", iters.List[int](), 0},
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
			input:    iters.List(1, 2, 3, 4, 5, 6),
			f:        func(i int) bool { return i%2 == 0 },
			expected: 3,
		},
		{
			name:     "no matches",
			input:    iters.List(1, 3, 5),
			f:        func(i int) bool { return i%2 == 0 },
			expected: 0,
		},
		{
			name:     "all match",
			input:    iters.List(2, 4, 6),
			f:        func(i int) bool { return i%2 == 0 },
			expected: 3,
		},
		{
			name:     "empty sequence",
			input:    iters.List[int](),
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
			got := iters.Min(iters.List(tt.seq...))
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
			got := iters.Max(iters.List(tt.seq...))
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
			gotMin, gotMax := iters.MinMax(iters.List(tt.seq...))
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
			got := iters.MinKey(iters.Enumerate2(tt.seq))
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
			got := iters.MaxKey(iters.Enumerate2(tt.seq))
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
			gotMin, gotMax := iters.MinMaxKey(iters.Enumerate2(tt.seq))
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
			got := iters.MinValue(iters.Enumerate2(tt.seq))
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
			got := iters.MaxValue(iters.Enumerate2(tt.seq))
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
			gotMin, gotMax := iters.MinMaxValue(iters.Enumerate2(tt.seq))
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("MinMaxValue() = (%v, %v), want (%v, %v)", gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

// Example tests

func ExampleMin() {
	min := iters.Min(iters.List(3, 1, 4, 1, 5, 9))
	fmt.Println(min)
	// Output: 1
}

func ExampleMax() {
	max := iters.Max(iters.List(3, 1, 4, 1, 5, 9))
	fmt.Println(max)
	// Output: 9
}

func ExampleMinMax() {
	min, max := iters.MinMax(iters.List(3, 1, 4, 1, 5, 9))
	fmt.Printf("Min: %d, Max: %d\n", min, max)
	// Output: Min: 1, Max: 9
}

func ExampleMinKey() {
	data := map[int]string{3: "three", 1: "one", 4: "four", 5: "five"}
	minKey := iters.MinKey(iters.Enumerate2(data))
	fmt.Println(minKey)
	// Output: 1
}

func ExampleMaxKey() {
	data := map[int]string{3: "three", 1: "one", 4: "four", 5: "five"}
	maxKey := iters.MaxKey(iters.Enumerate2(data))
	fmt.Println(maxKey)
	// Output: 5
}

func ExampleMinValue() {
	data := map[string]int{"three": 3, "one": 1, "four": 4, "five": 5}
	minValue := iters.MinValue(iters.Enumerate2(data))
	fmt.Println(minValue)
	// Output: 1
}

func ExampleMaxValue() {
	data := map[string]int{"three": 3, "one": 1, "four": 4, "five": 5}
	maxValue := iters.MaxValue(iters.Enumerate2(data))
	fmt.Println(maxValue)
	// Output: 5
}

func TestMinMaxWithInfinity(t *testing.T) {
	values := []float64{math.Inf(1), math.Inf(-1), 0, 1, -1}
	min, max := iters.MinMax(iters.List(values...))
	if min != math.Inf(-1) {
		t.Errorf("Min should be negative infinity, got %v", min)
	}
	if max != math.Inf(1) {
		t.Errorf("Max should be positive infinity, got %v", max)
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

func TestAppendSeq2(t *testing.T) {
	result := []pair.Pair[int, string]{pair.New(0, "")}
	seq := iters.Enumerate2(map[int]string{1: "a", 2: "b", 3: "c"})
	result = iters.AppendSeq2(result, seq)
	values := []string{"", "a", "b", "c"}
	for i := 0; i < len(result); i++ {
		if result[i].First != i {
			t.Errorf("First = %v, want %v", result[i].First, i)
		}
		if result[i].Second != values[i] {
			t.Errorf("Second = %v, want %v", result[i].Second, values[i])
		}
	}
}

func TestCollect2(t *testing.T) {
	seq := iters.Enumerate2(map[int]string{1: "a", 2: "b", 3: "c"})
	result := iters.Collect2(seq)
	values := []string{"a", "b", "c"}
	for i := 0; i < len(result); i++ {
		if result[i].First != i+1 {
			t.Errorf("First = %v, want %v", result[i].First, i+1)
		}
		if result[i].Second != values[i] {
			t.Errorf("Second = %v, want %v", result[i].Second, values[i])
		}
	}
}
