package markdown

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

	if r.spec.IsSensitive() {
		r.ren.paragraph(
			"⚠️ This variable is **sensitive**;",
			"its value may contain private information.",
		)()
	} else {
		r.renderExamples()
	}

	r.renderUnimportantDocumentation()
	r.renderSeeAlso()
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
	// Find the best constraint to use as the "primary" requirement, favoring
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

	var text string
	var args []any

	if r.spec.IsDeprecated() {
		text = "⚠️ The `%s` variable is **deprecated**; its use is **NOT RECOMMENDED** as it may be removed in a future version."
		args = []any{r.spec.Name()}

		relationships := variable.FilterRelationships[variable.IsSupersededBy](r.spec)
		if len(relationships) != 0 {
			for i, rel := range relationships {
				if i == len(relationships)-1 {
					text += " and"
				} else if i > 0 {
					text += ","
				}

				text += " " + r.ren.linkToSpec(rel.SupersededBy)
			}

			text += " **SHOULD** be used instead."
		}

		if req != "" {
			text += " If defined, the value %s."
			args = append(args, req)
		}
	} else if def, ok := r.renderDefaultValueFragment(); ok {
		text = "The `%s` variable **MAY** be left undefined, in which case %s is used."
		args = []any{r.spec.Name(), def}

		if req != "" {
			text += " Otherwise, the value %s."
			args = append(args, req)
		}
	} else if r.spec.IsRequired() {
		if req != "" {
			text = "The `%s` variable's value %s."
			args = []any{r.spec.Name(), req}
		} else {
			text = "The `%s` variable **MUST NOT** be left undefined."
			args = []any{r.spec.Name()}
		}
	} else {
		text = "The `%s` variable **MAY** be left undefined."
		args = []any{r.spec.Name()}

		if req != "" {
			text += " Otherwise, the value %s."
			args = append(args, req)
		}
	}

	r.ren.paragraph(text)(args...)
}

func (r *specRenderer) renderDefaultValueFragment() (string, bool) {
	def, ok := r.spec.Default()
	if !ok {
		return "", false
	}

	if r.spec.IsSensitive() {
		return "a default value", true
	}

	return fmt.Sprintf("the default value of `%s`", def.String), true
}
