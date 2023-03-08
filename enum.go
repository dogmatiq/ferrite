package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
	"github.com/dogmatiq/ferrite/variable"
	"golang.org/x/exp/slices"
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
		name: name,
		desc: desc,
		toLiteral: func(v T) variable.Literal {
			return variable.Literal{
				String: fmt.Sprint(v),
			}
		},
	}
}

// EnumBuilder is the specification for an enumeration.
type EnumBuilder[T any] struct {
	name, desc string
	def        maybe.Value[T]
	members    []variable.SetMember[T]
	toLiteral  func(T) variable.Literal
}

// WithMembers adds members to the enum.
//
// The environment variable must be set to the string representation of one of
// the member values. The values must not have an empty string representation.
func (b EnumBuilder[T]) WithMembers(values ...T) EnumBuilder[T] {
	for _, v := range values {
		b = b.WithMember(v, "")
	}
	return b
}

// WithMember adds a member to the enum.
//
// The environment variable must be set to the string representation of one of
// the member values. v must not have an empty string representation.
func (b EnumBuilder[T]) WithMember(v T, desc string) EnumBuilder[T] {
	b.members = append(
		slices.Clone(b.members),
		variable.SetMember[T]{
			Value:       v,
			Description: desc,
		},
	)
	return b
}

// WithRenderer sets the function used to generate the literal string
// representation of the enum's member values.
func (b EnumBuilder[T]) WithRenderer(fn func(T) variable.Literal) EnumBuilder[T] {
	b.toLiteral = fn
	return b
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b EnumBuilder[T]) WithDefault(v T) EnumBuilder[T] {
	b.def = maybe.Some(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b EnumBuilder[T]) Required(options ...variable.RegisterOption) Required[T] {
	v := variable.Register(b.spec(true), options)
	return requiredVar[T]{v}

}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b EnumBuilder[T]) Optional(options ...variable.RegisterOption) Optional[T] {
	v := variable.Register(b.spec(false), options)
	return optionalVar[T]{v}

}

func (b *EnumBuilder[T]) spec(req bool) variable.TypedSpec[T] {
	s, err := variable.NewSpec(
		b.name,
		b.desc,
		b.def,
		req,
		variable.TypedSet[T]{
			Members:   b.members,
			ToLiteral: b.toLiteral,
		},
	)
	if err != nil {
		panic(err.Error())
	}

	return s
}
