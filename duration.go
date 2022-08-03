package ferrite

import (
	"errors"
	"time"
)

// Duration configures an environment variable as a duration.
//
// name is the name of the environment variable to read. desc is a
// human-readable description of the environment variable.
func Duration(name, desc string) *DurationSpec {
	s := &DurationSpec{}
	s.init(s, name, desc)
	return s
}

// DurationSpec is the specification for a duration.
type DurationSpec struct {
	impl[time.Duration, *DurationSpec]
}

// parses parses and validates the value of the environment variable.
//
// validate() must be called on the result, as the parsed value does not
// necessarily meet all of the requirements.
func (s *DurationSpec) parse(value string) (time.Duration, error) {
	return time.ParseDuration(value)
}

// validate validates a parsed or default value.
func (s *DurationSpec) validate(value time.Duration) error {
	if value <= 0 {
		return errors.New("must be a positive duration")
	}

	return nil
}

// renderValidInput returns a string representation of the valid input values.
func (s *DurationSpec) renderValidInput() string {
	return "(1ns...)"
}

// renderParsed returns a string representation of the parsed value as it should
// appear in validation reports.
func (s *DurationSpec) renderParsed(value time.Duration) string {
	return value.String()
}

// renderRaw returns a string representation of the raw string value as it
// should appear in validation reports.
func (s *DurationSpec) renderRaw(value string) string {
	return value
}
