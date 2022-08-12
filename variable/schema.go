package variable

// Schema describes the valid values of an environment value.
type Schema interface {
	AcceptVisitor(SchemaVisitor)
}

// SchemaVisitor dispatches based on a variable's schema.
type SchemaVisitor interface {
	VisitSet(Set)
}

// SchemaFor describes the valid values of an environment varible value depicted
// by type T.
type SchemaFor[T any] interface {
	Schema

	// Marshal converts a value to its literal representation.
	//
	// It panics if v is not a valid value according to this schema.
	Marshal(v T) Literal

	// Unmarshal converts a literal value to it's native representation.
	Unmarshal(n Name, v Literal) (T, ValidationError)
}
