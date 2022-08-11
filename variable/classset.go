package variable

import (
	"errors"
	"fmt"
)

// Set is a class that only allows a specific set of static values.
type Set interface {
	// Members returns the members of the set.
	Members() []String
}

// SetOf is a Set containing values of type T.
type SetOf[T any] struct {
	ordered   []String
	marshal   func(T) String
	unmarshal map[String]T
}

// NewSet contains a new set containing the given values.
func NewSet[T any](
	marshal func(T) String,
	values ...T,
) (SetOf[T], error) {
	if len(values) == 0 {
		return SetOf[T]{}, errors.New("must have at least one member")
	}

	c := SetOf[T]{
		marshal:   marshal,
		unmarshal: make(map[String]T, len(values)),
	}

	for _, v := range values {
		s := marshal(v)
		if s == "" {
			return SetOf[T]{}, errors.New("members must not have empty string representations")
		}

		c.ordered = append(c.ordered, s)
		c.unmarshal[s] = v
	}

	return c, nil
}

// Members returns the members of the set.
func (c SetOf[T]) Members() []String {
	return c.ordered
}

// AcceptVisitor passes c to the appropriate method of v.
func (c SetOf[T]) AcceptVisitor(v ClassVisitor) {
	v.VisitSet(c)
}

// Marshal marshals v to its string representation.
func (c SetOf[T]) Marshal(v T) String {
	s := c.marshal(v)

	if _, ok := c.unmarshal[s]; !ok {
		panic(fmt.Sprintf("cannot marshal non-member (%s)", s))
	}

	return s
}

// Unmarshal unmarshals a string representation of a value.
//
// It returns the native value and the canonical string representation.
func (c SetOf[T]) Unmarshal(n Name, s String) (T, String, ValidationError) {
	if v, ok := c.unmarshal[s]; ok {
		return v, s, nil
	}

	var zero T
	return zero, "", SetError{
		name:     n,
		verbatim: s,
		set:      c,
	}
}

// SetError is a validation error that indicates a value is not a member of a
// specific set.
type SetError struct {
	name     Name
	verbatim String
	set      Set
}

// Name returns the name of the environment variable.
func (e SetError) Name() Name {
	return e.name
}

// Verbatim returns the offending value.
func (e SetError) Verbatim() String {
	return e.verbatim
}

// Reason returns a human-readable explanation of why the value is invalid.
func (e SetError) Reason() string {
	members := e.set.Members()

	if len(members) == 2 {
		return fmt.Sprintf("must be either %s or %s", members[0], members[1])
	}

	return "must be a member of the set"
}

func (e SetError) Error() string {
	return formatValidationError(e)
}
