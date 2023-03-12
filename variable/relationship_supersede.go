package variable

// Supersedes is a relationship that indicates that a variable supersedes
// another (usually deprecated) variable.
type Supersedes struct {
	Sub, Supersedes Spec
}

// Subject returns the spec of the variable that is the subject of the
// relationship.
func (r Supersedes) Subject() Spec {
	return r.Sub
}

// Related returns the spec of the related variable.
func (r Supersedes) Related() Spec {
	return r.Supersedes
}

// Inverse returns the inverse of r.
func (r Supersedes) Inverse() Relationship {
	return SupersededBy{r.Supersedes, r.Sub}
}

// ConflictsWith returns true if this relationship conflicts with r.
//
// The subject of both relationships must be the same specification.
func (r Supersedes) ConflictsWith(c Relationship) bool {
	if c, ok := c.(SupersededBy); ok {
		return c.SupersededBy == r.Supersedes
	}
	return false
}

// SupersededBy is a relationship that indicates that a (usually deprecated)
// variable is superseded by another variable.
type SupersededBy struct {
	Sub, SupersededBy Spec
}

// Subject returns the spec of the variable that is the subject of the
// relationship.
func (r SupersededBy) Subject() Spec {
	return r.Sub
}

// Related returns the spec of the related variable.
func (r SupersededBy) Related() Spec {
	return r.SupersededBy
}

// Inverse returns the inverse of r.
func (r SupersededBy) Inverse() Relationship {
	return Supersedes{r.SupersededBy, r.Sub}
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
