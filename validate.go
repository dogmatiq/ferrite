package ferrite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dogmatiq/ferrite/internal/table"
	"github.com/dogmatiq/ferrite/spec"
)

// validate parses and validates all environment variables.
func validate() (string, bool) {
	var (
		t  table.Table
		ok = true
	)

	for _, r := range spec.SortedResolvers() {
		s := r.Spec()
		v, err := r.Resolve()

		value, valid := renderValue(s, v, err)
		if !valid {
			ok = false
		}

		t.AddRow(
			renderName(s, valid),
			s.Description,
			renderSchema(s),
			value,
		)
	}

	return "Environment Variables:\n\n" + t.String(), ok
}

func renderName(s spec.Spec, valid bool) string {
	if valid {
		return fmt.Sprintf("   %s", s.Name)
	}

	return fmt.Sprintf(" %s %s", attentionIcon, s.Name)
}

func renderSchema(s spec.Spec) string {
	renderer := &validateSchemaRenderer{}
	s.Schema.AcceptVisitor(renderer)

	if s.HasDefault {
		return fmt.Sprintf(
			"[ %s ] = %s",
			renderer.Output.String(),
			spec.Escape(s.DefaultX),
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

func renderValue(s spec.Spec, v spec.Value, err error) (string, bool) {
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

type validateSchemaRenderer struct {
	Output strings.Builder
}

func (r *validateSchemaRenderer) VisitOneOf(s spec.OneOf) {
	for i, c := range s {
		if i > 0 {
			r.Output.WriteString(" | ")
		}

		c.AcceptVisitor(r)
	}
}

func (r *validateSchemaRenderer) VisitLiteral(s spec.Literal) {
	fmt.Fprintf(&r.Output, "%s", s)
}

func (r *validateSchemaRenderer) VisitType(s spec.Type) {
	fmt.Fprintf(&r.Output, "<%s>", s.Type)
}

func (r *validateSchemaRenderer) VisitRange(s spec.Range) {
	if s.Min != "" && s.Max != "" {
		fmt.Fprintf(&r.Output, "%s .. %s", s.Min, s.Max)
	} else if s.Max != "" {
		fmt.Fprintf(&r.Output, "... %s", s.Max)
	} else {
		fmt.Fprintf(&r.Output, "%s ...", s.Min)
	}
}
