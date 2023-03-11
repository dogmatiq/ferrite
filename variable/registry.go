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
func Register[T any](reg *Registry, spec *TypedSpec[T]) *OfType[T] {
	if reg == nil {
		reg = &DefaultRegistry
	}

	v := &OfType[T]{
		spec: spec,
		env:  reg.Environment,
	}

	reg.vars.Store(spec.name, v)

	return v
}
