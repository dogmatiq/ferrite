package ferrite

import (
	"errors"
)

var (
	// errUndefined is an error indicating that a mandatory environment variable
	// with no default value was not set explicitly in the environment.
	errUndefined = errors.New("must not be empty")
)
