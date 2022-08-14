package variable

import (
	"sync"

	"golang.org/x/exp/slices"
)

// Registry is a collection of environment variable specifications.
type Registry struct {
	Environment Environment

	vars sync.Map // map[String]Variable
}

// Specs returns the specs of the variables in the library, sorted by name.
func (r *Registry) Specs() []Spec {
	variables := r.Variables()
	specs := make([]Spec, len(variables))

	for i, v := range variables {
		specs[i] = v.Spec()
	}

	return specs
}

// Variables returns the variables in the registry, sorted by name.
func (r *Registry) Variables() []Any {
	var variables []Any

	r.vars.Range(func(_, v any) bool {
		variables = append(variables, v.(Any))
		return true
	})

	slices.SortFunc(
		variables,
		func(a, b Any) bool {
			return a.Spec().Name() < b.Spec().Name()
		},
	)

	return variables
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
	spec TypedSpec[T],
	options []RegisterOption,
) *OfType[T] {
	opts := registerOptions{
		Registry: &DefaultRegistry,
	}
	for _, opt := range options {
		opt(&opts)
	}

	v := &OfType[T]{
		spec: spec,
		env:  opts.Registry.Environment,
	}

	opts.Registry.vars.Store(spec.name, v)

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
