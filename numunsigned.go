package ferrite

import (
	"fmt"
	"strconv"

	"github.com/dogmatiq/ferrite/schema"
	"golang.org/x/exp/constraints"
)

// Unsigned configures an environment variable as a unsigned integer.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Unsigned[T constraints.Unsigned](name, desc string) *UnsignedSpec[T] {
	s := &UnsignedSpec[T]{
		min: 0,
		max: T(0) - 1,
	}

	s.init(s, name, desc)
	return s
}

// UnsignedSpec is the specification for a signed integer.
type UnsignedSpec[T constraints.Unsigned] struct {
	impl[T, *UnsignedSpec[T]]

	min, max T
}

// parses parses and validates the value of the environment variable.
//
// validate() must be called on the result, as the parsed value does not
// necessarily meet all of the requirements.
func (s *UnsignedSpec[T]) parse(value string) (T, error) {
	n, err := strconv.ParseUint(value, 10, bitSize[T]())
	v := T(n)

	if err != nil || v < s.min || v > s.max {
		return 0, fmt.Errorf(
			"must be an integer between %d and %d",
			s.min,
			s.max,
		)
	}

	return v, err
}

// validate validates a parsed or default value.
func (s *UnsignedSpec[T]) validate(value T) error {
	if value < s.min || value > s.max {
		return fmt.Errorf(
			"must be an integer between %d and %d",
			s.min,
			s.max,
		)
	}

	return nil
}

// schema returns the schema that describes the environment variable's
// valid values.
func (s *UnsignedSpec[T]) schema() schema.Schema {
	return schema.Range{
		Min: s.renderParsed(s.min),
		Max: s.renderParsed(s.max),
	}
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
