package ferrite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dogmatiq/ferrite/internal/table"
	"github.com/dogmatiq/ferrite/spec"
	"github.com/dogmatiq/ferrite/variable"
)

// validate parses and validates all environment variables.
func validate() (string, bool) {
	ok := true
	t := table.Table{
		Less: func(a, b []string) bool {
			ra := []rune(a[0])
			rb := []rune(b[0])

			// Trim the indicator icon/indenting from the start to sort by name.
			return string(ra[3:]) < string(rb[3:])
		},
	}

	variable.DefaultRegistry.Range(
		func(v variable.Any) bool {
			t.AddRow(
				renderName(v),
				v.Spec().Description(),
				renderSpec(v.Spec()),
				renderValue(v),
			)

			if !v.IsValid() {
				ok = false
			}

			return true
		},
	)

	for _, r := range spec.SortedResolvers() {
		s := r.Spec()
		v, err := r.Resolve()

		value, valid := renderValueXXX(s, v, err)
		if !valid {
			ok = false
		}

		t.AddRow(
			renderNameXXX(s, valid),
			s.Description,
			renderSchemaXXX(s),
			value,
		)
	}

	return "Environment Variables:\n\n" + t.String(), ok
}

func renderName(v variable.Any) string {
	if v.IsValid() {
		return fmt.Sprintf("   %s", v.Spec().Name())
	}

	return fmt.Sprintf(" %s %s", attentionIcon, v.Spec().Name())
}

func renderSpec(s variable.Spec) string {
	out := &strings.Builder{}
	s.Schema().AcceptVisitor(&schemaRenderer{
		Output: out,
	})

	if def, ok := s.Default().Get(); ok {
		return fmt.Sprintf("[ %s ] = %s", out, def)
	}

	if s.IsRequired() {
		return fmt.Sprintf("  %s  ", out)
	}

	return fmt.Sprintf("[ %s ]", out)
}

func renderValue(v variable.Any) string {
	m, err := v.Value()
	if err != nil {
		out := &strings.Builder{}
		err.AcceptVisitor(&errorRenderer{
			Output: out,
			Schema: v.Spec().Schema(),
			Error:  err,
		})

		return fmt.Sprintf(
			"%s set to %s, %s",
			invalidIcon,
			err.Literal(),
			out.String(),
		)
	}

	if value, ok := m.Get(); ok {
		if value.IsDefault() {
			return fmt.Sprintf("%s using default value", validIcon)
		}

		if value.Verbatim() == value.Canonical() {
			return fmt.Sprintf(
				"%s set to %s",
				validIcon,
				value.Canonical(),
			)
		}

		return fmt.Sprintf(
			"%s set to %s (specified non-canonically as %s)",
			validIcon,
			value.Canonical(),
			value.Verbatim(),
		)
	}

	icon := neutralIcon
	if v.Spec().IsRequired() {
		icon = invalidIcon
	}

	return fmt.Sprintf("%s undefined", icon)
}

func renderNameXXX(s spec.Spec, valid bool) string {
	if valid {
		return fmt.Sprintf("   %s", s.Name)
	}

	return fmt.Sprintf(" %s %s", attentionIcon, s.Name)
}

func renderSchemaXXX(s spec.Spec) string {
	renderer := &validateSchemaRendererXXX{}
	s.Schema.AcceptVisitor(renderer)

	if s.HasDefault {
		return fmt.Sprintf(
			"[ %s ] = %s",
			renderer.Output.String(),
			spec.Escape(s.Default),
		)
	}

	if s.IsOptional {
		return fmt.Sprintf(
			"[ %s ]",
			renderer.Output.String(),
		)
	}

	return fmt.Sprintf(
		"  %s  ",
		renderer.Output.String(),
	)
}

