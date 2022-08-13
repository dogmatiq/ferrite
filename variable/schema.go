package variable

// Schema describes the valid values of an environment value.
type Schema interface {
	// AcceptVisitor passes the schema to the appropriate method of v.
	AcceptVisitor(SchemaVisitor)
}

// SchemaError indicates that a value is invalid because it violates its schema.
type SchemaError interface {
	error

	// Schema returns the schema that was violated.
	Schema() Schema

	// AcceptVisitor passes the error to the appropriate method of v.
	AcceptVisitor(SchemaErrorVisitor)
}

// SchemaVisitor dispatches based on a variable's schema.
type SchemaVisitor interface {
	VisitSet(Set)
	VisitNumeric(Numeric)
}

// SchemaErrorVisitor dispatches based on the type of a SchemaError.
type SchemaErrorVisitor interface {
	VisitSetMembershipError(SetMembershipError)
	VisitMinError(MinError)
	VisitMaxError(MaxError)
}

// TypedSchema describes the valid values of an environment varible value
// depicted by type T.
type TypedSchema[T any] interface {
	Schema
	Marshaler[T]
}
