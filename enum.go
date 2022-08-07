package ferrite

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite/internal/optional"
	"github.com/dogmatiq/ferrite/schema"
	"github.com/dogmatiq/ferrite/spec"
	"golang.org/x/exp/slices"
)

// Enum configures an environment variable as an enumeration with members of
// type T.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Enum[T any](name, desc string) EnumBuilder[T] {
	return EnumBuilder[T]{
		name: name,
		desc: desc,
	}.WithRenderer(
		func(v T) string {
			return fmt.Sprint(v)
		},
	)
}

// EnumBuilder is the specification for an enumeration.
type EnumBuilder[T any] struct {
	name   string
	desc   string
	render func(T) string
	values []T
	def    optional.Optional[T]
}

// WithMembers adds members to the enum.
//
// The environment variable must be set to the string representation of one of
// the member values. WithMembers must not have an empty string representation.
func (b EnumBuilder[T]) WithMembers(values ...T) EnumBuilder[T] {
	b.values = values
	return b
}

// WithRenderer sets the function used to generate the literal string
// representation of the enum's member values.
func (b EnumBuilder[T]) WithRenderer(fn func(T) string) EnumBuilder[T] {
	b.render = fn
	return b
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b EnumBuilder[T]) WithDefault(v T) EnumBuilder[T] {
	b.def = optional.With(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b EnumBuilder[T]) Required() Required[T] {
	return registerRequired(b.spec(), b.resolve)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b EnumBuilder[T]) Optional() Optional[T] {
	return registerOptional(b.spec(), b.resolve)
}

func (b *EnumBuilder[T]) spec() spec.Spec {
	if len(b.values) == 0 {
		panic(fmt.Sprintf(
			"specification for %s is invalid: no enum members are defined",
			b.name,
		))
	}

	var (
		oneOf    schema.OneOf
		literals []string
	)

	for _, v := range b.values {
		lit := b.render(v)

		if slices.Contains(literals, lit) {
			panic(fmt.Sprintf(
				"specification for %s is invalid: multiple members use %q as their literal representation",
				b.name,
				lit,
			))
		}

		if lit == "" {
			b.def = b.def.Coalesce(v)
		}

		literals = append(literals, lit)
		oneOf = append(oneOf, schema.Literal(lit))
	}

	s := spec.Spec{
		Name:        b.name,
		Description: b.desc,
		Necessity:   spec.Required,
		Schema:      oneOf,
	}

	if v, ok := b.def.Get(); ok {
		s.Necessity = spec.Defaulted
		s.Default = b.render(v)

		if !slices.Contains(literals, s.Default) {
			panic(fmt.Sprintf(
				"specification for %s is invalid: the default value must be one of the enum members, got %q",
				b.name,
				s.Default,
			))
		}
	}

	return s
}

func (b EnumBuilder[T]) resolve() (spec.Value[T], error) {
	env := os.Getenv(b.name)

	if env == "" {
		if v, ok := b.def.Get(); ok {
			return spec.Value[T]{
				Go:        v,
				Env:       b.render(v),
				IsDefault: true,
			}, nil
		}

		return spec.Value[T]{}, UndefinedError{Name: b.name}
	}

	for _, v := range b.values {
		if b.render(v) == env {
			return spec.Value[T]{
				Go:  v,
				Env: env,
			}, nil
		}
	}

	return spec.Value[T]{}, fmt.Errorf(
		"%s must be one of one of the enum members, got %q",
		b.name,
		env,
	)
}
