// This module defines the types and interfaces for unique keys,
// which can be used in generic sets or maps.

package collections

// UniqueKey describes a comparable value that uniquely identifies another element
// of type T.
//
// The Go type checker cannot enforce that this type is actually comparable;
// if a non-comparable type is used, it will cause a runtime panic.
//
// This interface serves as a marker to type-safely define unique keys for
// collections (e.g., sets or maps).
type UniqueKey[T any] interface {
	// IsUniqueKey is never called; it only serves to mark the type
	// as a unique key for T.
	IsUniqueKey(T)
}

// UniqueKeyer is a type that can produce its own unique key.
//
// A correct implementation must return a different key for each distinct object of T.
// What "unique" means depends on the semantics of the type.
type UniqueKeyer[T any] interface {
	// UniqueKey returns the unique key of the object.
	UniqueKey() UniqueKey[T]
}

// comparableKey is a helper type that allows any comparable type
// to be used as a unique key in sets or maps.
//
// Internally, the value is stored in a single-element array to ensure
// that the type remains comparable with ==.
type comparableKey[T comparable] [1]T

// IsUniqueKey fulfills the UniqueKey interface without providing functionality.
func (comparableKey[T]) IsUniqueKey(T) {}

// comparableKeyFunc creates a UniqueKey instance from a comparable value.
func comparableKeyFunc[T comparable](v T) UniqueKey[T] {
	return comparableKey[T]{v}
}
