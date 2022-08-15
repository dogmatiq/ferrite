package markdownmode

import (
	"strings"

	"github.com/dogmatiq/ferrite/variable"
)

func (r *renderer) renderSpecs() {
	r.line("## Specification")
	r.line("")

	for _, s := range r.Specs {
		r.renderSpec(s)
	}
}

func (r *renderer) renderSpec(s variable.Spec) {
	r.line("### `%s`", s.Name())

	r.line("")
	r.line("> %s", s.Description())

	r.line("")
	s.Schema().AcceptVisitor(schemaRenderer{r, s})
}

type schemaRenderer struct {
	*renderer
	spec variable.Spec
}

func (r schemaRenderer) VisitNumeric(variable.Numeric) {
}

func (r schemaRenderer) VisitSet(s variable.Set) {
	literals := s.Literals()

	if len(literals) == 2 {
		if def, ok := r.spec.Default(); ok {
			r.line(
				"This variable **MAY** be set to either `%s` or `%s`. If it is undefined or",
				literals[0].Quote(),
				literals[1].Quote(),
			)
			r.line("empty a default value of `%s` is used.", def.Quote())
		} else if r.spec.IsRequired() {
			r.line(
				"This variable **MUST** be set to either `%s` or `%s`. If it is undefined or",
				literals[0].Quote(),
				literals[1].Quote(),
			)
			r.line("empty the application will print usage information to `STDERR` then exit with a")
			r.line("non-zero exit code.")
		} else {
			r.line(
				"This variable **MAY** be set to either `%s` or `%s`, or remain undefined.",
				literals[0].Quote(),
				literals[1].Quote(),
			)
		}
	}

	r.line("")
	r.line("```bash")

	for _, eg := range r.spec.Examples() {
		var comments []string
		if eg.Description != "" {
			comments = append(comments, eg.Description)
		}
		if variable.IsDefault(r.spec, eg.Canonical) {
			comments = append(comments, "default value")
		}

		if len(comments) == 0 {
			r.line(
				"export %s=%s",
				r.spec.Name(),
				eg.Canonical.Quote(),
			)
		} else {
			r.line(
				"export %s=%s # %s",
				r.spec.Name(),
				eg.Canonical.Quote(),
				strings.Join(comments, ", "),
			)
		}
	}

	r.line("```")
}

func (r schemaRenderer) VisitString(variable.String) {
}
