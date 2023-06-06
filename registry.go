package ferrite

import (
	"fmt"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// A Registry is a collection of environment variable specifications.
type Registry = *variable.Registry

// NewRegistry returns a new environment variable registry.
//
// Use the [WithRegistry] option to configure an environment variable definition
// or [Init] call to use a specific registry.
func NewRegistry(name string) Registry {
	if name == "" {
		panic("registry name must not be empty")
	}

	if name == variable.DefaultRegistry.Name {
		panic(fmt.Sprintf(
			"registry name must not be %q",
			variable.DefaultRegistry.Name,
		))
	}

	return &variable.Registry{
		Name: name,
	}
}
