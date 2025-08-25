package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// Bool configures an environment variable as a boolean.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Bool(name, desc string) *BoolBuilder[bool] {
	return BoolAs[bool](name, desc)
}

// BoolAs configures an environment variable as a boolean using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func BoolAs[T ~bool](name, desc string) *BoolBuilder[T] {
	b := &BoolBuilder[T]{
		schema: variable.TypedSet[T]{
			Members: []variable.SetMember[T]{
				{Value: true},
				{Value: false},
			},
			ToLiteral: func(v T) variable.Literal {
				return variable.Literal{String: fmt.Sprint(v)}
			},
		},
	}

	b.builder.Name(name)
	b.builder.Description(desc)

	return b
}

// BoolBuilder builds a specification for a boolean value.
type BoolBuilder[T ~bool] struct {
	schema  variable.TypedSet[T]
	builder variable.TypedSpecBuilder[T]
}

var _ isBuilderOfMinimal[
	bool,
	*BoolBuilder[bool],
]

// WithLiterals overrides the default literals used to represent true and false.
//
// The default literals "true" and "false" are no longer valid values when using
// custom literals.
func (b *BoolBuilder[T]) WithLiterals(t, f string) *BoolBuilder[T] {
	b.schema.ToLiteral = func(v T) variable.Literal {
		if v {
			return variable.Literal{String: t}
		}
		return variable.Literal{String: f}
	}

	return b
}

// WithDefault sets the default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *BoolBuilder[T]) WithDefault(v T) *BoolBuilder[T] {
	b.builder.Default(variable.ConstDefault(v))
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *BoolBuilder[T]) Required(options ...RequiredOption) Required[T] {
	return required(b.schema, &b.builder, options...)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *BoolBuilder[T]) Optional(options ...OptionalOption) Optional[T] {
	return optional(b.schema, &b.builder, options...)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *BoolBuilder[T]) Deprecated(options ...DeprecatedOption) Deprecated[T] {
	return deprecated(b.schema, &b.builder, options...)
}
