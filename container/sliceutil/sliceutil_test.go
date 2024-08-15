package sliceutil_test

import (
	"math"
	"reflect"
	"slices"
	"sort"
	"testing"

	"github.com/gopherd/core/container/sliceutil"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		f    func(int) int
		want []int
	}{
		{"Double", []int{1, 2, 3, 4, 5}, func(x int) int { return x * 2 }, []int{2, 4, 6, 8, 10}},
		{"Square", []int{1, 2, 3, 4, 5}, func(x int) int { return x * x }, []int{1, 4, 9, 16, 25}},
		{"EmptySlice", []int{}, func(x int) int { return x }, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceutil.Map(tt.s, tt.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want int
	}{
		{"PositiveNumbers", []int{1, 2, 3, 4, 5}, 15},
		{"NegativeNumbers", []int{-1, -2, -3, -4, -5}, -15},
		{"MixedNumbers", []int{-1, 2, -3, 4, -5}, -3},
		{"EmptySlice", []int{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceutil.Sum(tt.s); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccumulate(t *testing.T) {
	tests := []struct {
		name    string
		s       []int
		f       func(int, int) int
		initial int
		want    int
	}{
		{"Sum", []int{1, 2, 3, 4, 5}, func(acc, x int) int { return acc + x }, 0, 15},
		{"Product", []int{1, 2, 3, 4, 5}, func(acc, x int) int { return acc * x }, 1, 120},
		{"Max", []int{3, 1, 4, 1, 5, 9}, func(acc, x int) int { return max(acc, x) }, math.MinInt, 9},
		{"EmptySlice", []int{}, func(acc, x int) int { return acc + x }, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceutil.Accumulate(tt.s, tt.f, tt.initial); got != tt.want {
				t.Errorf("Accumulate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastIndex(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		v    int
		want int
	}{
		{"Found", []int{1, 2, 3, 4, 5}, 3, 2},
		{"NotFound", []int{1, 2, 3, 4, 5}, 6, -1},
		{"LastOccurrence", []int{1, 2, 3, 2, 4, 5}, 2, 3},
		{"EmptySlice", []int{}, 1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceutil.LastIndex(tt.s, tt.v); got != tt.want {
				t.Errorf("LastIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastIndexFunc(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		f    func(int) bool
		want int
	}{
		{"Found", []int{1, 2, 3, 4, 5}, func(x int) bool { return x > 3 }, 4},
		{"NotFound", []int{1, 2, 3, 4, 5}, func(x int) bool { return x > 5 }, -1},
		{"LastOccurrence", []int{1, 2, 3, 4, 5}, func(x int) bool { return x%2 == 0 }, 3},
		{"EmptySlice", []int{}, func(x int) bool { return true }, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceutil.LastIndexFunc(tt.s, tt.f); got != tt.want {
				t.Errorf("LastIndexFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want []int
	}{
		{"NoDuplicates", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"WithDuplicates", []int{1, 2, 2, 3, 3, 3, 4, 5, 5}, []int{1, 2, 3, 4, 5}},
		{"AllDuplicates", []int{1, 1, 1, 1, 1}, []int{1}},
		{"EmptySlice", []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceutil.Unique(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	shuffled := slices.Clone(original)
	sliceutil.Shuffle(shuffled)

	if reflect.DeepEqual(original, shuffled) {
		t.Errorf("Shuffle() did not change the order of elements")
	}

	if len(original) != len(shuffled) {
		t.Errorf("Shuffle() changed the length of the slice")
	}

	originalSet := make(map[int]bool)
	shuffledSet := make(map[int]bool)
	for i := range original {
		originalSet[original[i]] = true
		shuffledSet[shuffled[i]] = true
	}

	if !reflect.DeepEqual(originalSet, shuffledSet) {
		t.Errorf("Shuffle() changed the elements in the slice")
	}
}

func TestShuffleN(t *testing.T) {
	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	shuffled := slices.Clone(original)
	n := 5
	sliceutil.ShuffleN(shuffled, n)

	// Check that the length hasn't changed
	if len(original) != len(shuffled) {
		t.Errorf("ShuffleN() changed the length of the slice")
	}

	// Check that all original elements are still present
	sortedOriginal := slices.Clone(original)
	sortedShuffled := slices.Clone(shuffled)
	sort.Ints(sortedOriginal)
	sort.Ints(sortedShuffled)
	if !reflect.DeepEqual(sortedOriginal, sortedShuffled) {
		t.Errorf("ShuffleN() changed the elements in the slice")
	}

	// Check that at least one element in the first n positions has changed
	// (There's a very small chance this could fail even with correct implementation)
	changed := false
	for i := 0; i < n; i++ {
		if original[i] != shuffled[i] {
			changed = true
			break
		}
	}
	if !changed {
		t.Errorf("ShuffleN() did not change any of the first %d elements", n)
	}

	// Check that elements after n are not guaranteed to be in their original positions
	// (Again, there's a small chance this could fail even with correct implementation)
	allSame := true
	for i := n; i < len(original); i++ {
		if original[i] != shuffled[i] {
			allSame = false
			break
		}
	}
	if allSame {
		t.Errorf("ShuffleN() did not affect any elements after position %d", n)
	}

	t.Run("PanicOnNegativeN", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ShuffleN() did not panic on negative n")
			}
		}()
		sliceutil.ShuffleN([]int{1, 2, 3}, -1)
	})

	t.Run("PanicOnLargeN", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ShuffleN() did not panic when n > len(slice)")
			}
		}()
		sliceutil.ShuffleN([]int{1, 2, 3}, 4)
	})
}
