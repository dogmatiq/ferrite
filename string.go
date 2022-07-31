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
func (s *StringSpec[T]) Validate() error {
	raw := os.Getenv(s.name)

	if raw == "" {
		if s.useDefault() {
			return nil
		}

		m := `ENVIRONMENT VARIABLES
 ✗ %s [%s] (%s)
   ✗ must be set explicitly`
		return fmt.Errorf(
			m,
			s.name,
			typeName[T](),
			s.desc,
		)
	}

	s.useValue(T(raw))

	return nil
}
