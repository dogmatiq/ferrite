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
func Enum[T any](name, desc string) *EnumSpec[T] {
	s := &EnumSpec[T]{
		spec: spec[T]{
			name: name,
			desc: desc,
		},
	}

	Register[T](s)

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

// WithMembers adds members to the enum.
//
// The environment variable must be set to the string representation of one of
// the member values. WithMembers must not have an empty string representation.
func (s *EnumSpec[T]) WithMembers(members ...T) *EnumSpec[T] {
	for _, v := range members {
		k := enumKey(v)

		s.members = append(s.members, enumMember[T]{k, v})
	}

	return s
}

// WithDefault sets a default value to use when the environment variable is
// undefined.
//
// It must be one of the allowed members.
func (s *EnumSpec[T]) WithDefault(v T) *EnumSpec[T] {
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

	if v, ok := s.Default(); ok {
		res.DefaultValue = enumKey(v)
	}

	if valid {
		// nothing more to do
	} else if raw != "" {
		res.Error = fmt.Errorf("%s is not a member of the enum", raw)
	} else if s.useDefault() {
		res.UsingDefault = true
	} else {
		res.Error = errUndefined
	}

	return res
}
