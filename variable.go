package ferrite

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/ferrite/spec"
)

// UndefinedError is returned by a resolver when an environment variable has
// neither an explicit nor a default value.
type UndefinedError struct {
	Name string
}

func (e UndefinedError) Error() string {
	return fmt.Sprintf("%s is undefined and does not have a default value", e.Name)
}

// Optional is the application-facing representation of an environment variable
// that may optionally be undefined.
//
// It is the standard implementation of ferrite.Optional.
type Optional[T any] struct {
	resolver *spec.Resolver[T]
}

// Value returns the parsed value of the environment value, if it is defined.
//
// It panics if the environment variable is defined but invalid.
func (v Optional[T]) Value() (T, bool) {
	value, err := v.resolver.Resolve()
	if err != nil {
		if errors.As(err, &UndefinedError{}) {
			return value.Go, false
		}

		panic(err.Error())
	}

	return value.Go, true
}

// Required is the application-facing representation of an environment
// variable that must always have a valid value.
//
// It is the standard implementation of ferrite.Required.
type Required[T any] struct {
	resolver *spec.Resolver[T]
}

// Value returns the parsed value of the environment value.
//
// It panics if the environment variable is undefined or invalid.
func (v Required[T]) Value() T {
	value, err := v.resolver.Resolve()
	if err != nil {
		panic(err.Error())
	}

	return value.Go
}
