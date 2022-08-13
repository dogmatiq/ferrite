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

// TypedSpec builds a specification for a variable depicted by type T.
type TypedSpec[T any] struct {
	name       Name
	desc       string
	schema     TypedSchema[T]
	def        maybe.Value[valueOf[T]]
	isOptional bool
}

// finalizeSpec returns the completed specification.
//
// It panics if the specification is invalid.
func finalizeSpec[T any](s PendingSpec[T]) TypedSpec[T] {
	if s.Name == "" {
		s.Invalid("variable name must not be empty")
	}

	if s.Description == "" {
		s.Invalid("variable description must not be empty")
	}

	if s.Schema == nil {
		s.Invalid("a schema must be specified")
	}

	spec := TypedSpec[T]{
		name:       s.Name,
		desc:       s.Description,
		schema:     s.Schema,
		isOptional: s.IsOptional,
	}

	if v, ok := s.Default.Get(); ok {
		lit, err := s.Schema.Marshal(v)
		if err != nil {
			s.Invalid("default value: %w", err)
		}

		spec.def = maybe.Some(valueOf[T]{
			native:    v,
			canonical: lit,
			isDefault: true,
		})
	}

	return spec
}

// Name returns the name of the variable.
func (s TypedSpec[T]) Name() Name {
	return s.name
}

// Description returns a human-readable description of the variable.
func (s TypedSpec[T]) Description() string {
	return s.desc
}

// Schema returns the schema that applies to the variable's value.
func (s TypedSpec[T]) Schema() Schema {
	return s.schema
}

// Default returns the string representation of the default value.
func (s TypedSpec[T]) Default() maybe.Value[Literal] {
	return maybe.Map(s.def, valueOf[T].Canonical)
}

// IsOptional returns true if the application can handle the absence of a
// value for this variable.
func (s TypedSpec[T]) IsOptional() bool {
	return s.isOptional
}

// PendingSpec is a specification for a variable that is not yet complete.
type PendingSpec[T any] struct {
	Name        Name
	Description string
	Schema      TypedSchema[T]
	Default     maybe.Value[T]
	IsOptional  bool
}

// InvalidErr marks the specification as invalid.
func (s PendingSpec[T]) InvalidErr(err error) {
	panic(SpecError{s.Name, err}.Error())
}

// Invalid marks the specification as invalid.
func (s PendingSpec[T]) Invalid(f string, v ...any) {
	s.InvalidErr(fmt.Errorf(f, v...))
}

// SpecError represents a problem with a variable specification itself, rather
// than the variable's value.
type SpecError struct {
	name  Name
	cause error
}

// Name returns the name of the environment variable.
func (e SpecError) Name() Name {
	return e.name
}

func (e SpecError) Unwrap() error {
	return e.cause
}

func (e SpecError) Error() string {
	if e.name == "" {
		return fmt.Sprintf("invalid specification: %s", e.cause)
	}

	return fmt.Sprintf("specification for %s is invalid: %s", e.name, e.cause)
}
