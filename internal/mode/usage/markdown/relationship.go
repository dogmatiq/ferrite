package markdown

import (
	"github.com/dogmatiq/ferrite/variable"
)

func (r *specRenderer) renderSeeAlso() {
	relationships := variable.FilterRelationships[variable.RefersTo](r.spec)
	if len(relationships) == 0 {
		return
	}

	r.ren.gap()
	r.ren.line("#### See Also")
	r.ren.gap()

	for _, rel := range relationships {
		r.ren.renderIndexItem(rel.Related())
	}
}
