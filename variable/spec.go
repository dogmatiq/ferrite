package variable

import (
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
)

// Spec is a specification of a variable.
type Spec interface {
	// Name returns the name of the variable.
	Name() Name

	// Description returns a human-readable description of the variable.
	Description() string

	// Schema returns the schema that applies to the variable's value.
	Schema() Schema

	// Default returns the string representation of the default value.
	Default() maybe.Value[Literal]

	// IsOptional returns true if the application can handle the absence of a
	// value for this variable.
	IsOptional() bool
}

// SpecFor builds a specification for a variable depicted by type T.
type SpecFor[T any] struct {
	name       Name
	desc       string
	schema     SchemaFor[T]
	def        maybe.Value[T]
	isOptional bool
}

// NewSpec creates a new specification for a variable depicted by type T.
func NewSpec[T any](
	name string,
	desc string,
) SpecFor[T] {
	s := SpecFor[T]{
		name: Name(name),
		desc: desc,
	}

	if s.name == "" {
		s.Invalid("variable name must not be empty")
	}

	if s.desc == "" {
		s.Invalid("variable description must not be empty")
	}

	return s
}

// Name returns the name of the variable.
func (s SpecFor[T]) Name() Name {
	return s.name
}

// Description returns a human-readable description of the variable.
func (s SpecFor[T]) Description() string {
	return s.desc
}

// Schema returns the schema that applies to the variable's value.
func (s SpecFor[T]) Schema() Schema {
	return s.schema
}

// Default returns the string representation of the default value.
func (s SpecFor[T]) Default() maybe.Value[Literal] {
	return maybe.Map(s.def, s.schema.Marshal)
}

// IsOptional returns true if the application can handle the absence of a
// value for this variable.
func (s SpecFor[T]) IsOptional() bool {
	return s.isOptional
}

// InvalidErr marks the specification as invalid.
func (s SpecFor[T]) InvalidErr(err error) {
	panic(SpecError{s.name, err}.Error())
}

// Invalid marks the specification as invalid.
func (s SpecFor[T]) Invalid(f string, v ...any) {
	s.InvalidErr(fmt.Errorf(f, v...))
}

// SetSchema sets the schema for the variable's values.
func (s *SpecFor[T]) SetSchema(sc SchemaFor[T]) {
	s.schema = sc
}

// SetDefault sets the environment variable's default value.
func (s *SpecFor[T]) SetDefault(v T) {
	s.def = maybe.Some(v)
}

// MarkOptional marks the environment variable as optional, meaning that the
// application can operate without any value for the variable.
func (s *SpecFor[T]) MarkOptional() {
	s.isOptional = true
}

// Example is an example environment variable value.
type Example struct {
	// Description is a description of the example and/or the value.
	Description string

	// Literal is the example value.
	Literal Literal
}
