package random_test

import (
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"testing"

	"github.com/gopherd/core/math/random"
)

type mockRand struct {
}

func (r *mockRand) Uint64N(n uint64) uint64 {
	return rand.Uint64() % n
}

func TestShuffle(t *testing.T) {
	r := new(mockRand)
	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	shuffled := slices.Clone(original)
	random.Shuffle(r, shuffled)

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
	r := new(mockRand)
	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	shuffled := slices.Clone(original)
	n := 5
	random.ShuffleN(r, shuffled, n)

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

	// Check that elements after n are not guaranteed to be in their original positions
	// (Again, there's a small chance this could fail even with correct implementation)
	allSame := true
	for i := n; i < len(original); i++ {
		if original[i] != shuffled[i] {
			allSame = false
			break
		}
	}
	if !allSame {
		t.Errorf("ShuffleN() did not affect any elements after position %d", n)
	}

	t.Run("PanicOnNegativeN", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ShuffleN() did not panic on negative n")
			}
		}()
		random.ShuffleN(r, []int{1, 2, 3}, -1)
	})

	t.Run("PanicOnLargeN", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ShuffleN() did not panic when n > len(slice)")
			}
		}()
		random.ShuffleN(r, []int{1, 2, 3}, 4)
	})
}
