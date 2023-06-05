package validate

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/ferrite/variable"
)

// value renders a column describing the variable's value.
func value(v variable.Any) string {
	s := v.Spec()

	renderExplicit := func(icon string, lit variable.Literal, message string) string {
		out := &strings.Builder{}

		out.WriteString(icon)
		out.WriteByte(' ')

		if s.IsDeprecated() {
			out.WriteString("deprecated variable ")
		}

		out.WriteString("set to ")
		out.WriteString(renderValue(s, lit))

		if message != "" {
			out.WriteString(", ")
			out.WriteString(message)
		}

		return out.String()
	}

	switch v.Source() {
	case variable.SourceNone:
		if s.IsRequired() {
			return fmt.Sprintf("%s undefined", iconError)
		}
		return fmt.Sprintf("%s undefined", iconNeutral)

	case variable.SourceDefault:
		return fmt.Sprintf("%s using default value", iconOK)

	default:
		if err, ok := v.Error().(variable.ValueError); ok {
			return renderExplicit(
				iconError,
				err.Literal(),
				renderError(s, err),
			)
		}

		icon := iconOK
		if s.IsDeprecated() {
			icon = iconWarn
		}

		value := v.Value()
		message := ""

		if value.Verbatim() != value.Canonical() {
			message = fmt.Sprintf(
				"equivalent to %s",
				renderValue(s, value.Canonical()),
			)
		}

		return renderExplicit(icon, value.Verbatim(), message)
	}
}

func renderValue(s variable.Spec, v variable.Literal) string {
	vis := &valueRenderer{
		Spec: s,
		In:   v,
	}
	s.Schema().AcceptVisitor(vis)

	return vis.Out
}

type valueRenderer struct {
	Spec variable.Spec
	In   variable.Literal
	Out  string
}

func (r *valueRenderer) VisitBinary(s variable.Binary) {
	r.Out = fmt.Sprintf("%d-byte value", len(r.In.String))
}

func (r *valueRenderer) VisitNumeric(s variable.Numeric) {
	r.visitGeneric(s)
}

func (r *valueRenderer) VisitSet(s variable.Set) {
	r.visitGeneric(s)
}

func (r *valueRenderer) VisitString(s variable.String) {
	r.visitGeneric(s)
}

func (r *valueRenderer) VisitOther(s variable.Other) {
	r.visitGeneric(s)
}

func (r *valueRenderer) visitGeneric(s variable.Schema) {
	if r.Spec.IsSensitive() {
		r.Out = strings.Repeat("*", len(r.In.String))
	} else {
		r.Out = r.In.Quote()
	}
}
