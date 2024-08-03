package tuple

// Tuple holds n values
type Tuple[T any] interface {
	Len() int
	At(i int) T
}

// Equal reports whether t1 == t2
func Equal[T comparable](t1, t2 Tuple[T]) bool {
	var n = t1.Len()
	if n != t2.Len() {
		return false
	}
	for i := 0; i < n; i++ {
		if t1.At(i) != t2.At(i) {
			return false
		}
	}
	return true
}

// Make create a tuple by parameters
func Make[T any](x ...T) Tuple[T] {
	switch len(x) {
	case 0:
		return Empty[T]{}
	case 1:
		return tuple1[T]{x[0]}
	case 2:
		return tuple2[T]{x[0], x[1]}
	case 3:
		return tuple3[T]{x[0], x[1], x[2]}
	case 4:
		return tuple4[T]{x[0], x[1], x[2], x[3]}
	case 5:
		return tuple5[T]{x[0], x[1], x[2], x[3], x[4]}
	case 6:
		return tuple6[T]{x[0], x[1], x[2], x[3], x[4], x[5]}
	case 7:
		return tuple7[T]{x[0], x[1], x[2], x[3], x[4], x[5], x[6]}
	case 8:
		return tuple8[T]{x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7]}
	case 9:
		return tuple9[T]{x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7], x[8]}
	default:
		return tupleN[T](x)
	}
}

func T1[T any](a T) Tuple[T]                         { return tuple1[T]{a} }
func T2[T any](a, b T) Tuple[T]                      { return tuple2[T]{a, b} }
func T3[T any](a, b, c T) Tuple[T]                   { return tuple3[T]{a, b, c} }
func T4[T any](a, b, c, d T) Tuple[T]                { return tuple4[T]{a, b, c, d} }
func T5[T any](a, b, c, d, e T) Tuple[T]             { return tuple5[T]{a, b, c, d, e} }
func T6[T any](a, b, c, d, e, f T) Tuple[T]          { return tuple6[T]{a, b, c, d, e, f} }
func T7[T any](a, b, c, d, e, f, g T) Tuple[T]       { return tuple7[T]{a, b, c, d, e, f, g} }
func T8[T any](a, b, c, d, e, f, g, h T) Tuple[T]    { return tuple8[T]{a, b, c, d, e, f, g, h} }
func T9[T any](a, b, c, d, e, f, g, h, i T) Tuple[T] { return tuple9[T]{a, b, c, d, e, f, g, h, i} }

// Concat concats two tuple to one
func Concat[T any](a, b Tuple[T]) Tuple[T] {
	if a.Len() == 0 {
		return b
	}
	if b.Len() == 0 {
		return a
	}
	return combination[T]{a, b}
}

type combination[T any] struct {
	a, b Tuple[T]
}

func (t combination[T]) Len() int {
	return t.a.Len() + t.b.Len()
}

func (t combination[T]) At(i int) T {
	var m = t.a.Len()
	if i < m {
		return t.a.At(i)
	}
	return t.b.At(i - m)
}

// Slice slice the tuple by range [i, j)
func Slice[T any](t Tuple[T], i, j int) Tuple[T] {
	return sliced[T]{
		t: t,
		i: i,
		j: j,
	}
}

type sliced[T any] struct {
	t    Tuple[T]
	i, j int
}

func (t sliced[T]) Len() int   { return t.j - t.i }
func (t sliced[T]) At(i int) T { return t.t.At(i + t.i) }

//----------------------------------------
// special tuples

// Empty tuple
type Empty[T any] struct{}

func (Empty[T]) Len() int   { return 0 }
func (Empty[T]) At(i int) T { panic("out of range") }

type tuple1[T any] [1]T
type tuple2[T any] [2]T
type tuple3[T any] [3]T
type tuple4[T any] [4]T
type tuple5[T any] [5]T
type tuple6[T any] [6]T
type tuple7[T any] [7]T
type tuple8[T any] [8]T
type tuple9[T any] [9]T
type tupleN[T any] []T

func (tuple1[T]) Len() int   { return 1 }
func (tuple2[T]) Len() int   { return 2 }
func (tuple3[T]) Len() int   { return 3 }
func (tuple4[T]) Len() int   { return 4 }
func (tuple5[T]) Len() int   { return 5 }
func (tuple6[T]) Len() int   { return 6 }
func (tuple7[T]) Len() int   { return 7 }
func (tuple8[T]) Len() int   { return 8 }
func (tuple9[T]) Len() int   { return 9 }
func (t tupleN[T]) Len() int { return len(t) }

func (t tuple1[T]) At(i int) T { return t[i] }
func (t tuple2[T]) At(i int) T { return t[i] }
func (t tuple3[T]) At(i int) T { return t[i] }
func (t tuple4[T]) At(i int) T { return t[i] }
func (t tuple5[T]) At(i int) T { return t[i] }
func (t tuple6[T]) At(i int) T { return t[i] }
func (t tuple7[T]) At(i int) T { return t[i] }
func (t tuple8[T]) At(i int) T { return t[i] }
func (t tuple9[T]) At(i int) T { return t[i] }
func (t tupleN[T]) At(i int) T { return t[i] }
