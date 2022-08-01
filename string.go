package ferrite

import (
	"fmt"
	"os"
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

	Register[T](s)

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
func (s *StringSpec[T]) Validate(name string) ValidationResult {
	raw := os.Getenv(name)
	res := ValidationResult{
		Name:          name,
		Description:   s.desc,
		ValidInput:    fmt.Sprintf("[%T]", s.value),
		DefaultValue:  "",
		ExplicitValue: fmt.Sprintf("%q", raw),
		Error:         nil,
	}

	if v, ok := s.Default(); ok {
		res.DefaultValue = fmt.Sprintf("%q", v)
	}

	if raw != "" {
		s.useValue(T(raw))
	} else if s.useDefault() {
		res.UsingDefault = true
	} else {
		res.Error = errUndefined
	}

	return res
}
