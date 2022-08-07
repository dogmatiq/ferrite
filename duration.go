package ferrite

import (
	"fmt"
	"os"
	"time"

	"github.com/dogmatiq/ferrite/internal/optional"
	"github.com/dogmatiq/ferrite/schema"
	"github.com/dogmatiq/ferrite/spec"
)

// Duration configures an environment variable as a duration.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Duration(name, desc string) DurationBuilder {
	return DurationBuilder{
		name: name,
		desc: desc,
	}
}

// DurationBuilder builds a specification for a duration variable.
type DurationBuilder struct {
	name string
	desc string
	def  optional.Optional[time.Duration]
}

// WithDefault sets a default value of the variable.
//
// It is used when the environment variable is undefined or empty.
func (b DurationBuilder) WithDefault(v time.Duration) DurationBuilder {
	b.def = optional.With(v)
	return b
}

// Required completes the build process and registers a required variable with
// Ferrite's validation system.
func (b DurationBuilder) Required() Required[time.Duration] {
	return registerRequired(b.spec(), b.resolve)
}

// Optional completes the build process and registers an optional variable with
// Ferrite's validation system.
func (b DurationBuilder) Optional() Optional[time.Duration] {
	return registerOptional(b.spec(), b.resolve)
}

func (b DurationBuilder) spec() spec.Spec {
	s := spec.Spec{
		Name:        b.name,
		Description: b.desc,
		Necessity:   spec.Required,
		Schema: schema.Range{
			Min: time.Nanosecond.String(),
		},
	}

	if v, ok := b.def.Get(); ok {
		s.Necessity = spec.Defaulted
		s.Default = v.String()
	}

	return s
}

func (b DurationBuilder) resolve() (spec.Value[time.Duration], error) {
	env := os.Getenv(b.name)

	if env == "" {
		if v, ok := b.def.Get(); ok {
			return spec.Value[time.Duration]{
				Go:        v,
				Env:       v.String(),
				IsDefault: true,
			}, nil
		}

		return spec.Value[time.Duration]{}, UndefinedError{Name: b.name}
	}

	v, err := time.ParseDuration(env)
	if err != nil {
		return spec.Value[time.Duration]{}, fmt.Errorf("%s is invalid: %w", b.name, err)
	}

	min := time.Nanosecond
	if v < min {
		return spec.Value[time.Duration]{}, fmt.Errorf(
			"%s must be %s or greater, got %s",
			b.name,
			min,
			v,
		)
	}

	return spec.Value[time.Duration]{
		Go:  v,
		Env: env,
	}, nil
}
