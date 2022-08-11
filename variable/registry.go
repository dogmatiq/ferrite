package variable

import (
	"sync"
)

// Registry is a collection of environment variable specifications.
type Registry struct {
	Environment Environment

	vars sync.Map // map[String]Variable
}

// Variables returns the variables in the registry.
func (r *Registry) Variables() []Variable {
	var vars []Variable

	r.vars.Range(func(_, value any) bool {
		vars = append(vars, value.(Variable))
		return true
	})

	return vars
}

// Reset removes all variables from the registry.
func (r *Registry) Reset() {
	r.vars.Range(func(k, _ any) bool {
		r.vars.Delete(k)
		return true
	})
}

// DefaultRegistry is the default specification registry.
var DefaultRegistry = Registry{
	Environment: OSEnvironment,
}

// Register registers a new variable.
func Register[T any](
	spec Spec[T],
	options []RegisterOption,
) *TypedVariable[T] {
	opts := registerOptions{
		Registry: &DefaultRegistry,
	}
	for _, opt := range options {
		opt(&opts)
	}

	v := &TypedVariable[T]{
		spec: spec,
		env:  opts.Registry.Environment,
	}

	opts.Registry.vars.Store(spec.Name, v)

	return v
}

// RegisterOption is an option that controls how a specification is registered
// with an environment variable registry.
type RegisterOption func(*registerOptions)

// registerOptions contains options that are available to all specifications.
type registerOptions struct {
	Registry *Registry
}

// WithRegistry is an option that sets the registry that an environment variable
// specification is placed into.
func WithRegistry(r *Registry) RegisterOption {
	return func(o *registerOptions) {
		o.Registry = r
	}
}
