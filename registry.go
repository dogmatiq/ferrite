package ferrite

import (
	"fmt"
	urls "net/url"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// A Registry is a collection of environment variable specifications.
type Registry = *variable.Registry

// RegistryOption is an option that configures the behavior of a registry.
type RegistryOption interface {
	applyRegistryOption(*variable.Registry)
}

// NewRegistry returns a new environment variable registry.
//
// Use the [WithRegistry] option to configure an environment variable definition
// or [Init] call to use a specific registry.
func NewRegistry(name, url string, options ...RegistryOption) Registry {
	if name == "" {
		panic("registry name must not be empty")
	}

	if name == variable.DefaultRegistry.Name {
		panic(fmt.Sprintf(
			"registry name must not be %q",
			variable.DefaultRegistry.Name,
		))
	}

	u, err := urls.Parse(url)
	if err != nil {
		panic("invalid URL: " + err.Error())
	}

	reg := &variable.Registry{
		Name: name,
		URL:  u,
	}

	for _, opt := range options {
		opt.applyRegistryOption(reg)
	}

	return reg
}
