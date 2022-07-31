package ferrite

import (
	"fmt"
	"os"
)

// Bool declares a boolean environment variable.
func Bool(
	name, desc string,
	options ...SpecOption,
) *BoolSpec[bool] {
	return BoolKind[bool](name, desc, options...)
}

// BoolKind represents a boolean environment variable as some user-defined type.
func BoolKind[T ~bool](
	name, desc string,
	options ...SpecOption,
) *BoolSpec[T] {
	s := &BoolSpec[T]{
		spec: spec[T]{
			name: name,
			desc: desc,
		},
		t: "true",
		f: "false",
	}

	register(s, options)

	return s
}

// BoolSpec is a Spec for boolean types.
type BoolSpec[T ~bool] struct {
	spec[T]

	t, f string
}

// Literals sets a pair of custom stirng literals used to represent true and
// false. The default literals are "true" and "false".
func (s *BoolSpec[T]) Literals(t, f string) *BoolSpec[T] {
	s.t = t
	s.f = f
	return s
}

// Default sets a default value to use when the environment variable is
// undefined.
func (s *BoolSpec[T]) Default(v T) *BoolSpec[T] {
	s.setDefault(v)
	return s
}

// Resolve resolves the value of the environment variable from the environment.
func (s *BoolSpec[T]) Resolve() error {
	raw := os.Getenv(s.name)

	if raw == "" {
		if s.useDefault() {
			return nil
		}

		m := `ENVIRONMENT VARIABLES
 ✗ %s [bool] (%s)
   ✗ must be set explicitly
   ✗ must be either "%s" or "%s"`
		return fmt.Errorf(m, s.name, s.desc, s.t, s.f)
	}

	switch raw {
	case s.t:
		s.useValue(true)
	case s.f:
		s.useValue(false)
	default:
		m := `ENVIRONMENT VARIABLES
 ✗ %s [bool] (%s)
   ✓ must be set explicitly
   ✗ must be either "%s" or "%s", got "%s"`
		return fmt.Errorf(m, s.name, s.desc, s.t, s.f, raw)
	}

	return nil
}
