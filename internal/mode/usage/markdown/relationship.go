package markdown

import (
	"github.com/dogmatiq/ferrite/variable"
	"golang.org/x/exp/slices"
)

func (r *specRenderer) renderRelationships() {
	c := &relationshipCollator{
		spec: r.spec,
	}

	for _, rel := range r.spec.Relationships() {
		rel.AcceptVisitor(c)
	}

	if len(c.seeAlso) != 0 {
		r.ren.gap()
		r.ren.line("#### See Also")
		r.ren.gap()

		slices.SortFunc(
			c.seeAlso,
			func(a, b variable.SeeAlso) bool {
				return a.To.Name() < b.To.Name()
			},
		)

		for _, rel := range c.seeAlso {
			r.ren.renderIndexItem(rel.To)
		}
	}

}

type relationshipCollator struct {
	spec    variable.Spec
	seeAlso []variable.SeeAlso
}

func (r *relationshipCollator) VisitSeeAlso(rel variable.SeeAlso) {
	if r.spec == rel.From {
		r.seeAlso = append(r.seeAlso, rel)
	}
}
