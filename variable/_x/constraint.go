package spec2

// A Constraint places some requirement upon an environment variable's value.
type Constraint interface {
	// Test checks if a value adheres to this constraint.
	Test(RValue) error

	// AcceptVisitor calls the visitor method associated with this constraint
	// type.
	AcceptVisitor(ConstraintVisitor)
}

// ConstraintVisitor dispatches based on the type of a constraint.
type ConstraintVisitor interface {
	VisitLogicalAnd(LogicalAnd)
	VisitLogicalOr(LogicalOr)
	VisitType(Type)
	VisitRange(Range)
	VisitSet(Set)
}
