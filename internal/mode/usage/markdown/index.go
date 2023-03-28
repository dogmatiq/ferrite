package markdown

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/ferrite/variable"
)

func (r *renderer) renderIndex() {
	r.line("## Index")
	r.gap()

	for _, v := range r.Specs {
		r.renderIndexItem(v)
	}
}

func (r *renderer) renderIndexItem(s variable.Spec) {
	var w strings.Builder

	style := ""
	label := ""
	if s.IsDeprecated() {
		style = "~~"
		label = "deprecated"
	} else if def, ok := s.Default(); ok {
		label = fmt.Sprintf("defaults to `%s`", def.Quote())
	} else if !s.IsRequired() {
		label = "optional"
	} else if len(variable.Relationships[variable.DependsOn](s)) != 0 {
		label = "conditional"
	} else {
		style = "**"
		label = "required"
	}

	w.WriteString("- ")
	w.WriteString(style)
	w.WriteString(r.linkToSpec(s))
	w.WriteString(style)
	w.WriteString(" — ")
	w.WriteString(style)
	w.WriteString(s.Description())
	w.WriteString(style)
	w.WriteString(" (")
	w.WriteString(label)
	w.WriteString(")")

	r.line(w.String())
}
