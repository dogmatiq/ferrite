package schema

// Range is a schema that allows values within a specific range.
type Range struct {
	Min string
	Max string
}

// AcceptVisitor calls the method on v that is related to this schema type.
func (s Range) AcceptVisitor(v Visitor) {
	v.VisitRange(s)
}
