package variable

import (
	"errors"
	"fmt"
)

// Set is a schema that only allows a specific set of static values.
type Set interface {
	// Members returns the members of the set.
	Members() []Literal
}

// SetOf is a Set containing values of type T.
type SetOf[T any] struct {
	marshal func(T) (Literal, error)
	members []Literal
	values  map[Literal]T
}

// NewSet contains a new set containing the given values.
func NewSet[T any](
	marshal func(T) (Literal, error),
	values ...T,
) (SetOf[T], error) {
	if len(values) == 0 {
		return SetOf[T]{}, errors.New("must have at least one member")
	}

	s := SetOf[T]{
		marshal: marshal,
		values:  make(map[Literal]T, len(values)),
	}

	for _, v := range values {
		lit, err := marshal(v)
		if err != nil {
			return SetOf[T]{}, err
		}

		if lit == "" {
			return SetOf[T]{}, errors.New("members must not have empty string representations")
		}

		if _, ok := s.values[lit]; ok {
			return SetOf[T]{}, errors.New("members literals must be unique")
		}

		s.members = append(s.members, lit)
		s.values[lit] = v
	}

	return s, nil
}

// Members returns the members of the set.
func (s SetOf[T]) Members() []Literal {
	return s.members
}

// AcceptVisitor passes s to the appropriate method of v.
func (s SetOf[T]) AcceptVisitor(v SchemaVisitor) {
	v.VisitSet(s)
}

// Marshal converts a value to its literal representation.
func (s SetOf[T]) Marshal(v T) (Literal, error) {
	lit, err := s.marshal(v)
	if err != nil {
		return "", err
	}

	if _, ok := s.values[lit]; ok {
		return lit, nil
	}

	return "", SetMembershipError{
		Schema: s,
	}
}

// Unmarshal converts a literal value to it's native representation.
func (s SetOf[T]) Unmarshal(v Literal) (T, error) {
	if n, ok := s.values[v]; ok {
		return n, nil
	}

	var zero T
	return zero, SetMembershipError{
		Schema: s,
	}
}

// SetMembershipError is a validation error that indicates a value is not a
// member of a specific set.
type SetMembershipError struct {
	Schema Set
}

func (e SetMembershipError) Error() string {
	members := e.Schema.Members()

	if len(members) == 2 {
		return fmt.Sprintf("must be either %s or %s", members[0], members[1])
	}

	return "must be a member of the set"
}
