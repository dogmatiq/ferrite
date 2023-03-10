package validatemode

import "github.com/dogmatiq/ferrite/variable"

// description renders a column containing the variable's human-readable
// description.
func description(v variable.Any) string {
	return v.Spec().Description()
}
