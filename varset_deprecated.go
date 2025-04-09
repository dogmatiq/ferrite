package ferrite

import (
	"github.com/dogmatiq/ferrite/internal/variable"
)

// Deprecated is a VariableSet used to obtain a value that may be unavailable,
// due to the environment variables not being defined.
type Deprecated[T any] interface {
	VariableSet[T]

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
		cfg.Registries,
		b.Done(s),
	)

	return deprecatedFunc[T]{
		[]variable.Any{v},
		func() (T, bool, error) {
			return v.NativeValue(),
				v.Availability() == variable.AvailabilityOK,
				v.Error()
		},
		func(x T) ([]variable.Literal, error) {
			lit, err := v.TypedSpec.Marshal(variable.ConstraintContextExample, x)
			if err != nil {
				return nil, err
			}
			return []variable.Literal{lit}, nil
		},
	}
}

// deprecatedFunc is an implementation of Deprecated[T] that obtains the value
// from an arbitrary function.
type deprecatedFunc[T any] struct {
	vars       []variable.Any
	nativeFn   func() (T, bool, error)
	literalsFn func(T) ([]variable.Literal, error)
}

func (s deprecatedFunc[T]) DeprecatedValue() (T, bool) {
	n, ok, err := s.nativeFn()
	if err != nil {
		panic(err.Error())
	}
	return n, ok
}

func (s deprecatedFunc[T]) variables() []variable.Any {
	return s.vars
}

func (s deprecatedFunc[T]) native() (T, bool) {
	n, ok, err := s.nativeFn()
	return n, ok && err == nil
}

func (s deprecatedFunc[T]) literals(v T) ([]variable.Literal, error) {
	return s.literalsFn(v)
}
