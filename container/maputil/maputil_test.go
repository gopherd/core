package maputil_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/gopherd/core/container/maputil"
	"github.com/gopherd/core/container/pair"
)

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want []string
	}{
		{
			name: "empty map",
			m:    map[string]int{},
			want: []string{},
		},
		{
			name: "non-empty map",
			m:    map[string]int{"a": 1, "b": 2, "c": 3},
			want: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maputil.Keys(tt.m)
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want []int
	}{
		{
			name: "empty map",
			m:    map[string]int{},
			want: []int{},
		},
		{
			name: "non-empty map",
			m:    map[string]int{"a": 1, "b": 2, "c": 3},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maputil.Values(tt.m)
			sort.Ints(got)
			sort.Ints(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	want := []string{"a:1", "b:2", "c:3"}

	got := maputil.Map(m, func(k string, v int) string {
		return fmt.Sprintf("%s:%d", k, v)
	})

	sort.Strings(got)
	sort.Strings(want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Map() = %v, want %v", got, want)
	}
}

func TestMinKey(t *testing.T) {
	tests := []struct {
		name string
		m    map[int]string
		want int
	}{
		{
			name: "empty map",
			m:    map[int]string{},
			want: 0,
		},
		{
			name: "single element",
			m:    map[int]string{5: "five"},
			want: 5,
		},
		{
			name: "multiple elements",
			m:    map[int]string{3: "three", 1: "one", 4: "four"},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maputil.MinKey(tt.m); got != tt.want {
				t.Errorf("MinKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxKey(t *testing.T) {
	tests := []struct {
		name string
		m    map[int]string
		want int
	}{
		{
			name: "empty map",
			m:    map[int]string{},
			want: 0,
		},
		{
			name: "single element",
			m:    map[int]string{5: "five"},
			want: 5,
		},
		{
			name: "multiple elements",
			m:    map[int]string{3: "three", 1: "one", 4: "four"},
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maputil.MaxKey(tt.m); got != tt.want {
				t.Errorf("MaxKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinMaxKey(t *testing.T) {
	tests := []struct {
		name    string
		m       map[int]string
		wantMin int
		wantMax int
	}{
		{
			name:    "empty map",
			m:       map[int]string{},
			wantMin: 0,
			wantMax: 0,
		},
		{
			name:    "single element",
			m:       map[int]string{5: "five"},
			wantMin: 5,
			wantMax: 5,
		},
		{
			name:    "multiple elements",
			m:       map[int]string{3: "three", 1: "one", 4: "four"},
			wantMin: 1,
			wantMax: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := maputil.MinMaxKey(tt.m)
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("MinMaxKey() = (%v, %v), want (%v, %v)", gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestMinValue(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want pair.Pair[string, int]
	}{
		{
			name: "empty map",
			m:    map[string]int{},
			want: pair.Pair[string, int]{},
		},
		{
			name: "single element",
			m:    map[string]int{"five": 5},
			want: pair.New("five", 5),
		},
		{
			name: "multiple elements",
			m:    map[string]int{"three": 3, "one": 1, "four": 4},
			want: pair.New("one", 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maputil.MinValue(tt.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MinValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxValue(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want pair.Pair[string, int]
	}{
		{
			name: "empty map",
			m:    map[string]int{},
			want: pair.Pair[string, int]{},
		},
		{
			name: "single element",
			m:    map[string]int{"five": 5},
			want: pair.New("five", 5),
		},
		{
			name: "multiple elements",
			m:    map[string]int{"three": 3, "one": 1, "four": 4},
			want: pair.New("four", 4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maputil.MaxValue(tt.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaxValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinMaxValue(t *testing.T) {
	tests := []struct {
		name    string
		m       map[string]int
		wantMin pair.Pair[string, int]
		wantMax pair.Pair[string, int]
	}{
		{
			name:    "empty map",
			m:       map[string]int{},
			wantMin: pair.Pair[string, int]{},
			wantMax: pair.Pair[string, int]{},
		},
		{
			name:    "single element",
			m:       map[string]int{"five": 5},
			wantMin: pair.New("five", 5),
			wantMax: pair.New("five", 5),
		},
		{
			name:    "multiple elements",
			m:       map[string]int{"three": 3, "one": 1, "four": 4},
			wantMin: pair.New("one", 1),
			wantMax: pair.New("four", 4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := maputil.MinMaxValue(tt.m)
			if !reflect.DeepEqual(gotMin, tt.wantMin) || !reflect.DeepEqual(gotMax, tt.wantMax) {
				t.Errorf("MinMaxValue() = (%v, %v), want (%v, %v)", gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestMinKeyFunc(t *testing.T) {
	m := map[string]int{"a": 1, "bb": 2, "ccc": 3}
	want := 1

	got := maputil.MinKeyFunc(m, func(k string, v int) int {
		return len(k)
	})

	if got != want {
		t.Errorf("MinKeyFunc() = %v, want %v", got, want)
	}
}

func TestMaxKeyFunc(t *testing.T) {
	m := map[string]int{"a": 1, "bb": 2, "ccc": 3}
	want := 3

	got := maputil.MaxKeyFunc(m, func(k string, v int) int {
		return len(k)
	})

	if got != want {
		t.Errorf("MaxKeyFunc() = %v, want %v", got, want)
	}
}

func TestMinMaxKeyFunc(t *testing.T) {
	m := map[string]int{"a": 1, "bb": 2, "ccc": 3}
	wantMin, wantMax := 1, 3

	gotMin, gotMax := maputil.MinMaxKeyFunc(m, func(k string, v int) int {
		return len(k)
	})

	if gotMin != wantMin || gotMax != wantMax {
		t.Errorf("MinMaxKeyFunc() = (%v, %v), want (%v, %v)", gotMin, gotMax, wantMin, wantMax)
	}
}

func TestMinValueFunc(t *testing.T) {
	m := map[string]int{"a": 1, "bb": 2, "ccc": 3}
	want := pair.New("a", 1)

	got := maputil.MinValueFunc(m, func(k string, v int) int {
		return v
	})

	if !reflect.DeepEqual(got, want) {
		t.Errorf("MinValueFunc() = %v, want %v", got, want)
	}
}

func TestMaxValueFunc(t *testing.T) {
	m := map[string]int{"a": 1, "bb": 2, "ccc": 3}
	want := pair.New("ccc", 3)

	got := maputil.MaxValueFunc(m, func(k string, v int) int {
		return v
	})

	if !reflect.DeepEqual(got, want) {
		t.Errorf("MaxValueFunc() = %v, want %v", got, want)
	}
}

func TestMinMaxValueFunc(t *testing.T) {
	m := map[string]int{"a": 1, "bb": 2, "ccc": 3}
	wantMin := pair.New("a", 1)
	wantMax := pair.New("ccc", 3)

	gotMin, gotMax := maputil.MinMaxValueFunc(m, func(k string, v int) int {
		return v
	})

	if !reflect.DeepEqual(gotMin, wantMin) || !reflect.DeepEqual(gotMax, wantMax) {
		t.Errorf("MinMaxValueFunc() = (%v, %v), want (%v, %v)", gotMin, gotMax, wantMin, wantMax)
	}
}

func TestCopyFunc(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2, "c": 3}
	dst := make(map[int]string)
	want := map[int]string{1: "a", 2: "b", 3: "c"}

	got := maputil.CopyFunc(dst, src, func(k string, v int) (int, string) {
		return v, k
	})

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CopyFunc() = %v, want %v", got, want)
	}
}

func TestSumKey(t *testing.T) {
	tests := []struct {
		name string
		m    map[int]string
		want int
	}{
		{
			name: "empty map",
			m:    map[int]string{},
			want: 0,
		},
		{
			name: "single element",
			m:    map[int]string{5: "five"},
			want: 5,
		},
		{
			name: "multiple elements",
			m:    map[int]string{1: "one", 2: "two", 3: "three"},
			want: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maputil.SumKey(tt.m); got != tt.want {
				t.Errorf("SumKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumValue(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want int
	}{
		{
			name: "empty map",
			m:    map[string]int{},
			want: 0,
		},
		{
			name: "single element",
			m:    map[string]int{"five": 5},
			want: 5,
		},
		{
			name: "multiple elements",
			m:    map[string]int{"one": 1, "two": 2, "three": 3},
			want: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maputil.SumValue(tt.m); got != tt.want {
				t.Errorf("SumValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumFunc(t *testing.T) {
	m := map[string]int{"a": 1, "bb": 2, "ccc": 3}
	want := 6

	got := maputil.SumFunc(m, func(k string, v int) int {
		return len(k)
	})

	if got != want {
		t.Errorf("SumFunc() = %v, want %v", got, want)
	}
}

// TestEdgeCases tests some edge cases that weren't covered in the main tests
func TestEdgeCases(t *testing.T) {
	t.Run("MinMaxKey with negative numbers", func(t *testing.T) {
		m := map[int]string{-3: "minus three", 0: "zero", 5: "five"}
		gotMin, gotMax := maputil.MinMaxKey(m)
		if gotMin != -3 || gotMax != 5 {
			t.Errorf("MinMaxKey() = (%v, %v), want (-3, 5)", gotMin, gotMax)
		}
	})

	t.Run("SumKey with overflow", func(t *testing.T) {
		m := map[int8]string{127: "max", 1: "one"}
		sum := maputil.SumKey(m)
		if sum != -128 {
			t.Errorf("SumKey() = %v, want -128 (overflow)", sum)
		}
	})
}
