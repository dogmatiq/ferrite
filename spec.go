package ferrite

// Spec is a specification for an environment variable.
//
// It describes the environment variable itself, and how to construct valid
// values for the variable.
type Spec interface {
	Name() string
	Resolve(Environment) error
}

// spec provides common functionality for Spec implementations.
type spec struct {
	name string
	desc string
}

func (s spec) Name() string {
	return s.name
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
