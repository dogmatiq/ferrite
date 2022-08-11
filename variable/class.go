package variable

// Class describes a broad category of environment variable values.
type Class interface {
	AcceptVisitor(ClassVisitor)
}

// ClassVisitor dispatches based on a variable's class.
type ClassVisitor interface {
	VisitSet(Set)
}

// ClassOf describes a broad category of environment variable values depicted by
// type T.
type ClassOf[T any] interface {
	Class

	Marshal(v T) Literal
	Unmarshal(n Name, v Literal) (native T, canonical Literal, err ValidationError)
}
