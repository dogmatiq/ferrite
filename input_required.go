package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Required is a specialization of the Input interface for values obtained
// from required (mandatory) environment variables.
type Required[T any] interface {
	Input

	// Value returns the parsed and validated value of the environment variable.
	//
	// It panics if any of one of the constituent environment variable(s) is
	// undefined or has an invalid value.
	Value() T
}

// RequiredOption is an option that can be applied to a required variable.
type RequiredOption interface {
	applyRequiredOption(*requiredConfig)
}

// requiredConfig is the configuration for the deprecated inputs, built from
// RequiredOption values.
type requiredConfig struct {
	inputConfig
}

// buildRequiredConfig returns a new requiredConfig, built from the given
// options.
func buildRequiredConfig(
	spec variable.SpecBuilder,
	options ...RequiredOption,
) requiredConfig {
	var cfg requiredConfig
	cfg.Spec = spec

	for _, opt := range options {
		opt.applyRequiredOption(&cfg)
	}

	return cfg
}

// required is a convenience function that registers and returns a required[T]
// that maps one-to-one with an environment variable of the same type.
func required[T any, S variable.TypedSchema[T]](
	schema S,
	spec *variable.TypedSpecBuilder[T],
	options ...RequiredOption,
) Required[T] {
	spec.MarkRequired()

	cfg := buildRequiredConfig(spec, options...)

	v := variable.Register(
		cfg.Registry,
		spec.Done(schema),
	)

	return requiredFunc[T]{
		[]variable.Any{v},
		func() (T, error) {
			n, ok, err := v.NativeValue()
			if ok || err != nil {
				return n, err
			}
			return n, undefinedError(v)
		},
	}
}

// requiredFunc is an implementation of Required[T] that obtains the value from
// an arbitrary function.
type requiredFunc[T any] struct {
	vars []variable.Any
	fn   func() (T, error)
}

func (i requiredFunc[T]) Value() T {
	n, err := i.fn()
	if err != nil {
		panic(err.Error())
	}
	return n
}

func (i requiredFunc[T]) variables() []variable.Any {
	return i.vars
}
