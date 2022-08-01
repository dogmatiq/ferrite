package ferrite

// Spec is a specification for an environment variable.
//
// It describes the environment variable itself, and how to construct valid
// values for the variable.
type Spec interface {
	// Names returns the names of the environment variables constrained by this
	// spec.
	Names() []string

	// Validate validates the environment variable.
	Validate(name string) ValidationResult
}

// SpecFor is a specification for an environment variable that produces values
// of type T.
type SpecFor[T any] interface {
	Spec

	// Value returns the value of the environment variable.
	//
	// It panics if the variable is invalid.
	Value() T
}

// spec provides common functionality for Spec implementations.
type spec[T any] struct {
	name string
	desc string

	isValidated  bool
	hasDefault   bool
	defaultValue T
	value        T
}

func (s *spec[T]) Names() []string {
	return []string{s.name}
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
