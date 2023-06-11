package ferrite

import (
	"github.com/dogmatiq/ferrite/internal/variable"
)

// A Registry is a collection of environment variable specifications.
type Registry interface {
	variable.ProtectedRegistry
}

// RegistryOption is an option that configures the behavior of a registry.
type RegistryOption interface {
	applyRegistryOption(*variable.Registry)
}

// NewRegistry returns a new environment variable registry.
//
// Use the [WithRegistry] option to configure an environment variable definition
// or [Init] call to use a specific registry.
func NewRegistry(
	key, name string,
	options ...RegistryOption,
) Registry {
	if key == "" {
		panic("registry key must not be empty")
	}

	if name == "" {
		panic("registry name must not be empty")
	}

	reg := &variable.Registry{
		Key:  key,
		Name: name,
	}

	for _, opt := range options {
		opt.applyRegistryOption(reg)
	}

	return reg
}
