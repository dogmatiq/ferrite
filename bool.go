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

// Validate validates the environment variable.
func (s *BoolSpec[T]) Validate() VariableValidationResult {
	raw := os.Getenv(s.name)
	res := VariableValidationResult{
		Name:          s.name,
		Description:   s.desc,
		ValidInput:    fmt.Sprintf("%s|%s", s.t, s.f),
		ExplicitValue: raw,
	}

	if s.def != nil {
		if *s.def {
			res.DefaultValue = s.t
		} else {
			res.DefaultValue = s.f
		}
	}

	switch raw {
	case s.t:
		s.useValue(true)
	case s.f:
		s.useValue(false)
	case "":
		res.Error = s.useDefault()
	default:
		res.Error = errNotInList(s.t, s.f)
	}

	return res
}
