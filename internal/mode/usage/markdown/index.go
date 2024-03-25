package markdown

import (
	"fmt"

	"github.com/dogmatiq/ferrite/internal/variable"
)

func (r *renderer) renderIndex() {
	if len(r.Variables) == 0 {
		r.line("**`%s` does not appear to use any environment variables.**", r.App)
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
		t.AddRow("Name", "Usage", "Description", "Imported From")
	} else {
		t.AddRow("Name", "Usage", "Description")
	}

	for _, v := range r.Variables {
		s := v.Spec()
		name := r.linkToSpec(s)
		usage := "required"

		if s.IsDeprecated() {
			name = "~~" + name + "~~"
			usage = "optional, deprecated"
		} else if def, ok := s.Default(); ok {
			usage = fmt.Sprintf("defaults to `%s`", def.Quote())
		} else if !s.IsRequired() {
			usage = "optional"
		} else if len(variable.Relationships[variable.DependsOn](s)) != 0 {
			usage = "conditional"
		}

		if hasImportColumn {
			reg := ""
			if !v.Registry.IsDefault {
				reg = r.linkToRegistry(v.Registry)
			}
			t.AddRow(name, usage, s.Description(), reg)
		} else {
			t.AddRow(name, usage, s.Description())
		}
	}

	t.WriteTo(r.Output)
}
