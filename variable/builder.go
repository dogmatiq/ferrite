package variable

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
)

// Builder builds a variable schema, specification and ultimately registers the
// variable with a registry.
type Builder[T any] struct {
	spec     TypedSpec[T]
	def      maybe.Value[T]
	examples []TypedExample[T]
}

// Init initializes the builder.
func (b *Builder[T]) Init(name, desc string) {
	b.spec.name = name
	b.spec.desc = desc
}

// Default sets the default value for the variable.
func (b *Builder[T]) Default(v T) {
	b.def = maybe.Some(v)
}

// BuiltInConstraint adds a constraint to the variable's value.
//
// If fn was supplied by the application developer (as opposed to from within
// Ferrite itself), use UserConstraint() instead.
func (b *Builder[T]) BuiltInConstraint(
	desc string,
	fn func(T) ConstraintError,
) {
	b.spec.constraints = append(
		b.spec.constraints,
		ConstraintFunc[T]{desc, false, fn},
	)
}

// UserConstraint adds a user-defined constraint to the variable's value.
func (b *Builder[T]) UserConstraint(
	desc string,
	fn func(T) ConstraintError,
) {
	b.spec.constraints = append(
		b.spec.constraints,
		ConstraintFunc[T]{desc, true, fn},
	)
}

// MarkSensitive marks the variable's content as sensitive.
func (b *Builder[T]) MarkSensitive() {
	b.spec.sensitive = true
}

// MarkRequired marks the variable as required.
func (b *Builder[T]) MarkRequired() {
	b.spec.required = true
}

// NormativeExample adds a normative example to the variable.
//
// A normative example is one that is meaningful in the context of the
// variable's use.
func (b *Builder[T]) NormativeExample(v T, desc string) {
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
func (b *Builder[T]) NonNormativeExample(v T, desc string) {
	b.examples = append(b.examples, TypedExample[T]{v, desc, false})
}

// Documentation adds documentation to the variable.
//
// It returns the DocumentBuilder that can be used to add documentation content .
func (b *Builder[T]) Documentation(summary string) DocumentationBuilder {
	return DocumentationBuilder{
		&b.spec.docs,
		Documentation{
			Summary: summary,
		},
	}
}

// Done builds the specification and registers the variable.
func (b *Builder[T]) Done(
	schema TypedSchema[T],
	options []RegisterOption,
) *OfType[T] {
	b.spec.schema = schema

	if err := b.finalizeSpec(); err != nil {
		panic(err.Error())
	}

	return Register(b.spec, options)
}

func (b *Builder[T]) finalizeSpec() error {
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
			isDefault: true,
		})
	}

	if err := addExamples(&b.spec, b.examples); err != nil {
		return SpecError{
			name:  b.spec.name,
			cause: fmt.Errorf("example value: %w", err),
		}
	}

	return nil
}
