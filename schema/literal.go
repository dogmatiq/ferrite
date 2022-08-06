package schema

// Literal is a Schema that allows a specific string value.
type Literal string

// AcceptVisitor calls the method on v that is related to this schema type.
func (s Literal) AcceptVisitor(v Visitor) {
	v.VisitLiteral(s)
}
