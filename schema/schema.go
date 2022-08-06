package schema

// Schema describes the valid values of an environment variable.
type Schema interface {
	// AcceptVisitor calls the method on v that is related to this schema type.
	AcceptVisitor(v Visitor)
}

// Visitor dispatches logic based on a schema element.
type Visitor interface {
	VisitOneOf(OneOf)
	VisitLiteral(Literal)
	VisitType(TypeSchema)
	VisitRange(Range)
}
