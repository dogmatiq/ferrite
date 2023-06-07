package validate

import (
	"fmt"

	"github.com/dogmatiq/ferrite/internal/variable"
)

// name renders a column containing the variable's name.
func name(v variable.Any) string {
	s := v.Spec()

	icon := " "
	if attentionNeeded(v) != attentionNone {
		icon = iconAttention
	}

	return fmt.Sprintf(" %s %s", icon, s.Name())
}
