package set

import (
	"fmt"
	"sync"
)

type Set struct {
	mutex sync.RWMutex
	m     map[interface{}]struct{}
}

// New creates a new Set and initializes its internal map optionally adding an initial elements to the set
func New(values ...interface{}) *Set {
	s := Set{m: make(map[interface{}]struct{})}
	for _, v := range values {
		s.Add(v)
	}
	return &s
}

// Add a new element into the set and return a flag to determine whether the element was newly added or already existed.
func (s *Set) Add(value interface{}) bool {
	existed := s.Has(value)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.m[value] = struct{}{}
	return !existed
}

// Determine whether the value was already existed in the set.
func (s *Set) Has(value interface{}) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if _, ok := s.m[value]; ok {
		return true
	}
	return false
}

// Clone copies the current set into a new identical set.
func (s *Set) Clone() *Set {
	set := New()
	for _, v := range s.Enumerate() {
		set.Add(v)
	}
	return set
}

// Enumerate a list of unordered slice of all elements from the set.
func (s *Set) Enumerate() []interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	values := make([]interface{}, 0)
	for k := range s.m {
		values = append(values, k)
	}
	return values
}

// Returns a set containing all elements present in the set.
func (s *Set) Difference(t *Set) *Set {
	diff := New()
	for _, e := range s.Enumerate() {
		duplicate := false
		for _, p := range t.Enumerate() {
			if e == p {
				duplicate = true
			}
		}
		if !duplicate {
			diff.Add(e)
		}
	}
	return diff
}

// Returns whether or not two sets have the same length and no differences.
func (s *Set) Equal(t *Set) bool {
	return s.Size() == t.Size() && s.Difference(t).Size() == 0
}

// Returns all elements which return true when the function is applied.
func (s *Set) Filter(fn func(interface{}) bool) *Set {
	t := New()
	for _, e := range s.Enumerate() {
		if fn(e) {
			t.Add(e)
		}
	}
	return t
}

// Returns all elements present in both the compared set and the given set.
func (s *Set) Intersect(t *Set) *Set {
	c := s.Clone()
	for _, e := range s.Difference(t).Enumerate() {
		c.Remove(e)
	}
	return c
}

// Map applies a function over all elements of the set.
func (s *Set) Map(fn func(interface{}) interface{}) *Set {
	t := New()
	for _, e := range s.Enumerate() {
		t.Add(fn(e))
	}
	return t
}

// Reduce applies a function over all elements of the set.
func (s *Set) Reduce(value interface{}, fn func(interface{}, interface{}) interface{}) interface{} {
	for _, e := range s.Enumerate() {
		value = fn(value, e)
	}
	return value
}

// Destroys an element in the set, returning if the element was destroyed.
func (s *Set) Remove(value interface{}) bool {
	existed := s.Has(value)
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	delete(s.m, value)
	return existed
}

// Returns the size or cardinality of the set.
func (s *Set) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.m)
}

// Returns a string representation of the set.
func (s *Set) String() string {
	str := "{ "
	if s.Size() == 0 {
		return str + "Ø }"
	}
	for k := range s.m {
		str += fmt.Sprintf("%v ", k)
	}
	return str + "}"
}

// Returns a set containing all elements present in the set.
func (s *Set) Union(t *Set) *Set {
	c := s.Clone()
	for _, e := range t.Enumerate() {
		c.Add(e)
	}
	return c
}
