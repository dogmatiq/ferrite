package ferrite

import (
	"io"
	"os"
)

var (
	// registry is a global registry of environment variable specs.
	registry map[string]Spec

	// output is the writer to which the validation result is written.
	output io.Writer = os.Stderr

	// exit is called to exit the process when validation fails.
	exit = os.Exit
)

// Register adds a variable specification to the register.
func Register(s Spec) {
	if registry == nil {
		registry = map[string]Spec{}
	}

	registry[s.Name()] = s
}
