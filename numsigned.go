package ferrite

import (
	"fmt"
	"strconv"

	"github.com/dogmatiq/ferrite/schema"
	"golang.org/x/exp/constraints"
)

// Signed configures an environment variable as a signed integer.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Signed[T constraints.Signed](name, desc string) *SignedSpec[T] {
	shift := bitSize[T]() - 1

	s := &SignedSpec[T]{
		min: -1 << shift,
		max: (1 << shift) - 1,
	}

	s.init(s, name, desc)
	return s
}

// SignedSpec is the specification for a signed integer.
type SignedSpec[T constraints.Signed] struct {
	impl[T, *SignedSpec[T]]

	min, max T
}

// parses parses and validates the value of the environment variable.
//
// validate() must be called on the result, as the parsed value does not
// necessarily meet all of the requirements.
func (s *SignedSpec[T]) parse(value string) (T, error) {
	n, err := strconv.ParseInt(value, 10, bitSize[T]())
	v := T(n)

	if err != nil {
		return 0, fmt.Errorf(
			"must be an integer between %+d and %+d",
			s.min,
			s.max,
		)
	}

	return v, err
}

// validate validates a parsed or default value.
func (s *SignedSpec[T]) validate(value T) error {
	if value < s.min || value > s.max {
		return fmt.Errorf(
			"must be an integer between %+d and %+d",
			s.min,
			s.max,
		)
	}

	return nil
}

// schema returns the schema that describes the environment variable's
// valid values.
func (s *SignedSpec[T]) schema() schema.Schema {
	return schema.Range{
		Min: s.renderParsed(s.min),
		Max: s.renderParsed(s.max),
	}
}

// renderParsed returns a string representation of the parsed value as it should
// appear in validation reports.
func (s *SignedSpec[T]) renderParsed(value T) string {
	return fmt.Sprintf("%+d", value)
}

// renderRaw returns a string representation of the raw string value as it
// should appear in validation reports.
func (s *SignedSpec[T]) renderRaw(value string) string {
	return value
}
