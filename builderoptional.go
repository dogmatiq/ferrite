package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Optional is the application-facing interface for a value that is sourced
// from optional environment variables.
//
// It is obtained by calling Deprecated() on a variable builder.
type Optional[T any] interface {
	// Value returns the parsed and validated value of the environment variable,
	// if it is defined.
	//
	// If the environment variable is not defined (and there is no default
	// value), ok is false; otherwise, ok is true and v is the value.
	//
	// It panics if the environment variable is defined but invalid.
	Value() (T, bool)
}

// optional is a convenience function that registers and returns a required[T]
// that maps one-to-one with an environment variable of the same type.
func optional[T any, S variable.TypedSchema[T]](
	schema S,
	spec *variable.TypedSpecBuilder[T],
	options []VariableOption,
) Optional[T] {
	v := variable.Register(
		spec.Done(schema),
		options...,
	)

	return optionalFunc[T]{
		func() (T, bool, error) {
			return v.NativeValue()
		},
	}
}

// optionalFunc is an implementation of Optional[T] that obtains the value from an
// arbitrary function.
type optionalFunc[T any] struct {
	fn func() (T, bool, error)
}

func (d optionalFunc[T]) Value() (T, bool) {
	n, ok, err := d.fn()
	if err != nil {
		panic(err.Error())
	}
	return n, ok
}
