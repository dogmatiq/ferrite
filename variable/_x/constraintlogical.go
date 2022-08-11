package spec2

// LogicalAnd requires a value to match a set of other constraints.
type LogicalAnd []Constraint

// Test checks if a value adheres to this constraint.
func (c LogicalAnd) Test(r RValue) error {
	if len(c) == 0 {
		panic("empty compound constraint")
	}

	var errors LogicalOrError

	for _, x := range c {
		err := x.Test(r)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

// AcceptVisitor calls the visitor method associated with this constraint type.
func (c LogicalAnd) AcceptVisitor(v ConstraintVisitor) {
	v.VisitLogicalAnd(c)
}

// LogicalAndError indicates that one or more sub-constraints of an LogicalAnd
// constraint have been violated.
type LogicalAndError []error

func (e LogicalAndError) Error() string {
	return e[0].Error()
}

// LogicalOr permits a value to match any one of a number of other constraints.
type LogicalOr []Constraint

// Test checks if a value adheres to this constraint.
func (c LogicalOr) Test(r RValue) error {
	if len(c) == 0 {
		panic("empty compound constraint")
	}

	var errors LogicalOrError

	for _, x := range c {
		err := x.Test(r)
		if err == nil {
			return nil
		}

		errors = append(errors, err)
	}

	return errors
}

// AcceptVisitor calls the visitor method associated with this constraint type.
func (c LogicalOr) AcceptVisitor(v ConstraintVisitor) {
	v.VisitLogicalOr(c)
}

// LogicalOrError indicates that all sub-constraints of a LogicalOr constraint
// have been violated.
type LogicalOrError []error

func (e LogicalOrError) Error() string {
	return e[0].Error()
}
