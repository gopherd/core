//go:build go1.23

package iters_test

import (
	"fmt"
	"iter"
	"maps"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gopherd/core/container/iters"
	"github.com/gopherd/core/container/pair"
	"github.com/gopherd/core/op"
)

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
			seq1: iters.List(1, 2, 3),
			seq2: iters.List("a", "b", "c"),
			expected: []struct {
				v1 int
				v2 string
			}{{1, "a"}, {2, "b"}, {3, "c"}},
		},
		{
			name: "seq1 longer than seq2",
			seq1: iters.List(1, 2, 3, 4),
			seq2: iters.List("a", "b"),
			expected: []struct {
				v1 int
				v2 string
			}{{1, "a"}, {2, "b"}, {3, ""}, {4, ""}},
		},
		{
			name: "seq2 longer than seq1",
			seq1: iters.List(1, 2),
			seq2: iters.List("a", "b", "c", "d"),
			expected: []struct {
				v1 int
				v2 string
			}{{1, "a"}, {2, "b"}, {0, "c"}, {0, "d"}},
		},
		{
			name: "seq1 is empty",
			seq1: iters.List[int](),
			seq2: iters.List("a", "b", "c"),
			expected: []struct {
				v1 int
				v2 string
			}{{0, "a"}, {0, "b"}, {0, "c"}},
		},
		{
			name: "seq2 is empty",
			seq1: iters.List[int](1, 2, 3),
			seq2: iters.List[string](),
			expected: []struct {
				v1 int
				v2 string
			}{{1, ""}, {2, ""}, {3, ""}},
		},
		{
			name: "empty sequences",
			seq1: iters.List[int](),
			seq2: iters.List[string](),
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

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		expected []int
	}{
		{"no duplicates", iters.List(1, 2, 3, 4), []int{1, 2, 3, 4}},
		{"with duplicates", iters.List(1, 2, 2, 3, 3, 3, 4), []int{1, 2, 3, 4}},
		{"all duplicates", iters.List(1, 1, 1, 1), []int{1}},
		{"empty sequence", iters.List[int](), []int{}},
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
			input:    iters.List(1, 2, 3),
			f:        func(i int) string { return strconv.Itoa(i) },
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "empty sequence",
			input:    iters.List[int](),
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

func TestMap2(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq2[string, int]
		f        func(string, int) string
		expected []string
	}{
		{
			name:     "combine key and value",
			input:    maps.All(map[string]int{"a": 1, "b": 2, "c": 3}),
			f:        func(k string, v int) string { return fmt.Sprintf("%s:%d", k, v) },
			expected: []string{"a:1", "b:2", "c:3"},
		},
		{
			name:     "empty map",
			input:    maps.All(map[string]int{}),
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
			input:    iters.List(1, 2, 3, 4, 5, 6),
			f:        func(i int) bool { return i%2 == 0 },
			expected: []int{2, 4, 6},
		},
		{
			name:     "no matches",
			input:    iters.List(1, 3, 5),
			f:        func(i int) bool { return i%2 == 0 },
			expected: []int{},
		},
		{
			name:     "all match",
			input:    iters.List(2, 4, 6),
			f:        func(i int) bool { return i%2 == 0 },
			expected: []int{2, 4, 6},
		},
		{
			name:     "empty sequence",
			input:    iters.List[int](),
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

func TestFilter2(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq2[string, int]
		f        func(string, int) bool
		expected map[string]int
	}{
		{
			name:     "values greater than 2",
			input:    maps.All(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}),
			f:        func(k string, v int) bool { return v > 2 },
			expected: map[string]int{"c": 3, "d": 4},
		},
		{
			name:     "no matches",
			input:    maps.All(map[string]int{"a": 1, "b": 2}),
			f:        func(k string, v int) bool { return v > 5 },
			expected: map[string]int{},
		},
		{
			name:     "all match",
			input:    maps.All(map[string]int{"a": 1, "b": 2}),
			f:        func(k string, v int) bool { return v > 0 },
			expected: map[string]int{"a": 1, "b": 2},
		},
		{
			name:     "empty map",
			input:    maps.All(map[string]int{}),
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

func TestGroupBy(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		f        func(int) string
		expected map[string][]int
	}{
		{
			name:  "group by even/odd",
			input: iters.List(1, 2, 3, 4, 5, 6),
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
			input:    iters.List[int](),
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
		{"non-empty sequence", iters.List(3, 1, 4, 1, 5, 9, 2, 6, 5, 3), []int{1, 1, 2, 3, 3, 4, 5, 5, 6, 9}},
		{"empty sequence", iters.List[int](), []int{}},
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

func TestSort2(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []string
	}{
		{
			name:     "non-empty map",
			input:    map[string]int{"c": 3, "a": 1, "b": 2},
			expected: []string{"a:1", "b:2", "c:3"},
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]string, 0)
			for k, v := range iters.Sort2(maps.All(tt.input)) {
				result = append(result, fmt.Sprintf("%s:%d", k, v))
			}

			sort.Strings(result)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Sort2() = %v, want %v", result, tt.expected)
			}

			for range iters.Sort2(maps.All(tt.input)) {
				break
			}
		})
	}
}

func TestSortFunc2(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		compare  func(pair.Pair[string, int], pair.Pair[string, int]) int
		expected []string
	}{
		{
			name: "compare first",
			input: map[string]int{
				"c": 1,
				"b": 2,
				"a": 3,
			},
			compare: pair.CompareFirst[string, int],
			expected: []string{
				"a:3",
				"b:2",
				"c:1",
			},
		},
		{
			name: "compare second by reverse order",
			input: map[string]int{
				"c": 1,
				"b": 2,
				"a": 3,
			},
			compare: op.ReverseCompare(pair.CompareSecond[string, int]),
			expected: []string{
				"a:3",
				"b:2",
				"c:1",
			},
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			compare:  pair.Compare[string, int],
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]string, 0)
			for k, v := range iters.SortFunc2(maps.All(tt.input), tt.compare) {
				result = append(result, fmt.Sprintf("%s:%d", k, v))
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortFunc2() = %v, want %v", result, tt.expected)
			}

			for range iters.SortFunc2(maps.All(tt.input), tt.compare) {
				break
			}
		})
	}
}

func TestSortKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []string
	}{
		{
			name:     "non-empty map",
			input:    map[string]int{"c": 3, "a": 1, "b": 2},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]string, 0)
			for k := range iters.SortKeys(maps.All(tt.input)) {
				result = append(result, k)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortKeys() = %v, want %v", result, tt.expected)
			}

			for range iters.SortKeys(maps.All(tt.input)) {
				break
			}
		})
	}
}

func TestSortValues(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []int
	}{
		{
			name:     "non-empty map",
			input:    map[string]int{"c": 1, "a": 3, "b": 2},
			expected: []int{1, 2, 3},
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]int, 0)
			for _, v := range iters.SortValues(maps.All(tt.input)) {
				result = append(result, v)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortValues() = %v, want %v", result, tt.expected)
			}

			for range iters.SortValues(maps.All(tt.input)) {
				break
			}
		})
	}
}

// Example tests
func ExampleEnumerateMap() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range maps.All(m) {
		fmt.Printf("%s: %d\n", k, v)
	}
	// Unordered output:
	// a: 1
	// b: 2
	// c: 3
}

func ExampleZip() {
	s1 := iters.List(1, 2, 3)
	s2 := iters.List("a", "b", "c")
	for v1, v2 := range iters.Zip(s1, s2) {
		fmt.Printf("%d: %s\n", v1, v2)
	}
	// Output:
	// 1: a
	// 2: b
	// 3: c
}

func ExampleUnique() {
	s := iters.List(1, 2, 2, 3, 3, 3, 4)
	for v := range iters.Unique(s) {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4
}

func ExampleFilter() {
	s := iters.List(1, 2, 3, 4, 5, 6)
	evenNumbers := iters.Filter(s, func(i int) bool { return i%2 == 0 })
	for v := range evenNumbers {
		fmt.Printf("%d ", v)
	}
	// Output: 2 4 6
}

func ExampleGroupBy() {
	s := iters.List(1, 2, 3, 4, 5, 6)
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

func TestWithIndex(t *testing.T) {
	tests := []struct {
		name     string
		seq      []int
		expected []int
	}{
		{
			name:     "non-empty sequence",
			seq:      []int{1, 2, 3, 4, 5},
			expected: []int{0, 1, 1, 2, 2, 3, 3, 4, 4, 5},
		},
		{
			name:     "empty sequence",
			seq:      []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]int, 0)
			for i, v := range iters.WithIndex(iters.List(tt.seq...)) {
				result = append(result, i)
				result = append(result, v)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("WithIndex() = %v, want %v", result, tt.expected)
			}

			for range iters.WithIndex(iters.List(tt.seq...)) {
				break
			}
		})
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
			for k := range iters.Keys(maps.All(tt.seq)) {
				got = append(got, k)
			}

			// Sort both slices for comparison
			sort.Strings(got)
			sort.Strings(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}

			for range iters.Keys(maps.All(tt.seq)) {
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
			for v := range iters.Values(maps.All(tt.seq)) {
				got = append(got, v)
			}

			// Sort both slices for comparison
			sort.Ints(got)
			sort.Ints(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}

			for range iters.Values(maps.All(tt.seq)) {
				break
			}
		})
	}
}

