package spec2

// Set limits values to a concrete set of accepted values.
type Set interface {
	Constraint
}

// SetMember is a member of a Set constraint.
type SetMember interface {
	// String returns the string representation of the value.
	String() string

	// Description returns a human-readable description of the value.
	Description() string
}

type setOf[T any] struct {
	m map[string]setMember[T]
}

func (c setOf[T]) Test(RValue) error {
	panic("not implemented")
}

func (c setOf[T]) AcceptVisitor(v ConstraintVisitor) {
	v.VisitSet(c)
}

type setMember[T any] struct {
	K    string
	V    T
	Desc string
}

func (m setMember[T]) String() string {
	return m.K
}

func (m setMember[T]) Description() string {
	return m.Desc
}
