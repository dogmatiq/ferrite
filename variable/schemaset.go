package variable

import (
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

// Set is a schema that only allows a specific set of static values.
type Set interface {
	Schema

	// Literals returns the members of the set as literals.
	Literals() []Literal
}

// TypedSet is a Set containing values of type T.
type TypedSet[T any] struct {
	Members   []SetMember[T]
	ToLiteral func(T) Literal
}

// SetMember is a member of a TypedSet.
type SetMember[T any] struct {
	Value       T
	Description string
}

// Literals returns the members of the set as literals.
func (s TypedSet[T]) Literals() []Literal {
	literals := make([]Literal, len(s.Members))
	for i, m := range s.Members {
		literals[i] = s.ToLiteral(m.Value)
	}
	return literals
}

// Type returns the type of the native value.
func (s TypedSet[T]) Type() reflect.Type {
	return typeOf[T]()
}

// Finalize prepares the schema for use.
//
// It returns an error if schema is invalid.
func (s TypedSet[T]) Finalize() error {
	if len(s.Members) < 2 {
		return errors.New("must allow at least two distinct values")
	}

	var uniq []Literal

	for _, v := range s.Members {
		lit := s.ToLiteral(v.Value)

		if lit.String == "" {
			return errors.New("literals can not be an empty string")
		}

		if slices.Contains(uniq, lit) {
			return fmt.Errorf(
				"literals must be unique but multiple values are represented as %q",
				lit.String,
			)
		}

		uniq = append(uniq, lit)
	}

	return nil
}

// AcceptVisitor passes s to the appropriate method of v.
func (s TypedSet[T]) AcceptVisitor(v SchemaVisitor) {
	v.VisitSet(s)
}

// Marshal converts a value to its literal representation.
func (s TypedSet[T]) Marshal(v T) (Literal, error) {
	lit := s.ToLiteral(v)

	for _, v := range s.Members {
		if lit == s.ToLiteral(v.Value) {
			return lit, nil
		}
	}

	return Literal{}, SetMembershipError{s}
}

// Unmarshal converts a literal value to it's native representation.
func (s TypedSet[T]) Unmarshal(v Literal) (T, error) {
	for _, m := range s.Members {
		if v == s.ToLiteral(m.Value) {
			return m.Value, nil
		}
	}

	var zero T
	return zero, SetMembershipError{s}
}

// Examples returns a (possibly empty) set of examples of valid values.
func (s TypedSet[T]) Examples(hasOtherExamples bool) []TypedExample[T] {
	var examples []TypedExample[T]

	for _, m := range s.Members {
		examples = append(examples, TypedExample[T]{
			Native:      m.Value,
			Description: m.Description,
		})
	}

	return examples
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
	members := e.Set.Literals()

	switch n := len(members); n {
	case 2:
		return fmt.Sprintf(
			"expected either %s or %s",
			members[0].Quote(),
			members[1].Quote(),
		)
	case 3:
		return fmt.Sprintf(
			"expected %s, %s or %s",
			members[0].Quote(),
			members[1].Quote(),
			members[2].Quote(),
		)
	case 4:
		return fmt.Sprintf(
			"expected %s, %s, %s or %s",
			members[0].Quote(),
			members[1].Quote(),
			members[2].Quote(),
			members[3].Quote(),
		)
	default: // 5 or more
		return fmt.Sprintf(
			"expected %s, %s ... %s, or one of %d other values",
			members[0].Quote(),
			members[1].Quote(),
			members[n-1].Quote(),
			n-3,
		)
	}
}
