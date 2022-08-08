package spec

import (
	"fmt"
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

// Undefined returns a new UndefinedError.
func Undefined[T any](
	name string,
) (ValueOf[T], error) {
	return ValueOf[T]{}, UndefinedError{
		Name: name,
	}
}

// Invalid returns a new ValidationError.
func Invalid[T any](
	name string,
	value string,
	f string, v ...any,
) (ValueOf[T], error) {
	return ValueOf[T]{}, ValidationError{
		Name:  name,
		Value: value,
		Cause: fmt.Errorf(f, v...),
	}
}
