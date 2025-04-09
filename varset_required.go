package ferrite

import (
	"github.com/dogmatiq/ferrite/internal/variable"
)

// Required is a VariableSet used to obtain a value that must always be
// available, either from explicit environment variable values or by falling
// back to defaults.
type Required[T any] interface {
	VariableSet[T]

	// Value returns the parsed and validated value.
	//
	// It panics if any of one of the environment variables in the set is
	// undefined or has an invalid value.
	Value() T
}

// RequiredOption is an option that configures a "required" variable set. It may
// be passed to the Optional() method on each of the "builder" types.
type RequiredOption interface {
	applyRequiredOptionToConfig(*variableSetConfig)
	applyRequiredOptionToSpec(variable.SpecBuilder)
}

// required registers a variable that produces a value of type T and returns a
// Required[T] that maps one-to-one to that variable.
func required[T any, Schema variable.TypedSchema[T]](
	s Schema,
	b *variable.TypedSpecBuilder[T],
	options ...RequiredOption,
) Required[T] {
	b.MarkRequired()

	var cfg variableSetConfig
	for _, opt := range options {
		opt.applyRequiredOptionToConfig(&cfg)
		opt.applyRequiredOptionToSpec(b)
	}

	v := variable.Register(
		cfg.Registries,
		b.Done(s),
	)

	return requiredFunc[T]{
		[]variable.Any{v},
		func() (T, error) {
			return v.NativeValue(), v.Error()
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

// requiredFunc is an implementation of Required[T] that obtains the value from
// an arbitrary function.
type requiredFunc[T any] struct {
	vars       []variable.Any
	nativeFn   func() (T, error)
	literalsFn func(T) ([]variable.Literal, error)
}

func (s requiredFunc[T]) Value() T {
	n, err := s.nativeFn()
	if err != nil {
		panic(err.Error())
	}
	return n
}

func (s requiredFunc[T]) variables() []variable.Any {
	return s.vars
}

func (s requiredFunc[T]) native() (T, bool) {
	n, err := s.nativeFn()
	return n, err == nil
}

func (s requiredFunc[T]) literals(v T) ([]variable.Literal, error) {
	return s.literalsFn(v)
}
