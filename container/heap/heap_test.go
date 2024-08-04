package heap_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/gopherd/core/container/heap"
)

// intHeap is a min-heap of integers.
type intHeap []int

func (h intHeap) Len() int           { return len(h) }
func (h intHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x int)        { *h = append(*h, x) }
func (h *intHeap) Pop() int {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func TestHeapInit(t *testing.T) {
	h := &intHeap{3, 2, 1, 5, 6, 4}
	heap.Init(h)

	for i := 1; i < h.Len(); i++ {
		if h.Less(i, (i-1)/2) {
			t.Errorf("heap invariant violated: h[%d] = %d < h[%d] = %d", i, (*h)[i], (i-1)/2, (*h)[(i-1)/2])
		}
	}
}

func TestHeapPush(t *testing.T) {
	h := &intHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)

	if (*h)[0] != 1 {
		t.Errorf("expected minimum element 1, got %d", (*h)[0])
	}

	if h.Len() != 4 {
		t.Errorf("expected length 4, got %d", h.Len())
	}
}

func TestHeapPop(t *testing.T) {
	h := &intHeap{1, 2, 3, 4, 5}
	heap.Init(h)

	for i := 1; i <= 5; i++ {
		x := heap.Pop(h)
		if x != i {
			t.Errorf("pop got %d, want %d", x, i)
		}
	}

	if h.Len() != 0 {
		t.Errorf("expected heap to be empty, got length %d", h.Len())
	}
}

func TestHeapRemove(t *testing.T) {
	h := &intHeap{1, 2, 3, 4, 5}
	heap.Init(h)

	x := heap.Remove(h, 2)
	if x != 3 {
		t.Errorf("Remove(2) got %d, want 3", x)
	}

	if h.Len() != 4 {
		t.Errorf("expected length 4, got %d", h.Len())
	}

	// Verify heap invariant
	if !verifyHeap(h) {
		t.Errorf("heap invariant violated after Remove")
	}

	// Verify all original elements except the removed one are still in the heap
	remaining := []int{1, 2, 4, 5}
	for _, v := range remaining {
		found := false
		for _, hv := range *h {
			if hv == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find %d in heap after Remove, but it's missing", v)
		}
	}
}

// verifyHeap checks if the heap invariant is maintained
func verifyHeap(h *intHeap) bool {
	for i := 1; i < h.Len(); i++ {
		if h.Less(i, (i-1)/2) {
			return false
		}
	}
	return true
}

func TestHeapFix(t *testing.T) {
	h := &intHeap{1, 2, 3, 4, 5}
	heap.Init(h)

	(*h)[0] = 6
	heap.Fix(h, 0)

	expected := []int{2, 4, 3, 6, 5}
	for i, v := range *h {
		if v != expected[i] {
			t.Errorf("at index %d, got %d, want %d", i, v, expected[i])
		}
	}
}

func TestHeapFixUp(t *testing.T) {
	h := &intHeap{1, 3, 2, 4, 5}
	heap.Init(h)

	// Change a leaf node to a smaller value
	(*h)[4] = 0
	heap.Fix(h, 4)

	if !verifyHeap(h) {
		t.Errorf("heap invariant violated after Fix (up)")
	}

	// The smallest value should now be at the root
	if (*h)[0] != 0 {
		t.Errorf("expected root to be 0, got %d", (*h)[0])
	}

	// Verify the heap structure
	expected := []int{0, 1, 2, 4, 3}
	for i, v := range *h {
		if v != expected[i] {
			t.Errorf("at index %d, got %d, want %d", i, v, expected[i])
		}
	}
}

func TestHeapIntegration(t *testing.T) {
	h := &intHeap{}

	// Push elements
	for i := 20; i > 0; i-- {
		heap.Push(h, i)
	}

	// Verify heap property
	for i := 1; i < h.Len(); i++ {
		if h.Less(i, (i-1)/2) {
			t.Errorf("heap invariant violated: h[%d] = %d < h[%d] = %d", i, (*h)[i], (i-1)/2, (*h)[(i-1)/2])
		}
	}

	// Pop all elements
	for i := 1; h.Len() > 0; i++ {
		x := heap.Pop(h)
		if x != i {
			t.Errorf("pop got %d, want %d", x, i)
		}
	}
}

func TestHeapWithRandomData(t *testing.T) {
	h := &intHeap{}
	data := rand.Perm(1000)

	for _, v := range data {
		heap.Push(h, v)
	}

	sort.Ints(data)

	for i, want := range data {
		got := heap.Pop(h)
		if got != want {
			t.Errorf("pop %d got %d, want %d", i, got, want)
		}
	}
}

func BenchmarkHeapPush(b *testing.B) {
	h := &intHeap{}
	for i := 0; i < b.N; i++ {
		heap.Push(h, i)
	}
}

func BenchmarkHeapPop(b *testing.B) {
	h := &intHeap{}
	for i := 0; i < b.N; i++ {
		heap.Push(h, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop(h)
	}
}
