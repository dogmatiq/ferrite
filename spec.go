package ferrite

// Spec is a specification for an environment variable.
//
// It describes the environment variable itself, and how to construct valid
// values for the variable.
type Spec interface {
	// Name returns the name of the environment variable.
	Name() string

	// Describe human-readable descriptions of the variable's behavior.
	Describe() SpecDescription

	// Validate validates the environment variable.
	//
	// It returns a string representation of the value.
	Validate() (value string, isDefault bool, _ error)
}

// SpecDescription encapsulates descriptions of various aspects of a Spec.
type SpecDescription struct {
	Variable string
	Input    string
	Default  string
}

// spec provides common functionality for Spec implementations.
type spec[T any] struct {
	name string
	desc string

	isValidated bool
	def         *T
	value       T
}

func (s *spec[T]) Name() string {
	return s.name
}

func (s *spec[T]) Value() T {
	if !s.isValidated {
		panic("environment has not been validated")
	}

	return s.value
}

func (s *spec[T]) setDefault(v T) {
	s.def = &v
}

func (s *spec[T]) useValue(v T) {
	s.isValidated = true
	s.value = v
}

func (s *spec[T]) useDefault() (zero T, _ error) {
	if s.def == nil {
		return zero, errUndefined
	}

	s.isValidated = true
	s.value = *s.def

	return s.value, nil
}

// SpecOption is an option that alters the behavior of a variable specification.
type SpecOption func(*specOptions)

// WithRegistry is an option that sets the registry to use for a specific
// variable specification.
func WithRegistry(r *Registry) SpecOption {
	return func(opts *specOptions) {
		opts.Registry = r
	}
}

// specOptions is a set of options that alter the behavior of a variable
// specification.
type specOptions struct {
	Registry *Registry
}
