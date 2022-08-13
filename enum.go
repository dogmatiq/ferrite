package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
)

// Enum configures an environment variable as an enumeration.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Enum(name, desc string) EnumBuilder[string] {
	return EnumAs[string](name, desc)
}

// EnumAs configures an environment variable as an enumeration with members of
// type T.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func EnumAs[T any](name, desc string) EnumBuilder[T] {
	return EnumBuilder[T]{
		spec: variable.PendingSpec[T]{
			Name:        variable.Name(name),
			Description: desc,
		},
	}.WithRenderer(
		func(v T) variable.Literal {
			return variable.Literal(
				fmt.Sprint(v),
			)
		},
	)
}

// EnumBuilder is the specification for an enumeration.
type EnumBuilder[T any] struct {
	spec    variable.PendingSpec[T]
	render  func(T) variable.Literal
	members []T
}

// WithMembers adds members to the enum.
//
// The environment variable must be set to the string representation of one of
// the member values. WithMembers must not have an empty string representation.
func (b EnumBuilder[T]) WithMembers(members ...T) EnumBuilder[T] {
	b.members = members
	return b
}

// WithRenderer sets the function used to generate the literal string
// representation of the enum's member values.
func (b EnumBuilder[T]) WithRenderer(fn func(T) variable.Literal) EnumBuilder[T] {
	b.render = fn
	return b
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b EnumBuilder[T]) WithDefault(v T) EnumBuilder[T] {
	b.spec.Default = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b EnumBuilder[T]) Required(options ...variable.RegisterOption) Required[T] {
	b.finalize()
	return req[T]{variable.Register(b.spec, options)}
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b EnumBuilder[T]) Optional(options ...variable.RegisterOption) Optional[T] {
	b.spec.IsOptional = true
	b.finalize()
	return opt[T]{variable.Register(b.spec, options)}
}

func (b *EnumBuilder[T]) finalize() {
	s, err := variable.NewSet(b.members, b.render)
	if err != nil {
		b.spec.InvalidErr(err)
	}

	b.spec.Schema = s
}
