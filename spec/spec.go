package spec

import (
	"github.com/dogmatiq/ferrite/schema"
)

// Spec is a specification that describes an environment variable.
type Spec struct {
	// Name is the name of the environment variable.
	Name string

	// Description is a human-readable description of the environment variable.
	Description string

	// Necessity describes the necessity (or "optionalness") of the environment
	// variable.
	Necessity Necessity

	// Schema describes the valid values of the environment variable.
	Schema schema.Schema

	// Default is the environment variable's default value.
	Default string
}

// Necessity is an enumeration describing the necessity (or "optionalness") of
// an environment variable.
type Necessity string

const (
	// Required indicates that the application always requires a value and that
	// it must be set explicitly as an environment variable.
	Required Necessity = "required"

	// Defaulted indicates that the application always requires a value but an
	// explicit value is not mandatory because a default is provided.
	Defaulted Necessity = "defaulted"

	// Optional indicates that the application places meaning on the presence or
	// absence of the environment variable and it is therefore truly optional.
	Optional Necessity = "optional"
)
