package spec

// Spec is a specification that describes an environment variable.
type Spec struct {
	// Name is the name of the environment variable.
	Name string

	// Description is a human-readable description of the environment variable.
	Description string

	// IsOptional is true if the application does not require a value for this
	// variable.
	IsOptional bool

	// HasDefault is true if the user does not need to explicitly set a value
	// for this environment variable.
	HasDefault bool

	// Default is the environment variable's default value.
	DefaultX string

	// Schema describes the valid values of the environment variable.
	Schema Schema
}
