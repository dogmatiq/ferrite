package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dogmatiq/ferrite/internal/mode/internal/render"
	"github.com/dogmatiq/ferrite/internal/variable"
)

// spec renders a column describing the variable's specification and schema.
func spec(v variable.Any) string {
	s := v.Spec()

	out := &strings.Builder{}
	s.Schema().AcceptVisitor(&schemaRenderer{
		Output: out,
	})

	if def, ok := s.Default(); ok {
		return fmt.Sprintf(
			"[ %s ] = %s",
			out,
			render.Value(s, def),
		)
	}

	if s.IsRequired() {
		return fmt.Sprintf("  %s  ", out)
	}

	return fmt.Sprintf("[ %s ]", out)
}

type schemaRenderer struct {
	Output *strings.Builder
}

func (r *schemaRenderer) VisitBinary(s variable.Binary) {
	r.Output.WriteByte('<')
	r.Output.WriteString(s.EncodingDescription())
	r.Output.WriteByte('>')
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

func (r *schemaRenderer) VisitSet(s variable.Set) {
	for i, m := range s.Literals() {
		if i > 0 {
			r.Output.WriteString(" | ")
		}

		r.Output.WriteString(m.Quote())
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
