package collections

// testKey is a string type that can serve as its own unique key
// since strings are inherently comparable in Go.
type testKey string

// Ensure that testKey implements both UniqueKey and UniqueKeyer.
var (
	_ UniqueKey[testKey]   = testKey("")
	_ UniqueKeyer[testKey] = testKey("")
)

// IsUniqueKey is a no-op marker method that satisfies the UniqueKey interface.
func (k testKey) IsUniqueKey(testKey) {}

// UniqueKey implements the UniqueKeyer interface and returns the key itself.
func (k testKey) UniqueKey() UniqueKey[testKey] {
	return UniqueKey[testKey](k)
}
