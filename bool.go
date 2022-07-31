package ferrite

import "fmt"

// Bool declares a boolean environment variable.
func Bool(name, desc string, options ...SpecOption) RequiredBool {
	opts := resolveSpecOptions(options)

	spec := &boolSpec{
		name: name,
		desc: desc,
		t:    "true",
		f:    "false",
	}

	opts.Registry.Register(spec)

	return RequiredBool{spec}
}

type boolSpec struct {
	name  string
	desc  string
	t, f  string
	def   *bool
	value bool
}

func (s *boolSpec) Name() string {
	return s.name
}

func (s *boolSpec) Resolve(env Environment) error {
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
	case s.f:
		s.value = false
	default:
		m := `ENVIRONMENT VARIABLES
 ✗ %s [bool] (%s)
   ✓ must be set explicitly
   ✗ must be either "%s" or "%s", got "%s"`
		return fmt.Errorf(m, s.name, s.desc, s.t, s.f, raw)
	}

	return nil
}

type RequiredBool struct {
	spec *boolSpec
}

func (b RequiredBool) Literals(t, f string) RequiredBool {
	b.spec.t = t
	b.spec.f = f
	return b
}

func (b RequiredBool) Default(def bool) RequiredBool {
	b.spec.def = &def
	return b
}

func (b RequiredBool) Optional() OptionalBool {
	def := false
	b.spec.def = &def
	return OptionalBool{b.spec}
}

func (b RequiredBool) Value() bool {
	return b.spec.value
}

type OptionalBool struct {
	spec *boolSpec
}

func (b OptionalBool) Value() (bool, bool) {
	return b.spec.value, false
}
