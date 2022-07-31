package ferrite

import (
	"fmt"
	"os"
)

// Bool declares a boolean environment variable.
func Bool(name, desc string, options ...SpecOption) *BoolSpec {
	s := &BoolSpec{
		spec: spec[bool]{
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
type BoolSpec struct {
	spec[bool]

	t, f string
}

// Literals sets a pair of custom stirng literals used to represent true and
// false. The default literals are "true" and "false".
func (s *BoolSpec) Literals(t, f string) *BoolSpec {
	s.t = t
	s.f = f
	return s
}

// Default sets a default value to use when the environment variable is
// undefined.
func (s *BoolSpec) Default(v bool) *BoolSpec {
	s.setDefault(v)
	return s
}

// Resolve resolves the value of the environment variable from the environment.
func (s *BoolSpec) Resolve() error {
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
