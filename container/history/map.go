package history

import (
	"fmt"
	"maps"
	"strings"
)

type Map[K comparable, V any] struct {
	recorder Recorder
	data     map[K]V
}

func NewMap[K comparable, V any](recorder Recorder, size int) *Map[K, V] {
	return &Map[K, V]{
		recorder: recorder,
		data:     make(map[K]V, size),
	}
}

func (m Map[K, V]) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	i := 0
	for k, v := range m.data {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%v:%v", k, v)
		i++
	}
	sb.WriteByte('}')
	return sb.String()
}

func (m *Map[K, V]) Clone(recorder Recorder) *Map[K, V] {
	return &Map[K, V]{
		recorder: recorder,
		data:     maps.Clone(m.data),
	}
}

func (m *Map[K, V]) Len() int {
	return len(m.data)
}

func (m *Map[K, V]) Contains(k K) bool {
	_, ok := m.data[k]
	return ok
}

func (m *Map[K, V]) Get(k K) (v V, ok bool) {
	v, ok = m.data[k]
	return
}

func (m *Map[K, V]) Set(k K, v V) bool {
	old, replaced := m.data[k]
	if replaced {
		m.recorder.AddRecord(&mapUndoSet[K, V]{m: m, k: k, v: old, replaced: true})
	} else {
		m.recorder.AddRecord(&mapUndoSet[K, V]{m: m, k: k})
	}
	m.data[k] = v
	return replaced
}

type mapUndoSet[K comparable, V any] struct {
	m        *Map[K, V]
	k        K
	v        V
	replaced bool
}

func (r *mapUndoSet[K, V]) Undo() {
	if r.replaced {
		r.m.data[r.k] = r.v
	} else {
		delete(r.m.data, r.k)
	}
}

func (m *Map[K, V]) Delete(k K) (v V, deleted bool) {
	v, deleted = m.data[k]
	if deleted {
		m.recorder.AddRecord(&mapUndoDelete[K, V]{m: m, k: k, v: v})
	}
	delete(m.data, k)
	return
}

type mapUndoDelete[K comparable, V any] struct {
	m *Map[K, V]
	k K
	v V
}

func (r *mapUndoDelete[K, V]) Undo() {
	r.m.data[r.k] = r.v
}

func (m *Map[K, V]) Range(f func(K, V) bool) bool {
	for k, v := range m.data {
		if !f(k, v) {
			return false
		}
	}
	return true
}

func (m *Map[K, V]) Clear() {
	m.recorder.AddRecord(&mapUndoClear[K, V]{m: m, values: maps.Clone(m.data)})
	clear(m.data)
}

type mapUndoClear[K comparable, V any] struct {
	m      *Map[K, V]
	values map[K]V
}

func (r *mapUndoClear[K, V]) Undo() {
	r.m.data = r.values
}
