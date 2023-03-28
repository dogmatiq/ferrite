package markdown

import (
	"fmt"

	"github.com/dogmatiq/ferrite/variable"
)

func (r *renderer) renderIndex() {
	r.line("## Index")
	r.gap()

	var t table

	t.AddRow("Name", "Optionality", "Description")

	for _, s := range r.Specs {
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

// | Name                 | Optionality           | Description                            |
// | -------------------- | --------------------- | -------------------------------------- |
// | ~~[`BIND_ADDRESS`]~~ | optional, deprecated  | listen address for the HTTP server |
// | [`BIND_HOST`]        | defaults to `0.0.0.0` | listen host for the HTTP server        |
// | [`BIND_PORT`]        | defaults to `8080`    | listen port for the HTTP server        |
// | [`BIND_VERSION`]     | defaults to `4`       | IP version for the HTTP server         |
