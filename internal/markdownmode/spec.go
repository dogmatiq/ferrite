package markdownmode

import "github.com/dogmatiq/ferrite/variable"

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
		r.line(
			"This variable **MAY** be set to either `%s` or `%s`. If it is undefined or",
			literals[0].Quote(),
			literals[1].Quote(),
		)
		r.line("empty a default value of `%s` is used.", r.spec.Default().MustGet().Quote())
	}

	r.line("")
	r.line("```bash")

	def, _ := r.spec.Default().Get()

	for _, m := range s.Literals() {
		if m == def {
			r.line(
				"export %s=%s # default value",
				r.spec.Name(),
				m.Quote(),
			)
		} else {
			r.line(
				"export %s=%s",
				r.spec.Name(),
				m.Quote(),
			)
		}
	}

	r.line("```")
}

func (r schemaRenderer) VisitString(variable.String) {
}
