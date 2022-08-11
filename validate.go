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
			return a[0][3:] < b[0][3:]
		},
	}

	for _, v := range variable.DefaultRegistry.Variables() {
		if _, err := v.Canonical(); err != nil {
			ok = false
		}

		t.AddRow(
			renderName(v),
			v.Description(),
			renderSchema(v),
			renderValue(v),
		)
	}

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

func renderName(v variable.Variable) string {
	if _, err := v.Canonical(); err != nil {
		return fmt.Sprintf(" %s %s", attentionIcon, v.Name())
	}

	return fmt.Sprintf("   %s", v.Name())
}

func renderSchema(v variable.Variable) string {
	var r classRenderer
	v.Class().AcceptVisitor(&r)

	if def, ok := v.Default().Get(); ok {
		return fmt.Sprintf(
			"[ %s ] = %s",
			r.Output.String(),
			def,
		)
	}

	if v.IsOptional() {
		return fmt.Sprintf(
			"[ %s ]",
			r.Output.String(),
		)
	}

	return fmt.Sprintf(
		"  %s  ",
		r.Output.String(),
	)
}

func renderValue(v variable.Variable) string {
	m, err := v.Canonical()
	if err != nil {
		return fmt.Sprintf(
			"%s set to %s, %s",
			invalidIcon,
			v.Verbatim(),
			err.Reason(),
		)
	}

	if s, ok := m.Get(); ok {
		if v.IsDefault() {
			return fmt.Sprintf("%s using default value", validIcon)
		}

		return fmt.Sprintf(
			"%s set to %s",
			validIcon,
			s,
		)
	}

	if v.IsOptional() {
		return fmt.Sprintf("%s undefined", neutralIcon)
	}

	return fmt.Sprintf("%s undefined", invalidIcon)
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

type classRenderer struct {
	Output strings.Builder
}

func (r *classRenderer) VisitSet(s variable.Set) {
	for i, m := range s.Members() {
		if i > 0 {
			r.Output.WriteString(" | ")
		}

		r.Output.WriteString(m.String())
	}
}
