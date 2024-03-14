package markdown

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/ferrite/internal/variable"
)

type specRenderer struct {
	ren  *renderer
	spec variable.Spec
	reg  *variable.Registry
}

func (r *specRenderer) Render() {
	r.ren.gap()
	r.ren.line("### `%s`", r.spec.Name())

	r.ren.gap()
	r.ren.line("> %s", r.spec.Description())

	r.spec.Schema().AcceptVisitor(r)

	r.renderImportantDocumentation()

	if r.spec.IsSensitive() {
		r.ren.paragraphf(
			"⚠️ This variable is **sensitive**;",
			"its value may contain private information.",
		)()
	} else {
		r.renderExamples()
	}

	r.renderRegistry()
	r.renderUnimportantDocumentation()
	r.renderSeeAlso()
}

// VisitBinary renders the primary requirement for a spec that uses the
// "binary" schema type.
func (r *specRenderer) VisitBinary(s variable.Binary) {
	if con := r.bestConstraint(); con != nil {
		r.renderPrimaryRequirement(con.Description())
	} else {
		r.renderPrimaryRequirement("**MUST** be a binary value expressed using the `%s` encoding scheme", s.EncodingDescription())
	}
}

// VisitNumeric renders the primary requirement for a spec that uses the
// "numeric" schema type.
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

// VisitSet renders the primary requirement for a spec that uses the "set"
// schema type.
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
	if con := r.bestConstraint(); con != nil {
		r.renderPrimaryRequirement(con.Description())
	} else {
		r.renderPrimaryRequirement("")
	}
}

// VisitOther render the primary requirement for a spec that uses the "other"
// schema type.
func (r *specRenderer) VisitOther(variable.Other) {
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

	if r.spec.IsDeprecated() {
		r.renderPrimaryRequirementDeprecated(req)
	} else if def, ok := r.renderDefaultValueFragment(); ok {
		r.renderPrimaryRequirementDefault(req, def)
	} else if r.spec.IsRequired() {
		r.renderPrimaryRequirementRequired(req)
	} else {
		r.renderPrimaryRequirementOptional(req)
	}
}

func (r *specRenderer) renderPrimaryRequirementDefault(req, def string) {
	r.ren.paragraph(
		func(write func(string, ...any)) {
			write(
				"The `%s` variable **MAY** be left undefined, in which case %s is used.",
				r.spec.Name(),
				def,
			)

			if req != "" {
				write(" Otherwise, the value %s.", req)
			}

			if dependsOn := r.renderDependsOnClause(); dependsOn != "" {
				write(" The value is not used when %s.", dependsOn)
			}
		},
	)
}

func (r *specRenderer) renderPrimaryRequirementRequired(req string) {
	r.ren.paragraph(
		func(write func(string, ...any)) {
			length := r.renderLengthClause(r.spec.Schema(), req == "")

			if dependsOn := r.renderDependsOnClause(); dependsOn != "" {
				write(
					"The `%s` variable **MAY** be left undefined if and only if %s.",
					r.spec.Name(),
					dependsOn,
				)

				if req != "" && length != "" {
					write(" Otherwise, the value %s with %s.", req, length)
				} else if req != "" {
					write(" Otherwise, the value %s.", req)
				} else if length != "" {
					write(" Otherwise, the value **MUST** have %s.", length)
				}

				return
			}

			if req != "" && length != "" {
				write("The `%s` variable's value %s with %s.", r.spec.Name(), req, length)
			} else if req != "" {
				write("The `%s` variable's value %s.", r.spec.Name(), req)
			} else if length != "" {
				write("The `%s` variable's value **MUST** have %s.", r.spec.Name(), length)
			} else {
				write("The `%s` variable **MUST NOT** be left undefined.", r.spec.Name())
			}
		},
	)
}

func (r *specRenderer) renderPrimaryRequirementOptional(req string) {
	r.ren.paragraph(
		func(write func(string, ...any)) {
			write(
				"The `%s` variable **MAY** be left undefined.",
				r.spec.Name(),
			)

			length := r.renderLengthClause(r.spec.Schema(), false)

			if req != "" && length != "" {
				write(" Otherwise, the value %s with %s.", req, length)
			} else if req != "" {
				write(" Otherwise, the value %s.", req)
			} else if length != "" {
				write(" Otherwise, the value **MUST** have %s.", length)
			}

			if dependsOn := r.renderDependsOnClause(); dependsOn != "" {
				write(" The value is not used when %s.", dependsOn)
			}
		},
	)
}

func (r *specRenderer) renderPrimaryRequirementDeprecated(req string) {
	r.ren.paragraph(
		func(write func(string, ...any)) {
			write(
				"⚠️ The `%s` variable is **deprecated**; its use is **NOT RECOMMENDED** as it may be removed in a future version.",
				r.spec.Name(),
			)

			relationships := variable.InverseRelationships[variable.Supersedes](r.spec)
			if len(relationships) != 0 {
				write(
					" %s **SHOULD** be used instead.",
					andList(
						relationships,
						func(rel variable.Supersedes) string {
							return r.ren.linkToSpec(rel.Subject)
						},
					),
				)
			}

			length := r.renderLengthClause(r.spec.Schema(), false)

			if req != "" && length != "" {
				write(" If defined, the value %s with %s.", req, length)
			} else if req != "" {
				write(" If defined, the value %s.", req)
			} else if length != "" {
				write(" If defined, the value **MUST** have %s.", length)
			}

			if dependsOn := r.renderDependsOnClause(); dependsOn != "" {
				write(" The value is not used when %s.", dependsOn)
			}
		},
	)
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

func (r *specRenderer) renderDependsOnClause() string {
	return orList(
		variable.Relationships[variable.DependsOn](r.spec),
		func(rel variable.DependsOn) string {
			return fmt.Sprintf(
				"%s is `%s`",
				r.ren.linkToSpec(rel.DependsOn),
				rel.DependsOn.Zero().String,
			)
		},
	)
}

// bestConstraint returns the constraint to use as the "primary"
// requirement, favoring non-user-defined constraints.
func (r *specRenderer) bestConstraint() (con variable.Constraint) {
	constraints := r.spec.Constraints()

	for _, c := range constraints {
		if !c.IsUserDefined() {
			return c
		}
	}

	for _, c := range constraints {
		return c
	}

	return nil
}

func (r *specRenderer) renderLengthClause(s variable.Schema, implicitMin bool) string {
	lim, ok := s.(variable.LengthLimited)
	if !ok {
		return ""
	}

	min, hasMin := lim.MinLength()
	max, hasMax := lim.MaxLength()

	// In cases where there's not much else to say about the valid value, we
	// render an "implicit minimum" length of 1.
	//
	// This results in wording such as:
	// 	"The length must be between 1 and 10 bytes."
	// as opposed to:
	// 	"The variable must be defined with a length less than 10 bytes."
	if hasMax && implicitMin {
		if min < 1 {
			min = 1
		}
		hasMin = true
	}

	if hasMin && hasMax {
		if min == max {
			return fmt.Sprintf(
				"%s of exactly %d %s",
				lim.LengthDescription(),
				min,
				inflect("byte", min),
			)
		}

		return fmt.Sprintf(
			"%s between %d and %d bytes",
			lim.LengthDescription(),
			min,
			max,
		)
	}

	if hasMin {
		return fmt.Sprintf(
			"%s of at least %d %s",
			lim.LengthDescription(),
			min,
			inflect("byte", min),
		)
	}

	if hasMax {
		return fmt.Sprintf(
			"%s of %d %s or fewer",
			lim.LengthDescription(),
			max,
			inflect("byte", max),
		)
	}

	return ""
}
