# Go 1.23 Iterators: A Deep Dive

## Introduction

Go 1.23 introduces a game-changing feature to the language: built-in support for iterators. This addition brings a more functional and flexible approach to working with sequences of data, aligning Go with modern programming paradigms while maintaining its trademark simplicity and efficiency. In this comprehensive guide, we'll explore the new `iter` package, diving deep into its core concepts, practical applications, and some of the more advanced functions it offers.

## Understanding Iterators in Go 1.23

At its core, an iterator in Go 1.23 is a function that yields successive elements of a sequence to a callback function. The `iter` package defines two main types of iterators:

1. `Seq[V any]`: For sequences of single values
2. `Seq2[K, V any]`: For sequences of key-value pairs

These iterators are defined as function types:

```go
type Seq[V any] func(yield func(V) bool)
type Seq2[K, V any] func(yield func(K, V) bool)
```

The `yield` function is called for each element in the sequence. It returns a boolean indicating whether the iteration should continue (true) or stop (false).

## Basic Usage

Let's start with a simple example to illustrate how to use these iterators:

```go
func PrintAll[V any](seq iter.Seq[V]) {
    for v := range seq {
        fmt.Println(v)
    }
}
```

This function takes a `Seq[V]` and prints all its elements. The `range` keyword is overloaded in Go 1.23 to work with iterators, making their usage intuitive and familiar.

## Creating Iterators

The [github.com/gopherd/core/container/iters]() package provides several utility functions to create and manipulate iterators. Let's look at some of them:

### Enumerate

The `Enumerate` function creates an iterator from a slice:

```go
func Enumerate[S ~[]E, E any](s S) iter.Seq[E] {
    return func(yield func(E) bool) {
        for _, v := range s {
            if !yield(v) {
                return
            }
        }
    }
}
```

Usage:

```go
for v := range iters.Enumerate([]string{"a", "b", "c"}) {
    fmt.Println(v) // Output: a \n b \n c
}
```

### Range

The `Range` function generates a sequence of numbers:

```go
func Range[T cmp.Ordered](start, end, step T) iter.Seq[T] {
    // ... (error checking omitted for brevity)
    return func(yield func(T) bool) {
        if start < end {
            for i := start; i < end; i += step {
                if !yield(i) {
                    return
                }
            }
        } else {
            for i := start; i > end; i += step {
                if !yield(i) {
                    return
                }
            }
        }
    }
}
```

Usage:

```go
for v := range iters.Range(1, 10, 2) {
    fmt.Println(v) // Output: 1 3 5 7 9
}
```

## Advanced Iterator Functions

Now, let's dive into some of the more complex functions provided by the `iters` package.

### Zip

The `Zip` function combines two sequences into a single sequence of pairs:

```go
func Zip[T any, U any](s1 iter.Seq[T], s2 iter.Seq[U]) iter.Seq2[T, U] {
    return func(yield func(T, U) bool) {
        next, stop := iter.Pull(s2)
        defer stop()
        for v1 := range s1 {
            v2, _ := next()
            if !yield(v1, v2) {
                return
            }
        }
        var zero1 T
        for {
            if v2, ok := next(); !ok {
                return
            } else if !yield(zero1, v2) {
                return
            }
        }
    }
}
```

This function is particularly interesting because it demonstrates how to work with two iterators simultaneously. It uses the `Pull` function to convert `s2` into a pull-style iterator, allowing for more control over the iteration process.

The `Zip` function continues until both input sequences are exhausted. If one sequence is longer than the other, the remaining elements are paired with zero values of the other type.

Usage:

```go
seq1 := iters.Enumerate([]int{1, 2, 3})
seq2 := iters.Enumerate([]string{"a", "b", "c", "d"})
for v1, v2 := range iters.Zip(seq1, seq2) {
    fmt.Printf("(%d, %s)\n", v1, v2)
}
// Output:
// (1, a)
// (2, b)
// (3, c)
// (0, d)
```

### Split

The `Split` function divides a slice into multiple chunks:

```go
func Split[S ~[]T, T any](s S, n int) iter.Seq[[]T] {
    if n < 1 {
        panic("n must be positive")
    }
    return func(yield func([]T) bool) {
        total := len(s)
        size := total / n
        remainder := total % n
        i := 0
        for i < total {
            var chunk []T
            if remainder > 0 {
                chunk = s[i : i+size+1]
                remainder--
                i += size + 1
            } else {
                chunk = s[i : i+size]
                i += size
            }
            if !yield(chunk) {
                return
            }
        }
    }
}
```

This function is useful when you need to process a large slice in smaller, more manageable pieces. It ensures that the chunks are as evenly sized as possible, distributing any remainder elements among the first few chunks.

Usage:

```go
data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
for chunk := range iters.Split(data, 3) {
    fmt.Println(chunk)
}
// Output:
// [1 2 3 4]
// [5 6 7]
// [8 9 10]
```

### GroupBy

The `GroupBy` function is a powerful tool for organizing data based on a key function:

```go
func GroupBy[K comparable, V any](s iter.Seq[V], f func(V) K) iter.Seq2[K, []V] {
    return func(yield func(K, []V) bool) {
        groups := make(map[K][]V)
        for v := range s {
            k := f(v)
            groups[k] = append(groups[k], v)
        }
        for k, vs := range groups {
            if !yield(k, vs) {
                return
            }
        }
    }
}
```

This function takes a sequence and a key function. It groups the elements of the sequence based on the keys produced by the key function. The result is a sequence of key-value pairs, where each key is associated with a slice of all elements that produced that key.

Usage:

```go
data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
evenOdd := func(n int) string {
    if n%2 == 0 {
        return "even"
    }
    return "odd"
}
for key, group := range iters.GroupBy(iters.Enumerate(data), evenOdd) {
    fmt.Printf("%s: %v\n", key, group)
}
// Output:
// odd: [1 3 5 7 9]
// even: [2 4 6 8 10]
```

## The Power of Pull

The `iter` package introduces a powerful concept called "Pull", which converts a push-style iterator into a pull-style iterator. This is particularly useful when you need more control over the iteration process.

```go
func Pull[V any](seq Seq[V]) (next func() (V, bool), stop func())
```

The `Pull` function returns two functions:
1. `next`: Returns the next value in the sequence and a boolean indicating if the value is valid.
2. `stop`: Ends the iteration.

This allows for more complex iteration patterns, as demonstrated in the `Zip` function we saw earlier.

## Conclusion

Go 1.23's introduction of iterators marks a significant evolution in the language's capabilities for handling sequences of data. The `iter` package provides a rich set of tools for creating, manipulating, and consuming iterators, enabling more expressive and functional programming styles while maintaining Go's simplicity and performance.

From basic operations like `Enumerate` and `Range`, to more complex functions like `Zip`, `Split`, and `GroupBy`, the new iterator system offers powerful abstractions for working with data sequences. The introduction of pull-style iterators through the `Pull` function further extends the flexibility and control available to developers.

As you incorporate these new features into your Go programs, you'll likely find new, more elegant solutions to common programming problems. The iterator system in Go 1.23 opens up exciting possibilities for data processing, functional programming, and beyond.

Remember, while iterators provide powerful abstractions, they should be used judiciously. In many cases, simple loops or slices may still be the most readable and performant solution. As with all features, the key is to understand the tools available and choose the right one for each specific task.

Happy iterating!
