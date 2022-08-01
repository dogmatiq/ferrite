package ferrite

import (
	"fmt"
)

// String configures an environment variable as a string.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func String(name, desc string) *StringSpec[string] {
	return StringAs[string](name, desc)
}

// StringAs configures an environment variable as a string using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func StringAs[T ~string](name, desc string) *StringSpec[T] {
	s := &StringSpec[T]{
		spec: spec[T]{
			name: name,
			desc: desc,
		},
	}

	Register(name, s)

	return s
}

// StringSpec is the specification for a string.
type StringSpec[T ~string] struct {
	spec[T]
}

// WithDefault sets a default value to use when the environment variable is
// undefined.
func (s *StringSpec[T]) WithDefault(v T) *StringSpec[T] {
	s.setDefault(v)
	return s
}

// Validate validates the environment variable.
func (s *StringSpec[T]) Validate(name, value string) ValidationResult {
	res := ValidationResult{
		Name:          name,
		Description:   s.desc,
		ValidInput:    fmt.Sprintf("[%T]", s.value),
		DefaultValue:  "",
		ExplicitValue: fmt.Sprintf("%q", value),
		Error:         nil,
	}

	if v, ok := s.Default(); ok {
		res.DefaultValue = fmt.Sprintf("%q", v)
	}

	if value != "" {
		s.useValue(T(value))
	} else if s.useDefault() {
		res.UsingDefault = true
	} else {
		res.Error = errUndefined
	}

	return res
}
