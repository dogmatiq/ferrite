package ferrite

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dogmatiq/ferrite/internal/optional"
	"github.com/dogmatiq/ferrite/spec"
	"golang.org/x/exp/constraints"
)

// Unsigned configures an environment variable as a unsigned integer.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Unsigned[T constraints.Unsigned](name, desc string) UnsignedBuilder[T] {
	return UnsignedBuilder[T]{
		name: name,
		desc: desc,
		min:  0,
		max:  T(0) - 1,
	}
}

// UnsignedBuilder builds a specification for an unsigned integer value.
type UnsignedBuilder[T constraints.Unsigned] struct {
	name     string
	desc     string
	min, max T
	def      optional.Optional[T]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b UnsignedBuilder[T]) WithDefault(v T) UnsignedBuilder[T] {
	b.def = optional.With(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b UnsignedBuilder[T]) Required() Required[T] {
	return registerRequired(b.spec(), b.resolve)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b UnsignedBuilder[T]) Optional() Optional[T] {
	return registerOptional(b.spec(), b.resolve)
}

func (b UnsignedBuilder[T]) spec() spec.Spec {
	s := spec.Spec{
		Name:        b.name,
		Description: b.desc,
		Necessity:   spec.Required,
		Schema: spec.Range{
			Min: b.render(b.min),
			Max: b.render(b.max),
		},
	}

	if v, ok := b.def.Get(); ok {
		s.Necessity = spec.Defaulted
		s.Default = b.render(v)
	}

	return s
}

func (b UnsignedBuilder[T]) resolve() (spec.Value[T], error) {
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

	n, err := strconv.ParseUint(env, 10, bitSize[T]())
	v := T(n)
	if err != nil || v < b.min || v > b.max {
		return spec.Value[T]{}, fmt.Errorf(
			"%s must be an integer between %s and %s",
			b.name,
			b.render(b.min),
			b.render(b.max),
		)
	}

	return spec.Value[T]{
		Go:  v,
		Env: env,
	}, nil
}

func (b UnsignedBuilder[T]) render(v T) string {
	return fmt.Sprintf("%d", v)
}
