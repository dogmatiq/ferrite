package ferrite

import (
	"os"

	"github.com/dogmatiq/ferrite/internal/optional"
	"github.com/dogmatiq/ferrite/spec"
)

// String configures an environment variable as a string.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func String(name, desc string) StringBuilder[string] {
	return StringAs[string](name, desc)
}

// StringAs configures an environment variable as a string using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func StringAs[T ~string](name, desc string) StringBuilder[T] {
	return StringBuilder[T]{
		name: name,
		desc: desc,
	}
}

// StringBuilder builds a specification for a string variable.
type StringBuilder[T ~string] struct {
	name string
	desc string
	def  optional.Optional[T]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b StringBuilder[T]) WithDefault(v T) StringBuilder[T] {
	b.def = optional.With(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b StringBuilder[T]) Required() Required[T] {
	return registerRequired(b.spec(), b.resolve)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b StringBuilder[T]) Optional() Optional[T] {
	return registerOptional(b.spec(), b.resolve)
}

func (b StringBuilder[T]) spec() spec.Spec {
	s := spec.Spec{
		Name:        b.name,
		Description: b.desc,
		Schema:      spec.OfType[T](),
	}

	if v, ok := b.def.Get(); ok {
		s.HasDefault = true
		s.DefaultX = string(v)
	}

	return s
}

func (b StringBuilder[T]) resolve() (spec.ValueOf[T], error) {
	env := os.Getenv(b.name)

	if env == "" {
		if v, ok := b.def.Get(); ok {
			return spec.ValueOf[T]{
				Go:    v,
				Env:   string(v),
				IsDef: true,
			}, nil
		}

		return spec.Undefined[T](b.name)
	}

	return spec.ValueOf[T]{
		Go:  T(env),
		Env: env,
	}, nil
}
