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

// // NewBuilder returns a specification builder that builds a TypedSpec[T].
// func NewBuilder[T any, S TypedSchema[T]](
// 	name, desc string,
// ) *Builder[T, S] {
// 	return &Builder[T, S]{
// 		spec: TypedSpec[T]{
// 			name: name,
// 			desc: desc,
// 		},
// 	}
// }

func (b *Builder[T, S]) Begin(name, desc string) {
	b.spec.name = name
	b.spec.desc = desc
}

// Default sets the default value for the variable.
func (b *Builder[T, S]) Default(v T) {
	b.def = maybe.Some(v)
}

// Constraint adds a constraint to the variable's value.
func (b *Builder[T, S]) Constraint(c TypedConstraint[T]) {
	b.spec.constraints = append(b.spec.constraints, c)
}

// Sensitive marks the variable's content as sensitive.
func (b *Builder[T, S]) Sensitive() {
	b.spec.sensitive = true
}

// Required marks the variable as required.
func (b *Builder[T, S]) Required() {
	b.spec.required = true
}

// Register builds the specification and registers the variable.
func (b *Builder[T, S]) End(options []RegisterOption) *OfType[T] {
	spec, err := b.buildSpec()
	if err != nil {
		panic(err)
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
