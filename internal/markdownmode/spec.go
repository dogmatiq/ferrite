package markdownmode

import (
	"fmt"
	"reflect"

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

	s.Schema().AcceptVisitor(schemaRenderer{r, s})

	r.renderImportantDocumentation(s)

	r.renderExamples(s)

	r.renderUnimportantDocumentation(s)
}

func (r *renderer) renderImportantDocumentation(s variable.Spec) {
	for _, d := range s.Documentation() {
		if d.IsImportant {
			for _, p := range d.Paragraphs {
				r.paragraph(p)()
			}
		}
	}
}

func (r *renderer) renderUnimportantDocumentation(s variable.Spec) {
	for _, d := range s.Documentation() {
		if d.IsImportant {
			continue
		}

		r.line("")
		r.line("<details>")

		if d.Summary != "" {
			r.line("<summary>%s</summary>", d.Summary)
		}

		for _, p := range d.Paragraphs {
			r.paragraph(p)()
		}

		r.line("")
		r.line("</details>")
	}
}

type schemaRenderer struct {
	*renderer
	spec variable.Spec
}

func (r schemaRenderer) VisitNumeric(s variable.Numeric) {
	switch s.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r.renderSigned(s)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r.renderUnsigned(s)
	case reflect.Float32, reflect.Float64:
		r.renderFloat(s)
	default:
		panic("unsupported numeric type")
	}
}

func (r schemaRenderer) renderSigned(s variable.Numeric) {
	min, hasMin := s.Min()
	max, hasMax := s.Max()

	if hasMin && hasMax {
		r.renderPrimaryConstrant("**MUST** be between `%s` and `%s`", min.String, max.String)
	} else if hasMin {
		r.renderPrimaryConstrant("**MUST** be `%s` or greater", min.String)
	} else if hasMax {
		r.renderPrimaryConstrant("**MUST** be `%s` or less", max.String)
	} else {
		r.renderPrimaryConstrant("**MUST** be a whole number")
	}
}

func (r schemaRenderer) renderUnsigned(s variable.Numeric) {
	if min, ok := s.Min(); ok {
		if _, ok := r.spec.Default(); ok {
			r.renderPrimaryConstrant("")
			r.line("This variable **MAY** be set to `%s` or greater.", min.Quote())
		} else if r.spec.IsRequired() {
			r.line("This variable **MUST** be set to `%s` or greater.", min.Quote())
		} else {
			r.line("This variable **MAY** be set to `%s` or greater, or left undefined.", min.Quote())
		}
	} else {
		r.renderPrimaryConstrant("**MUST** be a non-negative whole number")
	}
}

func (r schemaRenderer) renderFloat(s variable.Numeric) {
	if min, ok := s.Min(); ok {
		if _, ok := r.spec.Default(); ok {
			r.renderPrimaryConstrant("")
			r.line("This variable **MAY** be set to `%s` or greater.", min.Quote())
		} else if r.spec.IsRequired() {
			r.line("This variable **MUST** be set to `%s` or greater.", min.Quote())
		} else {
			r.line("This variable **MAY** be set to `%s` or greater, or left undefined.", min.Quote())
		}
	} else {
		r.renderPrimaryConstrant("**MUST** be a number with an **OPTIONAL** fractional part")
	}
}

func (r schemaRenderer) VisitSet(s variable.Set) {
	if lits := s.Literals(); len(lits) == 2 {
		r.renderPrimaryConstrant(
			"**MUST** be either `%s` or `%s`",
			lits[0].String,
			lits[1].String,
		)
	} else {
		r.renderPrimaryConstrant("**MUST** be one of the values shown in the examples below")
	}
}

func (r schemaRenderer) VisitString(variable.String) {
	// Find the best constraint to use as the "primary" requirement, favouring
	// non-user-defined constraints.
	var con variable.Constraint
	for _, c := range r.spec.Constraints() {
		if !c.IsUserDefined() {
			con = c
			break
		} else if con == nil {
			con = c
		}
	}

	if con == nil {
		r.renderPrimaryConstrant("")
	} else {
		r.renderPrimaryConstrant(con.Description())
	}

}

func (r schemaRenderer) VisitOther(s variable.Other) {
	con := ""
	for _, c := range r.spec.Constraints() {
		con = c.Description()
		break
	}

	r.renderPrimaryConstrant(con)
}

// renderPrimaryConstrant renders information about the most important
// requirement of the variable's schema, this includes information about whether
// the variable is optional and the basic data type of the variable.
//
// The (optional) requirement text must complete the phrase "the value...". It
// should not include any trailing punctuation.
func (r schemaRenderer) renderPrimaryConstrant(f string, v ...any) {
	req := fmt.Sprintf(f, v...)

	if req != "" {
		if def, ok := r.spec.Default(); ok {
			r.paragraph(
				"The `%s` variable **MAY** be left undefined, in which case the default value of `%s` is used.",
				"Otherwise, the value %s.",
			)(
				r.spec.Name(),
				def.String,
				req,
			)
		} else if r.spec.IsRequired() {
			r.paragraph(
				"The `%s` variable's value %s.",
			)(
				r.spec.Name(),
				req,
			)
		} else {
			r.paragraph(
				"The `%s` variable **MAY** be left undefined.",
				"Otherwise, the value %s.",
			)(
				r.spec.Name(),
				req,
			)
		}

		return
	}

	if def, ok := r.spec.Default(); ok {
		r.paragraph(
			"The `%s` variable **MAY** be left undefined, in which case the default value of `%s` is used.",
		)(
			r.spec.Name(),
			def.String,
		)
	} else if r.spec.IsRequired() {
		r.paragraph(
			"The `%s` variable **MUST NOT** be left undefined.",
		)(
			r.spec.Name(),
		)
	} else {
		r.paragraph(
			"The `%s` variable **MAY** be left undefined.",
		)(
			r.spec.Name(),
		)
	}
}
