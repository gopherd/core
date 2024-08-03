package slices_test

import (
	"math"
	"testing"

	"github.com/gopherd/core/container/slices"
)

func TestSumFunc(t *testing.T) {
	var sum = slices.SumFunc([]int{1, 2, 3, 4}, func(x int) float64 {
		return math.Sqrt(float64(x))
	})
	t.Logf("sum: %v", sum)
}
