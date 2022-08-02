package ferrite

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

type spec[T any] struct {
	name string
	desc string

	isValidated  bool
	hasDefault   bool
	defaultValue T
	value        T
}

func (s *spec[T]) Default() (T, bool) {
	return s.defaultValue, s.hasDefault
}

func (s *spec[T]) Value() T {
	if !s.isValidated {
		panic("environment has not been validated")
	}

	return s.value
}

func (s *spec[T]) setDefault(v T) {
	s.hasDefault = true
	s.defaultValue = v
}

func (s *spec[T]) useValue(v T) {
	s.isValidated = true
	s.value = v
}

func (s *spec[T]) useDefault() bool {
	if s.hasDefault {
		s.isValidated = true
		s.value = s.defaultValue
		return true
	}

	return false
}

// facade is a constraint for specifications that parse environment variables
// into values of type T.
type facade[T any] interface {
	// parses parses and validates the value of the environment variable.
	parse(value string) (T, error)

	// validate validates a parsed or default value.
	validate(value T) error

	// renderParsed returns a string representation of the parsed value as it should
	// appear in validation reports.
	renderParsed(value T) string

	// renderRaw returns a string representation of the raw string value as it
	// should appear in validation reports.
	renderRaw(value string) string
}

// standard is the basis for a standard variable specification.
type standard[T any, F facade[T]] struct {
	facade F

	done      uint32
	m         sync.Mutex
	defaulted bool
	value     T
	result    ValidationResult
}

// init initializes the spec.
func (s *standard[T, F]) init(f F, name, desc string) {
	s.facade = f
	s.result.Name = name
	s.result.Description = desc
	s.result.ValidInput = fmt.Sprintf("[%T]", s.value)

	Register(name, s)
}

// WithDefault sets a default value to use when the environment variable is
// undefined.
func (s *standard[T, F]) WithDefault(v T) F {
	if err := s.facade.validate(v); err != nil {
		panic(fmt.Sprintf(
			"default value of %s is invalid: %s",
			s.result.Name,
			err,
		))
	}

	return s.update(func() {
		s.defaulted = true
		s.value = v
		s.result.DefaultValue = s.facade.renderParsed(v)
	})
}

// Value returns the environment variable's value.
//
// It panics if the value is invalid.
func (s *standard[T, F]) Value() T {
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
func (s *standard[T, F]) Validate() ValidationResult {
	s.resolve()
	return s.result
}

// resolve populates s.value and s.result.
func (s *standard[T, F]) resolve() {
	if atomic.LoadUint32(&s.done) != 0 {
		return
	}

	s.m.Lock()
	defer s.m.Unlock()

	if s.done != 0 {
		return
	}

	value := os.Getenv(s.result.Name)

	if value == "" {
		if s.defaulted {
			s.result.UsingDefault = true
		} else {
			s.result.Error = errUndefined
		}

		return
	}

	s.result.ExplicitValue = s.facade.renderRaw(value)

	v, err := s.facade.parse(value)
	if err != nil {
		s.result.Error = err
		return
	}

	if err := s.facade.validate(v); err != nil {
		s.result.Error = err
		return
	}

	s.value = v
	s.result.ExplicitValue = s.facade.renderParsed(v)
}

// update calls fn while holding a lock on s.
//
// It panics if s has already been populated by parsing the environment
// variable.
func (s *standard[T, F]) update(fn func()) F {
	if atomic.LoadUint32(&s.done) == 0 {
		s.m.Lock()
		defer s.m.Unlock()

		if s.done == 0 {
			fn()
			return s.facade
		}
	}

	panic("cannot modify spec after value has been used or validated")
}
