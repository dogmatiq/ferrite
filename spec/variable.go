package spec

import "errors"

// Optional is the application-facing representation of an environment variable
// that may optionally be undefined.
//
// It is the standard implementation of ferrite.Optional.
type Optional[T any] struct {
	Resolve func() (T, error)
}

// Value returns the parsed value of the environment value, if it is defined.
//
// It panics if the environment variable is defined but invalid.
func (v Optional[T]) Value() (T, bool) {
	value, err := v.Resolve()
	if err != nil {
		if errors.As(err, &UndefinedError{}) {
			return value, false
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
	Resolve func() (T, error)
}

// Value returns the parsed value of the environment value.
//
// It panics if the environment variable is undefined or invalid.
func (v Required[T]) Value() T {
	value, err := v.Resolve()
	if err != nil {
		panic(err.Error())
	}

	return value
}

func RegisterRequired[T any](
	s Spec,
	r func() (ValueOf[T], error),
) Required[T] {
	res := NewResolver(s, r)
	Register(res)

	return Required[T]{
		func() (T, error) {
			v, err := res.ResolveTyped()
			return v.Go, err
		},
	}
}

func RegisterOptional[T any](
	s Spec,
	r func() (ValueOf[T], error),
) Optional[T] {
	s.IsOptional = true

	res := NewResolver(s, r)
	Register(res)

	return Optional[T]{
		func() (T, error) {
			v, err := res.ResolveTyped()
			return v.Go, err
		},
	}
}
