package ferrite

import (
	"fmt"
	"os"

	"github.com/dogmatiq/ferrite/schema"
)

// Spec is a specification for an environment variable.
type Spec interface {
	// Describe returns a description of the environment variable(s) described
	// by this spec.
	Describe() []VariableXXX

	// Validate validates the environment variable(s) described by this spec.
	Validate() []ValidationResult
}

// VariableXXX describes an environment variable.
type VariableXXX struct {
	// Name is the name of the environment variable.
	Name string

	// Description is a human-readable description of the environment variable.
	Description string

	// Schema describes the valid values for this environment variable.
	Schema schema.Schema

	// Default is the environment variable's default value.
	//
	// It must be non-empty if the environment variable has a default value;
	// otherwise it must be empty.
	Default string
}

// ValidationResult is the result of validating an environment variable.
type ValidationResult struct {
	// Name is the name of the environment variable.
	Name string

	// Value is the environment variable's value.
	Value string

	// UsedDefault is true if the default value of the environment variable was
	// used to populate the value.
	UsedDefault bool

	// Error is an error describing why the validation failed.
	//
	// If it is nil, the validation is considered successful.
	Error error
}

// impl is the basis for a impl variable specification.
//
// S is the concrete type of the specification.
type impl[T any, S specXXX[T]] struct {
	self S

	name string
	desc string

	m     smutex
	def   *T
	value *T
	err   error
}

// init initializes the spec.
func (s *impl[T, S]) init(self S, name, desc string) {
	s.self = self
	s.name = name
	s.desc = desc

	Register(s)
}

// WithDefault sets a default value to use when the environment variable is
// undefined.
func (s *impl[T, S]) WithDefault(v T) S {
	if err := s.self.validate(v); err != nil {
		panic(fmt.Sprintf(
			"default value of %s is invalid: %s",
			s.name,
			err,
		))
	}

	return s.with(func() {
		s.def = &v
		s.value = s.def
	})
}

// Value returns the environment variable's value.
//
// It panics if the value is invalid.
func (s *impl[T, S]) Value() T {
	s.resolve()

	if s.err != nil {
		panic(fmt.Sprintf("%s is invalid: %s", s.name, s.err))
	}

	return *s.value
}

// Describe returns a description of the environment variable(s) described by
// this spec.
func (s *impl[T, S]) Describe() []VariableXXX {
	s.m.RLock()
	defer s.m.RUnlock()

	def := ""
	if s.def != nil {
		def = s.self.renderParsed(*s.def)
	}

	return []VariableXXX{
		{
			Name:        s.name,
			Description: s.desc,
			Schema:      s.self.schema(),
			Default:     def,
		},
	}
}

// Validate validates the environment variable.
func (s *impl[T, S]) Validate() []ValidationResult {
	s.resolve()

	if s.err != nil {
		return []ValidationResult{
			{
				Name:  s.name,
				Error: s.err,
			},
		}
	}

	return []ValidationResult{
		{
			Name:        s.name,
			Value:       s.self.renderParsed(*s.value),
			UsedDefault: s.value == s.def, // address comparison, not value
		},
	}
}

// resolve populates s.value and s.result, or returns immediately if they are
// already populated.
func (s *impl[T, S]) resolve() {
	s.m.Seal(func() {
		value := os.Getenv(s.name)

		if value == "" {
			if s.def != nil {
				s.value = s.def
			} else {
				s.err = errUndefined
			}

			return
		}

		v, err := s.self.parse(value)
		if err != nil {
			s.err = err
			return
		}

		if err := s.self.validate(v); err != nil {
			s.err = err
			return
		}

		s.value = &v
	})
}

// with calls fn while holding a lock on s.
//
// It panics if the value has already been resolved.
func (s *impl[T, S]) with(fn func()) S {
	s.m.Lock()
	defer s.m.Unlock()

	fn()

	return s.self
}

// specXXX is a constraint for concrete implementations of a specXXX that embed
// impl[T].
type specXXX[T any] interface {
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
