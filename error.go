package ferrite

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// errUndefined is an error indicating that a mandatory environment variable
	// with no default value was not set explicitly in the environment.
	errUndefined = errors.New("must be defined and not empty")
)

// errNotInList returns an error indicating that the environment variable value
// was expected to be a member of the given list but was some other value.
func errNotInList(values ...string) error {
	n := len(values)

	switch n {
	case 0:
		panic("must have at least one value")
	case 1:
		return fmt.Errorf("must be %q", values[0])
	case 2:
		return fmt.Errorf("must be either %q or %q", values[0], values[1])
	default:
		var w strings.Builder

		w.WriteString("must be one of ")
		fmt.Fprintf(&w, "%q", values[0])

		for _, v := range values[1 : n-1] {
			fmt.Fprintf(&w, ", %q", v)
		}

		fmt.Fprintf(&w, " or %q", values[n-1])

		return errors.New(w.String())
	}
}
