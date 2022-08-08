package ferrite

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dogmatiq/ferrite/internal/optional"
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
		Schema: spec.Range{
			Min: time.Nanosecond.String(),
		},
	}

	if v, ok := b.def.Get(); ok {
		s.HasDefault = true
		s.DefaultX = v.String()
	}

	return s
}

func (b DurationBuilder) resolve() (spec.ValueOf[time.Duration], error) {
	env := os.Getenv(b.name)

	if env == "" {
		if v, ok := b.def.Get(); ok {
			return spec.ValueOf[time.Duration]{
				Go:    v,
				Env:   v.String(),
				IsDef: true,
			}, nil
		}

		return spec.Undefined[time.Duration](b.name)
	}

	v, err := time.ParseDuration(env)
	if err != nil {
		if strings.Contains(err.Error(), "unit") || strings.Contains(err.Error(), "unit") {
			return spec.Invalid[time.Duration](
				b.name,
				env,
				"%s",
				strings.Replace(
					strings.TrimPrefix(err.Error(), "time: "),
					fmt.Sprintf(` in duration %q`, env),
					"",
					1,
				),
			)
		}

		return spec.Invalid[time.Duration](
			b.name,
			env,
			"must be a valid duration, e.g. 10m30s",
		)
	}

	min := time.Nanosecond
	if v < min {
		return spec.Invalid[time.Duration](
			b.name,
			v.String(),
			"must be %s or greater",
			min,
		)
	}

	return spec.ValueOf[time.Duration]{
		Go:  v,
		Env: env,
	}, nil
}
