package ferrite

import (
	"fmt"
	"os"
)

// ResolveEnvironment validates all environment variables.
//
// It panics if any of the defined variables are invalid.
func ResolveEnvironment() {
	DefaultRegistry.MustResolve(os.LookupEnv)
}

// LookupFn retrieves the value of an environment variable.
//
// If the variable is present in the environment, value (which may be empty)
// is returned and ok is true. Otherwise, the value is empty and ok is false.
type LookupFn func(name string) (value string, ok bool)

// DefaultRegistry is the default environment variable registry.
var DefaultRegistry = Registry{
	fatal: func(err error) {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	},
}

// Registry is a container of environment variable specifications.
type Registry struct {
	fatal func(error)
	specs map[string]Spec
}

// Spec is a specification for an environment variable.
//
// It describes the environment variable itself, and how to construct valid
// values for the variable.
type Spec interface {
	Name() string
	Resolve(LookupFn) error
}

func (r *Registry) Register(s Spec) {
	if r.specs == nil {
		r.specs = map[string]Spec{}
	}

	r.specs[s.Name()] = s
}

func (r *Registry) Resolve(lookup LookupFn) error {
	for _, s := range r.specs {
		if err := s.Resolve(lookup); err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry) MustResolve(lookup LookupFn) {
	if err := r.Resolve(lookup); err != nil {
		r.fatal(err)
	}
}
