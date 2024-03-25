package markdown

import (
	"strings"

	"github.com/dogmatiq/ferrite/internal/variable"
)

func (r *specRenderer) renderSeeAlso() {
	relationships := variable.Relationships[variable.RefersTo](r.spec)
	if len(relationships) == 0 {
		return
	}

	r.ren.gap()
	r.ren.line("### See Also")
	r.ren.gap()

	for _, rel := range relationships {
		r.ren.renderSeeAlsoItem(rel.RefersTo)
	}
}

func (r *renderer) renderSeeAlsoItem(s variable.Spec) {
	var w strings.Builder

	style := ""
	if s.IsDeprecated() {
		style = "~~"
	}

	w.WriteString("- ")
	w.WriteString(style)
	w.WriteString(r.linkToSpec(s))
	w.WriteString(style)
	w.WriteString(" — ")
	w.WriteString(style)
	w.WriteString(s.Description())
	w.WriteString(style)

	if s.IsDeprecated() {
		w.WriteString(" (deprecated)")
	}

	r.line(w.String())
}
