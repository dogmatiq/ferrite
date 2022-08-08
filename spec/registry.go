package spec

import (
	"fmt"
	"sync"
)

var registry sync.Map // map[string]Resolver

// Register adds a resolver to the global registry.
func Register(r Resolver) {
	if _, loaded := registry.LoadOrStore(r.Spec().Name, r); loaded {
		panic(fmt.Sprintf("%s has multiple specifications", r.Spec().Name))
	}
}

// ResetRegistry resets the global registry.
func ResetRegistry() {
	registry.Range(func(key, _ any) bool {
		registry.Delete(key)
		return true
	})
}

// RangeRegistry calls fn for each resolver in the registry.
//
// If fn returns false it stops the iteration.
func RangeRegistry(fn func(Resolver) bool) {
	registry.Range(func(_, r any) bool {
		return fn(r.(Resolver))
	})
}
