package variable

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
)

// SpecBuilder builds a specification for an environment variable.
type SpecBuilder interface {
	Name(string)
	Description(string)
	MarkRequired()
	MarkDeprecated()
	MarkSensitive()
	Documentation() DocumentationBuilder
	Relationship(Relationship)
	Peek() Spec
}

// TypedSpecBuilder builds a specification for a variable of type T.
type TypedSpecBuilder[T any] struct {
	spec     TypedSpec[T]
	def      maybe.Value[T]
	examples []TypedExample[T]
}

// Name sets the name of the environment variable.
func (b *TypedSpecBuilder[T]) Name(name string) {
	b.spec.name = name
}

// Description sets a human-readable description of the environment variable.
func (b *TypedSpecBuilder[T]) Description(desc string) {
	b.spec.desc = desc
}

// Default sets the default value for the variable.
func (b *TypedSpecBuilder[T]) Default(v T) {
	b.def = maybe.Some(v)
}

// BuiltInConstraint adds a constraint to the variable's value.
//
// If fn was supplied by the application developer (as opposed to from within
// Ferrite itself), use UserConstraint() instead.
func (b *TypedSpecBuilder[T]) BuiltInConstraint(
	desc string,
	fn func(T) ConstraintError,
) {
	b.spec.constraints = append(
		b.spec.constraints,
		constraint[T]{desc, false, fn},
	)
}

// UserConstraint adds a user-defined constraint to the variable's value.
func (b *TypedSpecBuilder[T]) UserConstraint(
	desc string,
	fn func(T) ConstraintError,
) {
	b.spec.constraints = append(
		b.spec.constraints,
		constraint[T]{desc, true, fn},
	)
}

// MarkRequired marks the variable as required.
func (b *TypedSpecBuilder[T]) MarkRequired() {
	b.spec.required = true
}

// MarkSensitive marks the variable's content as sensitive.
func (b *TypedSpecBuilder[T]) MarkSensitive() {
	b.spec.sensitive = true
}

// MarkDeprecated marks the variable as deprecated.
func (b *TypedSpecBuilder[T]) MarkDeprecated() {
	b.spec.deprecated = true
}

// NormativeExample adds a normative example to the variable.
//
// A normative example is one that is meaningful in the context of the
// variable's use.
func (b *TypedSpecBuilder[T]) NormativeExample(v T, desc string) {
	b.examples = append(b.examples, TypedExample[T]{v, desc, true})
}

// NonNormativeExample adds a non-normative example to the variable.
//
// A non-normative example is one that may not be meaningful in the context of
// the variable's use, but is included for illustrative purposes.
//
// For example, a variable that represents a URL may have a non-normative
// example of "https://example.org/path", even if the actual use-case for the
// variable requires an "ftp" URL.
func (b *TypedSpecBuilder[T]) NonNormativeExample(v T, desc string) {
	b.examples = append(b.examples, TypedExample[T]{v, desc, false})
}

// Documentation adds documentation to the variable.
//
// It returns the DocumentBuilder that can be used to add documentation content .
func (b *TypedSpecBuilder[T]) Documentation() DocumentationBuilder {
	return DocumentationBuilder{
		docs: &b.spec.docs,
	}
}

// Relationship adds a relationship that involves this variable.
func (b *TypedSpecBuilder[T]) Relationship(rel Relationship) {
	b.spec.AddRelationship(rel)
}

// Peek returns the (potentially invalid) spec that is being built.
func (b *TypedSpecBuilder[T]) Peek() Spec {
	return &b.spec
}

// Done builds the specification and registers the variable.
func (b *TypedSpecBuilder[T]) Done(schema TypedSchema[T]) *TypedSpec[T] {
	b.spec.schema = schema

	if err := b.finalize(); err != nil {
		panic(err.Error())
	}

	return &b.spec
}

func (b *TypedSpecBuilder[T]) finalize() error {
	if b.spec.name == "" {
		return SpecError{
			cause: errors.New("variable name must not be empty"),
		}
	}

	if b.spec.desc == "" {
		return SpecError{
			name:  b.spec.name,
			cause: errors.New("variable description must not be empty"),
		}
	}

	if err := b.spec.schema.Finalize(); err != nil {
		return SpecError{
			name:  b.spec.name,
			cause: err,
		}
	}

	if v, ok := b.def.Get(); ok {
		lit, err := b.spec.Marshal(v)
		if err != nil {
			return SpecError{
				name:  b.spec.name,
				cause: fmt.Errorf("default value: %w", err),
			}
		}

		b.spec.def = maybe.Some(valueOf[T]{
			native:    v,
			canonical: lit,
		})
	}

	if err := b.buildExamples(); err != nil {
		return SpecError{
			name:  b.spec.name,
			cause: fmt.Errorf("example value: %w", err),
		}
	}

	return nil
}

// buildExamples builds the examples for the spec from various sources.
func (b *TypedSpecBuilder[T]) buildExamples() error {
	uniq := map[Literal]struct{}{}

	// Add the examples provided directly to the builder.
	for _, eg := range b.examples {
		lit, err := b.spec.Marshal(eg.Native)
		if err != nil {
			return err
		}

		if _, ok := uniq[lit]; !ok {
			uniq[lit] = struct{}{}
			b.spec.examples = append(b.spec.examples, Example{
				Canonical:   lit,
				Description: eg.Description,
				IsNormative: eg.IsNormative,
				Source:      ExampleSourceSpecBuilder,
			})
		}
	}

	// Be conservative when generating examples if there are any explicitly
	// provided examples or a default value, which are likely to be better
	// examples.
	conservative := len(b.examples) != 0 || !b.spec.def.IsEmpty()

	// Generate examples from the schema and add each one only if it meets all
	// of the constraints (and there is no existing example of the same value).
	for _, eg := range b.spec.schema.Examples(conservative) {
		if lit, err := b.spec.Marshal(eg.Native); err == nil {
			if _, ok := uniq[lit]; !ok {
				uniq[lit] = struct{}{}
				b.spec.examples = append(b.spec.examples, Example{
					Canonical:   lit,
					Description: eg.Description,
					IsNormative: eg.IsNormative,
					Source:      ExampleSourceSchema,
				})
			}
		}
	}

	// Prepend an example of the default value if there is no existing example
	// of the same value.
	if def, ok := b.spec.def.Get(); ok {
		lit := def.Canonical()

		if _, ok := uniq[lit]; !ok {
			uniq[lit] = struct{}{}
			b.spec.examples = append(
				[]Example{
					{
						Canonical:   lit,
						IsNormative: true,
						Source:      ExampleSourceSpecDefault,
					},
				},
				b.spec.examples...,
			)
		}
	}

	return nil
}
