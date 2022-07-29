package ferrite

import "errors"

func Bool(name, desc string) RequiredBoolean {
	spec := &boolSpec{
		name: name,
	}

	DefaultRegistry.Register(spec)

	return RequiredBoolean{spec}
}

type boolSpec struct {
	name  string
	def   *bool
	value bool
}

func (s *boolSpec) Name() string {
	return s.name
}

func (s *boolSpec) Resolve(lookup LookupFn) error {
	raw, _ := lookup(s.name)

	if raw == "" {
		if s.def != nil {
			s.value = *s.def
			return nil
		}

		m := `
ENVIRONMENT VARIABLES
	✗ FERRITE_DEBUG [bool] (enable debug logging)
		✗ must be set explicitly
		✗ must be set to "true" or "false"
`
		return errors.New(m)
	}

	s.value = raw == "true"

	return nil
}

type RequiredBoolean struct {
	spec *boolSpec
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
