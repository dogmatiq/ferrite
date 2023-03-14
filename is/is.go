package is

import "github.com/dogmatiq/ferrite/maybe"

// Predicate is a function that returns true if a value matches some condition.
type Predicate[T any] func(maybe.Value[T]) bool

// Equal returns a predicate that asserts that a value is equal to v.
func Equal[T comparable](v T) Predicate[T] {
	return func(m maybe.Value[T]) bool {
		x, ok := m.Get()
		return ok && x == v
	}
}

// Not returns a predicate that asserts the logical inverse of p.
func Not[T any](p Predicate[T]) Predicate[T] {
	return func(v maybe.Value[T]) bool {
		return !p(v)
	}
}
