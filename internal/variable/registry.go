package variable

import (
	"net/url"
	"sync"

	"github.com/dogmatiq/ferrite/internal/environment"
)

// Registry is a collection of environment variable specifications.
type Registry struct {
	Key, Name string
	URL       *url.URL
	IsDefault bool

	vars sync.Map // map[string]Any
}

// Register registers a new variable with the registry.
func (r *Registry) Register(v Any) {
	name := v.Spec().Name()
	norm := environment.NormalizeName(name)

	if _, loaded := r.vars.LoadOrStore(norm, v); loaded {
		panic("a variable named " + name + " is already registered")
	}
}

// Assign copies the contents of reg into r.
func (r *Registry) Assign(reg *Registry) {
	r.Key = reg.Key
	r.Name = reg.Name
	r.URL = reg.URL
	r.IsDefault = reg.IsDefault

	r.vars.Range(func(k any, _ any) bool {
		DefaultRegistry.vars.Delete(k)
		return true
	})

	reg.vars.Range(func(k any, v any) bool {
		r.vars.Store(k, v)
		return true
	})
}

// Clone returns a copy of the registry.
func (r *Registry) Clone() *Registry {
	c := &Registry{}
	c.Assign(r)
	return c
}

// DefaultRegistry is the default specification registry.
var DefaultRegistry = &Registry{
	IsDefault: true,
}

// ResetDefaultRegistry removes all variables from [DefaultRegistry].
func ResetDefaultRegistry() {
	DefaultRegistry.Assign(
		&Registry{
			IsDefault: true,
		},
	)
}

// Register registers a new variable with one or more registries.
//
// If no registries are specified, [DefaultRegistry] is used.
func Register[T any](
	registries []*Registry,
	spec *TypedSpec[T],
) *OfType[T] {
	if len(registries) == 0 {
		registries = append(registries, DefaultRegistry)
	}

	v := &OfType[T]{
		TypedSpec: spec,
	}

	for _, reg := range registries {
		reg.Register(v)
	}

	return v
}

// ProtectedRegistry is an interface that allows access to the internals of a
// [Registry].
type ProtectedRegistry interface {
	expose() *Registry
}

// ExposeRegistry returns exposes the underlying registry of r.
func ExposeRegistry(r ProtectedRegistry) *Registry {
	return r.expose()
}

func (r *Registry) expose() *Registry {
	return r
}
