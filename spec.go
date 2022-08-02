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
	parse(value string, def *T) (T, ValidationResult)
}

// standard is the basis for a standard variable specification.
type standard[T any, F facade[T]] struct {
	facade F
	name   string
	desc   string

	done   uint32
	m      sync.Mutex
	def    *T
	value  T
	result ValidationResult
}

// init initializes the spec.
func (s *standard[T, F]) init(f F, name, desc string) {
	s.name = name
	s.desc = desc
	s.facade = f

	Register(name, s)
}

// WithDefault sets a default value to use when the environment variable is
// undefined.
func (s *standard[T, F]) WithDefault(v T) F {
	s.update(func() {
		s.def = &v
	})

	return s.facade
}

// Value returns the environment variable's value.
//
// It panics if the value is invalid.
func (s *standard[T, F]) Value() T {
	if res := s.Validate(s.name); res.Error != nil {
		panic(fmt.Sprintf("%s: %s", s.name, res.Error))
	}

	return s.value
}

// Validate validates the environment variable.
func (s *standard[T, F]) Validate(_ string) ValidationResult {
	if atomic.LoadUint32(&s.done) == 0 {
		s.m.Lock()
		defer s.m.Unlock()

		if s.done == 0 {
			value := os.Getenv(s.name)
			s.value, s.result = s.facade.parse(value, s.def)
		}
	}

	return s.result
}

// update calls fn while holding a lock on s.
//
// It panics if s has already been populated by parsing the environment
// variable.
func (s *standard[T, F]) update(fn func()) {
	if atomic.LoadUint32(&s.done) == 0 {
		s.m.Lock()
		defer s.m.Unlock()

		if s.done == 0 {
			fn()
			return
		}
	}

	panic("cannot modify spec after value has been used or validated")
}
