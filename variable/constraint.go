package variable

import "errors"

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

// ConstraintFunc is a function that implements the Constraint interface.
type ConstraintFunc[T any] struct {
	Desc string
	User bool
	Fn   func(T) ConstraintError
}

// Description returns a description of the constraint.
func (c ConstraintFunc[T]) Description() string {
	return c.Desc
}

// IsUserDefined returns true if this constraint was defined by the user.
func (c ConstraintFunc[T]) IsUserDefined() bool {
	return c.User
}

// Check returns an error if v does not satisfy the constraint.
func (c ConstraintFunc[T]) Check(v T) ConstraintError {
	return c.Fn(v)
}

// WithConstraint is a TypedSpecOption that adds a constraint to a variable.
// Deprecated: xxx
func WithConstraint[T any](
	desc string,
	fn func(T) ConstraintError,
) TypedSpecOption[T] {
	return withConstraint(
		ConstraintFunc[T]{
			desc,
			false,
			fn,
		},
	)
}

// WithUserConstraint is a TypedSpecOption that adds a user-defined constraint.
// Deprecated: xxx
func WithUserConstraint[T any](
	desc string,
	fn func(T) error,
) TypedSpecOption[T] {
	return withConstraint(
		ConstraintFunc[T]{
			desc,
			true,
			func(v T) ConstraintError {
				return fn(v)
			},
		},
	)
}

// Deprecated: xxx
func withConstraint[T any](c ConstraintFunc[T]) TypedSpecOption[T] {
	return func(opts *specOptions[T]) error {
		if c.Desc == "" {
			return errors.New("constraint description must not be empty")
		}

		opts.Constraints = append(opts.Constraints, c)
		return nil
	}
}
