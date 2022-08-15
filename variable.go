package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/variable"
)

// Required is the application-facing representation of an environment variable
// that must have a value.
type Required[T any] struct {
	resolve func() (T, error)
}

// Value returns the parsed and validated value of the environment variable.
//
// It panics if the environment variable is undefined or invalid.
func (r Required[T]) Value() T {
	x, err := r.resolve()
	if err != nil {
		panic(err.Error())
	}

	return x
}

// Optional is the application-facing representation of an environment variable
// that may optionally be undefined.
type Optional[T any] struct {
	resolve func() (T, bool, error)
}

// Value returns the parsed and validated value of the environment variable,
// if it is defined.
//
// If the environment variable is not defined (and there is no default
// value), ok is false; otherwise, ok is true and v is the value.
//
// It panics if the environment variable is defined but invalid.
func (o Optional[T]) Value() (v T, ok bool) {
	x, ok, err := o.resolve()
	if err != nil {
		panic(err.Error())
	}

	return x, ok
}

func registerRequired[T any](
	spec variable.TypedSpec[T],
	options []variable.RegisterOption,
) Required[T] {
	v := variable.Register(spec, options)

	return Required[T]{
		func() (zero T, _ error) {
			n, ok, err := v.NativeValue()
			if ok || err != nil {
				return n, err
			}

			return zero, undefinedError(v)
		},
	}
}

func registerOptional[T any](
	spec variable.TypedSpec[T],
	options []variable.RegisterOption,
) Optional[T] {
	v := variable.Register(spec, options)

	return Optional[T]{
		func() (T, bool, error) {
			return v.NativeValue()
		},
	}
}

func undefinedError(v variable.Any) error {
	return fmt.Errorf(
		"%s is undefined and does not have a default value",
		v.Spec().Name(),
	)
}
