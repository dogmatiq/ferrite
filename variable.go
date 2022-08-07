package ferrite

import (
	"errors"

	"github.com/dogmatiq/ferrite/spec"
)

// Optional is the application-facing representation of an environment variable
// that may optionally be undefined.
//
// It is the standard implementation of ferrite.Optional.
type Optional[T any] struct {
	resolve func() (T, error)
}

// Value returns the parsed value of the environment value, if it is defined.
//
// It panics if the environment variable is defined but invalid.
func (v Optional[T]) Value() (T, bool) {
	value, err := v.resolve()
	if err != nil {
		var undef UndefinedError
		if errors.As(err, &undef) {
			// Only treat "normal" undefined errors as ok for optional
			// variables. If there was some other causal error that's a problem.
			if undef.Cause == nil {
				return value, false
			}
		}

		panic(err.Error())
	}

	return value, true
}

// Required is the application-facing representation of an environment
// variable that must always have a valid value.
//
// It is the standard implementation of ferrite.Required.
type Required[T any] struct {
	resolve func() (T, error)
}

// Value returns the parsed value of the environment value.
//
// It panics if the environment variable is undefined or invalid.
func (v Required[T]) Value() T {
	value, err := v.resolve()
	if err != nil {
		panic(err.Error())
	}

	return value
}

func registerRequired[T any](
	s spec.Spec,
	r func() (spec.Value[T], error),
) Required[T] {
	res := spec.NewResolver(s, r)
	spec.Register(res)

	return Required[T]{
		func() (T, error) {
			v, err := res.Resolve()
			return v.Go, err
		},
	}
}

func registerOptional[T any](
	s spec.Spec,
	r func() (spec.Value[T], error),
) Optional[T] {
	s.Necessity = spec.Optional

	res := spec.NewResolver(s, r)
	spec.Register(res)

	return Optional[T]{
		func() (T, error) {
			v, err := res.Resolve()
			return v.Go, err
		},
	}
}
