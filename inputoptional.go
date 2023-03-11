package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Optional is a specialization of the Input interface for values obtained
// from deprecated environment variables.
type Optional[T any] interface {
	Input

	// Value returns the parsed and validated value built from the environment
	// variable(s).
	//
	// If the constituent environment variable(s) are not defined and there is
	// no default value, ok is false; otherwise, ok is true and v is the value.
	//
	// It panics if any of one of the constituent environment variable(s) has an
	// invalid value.
	Value() (T, bool)
}

// OptionalOption is an option that can be applied to an optional variable.
type OptionalOption interface {
	variable.RegisterOption
}

// optional is a convenience function that registers and returns a required[T]
// that maps one-to-one with an environment variable of the same type.
func optional[T any, S variable.TypedSchema[T]](
	schema S,
	spec *variable.TypedSpecBuilder[T],
	options []OptionalOption,
) Optional[T] {
	v := variable.Register(
		spec.Done(schema),
		options...,
	)

	return optionalFunc[T]{
		[]variable.Any{v},
		func() (T, bool, error) {
			return v.NativeValue()
		},
	}
}

// optionalFunc is an implementation of Optional[T] that obtains the value from
// an arbitrary function.
type optionalFunc[T any] struct {
	vars []variable.Any
	fn   func() (T, bool, error)
}

func (i optionalFunc[T]) Value() (T, bool) {
	n, ok, err := i.fn()
	if err != nil {
		panic(err.Error())
	}
	return n, ok
}

func (i optionalFunc[T]) variables() []variable.Any {
	return i.vars
}
