package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/variable"
)

// Enum configures an environment variable as an enumeration.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Enum(name, desc string) *EnumBuilder[string] {
	return EnumAs[string](name, desc)
}

// EnumAs configures an environment variable as an enumeration with members of
// type T.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func EnumAs[T any](name, desc string) *EnumBuilder[T] {
	b := &EnumBuilder[T]{
		schema: variable.TypedSet[T]{
			ToLiteral: func(v T) variable.Literal {
				return variable.Literal{
					String: fmt.Sprint(v),
				}
			},
		},
	}

	b.spec.Name(name)
	b.spec.Description(desc)

	return b
}

// EnumBuilder is the specification for an enumeration.
type EnumBuilder[T any] struct {
	schema variable.TypedSet[T]
	spec   variable.TypedSpecBuilder[T]
}

var _ isBuilderOf[any, *EnumBuilder[any]]

// WithMembers adds members to the enum.
//
// The environment variable must be set to the string representation of one of
// the member values. The values must not have an empty string representation.
func (b *EnumBuilder[T]) WithMembers(values ...T) *EnumBuilder[T] {
	for _, v := range values {
		b.WithMember(v, "")
	}
	return b
}

// WithMember adds a member to the enum.
//
// The environment variable must be set to the string representation of one of
// the member values. v must not have an empty string representation.
func (b *EnumBuilder[T]) WithMember(v T, desc string) *EnumBuilder[T] {
	b.schema.Members = append(
		b.schema.Members,
		variable.SetMember[T]{
			Value:       v,
			Description: desc,
		},
	)
	return b
}

// WithRenderer sets the function used to generate the literal string
// representation of the enum's member values.
func (b *EnumBuilder[T]) WithRenderer(fn func(T) variable.Literal) *EnumBuilder[T] {
	b.schema.ToLiteral = fn
	return b
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b *EnumBuilder[T]) WithDefault(v T) *EnumBuilder[T] {
	b.spec.Default(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b *EnumBuilder[T]) Required(options ...RequiredOption) Required[T] {
	return required(b.schema, &b.spec, options)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b *EnumBuilder[T]) Optional(options ...OptionalOption) Optional[T] {
	return optional(b.schema, &b.spec, options)
}

// Deprecated completes the build process and registers a deprecated variable
// with Ferrite's validation system.
func (b *EnumBuilder[T]) Deprecated(reason string, options ...DeprecatedOption) Deprecated[T] {
	return deprecated(b.schema, &b.spec, reason, options)
}
