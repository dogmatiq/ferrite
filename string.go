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

func (s *StringSpec[T]) Describe() SpecDescription {
	var def string
	if s.def != nil {
		def = fmt.Sprintf("%q", *s.def)
	}

	return SpecDescription{
		s.desc,
		typeName[T](),
		def,
	}
}

// Validate validates the environment variable.
func (s *StringSpec[T]) Validate() (value string, isDefault bool, _ error) {
	raw := os.Getenv(s.name)

	if raw == "" {
		v, err := s.useDefault()
		if err != nil {
			return "", false, err
		}

		return fmt.Sprintf("%q", v), true, nil
	}

	s.useValue(T(raw))

	return fmt.Sprintf("%q", raw), false, nil
}
