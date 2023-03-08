package variable

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
)

// Builder builds a variable schema, specification and ultimately registers the
// variable with a registry.
type Builder[T any, S TypedSchema[T]] struct {
	Schema S

	spec     TypedSpec[T]
	def      maybe.Value[T]
	examples []TypedExample[T]
}

// Init initializes the builder.
func (b *Builder[T, S]) Init(name, desc string) {
	b.spec.name = name
	b.spec.desc = desc
}

// Default sets the default value for the variable.
func (b *Builder[T, S]) Default(v T) {
	b.def = maybe.Some(v)
}

// BuiltInConstraint adds a constraint to the variable's value.
//
// If fn was supplied by the application developer (as opposed to from within
// Ferrite itself), use UserConstraint() instead.
func (b *Builder[T, S]) BuiltInConstraint(
	desc string,
	fn func(T) ConstraintError,
) {
	b.spec.constraints = append(
		b.spec.constraints,
		ConstraintFunc[T]{desc, false, fn},
	)
}

// UserConstraint adds a user-defined constraint to the variable's value.
func (b *Builder[T, S]) UserConstraint(
	desc string,
	fn func(T) ConstraintError,
) {
	b.spec.constraints = append(
		b.spec.constraints,
		ConstraintFunc[T]{desc, true, fn},
	)
}

// Sensitive marks the variable's content as sensitive.
func (b *Builder[T, S]) Sensitive() {
	b.spec.sensitive = true
}

// Required marks the variable as required.
func (b *Builder[T, S]) Required() {
	b.spec.required = true
}

// NormativeExample adds a normative example to the variable.
//
// A normative example is one that is meaningful in the context of the
// variable's use.
func (b *Builder[T, S]) NormativeExample(v T, desc string) {
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
func (b *Builder[T, S]) NonNormativeExample(v T, desc string) {
	b.examples = append(b.examples, TypedExample[T]{v, desc, false})
}

// Documentation adds documentation to the variable.
//
// It returns the DocumentBuilder that can be used to add documentation content .
func (b *Builder[T, S]) Documentation(summary string) DocumentationBuilder {
	return DocumentationBuilder{
		&b.spec.docs,
		Documentation{
			Summary: summary,
		},
	}
}

// Done builds the specification and registers the variable.
func (b *Builder[T, S]) Done(options []RegisterOption) *OfType[T] {
	spec, err := b.buildSpec()
	if err != nil {
		panic(err.Error())
	}

	return Register(spec, options)
}

func (b *Builder[T, S]) buildSpec() (TypedSpec[T], error) {
	if b.spec.name == "" {
		return TypedSpec[T]{}, SpecError{
			cause: errors.New("variable name must not be empty"),
		}
	}

	if b.spec.desc == "" {
		return TypedSpec[T]{}, SpecError{
			name:  b.spec.name,
			cause: errors.New("variable description must not be empty"),
		}
	}

	b.spec.schema = b.Schema

	if err := b.spec.schema.Finalize(); err != nil {
		return TypedSpec[T]{}, SpecError{
			name:  b.spec.name,
			cause: err,
		}
	}

	if v, ok := b.def.Get(); ok {
		lit, err := b.spec.Marshal(v)
		if err != nil {
			return TypedSpec[T]{}, SpecError{
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
		return TypedSpec[T]{}, SpecError{
			name:  b.spec.name,
			cause: fmt.Errorf("example value: %w", err),
		}
	}

	return b.spec, nil
}
