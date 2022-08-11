package variable

import (
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
)

// Spec builds a specification for a variable depicted by type T.
type Spec[T any] struct {
	Name        Name
	Description string
	Class       ClassOf[T]
	Default     maybe.Value[T]
	IsOptional  bool
}

// NewSpec creates a new specification for a variable depicted by type T.
func NewSpec[T any](
	name string,
	desc string,
) Spec[T] {
	s := Spec[T]{
		Name:        Name(name),
		Description: desc,
	}

	if s.Name == "" {
		s.Invalid("variable name must not be empty")
	}

	if s.Description == "" {
		s.Invalid("variable description must not be empty")
	}

	return s
}

// InvalidErr marks the specification as invalid.
func (s *Spec[T]) InvalidErr(err error) {
	panic(SpecError{s.Name, err}.Error())
}

// Invalid marks the specification as invalid.
func (s *Spec[T]) Invalid(f string, v ...any) {
	s.InvalidErr(fmt.Errorf(f, v...))
}

// Example is an example environment variable value.
type Example struct {
	// Description is a description of the example and/or the value.
	Description string

	// String is the example value.
	String String
}