func renderValueXXX(s spec.Spec, v spec.Value, err error) (string, bool) {
	if err != nil {
		var invalid spec.ValidationError
		if errors.As(err, &invalid) {
			return fmt.Sprintf(
				"%s set to %s, %s",
				invalidIcon,
				spec.Escape(invalid.Value),
				invalid.Cause,
			), false
		}

		if !errors.As(err, &spec.UndefinedError{}) {
			return fmt.Sprintf("%s %s", invalidIcon, err), false
		}

		if s.HasDefault || s.IsOptional {
			return fmt.Sprintf("%s undefined", neutralIcon), true
		}

		return fmt.Sprintf("%s undefined", invalidIcon), false
	}

	if v.IsDefault() {
		return fmt.Sprintf("%s using default value", validIcon), true
	}

	return fmt.Sprintf(
		"%s set to %s",
		validIcon,
		spec.Escape(v.String()),
	), true
}

const (
	validIcon     = "✓"
	invalidIcon   = "✗"
	neutralIcon   = "•"
	attentionIcon = "❯"
)

type validateSchemaRendererXXX struct {
	Output strings.Builder
}

func (r *validateSchemaRendererXXX) VisitOneOf(s spec.OneOf) {
	for i, c := range s {
		if i > 0 {
			r.Output.WriteString(" | ")
		}

		c.AcceptVisitor(r)
	}
}

func (r *validateSchemaRendererXXX) VisitLiteral(s spec.Literal) {
	fmt.Fprintf(&r.Output, "%s", s)
}

func (r *validateSchemaRendererXXX) VisitType(s spec.Type) {
	fmt.Fprintf(&r.Output, "<%s>", s.Type)
}

func (r *validateSchemaRendererXXX) VisitRange(s spec.Range) {
	if s.Min != "" && s.Max != "" {
		fmt.Fprintf(&r.Output, "%s .. %s", s.Min, s.Max)
	} else if s.Max != "" {
		fmt.Fprintf(&r.Output, "... %s", s.Max)
	} else {
		fmt.Fprintf(&r.Output, "%s ...", s.Min)
	}
}

type schemaRenderer struct {
	Output *strings.Builder
}

func (r *schemaRenderer) VisitSet(s variable.Set) {
	for i, m := range s.Literals() {
		if i > 0 {
			r.Output.WriteString(" | ")
		}

		r.Output.WriteString(m.String())
	}
}

func (r *schemaRenderer) VisitNumeric(s variable.Numeric) {
	min, hasMin := s.Min().Get()
	max, hasMax := s.Max().Get()

	if hasMin && hasMax {
		fmt.Fprintf(r.Output, "%s .. %s", min, max)
	} else if hasMin {
		fmt.Fprintf(r.Output, "%s ...", min)
	} else if hasMax {
		fmt.Fprintf(r.Output, "... %s", max)
	} else {
		fmt.Fprintf(r.Output, "<%s>", s.Type().Name())
	}
}

type errorRenderer struct {
	Output *strings.Builder
	Schema variable.Schema
	Error  variable.ValueError
}

func (r *errorRenderer) VisitGenericError(err error) {
	r.Schema.AcceptVisitor(r)
}

func (r *errorRenderer) VisitSet(s variable.Set) {
	r.Output.WriteString(r.Error.Unwrap().Error())
}

func (r *errorRenderer) VisitSetMembershipError(err variable.SetMembershipError) {
	r.Output.WriteString(err.Error())
}

func (r *errorRenderer) VisitNumeric(s variable.Numeric) {
	typeName := strings.ToLower(s.Type().Name())

	if s.Type().PkgPath() == "" {
		if strings.Contains(typeName, "int") {
			typeName = "integer"
		}
	}

	fmt.Fprintf(r.Output, "expected %s", typeName)

	const maxHumanReadableBits = 16
	min, max, explicit := s.Limits()
	if explicit || s.Bits() <= maxHumanReadableBits {
		fmt.Fprintf(r.Output, " between %s and %s", min, max)
	}
}

func (r *errorRenderer) VisitMinError(err variable.MinError) {
	r.Output.WriteString(err.Error())
}

func (r *errorRenderer) VisitMaxError(err variable.MaxError) {
	r.Output.WriteString(err.Error())
}
