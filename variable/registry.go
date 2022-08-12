package variable

import (
	"sync"
)

// Registry is a collection of environment variable specifications.
type Registry struct {
	Environment Environment

	vars sync.Map // map[String]Variable
}

// Range calls fn for each variable in the registry.
//
// It stops iterating if fn returns false.
func (r *Registry) Range(fn func(Any) bool) {
	r.vars.Range(func(_, value any) bool {
		return fn(value.(Any))
	})
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
	spec PendingSpecFor[T],
	options []RegisterOption,
) *OfType[T] {
	opts := registerOptions{
		Registry: &DefaultRegistry,
	}
	for _, opt := range options {
		opt(&opts)
	}

	v := &OfType[T]{
		spec: finalizeSpec(spec),
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
