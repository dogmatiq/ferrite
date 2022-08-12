package variable

// Schema describes the valid values of an environment value.
type Schema interface {
	AcceptVisitor(SchemaVisitor)
}

// SchemaVisitor dispatches based on a variable's schema.
type SchemaVisitor interface {
	VisitSet(Set)
	VisitNumeric(Number)
}

// SchemaFor describes the valid values of an environment varible value depicted
// by type T.
type SchemaFor[T any] interface {
	Schema
	Marshaler[T]
}
