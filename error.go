package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/spec"
)

// UndefinedError indicates that an environment variable has neither an explicit
// nor a default value.
type UndefinedError struct {
	Name  string
	Cause error
}

func (e UndefinedError) Unwrap() error {
	return e.Cause
}

func (e UndefinedError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s is undefined: %s", e.Name, e.Cause)
	}

	return fmt.Sprintf("%s is undefined and does not have a default value", e.Name)
}

// ValidationError indicates that an environment variable value is invalid.
type ValidationError struct {
	Name  string
	Value string
	Cause error
}

func (e ValidationError) Unwrap() error {
	return e.Cause
}

func (e ValidationError) Error() string {
	return fmt.Sprintf(
		"%s (%q) is invalid: %s",
		e.Name,
		e.Value,
		e.Cause,
	)
}

// undefined returns a new UndefinedError.
func undefined[T any](
	name string,
) (spec.ValueOf[T], error) {
	return spec.ValueOf[T]{}, UndefinedError{
		Name: name,
	}
}

// invalid returns a new ValidationError.
func invalid[T any](
	name string,
	value string,
	f string, v ...any,
) (spec.ValueOf[T], error) {
	return spec.ValueOf[T]{}, ValidationError{
		Name:  name,
		Value: value,
		Cause: fmt.Errorf(f, v...),
	}
}
