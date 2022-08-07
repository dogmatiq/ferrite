package spec

import (
	"fmt"
	"sync"
)

var registry sync.Map // map[string]*Resolver

// Register adds a resolver to the global registry.
func Register[T any](r *Resolver[T]) {
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
