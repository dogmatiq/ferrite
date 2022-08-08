package ferrite

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dogmatiq/ferrite/internal/optional"
	"github.com/dogmatiq/ferrite/spec"
	"golang.org/x/exp/constraints"
)

// Signed configures an environment variable as a signed integer.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Signed[T constraints.Signed](name, desc string) SignedBuilder[T] {
	shift := bitSize[T]() - 1

	return SignedBuilder[T]{
		name: name,
		desc: desc,
		min:  -1 << shift,
		max:  (1 << shift) - 1,
	}
}

// SignedBuilder builds a specification for a signed integer value.
type SignedBuilder[T constraints.Signed] struct {
	name     string
	desc     string
	min, max T
	def      optional.Optional[T]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b SignedBuilder[T]) WithDefault(v T) SignedBuilder[T] {
	b.def = optional.With(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b SignedBuilder[T]) Required() Required[T] {
	return registerRequired(b.spec(), b.resolve)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b SignedBuilder[T]) Optional() Optional[T] {
	return registerOptional(b.spec(), b.resolve)
}

func (b SignedBuilder[T]) spec() spec.Spec {
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

func (b SignedBuilder[T]) resolve() (spec.ValueOf[T], error) {
	env := os.Getenv(b.name)

	if env == "" {
		if v, ok := b.def.Get(); ok {
			return spec.ValueOf[T]{
				Go:    v,
				Env:   b.render(v),
				IsDef: true,
			}, nil
		}

		return spec.Undefined[T](b.name)
	}

	n, err := strconv.ParseInt(env, 10, bitSize[T]())
	v := T(n)
	if err != nil || v < b.min || v > b.max {
		return spec.ValueOf[T]{}, fmt.Errorf(
			"%s must be an integer between %s and %s",
			b.name,
			b.render(b.min),
			b.render(b.max),
		)
	}

	return spec.ValueOf[T]{
		Go:  v,
		Env: env,
	}, nil
}

func (b SignedBuilder[T]) render(v T) string {
	return fmt.Sprintf("%+d", v)
}
