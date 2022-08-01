package ferrite

import (
	"fmt"
	"os"
	"strings"
)

// Enum configures an environment variable as an enumeration with members of
// type T.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Enum[T any](
	name, desc string,
	options ...SpecOption,
) *EnumSpec[T] {
	s := &EnumSpec[T]{
		spec: spec[T]{
			name: name,
			desc: desc,
		},
	}

	register(s, options)

	return s
}

// EnumSpec is the specification for an enumeration.
type EnumSpec[T any] struct {
	spec[T]

	members []enumMember[T]
}

// enumKey returns the enum key to use for the given value.
func enumKey(v any) string {
	k := fmt.Sprint(v)
	if k == "" {
		panic("enum member must not have an empty string representation")
	}

	return k
}

// enumMember encapsulates an enum key and value.
type enumMember[T any] struct {
	Key   string
	Value T
}

// Members adds members to the enum.
//
// The environment variable must be set to the string representation of one of
// the member values. Members must not have an empty string representation.
func (s *EnumSpec[T]) Members(members ...T) *EnumSpec[T] {
	for _, v := range members {
		k := enumKey(v)

		s.members = append(s.members, enumMember[T]{k, v})
	}

	return s
}

// Default sets a default value to use when the environment variable is
// undefined.
//
// It must be one of the allowed members.
func (s *EnumSpec[T]) Default(v T) *EnumSpec[T] {
	for _, m := range s.members {
		k := enumKey(v)
		if k == m.Key {
			s.setDefault(v)
			return s
		}
	}

	panic("the default value must be one of the enum members")
}

// Validate validates the environment variable.
//
// It returns a string representation of the value.
func (s *EnumSpec[T]) Validate() VariableValidationResult {
	raw := os.Getenv(s.name)

	var keys []string
	valid := false
	for _, m := range s.members {
		keys = append(keys, m.Key)
		if raw == m.Key {
			valid = true
			s.useValue(m.Value)
		}
	}

	res := VariableValidationResult{
		Name:          s.name,
		Description:   s.desc,
		ValidInput:    strings.Join(keys, "|"),
		ExplicitValue: raw,
	}

	if s.def != nil {
		res.DefaultValue = enumKey(*s.def)
	}

	switch {
	case valid:
		return res
	case raw == "":
		res.Error = s.useDefault()
	default:
		res.Error = errNotInList(keys...)
	}

	return res
}
