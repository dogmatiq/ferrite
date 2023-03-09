package validatemode

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dogmatiq/ferrite/variable"
)

// Run validates the variables in the given registry.
//
// ok is true if all variables are valid.
//
// usage contains human-readable usage and validation information intended for
// display in the console.
func Run(reg *variable.Registry) (usage string, ok bool) {
	ok = true
	t := table{}

	for _, v := range reg.Variables() {
		s := v.Spec()

		t.AddRow(
			renderNameColumn(v),
			s.Description(),
			renderSpecColumn(s),
			renderValueColumn(v),
		)

		if !v.IsValid() {
			ok = false
		}
	}

	return "Environment Variables:\n\n" + t.String(), ok
}

func renderNameColumn(v variable.Any) string {
	if v.IsValid() {
		return fmt.Sprintf("   %s", v.Spec().Name())
	}

	return fmt.Sprintf(" %s %s", attentionIcon, v.Spec().Name())
}

func renderSpecColumn(s variable.Spec) string {
	out := &strings.Builder{}
	s.Schema().AcceptVisitor(&schemaRenderer{
		Output: out,
	})

	if def, ok := s.Default(); ok {
		return fmt.Sprintf(
			"[ %s ] = %s",
			out,
			renderValue(s, def),
		)
	}

	if s.IsRequired() {
		return fmt.Sprintf("  %s  ", out)
	}

	return fmt.Sprintf("[ %s ]", out)
}

func renderValueColumn(v variable.Any) string {
	spec := v.Spec()

	value, ok, err := v.Value()
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
			renderValue(spec, err.Literal()),
			out.String(),
		)
	}

	if ok {
		if value.IsDefault() {
			return fmt.Sprintf("%s using default value", validIcon)
		}

		if value.Verbatim() == value.Canonical() {
			return fmt.Sprintf(
				"%s set to %s",
				validIcon,
				renderValue(spec, value.Canonical()),
			)
		}

		return fmt.Sprintf(
			"%s set to %s (specified non-canonically as %s)",
			validIcon,
			renderValue(spec, value.Canonical()),
			renderValue(spec, value.Verbatim()),
		)
	}

	icon := neutralIcon
	if v.Spec().IsRequired() {
		icon = invalidIcon
	}

	return fmt.Sprintf("%s undefined", icon)
}

func renderValue(s variable.Spec, v variable.Literal) string {
	if s.IsSensitive() {
		return strings.Repeat("*", len(v.String))
	}

	return v.Quote()
}

const (
	validIcon     = "✓"
	invalidIcon   = "✗"
	neutralIcon   = "•"
	attentionIcon = "❯"
)

type schemaRenderer struct {
	Output *strings.Builder
}

func (r *schemaRenderer) VisitSet(s variable.Set) {
	for i, m := range s.Literals() {
		if i > 0 {
			r.Output.WriteString(" | ")
		}

		r.Output.WriteString(m.Quote())
	}
}

func (r *schemaRenderer) VisitNumeric(s variable.Numeric) {
	min, hasMin := s.Min()
	max, hasMax := s.Max()

	if hasMin && hasMax {
		fmt.Fprintf(
			r.Output,
			"%s .. %s",
			min.Quote(),
			max.Quote(),
		)
	} else if hasMin {
		fmt.Fprintf(
			r.Output,
			"%s ...",
			min.Quote(),
		)
	} else if hasMax {
		fmt.Fprintf(
			r.Output,
			"... %s",
			max.Quote(),
		)
	} else {
		fmt.Fprintf(
			r.Output,
			"<%s>",
			s.Type().Kind(),
		)
	}
}

func (r *schemaRenderer) VisitString(s variable.String) {
	fmt.Fprintf(r.Output, "<%s>", s.Type().Kind())
}

func (r *schemaRenderer) VisitOther(s variable.Other) {
	t := s.Type()

again:
	switch t.Kind() {
	case reflect.Pointer:
		t = t.Elem()
		goto again

	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		fmt.Fprintf(r.Output, "<%s>", t.Kind())

	default:
		r.Output.WriteString("<string>")
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
		fmt.Fprintf(
			r.Output,
			" between %s and %s",
			min.Quote(),
			max.Quote(),
		)
	}
}

func (r *errorRenderer) VisitMinError(err variable.MinError) {
	r.Output.WriteString(err.Error())
}

func (r *errorRenderer) VisitMaxError(err variable.MaxError) {
	r.Output.WriteString(err.Error())
}

func (r *errorRenderer) VisitString(s variable.String) {
	r.Output.WriteString(r.Error.Unwrap().Error())
}

func (r *errorRenderer) VisitMinLengthError(err variable.MinLengthError) {
	r.Output.WriteString(err.Error())
}

func (r *errorRenderer) VisitMaxLengthError(err variable.MaxLengthError) {
	r.Output.WriteString(err.Error())
}

func (r *errorRenderer) VisitOther(s variable.Other) {
	r.Output.WriteString(r.Error.Unwrap().Error())
}
