package ferrite

import (
	"fmt"
	"os"
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

	panic("default value must be one of the enum members")
}

// Validate validates the environment variable.
func (s *EnumSpec[T]) Validate() error {
	raw := os.Getenv(s.name)

	if raw == "" {
		if s.useDefault() {
			return nil
		}

		f := `ENVIRONMENT VARIABLES
 ✗ %s [%T enum] (%s)
   ✗ must be set explicitly
   • must be one of the enum members`

		for _, m := range s.members {
			f += fmt.Sprintf("\n     • %s", m.Key)
		}

		return fmt.Errorf(
			f,
			s.name,
			s.value,
			s.desc,
		)
	}

	for _, m := range s.members {
		if raw == m.Key {
			s.useValue(m.Value)
			return nil
		}
	}

	f := `ENVIRONMENT VARIABLES
 ✗ %s [%T enum] (%s)
   ✓ must be set explicitly
   ✗ must be one of the enum members, got %q`

	for _, m := range s.members {
		f += fmt.Sprintf("\n     • %s", m.Key)
	}

	return fmt.Errorf(
		f,
		s.name,
		s.value,
		s.desc,
		raw,
	)
}
