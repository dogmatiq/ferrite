package ferrite

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite/schema"
)

// impl is the basis for a impl variable specification.
//
// S is the concrete type of the specification.
type impl[T any, S spec[T]] struct {
	self S

	seal      seal
	defaulted bool
	value     T
	result    ValidationResult
}

// init initializes the spec.
func (s *impl[T, S]) init(self S, name, desc string) {
	s.self = self
	s.result.Name = name
	s.result.Description = desc

	Register(s)
}

// WithDefault sets a default value to use when the environment variable is
// undefined.
func (s *impl[T, S]) WithDefault(v T) S {
	if err := s.self.validate(v); err != nil {
		panic(fmt.Sprintf(
			"default value of %s is invalid: %s",
			s.result.Name,
			err,
		))
	}

	return s.with(func() {
		s.defaulted = true
		s.value = v
		s.result.DefaultValue = s.self.renderParsed(v)
	})
}

// Value returns the environment variable's value.
//
// It panics if the value is invalid.
func (s *impl[T, S]) Value() T {
	s.resolve()

	if s.result.Error != nil {
		panic(fmt.Sprintf(
			"%s is invalid: %s",
			s.result.Name,
			s.result.Error,
		))
	}

	return s.value
}

// Validate validates the environment variable.
func (s *impl[T, S]) Validate() []ValidationResult {
	s.resolve()
	return []ValidationResult{s.result}
}

// resolve populates s.value and s.result, or returns immediately if they are
// already populated.
func (s *impl[T, S]) resolve() {
	s.seal.Close(func() {
		s.result.Schema = s.self.schema()
		value := os.Getenv(s.result.Name)

		if value == "" {
			if s.defaulted {
				s.result.UsingDefault = true
			} else {
				s.result.Error = errUndefined
			}

			return
		}

		s.result.ExplicitValue = s.self.renderRaw(value)

		v, err := s.self.parse(value)
		if err != nil {
			s.result.Error = err
			return
		}

		if err := s.self.validate(v); err != nil {
			s.result.Error = err
			return
		}

		s.value = v
		s.result.ExplicitValue = s.self.renderParsed(v)
	})
}

// with calls fn while holding a lock on s.
//
// It panics if the value has already been resolved.
func (s *impl[T, S]) with(fn func()) S {
	s.seal.Do(fn)
	return s.self
}

// spec is a constraint for concrete implementations of a spec that embed
// impl[T].
type spec[T any] interface {
	// parses parses and validates the value of the environment variable.
	//
	// validate() must be called on the result, as the parsed value does not
	// necessarily meet all of the requirements.
	parse(value string) (T, error)

	// validate validates a parsed or default value.
	validate(value T) error

	// schema returns the schema that describes the environment variable's
	// valid values.
	schema() schema.Schema

	// renderParsed returns a string representation of the parsed value as it
	// should appear in validation reports.
	renderParsed(value T) string

	// renderRaw returns a string representation of the raw string value as it
	// should appear in validation reports.
	renderRaw(value string) string
}
