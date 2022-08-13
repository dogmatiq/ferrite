package variable

import (
	"fmt"

	"github.com/dogmatiq/ferrite/maybe"
	"golang.org/x/exp/constraints"
)

// Numeric is a schema that allows numeric values.
type Numeric interface {
	Schema

	// Min is the minimum allowed value.
	Min() maybe.Value[Literal]

	// Max is the maximum allowed value.
	Max() maybe.Value[Literal]
}

// TypedNumber is a numeric value depicted by type T.
type TypedNumber[T constraints.Ordered] struct {
	marshaler Marshaler[T]
	min, max  maybe.Value[valueOf[T]]
}

// NewNumber returns a new numeric class.
func NewNumber[T constraints.Ordered, M Marshaler[T]](
	m M,
	min, max maybe.Value[T],
) (TypedNumber[T], error) {
	vmin, err := marshal(m, min)
	if err != nil {
		return TypedNumber[T]{}, err
	}

	vmax, err := marshal(m, max)
	if err != nil {
		return TypedNumber[T]{}, err
	}

	return TypedNumber[T]{
		marshaler: m,
		min:       vmin,
		max:       vmax,
	}, nil
}

// AcceptVisitor passes s to the appropriate method of v.
func (s TypedNumber[T]) AcceptVisitor(v SchemaVisitor) {
	v.VisitNumeric(s)
}

// Min is the minimum allowed value.
func (s TypedNumber[T]) Min() maybe.Value[Literal] {
	return maybe.Map(s.min, valueOf[T].Canonical)
}

// Max is the maximum allowed value.
func (s TypedNumber[T]) Max() maybe.Value[Literal] {
	return maybe.Map(s.max, valueOf[T].Canonical)
}

// Marshal converts a value to its literal representation.
func (s TypedNumber[T]) Marshal(v T) (Literal, error) {
	if err := s.validate(v); err != nil {
		return "", err
	}

	return s.marshaler.Marshal(v)
}

// Unmarshal converts a literal value to it's native representation.
func (s TypedNumber[T]) Unmarshal(v Literal) (T, error) {
	n, err := s.marshaler.Unmarshal(v)
	if err != nil {
		return n, err
	}

	return n, s.validate(n)
}

// validate returns an error if v is invalid.
func (s TypedNumber[T]) validate(v T) error {
	min, hasMin := s.min.Get()
	max, hasMax := s.max.Get()

	if hasMin && v < min.native {
		return MinError{s}
	}

	if hasMax && v > max.native {
		return MaxError{s}
	}

	return nil
}

// MinError indicates that a numeric value was less than the minimum allowed
// value.
type MinError struct {
	Numeric Numeric
}

var _ SchemaError = MinError{}

// Schema returns the schema that was violated.
func (e MinError) Schema() Schema {
	return e.Numeric
}

// AcceptVisitor passes the error to the appropriate method of v.
func (e MinError) AcceptVisitor(v SchemaErrorVisitor) {
	v.VisitMinError(e)
}

func (e MinError) Error() string {
	min := e.Numeric.Min().MustGet()

	if max, ok := e.Numeric.Max().Get(); ok {
		return fmt.Sprintf("too low, expected between %s and %s", min, max)
	}

	return fmt.Sprintf("too low, expected %s or greater", min)
}

// MaxError indicates that a numeric value was greater than the maximum
// allowed value.
type MaxError struct {
	Numeric Numeric
}

var _ SchemaError = MaxError{}

// Schema returns the schema that was violated.
func (e MaxError) Schema() Schema {
	return e.Numeric
}

// AcceptVisitor passes the error to the appropriate method of v.
func (e MaxError) AcceptVisitor(v SchemaErrorVisitor) {
	v.VisitMaxError(e)
}

func (e MaxError) Error() string {
	max := e.Numeric.Max().MustGet()

	if min, ok := e.Numeric.Min().Get(); ok {
		return fmt.Sprintf("too high, expected between %s and %s", min, max)
	}

	return fmt.Sprintf("too high, expected %s or less", max)
}
