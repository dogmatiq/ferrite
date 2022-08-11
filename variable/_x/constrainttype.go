package spec2

// Type limits a value to a specific Go type.
//
// This is a very permissive constraint as it allows any value that can be
// unmarshaled.
type Type interface {
	Constraint
}

type typeOf[T any] struct {
	Marshaler Marshaler[T]
}

func (c typeOf[T]) Test(RValue) error {
	return nil
}

func (c typeOf[T]) AcceptVisitor(v ConstraintVisitor) {
	v.VisitType(c)
}
