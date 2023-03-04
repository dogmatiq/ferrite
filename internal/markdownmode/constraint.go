package markdownmode

import "github.com/dogmatiq/ferrite/variable"

func (r *renderer) renderConstraints(s variable.Spec) {
	con := s.Constraints()
	if len(con) == 0 {
		return
	}

	if _, ok := s.Default(); !ok && s.IsRequired() {
		if len(con) == 1 {
			r.renderConstraintSentence("The value %s.", con[0])
		} else {
			r.line("The value:")
			r.renderConstraintList(con)
		}
	} else {
		if len(con) == 1 {
			r.renderConstraintSentence("Otherwise, the value %s.", con[0])
		} else {
			r.line("Otherwise, the value:")
			r.renderConstraintList(con)
		}
	}

}

func (r *renderer) renderConstraintSentence(f string, c variable.Constraint) {
	r.line(f, asMarkdown(c.Description()))
}

func (r *renderer) renderConstraintList(constraints []variable.Constraint) {
	for i, c := range constraints {
		r.line("%d. %s", i+1, asMarkdown(c.Description()))
	}
}
