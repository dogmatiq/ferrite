package ferrite

import (
	"fmt"
)

// String configures an environment variable as a string.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func String(name, desc string) *StringSpec[string] {
	return StringAs[string](name, desc)
}

// StringAs configures an environment variable as a string using a user-defined
// type.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func StringAs[T ~string](name, desc string) *StringSpec[T] {
	s := &StringSpec[T]{}
	s.init(s, name, desc)
	return s
}

// StringSpec is the specification for a string.
type StringSpec[T ~string] struct {
	impl[T, *StringSpec[T]]
}

// parses parses and validates the value of the environment variable.
//
// validate() must be called on the result, as the parsed value does not
// necessarily meet all of the requirements.
func (s *StringSpec[T]) parse(value string) (T, error) {
	return T(value), nil
}

// validate validates a parsed or default value.
func (s *StringSpec[T]) validate(value T) error {
	return nil
}

// renderValidInput returns a string representation of the valid input values.
func (s *StringSpec[T]) renderValidInput() string {
	return inputType[T]()
}

// renderParsed returns a string representation of the parsed value as it should
// appear in validation reports.
func (s *StringSpec[T]) renderParsed(value T) string {
	return fmt.Sprintf("%q", value)
}

// renderRaw returns a string representation of the raw string value as it
// should appear in validation reports.
func (s *StringSpec[T]) renderRaw(value string) string {
	return fmt.Sprintf("%q", value)
}
