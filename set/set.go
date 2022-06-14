package set

import (
	"encoding/json"
	"fmt"
)

type Set[T comparable] struct {
	data map[T]bool
}

func New[T comparable](items ...T) *Set[T] {
	data := buildData(items)
	return &Set[T]{
		data: data,
	}
}

func buildData[T comparable](items []T) map[T]bool {
	data := make(map[T]bool, len(items))
	for _, item := range items {
		data[item] = true
	}
	return data
}

func (s *Set[T]) ContainedIn(superSet *Set[T]) bool {
	for item := range s.data {
		if !superSet.Contains(item) {
			return false
		}
	}
	return true
}

func (s *Set[T]) Contains(item T) bool {
	_, exists := s.data[item]
	return exists
}

func (s *Set[T]) Delete(item T) {
	delete(s.data, item)
}

func (s *Set[T]) Iter() []T {
	slice := make([]T, 0, len(s.data))
	for item := range s.data {
		slice = append(slice, item)
	}
	return slice
}

func (s *Set[T]) Len() int {
	return len(s.data)
}

func (s *Set[T]) Put(item T) {
	s.data[item] = true
}

func (s *Set[T]) PutAll(items ...T) {
	for _, item := range items {
		s.data[item] = true
	}
}

func (s *Set[T]) String() string {
	return fmt.Sprintf("%v", s.Iter())
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Iter())
}

func (s *Set[T]) UnmarshalJSON(b []byte) error {
	items := &[]T{}
	if err := json.Unmarshal(b, items); err != nil {
		return err
	}
	s.data = buildData(*items)
	return nil
}

func Intersection[T comparable](a, b *Set[T]) *Set[T] {
	intersection := New[T]()
	for _, item := range a.Iter() {
		if b.Contains(item) {
			intersection.Put(item)
		}
	}
	return intersection
}

// subtract a from b
func Subtract[T comparable](a, b *Set[T]) *Set[T] {
	subtraction := New[T]()
	for _, item := range b.Iter() {
		if !a.Contains(item) {
			subtraction.Put(item)
		}
	}
	return subtraction
}

func Union[T comparable](a, b *Set[T]) *Set[T] {
	union := New(a.Iter()...)
	union.PutAll(b.Iter()...)
	return union
}
