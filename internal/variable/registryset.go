package variable

import (
	"fmt"

	"golang.org/x/exp/slices"
)

// RegistrySet is a set of multiple environment variable registries
type RegistrySet struct {
	registries map[string]*Registry
	variables  []RegisteredVariable
}

// RegisteredVariable is a variable that is associated with a specific source
// [Registry].
type RegisteredVariable struct {
	key string
	Any
	Registry *Registry
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

	for _, v := range s.variables {
		if _, ok := r.vars.Load(v.key); ok {
			panic(fmt.Sprintf(
				"the %q environment variable is defined in both the %q and %q registries",
				v.Spec().Name(),
				v.Registry.Key,
				r.Key,
			))
		}
	}

	r.vars.Range(func(k, v any) bool {
		s.variables = append(
			s.variables,
			RegisteredVariable{
				key:      k.(string),
				Any:      v.(Any),
				Registry: r,
			},
		)
		return true
	})

	slices.SortFunc(
		s.variables,
		func(a, b RegisteredVariable) bool {
			return a.Spec().Name() < b.Spec().Name()
		},
	)

	if s.registries == nil {
		s.registries = make(map[string]*Registry)
	}

	s.registries[r.Key] = r
}

// Variables returns the variables in the registries, sorted by name.
func (s *RegistrySet) Variables() []RegisteredVariable {
	return s.variables
}

// IsEmpty returns true if the set contains no registries.
func (s *RegistrySet) IsEmpty() bool {
	return len(s.registries) == 0
}
