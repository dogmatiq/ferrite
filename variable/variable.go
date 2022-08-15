package variable

import (
	"sync"

	"github.com/dogmatiq/ferrite/maybe"
)

// Any is an interface for an environment variable of any type.
type Any interface {
	// Spec returns the variable's specification.
	Spec() Spec

	// IsValid returns true if the variable is valid.
	IsValid() bool

	// Value returns the variable's value.
	Value() (Value, bool, ValueError)
}

// OfType is an environment variable depicted by type T.
type OfType[T any] struct {
	spec TypedSpec[T]
	env  Environment

	once  sync.Once
	value maybe.Value[valueOf[T]]
	err   ValueError
}

// Spec returns the variable's specification.
func (v *OfType[T]) Spec() Spec {
	return v.spec
}

// IsValid returns true if the variable is valid.
func (v *OfType[T]) IsValid() bool {
	v.resolve()

	if v.err != nil {
		return false
	}

	if v.spec.required && v.value.IsEmpty() {
		return false
	}

	return true
}

// Value returns the variable's value.
func (v *OfType[T]) Value() (Value, bool, ValueError) {
	v.resolve()
	x, ok := v.value.Get()
	return x, ok, v.err
}

// NativeValue returns the variable's native value.
func (v *OfType[T]) NativeValue() (T, bool, ValueError) {
	v.resolve()
	x, ok := v.value.Get()
	return x.native, ok, v.err
}

func (v *OfType[T]) resolve() {
	v.once.Do(func() {
		lit := v.env.Get(v.spec.name)

		if lit.String == "" {
			v.value = v.spec.def
			return
		}

		n, c, err := v.spec.Unmarshal(lit)
		if err != nil {
			v.err = valueError{
				name:    v.spec.name,
				literal: lit,
				cause:   err,
			}
			return
		}

		v.value = maybe.Some(valueOf[T]{
			verbatim:  lit,
			native:    n,
			canonical: c,
		})
	})
}
