package slices_test

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/gopherd/core/container/slices"
)

func TestMin(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want int
	}{
		{"SingleElement", []int{5}, 5},
		{"MultipleElements", []int{3, 1, 4, 1, 5, 9}, 1},
		{"NegativeNumbers", []int{-3, -1, -4, -1, -5, -9}, -9},
		{"MixedNumbers", []int{-3, 1, -4, 1, 5, -9}, -9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Min(tt.s); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("EmptySlice", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Min() did not panic on empty slice")
			}
		}()
		slices.Min([]int{})
	})
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want int
	}{
		{"SingleElement", []int{5}, 5},
		{"MultipleElements", []int{3, 1, 4, 1, 5, 9}, 9},
		{"NegativeNumbers", []int{-3, -1, -4, -1, -5, -9}, -1},
		{"MixedNumbers", []int{-3, 1, -4, 1, 5, -9}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Max(tt.s); got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("EmptySlice", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Max() did not panic on empty slice")
			}
		}()
		slices.Max([]int{})
	})
}

func TestMinmax(t *testing.T) {
	tests := []struct {
		name    string
		s       []int
		wantMin int
		wantMax int
	}{
		{"SingleElement", []int{5}, 5, 5},
		{"MultipleElements", []int{3, 1, 4, 1, 5, 9}, 1, 9},
		{"NegativeNumbers", []int{-3, -1, -4, -1, -5, -9}, -9, -1},
		{"MixedNumbers", []int{-3, 1, -4, 1, 5, -9}, -9, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := slices.Minmax(tt.s)
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("Minmax() = (%v, %v), want (%v, %v)", gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}

	t.Run("EmptySlice", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Minmax() did not panic on empty slice")
			}
		}()
		slices.Minmax([]int{})
	})
}

func TestMinFunc(t *testing.T) {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	tests := []struct {
		name string
		s    []int
		f    func(int) int
		want int
	}{
		{"Identity", []int{3, 1, 4, 1, 5, 9}, func(x int) int { return x }, 1},
		{"Absolute", []int{-3, 1, -4, 1, 5, -9}, abs, 1},
		{"Square", []int{-3, 1, -4, 1, 5, -9}, func(x int) int { return x * x }, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.MinFunc(tt.s, tt.f); got != tt.want {
				t.Errorf("MinFunc() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("EmptySlice", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MinFunc() did not panic on empty slice")
			}
		}()
		slices.MinFunc([]int{}, func(x int) int { return x })
	})
}

func TestMaxFunc(t *testing.T) {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	tests := []struct {
		name string
		s    []int
		f    func(int) int
		want int
	}{
		{"Identity", []int{3, 1, 4, 1, 5, 9}, func(x int) int { return x }, 9},
		{"Absolute", []int{-3, 1, -4, 1, 5, -9}, abs, 9},
		{"Square", []int{-3, 1, -4, 1, 5, -9}, func(x int) int { return x * x }, 81},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.MaxFunc(tt.s, tt.f); got != tt.want {
				t.Errorf("MaxFunc() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("EmptySlice", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MaxFunc() did not panic on empty slice")
			}
		}()
		slices.MaxFunc([]int{}, func(x int) int { return x })
	})
}

