// Note: This implementation is not concurrent safe.
// If multiple goroutines access it simultaneously, external synchronization should be used.

package collections

import "iter"

// GenericSet represents a set of unique values of type T.
//
// The identity of an element is defined by a key function provided by the caller,
// which returns a comparable key for each element.
type Set[T any] struct {
	elements map[UniqueKey[T]]T
	keyFn    func(T) UniqueKey[T]
}

// NewSet creates a new set for types that can provide a unique key themselves.
// The type T must implement the UniqueKeyer[T] interface.
func NewSet[T UniqueKeyer[T]](items ...T) Set[T] {
	return NewSetFunc(T.UniqueKey, items...)
}

// NewSetFunc creates a new set where the keyFn function determines
// how a unique key is computed for each element.
//
// The key function must return values that are comparable using ==
// and must guarantee a unique key for each distinct element.
func NewSetFunc[T any](keyFn func(T) UniqueKey[T], items ...T) Set[T] {
	s := Set[T]{
		elements: make(map[UniqueKey[T]]T),
		keyFn:    keyFn,
	}
	for _, v := range items {
		s.Add(v)
	}
	return s
}

// NewSetCmp creates a new set for any comparable types.
// The equality of elements is determined by the built-in == operator.
func NewSetCmp[T comparable](items ...T) Set[T] {
	return NewSetFunc(comparableKeyFunc[T], items...)
}

// Has returns true if the specified element is contained in the set,
// otherwise false.
func (s Set[T]) Has(v T) bool {
	if len(s.elements) == 0 {
		// Uninitialized sets are treated like empty sets.
		return false
	}
	k := s.keyFn(v)
	_, exists := s.elements[k]
	return exists
}

// Add adds one or more elements to the set.
//
// If an element with the same key already exists, it is replaced.
// Duplicate entries within the input are handled such that only
// the last value remains in the set.
func (s Set[T]) Add(values ...T) {
	for _, v := range values {
		k := s.keyFn(v)
		s.elements[k] = v
	}
}

// Remove removes an element from the set if it exists.
func (s Set[T]) Remove(v T) {
	k := s.keyFn(v)
	delete(s.elements, k)
}

// All returns an iterator over all elements in the set in no particular order.
//
// Example:
//
//	for v := range s.All() {
//	    fmt.Println(v)
//	}
//
// Modifications to the set during iteration result in undefined behavior.
func (s Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s.elements {
			if !yield(v) {
				return
			}
		}
	}
}

// Len returns the number of unique elements in the set.
func (s Set[T]) Len() int {
	return len(s.elements)
}
