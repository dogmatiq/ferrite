package variable

import "github.com/dogmatiq/ferrite/internal/maybe"

// Relationship represents a relationship between two variables.
type Relationship interface {
	subject() Spec
	object() Spec
}

// EstablishRelationships establishes the given relationships.
func EstablishRelationships(relationships ...Relationship) {
	for _, rel := range relationships {
		sub := rel.subject()
		obj := rel.object()
		sub.addRelationship(rel)
		obj.addRelationship(rel)
	}
}

// Relationships returns a list of relationships where s is the subject.
func Relationships[T Relationship](s Spec) []T {
	var result []T

	for _, rel := range s.Relationships() {
		if rel.subject() == s {
			if rel, ok := rel.(T); ok {
				result = append(result, rel)
			}
		}
	}

	return result
}

// InverseRelationships returns a list of relationships where s is the object.
func InverseRelationships[T Relationship](s Spec) []T {
	var result []T

	for _, rel := range s.Relationships() {
		if rel.object() == s {
			if rel, ok := rel.(T); ok {
				result = append(result, rel)
			}
		}
	}

	return result
}

// Supersedes is a relationship type that indicates that a variable supersedes
// another (usually deprecated) variable.
type Supersedes struct {
	Subject, Supersedes Spec
}

func (r Supersedes) subject() Spec {
	return r.Subject
}

func (r Supersedes) object() Spec {
	return r.Supersedes
}

// RefersTo is a relationship type that indicates that a variable refers to
// another variable for information/documentation purposes.
type RefersTo struct {
	Subject, RefersTo Spec
}

func (r RefersTo) subject() Spec {
	return r.Subject
}

func (r RefersTo) object() Spec {
	return r.RefersTo
}

// DependsOn is a relationship type that indicates that a variable requires
// another variable to have a specific value in order be used.
type DependsOn struct {
	Subject, DependsOn Spec

	// Value is the value that the dependency must have in order for the subject
	// to be used. If it is absent the dependency must be any "truthy" value.
	Value maybe.Value[Literal]
}

func (r DependsOn) subject() Spec {
	return r.Subject
}

func (r DependsOn) object() Spec {
	return r.DependsOn
}
