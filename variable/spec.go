package variable

import (
	"errors"
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

	// IsRequired returns true if the application MUST have a value for this
	// variable (even if it is fulfilled by a default value).
	IsRequired() bool
}

// TypedSpec builds a specification for a variable depicted by type T.
type TypedSpec[T any] struct {
	name       Name
	desc       string
	def        maybe.Value[valueOf[T]]
	required   bool
	schema     TypedSchema[T]
	validators []Validator[T]
}

// NewSpec returns a new specification.
func NewSpec[T any, S TypedSchema[T]](
	name, desc string,
	def maybe.Value[T],
	req bool,
	schema S,
	validators ...Validator[T],
) (TypedSpec[T], error) {
	n := Name(name)

	if n == "" {
		return TypedSpec[T]{}, SpecError{
			cause: errors.New("variable name must not be empty"),
		}
	}

	if desc == "" {
		return TypedSpec[T]{}, SpecError{
			name:  n,
			cause: errors.New("variable description must not be empty"),
		}
	}

	if err := schema.Finalize(); err != nil {
		return TypedSpec[T]{}, SpecError{
			name:  n,
			cause: err,
		}
	}

	spec := TypedSpec[T]{
		name:       n,
		desc:       desc,
		schema:     schema,
		required:   req,
		validators: validators,
	}

	if v, ok := def.Get(); ok {
		for _, val := range validators {
			if err := val.Validate(v); err != nil {
				return TypedSpec[T]{}, SpecError{
					name:  n,
					cause: fmt.Errorf("default value: %w", err),
				}
			}
		}

		lit, err := schema.Marshal(v)
		if err != nil {
			return TypedSpec[T]{}, SpecError{
				name:  n,
				cause: fmt.Errorf("default value: %w", err),
			}
		}

		spec.def = maybe.Some(valueOf[T]{
			native:    v,
			canonical: lit,
			isDefault: true,
		})
	}

	return spec, nil
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

// IsRequired returns true if the application MUST have a value for this
// variable (even if it is fulfilled by a default value).
func (s TypedSpec[T]) IsRequired() bool {
	return s.required
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
