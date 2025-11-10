package collections

import (
	"testing"
)

func TestGenericSet(t *testing.T) {
	s := NewSet[testKey]()

	if s.Len() != 0 {
		t.Errorf("Expects empty set, reveived length %d", s.Len())
	}

	s.Add(testKey("a"))
	if s.Len() != 1 {
		t.Errorf("Expects length 1 after adding 'a', received %d", s.Len())
	}

	s.Add(testKey("a"))
	s.Add(testKey("b"))
	if s.Len() != 2 {
		t.Errorf("Expects length 2 after adding 'a' and 'b', received %d", s.Len())
	}

	s.Remove(testKey("a"))
	if s.Len() != 1 {
		t.Errorf("Expects length 1 after removing 'a', received %d", s.Len())
	}

	if s.Has(testKey("a")) {
		t.Errorf("'a' should not be present after removal")
	}
	if !s.Has(testKey("b")) {
		t.Errorf("'b' should be present in the set")
	}
}

func TestEmptySet(t *testing.T) {
	var unset Set[string]

	if unset.Len() != 0 {
		t.Errorf("An uninitialized set should have length 0")
	}

	if unset.Has("xyz") {
		t.Errorf("An uninitialized set should not contain any elements")
	}
}
