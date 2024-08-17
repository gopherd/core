package constraints_test

import (
	"fmt"

	"github.com/gopherd/core/constraints"
)

func sum[T constraints.Addable](values ...T) T {
	var result T
	for _, v := range values {
		result += v
	}
	return result
}

func ExampleAddable() {
	fmt.Println(sum(1, 2))
	fmt.Println(sum("hello", " ", "world"))
	// Output:
	// 3
	// hello world
}
