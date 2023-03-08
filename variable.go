package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/variable"
)

// Required is the application-facing interface for obtaining a value from an
// environment variable (or variables) when those variables are required to be
// defined.
type Required[T any] interface {
	// Value returns the parsed and validated value of the environment variable.
	//
	// It panics if the environment variable is undefined or invalid.
	Value() T
}

// Optional is the application-facing interface for obtaining a value from an
// environment variable (or variables) when those variables are optional.
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

// undefinedError returns an error that indicates that a variable is undefined.
func undefinedError(v variable.Any) error {
	return fmt.Errorf(
		"%s is undefined and does not have a default value",
		v.Spec().Name(),
	)
}

// req is a convenience function that registers and returns a required[T] that
// maps one-to-one with an environment variable of the same type.
func req[T any, S variable.TypedSchema[T]](
	schema S,
	spec variable.SpecBuilder[T],
	options []Option,
) Required[T] {
	spec.MarkRequired()

	opts := resolveOptions(options)

	v := variable.Register(
		spec.Done(schema),
		opts.RegisterOptions,
	)

	return required[T]{
		[]variable.Any{v},
		func() (T, error) {
			n, ok, err := v.NativeValue()
			if ok || err != nil {
				return n, err
			}
			return n, undefinedError(v)
		},
	}
}

// opt is a convenience function that registers and returns a required[T] that
// maps one-to-one with an environment variable of the same type.
func opt[T any, S variable.TypedSchema[T]](
	schema S,
	spec variable.SpecBuilder[T],
	options []Option,
) Optional[T] {
	opts := resolveOptions(options)

	v := variable.Register(
		spec.Done(schema),
		opts.RegisterOptions,
	)

	return optional[T]{
		[]variable.Any{v},
		func() (T, bool, error) {
			return v.NativeValue()
		},
	}
}

// required is an implementation of Required[T] that obtains the value
// from an arbitrary function.
type required[T any] struct {
	vars []variable.Any
	fn   func() (T, error)
}

func (d required[T]) Value() T {
	n, err := d.fn()
	if err != nil {
		panic(err.Error())
	}
	return n
}

func (d required[T]) variables() []variable.Any {
	return d.vars
}

// optional is an implementation of Optional[T] that obtains the value from an
// arbitrary function.
type optional[T any] struct {
	vars []variable.Any
	fn   func() (T, bool, error)
}

func (d optional[T]) Value() (T, bool) {
	n, ok, err := d.fn()
	if err != nil {
		panic(err.Error())
	}
	return n, ok
}

func (d optional[T]) variables() []variable.Any {
	return d.vars
}
