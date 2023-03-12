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

// ConflictsWith returns true if this relationship conflicts with r.
//
// The subject of both relationships must be the same specification.
func (r RefersTo) ConflictsWith(c Relationship) bool {
	return false
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

// ConflictsWith returns true if this relationship conflicts with r.
//
// The subject of both relationships must be the same specification.
func (r IsReferredToBy) ConflictsWith(c Relationship) bool {
	return false
}
