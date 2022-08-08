package spec

import (
	"sync"
)

// Resolver is an interface for a resolver of any type.
type Resolver interface {
	Spec() Spec
	Resolve() (Value, error)
}

// ResolverOf parses and validates an environment variable of type T.
type ResolverOf[T any] struct {
	spec Spec

	once  sync.Once
	impl  func() (ValueOf[T], error)
	value ValueOf[T]
	err   error
}

// NewResolver returns a new resolver for the given specification.
func NewResolver[T any](
	spec Spec,
	resolve func() (ValueOf[T], error),
) *ResolverOf[T] {
	return &ResolverOf[T]{
		spec: spec,
		impl: resolve,
	}
}

// Spec returns the specification for the variable resolved by r.
func (r *ResolverOf[T]) Spec() Spec {
	return r.spec
}

// Resolve parses and validates the environment value.
func (r *ResolverOf[T]) Resolve() (Value, error) {
	return r.ResolveTyped()
}

// ResolveTyped parses and validates the environment value.
func (r *ResolverOf[T]) ResolveTyped() (ValueOf[T], error) {
	r.once.Do(func() {
		r.value, r.err = r.impl()
		r.impl = nil
	})

	return r.value, r.err
}

// Value is the resolved value of an environment variable.
type Value interface {
	// String returns the normalized string representation of the actual
	// environment variable value.
	String() string

	// IsDefault returns true if this value was derived from a default value,
	// rather than from the environment itself.
	IsDefault() bool
}

// ValueOf is the resolved value of an environment variable that is represented
// by type T.
type ValueOf[T any] struct {
	// Go is the Go in-memory representation of the value.
	Go T

	// Env is the normalized string representation of the actual environment
	// variable value.
	Env string

	// IsDef is true if this value was derived from a default value, rather
	// than from the environment itself.
	IsDef bool
}

// String returns the normalized string representation of the actual
// environment variable value.
func (v ValueOf[T]) String() string {
	return v.Env
}

// IsDefault returns true if this value was derived from a default value,
// rather than from the environment itself.
func (v ValueOf[T]) IsDefault() bool {
	return v.IsDef
}
