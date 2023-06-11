package markdown

import (
	"fmt"

	"github.com/dogmatiq/ferrite/internal/variable"
)

func (r *renderer) renderIndex() {
	if len(r.Variables) == 0 {
		r.line("**There do not appear to be any environment variables.**")
		return
	}

	hasImportColumn := false
	for _, v := range r.Variables {
		if !v.Registry.IsDefault {
			hasImportColumn = true
			break
		}
	}

	var t table

	if hasImportColumn {
		t.AddRow("Name", "Optionality", "Description", "Imported From")
	} else {
		t.AddRow("Name", "Optionality", "Description")
	}

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

		if hasImportColumn {
			reg := ""
			if !v.Registry.IsDefault {
				reg = r.linkToRegistry(v.Registry)
			}
			t.AddRow(name, optionality, s.Description(), reg)
		} else {
			t.AddRow(name, optionality, s.Description())
		}
	}

	t.WriteTo(r.Output)
}
