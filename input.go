package ferrite

import (
	"github.com/dogmatiq/ferrite/variable"
)

// An Input is the application-facing interface for obtaining a value from
// environment variables.
//
// Typically each input is sourced from exactly one environment variable,
// however it is possible that a value collates values from multiple variables.
//
// An input is created using one of the various builder-functions, for example
// String().
type Input interface {
	variables() []variable.Any
}

type inputConfig struct {
	Registry *variable.Registry
}
