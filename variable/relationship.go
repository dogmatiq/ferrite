package variable

// Relationship represents a relationship between two variables.
type Relationship interface {
	// Specs returns the specs for the variables that are involved in the
	// relationship.
	Specs() []Spec

	// AcceptVisitor passes the schema to the appropriate method of v.
	AcceptVisitor(v RelationshipVisitor)
}

// RelationshipVisitor dispatches based on a relationship's type.
type RelationshipVisitor interface {
	VisitSeeAlso(SeeAlso)
}

// SeeAlso is a relationship that refers the user from one variable to another.
type SeeAlso struct {
	From, To Spec
}

// Specs returns the specs for the variables that are involved in the
// relationship.
func (r SeeAlso) Specs() []Spec {
	return []Spec{r.From, r.To}
}

// AcceptVisitor passes r to the appropriate method of v.
func (r SeeAlso) AcceptVisitor(v RelationshipVisitor) {
	v.VisitSeeAlso(r)
}
