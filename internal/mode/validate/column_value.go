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
	if s.IsSensitive() {
		return strings.Repeat("*", len(v.String))
	}

	return v.Quote()
}
