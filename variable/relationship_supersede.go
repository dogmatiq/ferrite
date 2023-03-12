package variable

// Supersedes is a relationship that indicates that a variable supersedes
// another (usually deprecated) variable.
type Supersedes struct{ Spec, Supersedes Spec }

// Subject returns the spec of the variable that is the subject of the
// relationship.
func (r Supersedes) Subject() Spec {
	return r.Spec
}

// Related returns the spec of the related variable.
func (r Supersedes) Related() Spec {
	return r.Supersedes
}

// Inverse returns the inverse of r.
func (r Supersedes) Inverse() Relationship {
	return IsSupersededBy{r.Supersedes, r.Spec}
}

// AcceptVisitor passes r to the appropriate method of v.
func (r Supersedes) AcceptVisitor(v RelationshipVisitor) {
	v.VisitSupersedes(r)
}

// IsSupersededBy is a relationship that indicates that a (usually deprecated)
// variable is superseded by another variable.
type IsSupersededBy struct{ Spec, SupersededBy Spec }

// Subject returns the spec of the variable that is the subject of the
// relationship.
func (r IsSupersededBy) Subject() Spec {
	return r.Spec
}

// Related returns the spec of the related variable.
func (r IsSupersededBy) Related() Spec {
	return r.SupersededBy
}

// Inverse returns the inverse of r.
func (r IsSupersededBy) Inverse() Relationship {
	return Supersedes{r.SupersededBy, r.Spec}
}

// AcceptVisitor passes r to the appropriate method of v.
func (r IsSupersededBy) AcceptVisitor(v RelationshipVisitor) {
	v.VisitIsSupersededBy(r)
}
