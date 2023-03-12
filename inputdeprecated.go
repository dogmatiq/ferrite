package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Deprecated is a specialization of the Input interface for values obtained
// from deprecated environment variables.
type Deprecated[T any] interface {
	Input

	// DeprecatedValue returns the parsed and validated value built from the
	// environment variable(s).
	//
	// If the constituent environment variable(s) are not defined and there is
	// no default value, ok is false; otherwise, ok is true and v is the value.
	//
	// It panics if any of one of the constituent environment variable(s) has an
	// invalid value.
	DeprecatedValue() (T, bool)
}

// DeprecatedOption is an option that can be applied to a deprecated variable.
type DeprecatedOption interface {
	applyDeprecatedOption(*deprecatedConfig)
}

// deprecatedConfig is the configuration for the deprecated inputs, built from
// DeprecatedOption values.
type deprecatedConfig struct {
	inputConfig
}

// buildDeprecatedConfig returns a new deprecatedConfig, built from the given
// options.
func buildDeprecatedConfig(
	spec variable.SpecBuilder,
	options ...DeprecatedOption,
) deprecatedConfig {
	var cfg deprecatedConfig
	cfg.Spec = spec

	for _, opt := range options {
		opt.applyDeprecatedOption(&cfg)
	}

	return cfg
}

// SupersededBy is a deprecation option that indicates i is a direct
// replacement for the deprecated variable(s).
func SupersededBy(inputs ...Input) DeprecatedOption {
	return option{
		Deprecated: func(cfg *deprecatedConfig) {
			for _, i := range inputs {
				for _, v := range i.variables() {
					variable.ApplyRelationship(
						variable.IsSupersededBy{
							Spec:         cfg.Spec.Peek(),
							SupersededBy: v.Spec(),
						},
					)
				}
			}
		},
	}
}

// deprecated is a convenience function that registers and returns a
// deprecated[T] that maps one-to-one with an environment variable of the same
// type.
func deprecated[T any, S variable.TypedSchema[T]](
	schema S,
	spec *variable.TypedSpecBuilder[T],
	options ...DeprecatedOption,
) Deprecated[T] {
	spec.MarkDeprecated()

	cfg := buildDeprecatedConfig(spec, options...)

	v := variable.Register(
		cfg.Registry,
		spec.Done(schema),
	)

	// interface is currently empty so we don't need an implementation
	return deprecatedFunc[T]{
		[]variable.Any{v},
		func() (T, bool, error) {
			return v.NativeValue()
		},
	}
}

// deprecatedFunc is an implementation of Deprecated[T] that obtains the value
// from an arbitrary function.
type deprecatedFunc[T any] struct {
	vars []variable.Any
	fn   func() (T, bool, error)
}

func (i deprecatedFunc[T]) DeprecatedValue() (T, bool) {
	n, ok, err := i.fn()
	if err != nil {
		panic(err.Error())
	}
	return n, ok
}

func (i deprecatedFunc[T]) variables() []variable.Any {
	return i.vars
}
