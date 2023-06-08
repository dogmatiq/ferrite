package markdown

import (
	"fmt"

	"github.com/dogmatiq/ferrite/internal/variable"
)

func (r *renderer) renderIndex() {
	if len(r.Variables) == 0 {
		r.line("**There do not appear to be any environment variables.**")
	} else {
		var t table

		t.AddRow("Name", "Optionality", "Description")

		for _, v := range r.Variables {
			s := v.Spec()
			name := r.linkToSpec(s)
			optionality := "required"

			if s.IsDeprecated() {
				name = "~~" + name + "~~"
				optionality = "optional, deprecated"
			} else if def, ok := s.Default(); ok {
				optionality = fmt.Sprintf("defaults to `%s`", def.Quote())
			} else if !s.IsRequired() {
				optionality = "optional"
			} else if len(variable.Relationships[variable.DependsOn](s)) != 0 {
				optionality = "conditional"
			}

			t.AddRow(name, optionality, s.Description())
		}

		t.WriteTo(r.Output)
	}
}
