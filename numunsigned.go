package ferrite

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

// Unsigned configures an environment variable as a unsigned integer.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Unsigned[T constraints.Unsigned](name, desc string) *UnsignedSpec[T] {
	s := &UnsignedSpec[T]{}
	s.init(s, name, desc)
	return s
}

// UnsignedSpec is the specification for a signed integer.
type UnsignedSpec[T constraints.Unsigned] struct {
	impl[T, *UnsignedSpec[T]]
}

// parses parses and validates the value of the environment variable.
//
// validate() must be called on the result, as the parsed value does not
// necessarily meet all of the requirements.
func (s *UnsignedSpec[T]) parse(value string) (T, error) {
	n, err := strconv.ParseUint(value, 10, bitSize[T]())
	return T(n), err
}

// validate validates a parsed or default value.
func (s *UnsignedSpec[T]) validate(value T) error {
	return nil
}

// renderValidInput returns a string representation of the valid input values.
func (s *UnsignedSpec[T]) renderValidInput() string {
	return inputType[T]()
}

// renderParsed returns a string representation of the parsed value as it should
// appear in validation reports.
func (s *UnsignedSpec[T]) renderParsed(value T) string {
	return strconv.FormatUint(uint64(value), 10)
}

// renderRaw returns a string representation of the raw string value as it
// should appear in validation reports.
func (s *UnsignedSpec[T]) renderRaw(value string) string {
	return value
}
