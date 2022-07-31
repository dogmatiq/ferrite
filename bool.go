package ferrite

import "fmt"

func Bool(name, desc string) RequiredBoolean {
	spec := &boolSpec{
		name: name,
		desc: desc,
		t:    "true",
		f:    "false",
	}

	DefaultRegistry.Register(spec)

	return RequiredBoolean{spec}
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

func (s *boolSpec) Resolve(lookup Lookup) error {
	raw, _ := lookup(s.name)

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

type RequiredBoolean struct {
	spec *boolSpec
}

func (b RequiredBoolean) Literals(t, f string) RequiredBoolean {
	b.spec.t = t
	b.spec.f = f
	return b
}

func (b RequiredBoolean) Default(def bool) RequiredBoolean {
	b.spec.def = &def
	return b
}

func (b RequiredBoolean) Optional() OptionalBoolean {
	def := false
	b.spec.def = &def
	return OptionalBoolean{b.spec}
}

func (b RequiredBoolean) Value() bool {
	return b.spec.value
}

type OptionalBoolean struct {
	spec *boolSpec
}

func (b OptionalBoolean) Value() (bool, bool) {
	return b.spec.value, false
}
