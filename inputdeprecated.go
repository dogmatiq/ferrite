package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Deprecated is a specialization of the Input interface for values obtained
// from deprecated environment variables.
type Deprecated[T any] interface {
	Input

	// DeprecatedValue returns the parsed and validated value built from the
	// environment variable(s).
	//
	// If the constituent environment variable(s) are not defined and there is
	// no default value, ok is false; otherwise, ok is true and v is the value.
	//
	// It panics if any of one of the constituent environment variable(s) has an
	// invalid value.
	DeprecatedValue() (T, bool)
}

// DeprecatedOption is an option that can be applied to a deprecated variable.
type DeprecatedOption interface {
	variable.RegisterOption
}

// deprecated is a convenience function that registers and returns a
// deprecated[T] that maps one-to-one with an environment variable of the same
// type.
func deprecated[T any, S variable.TypedSchema[T]](
	schema S,
	spec *variable.TypedSpecBuilder[T],
	options []DeprecatedOption,
) Deprecated[T] {
	spec.MarkDeprecated()

	v := variable.Register(
		spec.Done(schema),
		options...,
	)

	// interface is currently empty so we don't need an implementation
	return deprecatedFunc[T]{
		[]variable.Any{v},
		func() (T, bool, error) {
			return v.NativeValue()
		},
	}
}

// deprecatedFunc is an implementation of Deprecated[T] that obtains the value
// from an arbitrary function.
type deprecatedFunc[T any] struct {
	vars []variable.Any
	fn   func() (T, bool, error)
}

func (i deprecatedFunc[T]) DeprecatedValue() (T, bool) {
	n, ok, err := i.fn()
	if err != nil {
		panic(err.Error())
	}
	return n, ok
}

func (i deprecatedFunc[T]) variables() []variable.Any {
	return i.vars
}
