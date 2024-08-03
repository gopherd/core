package ordered_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gopherd/core/container/ordered"
	"github.com/gopherd/core/container/tree"
	"github.com/gopherd/core/operator"
)

func ExampleSet() {
	s := ordered.NewSet[int]()
	fmt.Print("empty:\n" + s.Stringify(nil))

	s.Insert(1)
	s.Insert(2)
	s.Insert(4)

	iter, ok := s.Insert(8)
	if !ok {
		fmt.Println("insert fail")
	} else {
		fmt.Println("inserted:", iter.Key())
	}

	fmt.Print("default:\n" + s.Stringify(nil))
	fmt.Print("custom:\n" + s.Stringify(&tree.Options{
		Prefix:     "... ",
		Parent:     "|   ",
		Branch:     "|-- ",
		LastBranch: "+-- ",
	}))
	fmt.Println("plain:\n" + s.String())

	// Output:
	// empty:
	// <nil>
	// inserted: 8
	// default:
	// 2
	// ├── 1
	// └── 4
	//     └── 8
	// custom:
	// ... 2
	// ... |-- 1
	// ... +-- 4
	// ...     +-- 8
	// plain:
	// [1 2 4 8]
}

func TestSet(t *testing.T) {
	s := ordered.NewSetFunc[int](operator.Greater[int])
	hashset := make(map[int]bool)

	rand.Seed(100)

	makeKey := func(i int) int {
		return i
	}

	add := func(k int) {
		_, found := hashset[k]
		hashset[k] = true
		_, ok := s.Insert(k)
		if ok != !found {
			t.Fatalf("set.Set: returned value want %v, but got %v", !found, found)
		}
	}

	remove := func(k int) {
		_, found := hashset[k]
		delete(hashset, k)
		ok := s.Remove(k)
		if ok != found {
			t.Fatalf("set.Remove: want %v, but got %v", found, ok)
		}
	}

	const (
		n    = 100
		keys = 30
	)
	for i := 0; i < n; i++ {
		for j := 0; j < keys/2; j++ {
			add(makeKey(j))
			add(makeKey(keys - 1 - j))
		}
		checkSet("add", t, s, hashset)
	}
	for j := 0; j < keys; j++ {
		key := makeKey(j)
		remove(key)
	}
	checkSet("remove", t, s, hashset)

	for i := 0; i < n; i++ {
		k := makeKey(rand.Intn(keys))
		var op string
		if rand.Intn(2) == 0 {
			op = "add"
			add(k)
		} else {
			op = "remove"
			remove(k)
		}
		checkSet(op, t, s, hashset)
	}
}

func checkSet[K comparable](op string, t *testing.T, s *ordered.Set[K], hashset map[K]bool) {
	if s.Len() != len(hashset) {
		t.Fatalf("[%s] len mismacthed: want %d, got %d", op, len(hashset), s.Len())
	}
	for k := range hashset {
		if !s.Contains(k) {
			t.Fatalf("[%s] key %v not found", op, k)
		}
	}
}
