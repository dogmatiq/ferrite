package variable

import (
	"sync"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Registry is a collection of environment variable specifications.
type Registry struct {
	Environment Environment

	m    sync.RWMutex
	vars map[string]Any
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
	r.m.RLock()
	variables := maps.Values(r.vars)
	r.m.RUnlock()

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
	r.m.Lock()
	r.vars = nil
	r.m.Unlock()
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

	name := normalizeVariableName(spec.name)

	reg.m.Lock()
	defer reg.m.Unlock()

	if reg.vars == nil {
		reg.vars = map[string]Any{}
	} else if _, ok := reg.vars[name]; ok {
		panic("a variable named " + spec.name + " is already registered")
	}

	v := &OfType[T]{
		spec: spec,
		env:  reg.Environment,
	}

	reg.vars[name] = v

	return v
}
