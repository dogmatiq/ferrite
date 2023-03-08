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

// requiredVar is a simple implementation of Required[T] that obtains the value
// from a single environment variable of the same type.
type requiredVar[T any] struct {
	v *variable.OfType[T]
}

func (r requiredVar[T]) Value() T {
	n, ok, err := r.v.NativeValue()
	if err != nil {
		panic(err.Error())
	}
	if !ok {
		panic(undefinedError(r.v).Error())
	}

	return n
}

// requiredFunc is an implementation of Required[T] that obtains the value
// from an arbitrary function.
type requiredFunc[T any] struct {
	fn func() (T, error)
}

func (r requiredFunc[T]) Value() T {
	n, err := r.fn()
	if err != nil {
		panic(err.Error())
	}
	return n
}

// optionalVar is a simple implementation of Optional[T] that obtains the value
// from a single environment variable of the same type.
type optionalVar[T any] struct {
	v *variable.OfType[T]
}

func (o optionalVar[T]) Value() (v T, ok bool) {
	n, ok, err := o.v.NativeValue()
	if err != nil {
		panic(err.Error())
	}

	return n, ok
}

// optionalFunc is an implementation of Optional[T] that obtains the value
// from an arbitrary function.
type optionalFunc[T any] struct {
	fn func() (T, bool, error)
}

func (r optionalFunc[T]) Value() (T, bool) {
	n, ok, err := r.fn()
	if err != nil {
		panic(err.Error())
	}
	return n, ok
}

// undefinedError returns an error that indicates that a variable is undefined.
func undefinedError(v variable.Any) error {
	return fmt.Errorf(
		"%s is undefined and does not have a default value",
		v.Spec().Name(),
	)
}