func ExampleKeys() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := make([]string, 0, 3)
	for k := range iters.Keys(maps.All(m)) {
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
	for v := range iters.Values(maps.All(m)) {
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
	for range iters.Keys(maps.All(largeMap)) {
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
	for range iters.Values(maps.All(largeMap)) {
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
			uniqueSeq := iters.UniqueFunc(iters.List(tt.seq...), tt.eq)
			got := make([]int, 0)
			for v := range uniqueSeq {
				got = append(got, v)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueFunc() = %v, want %v", got, tt.want)
			}

			for range iters.UniqueFunc(iters.List(tt.seq...), tt.eq) {
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

	uniqueSeq := iters.UniqueFunc(iters.List(words...), eq)
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

	uniqueSeq := iters.UniqueFunc(iters.List(numbers...), eq)
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

	uniqueSeq := iters.UniqueFunc(iters.List(people...), eq)
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

func TestConcat(t *testing.T) {
	tests := []struct {
		name     string
		input    []iter.Seq[int]
		expected []int
	}{
		{
			name:     "non-empty sequences",
			input:    []iter.Seq[int]{iters.List(1, 2), iters.List(3, 4), iters.List(5, 6)},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "empty sequences",
			input:    []iter.Seq[int]{iters.List[int](), iters.List[int](), iters.List[int]()},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]int, 0)
			for v := range iters.Concat(tt.input...) {
				result = append(result, v)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Concat() = %v, want %v", result, tt.expected)
			}

			for range iters.Concat(tt.input...) {
				break
			}
		})
	}
}

func TestConcat2(t *testing.T) {
	tests := []struct {
		name     string
		input    []iter.Seq2[string, int]
		expected map[string]int
	}{
		{
			name: "non-empty sequences",
			input: []iter.Seq2[string, int]{
				maps.All(map[string]int{"a": 1, "b": 2}),
				maps.All(map[string]int{"c": 3, "d": 4}),
				maps.All(map[string]int{"e": 5, "f": 6}),
			},
			expected: map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6},
		},
		{
			name:     "empty sequences",
			input:    []iter.Seq2[string, int]{maps.All(map[string]int{}), maps.All(map[string]int{}), maps.All(map[string]int{})},
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make(map[string]int)
			for k, v := range iters.Concat2(tt.input...) {
				result[k] = v
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Concat2() = %v, want %v", result, tt.expected)
			}

			for range iters.Concat2(tt.input...) {
				break
			}
		})
	}
}

func TestDistinct(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		expected []int
	}{
		{
			name:     "non-empty sequences",
			input:    iters.List(1, 2, 3, 4, 3, 4, 5, 6, 5, 6, 7, 8),
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:     "empty sequences",
			input:    iters.List[int](),
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]int, 0)
			for v := range iters.Distinct(tt.input) {
				result = append(result, v)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Distinct() = %v, want %v", result, tt.expected)
			}

			for range iters.Distinct(tt.input) {
				break
			}
		})
	}
}

func TestDistinctFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    iter.Seq[int]
		key      func(int) int
		expected []int
	}{
		{
			name:     "StandardKey",
			input:    iters.List(1, 2, 3, 4, 3, 4, 5, 6, 5, 6, 7, 8),
			key:      func(a int) int { return a },
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:     "EmptySequence",
			input:    iters.List[int](),
			key:      func(a int) int { return a },
			expected: []int{},
		},
		{
			name:     "AllUnique",
			input:    iters.List(1, 2, 3, 4, 5),
			key:      func(a int) int { return a },
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "AllSame",
			input:    iters.List(1, 2, 3, 4, 5),
			key:      func(a int) int { return 0 },
			expected: []int{0},
		},
		{
			name:     "ModKey",
			input:    iters.List(1, 11, 2, 22, 3, 33),
			key:      func(a int) int { return a % 10 },
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distinctSeq := iters.DistinctFunc(tt.input, tt.key)
			got := make([]int, 0)
			for v := range distinctSeq {
				got = append(got, v)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("DistinctFunc() = %v, want %v", got, tt.expected)
			}

			for range iters.DistinctFunc(tt.input, tt.key) {
				break
			}
		})
	}
}
