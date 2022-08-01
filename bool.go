package ferrite

import (
	"fmt"
	"os"
)

// Bool configures an environment variable as a boolean.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Bool(
	name, desc string,
	options ...SpecOption,
) *BoolSpec[bool] {
	return BoolAs[bool](name, desc, options...)
}

// BoolAs configures an environment variable as a boolean using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func BoolAs[T ~bool](
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

// BoolSpec is the specification for a boolean.
type BoolSpec[T ~bool] struct {
	spec[T]

	t, f string
}

// Literals overrides the default literals used to represent true and false.
//
// The default literals "true" and "false" are not considered valid values when
// using custom literals.
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

func (s *BoolSpec[T]) Describe() SpecDescription {
	var def string
	if s.def != nil {
		if *s.def {
			def = s.t
		} else {
			def = s.f
		}
	}

	return SpecDescription{
		s.desc,
		fmt.Sprintf("%s|%s", s.t, s.f),
		def,
	}
}

// Validate validates the environment variable.
func (s *BoolSpec[T]) Validate() (value string, isDefault bool, _ error) {
	switch os.Getenv(s.name) {
	case s.t:
		s.useValue(true)
		return s.t, false, nil

	case s.f:
		s.useValue(false)
		return s.f, false, nil

	case "":
		if v, err := s.useDefault(); err != nil {
			return "", false, err
		} else if v {
			return s.t, true, nil
		} else {
			return s.f, true, nil
		}

	default:
		return "", false, errNotInList(s.t, s.f)
	}
}
