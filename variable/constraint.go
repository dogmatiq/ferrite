package variable

// A Constraint represents a constraint on the variable value in addition to the
// schema's requirements.
type Constraint interface {
	// Description returns a description of the constraint.
	Description() string

	// IsUserDefined returns true if this constraint was defined by the user.
	IsUserDefined() bool
}

// TypedConstraint places a constraint on the variable value in addition to the
// schema's requirements.
type TypedConstraint[T any] interface {
	Constraint

	// Check returns an error if v does not satisfy the constraint.
	Check(T) ConstraintError
}

// ConstraintError indicates that a value does not satisfy a constraint.
type ConstraintError interface {
	error
}

// constraint is a function that implements the Constraint interface.
type constraint[T any] struct {
	desc  string
	user  bool
	check func(T) ConstraintError
}

// Description returns a description of the constraint.
func (c constraint[T]) Description() string {
	return c.desc
}

// IsUserDefined returns true if this constraint was defined by the user.
func (c constraint[T]) IsUserDefined() bool {
	return c.user
}

// Check returns an error if v does not satisfy the constraint.
func (c constraint[T]) Check(v T) ConstraintError {
	return c.check(v)
}
