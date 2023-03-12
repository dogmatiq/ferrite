package variable

// Relationship represents a relationship between two variables.
type Relationship interface {
	// Subject returns the spec of the variable that is the subject of the
	// relationship.
	Subject() Spec

	// Related returns the spec of the related variable.
	Related() Spec

	// Inverse returns the inverse of the relationship.
	Inverse() Relationship

	// AcceptVisitor passes the schema to the appropriate method of v.
	AcceptVisitor(v RelationshipVisitor)
}

// ApplyRelationship adds the relationship to the subject and related
// variable specifications.
func ApplyRelationship(rel Relationship) {
	rel.Subject().addRelationship(rel)
	rel.Related().addRelationship(rel.Inverse())
}

// RelationshipVisitor dispatches based on a relationship's type.
type RelationshipVisitor interface {
	VisitRefersTo(RefersTo)
	VisitIsReferredToBy(IsReferredToBy)
	VisitSupersedes(Supersedes)
	VisitIsSupersededBy(IsSupersededBy)
}

// FilterRelationships returns the relationships of s that are of type T.
func FilterRelationships[T Relationship](s Spec) []T {
	var result []T

	for _, rel := range s.Relationships() {
		if rel, ok := rel.(T); ok {
			result = append(result, rel)
		}
	}

	return result
}
