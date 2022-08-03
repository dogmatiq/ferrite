package ferrite

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

// Signed configures an environment variable as a signed integer.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Signed[T constraints.Signed](name, desc string) *SignedSpec[T] {
	s := &SignedSpec[T]{}
	s.init(s, name, desc)
	return s
}

// SignedSpec is the specification for a signed integer.
type SignedSpec[T constraints.Signed] struct {
	impl[T, *SignedSpec[T]]
}

// parses parses and validates the value of the environment variable.
//
// validate() must be called on the result, as the parsed value does not
// necessarily meet all of the requirements.
func (s *SignedSpec[T]) parse(value string) (T, error) {
	n, err := strconv.ParseInt(value, 10, bitSize[T]())
	return T(n), err
}

// validate validates a parsed or default value.
func (s *SignedSpec[T]) validate(value T) error {
	return nil
}

// renderValidInput returns a string representation of the valid input values.
func (s *SignedSpec[T]) renderValidInput() string {
	return inputType[T]()
}

// renderParsed returns a string representation of the parsed value as it should
// appear in validation reports.
func (s *SignedSpec[T]) renderParsed(value T) string {
	return strconv.FormatInt(int64(value), 10)
}

// renderRaw returns a string representation of the raw string value as it
// should appear in validation reports.
func (s *SignedSpec[T]) renderRaw(value string) string {
	return value
}
