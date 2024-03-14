package validate

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/ferrite/internal/variable"
)

func renderError(s variable.Spec, err variable.ValueError) string {
	out := &strings.Builder{}
	err.AcceptVisitor(&errorRenderer{
		Output: out,
		Schema: s.Schema(),
		Error:  err,
	})
	return out.String()
}

type errorRenderer struct {
	Output *strings.Builder
	Schema variable.Schema
	Error  variable.ValueError
}

func (r *errorRenderer) VisitGenericError(error) {
	r.Schema.AcceptVisitor(r)
}

func (r *errorRenderer) VisitBinary(variable.Binary) {
	r.Output.WriteString(r.Error.Unwrap().Error())
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

func (r *errorRenderer) VisitSet(variable.Set) {
	r.Output.WriteString(r.Error.Unwrap().Error())
}

func (r *errorRenderer) VisitSetMembershipError(err variable.SetMembershipError) {
	r.Output.WriteString(err.Error())
}

func (r *errorRenderer) VisitString(variable.String) {
	r.Output.WriteString(r.Error.Unwrap().Error())
}

func (r *errorRenderer) VisitMinLengthError(err variable.MinLengthError) {
	r.Output.WriteString(err.Error())
}

func (r *errorRenderer) VisitMaxLengthError(err variable.MaxLengthError) {
	r.Output.WriteString(err.Error())
}

func (r *errorRenderer) VisitOther(variable.Other) {
	r.Output.WriteString(r.Error.Unwrap().Error())
}
