package ferrite

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite/internal/optional"
	"github.com/dogmatiq/ferrite/schema"
	"github.com/dogmatiq/ferrite/spec"
)

// Bool configures an environment variable as a boolean.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Bool(name, desc string) BoolBuilder[bool] {
	return BoolAs[bool](name, desc)
}

// BoolAs configures an environment variable as a boolean using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func BoolAs[T ~bool](name, desc string) BoolBuilder[T] {
	b := &BoolBuilder[T]{
		name: name,
		desc: desc,
	}

	return b.WithLiterals(
		fmt.Sprint(T(true)),
		fmt.Sprint(T(false)),
	)
}

// BoolBuilder builds a specification for a boolean value.
type BoolBuilder[T ~bool] struct {
	name string
	desc string
	t, f string
	def  optional.Optional[T]
}

// WithLiterals overrides the default literals used to represent true and false.
//
// The default literals "true" and "false" are no longer valid values when using
// custom literals.
func (b BoolBuilder[T]) WithLiterals(t, f string) BoolBuilder[T] {
	if t == "" || f == "" {
		panic("boolean literals must not be zero-length")
	}

	b.t = t
	b.f = f

	return b
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b BoolBuilder[T]) WithDefault(v T) BoolBuilder[T] {
	b.def = optional.With(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b BoolBuilder[T]) Required() Required[T] {
	return Required[T]{
		b.resolver(false),
	}
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b BoolBuilder[T]) Optional() Optional[T] {
	return Optional[T]{
		b.resolver(true),
	}
}

func (b BoolBuilder[T]) resolve() (spec.Value[T], error) {
	switch str := os.Getenv(b.name); str {
	case b.t, b.f:
		return spec.Value[T]{
			Parsed:     str == b.t,
			Normalized: str,
		}, nil

	case "":
		if parsed, ok := b.def.TryGet(); ok {
			return spec.Value[T]{
				Parsed:     parsed,
				Normalized: b.render(parsed),
				IsDefault:  true,
			}, nil
		}

		return spec.Value[T]{}, UndefinedError{Name: b.name}

	default:
		return spec.Value[T]{}, fmt.Errorf(
			`%s must be either "%s" or "%s", got "%s"`,
			b.name,
			b.t,
			b.f,
			str,
		)
	}
}

func (b BoolBuilder[T]) resolver(opt bool) *spec.Resolver[T] {
	s := spec.Spec{
		Name:        b.name,
		Description: b.desc,
		Necessity:   spec.Required,
		Schema: schema.OneOf{
			schema.Literal(b.t),
			schema.Literal(b.f),
		},
	}

	if v, ok := b.def.TryGet(); ok {
		s.Necessity = spec.Defaulted
		s.Default = b.render(v)
	}

	if opt {
		s.Necessity = spec.Optional
	}

	r := spec.NewResolver(s, b.resolve)
	spec.Register(r)

	return r
}

func (b BoolBuilder[T]) render(v T) string {
	if v {
		return b.t
	}
	return b.f
}
