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

func requiredOne[T any](
	v *variable.OfType[T],
) Required[T] {
	return required[T]{
		vars: []variable.Any{v},
		fn: func() (T, error) {
			n, ok, err := v.NativeValue()
			if ok || err != nil {
				return n, err
			}
			return n, undefinedError(v)
		},
	}
}

func optionalOne[T any](
	v *variable.OfType[T],
) Optional[T] {
	return optional[T]{
		vars: []variable.Any{v},
		fn: func() (T, bool, error) {
			return v.NativeValue()
		},
	}
}

func requiredMany[T any](
	fn func() (T, error),
	vars ...variable.Any,
) Required[T] {
	return required[T]{fn, vars}
}

func optionalMany[T any](
	fn func() (T, bool, error),
	vars ...variable.Any,
) Optional[T] {
	return optional[T]{fn, vars}
}

// required is an implementation of Required[T] that obtains the value
// from an arbitrary function.
type required[T any] struct {
	fn   func() (T, error)
	vars []variable.Any
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
	fn   func() (T, bool, error)
	vars []variable.Any
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
