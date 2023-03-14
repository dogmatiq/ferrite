package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Deprecated is a VariableSet used to obtain a value that may be unavailable,
// due to the environment variables not being defined.
type Deprecated[T any] interface {
	VariableSet

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

// DeprecatedOption is an option that configures a "deprecated" variable set. It
// may be passed to the Deprecated() method on each of the "builder" types.
type DeprecatedOption interface {
	applyDeprecatedOptionToConfig(*variableSetConfig)
	applyDeprecatedOptionToSpec(variable.SpecBuilder)
}

// deprecated registers a variable that produces a value of type T and returns a
// Deprecated[T] that maps one-to-one to that variable.
func deprecated[T any, Schema variable.TypedSchema[T]](
	s Schema,
	b *variable.TypedSpecBuilder[T],
	options ...DeprecatedOption,
) Deprecated[T] {
	b.MarkDeprecated()

	var cfg variableSetConfig
	for _, opt := range options {
		opt.applyDeprecatedOptionToConfig(&cfg)
		opt.applyDeprecatedOptionToSpec(b)
	}

	v := variable.Register(
		cfg.Registry,
		b.Done(s),
	)

	return deprecatedFunc[T]{
		[]variable.Any{v},
		func() (T, bool, error) {
			return v.NativeValue(),
				v.Availability() == variable.AvailabilityOK,
				v.Error()
		},
	}
}

// deprecatedFunc is an implementation of Deprecated[T] that obtains the value
// from an arbitrary function.
type deprecatedFunc[T any] struct {
	vars []variable.Any
	fn   func() (T, bool, error)
}

func (s deprecatedFunc[T]) DeprecatedValue() (T, bool) {
	n, ok, err := s.fn()
	if err != nil {
		panic(err.Error())
	}
	return n, ok
}

func (s deprecatedFunc[T]) value() any {
	if n, ok, _ := s.fn(); ok {
		return n
	}
	return nil
}

func (s deprecatedFunc[T]) variables() []variable.Any {
	return s.vars
}
