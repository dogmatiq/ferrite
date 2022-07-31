package ferrite

import (
	"os"
)

// ResolveEnvironment validates all environment variables.
//
// It panics if any of the defined variables are invalid.
func ResolveEnvironment() {
	DefaultRegistry.MustResolve(os.LookupEnv)
}

// Lookup retrieves the value of an environment variable.
//
// If the variable is present in the environment, value (which may be empty)
// is returned and ok is true. Otherwise, the value is empty and ok is false.
type Lookup func(name string) (value string, ok bool)

// DefaultRegistry is the default environment variable registry.
var DefaultRegistry Registry

// Registry is a container of environment variable specifications.
type Registry struct {
	specs map[string]Spec
}

// Spec is a specification for an environment variable.
//
// It describes the environment variable itself, and how to construct valid
// values for the variable.
type Spec interface {
	Name() string
	Resolve(Lookup) error
}

// Register adds a variable specification to the register.
func (r *Registry) Register(s Spec) {
	if r.specs == nil {
		r.specs = map[string]Spec{}
	}

	r.specs[s.Name()] = s
}

// Reset removes all variable specifications from the registry.
func (r *Registry) Reset() {
	r.specs = nil
}

// Resolve parses and validates all environment variables in the registry,
// allowing their associated values to be obtained.
func (r *Registry) Resolve(lookup Lookup) error {
	for _, s := range r.specs {
		if err := s.Resolve(lookup); err != nil {
			return err
		}
	}

	return nil
}

// MustResolve parses and validates all environment variables in the registry,
// allowing their associated values to be obtained.
//
// It panics if any of the specifications are violated.
func (r *Registry) MustResolve(lookup Lookup) {
	if err := r.Resolve(lookup); err != nil {
		panic(err.Error())
	}
}
