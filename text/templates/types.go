package templates

import (
	"errors"
	"fmt"
	"slices"
)

// @api(Container/Vector) represents a slice of any type. Function `list` creates a Vector from the given values.
//
// Example:
// ```tmpl
// {{list 1 2 3}}
// ```
//
// Output:
// ```
// [1 2 3]
// ```
type Vector []any

// @api(Container/Vector.Push) pushes an item to the end of the vector.
//
// Example:
// ```tmpl
// {{- $v := list 1 2 3}}
// {{$v.Push 4}}
// ```
//
// Output:
// ```
// [1 2 3 4]
// ```
func (v *Vector) Push(item any) Vector {
	*v = append(*v, item)
	return *v
}

// @api(Container/Vector.Pop) pops an item from the end of the vector.
// It returns an error if the vector is empty.
//
// Example:
// ```tmpl
// {{- $v := list 1 2 3}}
// {{- $v.Pop}}
// {{$v}}
// ```
//
// Output:
// ```
// 3
// [1 2]
// ```
func (v *Vector) Pop() (any, error) {
	if len(*v) == 0 {
		return nil, errors.New("empty vector")
	}
	item := (*v)[len(*v)-1]
	*v = (*v)[:len(*v)-1]
	return item, nil
}

// @api(Container/Vector.Len) returns the length of the vector.
//
// Example:
// ```tmpl
// {{- $v := list 1 2 3 4 5}}
// {{- $v.Len}}
// ```
//
// Output:
// ```
// 5
// ```
func (v Vector) Len() int {
	return len(v)
}

// @api(Container/Vector.IsEmpty) reports whether the vector is empty.
//
// Example:
// ```tmpl
// {{- $v := list 1 2}}
// {{- $v.IsEmpty}}
// {{- $v.Pop | _ }}
// {{- $v.IsEmpty}}
// {{- $v.Pop | _ }}
// {{- $v.IsEmpty}}
// ```
//
// Output:
// ```
// false
// false
// true
// ```
func (v Vector) IsEmpty() bool {
	return len(v) == 0
}

// @api(Container/Vector.Get) returns the item at the specified index.
func (v Vector) Get(index int) (any, error) {
	if index < 0 || index >= len(v) {
		return nil, fmt.Errorf("index %d out of range [0, %d)", index, len(v))
	}
	return v[index], nil
}

// @api(Container/Vector.Set) sets the item at the specified index.
// It returns an error if the index is out of range.
func (v Vector) Set(index int, item any) error {
	if index < 0 || index >= len(v) {
		return fmt.Errorf("index %d out of range [0, %d)", index, len(v))
	}
	v[index] = item
	return nil
}

// @api(Container/Vector.Insert) inserts an item at the specified index.
func (v *Vector) Insert(index int, item any) error {
	if index < 0 || index > len(*v) {
		return fmt.Errorf("index %d out of range [0, %d]", index, len(*v))
	}
	*v = slices.Insert(*v, index, item)
	return nil
}

// @api(Container/Vector.RemoveAt) removes the item at the specified index.
func (v *Vector) RemoveAt(index int) error {
	if index < 0 || index >= len(*v) {
		return fmt.Errorf("index %d out of range [0, %d)", index, len(*v))
	}
	*v = append((*v)[:index], (*v)[index+1:]...)
	return nil
}

// @api(Container/Vector.Splice) removes count items starting from the specified index and inserts new items.
func (v *Vector) Splice(start, count int, items ...any) error {
	if start < 0 || start >= len(*v) {
		return fmt.Errorf("start index %d out of range [0, %d)", start, len(*v))
	}
	if count < 0 || start+count > len(*v) {
		return fmt.Errorf("count %d out of range [0, %d)", count, len(*v)-start)
	}
	*v = append((*v)[:start], append(items, (*v)[start+count:]...)...)
	return nil
}

// @api(Container/Vector.Clear) removes all items from the vector.
func (v *Vector) Clear() {
	*v = nil
}

// @api(Container/Dictionary) represents a dictionary of any key-value pairs.
type Dictionary map[any]any

// @api(Container/Dictionary.Get) returns the value for the specified key.
func (d Dictionary) Get(key any) any {
	return d[key]
}

// @api(Container/Dictionary.Set) sets the value for the specified key.
func (d Dictionary) Set(key, value any) {
	d[key] = value
}

// @api(Container/Dictionary.Has) reports whether the dictionary has the specified key.
func (d Dictionary) Has(key any) bool {
	_, ok := d[key]
	return ok
}

// @api(Container/Dictionary.Remove) removes the specified key from the dictionary.
func (d Dictionary) Remove(key any) {
	delete(d, key)
}

// @api(Container/Dictionary.Clear) removes all key-value pairs from the dictionary.
func (d Dictionary) Clear() {
	for key := range d {
		delete(d, key)
	}
}
