package variable

// Relationship represents a relationship between two variables.
type Relationship interface {
	subject() Spec
	object() Spec
}

// AddRelationship adds the relationship to the subject and related variable
// specifications.
func AddRelationship(rel Relationship) error {
	sub := rel.subject()
	obj := rel.object()

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
