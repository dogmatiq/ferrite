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
	return SupersededBy{r.Supersedes, r.Spec}
}

// AcceptVisitor passes r to the appropriate method of v.
func (r Supersedes) AcceptVisitor(v RelationshipVisitor) {
	v.VisitSupersedes(r)
}

// SupersededBy is a relationship that indicates that a (usually deprecated)
// variable is superseded by another variable.
type SupersededBy struct{ Spec, SupersededBy Spec }

// Subject returns the spec of the variable that is the subject of the
// relationship.
func (r SupersededBy) Subject() Spec {
	return r.Spec
}

// Related returns the spec of the related variable.
func (r SupersededBy) Related() Spec {
	return r.SupersededBy
}

// Inverse returns the inverse of r.
func (r SupersededBy) Inverse() Relationship {
	return Supersedes{r.SupersededBy, r.Spec}
}

// ConflictsWith returns true if this relationship conflicts with r.
//
// The subject of both relationships must be the same specification.
func (r SupersededBy) ConflictsWith(c Relationship) bool {
	if c, ok := c.(Supersedes); ok {
		return c.Supersedes == r.SupersededBy
	}
	return false
}

// AcceptVisitor passes r to the appropriate method of v.
func (r SupersededBy) AcceptVisitor(v RelationshipVisitor) {
	v.VisitSupersededBy(r)
}
