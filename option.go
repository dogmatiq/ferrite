package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// An Option changes the behavior of an environment variable specification.
//
// WARNING: The signature of this function is not considered part of Ferrite's
// public API. It may change at any time and without warning.
type Option func(variable.SpecBuilder) []variable.RegisterOption

func applyOptions(
	spec variable.SpecBuilder,
	opts []Option,
) []variable.RegisterOption {
	var registerOptions []variable.RegisterOption

	for _, opt := range opts {
		registerOptions = append(
			registerOptions,
			opt(spec)...,
		)
	}

	return registerOptions
}

// Sensitive is an Option that marks a variable as containing sensitive
// information.
func Sensitive() Option {
	return func(spec variable.SpecBuilder) []variable.RegisterOption {
		spec.MarkSensitive()
		return nil
	}
}
