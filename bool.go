package ferrite

import (
	"fmt"
)

// Bool configures an environment variable as a boolean.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Bool(name, desc string) *BoolSpec[bool] {
	return BoolAs[bool](name, desc)
}

// BoolAs configures an environment variable as a boolean using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func BoolAs[T ~bool](name, desc string) *BoolSpec[T] {
	s := &BoolSpec[T]{}
	s.init(s, name, desc)
	return s.WithLiterals("true", "false")
}

// BoolSpec is the specification for a boolean.
type BoolSpec[T ~bool] struct {
	impl[T, *BoolSpec[T]]
	t, f string
}

// WithLiterals overrides the default literals used to represent true and false.
//
// The default literals "true" and "false" are not considered valid values when
// using custom literals.
func (s *BoolSpec[T]) WithLiterals(t, f string) *BoolSpec[T] {
	if t == "" || f == "" {
		panic("boolean literals must not be zero-length")
	}

	return s.with(func() {
		s.t = t
		s.f = f
		s.result.ValidInput = fmt.Sprintf("%s|%s", t, f)
	})
}

// parses parses and validates the value of the environment variable.
//
// validate() must be called on the result, as the parsed value does not
// necessarily meet all of the requirements.
func (s *BoolSpec[T]) parse(value string) (T, error) {
	switch value {
	case s.t:
		return true, nil
	case s.f:
		return false, nil
	default:
		return false, fmt.Errorf("must be either %q or %q", s.t, s.f)
	}
}

// validate validates a parsed or default value.
func (s *BoolSpec[T]) validate(value T) error {
	return nil
}

// renderValidInput returns a string representation of the valid input values.
func (s *BoolSpec[T]) renderValidInput() string {
	return inputList(s.t, s.f)
}

// renderParsed returns a string representation of the parsed value as it should
// appear in validation reports.
func (s *BoolSpec[T]) renderParsed(value T) string {
	if value {
		return s.t
	}
	return s.f
}

// renderRaw returns a string representation of the raw string value as it
// should appear in validation reports.
func (s *BoolSpec[T]) renderRaw(value string) string {
	return value
}
