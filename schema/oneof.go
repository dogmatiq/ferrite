package schema

// OneOf is a schema that allows any of a number of other schemas.
type OneOf []Schema

// AcceptVisitor calls the method on v that is related to this schema type.
func (s OneOf) AcceptVisitor(v Visitor) {
	v.VisitOneOf(s)
}
