package ferrite

import (
	"fmt"
	"os"
)

// String configures an environment variable as a string.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func String(
	name, desc string,
	options ...SpecOption,
) *StringSpec[string] {
	return StringAs[string](name, desc, options...)
}

// StringAs configures an environment variable as a string using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func StringAs[T ~string](
	name, desc string,
	options ...SpecOption,
) *StringSpec[T] {
	s := &StringSpec[T]{
		spec: spec[T]{
			name: name,
			desc: desc,
		},
	}

	register(s, options)

	return s
}

// StringSpec is the specification for a string.
type StringSpec[T ~string] struct {
	spec[T]
}

// Default sets a default value to use when the environment variable is
// undefined.
func (s *StringSpec[T]) Default(v T) *StringSpec[T] {
	s.setDefault(v)
	return s
}

// Validate validates the environment variable.
func (s *StringSpec[T]) Validate() VariableValidationResult {
	raw := os.Getenv(s.name)
	res := VariableValidationResult{
		Name:          s.name,
		Description:   s.desc,
		ValidInput:    typeName[T](),
		DefaultValue:  "",
		ExplicitValue: fmt.Sprintf("%q", raw),
		Error:         nil,
	}

	if s.def != nil {
		res.DefaultValue = fmt.Sprintf("%q", *s.def)
	}

	if raw == "" {
		res.UsingDefault = true
		res.Error = s.useDefault()
	} else {
		s.useValue(T(raw))
	}

	return res
}
