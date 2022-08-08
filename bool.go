package ferrite

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite/internal/optional"
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
	return BoolBuilder[T]{
		name: name,
		desc: desc,
	}.WithLiterals(
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
		panic(fmt.Sprintf(
			"specification for %s is invalid: boolean literals must not be zero-length",
			b.name,
		))
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
	return registerRequired(b.spec(), b.resolve)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b BoolBuilder[T]) Optional() Optional[T] {
	return registerOptional(b.spec(), b.resolve)
}

func (b BoolBuilder[T]) spec() spec.Spec {
	s := spec.Spec{
		Name:        b.name,
		Description: b.desc,
		Necessity:   spec.Required,
		Schema: spec.OneOf{
			spec.Literal(b.t),
			spec.Literal(b.f),
		},
	}

	if v, ok := b.def.Get(); ok {
		s.Necessity = spec.Defaulted
		s.Default = b.render(v)
	}

	return s
}

func (b BoolBuilder[T]) resolve() (spec.ValueOf[T], error) {
	switch env := os.Getenv(b.name); env {
	case b.t, b.f:
		return spec.ValueOf[T]{
			Go:  env == b.t,
			Env: env,
		}, nil

	case "":
		if v, ok := b.def.Get(); ok {
			return spec.ValueOf[T]{
				Go:    v,
				Env:   b.render(v),
				IsDef: true,
			}, nil
		}

		return spec.Undefined[T](b.name)

	default:
		return spec.ValueOf[T]{}, fmt.Errorf(
			`%s must be either "%s" or "%s", got "%s"`,
			b.name,
			b.t,
			b.f,
			env,
		)
	}
}

func (b BoolBuilder[T]) render(v T) string {
	if v {
		return b.t
	}
	return b.f
}
