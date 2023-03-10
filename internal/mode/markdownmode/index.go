package markdownmode

import (
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
	r.line(
		"- [`%s`](#%s) — %s",
		s.Name(),
		s.Name(),
		s.Description(),
	)
}
