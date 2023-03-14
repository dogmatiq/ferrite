package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// Optional is a VariableSet used to obtain a value that may be unavailable, due
// to the environment variables not being defined.
type Optional[T any] interface {
	VariableSet

	// Value returns the parsed and validated value.
	//
	// It returns a non-nil error if any of one of the environment variables in
	// the set has an invalid value.
	//
	// If the environment variable(s) are not defined and there is no default
	// value, ok is false; otherwise, ok is true and v is the value.
	Value() (T, bool)
}

// OptionalOption is an option that configures an "optional" variable set. It
// may be passed to the Optional() method on each of the "builder" types.
type OptionalOption interface {
	applyOptionalOptionToConfig(*variableSetConfig)
	applyOptionalOptionToSpec(variable.SpecBuilder)
}

// required registers a variable that produces a value of type T and returns a
// Optional[T] that maps one-to-one to that variable.
func optional[T any, Schema variable.TypedSchema[T]](
	s Schema,
	b *variable.TypedSpecBuilder[T],
	options ...OptionalOption,
) Optional[T] {
	var cfg variableSetConfig
	for _, opt := range options {
		opt.applyOptionalOptionToConfig(&cfg)
		opt.applyOptionalOptionToSpec(b)
	}

	v := variable.Register(
		cfg.Registry,
		b.Done(s),
	)

	return optionalFunc[T]{
		[]variable.Any{v},
		func() (T, bool, error) {
			return v.NativeValue(),
				v.Availability() == variable.AvailabilityOK,
				v.Error()
		},
	}
}

// optionalFunc is an implementation of Optional[T] that obtains the value from
// an arbitrary function.
type optionalFunc[T any] struct {
	vars []variable.Any
	fn   func() (T, bool, error)
}

func (s optionalFunc[T]) Value() (T, bool) {
	n, ok, err := s.fn()
	if err != nil {
		panic(err.Error())
	}
	return n, ok
}

func (s optionalFunc[T]) variables() []variable.Any {
	return s.vars
}
