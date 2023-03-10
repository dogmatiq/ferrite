package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Required is the application-facing interface for a value that is sourced
// from required environment variables.
//
// It is obtained by calling Deprecated() on a variable builder.
type Required[T any] interface {
	// Value returns the parsed and validated value of the environment variable.
	//
	// It panics if the environment variable is undefined or invalid.
	Value() T
}

// RequiredOption is an option that can be applied to a required variable.
type RequiredOption interface {
	variable.RegisterOption
}

// required is a convenience function that registers and returns a required[T]
// that maps one-to-one with an environment variable of the same type.
func required[T any, S variable.TypedSchema[T]](
	schema S,
	spec *variable.TypedSpecBuilder[T],
	options []RequiredOption,
) Required[T] {
	spec.MarkRequired()

	v := variable.Register(
		spec.Done(schema),
		options...,
	)

	return requiredFunc[T]{
		func() (T, error) {
			n, ok, err := v.NativeValue()
			if ok || err != nil {
				return n, err
			}
			return n, undefinedError(v)
		},
	}
}

// requiredFunc is an implementation of Required[T] that obtains the value
// from an arbitrary function.
type requiredFunc[T any] struct {
	fn func() (T, error)
}

func (d requiredFunc[T]) Value() T {
	n, err := d.fn()
	if err != nil {
		panic(err.Error())
	}
	return n
}