func TestMinmaxFunc(t *testing.T) {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	tests := []struct {
		name    string
		s       []int
		f       func(int) int
		wantMin int
		wantMax int
	}{
		{"Identity", []int{3, 1, 4, 1, 5, 9}, func(x int) int { return x }, 1, 9},
		{"Absolute", []int{-3, 1, -4, 1, 5, -9}, abs, 1, 9},
		{"Square", []int{-3, 1, -4, 1, 5, -9}, func(x int) int { return x * x }, 1, 81},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := slices.MinmaxFunc(tt.s, tt.f)
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("MinmaxFunc() = (%v, %v), want (%v, %v)", gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}

	t.Run("EmptySlice", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MinmaxFunc() did not panic on empty slice")
			}
		}()
		slices.MinmaxFunc([]int{}, func(x int) int { return x })
	})
}

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
			if got := slices.Map(tt.s, tt.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyFunc(t *testing.T) {
	tests := []struct {
		name string
		dst  []int
		src  []int
		f    func(int) int
		want []int
	}{
		{"Double", make([]int, 5), []int{1, 2, 3, 4, 5}, func(x int) int { return x * 2 }, []int{2, 4, 6, 8, 10}},
		{"Square", make([]int, 5), []int{1, 2, 3, 4, 5}, func(x int) int { return x * x }, []int{1, 4, 9, 16, 25}},
		{"EmptySlice", []int{}, []int{}, func(x int) int { return x }, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slices.CopyFunc(tt.dst, tt.src, tt.f)
			if !reflect.DeepEqual(tt.dst, tt.want) {
				t.Errorf("CopyFunc() resulted in %v, want %v", tt.dst, tt.want)
			}
		})
	}

	t.Run("DestShorterThanSrc", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("CopyFunc() did not panic when destination is shorter than source")
			}
		}()
		slices.CopyFunc(make([]int, 2), []int{1, 2, 3}, func(x int) int { return x })
	})
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
			if got := slices.Sum(tt.s); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumFunc(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		f    func(int) int
		want int
	}{
		{"Identity", []int{1, 2, 3, 4, 5}, func(x int) int { return x }, 15},
		{"Double", []int{1, 2, 3, 4, 5}, func(x int) int { return x * 2 }, 30},
		{"Square", []int{1, 2, 3, 4, 5}, func(x int) int { return x * x }, 55},
		{"EmptySlice", []int{}, func(x int) int { return x }, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.SumFunc(tt.s, tt.f); got != tt.want {
				t.Errorf("SumFunc() = %v, want %v", got, tt.want)
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
			if got := slices.Accumulate(tt.s, tt.f, tt.initial); got != tt.want {
				t.Errorf("Accumulate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMean(t *testing.T) {
	tests := []struct {
		name string
		s    []float64
		want float64
	}{
		{"PositiveNumbers", []float64{1, 2, 3, 4, 5}, 3},
		{"NegativeNumbers", []float64{-1, -2, -3, -4, -5}, -3},
		{"MixedNumbers", []float64{-1, 2, -3, 4, -5}, -0.6},
		{"EmptySlice", []float64{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Mean(tt.s); math.Abs(got-tt.want) > 1e-6 {
				t.Errorf("Mean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeanFunc(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		f    func(int) float64
		want float64
	}{
		{"Identity", []int{1, 2, 3, 4, 5}, func(x int) float64 { return float64(x) }, 3},
		{"Double", []int{1, 2, 3, 4, 5}, func(x int) float64 { return float64(x * 2) }, 6},
		{"Square", []int{1, 2, 3, 4, 5}, func(x int) float64 { return float64(x * x) }, 11},
		{"EmptySlice", []int{}, func(x int) float64 { return float64(x) }, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.MeanFunc(tt.s, tt.f); math.Abs(got-tt.want) > 1e-6 {
				t.Errorf("MeanFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name string
		x    []int
		y    []int
		want bool
	}{
		{"EqualSlices", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"DifferentLength", []int{1, 2, 3}, []int{1, 2}, false},
		{"DifferentContent", []int{1, 2, 3}, []int{1, 2, 4}, false},
		{"EmptySlices", []int{}, []int{}, true},
		{"NilSlices", nil, nil, true},
		{"NilAndEmpty", nil, []int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Equal(tt.x, tt.y); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualFunc(t *testing.T) {
	tests := []struct {
		name string
		x    []int
		y    []int
		f    func(int, int) bool
		want bool
	}{
		{"EqualSlices", []int{1, 2, 3}, []int{1, 2, 3}, func(a, b int) bool { return a == b }, true},
		{"AbsoluteEqual", []int{1, -2, 3}, []int{1, 2, -3}, func(a, b int) bool { return abs(a) == abs(b) }, true},
		{"DifferentLength", []int{1, 2, 3}, []int{1, 2}, func(a, b int) bool { return a == b }, false},
		{"EmptySlices", []int{}, []int{}, func(a, b int) bool { return a == b }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.EqualFunc(tt.x, tt.y, tt.f); got != tt.want {
				t.Errorf("EqualFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		v    int
		want int
	}{
		{"Found", []int{1, 2, 3, 4, 5}, 3, 2},
		{"NotFound", []int{1, 2, 3, 4, 5}, 6, -1},
		{"FirstOccurrence", []int{1, 2, 3, 2, 4, 5}, 2, 1},
		{"EmptySlice", []int{}, 1, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Index(tt.s, tt.v); got != tt.want {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexFunc(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		f    func(int) bool
		want int
	}{
		{"Found", []int{1, 2, 3, 4, 5}, func(x int) bool { return x > 3 }, 3},
		{"NotFound", []int{1, 2, 3, 4, 5}, func(x int) bool { return x > 5 }, -1},
		{"FirstOccurrence", []int{1, 2, 3, 4, 5}, func(x int) bool { return x%2 == 0 }, 1},
		{"EmptySlice", []int{}, func(x int) bool { return true }, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.IndexFunc(tt.s, tt.f); got != tt.want {
				t.Errorf("IndexFunc() = %v, want %v", got, tt.want)
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
			if got := slices.LastIndex(tt.s, tt.v); got != tt.want {
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
			if got := slices.LastIndexFunc(tt.s, tt.f); got != tt.want {
				t.Errorf("LastIndexFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		v    int
		want bool
	}{
		{"Found", []int{1, 2, 3, 4, 5}, 3, true},
		{"NotFound", []int{1, 2, 3, 4, 5}, 6, false},
		{"EmptySlice", []int{}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Contains(tt.s, tt.v); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShrink(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want []int
	}{
		{"FullSlice", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"SliceWithExtraCapacity", append(make([]int, 0, 10), 1, 2, 3), []int{1, 2, 3}},
		{"EmptySlice", []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Shrink(tt.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shrink() = %v, want %v", got, tt.want)
			}
			if cap(got) != len(got) {
				t.Errorf("Shrink() capacity = %d, want %d", cap(got), len(got))
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
			if got := slices.Unique(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	shuffled := slices.Clone(original)
	slices.Shuffle(shuffled)

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
	slices.ShuffleN(shuffled, n)

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
		slices.ShuffleN([]int{1, 2, 3}, -1)
	})

	t.Run("PanicOnLargeN", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ShuffleN() did not panic when n > len(slice)")
			}
		}()
		slices.ShuffleN([]int{1, 2, 3}, 4)
	})
}

func TestClone(t *testing.T) {
	tests := []struct {
		name string
		s    []int
	}{
		{"NonEmptySlice", []int{1, 2, 3, 4, 5}},
		{"EmptySlice", []int{}},
		{"NilSlice", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloned := slices.Clone(tt.s)
			if !reflect.DeepEqual(cloned, tt.s) {
				t.Errorf("Clone() = %v, want %v", cloned, tt.s)
			}
			if len(tt.s) > 0 && &cloned[0] == &tt.s[0] {
				t.Errorf("Clone() returned a slice sharing memory with the original")
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name string
		n    int
		x    int
		want []int
	}{
		{"RepeatPositive", 5, 3, []int{3, 3, 3, 3, 3}},
		{"RepeatZero", 0, 3, []int{}},
		{"RepeatOne", 1, 3, []int{3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Repeat(tt.n, tt.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repeat() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("PanicOnNegativeCount", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Repeat() did not panic on negative count")
			}
		}()
		slices.Repeat(-1, 3)
	})
}

// Helper functions for testing
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
