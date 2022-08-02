package ferrite

import (
	"fmt"
)

// Enum configures an environment variable as an enumeration with members of
// type T.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Enum[T any](name, desc string) *EnumSpec[T] {
	s := &EnumSpec[T]{}
	s.init(s, name, desc)
	return s
}

// EnumSpec is the specification for an enumeration.
type EnumSpec[T any] struct {
	impl[T, *EnumSpec[T]]
	members map[string]T
	order   []string
}

// WithMembers adds members to the enum.
//
// The environment variable must be set to the string representation of one of
// the member values. WithMembers must not have an empty string representation.
func (s *EnumSpec[T]) WithMembers(members ...T) *EnumSpec[T] {
	return s.with(func() {
		for _, v := range members {
			k := s.keyOf(v)

			if k == "" {
				panic("enum member must not have an empty string representation")
			}

			if _, ok := s.members[k]; ok {
				panic(fmt.Sprintf("enum already has a %q member", k))
			}

			if s.members == nil {
				s.members = map[string]T{}
			}

			s.members[k] = v
			s.order = append(s.order, k)
		}
	})
}

// parses parses and validates the value of the environment variable.
//
// validate() must be called on the result, as the parsed value does not
// necessarily meet all of the requirements.
func (s *EnumSpec[T]) parse(value string) (T, error) {
	if v, ok := s.members[value]; ok {
		return v, nil
	}

	var zero T
	return zero, fmt.Errorf("%s is not a member of the enum", value)
}

// validate validates a parsed or default value.
func (s *EnumSpec[T]) validate(value T) error {
	k := s.keyOf(value)

	if _, ok := s.members[k]; !ok {
		return fmt.Errorf("%s is not a member of the enum", k)
	}

	return nil
}

// renderValidInput returns a string representation of the valid input values.
func (s *EnumSpec[T]) renderValidInput() string {
	return inputList(s.order...)
}

// renderParsed returns a string representation of the parsed value as it should
// appear in validation reports.
func (s *EnumSpec[T]) renderParsed(value T) string {
	return s.keyOf(value)
}

// renderRaw returns a string representation of the raw string value as it
// should appear in validation reports.
func (s *EnumSpec[T]) renderRaw(value string) string {
	return value
}

// keyOf returns the key to use for the given enum member.
func (s *EnumSpec[T]) keyOf(member T) string {
	return fmt.Sprint(member)
}
