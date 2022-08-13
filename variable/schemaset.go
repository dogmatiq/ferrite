package variable

import (
	"errors"
	"fmt"
)

// Set is a schema that only allows a specific set of static values.
type Set interface {
	Schema

	// Members returns the members of the set.
	Members() []Literal
}

// TypedSet is a Set containing values of type T.
type TypedSet[T any] struct {
	render  func(T) Literal
	order   []Literal
	members map[Literal]T
}

// NewSet contains a new set containing the given values.
func NewSet[T any](
	members []T,
	render func(T) Literal,
) (TypedSet[T], error) {
	if len(members) < 2 {
		return TypedSet[T]{}, errors.New("must allow at least two distinct values")
	}

	s := TypedSet[T]{
		render:  render,
		members: make(map[Literal]T, len(members)),
	}

	for _, v := range members {
		lit := render(v)

		if lit == "" {
			return TypedSet[T]{}, errors.New("literals can not be an empty string")
		}

		if _, ok := s.members[lit]; ok {
			return TypedSet[T]{}, fmt.Errorf("literals must be unique but multiple values are represented as %q", string(lit))
		}

		s.order = append(s.order, lit)
		s.members[lit] = v
	}

	return s, nil
}

// Members returns the members of the set.
func (s TypedSet[T]) Members() []Literal {
	return s.order
}

// AcceptVisitor passes s to the appropriate method of v.
func (s TypedSet[T]) AcceptVisitor(v SchemaVisitor) {
	v.VisitSet(s)
}

// Marshal converts a value to its literal representation.
func (s TypedSet[T]) Marshal(v T) (Literal, error) {
	lit := s.render(v)
	if _, ok := s.members[lit]; ok {
		return lit, nil
	}

	return "", SetMembershipError{s}
}

// Unmarshal converts a literal value to it's native representation.
func (s TypedSet[T]) Unmarshal(v Literal) (T, error) {
	if n, ok := s.members[v]; ok {
		return n, nil
	}

	var zero T
	return zero, SetMembershipError{s}
}

// SetMembershipError is a validation error that indicates a value is not a
// member of a specific set.
type SetMembershipError struct {
	Set Set
}

var _ SchemaError = SetMembershipError{}

// Schema returns the schema that was violated.
func (e SetMembershipError) Schema() Schema {
	return e.Set
}

// AcceptVisitor passes the error to the appropriate method of v.
func (e SetMembershipError) AcceptVisitor(v SchemaErrorVisitor) {
	v.VisitSetMembershipError(e)
}

func (e SetMembershipError) Error() string {
	members := e.Set.Members()

	switch n := len(members); n {
	case 2:
		return fmt.Sprintf(
			"expected either %s or %s",
			members[0],
			members[1],
		)
	case 3:
		return fmt.Sprintf(
			"expected %s, %s or %s",
			members[0],
			members[1],
			members[2],
		)
	case 4:
		return fmt.Sprintf(
			"expected %s, %s, %s or %s",
			members[0],
			members[1],
			members[2],
			members[3],
		)
	default: // 5 or more
		return fmt.Sprintf(
			"expected %s, %s ... %s, or one of %d other values",
			members[0],
			members[1],
			members[n-1],
			n-3,
		)
	}
}
