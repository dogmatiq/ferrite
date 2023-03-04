package variable

import (
	"errors"
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

	// Constraints returns a list of additional constraints on the variable's
	// value.
	Constraints() []Constraint

	// Examples returns a list of additional examples.
	//
	// The implementation MUST return at least one example.
	Examples() []Example

	// Documentation returns a list of chunks of documentation text.
	Documentation() []string
}

// IsDefault returns true if v is the default value of the given spec.
func IsDefault(s Spec, v Literal) bool {
	if def, ok := s.Default(); ok {
		return def == v
	}

	return false
}

// SpecOption is an option that changes the behavior of a spec.
type SpecOption[T any] func(*specOptions[T]) error

type specOptions[T any] struct {
	Constraints []TypedConstraint[T]
	Examples    []TypedExample[T]
	Docs        []string
}

// TypedSpec builds a specification for a variable depicted by type T.
type TypedSpec[T any] struct {
	name        string
	desc        string
	def         maybe.Value[valueOf[T]]
	required    bool
	schemax     TypedSchema[T]
	examples    []Example
	docs        []string
	constraints []TypedConstraint[T]
}

// NewSpec returns a new specification.
func NewSpec[T any, S TypedSchema[T]](
	name, desc string,
	def maybe.Value[T],
	req bool,
	schema S,
	options ...SpecOption[T],
) (TypedSpec[T], error) {
	if name == "" {
		return TypedSpec[T]{}, SpecError{
			cause: errors.New("variable name must not be empty"),
		}
	}

	if desc == "" {
		return TypedSpec[T]{}, SpecError{
			name:  name,
			cause: errors.New("variable description must not be empty"),
		}
	}

	if err := schema.Finalize(); err != nil {
		return TypedSpec[T]{}, SpecError{
			name:  name,
			cause: err,
		}
	}

	var opts specOptions[T]
	for _, opt := range options {
		opt(&opts)
	}

	spec := TypedSpec[T]{
		name:        name,
		desc:        desc,
		schemax:     schema,
		required:    req,
		constraints: opts.Constraints,
	}

	for _, eg := range opts.Examples {
		lit, err := spec.Marshal(eg.Native)
		if err != nil {
			return TypedSpec[T]{}, SpecError{
				name:  name,
				cause: fmt.Errorf("example value: %w", err),
			}
		}

		spec.examples = appendExample(spec.examples, Example{
			Canonical:   lit,
			Description: eg.Description,
			IsNormative: eg.IsNormative,
		})
	}

	hasOtherExamples := len(spec.examples) > 0 || !def.IsEmpty()

	for _, eg := range schema.Examples(hasOtherExamples) {
		lit, err := spec.Marshal(eg.Native)

		// Append the example if it meets all of the constraints, otherwise just
		// ignore it.
		if err == nil {
			spec.examples = appendExample(spec.examples, Example{
				Canonical:   lit,
				Description: eg.Description,
				IsNormative: eg.IsNormative,
			})
		}
	}

	if v, ok := def.Get(); ok {
		lit, err := spec.Marshal(v)
		if err != nil {
			return TypedSpec[T]{}, SpecError{
				name:  name,
				cause: fmt.Errorf("default value: %w", err),
			}
		}

		spec.def = maybe.Some(valueOf[T]{
			native:    v,
			canonical: lit,
			isDefault: true,
		})

		// Append an example of the default value if one is not already present.
		spec.examples = appendExample(spec.examples, Example{
			Canonical:   lit,
			IsNormative: true,
		})
	}

	spec.docs = opts.Docs

	return spec, nil
}

// Name returns the name of the variable.
func (s TypedSpec[T]) Name() string {
	return s.name
}

// Description returns a human-readable description of the variable.
func (s TypedSpec[T]) Description() string {
	return s.desc
}

// Schema returns the schema that applies to the variable's value.
func (s TypedSpec[T]) Schema() Schema {
	return s.schemax
}

// Default returns the string representation of the default value.
func (s TypedSpec[T]) Default() (Literal, bool) {
	return maybe.Map(s.def, valueOf[T].Canonical).Get()
}

// IsRequired returns true if the application MUST have a value for this
// variable (even if it is fulfilled by a default value).
func (s TypedSpec[T]) IsRequired() bool {
	return s.required
}

// Constraints returns a list of additional constraints on the variable's
// value.
func (s TypedSpec[T]) Constraints() []Constraint {
	constraints := make([]Constraint, len(s.constraints))
	for i, c := range s.constraints {
		constraints[i] = c
	}
	return constraints
}

// Examples returns a list of examples of valid values.
func (s TypedSpec[T]) Examples() []Example {
	return s.examples
}

// Documentation returns a list of chunks of documentation text.
func (s TypedSpec[T]) Documentation() []string {
	return s.docs
}

// CheckConstraints returns an error if v does not satisfy any one of the
// specification's constraints.
func (s TypedSpec[T]) CheckConstraints(v T) ConstraintError {
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
func (s TypedSpec[T]) Marshal(v T) (Literal, error) {
	if err := s.CheckConstraints(v); err != nil {
		return Literal{}, err
	}

	return s.schemax.Marshal(v)
}

// Unmarshal converts a literal value to it's native representation.
//
// It returns an error if v does not meet the specification's constraints or
// unmarshaling fails at the schema level.
func (s TypedSpec[T]) Unmarshal(v Literal) (T, Literal, error) {
	n, err := s.schemax.Unmarshal(v)
	if err != nil {
		return n, Literal{}, err
	}

	if err := s.CheckConstraints(n); err != nil {
		return n, Literal{}, err
	}

	c, err := s.schemax.Marshal(n)
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
