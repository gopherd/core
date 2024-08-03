package history

import (
	"fmt"
	"maps"
	"strings"
)

type Set[K comparable] struct {
	recorder Recorder
	data     map[K]struct{}
}

func NewSet[K comparable](recorder Recorder, size int) *Set[K] {
	return &Set[K]{
		recorder: recorder,
		data:     make(map[K]struct{}, size),
	}
}

func (s Set[K]) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	i := 0
	for k := range s.data {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, k)
		i++
	}
	sb.WriteByte('}')
	return sb.String()
}

func (s *Set[K]) Clone(recorder Recorder) *Set[K] {
	return &Set[K]{
		recorder: recorder,
		data:     maps.Clone(s.data),
	}
}

func (s *Set[K]) Len() int {
	return len(s.data)
}

func (s *Set[K]) Contains(k K) bool {
	_, ok := s.data[k]
	return ok
}

func (s *Set[K]) Insert(k K) (inserted bool) {
	_, found := s.data[k]
	if found {
		return false
	}
	s.data[k] = struct{}{}
	s.recorder.AddRecord(&setUndoInsert[K]{s: s, k: k})
	return
}

type setUndoInsert[K comparable] struct {
	s *Set[K]
	k K
}

func (r *setUndoInsert[K]) Undo() {
	delete(r.s.data, r.k)
}

func (s *Set[K]) Delete(k K) (deleted bool) {
	_, deleted = s.data[k]
	if deleted {
		s.recorder.AddRecord(&setUndoDelete[K]{s: s, k: k})
	}
	delete(s.data, k)
	return
}

type setUndoDelete[K comparable] struct {
	s *Set[K]
	k K
}

func (r *setUndoDelete[K]) Undo() {
	r.s.data[r.k] = struct{}{}
}

func (s *Set[K]) Range(f func(K) bool) bool {
	for k := range s.data {
		if !f(k) {
			return false
		}
	}
	return true
}

func (s *Set[K]) Clear() {
	s.recorder.AddRecord(&setUndoClear[K]{s: s, values: maps.Clone(s.data)})
	clear(s.data)
}

type setUndoClear[K comparable] struct {
	s      *Set[K]
	values map[K]struct{}
}

func (r *setUndoClear[K]) Undo() {
	r.s.data = r.values
}
