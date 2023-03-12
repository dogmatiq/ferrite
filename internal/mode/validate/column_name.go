package validate

import (
	"fmt"

	"github.com/dogmatiq/ferrite/variable"
)

// name renders a column containing the variable's name.
func name(v variable.Any) string {
	s := v.Spec()

	icon := " "
	if needsAttention(v) {
		icon = iconAttention
	}

	return fmt.Sprintf(" %s %s", icon, s.Name())
}
