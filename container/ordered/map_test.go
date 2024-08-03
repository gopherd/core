package ordered_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/gopherd/core/container/ordered"
	"github.com/gopherd/core/container/tree"
	"github.com/gopherd/core/operator"
)

func ExampleMap() {
	m := ordered.NewMap[int, int]()
	fmt.Print("empty:\n" + m.Stringify(nil))

	m.Insert(1, 2)
	m.Insert(2, 4)
	m.Insert(4, 8)

	iter, ok := m.Insert(8, 16)
	if !ok {
		fmt.Println("insert fail")
	} else {
		fmt.Println("inserted:", iter.Key(), iter.Value())
	}

	fmt.Print("default:\n" + m.Stringify(nil))
	fmt.Print("custom:\n" + m.Stringify(&tree.Options{
		Prefix:     "... ",
		Parent:     "|   ",
		Branch:     "|-- ",
		LastBranch: "+-- ",
	}))
	fmt.Println("plain:\n" + m.String())

	// Output:
	// empty:
	// <nil>
	// inserted: 8 16
	// default:
	// 2:4
	// ├── 1:2
	// └── 4:8
	//     └── 8:16
	// custom:
	// ... 2:4
	// ... |-- 1:2
	// ... +-- 4:8
	// ...     +-- 8:16
	// plain:
	// [1:2 2:4 4:8 8:16]
}

func TestMap(t *testing.T) {
	m := ordered.NewMapFunc[int, string](operator.Greater[int])
	hashmap := make(map[int]string)

	rand.Seed(100)

	makeKey := func(i int) int {
		return i
	}
	makeValue := func(i int) string {
		return strconv.Itoa(i)
	}

	add := func(k int, v string) {
		_, found := hashmap[k]
		hashmap[k] = v
		_, ok := m.Insert(k, v)
		if ok != !found {
			t.Fatalf("map.Set: returned value want %v, but got %v", !found, found)
		}
	}

	remove := func(k int) {
		_, found := hashmap[k]
		delete(hashmap, k)
		ok := m.Remove(k)
		if ok != found {
			t.Fatalf("map.Remove: want %v, but got %v", found, ok)
		}
	}

	const (
		n    = 100
		keys = 30
	)
	for i := 0; i < n; i++ {
		for j := 0; j < keys/2; j++ {
			key := makeKey(j)
			value := makeValue(j * (i + 1))
			add(key, value)
			key = makeKey(keys - 1 - j)
			value = makeValue((keys - 1 - j) * (i + 1))
			add(key, value)
		}
		checkMap("add", t, m, hashmap)
	}
	for j := 0; j < keys; j++ {
		key := makeKey(j)
		remove(key)
	}
	checkMap("remove", t, m, hashmap)

	for i := 0; i < n; i++ {
		k := makeKey(rand.Intn(keys))
		var op string
		if rand.Intn(2) == 0 {
			op = "add"
			v := makeValue(rand.Intn(99999999) * 3)
			add(k, v)
		} else {
			op = "remove"
			remove(k)
		}
		checkMap(op, t, m, hashmap)
	}
}

func checkMap[K comparable, V comparable](op string, t *testing.T, m *ordered.Map[K, V], hashmap map[K]V) {
	if m.Len() != len(hashmap) {
		t.Fatalf("[%s] len mismacthed: want %d, got %d", op, len(hashmap), m.Len())
	}
	for k, v := range hashmap {
		node := m.Find(k)
		if node == nil {
			t.Fatalf("[%s] key %v not found", op, k)
		}
		if node.Value() != v {
			t.Fatalf("[%s] value mismacthed for key %v: want %v, got %v", op, k, v, node.Value())
		}
	}
}
