package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/variable"
)

// Required is the application-facing representation of an environment variable
// that must have a value.
type Required[T any] interface {
	// Value returns the parsed and validated value of the environment variable.
	//
	// It panics if the environment variable is undefined or invalid.
	Value() T
}

// Optional is the application-facing representation of an environment variable
// that may optionally be undefined.
type Optional[T any] interface {
	// Value returns the parsed and validated value of the environment variable,
	// if it is defined.
	//
	// If the environment variable is not defined (and there is no default
	// value), ok is false; otherwise, ok is true and v is the value.
	//
	// It panics if the environment variable is defined but invalid.
	Value() (v T, ok bool)
}

type req[T any] struct {
	v *variable.TypedVariable[T]
}

func (r req[T]) Value() T {
	m, err := r.v.Value()
	if err != nil {
		panic(err.Error())
	}

	if v, ok := m.Get(); ok {
		return v
	}

	panic(fmt.Sprintf("%s is undefined and does not have a default value", r.v.Name()))
}

type opt[T any] struct {
	v *variable.TypedVariable[T]
}

func (o opt[T]) Value() (T, bool) {
	v, err := o.v.Value()
	if err != nil {
		panic(err.Error())
	}

	return v.Get()
}
