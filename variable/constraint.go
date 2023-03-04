package variable

import "errors"

// A Constraint represents a constraint on the variable value in addition to the
// schema's requirements.
type Constraint interface {
	// Description returns a description of the constraint.
	Description() string
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

// WithConstraint is a SpecOption that adds a validator.
func WithConstraint[T any](
	desc string,
	fn func(T) ConstraintError,
) SpecOption[T] {
	return func(opts *specOptions[T]) error {
		if desc == "" {
			return errors.New("constraint description must not be empty")
		}

		v := constraint[T]{desc, fn}
		opts.Constraints = append(opts.Constraints, v)
		return nil
	}
}

// constraint is a function that implements the Validator interface.
type constraint[T any] struct {
	desc  string
	check func(T) ConstraintError
}

func (c constraint[T]) Description() string {
	return c.desc
}

func (c constraint[T]) Check(v T) ConstraintError {
	return c.check(v)
}
