package variable

import (
	"fmt"
	"sync"

	"github.com/dogmatiq/ferrite/internal/environment"
	"golang.org/x/exp/slices"
)

// Registry is a collection of environment variable specifications.
type Registry struct {
	Name string
	vars sync.Map // map[string]Any
}

func (r *Registry) register(v Any) {
	name := v.Spec().Name()
	norm := environment.NormalizeName(name)

	if _, loaded := r.vars.LoadOrStore(norm, v); loaded {
		panic("a variable named " + name + " is already registered")
	}
}

func (r *Registry) clone() *Registry {
	c := &Registry{
		Name: r.Name,
	}

	r.vars.Range(func(k any, v any) bool {
		c.vars.Store(k, v)
		return true
	})

	return c
}

// ResetDefaultRegistry removes all variables from [DefaultRegistry].
func ResetDefaultRegistry() {
	DefaultRegistry.vars.Range(func(k any, _ any) bool {
		DefaultRegistry.vars.Delete(k)
		return true
	})
}

// DefaultRegistry is the default specification registry.
var DefaultRegistry = Registry{
	Name: "default",
}

// Register registers a new variable with one or more registries.
//
// If no registries are specified, [DefaultRegistry] is used.
func Register[T any](
	registries []*Registry,
	spec *TypedSpec[T],
) *OfType[T] {
	if len(registries) == 0 {
		registries = append(registries, &DefaultRegistry)
	}

	v := &OfType[T]{
		spec: spec,
	}

	for _, reg := range registries {
		reg.register(v)
	}

	return v
}

// RegistrySet is a set of multiple environment variable registries
type RegistrySet struct {
	registries []*Registry
}

// Add adds r to the set.
//
// It panics if r contains any variables that are already defined in another
// registry in the set.
//
// Any changes to r after it is added to the set are not reflected in the set.
func (s *RegistrySet) Add(r *Registry) {
	for _, x := range s.registries {
		x.vars.Range(func(k, _ any) bool {
			if v, ok := r.vars.Load(k); ok {
				panic(fmt.Sprintf(
					"the %q environment variable is defined in both the %q and %q registries",
					v.(Any).Spec().Name(),
					x.Name,
					r.Name,
				))
			}
			return true
		})
	}

	s.registries = append(s.registries, r.clone())
}

// Specs returns the specs of the variables in the registries, sorted by name.
func (s *RegistrySet) Specs() []Spec {
	variables := s.Variables()
	specs := make([]Spec, len(variables))

	for i, v := range variables {
		specs[i] = v.Spec()
	}

	return specs
}

// Variables returns the variables in the registries, sorted by name.
func (s *RegistrySet) Variables() []Any {
	var variables []Any

	for _, r := range s.registries {
		r.vars.Range(func(_, v any) bool {
			variables = append(variables, v.(Any))
			return true
		})
	}

	slices.SortFunc(
		variables,
		func(a, b Any) bool {
			return a.Spec().Name() < b.Spec().Name()
		},
	)

	return variables
}

// IsEmpty returns true if the set contains no registries.
func (s *RegistrySet) IsEmpty() bool {
	return len(s.registries) == 0
}
