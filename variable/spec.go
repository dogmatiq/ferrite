package variable

import (
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
)

// Spec is a specification of a variable.
type Spec interface {
	// Name returns the name of the variable.
	Name() string

	// Description returns a human-readable description of the variable.
	Description() string

	// Schema returns the schema that applies to the variable's value.
	Schema() Schema

	// Default returns the string representation of the default value.
	Default() (Literal, bool)

	// IsRequired returns true if the application MUST have a value for this
	// variable (even if it is fulfilled by a default value).
	IsRequired() bool

	// IsSensitive returns true if the variable's value contains sensitive
	// information.
	IsSensitive() bool

	// IsDeprecated returns true if the variable is deprecated.
	IsDeprecated() bool

	// Constraints returns a list of additional constraints on the variable's
	// value.
	Constraints() []Constraint

	// Examples returns a list of additional examples.
	//
	// The implementation MUST return at least one example.
	Examples() []Example

	// Documentation returns a list of chunks of documentation text.
	Documentation() []Documentation

	// Relationships returns a list of relationships that involve this variable.
	Relationships() []Relationship

	// addRelationship adds a relationship that involves this variable.
	addRelationship(r Relationship)
}

// IsDefault returns true if v is the default value of the given spec.
func IsDefault(s Spec, v Literal) bool {
	if def, ok := s.Default(); ok {
		return def == v
	}

	return false
}

// TypedSpec builds a specification for a variable depicted by type T.
type TypedSpec[T any] struct {
	name          string
	desc          string
	def           maybe.Value[valueOf[T]]
	required      bool
	sensitive     bool
	deprecated    bool
	schema        TypedSchema[T]
	examples      []Example
	docs          []Documentation
	constraints   []TypedConstraint[T]
	relationships []Relationship
}

// Name returns the name of the variable.
func (s *TypedSpec[T]) Name() string {
	return s.name
}

// Description returns a human-readable description of the variable.
func (s *TypedSpec[T]) Description() string {
	return s.desc
}

// Schema returns the schema that applies to the variable's value.
func (s *TypedSpec[T]) Schema() Schema {
	return s.schema
}

// Default returns the string representation of the default value.
func (s *TypedSpec[T]) Default() (Literal, bool) {
	return maybe.Map(s.def, valueOf[T].Canonical).Get()
}

// IsRequired returns true if the application MUST have a value for this
// variable (even if it is fulfilled by a default value).
func (s *TypedSpec[T]) IsRequired() bool {
	return s.required
}

// IsSensitive returns true if the variable's value contains sensitive
// information.
func (s *TypedSpec[T]) IsSensitive() bool {
	return s.sensitive
}

// IsDeprecated returns true if the variable is deprecated.
func (s *TypedSpec[T]) IsDeprecated() bool {
	return s.deprecated
}

// Constraints returns a list of additional constraints on the variable's
// value.
func (s *TypedSpec[T]) Constraints() []Constraint {
	constraints := make([]Constraint, len(s.constraints))
	for i, c := range s.constraints {
		constraints[i] = c
	}
	return constraints
}

// Examples returns a list of examples of valid values.
func (s *TypedSpec[T]) Examples() []Example {
	return s.examples
}

// Documentation returns a list of chunks of documentation text.
func (s *TypedSpec[T]) Documentation() []Documentation {
	return s.docs
}

// Relationships returns a list of relationships that involve this variable.
func (s TypedSpec[T]) Relationships() []Relationship {
	return s.relationships
}

// AddRelationship adds a relationship that involves this variable.
func (s *TypedSpec[T]) addRelationship(r Relationship) {
	s.relationships = append(s.relationships, r)
}

// CheckConstraints returns an error if v does not satisfy any one of the
// specification's constraints.
func (s *TypedSpec[T]) CheckConstraints(v T) ConstraintError {
	for _, c := range s.constraints {
		if err := c.Check(v); err != nil {
			return err
		}
	}

	return nil
}

// Marshal converts a value to its literal representation.
//
// It returns an error if v does not meet the specification's constraints or
// marshaling fails at the schema level.
func (s *TypedSpec[T]) Marshal(v T) (Literal, error) {
	if err := s.CheckConstraints(v); err != nil {
		return Literal{}, err
	}

	return s.schema.Marshal(v)
}

// Unmarshal converts a literal value to it's native representation.
//
// It returns an error if v does not meet the specification's constraints or
// unmarshaling fails at the schema level.
func (s *TypedSpec[T]) Unmarshal(v Literal) (T, Literal, error) {
	n, err := s.schema.Unmarshal(v)
	if err != nil {
		return n, Literal{}, err
	}

	if err := s.CheckConstraints(n); err != nil {
		return n, Literal{}, err
	}

	c, err := s.schema.Marshal(n)
	if err != nil {
		// Schema can't marshal a value it just successfully unmarshaled!
		panic(err)
	}

	return n, c, err
}

// SpecError represents a problem with a variable specification itself, rather
// than the variable's value.
type SpecError struct {
	name  string
	cause error
}

// Name returns the name of the environment variable.
func (e SpecError) Name() string {
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
