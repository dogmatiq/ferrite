package variable

import "errors"

// Relationship represents a relationship between two variables.
type Relationship interface {
	subject() Spec
	object() Spec
	conflictsWith(rel Relationship) bool
}

// AddRelationship adds the relationship to the subject and related variable
// specifications.
func AddRelationship(rel Relationship) error {
	sub := rel.subject()
	obj := rel.object()

	for _, x := range sub.Relationships() {
		if x.conflictsWith(rel) || rel.conflictsWith(x) {
			return errors.New("<error>")
		}
	}

	for _, x := range obj.Relationships() {
		if x.conflictsWith(rel) || rel.conflictsWith(x) {
			return errors.New("<error>")
		}
	}

	sub.addRelationship(rel)
	obj.addRelationship(rel)

	return nil
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

func (r Supersedes) conflictsWith(rel Relationship) bool {
	if rel, ok := rel.(Supersedes); ok {
		return rel.Subject == r.Supersedes &&
			rel.Supersedes == r.Subject
	}
	return false
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

func (r RefersTo) conflictsWith(rel Relationship) bool {
	return false
}
