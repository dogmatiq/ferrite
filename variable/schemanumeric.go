package variable

import (
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
	"golang.org/x/exp/constraints"
)

// Number is a schema that allows numeric values.
type Number interface {
	// Min is the minimum allowed value.
	Min() maybe.Value[Literal]

	// Max is the maximum allowed value.
	Max() maybe.Value[Literal]
}

// NumberAs is a numeric value depicted by type T.
type NumberAs[T constraints.Ordered] struct {
	marshaler Marshaler[T]
	min, max  maybe.Value[valueOf[T]]
}

// NewNumber returns a new numeric class.
func NewNumber[T constraints.Ordered, M Marshaler[T]](
	m M,
	min, max maybe.Value[T],
) (NumberAs[T], error) {
	vmin, err := marshal(m, min)
	if err != nil {
		return NumberAs[T]{}, err
	}

	vmax, err := marshal(m, max)
	if err != nil {
		return NumberAs[T]{}, err
	}

	return NumberAs[T]{
		marshaler: m,
		min:       vmin,
		max:       vmax,
	}, nil
}

// AcceptVisitor passes s to the appropriate method of v.
func (s NumberAs[T]) AcceptVisitor(v SchemaVisitor) {
	v.VisitNumeric(s)
}

// Min is the minimum allowed value.
func (s NumberAs[T]) Min() maybe.Value[Literal] {
	return maybe.Map(s.min, valueOf[T].Canonical)
}

// Max is the maximum allowed value.
func (s NumberAs[T]) Max() maybe.Value[Literal] {
	return maybe.Map(s.max, valueOf[T].Canonical)
}

// Marshal converts a value to its literal representation.
func (s NumberAs[T]) Marshal(v T) (Literal, error) {
	if err := s.validate(v); err != nil {
		return "", err
	}

	return s.marshaler.Marshal(v)
}

// Unmarshal converts a literal value to it's native representation.
func (s NumberAs[T]) Unmarshal(v Literal) (T, error) {
	n, err := s.marshaler.Unmarshal(v)
	if err != nil {
		return n, err
	}

	return n, s.validate(n)
}

// validate returns an error if v is invalid.
func (s NumberAs[T]) validate(v T) error {
	min, hasMin := s.min.Get()
	max, hasMax := s.max.Get()

	if (hasMin && v < min.native) || (hasMax && v > max.native) {
		return RangeError{
			Schema: s,
		}
	}

	return nil
}

// RangeError indicates that a numeric value was outside of the expected range.
type RangeError struct {
	Schema Number
}

func (e RangeError) Error() string {
	min, hasMin := e.Schema.Min().Get()
	max, hasMax := e.Schema.Max().Get()

	if !hasMax {
		return fmt.Sprintf("must be %s or greater", min)
	}

	if !hasMin {
		return fmt.Sprintf("must be %s or less", max)
	}

	return fmt.Sprintf("must be between %s and %s", min, max)
}
