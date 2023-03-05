package markdownmode

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/ferrite/variable"
)

type specRenderer struct {
	ren  *renderer
	spec variable.Spec
}

func (r *specRenderer) Render() {
	r.ren.gap()
	r.ren.line("### `%s`", r.spec.Name())

	r.ren.gap()
	r.ren.line("> %s", r.spec.Description())

	r.spec.Schema().AcceptVisitor(r)

	r.renderImportantDocumentation()
	r.renderExamples()
	r.renderUnimportantDocumentation()
}

// VisitNumeric renders the primary requirement for spec that uses the "numeric"
// schema type.
func (r *specRenderer) VisitNumeric(s variable.Numeric) {
	min, hasMin := s.Min()
	max, hasMax := s.Max()

	if hasMin && hasMax {
		r.renderPrimaryRequirement("**MUST** be between `%s` and `%s`", min.String, max.String)
	} else if hasMin {
		r.renderPrimaryRequirement("**MUST** be `%s` or greater", min.String)
	} else if hasMax {
		r.renderPrimaryRequirement("**MUST** be `%s` or less", max.String)
	} else {
		switch s.Type().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			r.renderPrimaryRequirement("**MUST** be a whole number")
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			r.renderPrimaryRequirement("**MUST** be a non-negative whole number")
		case reflect.Float32, reflect.Float64:
			r.renderPrimaryRequirement("**MUST** be a number with an **OPTIONAL** fractional part")
		}
	}
}

// VisitSet renders the primary requirement for spec that uses the "set" schema
// type.
func (r *specRenderer) VisitSet(s variable.Set) {
	if lits := s.Literals(); len(lits) == 2 {
		r.renderPrimaryRequirement(
			"**MUST** be either `%s` or `%s`",
			lits[0].String,
			lits[1].String,
		)
	} else {
		r.renderPrimaryRequirement("**MUST** be one of the values shown in the examples below")
	}
}

// VisitString renders the primary requirement for a spec that uses the "string"
// schema type.
func (r *specRenderer) VisitString(variable.String) {
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
		r.renderPrimaryRequirement("")
	} else {
		r.renderPrimaryRequirement(con.Description())
	}
}

// VisitOther render the primary requirement for a spec that uses the "other"
// schema type.
func (r *specRenderer) VisitOther(s variable.Other) {
	con := ""
	for _, c := range r.spec.Constraints() {
		con = c.Description()
		break
	}

	r.renderPrimaryRequirement(con)
}

// renderPrimaryRequirement renders information about the most important
// requirement of the variable's schema, this includes information about whether
// the variable is optional and the basic data type of the variable.
//
// The (optional) requirement text must complete the phrase "the value...". It
// should not include any trailing punctuation.
func (r *specRenderer) renderPrimaryRequirement(f string, v ...any) {
	req := fmt.Sprintf(f, v...)

	if req != "" {
		if def, ok := r.spec.Default(); ok {
			r.ren.paragraph(
				"The `%s` variable **MAY** be left undefined, in which case the default value of `%s` is used.",
				"Otherwise, the value %s.",
			)(
				r.spec.Name(),
				def.String,
				req,
			)
		} else if r.spec.IsRequired() {
			r.ren.paragraph(
				"The `%s` variable's value %s.",
			)(
				r.spec.Name(),
				req,
			)
		} else {
			r.ren.paragraph(
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
		r.ren.paragraph(
			"The `%s` variable **MAY** be left undefined, in which case the default value of `%s` is used.",
		)(
			r.spec.Name(),
			def.String,
		)
	} else if r.spec.IsRequired() {
		r.ren.paragraph(
			"The `%s` variable **MUST NOT** be left undefined.",
		)(
			r.spec.Name(),
		)
	} else {
		r.ren.paragraph(
			"The `%s` variable **MAY** be left undefined.",
		)(
			r.spec.Name(),
		)
	}
}
