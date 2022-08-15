package markdownmode

import (
	"github.com/dogmatiq/ferrite/variable"
	"github.com/mattn/go-runewidth"
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
	if def, ok := r.spec.Default(); ok {
		r.line("This variable **MAY** be set to one of the values below.")
		r.line("If it is undefined or empty a default value of `%s` is used.", def.Quote())
	} else if r.spec.IsRequired() {
		r.line("This variable **MUST** be set to one of the values below.")
		r.line("If it is undefined or empty the application will print usage information to")
		r.line("`STDERR` then exit with a non-zero exit code.")
	} else {
		r.line("This variable **MAY** be set to one of the values below, or remain undefined.")
	}

	r.line("")
	r.line("```bash")

	width := 0
	for _, eg := range r.spec.Examples() {
		w := runewidth.StringWidth(eg.Canonical.Quote())
		if w > width {
			width = w
		}
	}

	for _, eg := range r.spec.Examples() {
		comment := ""
		if variable.IsDefault(r.spec, eg.Canonical) {
			comment = "(default)"
		}
		if eg.Description != "" {
			if comment != "" {
				comment += " "
			}
			comment += eg.Description
		}

		if len(comment) == 0 {
			r.line(
				"export %s=%s",
				r.spec.Name(),
				eg.Canonical.Quote(),
			)
		} else {
			r.line(
				"export %s=%-*s # %s",
				r.spec.Name(),
				width,
				eg.Canonical.Quote(),
				comment,
			)
		}
	}

	r.line("```")
}

func (r schemaRenderer) VisitString(variable.String) {
}
