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

// WithConstraint is a SpecOption that adds a constraint to a variable.
func WithConstraint[T any](
	desc string,
	fn func(T) ConstraintError,
) SpecOption[T] {
	return withConstraint(
		constraint[T]{
			desc,
			false,
			fn,
		},
	)
}

// WithUserConstraint is a SpecOption that adds a user-defined constraint.
func WithUserConstraint[T any](
	desc string,
	fn func(T) error,
) SpecOption[T] {
	return withConstraint(
		constraint[T]{
			desc,
			true,
			func(v T) ConstraintError {
				return fn(v)
			},
		},
	)
}

func withConstraint[T any](c constraint[T]) SpecOption[T] {
	return func(opts *specOptions[T]) error {
		if c.desc == "" {
			return errors.New("constraint description must not be empty")
		}

		opts.Constraints = append(opts.Constraints, c)
		return nil
	}
}

// constraint is a function that implements the Constraint interface.
type constraint[T any] struct {
	desc  string
	user  bool
	check func(T) ConstraintError
}

func (c constraint[T]) Description() string {
	return c.desc
}

func (c constraint[T]) IsUserDefined() bool {
	return c.user
}

func (c constraint[T]) Check(v T) ConstraintError {
	return c.check(v)
}
