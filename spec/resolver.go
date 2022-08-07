package spec

import (
	"sync"
)

// Resolver parses and validates an environment variable.
type Resolver[T any] struct {
	spec Spec

	once  sync.Once
	impl  func() (Value[T], error)
	value Value[T]
	err   error
}

// NewResolver returns a new resolver for the given specification.
func NewResolver[T any](
	spec Spec,
	resolve func() (Value[T], error),
) *Resolver[T] {
	return &Resolver[T]{
		spec: spec,
		impl: resolve,
	}
}

func (r *Resolver[T]) Spec() Spec {
	return r.spec
}

func (r *Resolver[T]) Resolve() (Value[T], error) {
	r.once.Do(func() {
		r.value, r.err = r.impl()
		r.impl = nil
	})

	return r.value, r.err
}

// Value is the resolved value of an environment variable.
type Value[T any] struct {
	Go        T
	Env       string
	IsDefault bool
}

func (v Value[T]) String() string {
	return v.Env
}

func (v Value[T]) isDefault() bool {
	return v.IsDefault
}
