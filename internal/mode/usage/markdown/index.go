package markdown

import (
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
	strike := ""
	if s.IsDeprecated() {
		strike = "~~"
	}

	w.WriteString("- ")
	w.WriteString(strike)
	w.WriteString(r.linkToSpec(s))
	w.WriteString(strike)
	w.WriteString(" — ")
	w.WriteString(strike)
	w.WriteString(s.Description())
	w.WriteString(strike)

	if s.IsDeprecated() {
		w.WriteString(" (deprecated)")
	}

	r.line(w.String())
}
