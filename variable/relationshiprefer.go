package variable

// RefersTo is a relationship that indicates that a variable refers to
// another variable for information/documentation purposes.
type RefersTo struct{ Spec, RefersTo Spec }

// Subject returns the spec of the variable that is the subject of the
// relationship.
func (r RefersTo) Subject() Spec {
	return r.Spec
}

// Related returns the spec of the related variable.
func (r RefersTo) Related() Spec {
	return r.RefersTo
}

// Inverse returns the inverse of r.
func (r RefersTo) Inverse() Relationship {
	return IsReferredToBy{r.RefersTo, r.Spec}
}

// AcceptVisitor passes r to the appropriate method of v.
func (r RefersTo) AcceptVisitor(v RelationshipVisitor) {
	v.VisitRefersTo(r)
}

// IsReferredToBy is a relationship that indicates that a variable is
// referred to by another variable.
type IsReferredToBy struct{ Spec, IsReferredToBy Spec }

// Subject returns the spec of the variable that is the subject of the
// relationship.
func (r IsReferredToBy) Subject() Spec {
	return r.Spec
}

// Related returns the spec of the related variable.
func (r IsReferredToBy) Related() Spec {
	return r.IsReferredToBy
}

// Inverse returns the inverse of r.
func (r IsReferredToBy) Inverse() Relationship {
	return RefersTo{r.IsReferredToBy, r.Spec}
}

// AcceptVisitor passes r to the appropriate method of v.
func (r IsReferredToBy) AcceptVisitor(v RelationshipVisitor) {
	v.VisitIsReferredToBy(r)
}
