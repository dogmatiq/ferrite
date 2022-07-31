package ferrite

import "fmt"

// Bool declares a boolean environment variable.
func Bool(name, desc string, options ...SpecOption) *BoolSpec {
	return register(
		&BoolSpec{
			name: name,
			desc: desc,
			t:    "true",
			f:    "false",
		},
		options,
	)
}

// BoolSpec is a Spec for boolean types.
type BoolSpec struct {
	name     string
	desc     string
	t, f     string
	def      *bool
	explicit bool
	value    bool
}

func (s *BoolSpec) Name() string {
	return s.name
}

func (s *BoolSpec) Literals(t, f string) *BoolSpec {
	s.t = t
	s.f = f
	return s
}

func (s *BoolSpec) Default(def bool) *BoolSpec {
	s.def = &def
	return s
}

func (s *BoolSpec) Optional() *BoolSpec {
	return s.Default(false)
}

func (s *BoolSpec) Value() bool {
	return s.value
}

func (s *BoolSpec) IsExplicit() bool {
	return s.explicit
}

func (s *BoolSpec) Resolve(env Environment) error {
	raw, _ := env.Lookup(s.name)

	if raw == "" {
		if s.def != nil {
			s.value = *s.def
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
		s.value = true
		s.explicit = true
	case s.f:
		s.value = false
		s.explicit = true
	default:
		m := `ENVIRONMENT VARIABLES
 ✗ %s [bool] (%s)
   ✓ must be set explicitly
   ✗ must be either "%s" or "%s", got "%s"`
		return fmt.Errorf(m, s.name, s.desc, s.t, s.f, raw)
	}

	return nil
}
