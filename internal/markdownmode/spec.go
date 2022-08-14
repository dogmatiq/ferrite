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

	if s.IsRequired() {
		r.line("This variable is **required**, although a default is provided.")
	}

	s.Schema().AcceptVisitor(schemaRenderer{r, s})
}

type schemaRenderer struct {
	*renderer
	spec variable.Spec
}

func (r schemaRenderer) VisitNumeric(variable.Numeric) {
}

func (r schemaRenderer) VisitSet(s variable.Set) {
	r.line("It must be one of the values shown in the examples below.")
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
