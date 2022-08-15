package markdownmode

import (
	"github.com/dogmatiq/ferrite/variable"
	"github.com/mattn/go-runewidth"
)

func (r *renderer) renderSpecs() {
	r.line("## Specification")

	for _, s := range r.Specs {
		r.line("")
		r.renderSpec(s)
		r.line("")
		r.renderExamples(s)
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

func (r schemaRenderer) VisitNumeric(s variable.Numeric) {
	min, _ := s.Min()

	if def, ok := r.spec.Default(); ok {
		r.line("This variable **MAY** be set to `%s` or greater.", min.Quote())
		r.line("If left undefined the default value of `%s` is used.", def.Quote())
	} else if r.spec.IsRequired() {
		r.line("This variable **MUST** be set to `%s` or greater.", min.Quote())
		r.renderUndefinedFailureWarning()
	} else {
		r.line("This variable **MAY** be set to `%s` or greater, or left undefined.", min.Quote())
	}
}

func (r schemaRenderer) VisitSet(s variable.Set) {
	if def, ok := r.spec.Default(); ok {
		r.line("This variable **MAY** be set to one of the values below.")
		r.line("If left undefined the default value of `%s` is used.", def.Quote())
	} else if r.spec.IsRequired() {
		r.line("This variable **MUST** be set to one of the values below.")
		r.renderUndefinedFailureWarning()
	} else {
		r.line("This variable **MAY** be set to one of the values below or left undefined.")
	}
}

func (r schemaRenderer) VisitString(variable.String) {
	if _, ok := r.spec.Default(); ok {
		r.line("This variable **MAY** be set to a non-empty string.")
		r.line("If left undefined the default value is used (see below).")
	} else if r.spec.IsRequired() {
		r.line("This variable **MUST** be set to a non-empty string.")
		r.renderUndefinedFailureWarning()
	} else {
		r.line("This variable **MAY** be set to a non-empty string or left undefined.")
	}
}

func (r *renderer) renderUndefinedFailureWarning() {
	r.line("If left undefined the application will print usage information to `STDERR` then")
	r.line("exit with a non-zero exit code.")
}

func (r *renderer) renderExamples(s variable.Spec) {
	r.line("```bash")

	width := 0
	for _, eg := range s.Examples() {
		w := runewidth.StringWidth(eg.Canonical.Quote())
		if w > width {
			width = w
		}
	}

	for _, eg := range s.Examples() {
		comment := ""
		if variable.IsDefault(s, eg.Canonical) {
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
				s.Name(),
				eg.Canonical.Quote(),
			)
		} else {
			r.line(
				"export %s=%-*s # %s",
				s.Name(),
				width,
				eg.Canonical.Quote(),
				comment,
			)
		}
	}

	r.line("```")
}
