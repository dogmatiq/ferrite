package variable

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/dogmatiq/ferrite/internal/environment"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Registry is a collection of environment variable specifications.
type Registry struct {
	Key, Name string
	URL       *url.URL

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
		Key:  r.Key,
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
	Key:  "default",
	Name: "Default",
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
	registries         map[string]*Registry
	variables          map[string]Any
	registryByVariable map[Any]*Registry
}

// Add adds r to the set.
//
// It panics if r contains any variables that are already defined in another
// registry in the set.
//
// Any changes to r after it is added to the set are not reflected in the set.
func (s *RegistrySet) Add(r *Registry) {
	if _, ok := s.registries[r.Key]; ok {
		panic(fmt.Sprintf(
			"the set already contains the %s registry",
			r.Key,
		))
	}

	for k, v := range s.variables {
		if _, ok := r.vars.Load(k); ok {
			panic(fmt.Sprintf(
				"the %q environment variable is defined in both the %q and %q registries",
				v.Spec().Name(),
				s.registryByVariable[v].Key,
				r.Key,
			))
		}
	}

	if s.variables == nil {
		s.variables = map[string]Any{}
	}

	r.vars.Range(func(k, v any) bool {
		s.variables[k.(string)] = v.(Any)
		return true
	})

	if s.registries == nil {
		s.registries = make(map[string]*Registry)
	}

	s.registries[r.Key] = r
}

// Variables returns the variables in the registries, sorted by name.
func (s *RegistrySet) Variables() []Any {
	variables := maps.Values(s.variables)

	slices.SortFunc(
		variables,
		func(a, b Any) bool {
			return a.Spec().Name() < b.Spec().Name()
		},
	)

	return variables
}

// SourceRegistry returns the registry that contains v.
func (s *RegistrySet) SourceRegistry(v Any) *Registry {
	r, ok := s.registryByVariable[v]
	if !ok {
		panic("unrecognized variable")
	}
	return r
}

// IsEmpty returns true if the set contains no registries.
func (s *RegistrySet) IsEmpty() bool {
	return len(s.registries) == 0
}
