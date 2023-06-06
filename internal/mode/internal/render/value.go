package render

import (
	"io"
	"strconv"
	"strings"

	"github.com/dogmatiq/ferrite/variable"
)

// Value returns a value rendered for display within documentation and other
// human-readable output.
func Value(s variable.Spec, v variable.Literal) string {
	vis := &valueRenderer{
		Spec: s,
		In:   v,
		Out:  &strings.Builder{},
	}

	s.Schema().AcceptVisitor(vis)

	return vis.Out.String()
}

// WriteValue renders a value for display within documentation and other
// human-readable output.
func WriteValue(w io.Writer, s variable.Spec, v variable.Literal) {

}

type valueRenderer struct {
	Spec variable.Spec
	In   variable.Literal
	Out  *strings.Builder
}

func (r *valueRenderer) VisitBinary(s variable.Binary) {
	n := len(r.In.String)

	r.Out.WriteByte('{')
	r.Out.WriteString(strconv.Itoa(n))
	r.Out.WriteString(" byte")

	if n != 1 {
		r.Out.WriteByte('s')
	}

	r.Out.WriteByte('}')
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
		n := len(r.In.String)
		r.Out.Grow(n)
		for i := 0; i < n; i++ {
			r.Out.WriteByte('*')
		}
	} else {
		r.Out.WriteString(r.In.Quote())
	}
}
