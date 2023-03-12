package variable

import "errors"

// Relationship represents a relationship between two variables.
type Relationship interface {
	// Subject returns the spec of the variable that is the subject of the
	// relationship.
	Subject() Spec

	// Related returns the spec of the related variable.
	Related() Spec

	// Inverse returns the inverse of the relationship, such that the related
	// variable becomes the subject.
	Inverse() Relationship

	// ConflictsWith returns true if this relationship conflicts with r.
	ConflictsWith(c Relationship) bool
}

// ApplyRelationship adds the relationship to the subject and related variable
// specifications.
func ApplyRelationship(r Relationship) error {
	sub := r.Subject()
	rel := r.Related()
	inv := r.Inverse()

	for _, x := range sub.Relationships() {
		if x.ConflictsWith(r) {
			return errors.New("<error>")
		}
	}

	for _, x := range rel.Relationships() {
		if x.ConflictsWith(r) {
			return errors.New("<error>")
		}
	}

	sub.addRelationship(r)
	rel.addRelationship(inv)

	return nil
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
