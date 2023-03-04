package markdownmode

import (
	"github.com/dogmatiq/ferrite/variable"
)

func (r *renderer) renderSpecs() {
	r.line("## Specification")

	for _, s := range r.Specs {
		r.renderSpec(s)
	}
}

func (r *renderer) renderSpec(s variable.Spec) {
	r.line("")
	r.line("### `%s`", s.Name())

	r.line("")
	r.line("> %s", s.Description())

	r.line("")

	s.Schema().AcceptVisitor(schemaRenderer{r, s})

	for _, d := range s.Documentation() {
		r.line("")
		r.paragraph(d)
	}

	r.line("")
	r.renderExamples(s)
}

type schemaRenderer struct {
	*renderer
	spec variable.Spec
}

func (r schemaRenderer) VisitNumeric(s variable.Numeric) {
	if min, ok := s.Min(); ok {
		if def, ok := r.spec.Default(); ok {
			r.line("This variable **MAY** be set to `%s` or greater.", min.Quote())
			r.line("If left undefined the default value of `%s` is used.", def.Quote())
		} else if r.spec.IsRequired() {
			r.line("This variable **MUST** be set to `%s` or greater.", min.Quote())
		} else {
			r.line("This variable **MAY** be set to `%s` or greater, or left undefined.", min.Quote())
		}
	} else {
		if def, ok := r.spec.Default(); ok {
			r.line("This variable **MAY** be set to a `%s` value.", s.Type().Kind())
			r.line("If left undefined the default value of `%s` is used.", def.Quote())
		} else if r.spec.IsRequired() {
			r.line("This variable **MUST** be set to a `%s` value.", s.Type().Kind())
		} else {
			r.line("This variable **MAY** be set to a `%s` value or left undefined.", s.Type().Kind())
		}
	}
}

func (r schemaRenderer) VisitSet(s variable.Set) {
	if def, ok := r.spec.Default(); ok {
		r.line("This variable **MAY** be set to one of the values below.")
		r.line("If left undefined the default value of `%s` is used.", def.Quote())
	} else if r.spec.IsRequired() {
		r.line("This variable **MUST** be set to one of the values below.")
	} else {
		r.line("This variable **MAY** be set to one of the values below or left undefined.")
	}
}

func (r schemaRenderer) VisitString(variable.String) {
	if _, ok := r.spec.Default(); ok {
		r.renderDefaultClause()
	} else if r.spec.IsRequired() {
		r.line("This variable **MUST** be set to a non-empty string.")
	} else {
		r.line("This variable **MAY** be set to a non-empty string or left undefined.")
	}
}

func (r schemaRenderer) VisitOther(variable.Other) {
	if _, ok := r.spec.Default(); ok {
		r.renderDefaultClause()
		r.renderConstraints(r.spec)
	} else if r.spec.IsRequired() {
		r.renderConstraints(r.spec)
	} else {
		r.line("This variable **MAY** be left undefined.")
		r.renderConstraints(r.spec)
	}
}

func (r schemaRenderer) renderDefaultClause() {
	v, _ := r.spec.Default()
	r.line("This variable **MAY** be left undefined, in which case the default value")
	r.line("of `%s` is used.", v.String)
}
