package ferrite

import (
	"github.com/dogmatiq/ferrite/internal/variable"
)

// VariableSet is the application-facing interface for obtaining a value
// constructed from one or more environment variables.
//
// A variable set is obtained using one of the many "builder" types returned by
// functions such as Bool(), String(), etc.
//
// It is common for a variable set to contain a single variable. However, some
// builders produce sets containing multiple variables.
type VariableSet[T any] interface {
	variables() []variable.Any
	native() (T, bool)
	literals(T) ([]variable.Literal, error)
}

// variableSetConfig encapsulates configuration common to all variable sets.
type variableSetConfig struct {
	Registries []*variable.Registry
}
